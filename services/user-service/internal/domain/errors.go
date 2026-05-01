package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailExists       = errors.New("email already exists")
	ErrUsernameExists    = errors.New("username already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidInput      = errors.New("invalid input")
)
