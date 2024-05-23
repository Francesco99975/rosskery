package controllers

import (
	"context"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views/components"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

func GetCartItems(ctx context.Context) echo.HandlerFunc {
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
			sessionID = uuid.NewV4().String()
			sess.Values["sessionID"] = sessionID
			err = sess.Save(c.Request(), c.Response())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not create session")
			}
		}
		cart, err := models.GetCart(ctx, sessionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		html, err := helpers.GeneratePage(components.Badge(cart.Len()))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page index")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}

func AddToCart(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		productId := c.Param("id")

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

		if err := cart.AddItem(ctx, productId, 1); err != nil {
			return err
		}

		html, err := helpers.GeneratePage(components.Badge(cart.Len()))
		if err != nil {
			return err
		}

		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create session")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
