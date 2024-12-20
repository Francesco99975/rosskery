import PhotoSwipe from "photoswipe";
import PhotoSwipeLightbox from "photoswipe/lightbox";
import "photoswipe/style.css";

function initGallery() {
  const lightbox = new PhotoSwipeLightbox({
    gallery: "#gallery",
    children: "a",
    pswpModule: PhotoSwipe,
  });
  lightbox.init();
}

if (document.readyState !== "loading") {
  initGallery();
} else {
  document.addEventListener("DOMContentLoaded", function () {
    initGallery();
  });
}

document.addEventListener("htmx:afterSettle", function () {
  initGallery();
});
