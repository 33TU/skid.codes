import { h } from "preact";
import { Route, Router } from "preact-router";
import { useState } from "preact/hooks";

import Header from "./header";
import Home from "../routes/home";
import Paste from "../routes/paste";

import { AuthContext } from "../context";

const App = () => {
  const [authState, setAuthState] = useState({});

  return (
    <AuthContext.Provider value={{ authState, setAuthState }}>
      <div id="app">
        <div class="container">
          <Header />
        </div>

        <div class="container">
          <Router>
            <Route path="/" component={Home} />
            <Route path="/paste/:id" component={Paste} />
          </Router>
        </div>
      </div>
    </AuthContext.Provider>
  );
};

export default App;
