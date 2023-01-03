import "@fortawesome/fontawesome-free/css/all.min.css";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap/dist/js/bootstrap.min.js";

import "./style.css";

import { render } from "preact";
import { App } from "./app";

render(<App />, document.getElementById("app") as HTMLElement);
