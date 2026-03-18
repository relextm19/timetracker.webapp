package users

import "golang.org/x/crypto/bcrypt"

func GetHash(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), 14)
	return string(bytes), err
}
