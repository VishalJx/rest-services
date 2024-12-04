package models

import (
	"rest-services/config"
	"database/sql"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUserByEmail(email string) (*User, error) {
	row := config.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func CreateUser(email, hashedPassword string) error {
	_, err := config.DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, hashedPassword)
	return err
}
