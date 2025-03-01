package data

import "sync"

var ReadedLinks = make(map[string]bool)
var mu sync.Mutex

func CheckAndAddLink(link string) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := ReadedLinks[link]; exists {
		return false
	}

	ReadedLinks[link] = false
	return true
}
