package cache

import "sync"

type LatestSuggestionsCache struct {
	Suggestions []string
	Lock        sync.Mutex
}

func NewSuggestionsCache() *LatestSuggestionsCache {
	return &LatestSuggestionsCache{
		Suggestions: []string{},
		Lock:        sync.Mutex{},
	}
}

var suggestionsCache = NewSuggestionsCache()

func GetSuggestions() []string {
	suggestionsCache.Lock.Lock()
	defer suggestionsCache.Lock.Unlock()
	return suggestionsCache.Suggestions
}

func SetSuggestions(suggestions []string) {
	suggestionsCache.Lock.Lock()
	defer suggestionsCache.Lock.Unlock()
	suggestionsCache.Suggestions = suggestions
}
