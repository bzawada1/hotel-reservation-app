package api

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code int    `json: "code"`
	Err  string `json: "error"`
}

func NewError(code int, message string) Error {
	return Error{
		Code: 520,
		Err:  message,
	}
}

func ErrorInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id given",
	}
}

func ErrorUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized",
	}
}

func ErrorNotFound(entity string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  fmt.Sprintf("%s not found", entity),
	}
}

func ErrorBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid request payload",
	}
}

func (e Error) Error() string {
	return e.Err
}
