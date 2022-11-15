package factory

import (
	"fmt"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildAdminUser is the base maker function to create an adminuser
func BuildAdminUser(db *pop.Connection, customs []Customization, traits []Trait) models.AdminUser {
	customs = setupCustomizations(customs, traits)

	// Find adminuser assertion and convert to models adminuser
	var cAdminUser models.AdminUser
	if result := findValidCustomization(customs, AdminUser); result != nil {
		cAdminUser = result.Model.(models.AdminUser)
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

	// create adminuser
	adminUser := models.AdminUser{
		UserID:    &user.ID,
		User:      user,
		FirstName: "Leo",
		LastName:  "Spaceman",
		Email:     "leo_spaceman_admin@example.com",
		Role:      "SYSTEM_ADMIN",
	}
	// Overwrite values with those from assertions
	testdatagen.MergeModels(&adminUser, cAdminUser)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &adminUser)
	}

	return adminUser
}

// ------------------------
//        TRAITS
// ------------------------

// GetTraitAdminUserEmail helps comply with the uniqueness constraint on emails
func GetTraitAdminUserEmail() []Customization {
	// There's a uniqueness constraint on admin user emails so add some randomness
	email := fmt.Sprintf("leo_spaceman_admin_%s@example.com", makeRandomString(5))
	return []Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
			},
			Type: &User,
		},
		{
			Model: models.AdminUser{
				Email: email,
			},
			Type: &AdminUser,
		},
	}
}
