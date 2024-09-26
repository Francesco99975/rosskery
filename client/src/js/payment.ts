function initPayment() {
  const paymentForm = document.getElementById("stripe-form");
  const errors = document.getElementById("error-messages");
  const publishableKeyElem = document.getElementById("pk") as HTMLInputElement;

  if (publishableKeyElem) {
    const stripe = (window as any).Stripe(publishableKeyElem.value);
    publishableKeyElem.remove();

    const csrfElem = document.getElementById("_csrf") as HTMLInputElement;
    if (csrfElem) {
      const csrfToken = csrfElem.value;

      fetch("/intent", {
        method: "POST",
        headers: {
          "X-CSRF-Token": csrfToken,
          "Content-Type": "application/json",
        },
      })
        .then((res) => res.json())
        .then((response) => {
          const elements = stripe.elements({
            clientSecret: response.clientSecret,
          });
          const paymentElement = elements.create("payment");
          paymentElement.mount("#payment-element");

          if (paymentForm) {
            paymentForm.addEventListener("submit", (event) => {
              event.preventDefault();
              const { error } = stripe.confirmPayment({
                elements,
                confirmParams: {
                  return_url: window.location.origin + "/orders/success",
                },
              });

              if (error && errors) {
                errors.innerHTML = error.message;
              }
            });
          }
        });
    }
  }
}

if (document.readyState !== "loading") {
  initPayment();
}

document.addEventListener("DOMContentLoaded", function () {
  initPayment();
});
