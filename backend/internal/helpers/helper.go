package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmptyField = errors.New("empty field")

func ValidStringField(s string) bool { return strings.TrimSpace(s) != "" }

func GetHashFromPassword(b []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(b, 10)
	return string(bytes), err
}

func GetHashFromUUID(b []byte) (string, error) {
	hash := sha256.Sum256(b)

	return hex.EncodeToString(hash[:]), nil
}
