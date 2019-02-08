package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
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

// AppDetector detects which type of app is handling the request
type AppDetector interface {
	IsOfficeApp() bool
	IsTspApp() bool
}

// UserInitializer is the service object interface for CreateForm
type UserInitializer interface {
	InitializeUser(session AppDetector, openIDUser goth.User) (response InitializeUserResponse, verrs *validate.Errors, err error)
}
