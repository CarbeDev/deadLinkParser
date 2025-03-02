package logger

import (
	"log"
	"net/http"
)

func LogResponseResult(response *http.Response) {
	log.Print(formatLinkStatusMessage(response))
}
