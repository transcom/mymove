package testdatagen

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUser creates a single User.
func MakeUser(db *pop.Connection, assertions Assertions) models.User {

	user := models.User{
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
	lgu := uuid.Must(uuid.NewV4())
	return MakeUser(db, Assertions{
		User: models.User{
			LoginGovUUID: &lgu,
		},
	})
}

// MakeStubbedUser returns a user without hitting the DB
func MakeStubbedUser(db *pop.Connection) models.User {
	return MakeUser(db, Assertions{Stub: true})
}
