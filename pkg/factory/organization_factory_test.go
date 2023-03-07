package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildOrganization() {
	suite.Run("Successful creation of default Organization", func() {
		// Under test:      BuildOrganization
		// Mocked:          None
		// Set up:          Create an organization no customizations or traits
		// Expected outcome:organization should be created with default values

		// SETUP
		phone := "(510) 555-5555"
		email := "sample@organization.com"

		defaultOrganization := models.Organization{
			Name:     "Sample Organization",
			PocPhone: &phone,
			PocEmail: &email,
		}

		// CALL FUNCTION UNDER TEST
		organization := BuildOrganization(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultOrganization.Name, organization.Name)
		suite.Equal(phone, *organization.PocPhone)
		suite.Equal(email, *organization.PocEmail)
	})

	suite.Run("Successful creation of customized Organization", func() {
		// Under test:       BuildOrganization
		// Set up:           Create an organization and pass custom fields,
		// Expected outcome: organization should be created with custom fields

		// SETUP
		phone := "(555) 444-4444"
		email := "custom@organization.com"
		organizationID := uuid.Must(uuid.NewV4())

		customOrganization := models.Organization{
			ID:       organizationID,
			Name:     "Custom Organization",
			PocPhone: &phone,
			PocEmail: &email,
		}

		// CALL FUNCTION UNDER TEST
		organization := BuildOrganization(suite.DB(), []Customization{
			{Model: customOrganization},
		}, nil)

		suite.Equal(customOrganization.ID, organization.ID)
		suite.Equal(customOrganization.Name, organization.Name)
		suite.Equal(phone, *organization.PocPhone)
		suite.Equal(email, *organization.PocEmail)
	})

	suite.Run("Successful return of linkOnly Organization", func() {
		// Under test:       BuildOrganization
		// Set up:           Pass in a linkOnly organization
		// Expected outcome: No new organization should be created.

		// Check num organization records
		precount, err := suite.DB().Count(&models.Organization{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		organization := BuildOrganization(suite.DB(), []Customization{
			{
				Model: models.Organization{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.Organization{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, organization.ID)

	})
	suite.Run("Successful return of stubbed Organization", func() {
		// Under test:       BuildOrganization
		// Set up:           Pass in a linkOnly organization
		// Expected outcome: No new organization should be created.

		// Check num organization records
		precount, err := suite.DB().Count(&models.Organization{})
		suite.NoError(err)

		// Nil passed in as db
		customName := "Custom Organization"
		organization := BuildOrganization(nil, []Customization{
			{
				Model: models.Organization{
					Name: customName,
				},
			},
		}, nil)
		count, err := suite.DB().Count(&models.Organization{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(customName, organization.Name)
	})
}
