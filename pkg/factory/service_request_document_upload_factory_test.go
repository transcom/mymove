package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildServiceRequestDocumentUpload() {
	suite.Run("Successful creation of default ServiceRequestDocumentUpload", func() {
		// Under test:      BuildServiceRequestDocumentUpload
		// Set up:          Create a ServiceRequestDocumentUpload with no customizations or traits
		// Expected outcome: ServiceRequestDocumentUpload should be created with default values

		// SETUP
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), nil, nil)

		suite.False(serviceRequestDocumentUpload.ID.IsNil())
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocument.MTOServiceItemID.IsNil())
		suite.NotNil(serviceRequestDocumentUpload.ServiceRequestDocument.MTOServiceItem)
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocument.MTOServiceItem.ID.IsNil())
	})

	suite.Run("Successful creation of custom ServiceRequestDocumentUpload", func() {
		// Under test:      BuildServiceRequestDocumentUpload
		// Set up:          Create a ServiceRequestDocumentUpload and pass custom fields
		// Expected outcome: ServiceRequestDocumentUpload should be created with custom values

		// SETUP
		customMove := models.Move{
			Locator: "AAAA",
		}
		customMtoServiceItem := models.ServiceRequestDocument{}

		// CALL FUNCTION UNDER TEST
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
			{Model: customMove},
			{Model: customMtoServiceItem},
		}, nil)

		suite.False(serviceRequestDocumentUpload.ID.IsNil())
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocument.MTOServiceItemID.IsNil())
		suite.NotNil(serviceRequestDocumentUpload.ServiceRequestDocument.MTOServiceItem)
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocument.MTOServiceItem.ID.IsNil())
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocumentID.IsNil())
		suite.NotNil(serviceRequestDocumentUpload.ServiceRequestDocument)
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocument.ID.IsNil())
	})

	suite.Run("Successful creation of stubbed ServiceRequestDocumentUpload", func() {
		// Under test:      BuildServiceRequestDocumentUpload
		// Set up:          Create a stubbed ServiceRequestDocumentUpload
		// Expected outcome:No new ServiceRequestDocumentUpload should be created
		precount, err := suite.DB().Count(&models.ServiceRequestDocumentUpload{})
		suite.NoError(err)

		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(nil, nil, nil)

		// VALIDATE RESULTS
		suite.True(serviceRequestDocumentUpload.ID.IsNil())
		suite.True(serviceRequestDocumentUpload.ServiceRequestDocumentID.IsNil())
		suite.NotNil(serviceRequestDocumentUpload.ServiceRequestDocument)
		suite.True(serviceRequestDocumentUpload.ServiceRequestDocument.ID.IsNil())

		// Count how many notification are in the DB, no new
		// notifications should have been created
		count, err := suite.DB().Count(&models.ServiceRequestDocumentUpload{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful return of linkOnly ServiceRequestDocumentUpload", func() {
		// Under test:       BuildServiceRequestDocumentUpload
		// Set up:           Pass in a linkOnly ServiceRequestDocumentUpload
		// Expected outcome: No new ServiceRequestDocumentUpload should be created

		// Check num ServiceRequestDocumentUpload records
		precount, err := suite.DB().Count(&models.ServiceRequestDocumentUpload{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
			{
				Model: models.ServiceRequestDocumentUpload{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.ServiceRequestDocumentUpload{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, serviceRequestDocumentUpload.ID)
	})
}
