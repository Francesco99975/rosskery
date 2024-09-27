package middlewares

import (
	"io"
	"net/http"

	"github.com/google/brotli/go/cbrotli"
	"github.com/labstack/echo/v4"
)

// BrotliMiddleware compresses response using Brotli
func BrotliMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if client accepts Brotli encoding
		if c.Request().Header.Get(echo.HeaderAcceptEncoding) != "br" {
			return next(c) // Skip compression
		}

		// Capture the response
		res := c.Response()
		res.Header().Set(echo.HeaderContentEncoding, "br")

		// Create a Brotli writer
		bw := cbrotli.NewWriter(res.Writer, cbrotli.WriterOptions{Quality: 5})
		defer bw.Close()

		// Wrap the original writer with Brotli writer
		res.Writer = &brotliResponseWriter{Writer: bw, ResponseWriter: res.Writer}

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
