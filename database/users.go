package database

import (
	"database/sql"
	"errors"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmailTaken = errors.New("email déjà utilisé")
var ErrUserNotFound = errors.New("utilisateur introuvable")

func CreateUser(email, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)`,
		email, username, string(hash),
	)
	if err != nil && err.Error() == "UNIQUE constraint failed: users.email" {
		return ErrEmailTaken
	}
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	row := DB.QueryRow(
		`SELECT id, email, username, password_hash FROM users WHERE email = ?`,
		email,
	)

	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	return u, err
}
