import "./css/style.css";

import htmx from "htmx.org";

import { UCounter } from "./components/ucounter";

declare global {
  interface Window {
    htmx: typeof htmx;
    conn: WebSocket;
    visited: boolean;
  }

  interface Date {
    fp_incr(increment: number): Date; // Declare the new method
  }

  interface HTMLElementTagNameMap {
    "u-counter": UCounter;
  }
}

const SETTINGS_CHANGED_EVENT = "settingschanged";
const VISIT_EVENT = "visit";
const PRODUCT_ADDED_EVENT = "newproduct";
const PRODUCT_UPDATED_EVENT = "updateproduct";
const PRODUCT_REMOVED_EVENT = "removeproduct";
const CATEGORY_ADDED_EVENT = "newcategory";
const CATEGORY_REMOVED_EVENT = "removecategory";

window.htmx = htmx;
htmx.config.includeIndicatorStyles = false;
window.customElements.define("u-counter", UCounter);
const host = document.location.host;
const wsProtocol = document.location.protocol === "https:" ? "wss://" : "ws://";
window.conn = new WebSocket(wsProtocol + host + "/ws");
window.visited = false;

function addProductToDOM(html: string) {
  const productShop = document.getElementById("sp");
  // const productFeutured = document.getElementById("fp");
  const productNewArrivals = document.getElementById("na");
  const productDiv = document.createElement("div");
  const csrfBox = document.getElementById("csrf_store");
  productDiv.innerHTML = html;
  const content = productDiv.children[0];
  window.htmx.process(content);

  if (productShop && csrfBox) {
    const csrfToken = csrfBox.getAttribute("value") || "";
    (content.lastChild?.firstChild as HTMLElement).setAttribute(
      "value",
      csrfToken
    );

    productShop.prepend(content);
  }

  // if (productFeutured && csrfBox) {
  // const csrfToken = csrfBox.getAttribute("value") || "";
  // (content.lastChild?.firstChild as HTMLElement).setAttribute(
  //   "value",
  //   csrfToken
  // );
  //   productFeutured.prepend(content);
  // }

  if (productNewArrivals && csrfBox) {
    const csrfToken = csrfBox.getAttribute("value") || "";
    (content.lastChild?.firstChild as HTMLElement).setAttribute(
      "value",
      csrfToken
    );
    productNewArrivals.prepend(content);
  }
}

function updateProductInDOM(productId: string, html: string) {
  deleteProductFromDOM(productId);
  addProductToDOM(html);
}

function deleteProductFromDOM(productId: string) {
  const productDiv = document.getElementById(productId);
  if (productDiv) {
    productDiv.remove();
  } else {
    console.log("Product to delete not found");
  }
}

function sendVisit() {
  if (window.visited) return;

  // Onopen
  window.conn.onopen = function (evt) {
    window.conn.send(JSON.stringify({ type: VISIT_EVENT, payload: "" }));
    window.visited = true;

    window.conn.onmessage = function (evt) {
      const data = JSON.parse(evt.data);

      switch (data.type) {
        case SETTINGS_CHANGED_EVENT:
          if (!data.payload.online || !data.payload.operative) {
            window.location.reload();
          }

          const messageBox = document.getElementById("mbg");
          const messageSpan = document.getElementById("mtx");
          console.log("Received message:", data.payload.message);
          if (messageBox && messageSpan) {
            // Remove existing clones
            const clones = messageBox.querySelectorAll(".mtc");
            clones.forEach((clone) => clone.remove());

            if (data.payload.message === "") {
              messageBox.style.display = "none";
              return;
            }

            messageBox.style.display = "flex";
            messageSpan.innerText = data.payload.message;

            const divWidth = messageBox.offsetWidth;
            const spanWidth = messageSpan.offsetWidth;
            const clonesNeeded = Math.ceil(divWidth / spanWidth) + 1;

            for (let i = 0; i < clonesNeeded; i++) {
              const clone = messageSpan.cloneNode(true) as HTMLElement;
              clone.classList.add("mtc");
              if (messageSpan.parentNode) {
                messageSpan.parentNode.appendChild(clone);
              }
            }
          } else {
            console.log("Message box not found.");
          }
          break;
        case PRODUCT_ADDED_EVENT:
          addProductToDOM(data.payload.html);
          break;
        case PRODUCT_UPDATED_EVENT:
          updateProductInDOM(data.payload.id, data.payload.html);
          break;
        case PRODUCT_REMOVED_EVENT:
          deleteProductFromDOM(data.payload.id);
          break;
        case CATEGORY_ADDED_EVENT:
          break;
        case CATEGORY_REMOVED_EVENT:
          break;
        default:
          break;
      }
    };
    console.log("WebSocket connection is listening for messages...");
  };
}

function initMain() {
  document
    .querySelectorAll(".product-container")
    .forEach((productContainer) => {
      if (!(productContainer as any)._hasListener) {
        productContainer.addEventListener("click", (event: Event) => {
          const target = event.target as HTMLElement;

          // Check if the clicked element is a button with data attributes
          if (target.matches("[data-action]")) {
            const action = target.getAttribute("data-action"); // 'increment' or 'decrement'
            const type = target.getAttribute("data-type"); // 'weight' or 'quantity'
            const productId = target.getAttribute("data-product-id"); // Product ID

            if (!action || !type || !productId) return;

            // Determine the input element based on the product ID and type
            const inputId =
              type === "weight"
                ? `weightInput-${productId}`
                : `quantityInput-${productId}`;
            const input = document.getElementById(inputId) as HTMLInputElement;

            if (!input) return;

            // Parse the current value and adjust it
            let currentValue = parseFloat(input.value);
            const step = type === "weight" ? 0.1 : 1; // Step size based on type
            currentValue += action === "increment" ? step : -step;

            // Enforce minimum limits
            if (type === "weight") {
              currentValue = Math.max(0.1, currentValue); // Min weight is 0.1
            } else {
              currentValue = Math.max(1, currentValue); // Min quantity is 1
            }

            // Update the input value
            input.value = currentValue.toFixed(type === "weight" ? 1 : 0);
          }
        });
        // Mark the parent container as having a listener
        (productContainer as any)._hasListener = true;
      }
    });
}

if (document.readyState !== "loading") {
  sendVisit();
  initMain();
} else {
  document.addEventListener("DOMContentLoaded", function () {
    sendVisit();
    initMain();
  });
}

// HTMX dynamic content handling
document.addEventListener("htmx:afterSwap", (event) => {
  initMain();
});
