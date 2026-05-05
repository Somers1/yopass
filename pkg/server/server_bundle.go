package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhaals/yopass/pkg/yopass"
	"go.uber.org/zap"
)

const bundleKeyPrefix = "bundle:"

type createBundleRequest struct {
	FileKeys   []string `json:"file_keys"`
	Filenames  []string `json:"filenames"`
	Sizes      []int64  `json:"sizes"`
	Expiration int32    `json:"expiration"`
	OneTime    bool     `json:"one_time"`
}

func (y *Server) createBundle(w http.ResponseWriter, r *http.Request) {
	var req createBundleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		y.Logger.Debug("Unable to decode bundle request", zap.Error(err))
		http.Error(w, `{"message": "Unable to parse json"}`, http.StatusBadRequest)
		return
	}

	if len(req.FileKeys) == 0 {
		http.Error(w, `{"message": "file_keys must not be empty"}`, http.StatusBadRequest)
		return
	}

	if len(req.FileKeys) != len(req.Filenames) || len(req.FileKeys) != len(req.Sizes) {
		http.Error(w, `{"message": "file_keys, filenames, and sizes must have the same length"}`, http.StatusBadRequest)
		return
	}

	if !validExpiration(req.Expiration) {
		http.Error(w, `{"message": "Invalid expiration specified"}`, http.StatusBadRequest)
		return
	}

	if !req.OneTime && y.ForceOneTimeSecrets {
		http.Error(w, `{"message": "Secret must be one time download"}`, http.StatusBadRequest)
		return
	}

	for _, fk := range req.FileKeys {
		if _, err := y.DB.Status(streamKeyPrefix + fk); err != nil {
			y.Logger.Debug("Bundle references non-existent file", zap.String("key", fk), zap.Error(err))
			http.Error(w, `{"message": "Referenced file not found: `+fk+`"}`, http.StatusBadRequest)
			return
		}
	}

	bundle := yopass.Bundle{
		Expiration: req.Expiration,
		OneTime:    req.OneTime,
	}
	for i, fk := range req.FileKeys {
		bundle.Files = append(bundle.Files, yopass.BundleFile{
			Key:      fk,
			Filename: req.Filenames[i],
			Size:     req.Sizes[i],
		})
	}

	manifest, err := json.Marshal(bundle)
	if err != nil {
		y.Logger.Error("Failed to marshal bundle", zap.Error(err))
		http.Error(w, `{"message": "Failed to create bundle"}`, http.StatusInternalServerError)
		return
	}

	key, err := yopass.GenerateID()
	if err != nil {
		y.Logger.Error("Unable to generate bundle ID", zap.Error(err))
		http.Error(w, `{"message": "Unable to generate ID"}`, http.StatusInternalServerError)
		return
	}

	secret := yopass.Secret{
		Expiration: req.Expiration,
		Message:    string(manifest),
		OneTime:    req.OneTime,
	}
	if err := y.DB.Put(bundleKeyPrefix+key, secret); err != nil {
		y.Logger.Error("Failed to store bundle", zap.Error(err))
		http.Error(w, `{"message": "Failed to store bundle"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": key}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		y.Logger.Error("Failed to marshal bundle response", zap.Error(err))
	}
	if _, err := w.Write(jsonData); err != nil {
		y.Logger.Error("Failed to write bundle response", zap.Error(err))
	}
}

func (y *Server) getBundle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-cache")
	w.Header().Set("Content-Type", "application/json")

	key := mux.Vars(r)["key"]

	secret, err := y.DB.Get(bundleKeyPrefix + key)
	if err != nil {
		y.Logger.Debug("Bundle not found", zap.Error(err))
		http.Error(w, `{"message": "Bundle not found"}`, http.StatusNotFound)
		return
	}

	var bundle yopass.Bundle
	if err := json.Unmarshal([]byte(secret.Message), &bundle); err != nil {
		y.Logger.Error("Failed to decode bundle manifest", zap.Error(err))
		http.Error(w, `{"message": "Invalid bundle data"}`, http.StatusInternalServerError)
		return
	}

	type bundleResponse struct {
		Files      []yopass.BundleFile `json:"files"`
		OneTime    bool                `json:"one_time"`
		Expiration int32               `json:"expiration"`
	}
	resp := bundleResponse{
		Files:      bundle.Files,
		OneTime:    bundle.OneTime,
		Expiration: bundle.Expiration,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		y.Logger.Error("Failed to write bundle response", zap.Error(err))
		return
	}

	if bundle.OneTime {
		y.deleteBundleManifest(key)
	}
}

func (y *Server) getBundleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-cache")
	w.Header().Set("Content-Type", "application/json")

	key := mux.Vars(r)["key"]
	oneTime, err := y.DB.Status(bundleKeyPrefix + key)
	if err != nil {
		y.Logger.Debug("Bundle not found", zap.Error(err))
		http.Error(w, `{"message": "Bundle not found"}`, http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]bool{"oneTime": oneTime}); err != nil {
		y.Logger.Error("Failed to write bundle status response", zap.Error(err))
	}
}

func (y *Server) deleteBundle(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	secret, err := y.DB.Get(bundleKeyPrefix + key)
	if err != nil {
		http.Error(w, `{"message": "Bundle not found"}`, http.StatusNotFound)
		return
	}

	var bundle yopass.Bundle
	if err := json.Unmarshal([]byte(secret.Message), &bundle); err != nil {
		y.Logger.Error("Failed to decode bundle manifest for delete", zap.Error(err))
		http.Error(w, `{"message": "Invalid bundle data"}`, http.StatusInternalServerError)
		return
	}

	y.deleteBundleAndFiles(key, bundle)
	w.WriteHeader(http.StatusNoContent)
}

func (y *Server) deleteBundleAndFiles(key string, bundle yopass.Bundle) {
	ctx := context.Background()
	for _, f := range bundle.Files {
		if err := y.FileStore.Delete(ctx, f.Key); err != nil {
			y.Logger.Error("Failed to delete bundle file from store",
				zap.String("bundle", key), zap.String("file", f.Key), zap.Error(err))
		}
		if _, err := y.DB.Delete(streamKeyPrefix + f.Key); err != nil {
			y.Logger.Error("Failed to delete bundle file metadata",
				zap.String("bundle", key), zap.String("file", f.Key), zap.Error(err))
		}
	}
	y.deleteBundleManifest(key)
}

func (y *Server) deleteBundleManifest(key string) {
	if _, err := y.DB.Delete(bundleKeyPrefix + key); err != nil {
		y.Logger.Error("Failed to delete bundle manifest",
			zap.String("bundle", key), zap.Error(err))
	}
}
