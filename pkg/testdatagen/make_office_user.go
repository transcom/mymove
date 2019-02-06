package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeOfficeUser creates a single office user and associated TransportOffice
func MakeOfficeUser(db *pop.Connection, assertions Assertions) models.OfficeUser {

	user := assertions.OfficeUser.User
	if assertions.OfficeUser.UserID == nil || isZeroUUID(*assertions.OfficeUser.UserID) {
		if assertions.User.LoginGovEmail == "" {
			// There's a uniqueness constraint on office user emails so add some randomness
			email := fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5))
			assertions.User.LoginGovEmail = email
		}
		user = MakeUser(db, assertions)
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
		Email:                  user.LoginGovEmail,
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
