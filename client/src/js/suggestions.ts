function init() {
  document
    .getElementById("suggestion")
    .addEventListener("click", function (event) {
      document.getElementById("address").value = event.target.textContent;
      this.style.display = "none";
    });
}

if (document.readyState !== "loading") {
  init();
}

document.addEventListener("DOMContentLoaded", init);

document.addEventListener("htmx:afterSettle", function (_) {
  init();
});
