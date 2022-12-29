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
export function getUser(username: string): Promise<CreateUserResponse> {
  return axios.get(`/api/user/${encodeURI(username)}`);
}

/**
 * Creates user account.
 */
export function createUser(req: {
  username: string;
  email: string;
  password: string;
}): Promise<CreateUserResponse> {
  return axios.post("/api/user/create", req);
}

/**
 * Finds user by username.
 * Count is the limit of results.
 * Offset is the begin offset of search.
 */
export function findUser(req: {
  username: string;
  count: number;
  offset?: number;
}): Promise<CreateUserResponse> {
  return axios.post("/api/user/find", req);
}

/**
 * Updates user information.
 * Authorization required.
 */
export function updateUser(req: {
  username?: string;
  email?: string;
  password?: string;
}): Promise<UpdateUserResponse> {
  return axios.post("/api/user/update", req);
}
