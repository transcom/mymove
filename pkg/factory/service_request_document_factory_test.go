package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildServiceRequestDocument() {
	suite.Run("Successful creation of default ServiceRequestDocument", func() {
		// Under test:      BuildServiceRequestDocument
		// Set up:          Create a ServiceRequestDocument with no customizations or traits
		// Expected outcome: ServiceRequestDocument should be created with default values

		// SETUP
		serviceRequestDocument := BuildServiceRequestDocument(suite.DB(), nil, nil)

		suite.False(serviceRequestDocument.ID.IsNil())
		suite.False(serviceRequestDocument.MTOServiceItem.MoveTaskOrderID.IsNil())
		suite.NotNil(serviceRequestDocument.MTOServiceItem.MoveTaskOrder)
		suite.False(serviceRequestDocument.MTOServiceItem.MoveTaskOrder.ID.IsNil())
	})

	suite.Run("Successful creation of custom ServiceRequestDocument", func() {
		// Under test:      BuildServiceRequestDocument
		// Set up:          Create a ServiceRequestDocument and pass custom fields
		// Expected outcome: ServiceRequestDocument should be created with custom values

		// SETUP
		customMove := models.Move{
			Locator: "AAAA",
		}
		customMtoServiceItem := models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		}

		// CALL FUNCTION UNDER TEST
		serviceRequestDocument := BuildServiceRequestDocument(suite.DB(), []Customization{
			{Model: customMove},
			{Model: customMtoServiceItem},
		}, nil)

		suite.False(serviceRequestDocument.ID.IsNil())
		suite.False(serviceRequestDocument.MTOServiceItem.MoveTaskOrderID.IsNil())
		suite.NotNil(serviceRequestDocument.MTOServiceItem.MoveTaskOrder)
		suite.False(serviceRequestDocument.MTOServiceItem.MoveTaskOrder.ID.IsNil())
		suite.Equal(customMove.Locator, serviceRequestDocument.MTOServiceItem.MoveTaskOrder.Locator)
		suite.False(serviceRequestDocument.MTOServiceItemID.IsNil())
		suite.NotNil(serviceRequestDocument.MTOServiceItem)
		suite.False(serviceRequestDocument.MTOServiceItem.ID.IsNil())
		suite.Equal(customMtoServiceItem.Status, serviceRequestDocument.MTOServiceItem.Status)
	})

	suite.Run("Successful creation of stubbed ServiceRequestDocument", func() {
		// Under test:      BuildServiceRequestDocument
		// Set up:          Create a stubbed ServiceRequestDocument
		// Expected outcome:No new ServiceRequestDocument should be created
		precount, err := suite.DB().Count(&models.ServiceRequestDocument{})
		suite.NoError(err)

		serviceRequestDocument := BuildServiceRequestDocument(nil, nil, nil)

		// VALIDATE RESULTS
		suite.True(serviceRequestDocument.ID.IsNil())
		suite.True(serviceRequestDocument.MTOServiceItemID.IsNil())
		suite.NotNil(serviceRequestDocument.MTOServiceItem)
		suite.True(serviceRequestDocument.MTOServiceItem.ID.IsNil())

		// Count how many notification are in the DB, no new
		// notifications should have been created
		count, err := suite.DB().Count(&models.ServiceRequestDocument{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful return of linkOnly ServiceRequestDocument", func() {
		// Under test:       BuildServiceRequestDocument
		// Set up:           Pass in a linkOnly ServiceRequestDocument
		// Expected outcome: No new ServiceRequestDocument should be created

		// Check num ServiceRequestDocument records
		precount, err := suite.DB().Count(&models.ServiceRequestDocument{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		serviceRequestDocument := BuildServiceRequestDocument(suite.DB(), []Customization{
			{
				Model: models.ServiceRequestDocument{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.ServiceRequestDocument{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, serviceRequestDocument.ID)
	})
}
