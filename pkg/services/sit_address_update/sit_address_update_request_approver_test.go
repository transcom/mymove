package sitaddressupdate

import (
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movefetcher "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	query "github.com/transcom/mymove/pkg/services/query"
)

func (suite *SITAddressUpdateServiceSuite) TestApproveSITAddressUpdateRequest() {
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(
		query.NewQueryBuilder(),
		moverouter.NewMoveRouter(),
		mtoshipment.NewMTOShipmentFetcher(),
	)
	moveRouter := moverouter.NewMoveRouter()
	approve := NewSITAddressUpdateRequestApprover(serviceItemUpdater, moveRouter)
	officeRemarks := "I have chosen to approve this address update request"
	blankOfficeRemarks := ""

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
		updatedServiceItemPostApproval, err := approve.ApproveSITAddressUpdateRequest(suite.AppContextForTest(), sitAddressUpdate.ID, &officeRemarks, eTag)

		// Checking fields were updated as expected
		suite.NoError(err)
		suite.NotNil(updatedServiceItemPostApproval)

		// Make sure fields were updated
		suite.NotEqual(serviceItem.SITDestinationFinalAddressID, updatedServiceItemPostApproval.SITDestinationFinalAddressID)
		suite.NotEqual(serviceItem.SITDestinationFinalAddress, updatedServiceItemPostApproval.SITDestinationFinalAddress)
		suite.NotEqual(serviceItem.UpdatedAt, updatedServiceItemPostApproval.UpdatedAt)

		// Ensure updated field values are correct
		suite.Equal(sitAddressUpdate.NewAddressID.String(), updatedServiceItemPostApproval.SITDestinationFinalAddressID.String())
		suite.Equal(sitAddressUpdate.NewAddress.StreetAddress1, updatedServiceItemPostApproval.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(sitAddressUpdate.NewAddress.PostalCode, updatedServiceItemPostApproval.SITDestinationFinalAddress.PostalCode)

		// Grab the associated sit address update request and check it was properly updated
		updatedSitAddressUpdate, addressUpdateErr := models.FetchSITAddressUpdate(suite.AppContextForTest().DB(), sitAddressUpdate.ID)

		suite.Nil(addressUpdateErr)
		suite.NotNil(updatedSitAddressUpdate)
		suite.Equal(models.SITAddressUpdateStatusApproved, updatedSitAddressUpdate.Status)
		suite.Equal(officeRemarks, *updatedSitAddressUpdate.OfficeRemarks)
		suite.NotEqual(updatedSitAddressUpdate.UpdatedAt, sitAddressUpdate.UpdatedAt)

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

	suite.Run("Fails to approve SIT address update request due to missing remarks", func() {
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
		updatedServiceItemPostApproval, err := approve.ApproveSITAddressUpdateRequest(suite.AppContextForTest(), sitAddressUpdate.ID, &blankOfficeRemarks, eTag)

		suite.Error(err)
		suite.Nil(updatedServiceItemPostApproval)
		suite.ErrorContains(err, "OfficeRemarks can not be blank.")
	})
}
