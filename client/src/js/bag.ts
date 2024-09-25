function initBag() {
  const close = document.getElementById("close");
  if (close) {
    close.addEventListener("click", function (event) {
      var dialog = document.getElementById("preview") as HTMLDialogElement;

      if (dialog) {
        dialog.close();
        event.stopPropagation();
      }
    });
  }
}

if (document.readyState !== "loading") {
  initBag();
}

document.addEventListener("DOMContentLoaded", function () {
  initBag();
});

document.addEventListener("htmx:afterSettle", function () {
  initBag();
});
