package models

import (
	"context"
	"os"
	"time"

	"github.com/Francesco99975/rosskery/internal/storage"
)

type SEO struct {
	Description string
	Keywords    string
}
type Site struct {
	AppName      string
	Title        string
	Metatags     SEO
	Year         int
	Message      string
	ContactEmail string
	ContactPhone string
	ContactIG    string
	ContactFB    string
}

func GetDefaultSite(title string, ctx context.Context) Site {
	val, err := storage.Valkey.Get(ctx, string(storage.Message)).Result()
	if err != nil {
		return Site{
			AppName:  "Rosskery",
			Title:    title,
			Metatags: SEO{Description: "Sweets store", Keywords: "shop,pastries,sweets,cookies,biscuits,buy,store,purchase"},
			Year:     time.Now().Year(),
		}
	}

	return Site{
		AppName:      "Rosskery",
		Title:        title,
		Metatags:     SEO{Description: "Sweets store", Keywords: "shop,pastries,sweets,cookies,biscuits,buy,store,purchase"},
		Year:         time.Now().Year(),
		Message:      val,
		ContactEmail: os.Getenv("CONTACT_EMAIL"),
		ContactPhone: os.Getenv("CONTACT_PHONE"),
		ContactIG:    os.Getenv("CONTACT_IG"),
		ContactFB:    os.Getenv("CONTACT_FB"),
	}
}
