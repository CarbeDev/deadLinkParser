package logger

import (
	"net/http"
	"net/url"
	"testing"
)

func TestFormatLinkStatusMessage(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		want     string
	}{
		{
			name: "success response",
			response: &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Request: &http.Request{
					URL: &url.URL{
						Scheme: "https",
						Host:   "example.com",
						Path:   "/test",
					},
				},
			},
			want: "Link : https://example.com/test | Status : 200 OK ✅",
		},
		{
			name: "error response",
			response: &http.Response{
				Status:     "404 Not Found",
				StatusCode: 404,
				Request: &http.Request{
					URL: &url.URL{
						Scheme: "https",
						Host:   "example.com",
						Path:   "/notfound",
					},
				},
			},
			want: "Link : https://example.com/notfound | Status : 404 Not Found ❌",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatLinkStatusMessage(tt.response)
			if got != tt.want {
				t.Errorf("formatLinkStatusMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultEmoji(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		want     rune
	}{
		{
			name: "success response",
			response: &http.Response{
				StatusCode: 200,
			},
			want: '✅',
		},
		{
			name: "error response",
			response: &http.Response{
				StatusCode: 404,
			},
			want: '❌',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resultEmoji(tt.response)
			if got != tt.want {
				t.Errorf("resultEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}
