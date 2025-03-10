package storage

import (
	"slices"
	"sync"
)

var ReadedLinks []string
var mu sync.Mutex

func CheckAndAddLink(link string) bool {
	mu.Lock()
	defer mu.Unlock()

	if slices.Contains(ReadedLinks, link) {
		return false
	}

	ReadedLinks = append(ReadedLinks, link)
	return true
}
