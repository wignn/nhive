package domain

import "errors"

var (
	ErrEntryNotFound    = errors.New("library entry not found")
	ErrBookmarkNotFound = errors.New("bookmark not found")
	ErrAlreadyInLibrary = errors.New("novel already in library")
	ErrInvalidInput     = errors.New("invalid input")
)
