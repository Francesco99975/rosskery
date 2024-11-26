package middlewares

import (
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/labstack/echo/v4"
)

// BrotliMiddleware compresses response using Brotli
func BrotliMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if the client accepts Brotli encoding
		if !strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), "br") {
			return next(c) // Skip Brotli compression if not accepted
		}

		// Capture the response
		res := c.Response()
		res.Header().Set(echo.HeaderContentEncoding, "br")

		// Create a Brotli writer
		writer := brotli.NewWriter(res.Writer)
		defer writer.Close()

		// Wrap the original writer with the Brotli writer
		res.Writer = &brotliResponseWriter{Writer: writer, ResponseWriter: res.Writer}

		return next(c)
	}
}

// brotliResponseWriter wraps http.ResponseWriter with Brotli writer
type brotliResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write overrides the default ResponseWriter Write method
func (w *brotliResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
