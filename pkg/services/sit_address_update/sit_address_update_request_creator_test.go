package sitaddressupdate

import (
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movefetcher "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *SITAddressUpdateServiceSuite) TestCreateSITAddressUpdateRequest() {
	moveRouter := moverouter.NewMoveRouter()
	addressCreator := address.NewAddressCreator()
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(planner, query.NewQueryBuilder(), moveRouter, mtoshipment.NewMTOShipmentFetcher(), addressCreator)
	requestedMockedDistance := 55
	approvedMockedDistance := 45

	suite.Run("Successfully create SIT address update request for approved service item with a REQUESTED status", func() {
		// TESTCASE SCENARIO
		// Under test: CreateSITAddressUpdateRequest function
		// Set up:     We create an approved service item and successfully attempt to create a SITAddressUpdate REQUEST
		// Expected outcome: A SITAddressUpdate should be created with a REQUESTED status
		mockPlanner := &routemocks.Planner{}
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(requestedMockedDistance, nil)

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
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
			{
				Model: models.SITAddressUpdate{
					ContractorRemarks: models.StringPointer("Moving closer to family"),
				},
			},
		}, []factory.Trait{factory.GetTraitSITAddressUpdateWithMoveSetUp})

		creator := NewSITAddressUpdateRequestCreator(mockPlanner, addressCreator, serviceItemUpdater, moveRouter)

		createdAddressUpdateRequest, err := creator.CreateSITAddressUpdateRequest(suite.AppContextForTest(), &sitAddressUpdate)

		suite.NoError(err)
		suite.NotNil(createdAddressUpdateRequest)

		// Distance should exist on the created address request - this is something we calculate on our end
		suite.Equal(requestedMockedDistance, createdAddressUpdateRequest.Distance)

		// Status should be set as requested, since our created update request requires TOO approval due to it being over 50 miles
		suite.Equal(createdAddressUpdateRequest.Status, models.SITAddressUpdateStatusRequested)

		//Checking our set old address matches the final address on the service item the prime is requesting to update
		suite.Equal(createdAddressUpdateRequest.OldAddress.ID, serviceItem.SITDestinationFinalAddress.ID)
		suite.Equal(createdAddressUpdateRequest.OldAddressID, *serviceItem.SITDestinationFinalAddressID)
		suite.Equal(createdAddressUpdateRequest.OldAddress.StreetAddress1, serviceItem.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(createdAddressUpdateRequest.OldAddress.PostalCode, serviceItem.SITDestinationFinalAddress.PostalCode)

		// Contractor Remarks should match
		suite.Equal(*createdAddressUpdateRequest.ContractorRemarks, *sitAddressUpdate.ContractorRemarks)

		// Grabbing move to check the status was updated
		movefetcher := movefetcher.NewMoveTaskOrderFetcher()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: serviceItem.MoveTaskOrderID,
		}
		updatedMove, moveErr := movefetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)

		suite.Nil(moveErr)
		suite.Equal(updatedMove.Status, models.MoveStatusAPPROVALSREQUESTED)
	})

	suite.Run("Successfully create SIT address update request for approved service item with a APPROVED status", func() {
		// TESTCASE SCENARIO
		// Under test: CreateSITAddressUpdateRequest function
		// Set up:     We create an approved service item and successfully attempt to create a SITAddressUpdate REQUEST
		// Expected outcome: A SITAddressUpdate should be created with a REQUESTED status
		mockPlanner := &routemocks.Planner{}
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(approvedMockedDistance, nil)

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
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
			{
				Model: models.SITAddressUpdate{
					ContractorRemarks: models.StringPointer("Moving closer to family"),
				},
			},
		}, []factory.Trait{factory.GetTraitSITAddressUpdateWithMoveSetUp})

		creator := NewSITAddressUpdateRequestCreator(mockPlanner, addressCreator, serviceItemUpdater, moveRouter)

		createdAddressUpdateRequest, err := creator.CreateSITAddressUpdateRequest(suite.AppContextForTest(), &sitAddressUpdate)

		suite.NoError(err)
		suite.NotNil(createdAddressUpdateRequest)

		// Distance should exist on the created address request - this is something we calculate on our end
		suite.Equal(approvedMockedDistance, createdAddressUpdateRequest.Distance)

		// Status should be set as APPROVED, since our created update request requires TOO approval due to it being over 50 miles
		suite.Equal(models.SITAddressUpdateStatusApproved, createdAddressUpdateRequest.Status)

		//Checking our set old address matches the final address on the service item the prime is requesting to update
		suite.Equal(serviceItem.SITDestinationFinalAddress.ID, createdAddressUpdateRequest.OldAddress.ID)
		suite.Equal(*serviceItem.SITDestinationFinalAddressID, createdAddressUpdateRequest.OldAddressID)
		suite.Equal(serviceItem.SITDestinationFinalAddress.StreetAddress1, createdAddressUpdateRequest.OldAddress.StreetAddress1)
		suite.Equal(serviceItem.SITDestinationFinalAddress.PostalCode, createdAddressUpdateRequest.OldAddress.PostalCode)
		sitDestinationFinalAddress := *createdAddressUpdateRequest.MTOServiceItem.SITDestinationFinalAddress
		suite.Equal(createdAddressUpdateRequest.NewAddress.StreetAddress1, sitDestinationFinalAddress.StreetAddress1)
		suite.Equal(createdAddressUpdateRequest.NewAddress.PostalCode, sitDestinationFinalAddress.PostalCode)

		// Contractor Remarks should match
		suite.Equal(*sitAddressUpdate.ContractorRemarks, *createdAddressUpdateRequest.ContractorRemarks)
	})

	suite.Run("Fail to create SIT address update request for unapproved service item", func() {
		// TESTCASE SCENARIO
		// Under test: CreateSITAddressUpdateRequest function
		// Set up:     We create an unapproved service item and fail to create a SITAddressUpdate REQUEST
		// Expected outcome: Failure due to unapproved service item
		mockPlanner := &routemocks.Planner{}
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(requestedMockedDistance, nil)

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusRejected,
				},
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

		sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
			{
				Model: models.SITAddressUpdate{
					ContractorRemarks: models.StringPointer("Moving closer to family"),
				},
			},
		}, []factory.Trait{factory.GetTraitSITAddressUpdateWithMoveSetUp})

		creator := NewSITAddressUpdateRequestCreator(mockPlanner, addressCreator, serviceItemUpdater, moveRouter)

		createdAddressUpdateRequest, err := creator.CreateSITAddressUpdateRequest(suite.AppContextForTest(), &sitAddressUpdate)

		suite.Error(err)
		suite.Nil(createdAddressUpdateRequest)
	})

	suite.Run("Fail to create SIT address update request for service item of incorrect type", func() {
		// TESTCASE SCENARIO
		// Under test: CreateSITAddressUpdateRequest function
		// Set up:     We create a service item with the wrong type and fail to create a SITAddressUpdate REQUEST
		// Expected outcome: Failure due to wrong type of service item
		mockPlanner := &routemocks.Planner{}
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(requestedMockedDistance, nil)

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
			{
				Model: models.SITAddressUpdate{
					ContractorRemarks: models.StringPointer("Moving closer to family"),
				},
			},
		}, []factory.Trait{factory.GetTraitSITAddressUpdateWithMoveSetUp})

		creator := NewSITAddressUpdateRequestCreator(mockPlanner, addressCreator, serviceItemUpdater, moveRouter)

		createdAddressUpdateRequest, err := creator.CreateSITAddressUpdateRequest(suite.AppContextForTest(), &sitAddressUpdate)

		suite.Error(err)
		suite.Nil(createdAddressUpdateRequest)
	})
}
