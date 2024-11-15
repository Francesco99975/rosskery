package api

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views/components"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
)

func AddProduct(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {

		parsedPrice, err := strconv.Atoi(c.FormValue("price"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing price: %v", err), Errors: []string{err.Error()}})
		}

		parsedLv, err := strconv.Atoi(c.FormValue("lv"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing labor value: %v", err), Errors: []string{err.Error()}})
		}

		payload := models.ProductDto{
			Name:        c.FormValue("name"),
			Description: c.FormValue("description"),
			Price:       parsedPrice,
			CategoryId:  c.FormValue("category_id"),
			Weighed:     c.FormValue("weighed") == "true",
			Lv:          parsedLv,
		}

		if err := payload.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error product not valid: %v", err), Errors: []string{err.Error()}})
		}

		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		uploadedFiles := form.File["image"]
		if len(uploadedFiles) != 1 {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error uploading files: %v", err), Errors: nil})
		}

		file := uploadedFiles[0]

		id := uuid.NewV4().String()

		products, err := models.CreateProduct(id, payload.Name, payload.Description, payload.Price, file, payload.CategoryId, payload.Weighed, payload.Lv)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error creating product: %v", err), Errors: []string{err.Error()}})
		}

		newProduct, err := models.GetProduct(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching new product: %v", err), Errors: []string{err.Error()}})
		}

		csrfToken := c.Request().Header.Get("X-CSRF-Token")

		html, err := helpers.GeneratePage(components.ProductItem(*newProduct, csrfToken))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing html data: %v", err), Errors: []string{err.Error()}})
		}

		htmlData := models.HtmlData{Id: id, Html: string(html)}
		rawHtmlData, err := json.Marshal(htmlData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing html data: %v", err), Errors: []string{err.Error()}})
		}

		cm.BroadcastEvent(models.Event{Type: models.EventNewProduct, Payload: rawHtmlData})

		return c.JSON(http.StatusCreated, products)
	}
}

func Products() echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := models.GetProducts()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching products: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, products)
	}
}

func Product() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Product
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing data for product: %v", err), Errors: []string{err.Error()}})
		}

		product, err := models.GetProduct(payload.Id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching product: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, product)
	}
}

func UpdateProduct(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		parsedPrice, err := strconv.Atoi(c.FormValue("price"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing price: %v", err), Errors: []string{err.Error()}})
		}

		parsedLv, err := strconv.Atoi(c.FormValue("lv"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing labor value: %v", err), Errors: []string{err.Error()}})
		}

		payload := models.ProductDto{
			Name:        c.FormValue("name"),
			Description: c.FormValue("description"),
			Price:       parsedPrice,
			CategoryId:  c.FormValue("category_id"),
			Weighed:     c.FormValue("weighed") == "true",
			Lv:          parsedLv,
			Published:   c.FormValue("published") == "true",
			Featured:    c.FormValue("featured") == "true",
		}

		if err := payload.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error product not valid: %v", err), Errors: []string{err.Error()}})
		}

		image := c.FormValue("file")
		var file *multipart.FileHeader
		log.Info("image: ", image)

		if !strings.Contains(image, "/assets/") {
			file, err = c.FormFile("image")
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error uploading files: %v", err), Errors: nil})
			}
		} else {
			file = nil
		}

		product, err := models.GetProduct(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Product not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		products, err := product.Update(payload.Name, payload.Description, payload.Price, file, payload.Featured, payload.Published, payload.CategoryId, payload.Weighed, payload.Lv)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while updating product: %v", err), Errors: []string{err.Error()}})
		}

		updatedProduct, err := models.GetProduct(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Updated Product not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		csrfToken := c.Get("csrf").(string)

		html, err := helpers.GeneratePage(components.ProductItem(*updatedProduct, csrfToken))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
		}

		htmlData := models.HtmlData{Id: id, Html: string(html)}
		rawHtmlData, err := json.Marshal(htmlData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing html data: %v", err), Errors: []string{err.Error()}})
		}

		cm.BroadcastEvent(models.Event{Type: models.EventUpdateVisitsAdmin, Payload: rawHtmlData})

		return c.JSON(http.StatusOK, products)
	}
}

func DeleteProduct(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		product, err := models.GetProduct(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Product not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		products, err := product.Delete()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while deleting product: %v", err), Errors: []string{err.Error()}})
		}

		rawId, err := json.Marshal(struct{ Id string }{Id: id})
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing product removal: %v", err), Errors: []string{err.Error()}})
		}

		cm.BroadcastEvent(models.Event{Type: models.EventRemoveProduct, Payload: rawId})

		return c.JSON(http.StatusOK, products)
	}
}
