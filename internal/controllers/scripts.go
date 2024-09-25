package controllers

import (
	"fmt"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/views/components"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Scripts() echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(components.Script(fmt.Sprintf("/assets/js/%s", key), nonce))
		if err != nil {
			log.Error(err)
		}
		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
