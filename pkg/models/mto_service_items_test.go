package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMTOServiceItemValidation() {
	suite.Run("test valid MTOServiceItem", func() {
		moveTaskOrderID := uuid.Must(uuid.NewV4())
		mtoShipmentID := uuid.Must(uuid.NewV4())
		reServiceID := uuid.Must(uuid.NewV4())
		poeLocationID := uuid.Must(uuid.NewV4())
		podLocationID := uuid.Must(uuid.NewV4())

		validMTOServiceItem := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrderID,
			MTOShipmentID:   &mtoShipmentID,
			ReServiceID:     reServiceID,
			Status:          models.MTOServiceItemStatusSubmitted,
			POELocationID:   &poeLocationID,
			PODLocationID:   &podLocationID,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOServiceItem, expErrors)
	})
}

func (suite *ModelSuite) TestFetchRelatedDestinationSITServiceItems() {
	suite.Run("successful fetch of destination service items", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		ddfServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)
		dddServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		relatedServiceItems, fetchErr := models.FetchRelatedDestinationSITServiceItems(suite.DB(), dddServiceItem.ID)
		suite.NoError(fetchErr)
		suite.Len(relatedServiceItems, 2, "There should be two related service items")
		foundDDF := false
		foundDDD := false
		for _, serviceItem := range relatedServiceItems {
			if serviceItem.ID == ddfServiceItem.ID {
				foundDDF = true
			}
			if serviceItem.ID == dddServiceItem.ID {
				foundDDD = true
			}
		}
		suite.True(foundDDF)
		suite.True(foundDDD)
	})
	suite.Run("successful fetch of destination service items", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		msServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeMS,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeCS,
				},
			},
		}, nil)
		relatedServiceItems, fetchErr := models.FetchRelatedDestinationSITServiceItems(suite.DB(), msServiceItem.ID)
		suite.NoError(fetchErr)
		suite.Len(relatedServiceItems, 0, "There should be zero related destination service items")
	})
}
