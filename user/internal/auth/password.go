package user_auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (a *Auth) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error generate hash password: %w", err)
	}
	return string(hash), nil
}

func (a *Auth) ComparePassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("compare password : %w", err)
	}
	return nil
}
