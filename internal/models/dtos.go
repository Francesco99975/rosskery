package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type CategoryDto struct {
	Category string `json:"category"`
}

func (c *CategoryDto) Validate() error {
	if c.Category == "" {
		return fmt.Errorf("Category cannot be empty")
	}

	c.Category = strings.ToLower(c.Category)

	return nil
}

type ProductDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Featured    bool   `json:"featured"`
	Published   bool   `json:"published"`
	CategoryId  string `json:"categoryId"`
	Weighed     bool   `json:"weighed"`
	Lv          int    `json:"lv"`
}

func (p *ProductDto) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("Product Name cannot be empty")
	}

	if p.Description == "" {
		return fmt.Errorf("Product Description cannot be empty")
	}

	if p.Price < 0 {
		return fmt.Errorf("Product Price cannot be negative")
	}

	if p.CategoryId == "" {
		return fmt.Errorf("Product CategoryId cannot be empty")
	}

	if !CategoryExists(p.CategoryId) {
		return fmt.Errorf("Product CategoryId does not exist")
	}

	p.Name = strings.ToLower(p.Name)

	if p.Lv < 1 {
		return fmt.Errorf("Product Lv cannot be negative or zero")
	}

	return nil
}

type PurchasedItem struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func (p *PurchasedItem) Validate() error {
	if p.ProductId == "" {
		return fmt.Errorf("PurchasedItem ProductId cannot be empty")
	}

	if !ProductExists(p.ProductId) {
		return fmt.Errorf("PurchasedItem ProductId does not exist")
	}

	if p.Quantity < 1 {
		return fmt.Errorf("PurchasedItem Quantity cannot be negative or zero")
	}

	return nil
}

type OrderDto struct {
	Pickuptime time.Time     `json:"pickuptime"`
	Method     PaymentMethod `json:"method"`
	Fullname   string        `json:"fullname"`
	Email      string        `json:"email"`
	Address    string        `json:"address"`
	Phone      string        `json:"phone"`
}

func (o *OrderDto) Validate() error {

	if o.Pickuptime.IsZero() {
		return fmt.Errorf("pickuptime cannot be empty")
	}

	if o.Method == "" {
		return fmt.Errorf("method cannot be empty")
	}

	if o.Fullname == "" {
		return fmt.Errorf("fullname cannot be empty")
	}

	if o.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if o.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if o.Phone == "" {
		return fmt.Errorf("phone cannot be empty")
	}

	//Validate Email

	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegexPattern)

	if !re.MatchString(o.Email) {
		return fmt.Errorf("email is not a valid email address")
	}

	o.Email = strings.ToLower(o.Email)

	// Validate phone number
	phoneRegex := regexp.MustCompile(`^\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`)
	phoneCleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(o.Phone, "")
	if !phoneRegex.MatchString(phoneCleaned) {
		return fmt.Errorf("phone is not a valid phone number")
	}

	// Format phone number
	formattedPhone := fmt.Sprintf("(%s) %s-%s", phoneCleaned[:3], phoneCleaned[3:6], phoneCleaned[6:])

	o.Phone = formattedPhone

	// Validate address

	addressRegex := regexp.MustCompile(`^\d+\s[A-Za-z]+(?:\s[A-Za-z]+)*,?\s*[A-Za-z]+(?:\s[A-Za-z]+)*,?\s*(?:[A-Za-z]+\s*)?,?\s*(\d{5}|\b[ABCEGHJKLMNPRSTVXY]{1}\d{1}[A-Z]{1}\s?\d{1}[A-Z]{1}\d{1}\b)?(?:\s[A-Za-z]+)?$`)
	if !addressRegex.MatchString(o.Address) {
		return fmt.Errorf("address is not a valid address")
	}

	formattedAddress := strings.TrimSpace(o.Address)
	formattedAddress = regexp.MustCompile(`\s{2,}`).ReplaceAllString(formattedAddress, " ") // Replace multiple spaces with a single space
	formattedAddress = regexp.MustCompile(`,\s*`).ReplaceAllString(formattedAddress, ", ")  // Ensure a single space after commas

	type GeocodeResponse struct {
		Results []struct {
			FormattedAddress string `json:"formatted_address"`
			Geometry         struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
		Status string `json:"status"`
	}

	apiURL := "https://maps.googleapis.com/maps/api/geocode/json?address=" + url.QueryEscape(o.Address) + "&key=" + os.Getenv("GOOGLE_MAPS_API_KEY")

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Errorf("Failed to get geocode response: %v", err)
		return fmt.Errorf("address is not a valid address")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read geocode response")
		return fmt.Errorf("address is not a valid address")
	}

	var geocodeResponse GeocodeResponse
	if err := json.Unmarshal(body, &geocodeResponse); err != nil {
		log.Errorf("Failed to unmarshal geocode response: %v", err)
		return fmt.Errorf("address is not a valid address")
	}

	if geocodeResponse.Status != "OK" || len(geocodeResponse.Results) == 0 {
		log.Errorf("Failed to geocode address: %v", geocodeResponse.Status)
		return fmt.Errorf("address is not a valid address")
	}

	o.Address = formattedAddress

	return nil
}

type JSONErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}
