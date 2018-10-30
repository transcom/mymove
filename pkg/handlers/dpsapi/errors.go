package dpsapi

import "github.com/gobuffalo/uuid"

type errUserMissingData struct {
	userID     uuid.UUID
	errMessage string
}

func (e *errUserMissingData) Error() string {
	return e.errMessage
}
