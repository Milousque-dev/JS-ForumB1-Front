package fake

import "net/http"

func GetCurrentUser(r *http.Request) (string, bool) {
	return "Boss", true
}