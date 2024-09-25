import "./css/style.css";

import htmx from "htmx.org";

import { UCounter } from "./components/ucounter";

declare global {
  interface Window {
    htmx: typeof htmx;
    conn: WebSocket;
    visited: boolean;
  }

  interface Date {
    fp_incr(increment: number): Date; // Declare the new method
  }

  interface HTMLElementTagNameMap {
    "u-counter": UCounter;
  }
}

window.htmx = htmx;
window.customElements.define("u-counter", UCounter);
window.conn = new WebSocket("ws://" + document.location.host + "/ws");
window.visited = false;

function sendVisit() {
  if (window.visited) return;

  // Onopen
  window.conn.onopen = function (evt) {
    window.conn.send(JSON.stringify({ type: "visit", payload: "" }));
    window.visited = true;
  };
}

if (document.readyState !== "loading") {
  sendVisit();
}

document.addEventListener("DOMContentLoaded", function () {
  sendVisit();
});
