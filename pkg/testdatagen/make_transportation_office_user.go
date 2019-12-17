package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTransportationOrderingOfficer creates a single TransportationOrderingOfficer user and associated TransportOffice
func MakeTransportationOrderingOfficer(db *pop.Connection, assertions Assertions) models.TransportationOrderingOfficer {
	user := assertions.TransportationOrderingOfficer.User
	// There's a uniqueness constraint on office user emails so add some randomness
	email := fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5))

	if assertions.TransportationOrderingOfficer.UserID == nil || isZeroUUID(*assertions.TransportationOrderingOfficer.UserID) {
		if assertions.User.LoginGovEmail == "" {
			assertions.User.LoginGovEmail = email
		}
		//TODO check if User really needs to be a pointer??
		u := MakeUser(db, assertions)
		user = &u
	}

	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	transportationOrderingOfficer := models.TransportationOrderingOfficer{
		UserID: &user.ID,
		User:   user,
	}

	mergeModels(&transportationOrderingOfficer, assertions.TransportationOrderingOfficer)

	mustCreate(db, &transportationOrderingOfficer)

	return transportationOrderingOfficer
}

// MakeDefaultTransportationOrderingOfficer makes an TransportationOrderingOfficer with default values
func MakeDefaultTransportationOrderingOfficer(db *pop.Connection) models.TransportationOrderingOfficer {
	return MakeTransportationOrderingOfficer(db, Assertions{})
}
