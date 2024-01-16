package models

import "time"

type SEO struct {
	Description string
	Keywords    string
}
type Site struct {
	AppName  string
	Title    string
	Metatags SEO
	Year     int
}

func GetDefaultSite(title string) Site {
	return Site{
		AppName:  "Rosskery",
		Title:    title,
		Metatags: SEO{Description: "Sweets store", Keywords: "shop,pastries,sweets,cookies,biscuits,buy,store,purchase"},
		Year:     time.Now().Year(),
	}
}
