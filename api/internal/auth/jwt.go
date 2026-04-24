package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	secret []byte
}

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewValidator(secret string) *Validator {
	return &Validator{secret: []byte(secret)}
}

func (v *Validator) ValidateToken(tokenString string) (int64, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return v.secret, nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	if claims.UserID == 0 {
		return 0, errors.New("user_id is required in token")
	}

	return claims.UserID, nil
}
