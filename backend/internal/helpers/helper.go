package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmptyField = errors.New("empty field")

func ValidStringField(s string) bool { return strings.TrimSpace(s) != "" }

func GetHashFromString(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), 14)
	return string(bytes), err
}

func GetHash() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(b)), nil
}
