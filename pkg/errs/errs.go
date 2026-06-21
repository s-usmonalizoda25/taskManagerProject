package errs

import "errors"

var (
	ErrInvalidID     = errors.New("invalid ID format")
	ErrInvalidStatus = errors.New("status must be 'new', 'done' or 'canceled'")

	ErrEmptyTitle   = errors.New("title cannot be empty")
	ErrTaskNotFound = errors.New("task not found")

	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrEmptyEmail    = errors.New("email cannot be empty")
	ErrUserNotFound  = errors.New("user not found")
	ErrUsernameTaken = errors.New("username is already taken")
	ErrEmailTaken    = errors.New("email is already registered")
)
