package logger

import (
	"fmt"
	"net/http"
)

func formatLinkStatusMessage(response *http.Response) string {
	return fmt.Sprintf("Link : %v | Status : %v %s", response.Request.URL, response.Status, string(resultEmoji(response)))
}

func resultEmoji(response *http.Response) rune {
	if isError(response) {
		return '❌'
	}

	return '✅'
}

func isError(response *http.Response) bool {
	return response.StatusCode >= 400
}
