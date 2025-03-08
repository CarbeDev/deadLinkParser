package crawler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
)

// MockHTTPClient implements HTTPClient interface for testing
type MockHTTPClient struct {
	internalResponses map[string]*http.Response
	externalResponses map[string]*http.Response
	internalErrors    map[string]error
	externalErrors    map[string]error
}

func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		internalResponses: make(map[string]*http.Response),
		externalResponses: make(map[string]*http.Response),
		internalErrors:    make(map[string]error),
		externalErrors:    make(map[string]error),
	}
}

func (m *MockHTTPClient) InternalRequest(link, baseURL string) (*http.Response, error) {
	if err, exists := m.internalErrors[link]; exists && err != nil {
		return nil, err
	}
	if resp, exists := m.internalResponses[link]; exists {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString("")),
		Request:    &http.Request{URL: nil},
	}, nil
}

func (m *MockHTTPClient) ExternalRequest(link string) (*http.Response, error) {
	if err, exists := m.externalErrors[link]; exists && err != nil {
		return nil, err
	}
	if resp, exists := m.externalResponses[link]; exists {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString("")),
		Request:    &http.Request{URL: nil},
	}, nil
}

func createMockResponse(statusCode int, body string, urlPath string) *http.Response {
	parsedURL, _ := url.Parse(urlPath)
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Request: &http.Request{
			URL: parsedURL,
		},
	}
}

func TestCrawler_isInternal(t *testing.T) {
	mockClient := NewMockHTTPClient()
	crawler := NewCrawler(mockClient)
	crawler.baseUrl = "https://example.com"

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "Relative path should be internal",
			link:     "/some/path",
			expected: true,
		},
		{
			name:     "Same domain should be internal",
			link:     "https://example.com/path",
			expected: true,
		},
		{
			name:     "Different domain should be external",
			link:     "https://other-domain.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crawler.isInternal(tt.link)
			if result != tt.expected {
				t.Errorf("isInternal(%s) = %v; want %v", tt.link, result, tt.expected)
			}
		})
	}
}

func TestCrawler_makeRequest(t *testing.T) {
	mockClient := NewMockHTTPClient()
	crawler := NewCrawler(mockClient)
	crawler.baseUrl = "https://example.com"

	internalURL := "/internal-path"
	externalURL := "https://external-domain.com"

	expectedInternalResp := createMockResponse(http.StatusOK, "internal content", internalURL)
	expectedExternalResp := createMockResponse(http.StatusOK, "external content", externalURL)

	mockClient.internalResponses[internalURL] = expectedInternalResp
	mockClient.externalResponses[externalURL] = expectedExternalResp

	tests := []struct {
		name          string
		link          string
		expectedResp  *http.Response
		expectedError error
	}{
		{
			name:          "Internal request success",
			link:          internalURL,
			expectedResp:  expectedInternalResp,
			expectedError: nil,
		},
		{
			name:          "External request success",
			link:          externalURL,
			expectedResp:  expectedExternalResp,
			expectedError: nil,
		},
		{
			name:          "Internal request error",
			link:          "/error-path",
			expectedResp:  nil,
			expectedError: errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError != nil {
				mockClient.internalErrors[tt.link] = tt.expectedError
			}

			resp, err := crawler.makeRequest(tt.link)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if resp != tt.expectedResp {
					t.Errorf("Expected response %v, got %v", tt.expectedResp, resp)
				}
			}
		})
	}
}
