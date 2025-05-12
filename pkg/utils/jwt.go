package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var secret = []byte("your-secret-key")
var jwtKey = []byte("your-secret")

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}
	return claims, nil
}
