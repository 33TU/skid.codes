import { Route, Router } from "preact-router";
import { refreshSession } from "./store/session";

import Header from "./components/header";
import Home from "./routes/home";
import Signin from "./routes/signin";
import Signup from "./routes/signup";

export function App() {
  refreshSession();

  return (
    <div
      data-bs-theme="dark"
      class="container-fluid p-0 text-body bg-body h-100"
    >
      <Header></Header>

      <Router>
        <Route path="/" component={Home} />
        <Route path="/signin" component={Signin} />
        <Route path="/signup" component={Signup} />
      </Router>
    </div>
  );
}
