package errors

import (
	"errors"
)

var (
	ErrBrokerNotFound        = errors.New("broker not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrDeviceNotFound        = errors.New("device not found")
	ErrDeviceControlNotFound = errors.New("device control not found")
	ErrInvalidPassword       = errors.New("invalid user password")
	ErrEmailExists           = errors.New("email is already taken")
	ErrBrokerServerExists    = errors.New("provided server already exists")
	ErrMissingAuthToken      = errors.New("missing authorization token")
	ErrControlStateExists    = errors.New("one state control per device is allowed")
	ErrMissingParams         = errors.New("no valid properties were provided")
	ErrInvalidRefreshToken   = errors.New("refresh token is invalid")
	ErrNoBrokerCredentials   = errors.New("broker credentials are not set")
)

type AppError struct {
	InternalError error
}

type HTTPError struct {
	Message string `json:"message"`
}

func NewHTTPError(err error) *HTTPError {
	return &HTTPError{Message: err.Error()}
}
