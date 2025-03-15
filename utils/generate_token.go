package utils

import (
	"errors"
	"time"

	"maps"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riskiapl/fiber-app/config"
)

var secretKey = func() []byte {
	secret := config.GetEnv("JWT_SECRET", "fallback-secret-key")
	return []byte(secret)
}()

func GenerateToken(data map[string]any) (string, error) {
	claims := jwt.MapClaims{}

	// Add all data fields to the claims
	maps.Copy(claims, data)

	// Check if expiration time is provided in data
	if exp, exists := data["expired"]; exists {
		claims["exp"] = exp
	} else {
		// Default: Token valid for 24 hours
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// Parse token with claims and verify using the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and verify claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
