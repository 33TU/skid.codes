import { Route, Router } from "preact-router";

import Header from "./components/header";
import Home from "./routes/home";
import Signin from "./routes/signin";
import Signup from "./routes/signup";

export function App() {
  return (
    <div
      id="app"
      data-bs-theme="dark"
      class="container-fluid p-0 min-vh-100 text-body bg-body "
    >
      <Header></Header>

      <div class="container">
        <Router>
          <Route path="/" component={Home} />
          <Route path="/signin" component={Signin} />
          <Route path="/signup" component={Signup} />
        </Router>
      </div>
    </div>
  );
}
