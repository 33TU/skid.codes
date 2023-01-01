import { Link } from "preact-router/match";
import { authSession, clearSession } from "../../store/session";

import Icon from "../../assets/icon.png";
import "./style.css";

export default () => {
  const unauthElement = (
    <div class="text-end px-3">
      <Link
        href="/signin"
        type="button"
        class="btn btn-outline-light me-2"
      >
        Login
      </Link>

      <Link href="/signup" type="button" class="btn btn-primary">
        Sign-up
      </Link>
    </div>
  );

  const authElement = (
    <li class="nav-item dropdown">
      <a
        class="nav-link dropdown-toggle"
        href="#"
        role="button"
        data-bs-toggle="dropdown"
        aria-expanded="true"
      >
        <span>
          <strong>{authSession.value?.username}</strong>
        </span>
      </a>

      <ul class="dropdown-menu">
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
          <a class="dropdown-item" href="#" onClick={() => clearSession()}>
            Logout
          </a>
        </li>
      </ul>
    </li>
  );

  return (
    <header class="border-bottom">
      <nav class="navbar navbar-expand-lg bg-body-tertiary">
        <div class="container">
          <a class="navbar-brand" href="#">
            <Link class="navbar-brand" href="/">
              <img class="logo" src={Icon} />
            </Link>
          </a>

          <ul class="navbar-nav d-block d-lg-none d-xl-block d-xl-none">
            {!authSession.value ? unauthElement : (
              <span>
                Logged in as <strong>{authSession.value?.username}</strong>
              </span>
            )}
          </ul>

          <button
            class="navbar-toggler"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#navbarSupportedContent"
            aria-controls="navbarSupportedContent"
            aria-expanded="false"
            aria-label="Toggle navigation"
          >
            <span class="navbar-toggler-icon"></span>
          </button>

          <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav me-auto mb-2 mb-lg-0">
              <li class="nav-item">
                <Link
                  class="nav-link"
                  activeClassName="nav-link active"
                  href="/"
                >
                  Home
                </Link>
              </li>

              <li class="nav-item">
                <Link
                  class="nav-link"
                  activeClassName="nav-link active"
                  href="/members"
                >
                  Members
                </Link>
              </li>

              <li class="nav-item">
                <Link
                  class="nav-link"
                  activeClassName="nav-link active"
                  href="/pastes"
                >
                  Pastes
                </Link>
              </li>
            </ul>

            <form class="d-flex" role="search">
              <input
                class="form-control form-control-dark"
                type="search"
                placeholder="Quick Search"
                aria-label="Quick Search"
              />
            </form>

            <ul class="navbar-nav">
              {authSession.value
                ? authElement
                : (
                  <div class="d-none d-lg-block d-xl-block">
                    {unauthElement}
                  </div>
                )}
            </ul>
          </div>
        </div>
      </nav>
    </header>
  );
};
