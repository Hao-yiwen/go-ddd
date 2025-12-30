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

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// yiwen@gmail.com -> gmail.com
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// yiwen@gmail.com -> yiwen
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}
