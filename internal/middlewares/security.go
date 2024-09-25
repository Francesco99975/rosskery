package middlewares

import (
	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/labstack/echo/v4"
)

func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nonce, err := helpers.GenerateNonce()
			if err != nil {
				return err
			}

			// Store nonce in the request context so it can be accessed in your templates
			c.Set("nonce", nonce)

			// Set security headers
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'nonce-"+nonce+"' 'strict-dynamic' https://js.stripe.com; connect-src 'self' https://api.stripe.com; style-src 'self'; frame-src https://js.stripe.com")
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "no-referrer")
			c.Response().Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")
			return next(c)
		}
	}
}

func SecurityHeadersDev() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			nonce, err := helpers.GenerateNonce()
			if err != nil {
				return err
			}

			// Store nonce in the request context so it can be accessed in your templates
			c.Set("nonce", nonce)

			// Set security headers
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'nonce-"+nonce+"' 'strict-dynamic' 'unsafe-eval' https://js.stripe.com https://cdn.jsdelivr.net/npm/flatpickr/dist/flatpickr.min.css https://cdn.jsdelivr.net/npm/flatpickr; connect-src 'self' https://api.stripe.com https://cdn.jsdelivr.net/npm/flatpickr/dist/flatpickr.min.css https://cdn.jsdelivr.net/npm/flatpickr; style-src 'self' 'unsafe-inline'; frame-src https://js.stripe.com")
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "SAMEORIGIN") // Allow iframes for easier testing
			// No Strict-Transport-Security for local development
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "no-referrer-when-downgrade") // Less strict for development
			return next(c)
		}
	}
}
