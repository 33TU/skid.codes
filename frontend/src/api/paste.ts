import { authAxios, authorized, baseAxios } from "../axios";

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

export interface FindPasteResult {
  count: number;
  offset: number;
  pastes: Paste[];
}

export interface FetchPasteResult {
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

export interface CreatePasteResult {
  id: string;
}

export interface UpdatePasteResult {
  id: string;
}

export interface DeletePasteResult {
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
}): Promise<FetchPasteResult> {
  const axios = authorized() ? authAxios : baseAxios;
  const path = authorized() ? "/api/paste/ufetch" : "/api/paste/fetch"
  const r = await axios.post(path, JSON.stringify(req));
  return JSON.parse(r.data);
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
}): Promise<FetchPasteResult> {
  const axios = authorized() ? authAxios : baseAxios;
  const path = authorized() ? "/api/paste/ufind" : "/api/paste/find"
  const r = await axios.post(path, JSON.stringify(req));
  return JSON.parse(r.data);
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
}): Promise<CreatePasteResult> {
  const r = await authAxios.post("/api/paste/create", JSON.stringify(req));
  return JSON.parse(r.data);
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
}): Promise<UpdatePasteResult> {
  const r = await authAxios.post("/api/paste/update", JSON.stringify(req));
  return JSON.parse(r.data);
}

/**
 * Delete paste.
 * Authorization is required.
 */
export async function deletePaste(req: {
  id: string;
}): Promise<DeletePasteResult> {
  const r = await authAxios.post("/api/paste/delete", JSON.stringify(req));
  return JSON.parse(r.data);
}
