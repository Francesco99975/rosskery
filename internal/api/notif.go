package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetNotifConfig() echo.HandlerFunc {
	return func(c echo.Context) error {
		notifServerUrl := os.Getenv("GOTIFY_SERVER")
		notifToken := os.Getenv("GOTIFY_TOKEN")

		//replace https:// with wss://
		notifServerUrl = strings.Replace(notifServerUrl, "https://", "wss://", 1)

		return c.JSON(http.StatusOK, map[string]string{
			"server": notifServerUrl,
			"token":  notifToken,
		})
	}
}
