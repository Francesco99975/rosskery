package controllers

import (
	"context"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

func Checkout(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.GetDefaultSite("Checkout", ctx)

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
			sessionID = uuid.NewV4().String()
			sess.Values["sessionID"] = sessionID
			err = sess.Save(c.Request(), c.Response())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not create session")
			}
		}
		cart, err := models.GetCart(ctx, sessionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get cart")
		}

		preview, err := cart.Preview(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get cart preview")
		}

		overbookedData, err := models.GetOverbooked()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get overbooked data")
		}

		csrfToken := c.Get("csrf").(string)
		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(views.Checkout(data, &preview, overbookedData, csrfToken, nonce))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
