package failure

import (
	"net/http"
)

type CustomError struct {
	Message string
	Code    int
}

func newCustomError(message string, code int) *CustomError {
	return &CustomError{Message: message, Code: code}
}

func (c CustomError) Error() string {
	return c.Message
}

var InternalServerError = newCustomError("internal server error", http.StatusInternalServerError)
var BadRequestWithString = func(message string) error {
	return newCustomError(message, http.StatusBadRequest)
}
var Error = func(err error, code int) error {
	return newCustomError(err.Error(), code)
}
var NotFound = func(message string) error {
	return newCustomError(message, http.StatusNotFound)
}

func GetCode(err error) int {
	if f, ok := err.(*CustomError); ok {
		return f.Code
	}
	return http.StatusInternalServerError
}
