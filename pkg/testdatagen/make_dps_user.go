package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDpsUser grants DPS user permissions to the given user, or creates a new
// email associated with a DPS user. Note that this does not create the actual
// user, only gives permissions.
func MakeDpsUser(db *pop.Connection, assertions Assertions) models.DpsUser {
	email := "first.last@login.gov.test"
	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	user := models.DpsUser{
		LoginGovEmail: email,
	}

	mustCreate(db, &user)
	return user
}
