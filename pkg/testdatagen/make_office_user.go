package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v6"
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
	if user.Roles == nil {
		officeRole := roles.Role{
			ID:       uuid.Must(uuid.NewV4()),
			RoleType: roles.RoleTypePPMOfficeUsers,
			RoleName: "PPM Office Users",
		}

		user.Roles = []roles.Role{officeRole}
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
	tioRole, _ := LookupRole(db, roles.RoleTypeTIO)

	tioUser := models.User{
		Roles: []roles.Role{tioRole},
	}

	officeUser := MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			ID:   uuid.Must(uuid.NewV4()),
			User: tioUser,
		},
		Stub: assertions.Stub,
	})

	return officeUser
}

// MakeActiveOfficeUser returns an active office user
func MakeActiveOfficeUser(db *pop.Connection) models.OfficeUser {
	officeUser := MakeOfficeUser(db, Assertions{
		User: models.User{
			Active: true, // an active office user should also have an active user
		},
	})

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

// MakeServicesCounselorOfficeUser makes an OfficeUser with the ServicesCounselor role
func MakeServicesCounselorOfficeUser(db *pop.Connection, assertions Assertions) models.OfficeUser {
	servicesRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeServicesCounselor,
		RoleName: "Services Counselor",
	}

	servicesUser := models.User{
		Roles: []roles.Role{servicesRole},
	}

	officeUser := MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			ID:   uuid.Must(uuid.NewV4()),
			User: servicesUser,
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

// MakeQAECSROfficeUser makes an OfficeUser with the QAECSR role
func MakeQAECSROfficeUser(db *pop.Connection, assertions Assertions) models.OfficeUser {
	qaeCsrRole, _ := LookupRole(db, roles.RoleTypeQaeCsr)

	qaeCsrUser := models.User{
		Roles: []roles.Role{qaeCsrRole},
	}

	officeUser := MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			ID:   uuid.Must(uuid.NewV4()),
			User: qaeCsrUser,
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

// MakeServicesCounselorOfficeUserWithUSMCGBLOC makes a Services Counselor tied to the USMC GBLOC
func MakeServicesCounselorOfficeUserWithUSMCGBLOC(db *pop.Connection) models.OfficeUser {
	officeUUID, _ := uuid.NewV4()
	transportationOffice := MakeTransportationOffice(db, Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "USMC",
			ID:    officeUUID,
		},
	})

	servicesRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeServicesCounselor,
		RoleName: "Services Counselor",
	}

	servicesUser := models.User{
		Roles: []roles.Role{servicesRole},
	}

	return MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			ID:                   uuid.Must(uuid.NewV4()),
			User:                 servicesUser,
			TransportationOffice: transportationOffice,
		},
	})
}

// MakeOfficeUserWithMultipleRoles makes an OfficeUser with Counselor and TXO roles
func MakeOfficeUserWithMultipleRoles(db *pop.Connection, assertions Assertions) models.OfficeUser {
	tooRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeTOO,
		RoleName: "Transportation Ordering Officer",
	}

	servicesRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeServicesCounselor,
		RoleName: "Services Counselor",
	}

	tioRole := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeTIO,
		RoleName: "Transportation Invoicing Officer",
	}

	multipleRoleUser := models.User{
		Roles: []roles.Role{tooRole, tioRole, servicesRole},
	}

	officeUser := MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			ID:   uuid.Must(uuid.NewV4()),
			User: multipleRoleUser,
		},
		Stub: assertions.Stub,
	})

	// save roles to db
	rolesList := officeUser.User.Roles
	for _, role := range rolesList {
		newRole := MakeRole(db, Assertions{
			Role: role,
			Stub: assertions.Stub,
		})
		MakeUsersRoles(db, Assertions{
			UsersRoles: models.UsersRoles{
				UserID: officeUser.User.ID,
				RoleID: newRole.ID,
			},
			Stub: assertions.Stub,
		})
	}

	return officeUser
}

// MakeStubbedOfficeUser returns a user without hitting the DB
func MakeStubbedOfficeUser(db *pop.Connection) models.OfficeUser {
	return MakeOfficeUser(db, Assertions{
		OfficeUser: models.OfficeUser{
			ID: uuid.Must(uuid.NewV4()),
		},
		User: models.User{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})
}
