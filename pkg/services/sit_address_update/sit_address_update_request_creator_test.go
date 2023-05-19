package sitaddressupdate

import (
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
)

func (suite *SITAddressUpdateServiceSuite) TestCreateSITAddressUpdateRequest() {
	addressCreator := address.NewAddressCreator()
	mockPlanner := &routemocks.Planner{}
	mockedDistance := 55
	mockPlanner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*models.Address"),
		mock.AnythingOfType("*models.Address"),
	).Return(mockedDistance, nil)

	suite.Run("Successfully create SIT address update request for approved service item", func() {
		// TESTCASE SCENARIO
		// Under test: CreateSITAddressUpdateRequest function
		// Set up:     We create an approved service item and successfully attempt to create a SITAddressUpdate REQUEST
		// Expected outcome: A SITAddressUpdate should be created with a REQUESTED status

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

		creator := NewSITAddressUpdateRequestCreator(mockPlanner, addressCreator)

		createdAddressUpdateRequest, err := creator.CreateSITAddressUpdateRequest(suite.AppContextForTest(), &sitAddressUpdate)

		suite.NoError(err)
		suite.NotNil(createdAddressUpdateRequest)

		// Distance should exist on the created address request - this is something we calculate on our end
		suite.Equal(mockedDistance, createdAddressUpdateRequest.Distance)

		// Status should be set as requested, since our created update request requires TOO approval due to it being over 50 miles
		suite.Equal(createdAddressUpdateRequest.Status, models.SITAddressUpdateStatusRequested)

		//Checking our set old address matches the final address on the service item the prime is requesting to update
		suite.Equal(createdAddressUpdateRequest.OldAddress.ID, serviceItem.SITDestinationFinalAddress.ID)
		suite.Equal(createdAddressUpdateRequest.OldAddressID, *serviceItem.SITDestinationFinalAddressID)
		suite.Equal(createdAddressUpdateRequest.OldAddress.StreetAddress1, serviceItem.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(createdAddressUpdateRequest.OldAddress.PostalCode, serviceItem.SITDestinationFinalAddress.PostalCode)

		// Contractor Remarks should match
		suite.Equal(*createdAddressUpdateRequest.ContractorRemarks, *sitAddressUpdate.ContractorRemarks)
	})

	suite.Run("Failed to create SIT address update request for unapproved service item", func() {
		// TESTCASE SCENARIO
		// Under test: CreateSITAddressUpdateRequest function
		// Set up:     We create an unapproved service item and fail to create a SITAddressUpdate REQUEST
		// Expected outcome: Failure due to unapproved service item

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

		creator := NewSITAddressUpdateRequestCreator(mockPlanner, addressCreator)

		createdAddressUpdateRequest, err := creator.CreateSITAddressUpdateRequest(suite.AppContextForTest(), &sitAddressUpdate)

		suite.Error(err)
		suite.Nil(createdAddressUpdateRequest)
	})
}
