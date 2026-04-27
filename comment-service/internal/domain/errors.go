package domain

import "errors"

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrAlreadyLiked    = errors.New("already liked")
	ErrUnauthorized    = errors.New("unauthorized")
)
