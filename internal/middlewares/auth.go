package middlewares

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
)

func IsAuthenticatedAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: "Unauthorized. Cause -> Token not provided"})
			}

			token := cookie.Value

			if token == "" {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: "Unauthorized. Cause -> Token not provided"})
			}

			userid, err := helpers.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized. Cause -> %v", err), Errors: []string{err.Error()}})
			}

			user, err := models.GetUserById(userid)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized. Cause -> %v", err), Errors: []string{err.Error()}})
			}

			basic, err := user.ToUser()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized. Cause -> %v", err), Errors: []string{err.Error()}})
			}

			if basic.Role.Id != "3" {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: "Unauthorized. Cause -> User not an Admin"})
			}

			c.Set("userid", userid)

			return next(c)
		}
	}
}

func IsAuthenticatedModerator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: "Unauthorized. Cause -> Token not provided"})
			}

			token := cookie.Value

			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized. Cause -> Token not provided")
			}

			userid, err := helpers.ValidateToken(token)
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
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: "Unauthorized. Cause -> Token not provided"})
			}

			token := cookie.Value

			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized. Cause -> Token not provided")
			}

			userid, err := helpers.ValidateToken(token)

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized. Cause -> %v", err))
			}

			c.Set("userid", userid)

			return next(c)
		}
	}
}
