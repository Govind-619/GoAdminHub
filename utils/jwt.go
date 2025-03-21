package utils

import (
	"time"

	"github.com/Govind-619/GoAdminHub/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key")

// GenerateToken creates a JWT token for a given user email and ID.
func GenerateToken(email string, id uint) (string, error) {
	claims := &models.Claims{
		Id:        id,
		UserEmail: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseToken validates the token string and returns the user email if valid.
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims.UserEmail, nil
	}
	return "", err
}
