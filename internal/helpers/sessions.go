package helpers

import (
	"fmt"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetSessionId(c echo.Context) (string, error) {
	sess, err := session.Get("session", c)

	if err != nil {
		return "", err
	}

	sessionID, ok := sess.Values["sessionID"].(string)
	if !ok || sessionID == "" {
		return "", fmt.Errorf("sessionID not found in session")
	}

	return sessionID, nil
}
