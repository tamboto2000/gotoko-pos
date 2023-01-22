package jwtparse

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// BuildToken create a JWT using EdDSA signing method.
// secret is obtained from config
func BuildJWT(claims jwt.RegisteredClaims) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseJWT parse token string usinng EdDSA signing method
func ParseJWT(tokenStr string) (*jwt.RegisteredClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := new(jwt.RegisteredClaims)
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	})

	if token == nil {
		return claims, err
	}

	claims = token.Claims.(*jwt.RegisteredClaims)

	return claims, nil
}
