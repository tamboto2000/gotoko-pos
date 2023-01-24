package httpresponse

import (
	"net/http"

	"github.com/tamboto2000/gotoko-pos/apperror"
)

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Error      interface{} `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"-"`
}

type defaultErrObj struct{}

// errStatusMap map error types with their appropriate
// status code
var errStatusMap map[string]int = map[string]int{
	apperror.AnyRequired:   http.StatusBadRequest,
	apperror.AnyInvalid:    http.StatusBadRequest,
	apperror.NotFound:      http.StatusNotFound,
	apperror.InternalError: http.StatusInternalServerError,
	apperror.ObjecMissing:  http.StatusBadRequest,
	apperror.BadRequest:    http.StatusBadRequest,
	apperror.InvalidAuth:   http.StatusUnauthorized,
	apperror.ArrayBase:     http.StatusBadRequest,
}

func Success(data interface{}) Response {
	return Response{
		Success:    true,
		Message:    "Success",
		Data:       data,
		StatusCode: http.StatusOK,
	}
}

func FromError(err error) Response {
	if ok, aerr := apperror.FromError(err); ok {
		return Response{
			Success:    false,
			Message:    err.Error(),
			Error:      defaultErrObj{},
			StatusCode: errStatusMap[aerr.Type],
		}
	} else if ok, aerr := apperror.ListFromError(err); ok {
		return Response{
			Success:    false,
			Message:    err.Error(),
			Error:      aerr.Errors(),
			StatusCode: errStatusMap[aerr.Type()],
		}
	}

	return Response{
		Success:    false,
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
		Error:      defaultErrObj{},
	}
}
