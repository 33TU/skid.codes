import { authAxios, baseAxios } from "../axios";
import { RevokeSessionResult } from "./session";

export interface AuthResult {
  auth: string;
  res: {
    sid: number;
    uid: number;
    username: string;
    email: string;
    role: string;
    exp: number;
  };
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
export async function login(req: LoginRequest): Promise<AuthResult> {
  const r = await baseAxios.post("/api/auth/login", JSON.stringify(req));
  return JSON.parse(r.data);
}

/**
 * Refresh auth and refresh token.
 * NOTE: the refresh_token is always stores as httponly cookie.
 *  Returns auth token and user information.
 */
export async function refresh(): Promise<AuthResult> {
  const r = await baseAxios.get("/api/auth/refresh");
  return JSON.parse(r.data);
}

/**
 * Logout and revoke current session.
 * Returns revoked session.
 */
export async function logout(): Promise<RevokeSessionResult> {
  const r = await authAxios.get("/api/auth/logout");
  return JSON.parse(r.data);
}
