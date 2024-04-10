package hlsoc

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID
	Password   string
	FirstName  string
	SecondName string
	BirthDate  time.Time
	Gender     Gender
	Biography  string
	City       string
}

type Gender uint8

const (
	Male   Gender = 1
	Female Gender = 2
)

type UserFilter struct {
	ID                    uuid.UUID
	FirstName, SecondName string
	OrderBy               string
	Limit, Offset         int
}

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, filter UserFilter) (*User, error)
	SearchUsers(ctx context.Context, filter UserFilter) ([]*User, error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type UserClaims struct {
	UserID string `json:"user_id"`
}

type Tokenizer interface {
	CreateToken(user *User, duration time.Duration) (string, error)
	Verify(accessToken string) (UserClaims, error)
}
