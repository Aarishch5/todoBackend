package utils

import (
	"errors"
	"regexp"
)

var (
	upperCase = regexp.MustCompile(`[A-Z]`)
	lowerCase = regexp.MustCompile(`[a-z]`)
	number    = regexp.MustCompile(`[0-9]`)
	special   = regexp.MustCompile(`[^a-zA-Z0-9]`)
	emailExp  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if !upperCase.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !lowerCase.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !number.MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	if !special.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func ValidateEmail(email string) error {
	if !emailExp.MatchString(email) {
		return errors.New("invalid email address")
	}
	return nil
}
