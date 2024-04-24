package mtoserviceitem

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *MTOServiceItemServiceSuite) TestGetMTOServiceItem() {
	mtoServiceItemFetcher := NewMTOServiceItemFetcher()

	// Test successful fetch
	suite.Run("Returns a service item successfully with correct ID", func() {
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)

		fetchedServiceItem, err := mtoServiceItemFetcher.GetServiceItem(suite.AppContextForTest(), serviceItem.ID)
		suite.NoError(err)
		suite.Equal(serviceItem.ID, fetchedServiceItem.ID)
	})

	// Test 404 fetch
	suite.Run("Returns not found error when shipment id doesn't exist", func() {
		serviceItemID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(serviceItemID, "while looking for service item")

		mtoServiceItem, err := mtoServiceItemFetcher.GetServiceItem(suite.AppContextForTest(), serviceItemID)

		suite.Nil(mtoServiceItem)
		suite.Equalf(err, expectedError, "while looking for service item")
	})
}
