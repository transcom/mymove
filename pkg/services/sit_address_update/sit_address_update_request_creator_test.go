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

	suite.Run("Successfully create SIT update request", func() {
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

		sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitSITAddressUpdateWithMoveSetUp})

		creator := NewSITAddressUpdateRequestCreator(mockPlanner, addressCreator)

		createdAddressUpdate, err := creator.CreateSITAddressUpdateRequest(suite.AppContextForTest(), &sitAddressUpdate)

		suite.NoError(err)
		suite.NotNil(createdAddressUpdate)
		suite.Equal(mockedDistance, createdAddressUpdate.Distance)
		suite.Equal(createdAddressUpdate.Status, models.SITAddressUpdateStatusRequested)
		suite.Equal(createdAddressUpdate.OldAddress, serviceItem.SITDestinationFinalAddress)
		suite.Equal(createdAddressUpdate.OldAddressID, serviceItem.SITDestinationFinalAddressID)
	})
}
