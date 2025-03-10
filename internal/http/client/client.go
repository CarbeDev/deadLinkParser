package client

import (
	"net/http"
)

type RealHTTPClient struct{}

func NewRealHTTPClient() *RealHTTPClient {
	return &RealHTTPClient{}
}

func (c *RealHTTPClient) InternalRequest(path string, baseUrl string) (*http.Response, error) {
	if path == baseUrl {
		return http.Get(path)
	}

	fullUrl := constructFullUrl(path, baseUrl)
	return HttpClient().Get(fullUrl)
}

func (c *RealHTTPClient) ExternalRequest(url string) (*http.Response, error) {
	return HttpClient().Get(url)
} 