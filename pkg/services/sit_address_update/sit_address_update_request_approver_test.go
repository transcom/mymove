package sitaddressupdate

import (
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	// movefetcher "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *SITAddressUpdateServiceSuite) TestApproveSITAddressUpdateRequest() {
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(query.NewQueryBuilder(), moverouter.NewMoveRouter())
	officeRemarks := "I have chosen to reject this address update request"
	moveRouter := moverouter.NewMoveRouter()

	suite.Run("Successfully Updates the sit address update status to APPROVED", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		sitAddressUpdate := factory.BuildSITAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
			{
				Model: models.SITAddressUpdate{
					ContractorRemarks: models.StringPointer("Moving closer to family"),
					Status:            models.SITAddressUpdateStatusRequested,
				},
			},
		}, []factory.Trait{factory.GetTraitSITAddressUpdateWithMoveSetUp})

		approve := NewSITAddressUpdateRequestApprover(serviceItemUpdater, moveRouter)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		updatedServiceItemPostApproval, err := approve.ApproveSITAddressUpdateRequest(suite.AppContextForTest(), serviceItem.ID, sitAddressUpdate.ID, &officeRemarks, eTag)

		// Checking fields were updated as expected
		suite.NoError(err)
		suite.NotNil(updatedServiceItemPostApproval)

		// Grab the associated move and check its status was properly updated
		// movefetcher := movefetcher.NewMoveTaskOrderFetcher()
		// searchParams := services.MoveTaskOrderFetcherParams{
		// 	IncludeHidden:   false,
		// 	MoveTaskOrderID: serviceItem.MoveTaskOrderID,
		// }
		// updatedMove, moveErr := movefetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)

		// suite.Nil(moveErr)
		// suite.Equal(updatedMove.Status, models.MoveStatusAPPROVED)
	})
}
