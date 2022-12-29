import { Link } from "preact-router";
import { useState } from "preact/hooks";

import { login } from "../../api/auth";
import { isApiError } from "../../api/error";
import { setAuthToken } from "../../axios";
import { authSession } from "../../store";

import "./style.css";
import Icon from "../../assets/icon.png";

export default () => {
  const [user, setUser] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");

  const submitLogin = async () => {
    try {
      const emailRegex = (/\S+@\S+\.\S+/);

      const loginReq = await login(
        emailRegex.test(user)
          ? { email: user, password }
          : { username: user, password },
      );

      // Set auth token
      setAuthToken(loginReq.auth);

      // Set authSession
      authSession.value = loginReq.res;
    } catch (err) {
      setError(isApiError(err)?.message ?? String(err));
    }
  };

  return (
    <div class="text-center body-signin">
      <main class="form-signin">
        {error
          ? (
            <div class="alert alert-danger alert-dismissible" role="alert">
              <button
                type="button"
                class="btn-close"
                data-bs-dismiss="alert"
                onClick={() => setError("")}
              >
              </button>
              <strong>{error}</strong>
            </div>
          )
          : ""}

        <div class="card card-signin">
          <div class="card-body">
            <Link href="/">
              <img src={Icon} class="mb-3" />
            </Link>

            <div class="form-floating">
              <input
                type="user"
                class="form-control form-control-dark"
                id="floatingLoginUser"
                placeholder="."
                value={user}
                onChange={(e: any) => setUser(e.target.value)}
              />
              <label for="floatingLoginUser">Username or email</label>
            </div>
            <div class="form-floating">
              <input
                type="password"
                class="form-control form-control-dark"
                id="floatingLoginPassword"
                placeholder="."
                value={password}
                onChange={(e: any) => setPassword(e.target.value)}
              />
              <label for="floatingLoginPassword">Password</label>
            </div>

            <div class="checkbox mb-3">
              <label>
                <input type="checkbox" value="remember-me" /> Remember me
              </label>
            </div>
            <button
              class="w-100 btn btn-lg btn-primary"
              onClick={() => submitLogin()}
            >
              Sign in
            </button>
            <p class="mt-5 mb-3 text-muted">&copy; 2022-2023</p>
          </div>
        </div>
      </main>
    </div>
  );
};
