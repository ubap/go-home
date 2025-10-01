package reku

import "errors"

// ErrUnauthorized is a sentinel error returned when a user's
// authentication fails or is missing.
var ErrUnauthorized = errors.New("user is unauthorized")
