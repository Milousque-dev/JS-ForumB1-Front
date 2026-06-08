package database

import (
	"database/sql"
	"forum/models"
	"time"

	"github.com/google/uuid"
)

func CreateSession(userID int) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err := DB.Exec(
		`INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`,
		sessionID, userID, expiresAt,
	)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func GetUserFromSession(sessionID string) (models.User, bool) {
	row := DB.QueryRow(`
		SELECT u.id, u.email, u.username
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = ? AND s.expires_at > CURRENT_TIMESTAMP
	`, sessionID)

	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Username)
	if err == sql.ErrNoRows {
		return models.User{}, false
	}
	if err != nil {
		return models.User{}, false
	}
	return u, true
}

func DeleteSession(sessionID string) {
	DB.Exec(`DELETE FROM sessions WHERE id = ?`, sessionID)
}
