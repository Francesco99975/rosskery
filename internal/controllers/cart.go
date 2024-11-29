package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views/components"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
)

func GetCartItems(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = helpers.GetSessionOptions()
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
			log.Errorf("Could Not get cart preview <- %w", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get cart preview")
		}

		csrfToken := c.Get("csrf").(string)
		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(components.Badge(cart.Len(), &preview, false, csrfToken, nonce))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page index")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}

func AddToCart(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		productId := c.Param("id")

		var quantity int

		qty, err := strconv.Atoi(c.FormValue(fmt.Sprintf("quantityInput-%s", productId)))
		if err != nil {
			weightStr := c.FormValue(fmt.Sprintf("weightInput-%s", productId))
			weight, err := strconv.ParseFloat(weightStr, 64)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Could not get weight")
			}

			quantity = int(weight * 10)
		}

		if !(quantity > 0) {
			quantity = qty
		}

		if !(quantity > 0) {
			return echo.NewHTTPError(http.StatusBadRequest, "Quantity must be greater than 0")
		}

		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = helpers.GetSessionOptions()

		sessionID, ok := sess.Values["sessionID"].(string)
		if !ok || sessionID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not get session id")
		}

		cart, err := models.GetCart(ctx, sessionID)
		if err != nil {
			return err
		}

		if err := cart.AddItem(ctx, productId, quantity); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could add to cart")
		}

		preview, err := cart.Preview(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get cart preview")
		}

		csrfToken := c.Get("csrf").(string)
		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(components.Badge(cart.Len(), &preview, false, csrfToken, nonce))
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

func RemoveOneFromCart(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		productId := c.Param("id")

		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = helpers.GetSessionOptions()

		sessionID, ok := sess.Values["sessionID"].(string)
		if !ok || sessionID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not get session id")
		}

		cart, err := models.GetCart(ctx, sessionID)
		if err != nil {
			return err
		}

		if err := cart.RemoveItem(ctx, productId, 1); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could remove from cart")
		}

		preview, err := cart.Preview(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get cart preview")
		}

		csrfToken := c.Get("csrf").(string)
		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(components.Badge(cart.Len(), &preview, true, csrfToken, nonce))
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

func RemoveItemFromCart(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		productId := c.Param("id")

		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = helpers.GetSessionOptions()

		sessionID, ok := sess.Values["sessionID"].(string)
		if !ok || sessionID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not get session id")
		}

		cart, err := models.GetCart(ctx, sessionID)
		if err != nil {
			return err
		}

		if err := cart.DeleteItem(ctx, productId); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could delete product from cart")
		}

		preview, err := cart.Preview(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get cart preview")
		}

		csrfToken := c.Get("csrf").(string)
		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(components.Badge(cart.Len(), &preview, true, csrfToken, nonce))
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

func ClearCart(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = helpers.GetSessionOptions()

		sessionID, ok := sess.Values["sessionID"].(string)
		if !ok || sessionID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not get session id")
		}

		cart, err := models.GetCart(ctx, sessionID)
		if err != nil {
			return err
		}

		if err := cart.Clear(ctx); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could clear cart")
		}

		preview, err := cart.Preview(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get cart preview")
		}

		csrfToken := c.Get("csrf").(string)
		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(components.Badge(cart.Len(), &preview, true, csrfToken, nonce))
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
