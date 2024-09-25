import PhotoSwipeLightbox from "photoswipe/lightbox";
import PhotoSwipe from "photoswipe";
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
}

document.addEventListener("DOMContentLoaded", function () {
  initGallery();
});
