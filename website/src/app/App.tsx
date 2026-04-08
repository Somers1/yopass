import FeaturesSection from '@shared/components/FeaturesSection';
import CreateSecret from '@features/CreateSecret';
import { Routes, Route, HashRouter } from 'react-router-dom';
import { useConfig } from '@shared/hooks/useConfig';
import Navbar from '@shared/components/Navbar';
import Prefetcher from '@features/display-secret/Prefetcher';
import StreamingUpload from '@features/StreamingUpload';
import ReadOnlyLanding from '@features/ReadOnlyLanding';
import { useTranslation } from 'react-i18next';

export default function App() {
  const { DISABLE_UPLOAD, READ_ONLY, PRIVACY_NOTICE_URL, IMPRINT_URL } =
    useConfig();
  const { t } = useTranslation();
  return (
    <div className="min-h-screen flex flex-col overflow-x-hidden" style={{ background: 'radial-gradient(ellipse, #f5f3f0 50%, #ececec 100%)' }}>
      <HashRouter>
        <Navbar />

        {/* Main Content */}
        <div className="w-full max-w-3xl mx-auto mb-auto px-4 py-12 sm:py-16">
          <div className="card bg-base-100 shadow-sm border border-base-300">
            <div className="card-body p-6">
              <Routes>
                <Route
                  path="/"
                  element={READ_ONLY ? <ReadOnlyLanding /> : <CreateSecret />}
                />
                {READ_ONLY ? (
                  <Route path="/upload" element={<ReadOnlyLanding />} />
                ) : (
                  !DISABLE_UPLOAD && (
                    <Route path="/upload" element={<StreamingUpload />} />
                  )
                )}
                <Route
                  path="/:format/:key/:password"
                  element={<Prefetcher />}
                />
                <Route path="/:format/:key" element={<Prefetcher />} />
              </Routes>
            </div>
          </div>
          <FeaturesSection />
        </div>
      </HashRouter>
      {/* Footer */}
      <footer className="border-base-300">
        <div className="container mx-auto px-4 py-4">
          <div className="flex flex-wrap items-center justify-center gap-2 text-xs text-base-content/40">
            {PRIVACY_NOTICE_URL && PRIVACY_NOTICE_URL.trim() && (
              <>
                <a href={PRIVACY_NOTICE_URL} className="hover:text-base-content/60 transition-colors" target="_blank" rel="noopener noreferrer">
                  {t('footer.privacyNotice')}
                </a>
                <span>•</span>
              </>
            )}
            {IMPRINT_URL && IMPRINT_URL.trim() && (
              <>
                <a href={IMPRINT_URL} className="hover:text-base-content/60 transition-colors" target="_blank" rel="noopener noreferrer">
                  {t('footer.imprint')}
                </a>
                <span>•</span>
              </>
            )}
            <span>
              Forked from open source{' '}
              <a href="https://github.com/jhaals/yopass" className="hover:text-base-content/60 transition-colors" target="_blank" rel="noopener noreferrer">
                Yopass
              </a>
              {' '}by Johan Haals
            </span>
          </div>
        </div>
      </footer>
    </div>
  );
}
