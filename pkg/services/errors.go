package services

import "errors"

// ErrFetchForbidden means that the record exists but that the user does not have access to it
var ErrFetchForbidden = errors.New("FETCH_FORBIDDEN")
