package data

import "sync"

var ReadedLinks = make(map[string]bool)
var mu sync.Mutex

func LinkHasBeenRead(link string) bool {
	mu.Lock()
	defer mu.Unlock()

	_, exists := ReadedLinks[link]
	return exists
}

func AddReadedLink(link string, accessible bool) {
	mu.Lock()
	defer mu.Unlock()

	ReadedLinks[link] = accessible
}
