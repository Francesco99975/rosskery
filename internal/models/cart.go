package models

import (
	"context"
	"encoding/json"

	"github.com/Francesco99975/rosskery/internal/storage"
	"github.com/labstack/gommon/log"
)

type Cart struct {
	Id    string         `json:"id"`
	Items map[string]int `json:"items"`
}

type CartPreview struct {
	Items []struct {
		Product  *Product `json:"product"`
		Quantity int      `json:"quantity"`
		Subtotal int      `json:"subtotal"`
	} `json:"items"`
	Total int `json:"subtotal"`
}

func (c *Cart) Preview(ctx context.Context) (CartPreview, error) {
	preview := CartPreview{
		Items: make([]struct {
			Product  *Product `json:"product"`
			Quantity int      `json:"quantity"`
			Subtotal int      `json:"subtotal"`
		}, 0, len(c.Items)),
		Total: 0,
	}

	for productId, quantity := range c.Items {
		product, err := GetProduct(productId)
		if err != nil {
			c.Clear(ctx)
			return CartPreview{}, err
		}

		var subtotal int
		var displayedQuantity int

		if product.Weighed {
			displayedQuantity = quantity
			subtotal = product.Price * quantity / 10
		} else {
			displayedQuantity = quantity
			subtotal = product.Price * quantity
		}

		preview.Items = append(preview.Items, struct {
			Product  *Product `json:"product"`
			Quantity int      `json:"quantity"`
			Subtotal int      `json:"subtotal"`
		}{product, displayedQuantity, subtotal})

		preview.Total += subtotal
	}

	return preview, nil
}

func (c *Cart) Purchases() ([]PurchasedItem, error) {
	purchases := make([]PurchasedItem, 0)

	for productId, quantity := range c.Items {
		product, err := GetProduct(productId)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, PurchasedItem{
			ProductId: product.Id,
			Quantity:  quantity,
		})
	}

	return purchases, nil
}

func (c *Cart) Save(ctx context.Context) error {
	cartData, err := json.Marshal(c.Items)
	if err != nil {
		return err
	}

	err = storage.Valkey.Set(ctx, c.Id, cartData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cart) AddItem(ctx context.Context, productId string, quantity int) error {
	c.Items[productId] += quantity

	return c.Save(ctx)
}

func (c *Cart) RemoveItem(ctx context.Context, productId string, quantity int) error {
	c.Items[productId] -= quantity

	if c.Items[productId] <= 0 {
		delete(c.Items, productId)
	}

	return c.Save(ctx)
}

func (c *Cart) DeleteItem(ctx context.Context, productId string) error {

	delete(c.Items, productId)

	return c.Save(ctx)
}

func (c *Cart) Clear(ctx context.Context) error {

	for productId := range c.Items {
		delete(c.Items, productId)
	}

	return c.Save(ctx)
}

func (c *Cart) Len() int {
	len := 0
	for _, quantity := range c.Items {
		len += quantity
	}

	return len
}

func GetCart(ctx context.Context, sessionID string) (*Cart, error) {
	cart := &Cart{Id: sessionID, Items: make(map[string]int, 0)}
	cartData, err := storage.Valkey.Get(ctx, sessionID).Result()

	if err != nil {
		if err := cart.Save(ctx); err != nil {
			log.Error(err)
			return nil, err
		}

		return cart, nil
	} else {

		if err := json.Unmarshal([]byte(cartData), &cart.Items); err != nil {

			return nil, err
		}

		return cart, nil
	}
}
