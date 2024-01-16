package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/views"

	"github.com/labstack/echo/v4"
)

func Gallery() echo.HandlerFunc {
	return func(c echo.Context) error {
		cache := models.GetCachePhotos()
		var images []models.Photo

		err := filepath.Walk("static/gallery", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if it's a regular file (not a directory)
			if !info.IsDir() {
				// Print the filename
				path := info.Name()
				file, err := os.Open(fmt.Sprintf("static/gallery/%s", path))
				if err != nil {
					return err
				}
				defer file.Close()

				// Decode the image
				img, err := webp.Decode(file)
				if err != nil {
					return err
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
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Did not find any images in gallery. Error: %s", err.Error()))
		}

		cache.Append(images)
		loadedImages, err := cache.Take(5)

		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Could not load images in gallery. Error: %s", err.Error()))
		}

		html, err := helpers.GeneratePage(views.Gallery(models.GetDefaultSite("Gallery"), loadedImages))

		if err != nil {
			echo.NewHTTPError(http.StatusBadRequest, "Could not parse page index")
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
			echo.NewHTTPError(http.StatusBadRequest, "Could not parse page index")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
