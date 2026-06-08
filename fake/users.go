package fake

import (
	"forum/database"
	"forum/models"
	"net/http"
)

func GetCurrentUser(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", false
	}

	user, ok := database.GetUserFromSession(cookie.Value)
	if !ok {
		return "", false
	}
	return user.Username, true
}

func GetCurrentUserFull(r *http.Request) (models.User, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return models.User{}, false
	}
	return database.GetUserFromSession(cookie.Value)
}
