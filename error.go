package form

import (
	"fmt"
)

//ErrUnexpectedStatusCode is returned when any request to form API returns response status code that can't be translated into more meaningful error
type ErrUnexpectedStatusCode struct {
	StatusCode int
}

//Error as in error interface implementation
func (err ErrUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexepcterd error code %d", err.StatusCode)
}

//ErrNotFound is returned when resource doesn't exists
type ErrNotFound struct{}

//Error as in error interface implementation
func (err ErrNotFound) Error() string {
	return "not found"
}

//ErrValidationError is returned when parameters passed to the client are known to be wrong
type ErrValidationError struct {
	Reason string
}

//Error as in error interface implementation
func (err ErrValidationError) Error() string {
	return fmt.Sprintf("request body isn't valid: %s", err.Reason)
}

//ErrConflict when requested resource/action can't be completed due to logic conflict
type ErrConflict struct {
	Reason string
}

//Error as in error interface implementation
func (err ErrConflict) Error() string {
	return err.Reason
}

//ErrBadRequest is returned when form API returns BadRequest status code
type ErrBadRequest struct {
	Reason string
}

//Error as in error interface implementation
func (err ErrBadRequest) Error() string {
	return fmt.Sprintf("bad request: %s", err.Reason)
}
