import { axios } from "../axios";
import { RevokeSessionResult } from "./session";

export interface AuthResultSession {
  sid: number;
  uid: number;
  username: string;
  email: string;
  role: string;
  exp: number;
}

export interface AuthResult {
  auth: string;
  res: AuthResultSession;
}

export interface LoginRequestWithEmail {
  email: string;
  password: string;
}

export interface LoginRequestWithUsername {
  username: string;
  password: string;
}

export type LoginRequest = LoginRequestWithUsername | LoginRequestWithEmail;

/**
 * Sends login request with credentials.
 * Returns auth token and user information.
 */
export function login(req: LoginRequest): Promise<AuthResult> {
  return axios.post("/api/auth/login", req);
}

/**
 * Refresh auth and refresh token.
 * NOTE: the refresh_token is always stores as httponly cookie.
 *  Returns auth token and user information.
 */
export function refresh(): Promise<AuthResult> {
  return axios.get("/api/auth/refresh");
}

/**
 * Logout and revoke current session.
 * Returns revoked session.
 */
export function logout(): Promise<RevokeSessionResult> {
  return axios.get("/api/auth/logout");
}
