package models

import "time"

type Purchase struct {
	Id string `json:"id"`
	Product Product `json:"product"`
	Quantity int `json:"quantity"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type Order struct {
	Id string `json:"id"`
	Customer Customer `json:"customer"`
	Purchases []Purchase `json:"purchases"`
	Pickuptime time.Time `json:"pickuptime"`
	Fulfilled bool `json:"fulfilled"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
