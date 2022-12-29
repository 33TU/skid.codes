import { Link } from "preact-router";
import Icon from "../../assets/icon.png";
import "./style.css";

export default () => (
  <header class="p-3 border-bottom">
    <div class="container">
      <div class="d-flex flex-wrap align-items-center justify-content-center justify-content-lg-start">
        <Link class="navbar-brand" href="/">
          <img class="logo" src={Icon} />
        </Link>

        <ul class="nav col-12 col-lg-auto me-lg-auto mb-2 justify-content-center mb-md-0">
        </ul>

        <form class="col-12 col-lg-auto mb-3 mb-lg-0 me-lg-3">
          <input
            type="search"
            class="form-control form-control-dark"
            placeholder="Search..."
            aria-label="Search"
          />
        </form>

        <div class="text-end">
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
      </div>
    </div>
  </header>
);
