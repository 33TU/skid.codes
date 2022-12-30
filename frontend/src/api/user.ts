import { axios } from "../axios";

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
  const res = await axios.get(`/api/user/${encodeURI(username)}`);
  return res.data;
}

/**
 * Creates user account.
 */
export async function createUser(req: {
  username: string;
  email: string;
  password: string;
}): Promise<CreateUserResponse> {
  const res = await axios.post("/api/user/create", req);
  return res.data;
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
  const res = await axios.post("/api/user/find", req);
  return res.data;
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
  const res = await axios.post("/api/user/update", req);
  return res.data;
}
