package errors

import (
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Reason     string
}

func NewHTTPError(statusCode int, message string) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
		Reason:     message,
	}
}

func InternalError(err string) HTTPError {
	return NewHTTPError(http.StatusInternalServerError, err)
}

func BadRequest(err string) HTTPError {
	return NewHTTPError(http.StatusBadRequest, err)
}

func NotFound(err string) HTTPError {
	return NewHTTPError(http.StatusNotFound, err)
}

func ToHTTPError(err error) HTTPError {
	contextualError, ok := err.(DetailedError)

	if !ok {
		return InternalError("something unexpected happened")
	}

	switch contextualError.ErrorType() {
	case ErrorInvalidInput:
		return BadRequest(contextualError.Error())
	case ErrorNotFound:
		return NotFound(contextualError.Error())
	default:
		return InternalError(contextualError.Error())
	}
}
