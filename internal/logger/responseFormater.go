package logger

import (
	responseUtils "deadLinkParser/internal/http/utils"
	"fmt"
	"net/http"
)

func formatLinkStatusMessage(response *http.Response) string {
	return fmt.Sprintf("Link : %v | Status : %v %s", response.Request.URL, response.Status, string(resultEmoji(response)))
}

func resultEmoji(response *http.Response) rune {
	if responseUtils.IsSuccess(response) {
		return '✅'
	}

	return '❌'
}
