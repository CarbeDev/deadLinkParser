package storage

import (
	"sync"
	"testing"
)

func TestCheckAndAddLink(t *testing.T) {
	// Reset the ReadedLinks map before testing
	ReadedLinks = make(map[string]bool)

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "Add new link",
			link:     "https://example.com",
			expected: true,
		},
		{
			name:     "Add duplicate link",
			link:     "https://example.com",
			expected: false,
		},
		{
			name:     "Add different link",
			link:     "https://example.org",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckAndAddLink(tt.link)
			if result != tt.expected {
				t.Errorf("CheckAndAddLink(%s) = %v, want %v", tt.link, result, tt.expected)
			}
		})
	}
}

func TestConcurrentCheckAndAddLink(t *testing.T) {
	// Reset the ReadedLinks map before testing
	ReadedLinks = make(map[string]bool)

	const numGoroutines = 100
	const link = "https://example.com"

	var wg sync.WaitGroup
	trueCount := 0
	var countMu sync.Mutex

	// Launch multiple goroutines to test concurrent access
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := CheckAndAddLink(link)
			if result {
				countMu.Lock()
				trueCount++
				countMu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Only one goroutine should have successfully added the link
	if trueCount != 1 {
		t.Errorf("Expected exactly one successful link addition, got %d", trueCount)
	}

	// Verify the link exists in the map
	if _, exists := ReadedLinks[link]; !exists {
		t.Error("Link should exist in the map")
	}
}
