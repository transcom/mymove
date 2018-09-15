package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeOfficeUser creates a single office user and associated TransportOffice
func MakeOfficeUser(db *pop.Connection, assertions Assertions) models.OfficeUser {
	user := assertions.OfficeUser.User
	email := "leo_spaceman_office@example.com"

	if assertions.OfficeUser.UserID == nil || isZeroUUID(*assertions.OfficeUser.UserID) {
		if assertions.User.LoginGovEmail == "" {
			assertions.User.LoginGovEmail = email
		}
		user = MakeUser(db, assertions)
	}

	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	office := assertions.OfficeUser.TransportationOffice
	if isZeroUUID(office.ID) {
		office = MakeTransportationOffice(db, assertions)
	}

	officeUser := models.OfficeUser{
		UserID:                 &user.ID,
		User:                   user,
		TransportationOffice:   office,
		TransportationOfficeID: office.ID,
		FirstName:              "Leo",
		LastName:               "Spaceman",
		Email:                  email,
		Telephone:              "415-555-1212",
	}

	mergeModels(&officeUser, assertions.OfficeUser)

	mustCreate(db, &officeUser)

	return officeUser
}

// MakeDefaultOfficeUser makes an OfficeUser with default values
func MakeDefaultOfficeUser(db *pop.Connection) models.OfficeUser {
	return MakeOfficeUser(db, Assertions{})
}
