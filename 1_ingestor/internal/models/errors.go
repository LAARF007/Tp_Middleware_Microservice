package models

import "errors"

// Errors a service can return to explain the problem.
// Any other error means it's an internal error and becomes a 500.
//
// You can use wrappers to explain the problem:
//
//	return nil, fmt.Errorf("source %s: %w", id, models.ErrNotFound)
//
// The %w keeps ErrNotFound inside the new error, so errors.Is can still find it.
var (
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
)
