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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expired := int64(claims["exp"].(float64))

		if time.Now().Unix() > expired {
			return "", fmt.Errorf("Token expired at %v, now is %v", time.Unix(expired, 0), time.Now())
		}

		return claims["sub"].(string), nil
	} else {
		return "", fmt.Errorf("Invalid Token")
	}
}



func IsAuthenticatedAdmin() echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")

			if token == "" {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{ Code: http.StatusUnauthorized, Message: "Unauthorized. Cause -> Token not provided"})
			}

			userid, err := validateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{ Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized. Cause -> %v", err), Errors: []string{err.Error()}})
			}

			user, err := models.GetUserById(userid)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{ Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized. Cause -> %v", err), Errors: []string{err.Error()}})
			}

			basic, err := user.ToUser()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{ Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized. Cause -> %v", err), Errors: []string{err.Error()}})
			}

			if basic.Role.Id != "3" {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{ Code: http.StatusUnauthorized, Message: "Unauthorized. Cause -> User not an Admin"})
			}

			c.Set("userid", userid)

			return next(c)
		}
	}
}



func IsAuthenticatedModerator() echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
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
}



func IsAuthenticated() echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
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
}


