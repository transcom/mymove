package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUser creates a single User
// It will not replace a true assertion with false.
func MakeUser(db *pop.Connection, assertions Assertions) models.User {

	loginGovUUID := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "first.last@login.gov.test",
		Active:        false,
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
			Active:       true,
		},
	})
}

// MakeStubbedUser returns a user without hitting the DB
func MakeStubbedUser(db *pop.Connection) models.User {
	return MakeUser(db, Assertions{
		User: models.User{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})
}
