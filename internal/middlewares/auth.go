package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/Francesco99975/rosskery/internal/models"
)

func validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(tokenString *jwt.Token) (interface{}, error) {
		if _, ok := tokenString.Method.(*jwt.SigningMethodHMAC);!ok {
			return nil, fmt.Errorf("There was an error")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expired := claims["exp"].(int64)

		if time.Unix(expired, 0).After(time.Now()) {
			return "", fmt.Errorf("Token has expired")
		}

		return claims["sub"].(string), nil
	} else {
		return "", fmt.Errorf("Invalid Token")
	}
}

func IsAuthenticatedAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized. Cause -> Token not provided")
		}

		userid, err := validateToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized. Cause -> %v", err))
		}

		user, err := models.GetUserById(userid)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized. Cause -> %v", err))
		}

		basic, err := user.ToUser()
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized. Cause -> %v", err))
		}

		if basic.Role.Id != "3" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized. Cause -> User is not an admin")
		}

		c.Set("userid", userid)

		return next(c)
	}
}

func IsAuthenticatedModerator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized. Cause -> Token not provided")
		}

		userid, err := validateToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized. Cause -> %v", err))
		}

		user, err := models.GetUserById(userid)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized. Cause -> %v", err))
		}

		basic, err := user.ToUser()
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized. Cause -> %v", err))
		}

		if basic.Role.Id != "2" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized. Cause -> User is not a moderator")
		}

		c.Set("userid", userid)

		return next(c)
	}
}

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized. Cause -> Token not provided")
		}

		userid, err := validateToken(token)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized. Cause -> %v", err))
		}

		c.Set("userid", userid)

		return next(c)
	}
}
