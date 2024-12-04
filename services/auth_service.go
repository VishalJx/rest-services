package services

import (
	"rest-services/models"
	"rest-services/utils"
	"golang.org/x/crypto/bcrypt"
)

// SignUp creates a new user
func SignUp(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return models.CreateUser(email, string(hashedPassword))
}

// SignIn authenticates a user and returns access and refresh tokens
func SignIn(email, password string) (string, string, error) {
	user, err := models.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := utils.GenerateToken(user.Email)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
