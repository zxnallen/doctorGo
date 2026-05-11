package jwt

import (
	"errors"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID    uint64 `json:"user_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	golangjwt.RegisteredClaims
}

func Generate(secret string, userID uint64, expireSeconds int) (string, error) {
	return GenerateWithRole(secret, userID, "user", expireSeconds)
}

func GenerateWithRole(secret string, userID uint64, role string, expireSeconds int) (string, error) {
	return GenerateWithRoleAndType(secret, userID, role, "access", expireSeconds)
}

func GenerateWithRoleAndType(secret string, userID uint64, role string, tokenType string, expireSeconds int) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: golangjwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: golangjwt.NewNumericDate(now.Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  golangjwt.NewNumericDate(now),
		},
	}
	token := golangjwt.NewWithClaims(golangjwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func Parse(secret string, tokenString string) (*Claims, error) {
	token, err := golangjwt.ParseWithClaims(tokenString, &Claims{}, func(token *golangjwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, golangjwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
