package users

import (
	"errors"
	"net/mail"

	"github.com/google/uuid"
	"github.com/relextm19/tracker.nvim/internal/helpers"
)

type ClientUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewClientUserBody() *ClientUserBody {
	return &ClientUserBody{}
}

type User struct {
	Email        string
	PasswordHash string
	Token        uuid.UUID
}

func NewUser(cub *ClientUserBody) (*User, error) {
	ph, err := HashPassword(cub.Password)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:        cub.Email,
		PasswordHash: ph,
		Token:        uuid.New(),
	}, nil
}

func (u *User) Valid() error {
	if ok := helpers.ValidStringField(u.Email); !ok {
		return helpers.ErrEmptyField
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	if ok := helpers.ValidStringField(u.PasswordHash); !ok {
		return helpers.ErrEmptyField
	}

	if u.Token == uuid.Nil {
		return errors.New("invalid user token")
	}

	return nil
}
