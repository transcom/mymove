package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v5"
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

// MakeActiveOfficeUser returns an active office user
func MakeActiveOfficeUser(db *pop.Connection) models.OfficeUser {
	officeUser := MakeDefaultOfficeUser(db)

	officeUser.Active = true

	err := db.Update(&officeUser)

	if err != nil {
		log.Fatal(err)
	}

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
			ID:   uuid.Must(uuid.NewV4()),
			User: tooUser,
		},
		Stub: assertions.Stub,
	})

	return officeUser
}

// MakePPMOfficeUser makes an OfficeUser with the PPM role
func MakePPMOfficeUser(db *pop.Connection, assertions Assertions) models.OfficeUser {
	ppmRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypePPMOfficeUsers,
		RoleName: "PPP Office User",
	}

	ppmUser := models.User{
		Roles: []roles.Role{ppmRole},
	}

	officeUser := MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			ID:   uuid.Must(uuid.NewV4()),
			User: ppmUser,
		},
		Stub: assertions.Stub,
	})

	return officeUser
}

// MakeOfficeUserWithUSMCGBLOC makes an OfficeUser tied to the USMC GBLOC
func MakeOfficeUserWithUSMCGBLOC(db *pop.Connection) models.OfficeUser {
	officeUUID, _ := uuid.NewV4()
	transportationOffice := MakeTransportationOffice(db, Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "USMC",
			ID:    officeUUID,
		},
	})

	tooRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeTOO,
		RoleName: "Transportation Ordering Officer",
	}

	tioRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeTIO,
		RoleName: "Transportation Invoicing Officer",
	}

	txoUser := models.User{
		Roles: []roles.Role{tooRole, tioRole},
	}

	return MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			ID:                   uuid.Must(uuid.NewV4()),
			User:                 txoUser,
			TransportationOffice: transportationOffice,
		},
	})
}
