package apperror

import (
	"reflect"
	"strings"
)

type ErrorList struct {
	ty     string
	prefix string
	errs   []Error
}

func NewErrorList() ErrorList {
	return ErrorList{}
}

func ListFromError(err error) (bool, ErrorList) {
	val := reflect.ValueOf(err)
	if val.Type().String() != "apperror.ErrorList" {
		return false, ErrorList{}
	}

	return true, val.Interface().(ErrorList)
}

func (errl *ErrorList) SetType(t string) {
	errl.ty = t
}

func (errl *ErrorList) Type() string {
	return errl.ty
}

func (errl *ErrorList) Add(err Error) {
	errl.errs = append(errl.errs, err)
}

func (errl *ErrorList) SetPrefix(p string) {
	errl.prefix = p
}

// Errors will return nil if no errors in errs
func (errl *ErrorList) Errors() []Error {
	if len(errl.errs) != 0 {
		return errl.errs
	}

	return nil
}

func (errl *ErrorList) BuildError() error {
	if len(errl.errs) == 0 {
		return nil
	}

	return *errl
}

func (errl ErrorList) Error() string {
	l := make([]string, 0)
	for _, err := range errl.errs {
		l = append(l, err.Message)
	}

	str := strings.Join(l, ". ")
	if errl.prefix != "" {
		str = errl.prefix + str
	}

	return str
}
