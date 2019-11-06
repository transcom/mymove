package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUser creates a single User.
func MakeUser(db *pop.Connection, assertions Assertions) models.User {
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "first.last@login.gov.test",
		Active:        true,
	}

	// Overwrite values with those from assertions
	mergeModels(&user, assertions.User)

	mustCreate(db, &user)

	return user
}

// MakeDefaultUser makes a user with default values
func MakeDefaultUser(db *pop.Connection) models.User {
	return MakeUser(db, Assertions{})
}
