import { authAxios, baseAxios } from "../axios";

export interface User {
  id: number;
  username: string;
  role: string;
  pastes: number;
  online: string;
}

export interface FindUserResponse {
  count: number;
  offset: number;
  users: User[];
}

export interface CreateUserResponse {
  id: number;
  username: string;
  email: string;
  role: string;
}

export interface GetUserResponse {
  id: number;
  username: string;
  role: string;
  pastes: number;
  online: string;
}

export interface UpdateUserResponse {
  id: number;
  username: string;
  email: string;
  role: string;
}

/**
 * Gets user information.
 */
export async function getUser(username: string): Promise<CreateUserResponse> {
  const r = await baseAxios.get(`/api/user/${encodeURI(username)}`);
  return JSON.parse(r.data);
}

/**
 * Creates user account.
 */
export async function createUser(req: {
  username: string;
  email: string;
  password: string;
}): Promise<CreateUserResponse> {
  const r = await baseAxios.post("/api/user/create", JSON.stringify(req));
  return JSON.parse(r.data);
}

/**
 * Finds user by username.
 * Count is the limit of results.
 * Offset is the begin offset of search.
 */
export async function findUser(req: {
  username: string;
  count: number;
  offset?: number;
}): Promise<CreateUserResponse> {
  const r = await baseAxios.post("/api/user/find", JSON.stringify(req));
  return JSON.parse(r.data);
}

/**
 * Updates user information.
 * Authorization required.
 */
export async function updateUser(req: {
  username?: string;
  email?: string;
  password?: string;
}): Promise<UpdateUserResponse> {
  const r = await authAxios.post("/api/user/update", JSON.stringify(req));
  return JSON.parse(r.data);
}
