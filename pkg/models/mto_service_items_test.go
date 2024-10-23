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

		validMTOServiceItem := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrderID,
			MTOShipmentID:   &mtoShipmentID,
			ReServiceID:     reServiceID,
			Status:          models.MTOServiceItemStatusSubmitted,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOServiceItem, expErrors)
	})
	suite.Run("can add an OCONUS ReService to an OCONUS shipment", func() {
		moveTaskOrderID := uuid.Must(uuid.NewV4())
		reServiceID := uuid.Must(uuid.NewV4())

		gb := factory.BuildCountry(suite.DB(), []factory.Customization{
			{
				Model: models.Country{
					Country:     "GB",
					CountryName: "Great Britain",
				},
			},
		}, nil)

		oconusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus: models.BoolPointer(true),
				},
			},
			{
				Model:    gb,
				LinkOnly: true,
			},
		}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ID:                   uuid.Must(uuid.NewV4()),
					PickupAddressID:      &oconusAddress.ID,
					DestinationAddressID: &oconusAddress.ID,
				},
			},
		}, nil)

		oconusReService := models.ReService{
			Code: models.ReServiceCodeIOFSIT,
		}

		oconusMTOServiceItem := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrderID,
			MTOShipmentID:   &shipment.ID,
			ReServiceID:     reServiceID,
			ReService:       oconusReService,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		verrs, err := models.ValidateableModel.Validate(&oconusMTOServiceItem, suite.DB())
		suite.False(verrs.HasAny())
		suite.Nil(err)
	})

	suite.Run("cannot add a CONUS ReService to an OCONUS shipment", func() {
		moveTaskOrderID := uuid.Must(uuid.NewV4())
		reServiceID := uuid.Must(uuid.NewV4())

		gb := factory.BuildCountry(suite.DB(), []factory.Customization{
			{
				Model: models.Country{
					Country:     "GB",
					CountryName: "Great Britain",
				},
			},
		}, nil)
		oconusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus: models.BoolPointer(true),
				},
			},
			{
				Model:    gb,
				LinkOnly: true,
			},
		}, nil)

		// Create shipment with OCONUS destination address
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ID:              uuid.Must(uuid.NewV4()),
					PickupAddressID: &oconusAddress.ID,
				},
			},
		}, nil)

		// CONUS ReService Code
		domesticReService := models.ReService{
			Code: models.ReServiceCodeDOFSIT,
		}

		conusMTOServiceItem := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrderID,
			MTOShipmentID:   &shipment.ID,
			ReServiceID:     reServiceID,
			ReService:       domesticReService,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		verrs, err := models.ValidateableModel.Validate(&conusMTOServiceItem, suite.DB())
		suite.True(verrs.HasAny())
		suite.Contains(verrs.Get("ReService"), "A domestic ReService Code cannot be applied to an OCONUS shipment")
		suite.Nil(err)
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
