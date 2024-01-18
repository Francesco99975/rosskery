package models

import "time"

type Product struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price int `json:"price"`
	Image string `json:"image"`
	Featured bool `json:"featured"`
	Published bool `json:"published"`
	Category Category `json:"category"`
	Weighed bool `json:"weighed"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
