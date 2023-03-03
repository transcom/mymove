package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildServiceMember creates a single ServiceMember
// Also creates, if not provided:
// - Residential Address of the ServiceMember
// - User
func BuildServiceMember(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceMember {
	customs = setupCustomizations(customs, traits)

	// Find ServiceMember customization and extract the custom ServiceMember
	var cServiceMember models.ServiceMember
	if result := findValidCustomization(customs, ServiceMember); result != nil {
		cServiceMember = result.Model.(models.ServiceMember)
		if result.LinkOnly {
			return cServiceMember
		}
	}

	// Find/create the user model
	user := BuildUser(db, customs, traits)

	// Find/create the address model
	currentAddress := BuildAddress(db, customs, traits)

	email := "leo_spaceman_sm@example.com"
	agency := models.AffiliationARMY
	rank := models.ServiceMemberRankE1

	randomEdipi := RandomEdipi()

	serviceMember := models.ServiceMember{
		UserID:               user.ID,
		User:                 user,
		Edipi:                models.StringPointer(randomEdipi),
		Affiliation:          &agency,
		FirstName:            models.StringPointer("Leo"),
		LastName:             models.StringPointer("Spacemen"),
		Telephone:            models.StringPointer("212-123-4567"),
		PersonalEmail:        &email,
		ResidentialAddressID: &currentAddress.ID,
		ResidentialAddress:   &currentAddress,
		Rank:                 &rank,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&serviceMember, cServiceMember)

	if db != nil {
		mustCreate(db, &serviceMember)
	}

	return serviceMember
}

// ServiceMemberSetIDs is a sample GetTraitFunc
// that sets ids for both ServiceMember and User models
func ServiceMemberSetIDs() []Customization {
	return []Customization{
		{
			Model: models.ServiceMember{
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
