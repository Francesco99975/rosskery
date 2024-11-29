package helpers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
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

func GetSessionOptions() *sessions.Options {
	sameSite := http.SameSiteDefaultMode
	maxAge := 86400 * 7
	if os.Getenv("GO_ENV") != "production" {
		sameSite = http.SameSiteNoneMode
		maxAge = 0
	}

	return &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   !(os.Getenv("ENVIRON") == "DEV"),
		Domain:   os.Getenv("DOMAIN"),
		SameSite: sameSite,
	}

}
