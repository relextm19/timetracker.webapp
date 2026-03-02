package users

import "github.com/google/uuid"

type User struct {
	Email        string
	PasswordHash string
	Token        string
}

func NewUser(email, password string) (*User, error) {
	hp, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:        email,
		PasswordHash: hp,
		Token:        uuid.New(),
	}, nil
}
