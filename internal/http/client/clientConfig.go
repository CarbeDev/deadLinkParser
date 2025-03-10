package client

import (
	"net/http"
	"time"
)

var instance *http.Client

func HttpClient() *http.Client {
	if instance != nil {
		return instance
	}

	instance = &http.Client{Timeout: 10 * time.Second}

	return instance
}
