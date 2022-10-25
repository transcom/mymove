package factory

import (
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// UserMaker is the base maker function to create a user
// MYTODO Instead of error (not useful) can we return a list of the created objects?
func BuildUser(db *pop.Connection, customs []Customization, traits []Trait) (models.User, error) {
	customs = mergeCustomization(customs, traits)
	customs, err := validateCustomizations(customs)
	if err != nil {
		log.Panic(err)
	}

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

	// MYTODO: Add back stub functionality
	mustCreate(db, &user, false)

	return user, nil
}
