package services

import (
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"

	"github.com/transcom/mymove/pkg/models"
)

// InitializeUserResponse encapsulates a response from InitializeUser
type InitializeUserResponse struct {
	UserID          uuid.UUID
	ServiceMemberID uuid.UUID
	OfficeUserID    uuid.UUID
	TspUserID       uuid.UUID
	FirstName       string
	Middle          string
	LastName        string
}

// UserInitializer is the service object interface for CreateForm
type UserInitializer interface {
	InitializeUser(openIDUser goth.User) (response *models.UserIdentity, err error)
}
