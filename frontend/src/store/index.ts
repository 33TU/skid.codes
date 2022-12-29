import { signal } from "@preact/signals";
import { AuthResultSession } from "../api/auth";

export const authSession = signal<AuthResultSession | undefined>(undefined);
