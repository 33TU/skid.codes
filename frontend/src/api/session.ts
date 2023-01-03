import { axios } from "../axios";

export interface Session {
  id: number;
  country: string;
  created: Date;
  updated: Date;
  revoked: boolean;
}

export interface FindSessionResponse {
  count: number;
  offset: number;
  sessions: Session[];
}

export interface RevokeSessionResponse {
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
}): Promise<FindSessionResponse> {
  const res = await axios.post("/api/session/find", req);
  return res.data;
}

/**
 * Revokes a session.
 * Returned order is by latest.
 * Authorization required.
 */
export async function revokeSession(req: {
  id: string;
}): Promise<RevokeSessionResponse> {
  const res = await axios.post("/api/session/revoke", req);
  return res.data;
}
