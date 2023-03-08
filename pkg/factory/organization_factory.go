package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildOrganization creates a Organization
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildOrganization(db *pop.Connection, customs []Customization, traits []Trait) models.Organization {
	customs = setupCustomizations(customs, traits)

	// Find Organization assertion and convert to models organization
	var cOrganization models.Organization
	if result := findValidCustomization(customs, Organization); result != nil {
		cOrganization = result.Model.(models.Organization)
		if result.LinkOnly {
			return cOrganization
		}
	}

	// create organization
	phone := "(510) 555-5555"
	email := "sample@organization.com"

	organization := models.Organization{
		Name:     "Sample Organization",
		PocPhone: &phone,
		PocEmail: &email,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&organization, cOrganization)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &organization)
	}

	return organization
}
