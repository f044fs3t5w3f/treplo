package errors

import (
	"errors"
)

var ErrNoField = errors.New("required fields is missing")

// FileProcessingError is an error that contains information about file and message that we can send to the owner of it
type FileProcessingError struct {
	error
	UserMessage string
}

func NewErrorForUser(s string, userMessage string) FileProcessingError {
	return FileProcessingError{
		error:       errors.New(s),
		UserMessage: userMessage,
	}
}
