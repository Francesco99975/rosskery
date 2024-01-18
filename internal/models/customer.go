package models

import "time"

type Customer struct {
	Id string `json:"id"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	Address string `json:"address"`
	Phone string `json:"phone"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
