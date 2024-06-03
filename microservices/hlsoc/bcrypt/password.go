package bcrypt

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct {
	cost int
}

func New(cost int) *PasswordHasher {
	return &PasswordHasher{cost: cost}
}

func (ph *PasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	return string(bytes), err
}

func (ph *PasswordHasher) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
