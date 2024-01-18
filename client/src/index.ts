import "./css/style.css";

import htmx from "htmx.org";
import PhotoSwipeLightbox from "photoswipe/lightbox";
import PhotoSwipe from "photoswipe";
import "photoswipe/style.css";

import { UCounter } from "./components/ucounter";

declare global {
  interface Window {
    htmx: typeof htmx;
    PhotoSwipe: typeof PhotoSwipe;
    PhotoSwipeLightbox: typeof PhotoSwipeLightbox;
    conn: WebSocket;
    visited: boolean;
  }
}

window.htmx = htmx;
window.PhotoSwipe = PhotoSwipe;
window.PhotoSwipeLightbox = PhotoSwipeLightbox;
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
