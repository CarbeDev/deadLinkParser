package parser

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestGetLinksFromResponse(t *testing.T) {
	type args struct {
		response *http.Response
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Should return all links in response body",
			args: args{
				response: &http.Response{
					Body: func() io.ReadCloser {
						content, err := os.ReadFile("parsing_test.html")
						if err != nil {
							panic(err)
						}
						return io.NopCloser(bytes.NewReader(content))
					}(),
				},
			},
			want: []string{
				"https://example.com", "/test", "https://nested.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := GetLinksFromResponse(tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinksFromResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
