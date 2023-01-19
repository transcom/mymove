package factory

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildOfficeUser creates an OfficeUser
// Also creates, if not provided
//   - User
//   - TransportationOffice
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
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
	user = BuildUser(db, customs, nil)
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

// BuildDefaultOfficeUser returns an office user with appropriate email
// Also creates
//   - User
func BuildDefaultOfficeUser(db *pop.Connection) models.OfficeUser {
	return BuildOfficeUser(db, nil, []Trait{GetTraitOfficeUserEmail})
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

// GetTraitOfficeUserTIO helps comply with the uniqueness constraint on emails
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

// GetTraitOfficeUserTIO helps comply with the uniqueness constraint on emails
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
