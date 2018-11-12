package dpsapi

import "github.com/gofrs/uuid"

type errUserMissingData struct {
	userID     uuid.UUID
	errMessage string
}

func (e *errUserMissingData) Error() string {
	return e.errMessage
}
