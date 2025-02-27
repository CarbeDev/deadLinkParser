package client

import (
	"net/http"
)

func InternalRequest(path string, baseUrl string) (*http.Response, error) {
	if path == baseUrl {
		return http.Get(path)
	}

	fullUrl := constructFullUrl(path, baseUrl)

	return HttpClient().Get(fullUrl)
}

func ExternalRequest(url string) (*http.Response, error) {
	return HttpClient().Get(url)
}

func constructFullUrl(path string, baseUrl string) string {
	return baseUrl + path[1:]
}
