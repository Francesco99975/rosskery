package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/views/components"
	"github.com/labstack/echo/v4"
)

func AddressAutocomplete() echo.HandlerFunc {
	return func(c echo.Context) error {

		type Prediction struct {
			Description string `json:"description"`
		}

		type AutocompleteResponse struct {
			Predictions []Prediction `json:"predictions"`
		}

		query := c.QueryParam("address")
		if query == "" {
			return c.JSON(http.StatusOK, []string{})
		}

		const torontoLatitude = "43.65107"
		const torontoLongitude = "-79.347015"
		const radius = "100000" // 100 kilometers in meters

		apiURL := "https://maps.googleapis.com/maps/api/place/autocomplete/json?input=" + url.QueryEscape(query) + "&types=address&components=country:CA&location=" + torontoLatitude + "," + torontoLongitude + "&radius=" + radius + "&key=" + os.Getenv("GOOGLE_MAPS_API_KEY")
		resp, err := http.Get(apiURL)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get suggestions"})
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read response"})
		}

		var autocompleteResponse AutocompleteResponse
		if err := json.Unmarshal(body, &autocompleteResponse); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to unmarshal response"})
		}

		var suggestions []string
		for _, prediction := range autocompleteResponse.Predictions {
			suggestions = append(suggestions, prediction.Description)
		}

		html, err := helpers.GeneratePage(components.Suggestions(suggestions))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
