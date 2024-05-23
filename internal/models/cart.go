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
