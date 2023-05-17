package sitaddressupdate

import (
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *SITAddressUpdateServiceSuite) TestRejectSITAddressUpdateRequest() {
	officeRemarks := "I have chosen to reject this address request"

	suite.Run("Successfully reject SIT address update request and update it's status", func() {
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)

		sitAddressUpdate := factory.BuildSITAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, nil)

		rejector := NewSITAddressUpdateRequestRejector()

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		updatedSITAddressUpdateRequestPostRejection, err := rejector.RejectSITAddressUpdateRequest(suite.AppContextForTest(), serviceItem.ID, sitAddressUpdate.ID, &officeRemarks, eTag)

		suite.NoError(err)
		suite.NotNil(updatedSITAddressUpdateRequestPostRejection)
		suite.Equal(updatedSITAddressUpdateRequestPostRejection.Status, models.SITAddressUpdateStatusRejected)
		suite.Equal(updatedSITAddressUpdateRequestPostRejection.OfficeRemarks, officeRemarks)
	})
}
