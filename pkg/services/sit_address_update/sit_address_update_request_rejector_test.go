package sitaddressupdate

import (
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movefetcher "github.com/transcom/mymove/pkg/services/move_task_order"
)

func (suite *SITAddressUpdateServiceSuite) TestRejectSITAddressUpdateRequest() {
	officeRemarks := "I have chosen to reject this address update request"
	blankOfficeRemarks := ""
	moveRouter := moverouter.NewMoveRouter()
	reject := NewSITAddressUpdateRequestRejector(moveRouter)

	suite.Run("Successfully Updates the sit address update status to REJECTED", func() {
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

		eTag := etag.GenerateEtag(sitAddressUpdate.UpdatedAt)
		updatedSITAddressUpdateRequestPostRejection, err := reject.RejectSITAddressUpdateRequest(suite.AppContextForTest(), sitAddressUpdate.ID, &officeRemarks, eTag)

		// Checking fields were updated as expected
		suite.NoError(err)
		suite.NotNil(updatedSITAddressUpdateRequestPostRejection)
		suite.Equal(updatedSITAddressUpdateRequestPostRejection.Status, models.SITAddressUpdateStatusRejected)
		suite.Equal(updatedSITAddressUpdateRequestPostRejection.OfficeRemarks, &officeRemarks)

		// Timestamp should differ after update
		suite.NotEqual(sitAddressUpdate.UpdatedAt, updatedSITAddressUpdateRequestPostRejection.UpdatedAt)

		// Grab the associated move and check its status was properly updated
		movefetcher := movefetcher.NewMoveTaskOrderFetcher()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: serviceItem.MoveTaskOrderID,
		}
		updatedMove, moveErr := movefetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)

		suite.Nil(moveErr)
		suite.Equal(updatedMove.Status, models.MoveStatusAPPROVED)
	})

	suite.Run("Fails to REJECT SIT address update due to missing remarks", func() {
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

		eTag := etag.GenerateEtag(sitAddressUpdate.UpdatedAt)
		updatedSITAddressUpdateRequestPostRejection, err := reject.RejectSITAddressUpdateRequest(suite.AppContextForTest(), sitAddressUpdate.ID, &blankOfficeRemarks, eTag)

		suite.Error(err)
		suite.Nil(updatedSITAddressUpdateRequestPostRejection)
		suite.ErrorContains(err, "OfficeRemarks can not be blank.")
	})
}
