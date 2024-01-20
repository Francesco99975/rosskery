package models

import "time"

type ProductDto struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Price int `json:"price"`
	Image string `json:"image"`
	Featured bool `json:"featured"`
	Published bool `json:"published"`
	CategoryId string `json:"categoryId"`
	Weighed bool `json:"weighed"`
}

type PurchasedItem struct {
	Product Product `json:"product"`
	Quantity int `json:"quantity"`
}

type OrderDto struct {
	PurchasedItems []PurchasedItem `json:"purchasedItems"`
	Pickuptime time.Time `json:"pickuptime"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	Address string `json:"address"`
	Phone string `json:"phone"`
}
