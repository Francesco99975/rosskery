function initSuggestions() {
  const suggestions = document.getElementById("suggestions");
  if (suggestions && !(suggestions as any).hasListeners) {
    suggestions.addEventListener("click", function (event) {
      const target = event.target as HTMLElement;
      const addressElem = document.getElementById(
        "address"
      ) as HTMLInputElement;
      if (addressElem) {
        addressElem.value = target.textContent || "";
        this.style.display = "none";
      }
    });
    (suggestions as any).hasListeners = true;
  }
}

if (document.readyState !== "loading") {
  initSuggestions();
} else {
  document.addEventListener("DOMContentLoaded", initSuggestions);
}

document.addEventListener("htmx:afterSettle", function (_) {
  initSuggestions();
});
