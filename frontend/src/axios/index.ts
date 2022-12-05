import { Axios, AxiosRequestConfig } from "axios";

// Base axios config
const baseConfig: AxiosRequestConfig = {
  baseURL: process.env.PREACT_APP_BASE_URL,
  headers: {
    "Content-Type": "application/json"
  }
};

// Authorization token
let authToken: string | undefined;

/**
 * Axios for for API points which do not require auth_token.
 */
export const baseAxios = new Axios(baseConfig);

/**
 * Axios for for API points which require auth_token.
 */
export const authAxios = new Axios(baseConfig);

// Adds auth token to header.
authAxios.interceptors.request.use((config) => {
  Object.assign(config, baseConfig);

  if (!authToken) {
    throw new Error("not authorized");
  }

  if (config.headers) {
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