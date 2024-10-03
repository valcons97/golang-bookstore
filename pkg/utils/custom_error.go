package utils

import "errors"

var (
	ErrBookNotFound = errors.New("book not found")

	ErrDuplicateEmail       = errors.New("duplicate email")
	ErrWrongPassword        = errors.New("wrong password")
	ErrEmptyEmailOrPassword = errors.New("empty")
	ErrEmailNotFound        = errors.New("email not found")
)
