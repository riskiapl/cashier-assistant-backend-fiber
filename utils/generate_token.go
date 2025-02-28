package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riskiapl/fiber-app/config"
)

var secretKey = func() []byte {
    secret := config.GetEnv("JWT_SECRET", "fallback-secret-key")
    return []byte(secret)
}()

func GenerateToken(data map[string]interface{}) (string, error) {
    claims := jwt.MapClaims{
        "exp": time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
    }
    
    // Add all data fields to the claims
    for key, value := range data {
        claims[key] = value
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}