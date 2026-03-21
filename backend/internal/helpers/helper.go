package helpers

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmptyField = errors.New("empty field")

func ValidStringField(s string) bool { return strings.TrimSpace(s) != "" }

func GetHash(b []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(b, 14)
	return string(bytes), err
}
