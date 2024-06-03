package api

import (
	"context"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

func CreatePaymentIntent(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			// Secure:   true,
			// Domain:   "",
			// SameSite: http.SameSiteDefaultMode,
		}
		sessionID, ok := sess.Values["sessionID"].(string)
		if !ok || sessionID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not get session id")
		}
		cart, err := models.GetCart(ctx, sessionID)
		if err != nil {
			return err
		}

		preview, err := cart.Preview()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get cart preview")
		}
		amountToPay := preview.Total
		if amountToPay <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "No items in cart")
		}

		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(int64(amountToPay)),
			Currency: stripe.String(string(stripe.CurrencyCAD)),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
		}

		pi, err := paymentintent.New(params)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error creating payment intent")
		}

		return c.JSON(http.StatusAccepted, struct {
			ClientSecret string `json:"clientSecret"`
		}{ClientSecret: pi.ClientSecret})
	}
}
