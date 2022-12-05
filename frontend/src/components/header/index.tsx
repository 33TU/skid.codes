import { h } from "preact";
import { Link } from "preact-router/match";

import style from "./style.css";
import Login from "../login";

import { AuthContext } from "../../context";
import { authorized, clearAuthToken, setAuthToken } from "../../axios";
import { AuthResult, logout, refresh } from "../../api/auth";
import { useContext, useEffect } from "preact/hooks";

const Header = () => {
  const { authState, setAuthState } = useContext(AuthContext);

  // Refresh refresh_token and auth_token
  useEffect(() => {
    const r = async () => {
      try {
        const r = await refresh();

        setAuthToken(r.auth);
        setAuthState({
          authUser: r,
        });
      } catch (_) {
        clearAuthToken();
        setAuthState({
          authUser: undefined,
        });
      }
    };

    // Refresh once
    r();

    // Refresh every 30 seconds
    setInterval(() => {
      if (!authorized) return;
      r();
    }, 1000 * 10);
  }, [setAuthState]);

  // doLogout logsout
  const doLogout = async () => {
    await logout();

    clearAuthToken();
    setAuthState({
      authUser: undefined,
    });
  };

  const notAuthorizedElement = () => (
    <ul class="navbar-nav ml-auto dropdown">
      <Login />
      <li class="nav-item mt-1">
        <Link
          class="nav-link btn btn-primary btn-sm py-1 font-display font-sm font-weight:600 text-white"
          href="/register"
        >
          Register
        </Link>
      </li>
    </ul>
  );

  const authorizedElement = (authUser: AuthResult) => (
    <ul class="navbar-nav">
      <li class="nav-item dropdown">
        <a
          class="nav-link dropdown-toggle"
          href="#"
          role="button"
          data-bs-toggle="dropdown"
          aria-expanded="false"
        >
          <span>
            Logged in as <strong>{authUser.res.username}</strong>
          </span>
        </a>
        <ul class="dropdown-menu dropdown-menu-dark">
          <li>
            <a class="dropdown-item" href="#">Manage pastes</a>
          </li>
          <li>
            <a class="dropdown-item" href="#">Manage sessions</a>
          </li>
          <li>
            <hr class="dropdown-divider" />
          </li>
          <li>
            <a class="dropdown-item" href="#" onClick={doLogout}>Logout</a>
          </li>
        </ul>
      </li>
    </ul>
  );

  return (
    <header class={style.header}>
      <nav class="navbar navbar-expand-md navbar-dark bg-dark main-nav">
        <div class="container-fluid">
          <Link class="navbar-brand" href="/">
            <img src="/assets/favicon.png" />
            <span class="d-none d-lg-inline">skid.codes</span>
          </Link>

          <button
            class="navbar-toggler"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#navbarSupportedContent"
            aria-controls="navbarSupportedContent"
            aria-expanded="false"
            aria-label="Toggle navigation"
          >
            <span class="navbar-toggler-icon" />
          </button>
          <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav me-auto mb-2 mb-lg-0">
              <li>
                <Link
                  class="nav-link"
                  activeClassName="nav-link active"
                  href="/"
                >
                  Home
                </Link>
              </li>
            </ul>

            {authState.authUser
              ? authorizedElement(authState.authUser)
              : notAuthorizedElement()}
          </div>
        </div>
      </nav>
    </header>
  );
};

export default Header;
