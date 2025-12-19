/**
 * Environment detection utilities for dual-mode operation
 * Determines whether the app is running in desktop (Wails) or web (browser) mode
 */

/**
 * Check if the app is running in desktop mode (Wails)
 * @returns {boolean} True if running in Wails desktop app
 */
export const isDesktopMode = () => {
  return typeof window !== 'undefined' && window.wails !== undefined;
};

/**
 * Check if the app is running in web mode (browser)
 * @returns {boolean} True if running in web browser
 */
export const isWebMode = () => {
  return !isDesktopMode();
};

/**
 * Get the API base URL for HTTP requests
 * Only relevant in web mode; returns null for desktop mode
 * @returns {string|null} Base URL for API requests or null if in desktop mode
 */
export const getAPIBaseURL = () => {
  if (isDesktopMode()) {
    return null; // Desktop mode uses Wails bindings, not HTTP
  }

  // In web mode, use environment variable or default to localhost:8080
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
};

/**
 * Get the current mode as a string
 * @returns {string} 'desktop' or 'web'
 */
export const getMode = () => {
  return isDesktopMode() ? 'desktop' : 'web';
};
