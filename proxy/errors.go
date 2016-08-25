package proxy

import (
	"errors"
)

var (
	ErrInvalidAuthentication  = errors.New("Invalid authentication.")
	ErrInvalidProtocolVersion = errors.New("Invalid protocol version.")
	ErrAlreadyListening       = errors.New("Already listening on another port.")
)
