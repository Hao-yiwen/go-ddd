package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	email := strings.TrimSpace(strings.ToLower(value))
	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: email}, nil
}

func (e *Email) String() string {
	return e.value
}
