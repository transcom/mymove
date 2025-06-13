package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

// MakeOfficeUser creates a single office user and associated TransportOffice
func MakeOfficeUser(db *pop.Connection, assertions Assertions) (models.OfficeUser, error) {
	user := assertions.OfficeUser.User
	// There's a uniqueness constraint on office user emails so add some randomness
	email := fmt.Sprintf("leo_spaceman_office_%s@example.com", MakeRandomString(5))

	if assertions.OfficeUser.UserID == nil || isZeroUUID(*assertions.OfficeUser.UserID) {

		if assertions.User.OktaEmail == "" {
			assertions.User.OktaEmail = email
		}

		user = MakeUser(db, assertions)
	}

	if assertions.User.OktaEmail != "" {
		email = assertions.User.OktaEmail
	}
	if user.Roles == nil {
		officeRole := roles.Role{
			ID:       uuid.Must(uuid.NewV4()),
			RoleType: roles.RoleTypeTOO,
			RoleName: "TOO Users",
		}

		user.Roles = []roles.Role{officeRole}
	}

	office := assertions.OfficeUser.TransportationOffice
	if isZeroUUID(office.ID) {
		var err error
		office, err = MakeTransportationOffice(db, assertions)
		if err != nil {
			return models.OfficeUser{}, err
		}
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

	mustCreate(db, &officeUser, assertions.Stub)

	return officeUser, nil
}
