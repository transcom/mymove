package sitaddressupdate

import (
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *SITAddressUpdateServiceSuite) TestCreateApprovedSITAddressUpdate() {
	addressCreator := address.NewAddressCreator()
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(query.NewQueryBuilder(), moverouter.NewMoveRouter())
	mockPlanner := &routemocks.Planner{}
	mockedDistance := 55
	mockPlanner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*models.Address"),
		mock.AnythingOfType("*models.Address"),
	).Return(mockedDistance, nil)

	suite.Run("Successfully create SITAddressUpdate", func() {
		// TESTCASE SCENARIO
		// Under test: CreateApprovedSITAddressUpdate function
		// Set up:     We create an approved service item and attempt to create a SITAddressUpdate
		// Expected outcome:
		//             SITAddressUpdate is created and SITDestinationFinalAddress on the service item is updated

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
					OfficeRemarks: models.StringPointer("office remarks"),
				},
			},
		}, []factory.Trait{factory.GetTraitSITAddressUpdateWithMoveSetUp})

		creator := NewApprovedOfficeSITAddressUpdateCreator(mockPlanner, addressCreator, serviceItemUpdater)

		createdAddressUpdate, err := creator.CreateApprovedSITAddressUpdate(suite.AppContextForTest(), &sitAddressUpdate)
		suite.NoError(err)
		suite.NotNil(createdAddressUpdate)
		suite.Equal(mockedDistance, createdAddressUpdate.Distance)
		suite.Equal(models.SITAddressUpdateStatusApproved, createdAddressUpdate.Status)
		suite.Equal(*sitAddressUpdate.OfficeRemarks, *createdAddressUpdate.OfficeRemarks)
		suite.Equal(*sitAddressUpdate.OfficeRemarks, *createdAddressUpdate.OfficeRemarks)
		suite.Equal(*serviceItem.SITDestinationFinalAddressID, createdAddressUpdate.OldAddressID)
		suite.Equal(serviceItem.SITDestinationFinalAddress.ID, createdAddressUpdate.OldAddress.ID)
		suite.Equal(serviceItem.SITDestinationFinalAddress.StreetAddress1, createdAddressUpdate.OldAddress.StreetAddress1)
		suite.Equal(serviceItem.SITDestinationFinalAddress.PostalCode, createdAddressUpdate.OldAddress.PostalCode)
		suite.Equal(sitAddressUpdate.NewAddress.StreetAddress1, createdAddressUpdate.NewAddress.StreetAddress1)
		suite.Equal(sitAddressUpdate.NewAddress.PostalCode, createdAddressUpdate.NewAddress.PostalCode)
		suite.Equal(sitAddressUpdate.MTOServiceItemID, createdAddressUpdate.MTOServiceItemID)
		suite.Equal(sitAddressUpdate.MTOServiceItem.ID, createdAddressUpdate.MTOServiceItem.ID)
		suite.Equal(sitAddressUpdate.MTOServiceItem.ReServiceID, createdAddressUpdate.MTOServiceItem.ReServiceID)
		suite.Equal(sitAddressUpdate.MTOServiceItem.ReService.Code, createdAddressUpdate.MTOServiceItem.ReService.Code)
		sitDestinationFinalAddress := *createdAddressUpdate.MTOServiceItem.SITDestinationFinalAddress
		suite.Equal(createdAddressUpdate.NewAddress.StreetAddress1, sitDestinationFinalAddress.StreetAddress1)
		suite.Equal(createdAddressUpdate.NewAddress.PostalCode, sitDestinationFinalAddress.PostalCode)
	})

	suite.Run("Error creating SITAddressUpdate", func() {
		// TESTCASE SCENARIO
		// Under test: CreateApprovedSITAddressUpdate function
		// Set up:     We create an unapproved service item and attempt to create a SITAddressUpdate
		// Expected outcome:
		//             Error because we cannot create a SITAddressUpdate before service item is approved

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
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

		oldAddress := factory.BuildAddress(suite.DB(), nil, nil)
		sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
			{
				Model:    oldAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITAddressUpdateOldAddress,
			},
			{
				Model: models.SITAddressUpdate{
					OfficeRemarks: models.StringPointer("office remarks"),
				},
			},
		}, []factory.Trait{factory.GetTraitSITAddressUpdateWithMoveSetUp})

		creator := NewApprovedOfficeSITAddressUpdateCreator(mockPlanner, addressCreator, serviceItemUpdater)

		createdSITAddressUpdate, err := creator.CreateApprovedSITAddressUpdate(suite.AppContextForTest(), &sitAddressUpdate)
		suite.Error(err)
		suite.Nil(createdSITAddressUpdate)
	})
}
