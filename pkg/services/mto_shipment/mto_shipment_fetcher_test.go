package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestListMTOShipments() {
	mtoShipmentFetcher := NewMTOShipmentFetcher()

	suite.Run("Returns not found error when move id doesn't exist", func() {
		moveID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(moveID, "move not found")

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(suite.AppContextForTest(), moveID)

		suite.Equalf(err, expectedError, "Expected not found error for non-existent move id")
		suite.Nil(mtoShipments, "Expected shipment slice to be nil")
	})

	suite.Run("Returns an empty shipment list when no shipments exist", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(suite.AppContextForTest(), move.ID)

		suite.NoError(err, "Expected no error for a move without shipments")
		suite.Len(mtoShipments, 0, "Expected a zero length shipment list")
	})

	suite.Run("Returns external vendor shipments last", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		externalVendorShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor: true,
				},
			},
		}, nil)
		firstShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		secondShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(suite.AppContextForTest(), move.ID)

		suite.NoError(err, "Expected no error for a move with 3 shipments")
		suite.Len(mtoShipments, 3, "Expected a shipment list of length 3")

		suite.Equal(firstShipment.ID.String(), mtoShipments[0].ID.String())
		suite.Equal(secondShipment.ID.String(), mtoShipments[1].ID.String())
		suite.Equal(externalVendorShipment.ID.String(), mtoShipments[2].ID.String())

	})

	suite.Run("Returns multiple shipments for move ordered by created date", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		firstShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		secondShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(suite.AppContextForTest(), move.ID)

		suite.NoError(err, "Expected no error for a move with two shipments")
		suite.Len(mtoShipments, 2, "Expected a shipment list of length 2")

		suite.Equal(firstShipment.ID.String(), mtoShipments[0].ID.String())
		suite.Equal(secondShipment.ID.String(), mtoShipments[1].ID.String())

	})

	suite.Run("Returns only non-deleted shipments", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					DeletedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		secondShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(suite.AppContextForTest(), move.ID)

		suite.NoError(err, "Expected no error for a move with one deleted and one not deleted shipment")
		suite.Len(mtoShipments, 1, "Expected a shipment list of length 1")

		suite.Equal(secondShipment.ID.String(), mtoShipments[0].ID.String())

	})

	suite.Run("Loads all shipment associations", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		secondaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    secondaryPickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model:    secondaryDeliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		shipmentAddressUpdate := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})

		serviceItem := testdatagen.MakeMTOServiceItemDomesticCrating(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDCRT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		agents := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)
		SITExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovedSITDurationUpdate})

		reweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			MTOShipment: shipment,
		})

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(suite.AppContextForTest(), move.ID)

		suite.NoError(err, "Expected no error for a move with shipment associations")
		suite.Len(mtoShipments, 1, "Expected a single shipment with associations")

		actualShipment := mtoShipments[0]

		suite.Equal(serviceItem.ReService.Code, actualShipment.MTOServiceItems[0].ReService.Code)
		suite.Equal(agents.ID.String(), actualShipment.MTOAgents[0].ID.String())
		suite.Equal(shipment.PickupAddress.ID.String(), actualShipment.PickupAddress.ID.String())
		suite.Equal(secondaryPickupAddress.ID.String(), actualShipment.SecondaryPickupAddress.ID.String())
		suite.Equal(shipment.DestinationAddress.ID.String(), actualShipment.DestinationAddress.ID.String())
		suite.Equal(secondaryDeliveryAddress.ID.String(), actualShipment.SecondaryDeliveryAddress.ID.String())
		suite.Len(actualShipment.MTOServiceItems[0].Dimensions, 2)
		suite.Equal(SITExtension.ID.String(), actualShipment.SITDurationUpdates[0].ID.String())
		suite.Equal(reweigh.ID.String(), actualShipment.Reweigh.ID.String())
		suite.Equal(shipmentAddressUpdate.ID.String(), actualShipment.DeliveryAddressUpdate.ID.String())
		suite.Equal(shipmentAddressUpdate.NewAddress.ID.String(), actualShipment.DeliveryAddressUpdate.NewAddress.ID.String())
		suite.Equal(shipmentAddressUpdate.OriginalAddress.ID.String(), actualShipment.DeliveryAddressUpdate.OriginalAddress.ID.String())
	})

	suite.Run("Loads PPM associations", func() {
		// not reusing the test above because the fetcher only loads PPM associations if the shipment type is PPM
		move := factory.BuildMove(suite.DB(), nil, nil)
		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    move.Orders.ServiceMember,
				LinkOnly: true,
			},
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		userUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    move.Orders.ServiceMember,
				LinkOnly: true,
			},
		}, nil)

		movingExpense := &models.MovingExpense{
			PPMShipmentID: ppmShipment.ID,
			Document:      userUpload.Document,
			DocumentID:    userUpload.Document.ID,
		}

		proGear := &models.ProgearWeightTicket{
			PPMShipmentID: ppmShipment.ID,
			Document:      userUpload.Document,
			DocumentID:    userUpload.Document.ID,
		}

		err := suite.DB().Create(movingExpense)
		suite.NoError(err)

		err = suite.DB().Create(proGear)
		suite.NoError(err)

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(suite.AppContextForTest(), move.ID)
		suite.NoError(err)

		actualPPMShipment := mtoShipments[0].PPMShipment

		suite.NotNil(actualPPMShipment)
		suite.Equal(ppmShipment.ID.String(), actualPPMShipment.ID.String())
		suite.Equal(ppmShipment.ShipmentID.String(), mtoShipments[0].ID.String())

		suite.Len(actualPPMShipment.WeightTickets, 1)
		suite.Len(actualPPMShipment.WeightTickets[0].EmptyDocument.UserUploads, 1)
		suite.Len(actualPPMShipment.WeightTickets[0].FullDocument.UserUploads, 1)

		suite.Len(actualPPMShipment.MovingExpenses, 1)
		suite.Len(actualPPMShipment.MovingExpenses[0].Document.UserUploads, 1)

		suite.Len(actualPPMShipment.ProgearWeightTickets, 1)
		suite.Len(actualPPMShipment.ProgearWeightTickets[0].Document.UserUploads, 1)
	})
}

