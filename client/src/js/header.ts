function initHeader() {
  var burgerMenu = document.getElementById("burgerMenu");
  var navLinks = document.getElementById("mobileNavLinks");
  var bar1 = document.getElementById("bar1");
  var bar2 = document.getElementById("bar2");
  var bar3 = document.getElementById("bar3");
  var bagic = document.getElementById("bagic");

  if (burgerMenu) {
    burgerMenu.addEventListener("click", function () {
      if (navLinks && bar1 && bar2 && bar3) {
        navLinks.classList.toggle("hidden");
        if (bar1.classList.contains("rotate-0")) {
          bar1.classList.remove("rotate-0");
          bar1.classList.add("rotate-45", "translate-y-2");

          bar2.classList.remove("rotate-0");
          bar2.classList.add("opacity-0");

          bar3.classList.remove("rotate-0");
          bar3.classList.add("-rotate-45", "-translate-y-2");
        } else {
          bar1.classList.remove("rotate-45", "translate-y-2");
          bar1.classList.add("rotate-0");

          bar2.classList.remove("opacity-0");
          bar3.classList.remove("-rotate-45", "-translate-y-2");
          bar3.classList.add("rotate-0");
        }
      }
    });
  }

  // Adjusting span width to wrap the text continuously
  const span = document.getElementById("mtx");
  const div = document.getElementById("mbg");
  if (div && span) {
    const divWidth = div.offsetWidth;
    const spanWidth = span.offsetWidth;
    const clonesNeeded = Math.ceil(divWidth / spanWidth) + 1;

    for (let i = 0; i < clonesNeeded; i++) {
      const clone = span.cloneNode(true);
      if (span.parentNode) {
        span.parentNode.appendChild(clone);
      }
    }
  }

  if (bagic) {
    bagic.addEventListener("click", function () {
      var dialog = document.getElementById("preview") as HTMLDialogElement;

      if (dialog && !dialog.open) {
        dialog.showModal();
      }
    });
  }
}

if (document.readyState !== "loading") {
  initHeader();
}

document.addEventListener("DOMContentLoaded", function () {
  initHeader();
});

document.addEventListener("htmx:afterSettle", function () {
  initHeader();
});
