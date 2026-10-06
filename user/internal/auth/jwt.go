package user_auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *Auth) SignAccess(userID, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"sid": sessionID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(a.accessTTL).Unix(),
		"iss": "url-shortener-user",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return signed, nil
}

func (a *Auth) SignRefresh(userID, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"sid": sessionID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(a.refreshTTL).Unix(),
		"iss": "url-shortener-user",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}
	return signed, nil
}

func (a *Auth) Parse(tokenString string) (userID, sessionID string, err error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(a.secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return "", "", fmt.Errorf("parse tokenString: %w", err)
	}
	if !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid sub claim")
	}

	sessionID, ok = claims["sid"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid sid claim")
	}

	return userID, sessionID, nil
}
