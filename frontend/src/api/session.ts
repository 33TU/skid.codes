import { authAxios } from "../axios";

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
export async function findSessions(req: {
  count: number;
  offset?: number;
}): Promise<FindSessionResult> {
  const r = await authAxios.post("/api/session/find", JSON.stringify(req));
  return JSON.parse(r.data);
}

/**
 * Revokes a session.
 * Returned order is by latest.
 * Authorization required.
 */
export async function revokeSession(req: {
  id: string;
}): Promise<RevokeSessionResult> {
  const r = await authAxios.post("/api/session/revoke", JSON.stringify(req));
  return JSON.parse(r.data);
}
