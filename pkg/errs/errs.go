package errs

import "errors"

var (
	ErrEmptyTitle    = errors.New("title cannot be empty")
	ErrInvalidID     = errors.New("invalid task ID")
	ErrInvalidStatus = errors.New("status must be 'new', 'done' or 'canceled'")
	ErrTaskNotFound  = errors.New("task not found")
)
