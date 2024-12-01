/**
 * Resets the mobile navigation menu state.
 */
function resetMenuState() {
  console.log("Resetting menu state...");
  var burgerMenu = document.getElementById("burgerMenu");
  var navLinks = document.getElementById("mobileNavLinks");
  var bar1 = document.getElementById("bar1");
  var bar2 = document.getElementById("bar2");
  var bar3 = document.getElementById("bar3");

  if (navLinks) {
    navLinks.classList.add("hidden"); // Forcefully hide the menu
    navLinks.setAttribute("aria-hidden", "true"); // Set ARIA attributes for accessibility
    console.log("Menu hidden:", navLinks.classList.contains("hidden"));
  }

  if (burgerMenu) {
    burgerMenu.setAttribute("aria-expanded", "false"); // Indicate menu is closed
  }

  if (bar1 && bar2 && bar3) {
    bar1.classList.remove("rotate-45", "translate-y-2");
    bar1.classList.add("rotate-0");

    bar2.classList.remove("opacity-0");

    bar3.classList.remove("-rotate-45", "-translate-y-2");
    bar3.classList.add("rotate-0");
    console.log("Burger menu reset.");
  }
}

/**
 * Resets the dialog state.
 */
function resetDialogState() {
  console.log("Resetting dialog state...");
  var bagic = document.getElementById("bagic") as HTMLElement;
  console.log("HTMX beforeSwap event triggered");
  bagic.lastChild!.remove();
  const content = document.createElement("div");
  content.innerHTML = `<div hx-get="/bag" hx-trigger="load" hx-swap="outerHTML"></div>`;
  bagic.appendChild(content);
}

function initHeader() {
  var burgerMenu = document.getElementById("burgerMenu");
  var navLinks = document.getElementById("mobileNavLinks");
  var bar1 = document.getElementById("bar1");
  var bar2 = document.getElementById("bar2");
  var bar3 = document.getElementById("bar3");
  var bagic = document.getElementById("bagic");

  if (burgerMenu && !(burgerMenu as any)._hasListener) {
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
    (burgerMenu as any)._hasListener = true;
  }

  // Adjusting span width to wrap the text continuously
  const span = document.getElementById("mtx");
  const div = document.getElementById("mbg");
  if (div && span) {
    const divWidth = div.offsetWidth;
    const spanWidth = span.offsetWidth;
    const clonesNeeded = Math.ceil(divWidth / spanWidth) + 1;

    for (let i = 0; i < clonesNeeded; i++) {
      const clone = span.cloneNode(true) as HTMLElement;
      clone.classList.add("mtc");
      if (span.parentNode) {
        span.parentNode.appendChild(clone);
      }
    }
  }

  if (bagic && !(bagic as any)._hasListener) {
    bagic.addEventListener("click", function () {
      var dialog = document.getElementById("preview") as HTMLDialogElement;

      if (dialog && !dialog.open) {
        dialog.showModal();
      }
    });
    (bagic as any)._hasListener = true;
  }
}

if (document.readyState !== "loading") {
  initHeader();
} else {
  document.addEventListener("DOMContentLoaded", function () {
    initHeader();
  });
}

document.addEventListener("htmx:afterSettle", function () {
  console.log("HTMX afterSettle event triggered");
  initHeader();
});

document.body.addEventListener("htmx:historyRestore", () => {
  console.log("HTMX historyRestore event triggered");
  resetMenuState();
  resetDialogState();
  initHeader();
});
