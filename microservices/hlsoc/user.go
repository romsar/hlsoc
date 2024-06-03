package hlsoc

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Password   string    `json:"password,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	SecondName string    `json:"second_name,omitempty"`
	BirthDate  time.Time `json:"birth_date"`
	Gender     Gender    `json:"gender,omitempty"`
	Biography  string    `json:"biography,omitempty"`
	City       string    `json:"city,omitempty"`
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
	GetFriends(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type Tokenizer interface {
	CreateToken(user *User, duration time.Duration) (string, error)
	Verify(accessToken string) (*User, error)
}
