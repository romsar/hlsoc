package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/romsar/hlsoc"
	"time"
)

type Tokenizer struct {
	secret string
}

func New(secret string) *Tokenizer {
	return &Tokenizer{secret: secret}
}

type UserClaims struct {
	jwt.RegisteredClaims
	hlsoc.UserClaims
}

func (t *Tokenizer) CreateToken(user *hlsoc.User, duration time.Duration) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserClaims: hlsoc.UserClaims{
			UserID: user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (t *Tokenizer) Verify(accessToken string) (hlsoc.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(t.secret), nil
		},
	)

	if err != nil {
		return hlsoc.UserClaims{}, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return hlsoc.UserClaims{}, fmt.Errorf("invalid token claims")
	}

	return claims.UserClaims, nil
}
