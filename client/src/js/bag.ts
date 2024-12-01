function initBag() {
  const close = document.getElementById("close");
  if (close && !(close as any)._hasListener) {
    close.addEventListener("click", function (event) {
      var dialog = document.getElementById("preview") as HTMLDialogElement;

      if (dialog) {
        dialog.close();
        event.stopPropagation();
      }
    });
    (close as any)._hasListener = true;
  }
}

if (document.readyState !== "loading") {
  initBag();
} else {
  document.addEventListener("DOMContentLoaded", function () {
    initBag();
  });
}

document.addEventListener("htmx:afterSettle", function () {
  initBag();
});
