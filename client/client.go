package client

import "net/http"

func InternalRequest(path string, baseUrl string) (*http.Response, error) {
	fullUrl := constructFullUrl(path, baseUrl)

	return http.Get(fullUrl)
}

func ExternalRequest(url string) (*http.Response, error) {
	return http.Get(url)
}

func constructFullUrl(path string, baseUrl string) string {
	return baseUrl + path[1:]
}
