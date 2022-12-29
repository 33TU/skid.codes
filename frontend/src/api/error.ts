import { AxiosError } from "axios";

export interface ApiError {
  code: number;
  message: string;
}

export function isApiError(err: unknown): ApiError | undefined {
  if (
    err instanceof AxiosError &&
    err.response &&
    err.response.data &&
    err.response.data.code &&
    err.response.data.message
  ) {
    return err.response?.data as ApiError;
  }

  return undefined;
}
