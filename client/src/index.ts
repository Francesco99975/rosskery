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
window.customElements.define("u-counter", UCounter);
window.conn = new WebSocket("ws://" + document.location.host + "/ws");
window.visited = false;

function addProductToDOM(html: string) {
  const productShop = document.getElementById("sp");
  const productFeutured = document.getElementById("fp");
  const productNewArrivals = document.getElementById("na");
  const productDiv = document.createElement("div");
  productDiv.innerHTML = html;
  if (productShop) {
    productShop.appendChild(productDiv);
  }

  if (productFeutured) {
    productFeutured.appendChild(productDiv);
  }

  if (productNewArrivals) {
    productNewArrivals.appendChild(productDiv);
  }
}

function updateProductInDOM(productId: string, html: string) {
  const productDiv = document.getElementById(`${productId}`);
  if (productDiv) {
    productDiv.innerHTML = html;
  }
}

function deleteProductFromDOM(productId: string) {
  const productDiv = document.getElementById(`${productId}`);
  if (productDiv) {
    productDiv.remove();
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

if (document.readyState !== "loading") {
  sendVisit();
}

document.addEventListener("DOMContentLoaded", function () {
  sendVisit();
});
