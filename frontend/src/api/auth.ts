import { axios } from "../axios";
import { RevokeSessionResponse } from "./session";

export interface AuthResponseSession {
  sid: number;
  uid: number;
  username: string;
  email: string;
  role: string;
  exp: number;
}

export interface AuthResponse {
  auth: string;
  res: AuthResponseSession;
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
export async function login(req: LoginRequest): Promise<AuthResponse> {
  const res = await axios.post("/api/auth/login", req);
  return res.data;
}

/**
 * Refresh auth and refresh token.
 * NOTE: the refresh_token is always stores as httponly cookie.
 *  Returns auth token and user information.
 */
export async function refresh(): Promise<AuthResponse> {
  const res = await axios.get("/api/auth/refresh");
  return res.data;
}

/**
 * Logout and revoke current session.
 * Returns revoked session.
 */
export async function logout(): Promise<RevokeSessionResponse> {
  const res = await axios.get("/api/auth/logout");
  return res.data;
}
