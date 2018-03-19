package models

// These are errors that are returned by various model functions

type FetchError string

const (
	FetchErrorNotFound  FetchError = "NOT_FOUND"
	FetchErrorForbidden FetchError = "NOT_AUTHORIZED"
)
