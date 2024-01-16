package models

import "time"

type Visit struct {
	Id string
	Ip string
	Views int
	Duration int
	Sauce string
	Agent string
	Date time.Time
}




