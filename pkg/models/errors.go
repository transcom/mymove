package models

// These are errors that are returned by various model functions

// FetchError is a base type that can typecast for specific APIs,
// It indicicates why an attempted db fetch failed.
type FetchError string

const (
	// FetchErrorNotFound means that the requested record does not exist
	FetchErrorNotFound FetchError = "NOT_FOUND"
	// FetchErrorForbidden means that the record exists but that the user does not have access to it
	FetchErrorForbidden FetchError = "NOT_AUTHORIZED"
)
