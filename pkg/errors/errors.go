package errors

import (
	"errors"
	"log"
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
	ErrUserCreation          = errors.New("could not create a user")
	ErrNoAwaitingConfirm     = errors.New("account does not await to be confirmed")
	ErrConfirmationRequired  = errors.New("email has to be verified")
	ErrUnexpected            = errors.New("something went wrong")
	ErrUnauthorized          = errors.New("could not authorize a user")
	ErrTokenNotFound         = errors.New("token no longer applies")
	ErrEndpointDisabled      = errors.New("the endpoint has been temporarily disabled")
)

type HTTPError struct {
	Message string `json:"message"`
}

func NewHTTPError(err error) *HTTPError {
	log.Printf("Error: %s", err.Error())
	return &HTTPError{Message: err.Error()}
}
