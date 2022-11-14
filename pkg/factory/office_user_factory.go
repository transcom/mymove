package factory

import (
	"fmt"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildOfficeUser is the base maker function to create an officeuser
func BuildOfficeUser(db *pop.Connection, customs []Customization, traits []Trait) models.OfficeUser {
	customs = setupCustomizations(customs, traits)

	// Find officeuser assertion and convert to models officeuser
	var cOfficeUser models.OfficeUser
	if result := findValidCustomization(customs, OfficeUser); result != nil {
		cOfficeUser = result.Model.(models.OfficeUser)
	}

	// Find/create the required user model
	var user models.User
	linkOnly := false
	result := findValidCustomization(customs, User)
	if result != nil {
		user = result.Model.(models.User)
		linkOnly = result.LinkOnly
	}
	if !linkOnly {
		user = BuildUser(db, customs, nil)
	}
	// At this point, user exists. It's either the provided or created user

	office := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{})

	// create officeuser
	officeUser := models.OfficeUser{
		UserID:                 &user.ID,
		User:                   user,
		TransportationOffice:   office,
		TransportationOfficeID: office.ID,
		FirstName:              "Leo",
		LastName:               "Spaceman",
		Email:                  "leo_spaceman_office@example.com",
		Telephone:              "415-555-1212",
	}
	// Overwrite values with those from assertions
	testdatagen.MergeModels(&officeUser, cOfficeUser)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &officeUser)
	}

	return officeUser
}

// ------------------------
//        TRAITS
// ------------------------

func GetTraitOfficeUserEmail() []Customization {
	email := fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5))
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
