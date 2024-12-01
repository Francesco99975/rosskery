package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views"
	"github.com/chai2010/webp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Gallery(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		cache := models.GetCachePhotos()
		var images []models.Photo

		err := filepath.Walk("static/gallery", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Errorf("Error accessing path %s: %v\n", path, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Error accessing path")
			}

			// Check if it's a regular file (not a directory)
			if !info.IsDir() {
				// Print the filename
				path := info.Name()
				file, err := os.Open(fmt.Sprintf("static/gallery/%s", path))
				if err != nil {
					log.Errorf("Error opening file %s: %v\n", path, err)
					return echo.NewHTTPError(http.StatusInternalServerError, "Error opening file")
				}
				defer file.Close()

				// Decode the image
				img, err := webp.Decode(file)
				if err != nil {
					log.Errorf("Error decoding image %s: %w\n", path, err)
					return echo.NewHTTPError(http.StatusInternalServerError, "Error decoding image")
				}

				// Get the dimensions of the image
				width := img.Bounds().Dx()
				height := img.Bounds().Dy()

				if width < height {
					width = 1080
					height = 1920
				} else if width > height {
					width = 1920
					height = 1080
				} else {
					width = 1440
					height = 1440
				}

				images = append(images, models.Photo{Path: path, Height: height, Width: width})

			}

			return nil
		})

		if err != nil {
			log.Errorf("Error walking through gallery: %v\n", err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Did not find any images in gallery. Error: %s", err.Error()))
		}

		cache.Append(images)
		loadedImages, err := cache.Take(5)

		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Could not load images in gallery. Error: %s", err.Error()))
		}

		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(views.Gallery(models.GetDefaultSite("Gallery", ctx), loadedImages, nonce))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page index")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}

func Photos() echo.HandlerFunc {
	return func(c echo.Context) error {
		cache := models.GetCachePhotos()
		loadedImages, err := cache.Take(5)

		if err != nil {
			return c.NoContent(http.StatusOK)
		}

		html, err := helpers.GeneratePage(views.Photos(loadedImages))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page index")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
