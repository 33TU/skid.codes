import { createContext } from "preact";
import { AuthResult } from "../api/auth";

export interface AuthState {
    authUser?: AuthResult;
}

export const AuthContext = createContext({
    authState: {} as AuthState,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    setAuthState: (auth: AuthState) => { },
});
