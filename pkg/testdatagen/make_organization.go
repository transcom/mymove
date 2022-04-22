package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeOrganization creates a single Organization.
func MakeOrganization(db *pop.Connection, assertions Assertions) models.Organization {

	organizationID := assertions.Organization.ID
	if isZeroUUID(organizationID) {
		organizationID = uuid.Must(uuid.NewV4())
	}

	phone := "(510) 555-5555"
	email := "sample@organization.com"

	organization := models.Organization{
		ID:       organizationID,
		Name:     "Sample Organization",
		PocPhone: &phone,
		PocEmail: &email,
	}

	mergeModels(&organization, assertions.Organization)

	mustCreate(db, &organization, assertions.Stub)

	return organization
}

// MakeDefaultOrganization makes a default Organization
func MakeDefaultOrganization(db *pop.Connection) models.Organization {
	return MakeOrganization(db, Assertions{})
}
