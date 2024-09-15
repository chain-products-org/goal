package errorx

import (
	"fmt"
	"net/http"
	"reflect"
)

var (
	ServerBusy       = NewPreferredCodeErrf(http.StatusInternalServerError, "server is busy")
	BadReq           = NewPreferredCodeErrf(http.StatusBadRequest, "bad request")
	NoAuth           = NewPreferredCodeErrf(http.StatusUnauthorized, "unauthorized")
	FrequentReq      = NewPreferredCodeErrf(http.StatusTooManyRequests, "frequent requests")
	InvalidOperation = NewPreferredCodeErrf(http.StatusForbidden, "invalid operation")
)

// PreferredError is a struct to descibe a preferred error, which will be processed specially, such as
// response it to http
type PreferredError struct {
	error
	code int
}

// SetCode set a http status code to PreferredError. To maintain encapsulation, prefer to use
// NewPreferredErrCode/NewPreferredCodeErrf to create a PreferredError
func (e *PreferredError) SetCode(code int) *PreferredError {
	e.code = code
	return e
}

func (e *PreferredError) Error() string {
	return e.error.Error()
}

func (e *PreferredError) Code() int {
	return e.code
}

// NewPreferredErr create a PreferredError with given error, code is 0, call SetCode to set a http status code
// Deprecated: prefer to NewPreferredErrCode
func NewPreferredErr(err error) error {
	return &PreferredError{error: err}
}

// NewPreferredErrCode create a PreferredError with given error and code
func NewPreferredErrCode(err error, code int) error {
	return &PreferredError{error: err, code: code}
}

// NewPreferredErrf create a PreferredError with error message, code is 0, call SetCode to set a http status code
// Deprecated: prefer to NewPreferredErrCodeErrf
func NewPreferredErrf(format string, args ...any) error {
	return &PreferredError{error: fmt.Errorf(format, args...)}
}

// NewPreferredCodeErrf create a PreferredError with given code and error message
func NewPreferredCodeErrf(code int, format string, args ...any) error {
	return &PreferredError{code: code, error: fmt.Errorf(format, args...)}
}

// Prefer400 create a PreferredError with http.StatusBadRequest
func Prefer400(format string, args ...any) error {
	return &PreferredError{code: http.StatusBadRequest, error: fmt.Errorf(format, args...)}
}

// Prefer403 create a PreferredError with http.StatusForbidden
func Prefer403(format string, args ...any) error {
	return &PreferredError{code: http.StatusForbidden, error: fmt.Errorf(format, args...)}
}

// Prefer415 create a PreferredError with http.StatusUnsupportedMediaType
func Prefer415(format string, args ...any) error {
	return &PreferredError{code: http.StatusUnsupportedMediaType, error: fmt.Errorf(format, args...)}
}

// Prefer423 create a PreferredError with http.StatusLocked
func Prefer423(format string, args ...any) error {
	return &PreferredError{code: http.StatusLocked, error: fmt.Errorf(format, args...)}
}

// Prefer426 create a PreferredError with http.StatusUpgradeRequired
func Prefer426(format string, args ...any) error {
	return &PreferredError{code: http.StatusUpgradeRequired, error: fmt.Errorf(format, args...)}
}

// Prefer429 create a PreferredError with http.StatusTooManyRequests
func Prefer429(format string, args ...any) error {
	return &PreferredError{code: http.StatusForbidden, error: fmt.Errorf(format, args...)}
}

// Prefer500 create a PreferredError with http.StatusInternalServerError
func Prefer500(format string, args ...any) error {
	return &PreferredError{code: http.StatusInternalServerError, error: fmt.Errorf(format, args...)}
}

// Prefer501 create a PreferredError with http.StatusNotImplemented
func Prefer501(format string, args ...any) error {
	return &PreferredError{code: http.StatusNotImplemented, error: fmt.Errorf(format, args...)}
}

// Prefer503 create a PreferredError with http.StatusServiceUnavailable
func Prefer503(format string, args ...any) error {
	return &PreferredError{code: http.StatusServiceUnavailable, error: fmt.Errorf(format, args...)}
}

// IsPreferred checks if the given err is a PreferredError
func IsPreferred(err error) bool {
	if err == nil {
		return false
	}
	return reflect.TypeOf(err).AssignableTo(reflect.TypeOf(&PreferredError{}))
}
