package errors

import (
	"github.com/tamboto2000/gotoko-pos/apperror"
)

var (
	ErrInternal     = apperror.New("internal server error", apperror.InternalError, "")
	ErrBadRequest   = apperror.New("bad request", apperror.BadRequest, "")
	ErrUnauthorized = apperror.New("Unauthorized", apperror.InvalidAuth, "")
)
