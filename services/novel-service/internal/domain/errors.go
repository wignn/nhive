package domain

import "errors"

var (
	ErrNovelNotFound   = errors.New("novel not found")
	ErrChapterNotFound = errors.New("chapter not found")
	ErrSlugExists      = errors.New("novel slug already exists")
	ErrChapterExists   = errors.New("chapter number already exists for this novel")
	ErrInvalidInput    = errors.New("invalid input")
)
