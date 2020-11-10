package testdatagen

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUser creates a single User.
func MakeUser(db *pop.Connection, assertions Assertions) models.User {

	lgu := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  &lgu,
		LoginGovEmail: "first.last@login.gov.test",
		Active:        true,
	}

	// Overwrite values with those from assertions
	mergeModels(&user, assertions.User)

	mustCreate(db, &user, assertions.Stub)

	return user
}

// MakeDefaultUser makes a user with default values
func MakeDefaultUser(db *pop.Connection) models.User {
	return MakeUser(db, Assertions{})
}

// MakeStubbedUser returns a user without hitting the DB
func MakeStubbedUser(db *pop.Connection) models.User {
	return MakeUser(db, Assertions{Stub: true})
}
