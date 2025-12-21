package domain

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
	Entity string
	ID     int64
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %d not found", e.Entity, e.ID)
}

func ErrNotFound(entity string, id int64) NotFoundError {
	return NotFoundError{Entity: entity, ID: id}
}
