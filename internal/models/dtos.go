package models

import (
	"fmt"
	"strings"
	"time"
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
	Image       string `json:"image"`
	Featured    bool   `json:"featured"`
	Published   bool   `json:"published"`
	CategoryId  string `json:"categoryId"`
	Weighed     bool   `json:"weighed"`
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

	if p.Image == "" {
		return fmt.Errorf("Product Image cannot be empty")
	}

	if p.CategoryId == "" {
		return fmt.Errorf("Product CategoryId cannot be empty")
	}

	if !CategoryExists(p.CategoryId) {
		return fmt.Errorf("Product CategoryId does not exist")
	}

	p.Name = strings.ToLower(p.Name)

	p.Image = strings.ToLower(p.Image)

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
	PurchasedItems []PurchasedItem `json:"purchasedItems"`
	Pickuptime     time.Time       `json:"pickuptime"`
	Method         PaymentMethod   `json:"method"`
	Fullname       string          `json:"fullname"`
	Email          string          `json:"email"`
	Address        string          `json:"address"`
	Phone          string          `json:"phone"`
}

func (o *OrderDto) Validate() error {
	if len(o.PurchasedItems) == 0 {
		return fmt.Errorf("OrderDto PurchasedItems cannot be empty")
	}

	for _, item := range o.PurchasedItems {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	if o.Pickuptime.IsZero() {
		return fmt.Errorf("OrderDto Pickuptime cannot be empty")
	}

	if o.Method == "" {
		return fmt.Errorf("OrderDto Method cannot be empty")
	}

	if o.Fullname == "" {
		return fmt.Errorf("OrderDto Fullname cannot be empty")
	}

	if o.Email == "" {
		return fmt.Errorf("OrderDto Email cannot be empty")
	}

	if o.Address == "" {
		return fmt.Errorf("OrderDto Address cannot be empty")
	}

	if o.Phone == "" {
		return fmt.Errorf("OrderDto Phone cannot be empty")
	}

	o.Email = strings.ToLower(o.Email)

	return nil
}

type JSONErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}
