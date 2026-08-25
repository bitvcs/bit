package domain

import (
	"fmt"
)

var ()

type Error struct {
	Code            int    // the code uses the HTTP just to make it simpler
	Message         string // message show to the user
	InternalMessage string // uses for internal logging
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) String() string {
	return fmt.Sprintf("Error{Code: %d, Message: %s, InternalMessage: %s}", e.Code, e.Message, e.InternalMessage)
}

func IsErrorNotFound(err error) bool {
	return isDomainError(err, 404)
}

func IsErrorNoPermission(err error) bool {
	return isDomainError(err, 403)
}

func isDomainError(err error, code int) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == code
}

func NewErrorRecordNotFound() *Error {
	return &Error{
		Code:    404,
		Message: "record not found",
	}
}

func NewErrorDatabase(internalMessage string) *Error {
	return &Error{
		Code:            500,
		Message:         "database error",
		InternalMessage: internalMessage,
	}
}

func NewErrorUser(message string) *Error {
	return &Error{
		Code:    400,
		Message: message,
	}
}

func NewErrorNoPermission() *Error {
	return &Error{
		Code:    403,
		Message: "no permission",
	}
}

func NewErrorInternalServer(internalMessage string) *Error {
	return &Error{
		Code:            500,
		Message:         "internal server error",
		InternalMessage: internalMessage,
	}
}
