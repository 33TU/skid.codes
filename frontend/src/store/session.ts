import { signal } from "@preact/signals";
import { AuthResultSession, logout, refresh } from "../api/auth";
import { clearAuthToken, setAuthToken } from "../axios";

export const authSession = signal<AuthResultSession | undefined>(undefined);

export async function refreshSession(): Promise<void> {
  try {
    const session = await refresh();
    setAuthToken(session.auth);
    authSession.value = session.res;
  } catch (err) {
    console.error(err);
  }
}

export async function clearSession(): Promise<void> {
  try {
    await logout();
    clearAuthToken();
    authSession.value = undefined;
  } catch (err) {
    console.error(err);
  }
}

const refreshIntervalRate = 60 * 1000;
let refreshInterval: number | undefined;

authSession.subscribe((session) => {
  clearInterval(refreshInterval);

  if (!session) {
    refreshInterval = undefined;
    return;
  }

  refreshInterval = window.setInterval(refreshSession, refreshIntervalRate);
});
