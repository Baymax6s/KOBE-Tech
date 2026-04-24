package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultTokenExpiration = 24 * time.Hour

type Issuer struct {
	secret     []byte
	expiration time.Duration
}

func NewIssuer(secret string) *Issuer {
	return &Issuer{secret: []byte(secret), expiration: defaultTokenExpiration}
}

func (i *Issuer) Issue(userID int64) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(i.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(i.secret)
}
