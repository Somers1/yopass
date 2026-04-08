import { useEffect } from 'react';
import {
  getInitialLogicalTheme,
  logicalToDaisyTheme,
  THEME_STORAGE_KEY,
} from '../theme/theme';
import { useConfig } from '../hooks/useConfig';
import { useTranslation } from 'react-i18next';
import { useLocation } from 'react-router-dom';

export default function Navbar() {
  const { DISABLE_UPLOAD, READ_ONLY } = useConfig();
  const { t } = useTranslation();
  const location = useLocation();
  useEffect(() => {
    const mode = getInitialLogicalTheme();
    const daisy = logicalToDaisyTheme(mode);
    document.documentElement.setAttribute('data-theme', daisy);
    try {
      localStorage.setItem(THEME_STORAGE_KEY, mode);
    } catch {
      void 0;
    }
  }, []);

  return (
    <nav className="sticky top-0 z-50 bg-base-100/80 backdrop-blur-lg border-b border-base-300">
      <div className="max-w-4xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center">
            <a
              className="flex items-center text-lg font-bold tracking-tight text-base-content hover:text-primary transition-colors duration-200 px-2 py-1 rounded-md hover:bg-base-200"
              href="/"
            >
              <img
                src="/logo-transparent.png"
                alt="Yopass logo"
                className="h-9"
              />
            </a>
          </div>
          <div className="flex items-center gap-2">
            {!READ_ONLY &&
              (!DISABLE_UPLOAD && location.pathname === '/upload' ? (
                <a
                  className="flex items-center gap-2 px-3 py-2 text-sm font-medium text-base-content/70 hover:text-base-content hover:bg-base-200 rounded-md transition-all duration-200"
                  href="#/"
                  title={t('header.buttonText')}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-5 h-5"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25h-15a2.25 2.25 0 0 1-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75"
                    />
                  </svg>
                  {t('header.buttonText')}
                </a>
              ) : (
                !DISABLE_UPLOAD && (
                  <a
                    className="flex items-center gap-2 px-3 py-2 text-sm font-medium text-base-content/70 hover:text-base-content hover:bg-base-200 rounded-md transition-all duration-200"
                    href="#/upload"
                    title={t('header.buttonUpload')}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      strokeWidth={1.5}
                      stroke="currentColor"
                      className="w-5 h-5"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5m-13.5-9L12 3m0 0 4.5 4.5M12 3v13.5"
                      />
                    </svg>
                    {t('header.buttonUpload')}
                  </a>
                )
              ))}
          </div>
        </div>
      </div>
    </nav>
  );
}
