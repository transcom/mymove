package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// UserMaker is the base maker function to create a user
func BuildUser(db *pop.Connection, customs []Customization, traits []Trait) models.User {
	customs = setupCustomizations(customs, traits)

	// Find user assertion and convert to models user
	var cUser models.User
	if result := findValidCustomization(customs, User); result != nil {
		cUser = result.Model.(models.User)
	}

	// create user
	// MYTODO: Add forceUUID functionality
	loginGovUUID := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "first.last@login.gov.test",
		Active:        false,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&user, cUser)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &user)
	}

	return user
}

// BuildDefaultUser creates an active user
func BuildDefaultUser(db *pop.Connection) models.User {
	return BuildUser(db, nil, []Trait{GetTraitActiveUser})
}

// ------------------------
//      TRAITS
// ------------------------

// GetTraitActiveUser returns a customization to enable active on a user
func GetTraitActiveUser() []Customization {
	return []Customization{
		{
			Model: models.User{
				Active: true,
			},
		},
	}
}
