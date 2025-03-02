package responseUtils

import "net/http"

func IsSuccess(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 400
}
