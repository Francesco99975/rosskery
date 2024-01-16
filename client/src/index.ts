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
  }
}

window.htmx = htmx;
window.PhotoSwipe = PhotoSwipe;
window.PhotoSwipeLightbox = PhotoSwipeLightbox;
window.customElements.define("u-counter", UCounter);
