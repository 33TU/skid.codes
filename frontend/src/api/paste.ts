import { authorized, axios } from "../axios";

export interface Language {
  ext: string[];
  mime: string;
  mode: string;
  name: string;
}

export interface Paste {
  id: string;
  uid: number;
  title: string | null;
  content: string | null;
  created: Date;
  private: boolean;
  language: Language;
  password: boolean;
  unlisted: boolean;
  username: string;
}

export interface FindPasteResponse {
  count: number;
  offset: number;
  pastes: Paste[];
}

export interface FetchPasteResponse {
  id: string;
  uid: number;
  title: string | null;
  private: boolean;
  unlisted: boolean;
  created: Date;
  password: boolean;
  content: string | null;
  language: Language;
}

export interface CreatePasteResponse {
  id: string;
}

export interface UpdatePasteResponse {
  id: string;
}

export interface DeletePasteResponse {
  id: string;
}

/**
 * Fetches a paste.
 * Pastes with password will return null content on wrong password.
 * Authorization is required for private pastes.
 */
export async function fetchPaste(req: {
  id: string;
  password?: string;
}): Promise<FetchPasteResponse> {
  const path = authorized() ? "/api/paste/ufetch" : "/api/paste/fetch";
  const res = await axios.post(path, req);
  return res.data;
}

/**
 * Finds pastes.
 * Returned order is by latest.
 * Authorization is required for private pastes.
 */
export async function findPaste(req: {
  uid?: number;
  username?: string;
  language?: string;
  title?: string;
  content?: string;
  private?: boolean;
  unlisted?: boolean;
  password?: boolean;
  createdBegin?: boolean;
  createdEnd?: boolean;
  offset: number;
  count: number;
}): Promise<FetchPasteResponse> {
  const path = authorized() ? "/api/paste/ufind" : "/api/paste/find";
  const res = await axios.post(path, req);
  return res.data;
}

/**
 * Creates paste.
 * Authorization is required.
 */
export async function createPaste(req: {
  language: string;
  content: string;
  title?: string;
  password?: string;
  private: boolean;
  unlisted: boolean;
}): Promise<CreatePasteResponse> {
  const res = await axios.post("/api/paste/create", req);
  return res.data;
}

/**
 * Update paste.
 * Authorization is required.
 */
export async function updatePaste(req: {
  id: string;
  language?: string;
  content?: string;
  title?: string;
  password?: string;
  private?: boolean;
  unlisted?: boolean;
}): Promise<UpdatePasteResponse> {
  const res = await axios.post("/api/paste/update", req);
  return res.data;
}

/**
 * Delete paste.
 * Authorization is required.
 */
export async function deletePaste(req: {
  id: string;
}): Promise<DeletePasteResponse> {
  const res = await axios.post("/api/paste/delete", req);
  return res.data;
}
