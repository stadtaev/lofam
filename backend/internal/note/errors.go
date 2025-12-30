package note

import "fmt"

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func ErrValidation(msg string) ValidationError {
	return ValidationError{Message: msg}
}

type NotFoundError struct {
	ID int64
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("note with id %d not found", e.ID)
}

func ErrNotFound(id int64) NotFoundError {
	return NotFoundError{ID: id}
}
