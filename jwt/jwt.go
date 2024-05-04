package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (t *Tokenizer) CreateToken(user *hlsoc.User, duration time.Duration) (string, error) {
	claims := UserClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (t *Tokenizer) Verify(accessToken string) (*hlsoc.User, error) {
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
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse uuid %s err: %w", claims.UserID, err)
	}

	return &hlsoc.User{ID: userID}, nil
}
