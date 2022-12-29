import { axios } from "../axios";

export interface Session {
  id: number;
  country: string;
  created: Date;
  updated: Date;
  revoked: boolean;
}

export interface FindSessionResult {
  count: number;
  offset: number;
  sessions: Session[];
}

export interface RevokeSessionResult {
  id: number;
  country: string;
  created: Date;
  updated: Date;
  revoked: boolean;
}

/**
 * Finds all sessions for authorized session.
 * Returned order is by latest.
 * Authorization required.
 */
export function findSessions(req: {
  count: number;
  offset?: number;
}): Promise<FindSessionResult> {
  return axios.post("/api/session/find", req);
}

/**
 * Revokes a session.
 * Returned order is by latest.
 * Authorization required.
 */
export function revokeSession(req: {
  id: string;
}): Promise<RevokeSessionResult> {
  return axios.post("/api/session/revoke", req);
}
