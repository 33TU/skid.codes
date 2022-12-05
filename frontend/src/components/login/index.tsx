import { h } from "preact";
import { useContext, useState } from "preact/hooks";
import { AuthContext } from "../../context";
import { login } from "../../api/auth";
import { setAuthToken } from "../../axios";

const Login = () => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const { setAuthState } = useContext(AuthContext);

  async function submitLogin(): Promise<void> {
    try {
      const res = await login({
        username,
        password,
      });

      setAuthToken(res.auth);

      setAuthState({
        authUser: res,
      });
    } catch (err) {
    }
  }

  return (
    <li class="nav-item dropdown">
      <a
        class="nav-link dropdown-toggle"
        role="button"
        data-bs-toggle="dropdown"
        aria-expanded="true"
        aria-haspopup="true"
      >
        <strong class="text-white">Login</strong>
      </a>
      <div class="dropdown-menu dropdown-menu-end dropdown-menu-dark">
        <form class="px-4 py-3" style="min-width:300px;">
          <div class="form-group">
            <label for="exampleDropdownFormEmail1">Email or username.</label>
            <input
              value={username}
              onChange={(e: any) => setUsername(e.target.value)}
              class="form-control"
              placeholder="Username or email."
            />
          </div>
          <div class="form-group">
            <label for="exampleDropdownFormPassword1">Password</label>
            <input
              value={password}
              onChange={(e: any) => setPassword(e.target.value)}
              type="password"
              class="form-control"
              id="exampleDropdownFormPassword1"
              placeholder="Password"
            />
          </div>
          <button
            type="button"
            onClick={() => submitLogin()}
            class="btn btn-primary"
          >
            Sign in
          </button>
        </form>
      </div>
    </li>
  );
};

export default Login;
