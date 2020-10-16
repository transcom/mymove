package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

// MakeOfficeUser creates a single office user and associated TransportOffice
func MakeOfficeUser(db *pop.Connection, assertions Assertions) models.OfficeUser {
	user := assertions.OfficeUser.User
	// There's a uniqueness constraint on office user emails so add some randomness
	email := fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5))

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

	mustCreate(db, &officeUser, assertions.Stub)

	return officeUser
}

// MakeOfficeUserWithNoUser creates a single office user and associated TransportOffice
func MakeOfficeUserWithNoUser(db *pop.Connection, assertions Assertions) models.OfficeUser {
	// There's a uniqueness constraint on office user emails so add some randomness
	email := fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5))

	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	office := assertions.OfficeUser.TransportationOffice
	if isZeroUUID(office.ID) {
		office = MakeTransportationOffice(db, assertions)
	}

	officeUser := models.OfficeUser{
		TransportationOffice:   office,
		TransportationOfficeID: office.ID,
		FirstName:              "Leo",
		LastName:               "Spaceman",
		Email:                  email,
		Telephone:              "415-555-1212",
	}

	mergeModels(&officeUser, assertions.OfficeUser)

	mustCreate(db, &officeUser, assertions.Stub)

	return officeUser
}

// MakeDefaultOfficeUser makes an OfficeUser with default values
func MakeDefaultOfficeUser(db *pop.Connection) models.OfficeUser {
	return MakeOfficeUser(db, Assertions{})
}

// MakeTIOOfficeUser makes an OfficeUser with the TIO role
func MakeTIOOfficeUser(db *pop.Connection, assertions Assertions) models.OfficeUser {
	tioRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeTIO,
		RoleName: "Transportation Invoicing Officer",
	}

	tioUser := models.User{
		Roles: []roles.Role{tioRole},
	}

	officeUser := MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			User: tioUser,
		},
		Stub: assertions.Stub,
	})

	return officeUser
}

// MakeTOOOfficeUser makes an OfficeUser with the TOO role
func MakeTOOOfficeUser(db *pop.Connection, assertions Assertions) models.OfficeUser {
	tooRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeTOO,
		RoleName: "Transportation Ordering Officer",
	}

	tooUser := models.User{
		Roles: []roles.Role{tooRole},
	}

	officeUser := MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			User: tooUser,
		},
		Stub: assertions.Stub,
	})

	return officeUser
}
