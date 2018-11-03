package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// GetUser returns a single User
func GetUser(assertions Assertions) *models.User {
	user := &models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "first.last@login.gov.test",
	}

	// Overwrite values with those from assertions
	mergeModels(user, assertions.User)
	return user
}

// GetDefaultUser returns a user with default values {
func GetDefaultUser() *models.User {
	return GetUser(Assertions{})
}

// MakeUser creates a single User.
func MakeUser(db *pop.Connection, assertions Assertions) models.User {
	user := GetUser(assertions)
	mustCreate(db, user)
	return *user
}

// MakeDefaultUser makes a user with default values
func MakeDefaultUser(db *pop.Connection) models.User {
	return MakeUser(db, Assertions{})
}
