import axiosStatic from "axios";

// Axios default export
export const axios = axiosStatic;

// Authorization token
let authToken: string | undefined;

// Axios defaults
axios.defaults.baseURL = import.meta.env.VITE_BASE_URL;

// Adds auth token to header.
axios.interceptors.request.use((config) => {
  if (authToken && config.headers) {
    config.headers.Authorization = authToken;
  }

  return config;
}, (error) => {
  // Do something with request error
  return Promise.reject(error);
});

/**
 * Checks if authorized.
 */
export function authorized(): boolean {
  return authToken !== undefined;
}

/**
 * Sets auth token.
 */
export function setAuthToken(token: string): void {
  authToken = `Bearer ${token}`;
}

/**
 * Clears auth token.
 */
export function clearAuthToken(): void {
  authToken = undefined;
}
