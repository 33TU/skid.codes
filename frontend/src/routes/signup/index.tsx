import { Link } from "preact-router";
import { useState } from "preact/hooks";

import { createUser } from "../../api/user";
import { isApiError } from "../../api/error";

import "./style.css";
import Icon from "../../assets/icon.png";

export default () => {
  const [username, setUsername] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");

  const submitSignup = async () => {
    try {
      const loginReq = await createUser({
        username,
        email,
        password,
      });

      setError(JSON.stringify(loginReq));
    } catch (err) {
      setError(isApiError(err)?.message ?? String(err));
    }
  };

  return (
    <div class="text-center body-signup">
      {error
        ? (
          <div
            class="alert alert-signup alert-danger alert-dismissible position-fixed start-50 translate-middle"
            role="alert"
          >
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

      <main class="form-signup">
        <div class="card card-signup">
          <div class="card-body">
            <Link href="/">
              <img src={Icon} class="mb-3" />
            </Link>

            <div class="form-floating">
              <input
                type="text"
                class="form-control form-control-dark"
                id="floatingInput"
                placeholder="."
                value={username}
                onChange={(e: any) => setUsername(e.target.value)}
              />
              <label for="floatingInput">Username</label>
            </div>
            <div class="form-floating">
              <input
                type="email"
                class="form-control form-control-dark"
                id="floatingInput"
                placeholder="."
                value={email}
                onChange={(e: any) => setEmail(e.target.value)}
              />
              <label for="floatingInput">Email address</label>
            </div>
            <div class="form-floating">
              <input
                type="password"
                class="form-control form-control-dark"
                id="floatingPassword"
                placeholder="."
                value={password}
                onChange={(e: any) => setPassword(e.target.value)}
              />
              <label for="floatingPassword">Password</label>
            </div>

            <div class="checkbox mb-3">
              <label>
                <input type="checkbox" value="remember-me" />{" "}
                Agree to terms of services
              </label>
            </div>
            <button
              class="w-100 btn btn-lg btn-primary"
              type="submit"
              onClick={() => submitSignup()}
            >
              Sign up
            </button>
            <p class="mt-5 mb-3 text-muted">&copy; 2022-2023</p>
          </div>
        </div>
      </main>
    </div>
  );
};
