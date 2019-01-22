package dpsauth

// ErrInvalidCookie is an error for invalid DPS authentication cookies
type ErrInvalidCookie struct {
	errMessage string
}

func (e *ErrInvalidCookie) Error() string {
	return e.errMessage
}
