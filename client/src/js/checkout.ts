import flatpickr from "flatpickr";
import "flatpickr/dist/flatpickr.css";

interface HtmxResposeErrorEvent extends Event {
  detail: {
    request: XMLHttpRequest; // The XMLHttpRequest object associated with the request
    xhr: XMLHttpRequest; // The XMLHttpRequest used for the request
    swap: string; // The swap method used (e.g., "innerHTML", "outerHTML")
    status: number; // The HTTP status code of the response
    url: string; // The URL of the request
    // You can add other properties as needed
  };
}

function initCheckout() {
  const disabledDatesElem = document.getElementById("dd") as HTMLInputElement;
  var disabledDates: string[] = [];
  if (disabledDatesElem && disabledDatesElem.value.includes(",")) {
    disabledDates = disabledDatesElem.value.split(",");
  } else {
    disabledDates = [disabledDatesElem.value];
  }
  disabledDatesElem.remove();
  flatpickr("#pickuptime", {
    enableTime: true,
    minTime: "11:00",
    maxTime: "18:30",
    minDate: new Date().fp_incr(2), // 2 days from now
    disable: [
      function (date) {
        // Disable Weekends
        return (
          date.getDay() === 0 ||
          date.getDay() === 6 ||
          disabledDates.includes(date.toISOString().slice(0, 10))
        );
      },
    ],
    altInput: true,
    altFormat: "F j, Y, h:i K",
    dateFormat: "Y-m-d H:i",
  });

  const form = window.document.getElementById("checkout-form");
  if (form && !(form as any)._hasListener) {
    form.addEventListener("htmx:responseError", function (evt: any) {
      const event = evt as HtmxResposeErrorEvent;
      const errorBox = window.document.getElementById("errors");
      if (errorBox) {
        errorBox.innerHTML = event.detail.xhr.responseText;
        errorBox.style.display = "block";
      }
    });
    (form as any)._hasListener = true;
  }

  document.addEventListener("htmx:afterOnLoad", function (event) {
    const target = event.target as HTMLElement;
    if (target && target.id === "suggestions") {
      var suggestionsBox = document.getElementById("suggestions");
      if (suggestionsBox) {
        if (suggestionsBox.innerHTML.trim() !== "") {
          suggestionsBox.style.display = "block";
        } else {
          suggestionsBox.style.display = "none";
        }
      }
    }
  });

  const address = document.getElementById("address");
  const suggestionsBox = document.getElementById("suggestions");
  if (address && suggestionsBox && !(address as any)._hasListener) {
    address.addEventListener("keyup", function (event) {
      if (suggestionsBox.style.display === "none") {
        suggestionsBox.style.display = "block";
      }
    });
    (address as any)._hasListener = true;
  }

  document.addEventListener("click", function (event) {
    var suggestionsBox = document.getElementById("suggestions");
    const target = event.target as HTMLElement;
    if (suggestionsBox && target) {
      if (!suggestionsBox.contains(target) && target.id !== "address") {
        suggestionsBox.style.display = "none";
      }
    }
  });
}

if (document.readyState !== "loading") {
  initCheckout();
} else {
  document.addEventListener("DOMContentLoaded", initCheckout);
}
