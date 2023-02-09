package factory

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildOfficeUser creates an OfficeUser
// Also creates, if not provided
//   - User
//   - TransportationOffice
//   - calls BuildUserAndUsersRoles which creates Roles and UsersRoles
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
// Notes:
//   - There's a uniqueness constraint on office user emails so use the GetTraitOfficeUserEmail trait
//     when creating a test with multiple office users
//   - The OfficeUser returned won't have an ID if the db is nil. If an ID is needed for a stubbed user,
//     use trait GetTraitOfficeUserWithID
//   - Don't use traits to build an office user with multiple roles. Either use BuildOfficeUserWithRoles
//     or pass in a list of roles into the User customization []roles.Role{roles.Role{RoleType: roles.RoleTypeTOO}}
func BuildOfficeUser(db *pop.Connection, customs []Customization, traits []Trait) models.OfficeUser {
	customs = setupCustomizations(customs, traits)

	// Find officeuser assertion and convert to models officeuser
	var cOfficeUser models.OfficeUser
	if result := findValidCustomization(customs, OfficeUser); result != nil {
		cOfficeUser = result.Model.(models.OfficeUser)
		if result.LinkOnly {
			return cOfficeUser
		}
	}

	// Find/create the user model
	var user models.User
	result := findValidCustomization(customs, User)
	if result != nil {
		user = result.Model.(models.User)
	}
	user = BuildUserAndUsersRoles(db, customs, nil)
	// At this point, user exists. It's either the provided or created user

	// Find/create the TransportationOffice model
	var transportationOffice models.TransportationOffice
	result = findValidCustomization(customs, TransportationOffice)
	if result != nil {
		transportationOffice = result.Model.(models.TransportationOffice)
	}
	transportationOffice = BuildTransportationOffice(db, customs, nil)
	// At this point, TransportationOffice exists. It's either the provided or created TransportationOffice

	// create officeuser
	officeUser := models.OfficeUser{
		UserID:                 &user.ID,
		User:                   user,
		FirstName:              "Leo",
		LastName:               "Spaceman",
		Email:                  "leo_spaceman_office@example.com",
		Telephone:              "415-555-1212",
		TransportationOffice:   transportationOffice,
		TransportationOfficeID: transportationOffice.ID,
	}
	// Overwrite values with those from assertions
	testdatagen.MergeModels(&officeUser, cOfficeUser)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &officeUser)
	}

	return officeUser
}

// BuildOfficeUserWithRoles returns an office user with an ID, unique email, and roles
// Also creates
//   - User
//   - Role
//   - UsersRoles
func BuildOfficeUserWithRoles(db *pop.Connection, roleTypes []roles.RoleType) models.OfficeUser {
	var rolesList []roles.Role
	for _, roleType := range roleTypes {
		role := roles.Role{
			RoleType: roleType,
		}
		rolesList = append(rolesList, role)
	}

	traits := []Trait{GetTraitOfficeUserEmail}
	if db == nil {
		// UUIDs are only set when saving to a DB, but they're necessary when checking session auths
		traits = append(traits, GetTraitOfficeUserWithID)
	}
	return BuildOfficeUser(db, []Customization{
		{
			Model: models.User{
				Roles: rolesList,
			},
		},
	}, traits)
}

// ------------------------
//        TRAITS
// ------------------------

// GetTraitOfficeUserEmail helps comply with the uniqueness constraint on emails
func GetTraitOfficeUserEmail() []Customization {
	// There's a uniqueness constraint on office user emails so add some randomness
	email := strings.ToLower(fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5)))
	return []Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
			},
			Type: &User,
		},
		{
			Model: models.OfficeUser{
				Email: email,
			},
			Type: &OfficeUser,
		},
	}
}

// GetTraitOfficeUserWithID adds a UUID to the record regardless of whether it's stubbed or not
func GetTraitOfficeUserWithID() []Customization {
	return []Customization{
		{
			Model: models.OfficeUser{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
		{
			Model: models.User{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
	}
}

// GetTraitOfficeUserTIO creates an office user with a unique email and the TIO role
func GetTraitOfficeUserTIO() []Customization {
	// There's a uniqueness constraint on office user emails so add some randomness
	email := strings.ToLower(fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5)))

	tioRole := roles.Role{
		RoleType: roles.RoleTypeTIO,
	}

	return []Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Roles:         []roles.Role{tioRole},
			},
		},
		{
			Model: models.OfficeUser{
				Email: email,
			},
		},
	}
}

// GetTraitOfficeUserTOO creates an office user with a unique email and the TOO role
func GetTraitOfficeUserTOO() []Customization {
	// There's a uniqueness constraint on office user emails so add some randomness
	email := strings.ToLower(fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5)))

	tooRole := roles.Role{
		RoleType: roles.RoleTypeTOO,
	}

	return []Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Roles:         []roles.Role{tooRole},
			},
		},
		{
			Model: models.OfficeUser{
				Email: email,
			},
		},
	}
}

// GetTraitOfficeUserServicesCounselor creates an office user with a unique email and the Services counselor role
func GetTraitOfficeUserServicesCounselor() []Customization {
	// There's a uniqueness constraint on office user emails so add some randomness
	email := strings.ToLower(fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5)))

	scRole := roles.Role{
		RoleType: roles.RoleTypeServicesCounselor,
	}

	return []Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Roles:         []roles.Role{scRole},
			},
		},
		{
			Model: models.OfficeUser{
				Email: email,
			},
		},
	}
}

// GetTraitOfficeUserQAECSR creates an office user with a unique email and the Services counselor role
func GetTraitOfficeUserQAECSR() []Customization {
	// There's a uniqueness constraint on office user emails so add some randomness
	email := strings.ToLower(fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5)))

	qaeCsrRole := roles.Role{
		RoleType: roles.RoleTypeQaeCsr,
	}

	return []Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Roles:         []roles.Role{qaeCsrRole},
			},
		},
		{
			Model: models.OfficeUser{
				Email: email,
			},
		},
	}
}
