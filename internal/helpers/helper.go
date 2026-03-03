package helpers

import (
	"errors"
	"strings"
)

var ErrEmptyField = errors.New("empty field")

func ValidStringField(s string) bool { return strings.TrimSpace(s) == "" }