func (suite *MTOShipmentServiceSuite) TestGetMTOShipment() {
	mtoShipmentFetcher := NewMTOShipmentFetcher()

	// Test successful fetch
	suite.Run("Returns a shipment successfully with correct ID", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		fetchedShipment, err := mtoShipmentFetcher.GetShipment(suite.AppContextForTest(), shipment.ID)
		suite.NoError(err)
		suite.Equal(shipment.ID, fetchedShipment.ID)
	})

	// Test 404 fetch
	suite.Run("Returns not found error when shipment id doesn't exist", func() {
		shipmentID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(shipmentID, "while looking for shipment")

		mtoShipment, err := mtoShipmentFetcher.GetShipment(suite.AppContextForTest(), shipmentID)

		suite.Nil(mtoShipment)
		suite.Equalf(err, expectedError, "while looking for shipment")
	})
}

func (suite *MTOShipmentServiceSuite) TestFindMTOShipment() {
	// Test successful fetch
	suite.Run("Returns a shipment successfully with correct ID", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		fetchedShipment, err := FindShipment(suite.AppContextForTest(), shipment.ID)
		suite.NoError(err)
		suite.Equal(shipment.ID, fetchedShipment.ID)
	})

	// Test 404 fetch
	suite.Run("Returns not found error when shipment id doesn't exist", func() {
		shipmentID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(shipmentID, "while looking for shipment")

		mtoShipment, err := FindShipment(suite.AppContextForTest(), shipmentID)

		suite.Nil(mtoShipment)
		suite.Equalf(err, expectedError, "while looking for shipment")
	})

	suite.Run("404 Not Found Error - shipment can only be created for service member associated with the current session", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		})

		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		mtoShipment, err := FindShipment(session, shipment.ID)
		suite.Error(err)
		suite.Nil(mtoShipment)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
