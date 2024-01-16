package models

import (
	"fmt"
	"sync"
)

type Photo struct {
	Path   string
	Height int
	Width  int
}

type CachedPhotos struct {
	Photos     []Photo
	startIndex int
	endIndex   int
}

var instance *CachedPhotos
var once sync.Once

func GetCachePhotos() *CachedPhotos {
	once.Do(func() {
		instance = &CachedPhotos{Photos: make([]Photo, 0), startIndex: 0, endIndex: 0}
	})

	return instance
}

func (cache *CachedPhotos) Append(photos []Photo) {
	cache.Photos = append(cache.Photos, photos...)
}

func (cache *CachedPhotos) Take(amount int) ([]Photo, error) {
	cache.endIndex = cache.startIndex + amount

	if cache.endIndex >= len(cache.Photos) {
		cache.endIndex = len(cache.Photos)
	}

	if cache.startIndex >= cache.endIndex {
		return []Photo{}, fmt.Errorf("out of bounds")
	}

	payload := cache.Photos[cache.startIndex:cache.endIndex]
	cache.startIndex = cache.endIndex

	return payload, nil
}

