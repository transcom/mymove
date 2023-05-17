package sitaddressupdate

import (
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *SITAddressUpdateServiceSuite) TestApproveSITAddressUpdateRequest() {
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(query.NewQueryBuilder(), moverouter.NewMoveRouter())
	officeRemarks := "I have chosen to approve this address request"

	suite.Run("Successfully approve SIT address update request and update service item address", func() {
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

		approver := NewSITAddressUpdateRequestApprover(serviceItemUpdater)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		updatedServiceItemPostApproval, err := approver.ApproveSITAddressUpdateRequest(suite.AppContextForTest(), serviceItem.ID, sitAddressUpdate.ID, &officeRemarks, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItemPostApproval)
		suite.Equal(updatedServiceItemPostApproval.SITDestinationFinalAddress, sitAddressUpdate.NewAddress)
		suite.Equal(updatedServiceItemPostApproval.SITDestinationFinalAddressID, sitAddressUpdate.NewAddressID)
	})
}
