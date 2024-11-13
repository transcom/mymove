package mtoshipment

import (
	"fmt"
	"math"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	notificationMocks "github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mockservices "github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func setUpMockNotificationSender() notifications.NotificationSender {
	// The NewMTOShipmentUpdater needs a NotificationSender for sending notification emails to the customer.
	// This function allows us to set up a fresh mock for each test so we can check the number of calls it has.
	mockSender := notificationMocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*notifications.ReweighRequested"),
	).Return(nil)

	return &mockSender
}

func (suite *MTOShipmentServiceSuite) TestMTOShipmentUpdater() {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(1000, nil)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester())
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
	mockSender := setUpMockNotificationSender()
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()

	mtoShipmentUpdaterOffice := NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)
	mtoShipmentUpdaterCustomer := NewCustomerMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)
	mtoShipmentUpdaterPrime := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)
	scheduledPickupDate := now.Add(time.Hour * 24 * 3)
	firstAvailableDeliveryDate := now.Add(time.Hour * 24 * 4)
	actualPickupDate := now.Add(time.Hour * 24 * 3)
	scheduledDeliveryDate := now.Add(time.Hour * 24 * 4)
	actualDeliveryDate := now.Add(time.Hour * 24 * 4)
	primeActualWeight := unit.Pound(1234)
	primeEstimatedWeight := unit.Pound(1234)

	var mtoShipment models.MTOShipment
	var oldMTOShipment models.MTOShipment
	var secondaryPickupAddress models.Address
	var secondaryDeliveryAddress models.Address
	var tertiaryPickupAddress models.Address
	var tertiaryDeliveryAddress models.Address
	var newDestinationAddress models.Address
	var newPickupAddress models.Address

	setupTestData := func() {
		oldMTOShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
					ScheduledPickupDate:        &scheduledPickupDate,
					ApprovedDate:               &firstAvailableDeliveryDate,
				},
			},
		}, nil)

		requestedPickupDate := *oldMTOShipment.RequestedPickupDate
		secondaryDeliveryAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
		secondaryPickupAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		tertiaryPickupAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		tertiaryDeliveryAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
		newDestinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "987 Other Avenue",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50309",
					County:         "POLK",
				},
			},
		}, nil)

		newPickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "987 Over There Avenue",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Houston",
					State:          "TX",
					PostalCode:     "77083",
				},
			},
		}, []factory.Trait{factory.GetTraitAddress4})

		mtoShipment = models.MTOShipment{
			ID:                         oldMTOShipment.ID,
			MoveTaskOrderID:            oldMTOShipment.MoveTaskOrderID,
			MoveTaskOrder:              oldMTOShipment.MoveTaskOrder,
			DestinationAddress:         oldMTOShipment.DestinationAddress,
			DestinationAddressID:       oldMTOShipment.DestinationAddressID,
			PickupAddress:              oldMTOShipment.PickupAddress,
			PickupAddressID:            oldMTOShipment.PickupAddressID,
			RequestedPickupDate:        &requestedPickupDate,
			ScheduledPickupDate:        &scheduledPickupDate,
			ShipmentType:               "UNACCOMPANIED_BAGGAGE",
			PrimeActualWeight:          &primeActualWeight,
			PrimeEstimatedWeight:       &primeEstimatedWeight,
			FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
			Status:                     oldMTOShipment.Status,
			ActualPickupDate:           &actualPickupDate,
			ApprovedDate:               &firstAvailableDeliveryDate,
			MarketCode:                 oldMTOShipment.MarketCode,
		}

		primeEstimatedWeight = unit.Pound(9000)
	}

	suite.Run("Etag is stale", func() {
		setupTestData()

		eTag := etag.GenerateEtag(time.Now())

		session := auth.Session{}
		_, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &mtoShipment, eTag, "test")
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.Run("404 Not Found Error - shipment can only be created for service member associated with the current session", func() {
		setupTestData()

		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: mtoShipment.MoveTaskOrder.Orders.ServiceMemberID,
		})

		eTag := etag.GenerateEtag(oldMTOShipment.UpdatedAt)
		move := factory.BuildMove(suite.DB(), nil, nil)

		shipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		updatedShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(session, &shipment, eTag, "test")
		suite.Error(err)
		suite.Nil(updatedShipment)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("If-Unmodified-Since is equal to the updated_at date", func() {
		setupTestData()

		eTag := etag.GenerateEtag(oldMTOShipment.UpdatedAt)

		session := auth.Session{}
		updatedMTOShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &mtoShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.Equal(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeUnaccompaniedBaggage)

		suite.Equal(updatedMTOShipment.PickupAddressID, oldMTOShipment.PickupAddressID)

		suite.Equal(updatedMTOShipment.PrimeActualWeight, &primeActualWeight)
		suite.True(actualPickupDate.Equal(*updatedMTOShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*updatedMTOShipment.FirstAvailableDeliveryDate))
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.Run("Updater can handle optional queries set as nil", func() {
		setupTestData()

		oldMTOShipment2 := factory.BuildMTOShipment(suite.DB(), nil, nil)
		mtoShipment2 := models.MTOShipment{
			ID:           oldMTOShipment2.ID,
			ShipmentType: "UNACCOMPANIED_BAGGAGE",
		}

		eTag := etag.GenerateEtag(oldMTOShipment2.UpdatedAt)
		session := auth.Session{}
		updatedMTOShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &mtoShipment2, eTag, "test")

		suite.Require().NoError(err)
		suite.Equal(updatedMTOShipment.ID, oldMTOShipment2.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment2.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeUnaccompaniedBaggage)
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.Run("Successfully remove a secondary pickup address", func() {
		setupTestData()

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
			{
				Model:    secondaryPickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryPickupAddress,
			},
		}, nil)
		suite.FatalNotNil(oldShipment.SecondaryPickupAddress)
		suite.FatalNotNil(oldShipment.SecondaryPickupAddressID)
		suite.FatalNotNil(oldShipment.HasSecondaryPickupAddress)
		suite.True(*oldShipment.HasSecondaryPickupAddress)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		no := false
		updatedShipment := models.MTOShipment{
			ID:                        oldShipment.ID,
			HasSecondaryPickupAddress: &no,
		}

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.FatalNotNil(newShipment.HasSecondaryPickupAddress)
		suite.False(*newShipment.HasSecondaryPickupAddress)
		suite.Nil(newShipment.SecondaryPickupAddress)
	})
	suite.Run("Successfully remove a secondary delivery address", func() {
		setupTestData()

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
			{
				Model:    secondaryDeliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
			},
		}, nil)
		suite.FatalNotNil(oldShipment.SecondaryDeliveryAddress)
		suite.FatalNotNil(oldShipment.SecondaryDeliveryAddressID)
		suite.FatalNotNil(oldShipment.HasSecondaryDeliveryAddress)
		suite.True(*oldShipment.HasSecondaryDeliveryAddress)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		no := false
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			HasSecondaryDeliveryAddress: &no,
		}

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.FatalNotNil(newShipment.HasSecondaryDeliveryAddress)
		suite.False(*newShipment.HasSecondaryDeliveryAddress)
		suite.Nil(newShipment.SecondaryDeliveryAddress)
	})

	suite.Run("Successfully remove a tertiary pickup address", func() {
		setupTestData()

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
			{
				Model:    tertiaryPickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.TertiaryPickupAddress,
			},
		}, nil)
		suite.FatalNotNil(oldShipment.TertiaryPickupAddress)
		suite.FatalNotNil(oldShipment.TertiaryPickupAddressID)
		suite.FatalNotNil(oldShipment.HasTertiaryPickupAddress)
		suite.True(*oldShipment.HasTertiaryPickupAddress)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		no := false
		updatedShipment := models.MTOShipment{
			ID:                       oldShipment.ID,
			HasTertiaryPickupAddress: &no,
		}

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.FatalNotNil(newShipment.HasTertiaryPickupAddress)
		suite.False(*newShipment.HasTertiaryPickupAddress)
		suite.Nil(newShipment.TertiaryPickupAddress)
	})
	suite.Run("Successfully remove a tertiary delivery address", func() {
		setupTestData()

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
			{
				Model:    tertiaryDeliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.TertiaryDeliveryAddress,
			},
		}, nil)
		suite.FatalNotNil(oldShipment.TertiaryDeliveryAddress)
		suite.FatalNotNil(oldShipment.TertiaryDeliveryAddressID)
		suite.FatalNotNil(oldShipment.HasTertiaryDeliveryAddress)
		suite.True(*oldShipment.HasTertiaryDeliveryAddress)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		no := false
		updatedShipment := models.MTOShipment{
			ID:                         oldShipment.ID,
			HasTertiaryDeliveryAddress: &no,
		}

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.FatalNotNil(newShipment.HasTertiaryDeliveryAddress)
		suite.False(*newShipment.HasTertiaryDeliveryAddress)
		suite.Nil(newShipment.TertiaryDeliveryAddress)
	})

	suite.Run("Successful update to all address fields for domestic shipment", func() {
		setupTestData()

		// Ensure we can update every address field on the shipment
		// Create an mtoShipment to update that has every address populated
		oldMTOShipment3 := factory.BuildMTOShipment(suite.DB(), nil, nil)

		eTag := etag.GenerateEtag(oldMTOShipment3.UpdatedAt)

		updatedShipment := &models.MTOShipment{
			ID:                          oldMTOShipment3.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			HasSecondaryPickupAddress:   models.BoolPointer(true),
			SecondaryPickupAddress:      &secondaryPickupAddress,
			SecondaryPickupAddressID:    &secondaryDeliveryAddress.ID,
			HasSecondaryDeliveryAddress: models.BoolPointer(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			SecondaryDeliveryAddressID:  &secondaryDeliveryAddress.ID,
			HasTertiaryPickupAddress:    models.BoolPointer(true),
			TertiaryPickupAddress:       &tertiaryPickupAddress,
			TertiaryPickupAddressID:     &tertiaryPickupAddress.ID,
			HasTertiaryDeliveryAddress:  models.BoolPointer(true),
			TertiaryDeliveryAddress:     &tertiaryDeliveryAddress,
			TertiaryDeliveryAddressID:   &tertiaryDeliveryAddress.ID,
		}
		session := auth.Session{}
		updatedShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.Equal(newDestinationAddress.ID, *updatedShipment.DestinationAddressID)
		suite.Equal(newDestinationAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(newPickupAddress.ID, *updatedShipment.PickupAddressID)
		suite.Equal(newPickupAddress.StreetAddress1, updatedShipment.PickupAddress.StreetAddress1)
		suite.Equal(secondaryPickupAddress.ID, *updatedShipment.SecondaryPickupAddressID)
		suite.Equal(secondaryPickupAddress.StreetAddress1, updatedShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(secondaryDeliveryAddress.ID, *updatedShipment.SecondaryDeliveryAddressID)
		suite.Equal(secondaryDeliveryAddress.StreetAddress1, updatedShipment.SecondaryDeliveryAddress.StreetAddress1)

		suite.Equal(tertiaryPickupAddress.ID, *updatedShipment.TertiaryPickupAddressID)
		suite.Equal(tertiaryPickupAddress.StreetAddress1, updatedShipment.TertiaryPickupAddress.StreetAddress1)
		suite.Equal(tertiaryDeliveryAddress.ID, *updatedShipment.TertiaryDeliveryAddressID)
		suite.Equal(tertiaryDeliveryAddress.StreetAddress1, updatedShipment.TertiaryDeliveryAddress.StreetAddress1)
		suite.Equal(updatedShipment.MarketCode, models.MarketCodeDomestic)
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.Run("Successful update to all address fields resulting in change of market code", func() {
		setupTestData()

		previousShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		newDestinationAddress.State = "AK"
		newPickupAddress.State = "HI"
		// this should be "d" since it is default
		suite.Equal(previousShipment.MarketCode, models.MarketCodeDomestic)

		eTag := etag.GenerateEtag(previousShipment.UpdatedAt)

		updatedShipment := &models.MTOShipment{
			ID:                   previousShipment.ID,
			DestinationAddress:   &newDestinationAddress,
			DestinationAddressID: &newDestinationAddress.ID,
			PickupAddress:        &newPickupAddress,
			PickupAddressID:      &newPickupAddress.ID,
		}
		session := auth.Session{}
		updatedShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), updatedShipment, eTag, "test")

		suite.NoError(err)
		suite.Equal(newDestinationAddress.ID, *updatedShipment.DestinationAddressID)
		suite.True(*updatedShipment.DestinationAddress.IsOconus)
		suite.Equal(newPickupAddress.ID, *updatedShipment.PickupAddressID)
		suite.True(*updatedShipment.PickupAddress.IsOconus)
		suite.Equal(updatedShipment.MarketCode, models.MarketCodeInternational)
	})

	suite.Run("Successful update to a minimal MTO shipment", func() {
		setupTestData()

		// Minimal MTO Shipment has no associated addresses created by default.
		// Part of this test ensures that if an address doesn't exist on a shipment,
		// the updater can successfully create it.
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		scheduledPickupDate := now.Add(time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		primeEstimatedWeightRecordedDate := now.Add(time.Hour * 24 * 3)
		customerRemarks := "I have a grandfather clock"
		counselorRemarks := "Counselor approved"
		actualProGearWeight := unit.Pound(400)
		actualSpouseProGearWeight := unit.Pound(125)
		updatedShipment := models.MTOShipment{
			ID:                               oldShipment.ID,
			DestinationAddress:               &newDestinationAddress,
			DestinationAddressID:             &newDestinationAddress.ID,
			PickupAddress:                    &newPickupAddress,
			PickupAddressID:                  &newPickupAddress.ID,
			SecondaryPickupAddress:           &secondaryPickupAddress,
			HasSecondaryPickupAddress:        handlers.FmtBool(true),
			SecondaryDeliveryAddress:         &secondaryDeliveryAddress,
			HasSecondaryDeliveryAddress:      handlers.FmtBool(true),
			TertiaryPickupAddress:            &tertiaryPickupAddress,
			HasTertiaryPickupAddress:         handlers.FmtBool(true),
			TertiaryDeliveryAddress:          &tertiaryDeliveryAddress,
			HasTertiaryDeliveryAddress:       handlers.FmtBool(true),
			RequestedPickupDate:              &requestedPickupDate,
			ScheduledPickupDate:              &scheduledPickupDate,
			RequestedDeliveryDate:            &requestedDeliveryDate,
			ActualPickupDate:                 &actualPickupDate,
			ActualDeliveryDate:               &actualDeliveryDate,
			ScheduledDeliveryDate:            &scheduledDeliveryDate,
			PrimeActualWeight:                &primeActualWeight,
			PrimeEstimatedWeight:             &primeEstimatedWeight,
			FirstAvailableDeliveryDate:       &firstAvailableDeliveryDate,
			PrimeEstimatedWeightRecordedDate: &primeEstimatedWeightRecordedDate,
			Status:                           models.MTOShipmentStatusSubmitted,
			CustomerRemarks:                  &customerRemarks,
			CounselorRemarks:                 &counselorRemarks,
			ActualProGearWeight:              &actualProGearWeight,
			ActualSpouseProGearWeight:        &actualSpouseProGearWeight,
		}

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.True(requestedPickupDate.Equal(*newShipment.RequestedPickupDate))
		suite.True(scheduledPickupDate.Equal(*newShipment.ScheduledPickupDate))
		suite.True(requestedDeliveryDate.Equal(*newShipment.RequestedDeliveryDate))
		suite.True(actualPickupDate.Equal(*newShipment.ActualPickupDate))
		suite.True(actualDeliveryDate.Equal(*newShipment.ActualDeliveryDate))
		suite.True(scheduledDeliveryDate.Equal(*newShipment.ScheduledDeliveryDate))
		suite.True(firstAvailableDeliveryDate.Equal(*newShipment.FirstAvailableDeliveryDate))
		suite.True(primeEstimatedWeightRecordedDate.Equal(*newShipment.PrimeEstimatedWeightRecordedDate))
		suite.Equal(primeEstimatedWeight, *newShipment.PrimeEstimatedWeight)
		suite.Equal(primeActualWeight, *newShipment.PrimeActualWeight)
		suite.Equal(customerRemarks, *newShipment.CustomerRemarks)
		suite.Equal(counselorRemarks, *newShipment.CounselorRemarks)
		suite.Equal(models.MTOShipmentStatusSubmitted, newShipment.Status)
		suite.Equal(newDestinationAddress.ID, *newShipment.DestinationAddressID)
		suite.Equal(newPickupAddress.ID, *newShipment.PickupAddressID)
		suite.Equal(secondaryPickupAddress.ID, *newShipment.SecondaryPickupAddressID)
		suite.Equal(secondaryDeliveryAddress.ID, *newShipment.SecondaryDeliveryAddressID)
		suite.Equal(tertiaryPickupAddress.ID, *newShipment.TertiaryPickupAddressID)
		suite.Equal(tertiaryDeliveryAddress.ID, *newShipment.TertiaryDeliveryAddressID)
		suite.Equal(actualProGearWeight, *newShipment.ActualProGearWeight)
		suite.Equal(actualSpouseProGearWeight, *newShipment.ActualSpouseProGearWeight)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Updating a shipment does not nullify ApprovedDate", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		customerRemarks := "I have a grandfather clock"
		counselorRemarks := "Counselor approved"
		updatedShipment := models.MTOShipment{
			ID:                       oldShipment.ID,
			DestinationAddress:       &newDestinationAddress,
			DestinationAddressID:     &newDestinationAddress.ID,
			PickupAddress:            &newPickupAddress,
			PickupAddressID:          &newPickupAddress.ID,
			SecondaryPickupAddress:   &secondaryPickupAddress,
			SecondaryDeliveryAddress: &secondaryDeliveryAddress,
			TertiaryPickupAddress:    &tertiaryPickupAddress,
			TertiaryDeliveryAddress:  &tertiaryDeliveryAddress,
			RequestedPickupDate:      &requestedPickupDate,
			RequestedDeliveryDate:    &requestedDeliveryDate,
			CustomerRemarks:          &customerRemarks,
			CounselorRemarks:         &counselorRemarks,
		}
		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotEmpty(newShipment.ApprovedDate)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Can update destination address type on shipment", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		customerRemarks := "I have a grandfather clock"
		counselorRemarks := "Counselor approved"
		destinationType := models.DestinationTypeHomeOfRecord
		updatedShipment := models.MTOShipment{
			ID:                       oldShipment.ID,
			DestinationAddress:       &newDestinationAddress,
			DestinationAddressID:     &newDestinationAddress.ID,
			DestinationType:          &destinationType,
			PickupAddress:            &newPickupAddress,
			PickupAddressID:          &newPickupAddress.ID,
			SecondaryPickupAddress:   &secondaryPickupAddress,
			SecondaryDeliveryAddress: &secondaryDeliveryAddress,
			TertiaryPickupAddress:    &tertiaryPickupAddress,
			TertiaryDeliveryAddress:  &tertiaryDeliveryAddress,
			RequestedPickupDate:      &requestedPickupDate,
			RequestedDeliveryDate:    &requestedDeliveryDate,
			CustomerRemarks:          &customerRemarks,
			CounselorRemarks:         &counselorRemarks,
		}
		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		session.Roles = append(session.Roles, too.User.Roles...)
		newShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.Equal(destinationType, *newShipment.DestinationType)
	})

	suite.Run("Successfully update MTO Agents", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		mtoAgent1 := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOAgent{
					FirstName:    models.StringPointer("Test"),
					LastName:     models.StringPointer("Agent"),
					Email:        models.StringPointer("test@test.email.com"),
					MTOAgentType: models.MTOAgentReleasing,
				},
			},
		}, nil)
		mtoAgent2 := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOAgent{
					FirstName:    models.StringPointer("Test2"),
					LastName:     models.StringPointer("Agent2"),
					Email:        models.StringPointer("test2@test.email.com"),
					MTOAgentType: models.MTOAgentReceiving,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedAgents := make(models.MTOAgents, 2)
		updatedAgents[0] = mtoAgent1
		updatedAgents[1] = mtoAgent2
		newFirstName := "hey this is new"
		newLastName := "new thing"
		phone := "555-666-7777"
		email := "updatedemail@test.email.com"
		updatedAgents[0].FirstName = &newFirstName
		updatedAgents[0].Phone = &phone
		updatedAgents[1].LastName = &newLastName
		updatedAgents[1].Email = &email

		updatedShipment := models.MTOShipment{
			ID:        shipment.ID,
			MTOAgents: updatedAgents,
		}

		session := auth.Session{}
		updatedMTOShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(phone, *updatedMTOShipment.MTOAgents[0].Phone)
		suite.Equal(newFirstName, *updatedMTOShipment.MTOAgents[0].FirstName)
		suite.Equal(email, *updatedMTOShipment.MTOAgents[1].Email)
		suite.Equal(newLastName, *updatedMTOShipment.MTOAgents[1].LastName)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Successfully add new MTO Agent and edit another", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		existingAgent := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOAgent{
					FirstName:    models.StringPointer("Test"),
					LastName:     models.StringPointer("Agent"),
					Email:        models.StringPointer("test@test.email.com"),
					MTOAgentType: models.MTOAgentReleasing,
				},
			},
		}, nil)
		mtoAgentToCreate := models.MTOAgent{
			MTOShipment:   shipment,
			MTOShipmentID: shipment.ID,
			FirstName:     models.StringPointer("Ima"),
			LastName:      models.StringPointer("Newagent"),
			Email:         models.StringPointer("test2@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		}
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedAgents := make(models.MTOAgents, 2)
		phone := "555-555-5555"
		existingAgent.Phone = &phone
		updatedAgents[1] = existingAgent
		updatedAgents[0] = mtoAgentToCreate

		updatedShipment := models.MTOShipment{
			ID:        shipment.ID,
			MTOAgents: updatedAgents,
		}

		session := auth.Session{}
		updatedMTOShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		// the returned updatedMTOShipment does not guarantee the same
		// order of MTOAgents
		suite.Equal(len(updatedAgents), len(updatedMTOShipment.MTOAgents))
		for i := range updatedMTOShipment.MTOAgents {
			agent := updatedMTOShipment.MTOAgents[i]
			if agent.ID == existingAgent.ID {
				suite.Equal(phone, *agent.Phone)
			} else {
				// this must be the newly created agent
				suite.Equal(*mtoAgentToCreate.FirstName, *agent.FirstName)
				suite.Equal(*mtoAgentToCreate.LastName, *agent.LastName)
				suite.Equal(*mtoAgentToCreate.Email, *agent.Email)
			}
		}

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Successfully remove MTO Agent", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		existingAgent := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOAgent{
					FirstName:    models.StringPointer("Test"),
					LastName:     models.StringPointer("Agent"),
					Email:        models.StringPointer("test@test.email.com"),
					MTOAgentType: models.MTOAgentReleasing,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedAgents := make(models.MTOAgents, 1)
		blankFirstName := ""
		blankLastName := ""
		blankPhone := ""
		blankEmail := ""
		existingAgent.FirstName = &blankFirstName
		existingAgent.LastName = &blankLastName
		existingAgent.Email = &blankEmail
		existingAgent.Phone = &blankPhone
		updatedAgents[0] = existingAgent

		updatedShipment := models.MTOShipment{
			ID:        shipment.ID,
			MTOAgents: updatedAgents,
		}

		session := auth.Session{}
		updatedMTOShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		// Verify that there are no returned MTO Agents
		suite.Equal(0, len(updatedMTOShipment.MTOAgents))

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Successfully add storage facility to shipment", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)

		updatedShipment := models.MTOShipment{
			ID:              shipment.ID,
			StorageFacility: &storageFacility,
		}
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		session.Roles = append(session.Roles, too.User.Roles...)
		updatedMTOShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.StorageFacility)
	})

	suite.Run("Successfully edit storage facility on shipment", func() {
		setupTestData()

		// Create initial shipment data
		storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.StorageFacility{
					Email: models.StringPointer("old@email.com"),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "1234 Over Here Street",
					City:           "Houston",
					State:          "TX",
					PostalCode:     "77083",
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)

		// Make updates to previously persisted data (don't need to create these in the DB first)
		newStorageFacilityAddress := models.Address{
			StreetAddress1: "987 Over There Avenue",
			City:           "Houston",
			State:          "TX",
			PostalCode:     "77083",
		}

		newEmail := "new@email.com"
		newStorageFacility := models.StorageFacility{
			Address: newStorageFacilityAddress,
			Email:   &newEmail,
		}

		newShipment := models.MTOShipment{
			ID:              shipment.ID,
			StorageFacility: &newStorageFacility,
		}

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		session.Roles = append(session.Roles, too.User.Roles...)
		updatedShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &newShipment, eTag, "test")
		suite.Require().NoError(err)
		suite.NotEqual(uuid.Nil, updatedShipment.ID)
		suite.Equal(&newEmail, updatedShipment.StorageFacility.Email)
		suite.Equal(newStorageFacilityAddress.StreetAddress1, updatedShipment.StorageFacility.Address.StreetAddress1)
	})

	suite.Run("Successfully update NTS previously recorded weight to shipment", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		ntsRecorededWeight := unit.Pound(980)
		updatedShipment := models.MTOShipment{
			ShipmentType:      models.MTOShipmentTypeHHGOutOfNTSDom,
			ID:                shipment.ID,
			NTSRecordedWeight: &ntsRecorededWeight,
		}
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		session.Roles = append(session.Roles, too.User.Roles...)
		updatedMTOShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(ntsRecorededWeight, *updatedMTOShipment.NTSRecordedWeight)

	})

	suite.Run("Unable to update NTS previously recorded weight due to shipment type", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		ntsRecorededWeight := unit.Pound(980)
		updatedShipment := models.MTOShipment{
			ID:                shipment.ID,
			NTSRecordedWeight: &ntsRecorededWeight,
		}
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		session.Roles = append(session.Roles, too.User.Roles...)
		updatedMTOShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().Error(err)
		suite.Nil(updatedMTOShipment)
		suite.Equal("Could not complete query related to object of type: mtoShipment.", err.Error())

		suite.IsType(apperror.QueryError{}, err)
		queryErr := err.(apperror.QueryError)
		wrappedErr := queryErr.Unwrap()
		suite.IsType(apperror.InvalidInputError{}, wrappedErr)
		suite.Equal("field NTSRecordedWeight cannot be set for shipment type HHG", wrappedErr.Error())
	})

	suite.Run("Successfully divert a shipment and transition statuses", func() {
		setupTestData()

		// A diverted shipment should transition to the SUBMITTED status.
		// If the move it is connected to is APPROVED, that move should transition to APPROVALS REQUESTED
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:    models.MTOShipmentStatusApproved,
					Diversion: false,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		shipmentInput := models.MTOShipment{
			ID:        shipment.ID,
			Diversion: true,
		}
		session := auth.Session{}
		updatedShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &shipmentInput, eTag, "test")

		suite.Require().NotNil(updatedShipment)
		suite.NoError(err)
		suite.Equal(shipment.ID, updatedShipment.ID)
		suite.Equal(move.ID, updatedShipment.MoveTaskOrderID)
		suite.Equal(true, updatedShipment.Diversion)
		suite.Equal(models.MTOShipmentStatusSubmitted, updatedShipment.Status)

		var updatedMove models.Move
		err = suite.DB().Find(&updatedMove, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	// Test UpdateMTOShipmentPrime
	// TODO: Add more tests, such as making sure this function fails if the
	// move is not available to the prime.
	suite.Run("Updating a shipment does not nullify ApprovedDate", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		scheduledPickupDate := now.Add(time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			HasSecondaryPickupAddress:   handlers.FmtBool(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			TertiaryPickupAddress:       &tertiaryPickupAddress,
			HasTertiaryPickupAddress:    handlers.FmtBool(true),
			TertiaryDeliveryAddress:     &tertiaryDeliveryAddress,
			HasTertiaryDeliveryAddress:  handlers.FmtBool(true),
			RequestedPickupDate:         &requestedPickupDate,
			ScheduledPickupDate:         &scheduledPickupDate,
			RequestedDeliveryDate:       &requestedDeliveryDate,
			ActualPickupDate:            &actualPickupDate,
			PrimeActualWeight:           &primeActualWeight,
			PrimeEstimatedWeight:        &primeEstimatedWeight,
			FirstAvailableDeliveryDate:  &firstAvailableDeliveryDate,
		}

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.False(verrs.HasAny())
		suite.FatalNoError(err)

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotEmpty(newShipment.ApprovedDate)
		suite.True(requestedPickupDate.Equal(*newShipment.RequestedPickupDate))
		suite.True(scheduledPickupDate.Equal(*newShipment.ScheduledPickupDate))
		suite.True(requestedDeliveryDate.Equal(*newShipment.RequestedDeliveryDate))
		suite.True(actualPickupDate.Equal(*newShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*newShipment.FirstAvailableDeliveryDate))
		suite.Equal(primeEstimatedWeight, *newShipment.PrimeEstimatedWeight)
		suite.Equal(primeActualWeight, *newShipment.PrimeActualWeight)
		suite.Equal(newDestinationAddress.ID, *newShipment.DestinationAddressID)
		suite.Equal(newPickupAddress.ID, *newShipment.PickupAddressID)
		suite.Equal(secondaryPickupAddress.ID, *newShipment.SecondaryPickupAddressID)
		suite.Equal(secondaryDeliveryAddress.ID, *newShipment.SecondaryDeliveryAddressID)
		suite.Equal(tertiaryPickupAddress.ID, *newShipment.TertiaryPickupAddressID)
		suite.Equal(tertiaryDeliveryAddress.ID, *newShipment.TertiaryDeliveryAddressID)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Prime not able to update an existing prime estimated weight", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: &primeEstimatedWeight,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		scheduledPickupDate := now.Add(time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			HasSecondaryPickupAddress:   handlers.FmtBool(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			RequestedPickupDate:         &requestedPickupDate,
			ScheduledPickupDate:         &scheduledPickupDate,
			RequestedDeliveryDate:       &requestedDeliveryDate,
			ActualPickupDate:            &actualPickupDate,
			PrimeActualWeight:           &primeActualWeight,
			PrimeEstimatedWeight:        &primeEstimatedWeight,
			FirstAvailableDeliveryDate:  &firstAvailableDeliveryDate,
		}

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.False(verrs.HasAny())
		suite.FatalNoError(err)

		session := auth.Session{}
		_, err = mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Error(err)
		suite.Contains(err.Error(), "cannot be updated after initial estimation")
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Updating a shipment with a Reweigh returns the Reweigh", func() {
		setupTestData()

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		reweigh := testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, oldShipment, unit.Pound(3000))

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		updatedShipment := models.MTOShipment{
			ID:                oldShipment.ID,
			PrimeActualWeight: &primeActualWeight,
		}

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotEmpty(newShipment.Reweigh)
		suite.Equal(newShipment.Reweigh.ID, reweigh.ID)
	})

	suite.Run("Prime cannot update estimated weights outside of required timeframe", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		scheduledPickupDate := now.Add(-time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			HasSecondaryPickupAddress:   handlers.FmtBool(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			TertiaryPickupAddress:       &tertiaryPickupAddress,
			HasTertiaryPickupAddress:    handlers.FmtBool(true),
			TertiaryDeliveryAddress:     &tertiaryDeliveryAddress,
			HasTertiaryDeliveryAddress:  handlers.FmtBool(true),
			RequestedPickupDate:         &requestedPickupDate,
			ScheduledPickupDate:         &scheduledPickupDate,
			RequestedDeliveryDate:       &requestedDeliveryDate,
			ActualPickupDate:            &actualPickupDate,
			PrimeActualWeight:           &primeActualWeight,
			PrimeEstimatedWeight:        &primeEstimatedWeight,
			FirstAvailableDeliveryDate:  &firstAvailableDeliveryDate,
		}

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.False(verrs.HasAny())
		suite.FatalNoError(err)

		session := auth.Session{}
		_, err = mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Error(err)
		suite.Contains(err.Error(), "the time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight")
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Prime cannot add MTO agents", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		scheduledPickupDate := now.Add(time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		firstName := "John"
		lastName := "Ash"
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			HasSecondaryPickupAddress:   handlers.FmtBool(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			TertiaryPickupAddress:       &tertiaryPickupAddress,
			HasTertiaryPickupAddress:    handlers.FmtBool(true),
			TertiaryDeliveryAddress:     &tertiaryDeliveryAddress,
			HasTertiaryDeliveryAddress:  handlers.FmtBool(true),
			RequestedPickupDate:         &requestedPickupDate,
			ScheduledPickupDate:         &scheduledPickupDate,
			RequestedDeliveryDate:       &requestedDeliveryDate,
			ActualPickupDate:            &actualPickupDate,
			PrimeActualWeight:           &primeActualWeight,
			PrimeEstimatedWeight:        &primeEstimatedWeight,
			FirstAvailableDeliveryDate:  &firstAvailableDeliveryDate,
			MTOAgents: models.MTOAgents{
				models.MTOAgent{
					FirstName: &firstName,
					LastName:  &lastName,
				},
			},
		}

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.False(verrs.HasAny())
		suite.FatalNoError(err)

		session := auth.Session{}
		_, err = mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Error(err)
		suite.Contains(err.Error(), "cannot add or update MTO agents to a shipment")
	})

	suite.Run("Prime cannot update existing pickup or destination address", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		scheduledPickupDate := now.Add(time.Hour * 24 * 7)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			HasSecondaryPickupAddress:   handlers.FmtBool(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			TertiaryPickupAddress:       &secondaryPickupAddress,
			HasTertiaryPickupAddress:    handlers.FmtBool(true),
			TertiaryDeliveryAddress:     &secondaryDeliveryAddress,
			HasTertiaryDeliveryAddress:  handlers.FmtBool(true),
			RequestedPickupDate:         &requestedPickupDate,
			ScheduledPickupDate:         &scheduledPickupDate,
			RequestedDeliveryDate:       &requestedDeliveryDate,
			ActualPickupDate:            &actualPickupDate,
			PrimeActualWeight:           &primeActualWeight,
			PrimeEstimatedWeight:        &primeEstimatedWeight,
			FirstAvailableDeliveryDate:  &firstAvailableDeliveryDate,
		}

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.False(verrs.HasAny())
		suite.FatalNoError(err)

		session := auth.Session{}
		_, err = mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Error(err)
		suite.Contains(err.Error(), "the pickup address already exists and cannot be updated with this endpoint")
		suite.Contains(err.Error(), "the destination address already exists and cannot be updated with this endpoint")
	})

	suite.Run("Prime cannot update shipment if parameters are outside of transit data", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		scheduledPickupDate := now.Add(time.Hour * 24 * 7)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			HasSecondaryPickupAddress:   handlers.FmtBool(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			TertiaryPickupAddress:       &tertiaryPickupAddress,
			HasTertiaryPickupAddress:    handlers.FmtBool(true),
			TertiaryDeliveryAddress:     &tertiaryDeliveryAddress,
			HasTertiaryDeliveryAddress:  handlers.FmtBool(true),
			RequestedPickupDate:         &requestedPickupDate,
			ScheduledPickupDate:         &scheduledPickupDate,
			RequestedDeliveryDate:       &requestedDeliveryDate,
			ActualPickupDate:            &actualPickupDate,
			PrimeActualWeight:           &primeActualWeight,
			PrimeEstimatedWeight:        &primeEstimatedWeight,
			FirstAvailableDeliveryDate:  &firstAvailableDeliveryDate,
		}

		session := auth.Session{}
		_, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Error(err)
		suite.Contains(err.Error(), "failed to find transit time for shipment of 9000 lbs weight and 1000 mile distance")
	})

	suite.Run("Prime can add an estimated weight up to the same date as the scheduled pickup", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		scheduledPickupDate := now.Add(time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			HasSecondaryPickupAddress:   handlers.FmtBool(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			ScheduledPickupDate:         &scheduledPickupDate,
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			TertiaryPickupAddress:       &tertiaryPickupAddress,
			HasTertiaryPickupAddress:    handlers.FmtBool(true),
			TertiaryDeliveryAddress:     &tertiaryDeliveryAddress,
			HasTertiaryDeliveryAddress:  handlers.FmtBool(true),
			RequestedPickupDate:         &requestedPickupDate,
			RequestedDeliveryDate:       &requestedDeliveryDate,
			ActualPickupDate:            &actualPickupDate,
			PrimeActualWeight:           &primeActualWeight,
			PrimeEstimatedWeight:        &primeEstimatedWeight,
			FirstAvailableDeliveryDate:  &firstAvailableDeliveryDate,
		}

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.False(verrs.HasAny())
		suite.FatalNoError(err)

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotEmpty(newShipment.ApprovedDate)
		suite.True(requestedPickupDate.Equal(*newShipment.RequestedPickupDate))
		suite.True(scheduledPickupDate.Equal(*newShipment.ScheduledPickupDate))
		suite.True(requestedDeliveryDate.Equal(*newShipment.RequestedDeliveryDate))
		suite.True(actualPickupDate.Equal(*newShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*newShipment.FirstAvailableDeliveryDate))
		suite.Equal(primeEstimatedWeight, *newShipment.PrimeEstimatedWeight)
		suite.Equal(primeActualWeight, *newShipment.PrimeActualWeight)
		suite.Equal(newDestinationAddress.ID, *newShipment.DestinationAddressID)
		suite.Equal(newPickupAddress.ID, *newShipment.PickupAddressID)
		suite.Equal(secondaryPickupAddress.ID, *newShipment.SecondaryPickupAddressID)
		suite.Equal(secondaryDeliveryAddress.ID, *newShipment.SecondaryDeliveryAddressID)
		suite.Equal(tertiaryPickupAddress.ID, *newShipment.TertiaryPickupAddressID)
		suite.Equal(tertiaryDeliveryAddress.ID, *newShipment.TertiaryDeliveryAddressID)
	})

	suite.Run("Prime can update the weight estimate if scheduled pickup date in nil", func() {
		setupTestData()

		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ScheduledPickupDate: nil,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := now.Add(time.Hour * 24 * 3)
		requestedDeliveryDate := now.Add(time.Hour * 24 * 4)
		updatedShipment := models.MTOShipment{
			ID:                          oldShipment.ID,
			DestinationAddress:          &newDestinationAddress,
			DestinationAddressID:        &newDestinationAddress.ID,
			PickupAddress:               &newPickupAddress,
			PickupAddressID:             &newPickupAddress.ID,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			HasSecondaryPickupAddress:   handlers.FmtBool(true),
			SecondaryDeliveryAddress:    &secondaryDeliveryAddress,
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			TertiaryPickupAddress:       &tertiaryPickupAddress,
			HasTertiaryPickupAddress:    handlers.FmtBool(true),
			TertiaryDeliveryAddress:     &tertiaryDeliveryAddress,
			HasTertiaryDeliveryAddress:  handlers.FmtBool(true),
			RequestedPickupDate:         &requestedPickupDate,
			RequestedDeliveryDate:       &requestedDeliveryDate,
			ActualPickupDate:            &actualPickupDate,
			PrimeActualWeight:           &primeActualWeight,
			PrimeEstimatedWeight:        &primeEstimatedWeight,
			FirstAvailableDeliveryDate:  &firstAvailableDeliveryDate,
		}
		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.False(verrs.HasAny())
		suite.FatalNoError(err)

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")
		suite.Require().NoError(err)
		suite.NotEmpty(newShipment.ApprovedDate)
		suite.True(requestedPickupDate.Equal(*newShipment.RequestedPickupDate))
		suite.True(requestedDeliveryDate.Equal(*newShipment.RequestedDeliveryDate))
		suite.True(actualPickupDate.Equal(*newShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*newShipment.FirstAvailableDeliveryDate))
		suite.Equal(primeEstimatedWeight, *newShipment.PrimeEstimatedWeight)
		suite.Equal(primeActualWeight, *newShipment.PrimeActualWeight)
		suite.Equal(newDestinationAddress.ID, *newShipment.DestinationAddressID)
		suite.Equal(newPickupAddress.ID, *newShipment.PickupAddressID)
		suite.Equal(secondaryPickupAddress.ID, *newShipment.SecondaryPickupAddressID)
		suite.Equal(secondaryDeliveryAddress.ID, *newShipment.SecondaryDeliveryAddressID)
		suite.Equal(tertiaryPickupAddress.ID, *newShipment.TertiaryPickupAddressID)
		suite.Equal(tertiaryDeliveryAddress.ID, *newShipment.TertiaryDeliveryAddressID)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentStatus() {
	estimatedWeight := unit.Pound(2000)
	status := models.MTOShipmentStatusApproved
	// need the re service codes to update status
	expectedReServiceCodes := []models.ReServiceCode{
		models.ReServiceCodeDLH,
		models.ReServiceCodeFSC,
		models.ReServiceCodeDOP,
		models.ReServiceCodeDDP,
		models.ReServiceCodeDPK,
		models.ReServiceCodeDUPK,
	}

	var shipmentForAutoApprove models.MTOShipment
	var draftShipment models.MTOShipment
	var shipment2 models.MTOShipment
	var shipment3 models.MTOShipment
	var shipment4 models.MTOShipment
	var approvedShipment models.MTOShipment
	var rejectedShipment models.MTOShipment
	var eTag string
	var mto models.Move

	setupTestData := func() {
		for i := range expectedReServiceCodes {
			factory.FetchReServiceByCode(suite.DB(), expectedReServiceCodes[i])
		}

		mto = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		draftShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusDraft,
				},
			},
		}, nil)
		shipment2 = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		shipment3 = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		shipment4 = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		shipmentForAutoApprove = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		approvedShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)
		rejectionReason := "exotic animals are banned"
		rejectedShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:          models.MTOShipmentStatusRejected,
					RejectionReason: &rejectionReason,
				},
			},
		}, nil)
		shipment.Status = models.MTOShipmentStatusSubmitted
		eTag = etag.GenerateEtag(shipment.UpdatedAt)
	}

	builder := query.NewQueryBuilder()
	moveRouter := moveservices.NewMoveRouter()
	planner := &mocks.Planner{}
	var TransitDistancePickupArg string
	var TransitDistanceDestinationArg string
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(500, nil).Run(func(args mock.Arguments) {
		TransitDistancePickupArg = args.Get(1).(string)
		TransitDistanceDestinationArg = args.Get(2).(string)
	})
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	updater := NewMTOShipmentStatusUpdater(builder, siCreator, planner)

	suite.Run("If the mtoShipment is approved successfully it should create approved mtoServiceItems", func() {
		setupTestData()

		appCtx := suite.AppContextForTest()
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		serviceItems := models.MTOServiceItems{}

		preApprovalTime := time.Now()
		_, err := updater.UpdateMTOShipmentStatus(appCtx, shipmentForAutoApprove.ID, status, nil, nil, shipmentForAutoApproveEtag)
		suite.NoError(err)

		err = appCtx.DB().Find(&fetchedShipment, shipmentForAutoApprove.ID)
		suite.NoError(err)

		// Let's make sure the status is approved
		suite.Equal(models.MTOShipmentStatusApproved, fetchedShipment.Status)

		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipmentForAutoApprove.ID).All(&serviceItems)
		suite.NoError(err)

		suite.Equal(6, len(serviceItems))

		// All ApprovedAt times for service items should be the same, so just get the first one
		// Test that service item was approved within a few seconds of the current time
		suite.Assertions.WithinDuration(preApprovalTime, *serviceItems[0].ApprovedAt, 2*time.Second)

		// If we've gotten the shipment updated and fetched it without error then we can inspect the
		// service items created as a side effect to see if they are
		// approved.
		missingReServiceCodes := make([]models.ReServiceCode, len(expectedReServiceCodes))
		copy(missingReServiceCodes, expectedReServiceCodes)
		for _, serviceItem := range serviceItems {
			suite.Equal(models.MTOServiceItemStatusApproved, serviceItem.Status)

			// Want to make sure each of the expected service codes is included at some point.
			codeFound := false
			for i, reServiceCodeToCheck := range missingReServiceCodes {
				if reServiceCodeToCheck == serviceItem.ReService.Code {
					missingReServiceCodes[i] = missingReServiceCodes[len(missingReServiceCodes)-1]
					missingReServiceCodes = missingReServiceCodes[:len(missingReServiceCodes)-1]
					codeFound = true
					break
				}
			}

			if !codeFound {
				suite.Fail("Unexpected service code", "unexpected ReService code: %s", string(serviceItem.ReService.Code))
			}
		}

		suite.Empty(missingReServiceCodes)
	})

	suite.Run("If we act on a shipment with a weight that has a 0 upper weight it should still work", func() {
		setupTestData()

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.Assert().False(verrs.HasAny())
		suite.NoError(err)

		// Let's also create a transit time object with a zero upper bound for weight (this can happen in the table).
		ghcDomesticTransitTime0LbsUpper := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     10001,
			WeightLbsUpper:     0,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime0LbsUpper)
		suite.Assert().False(verrs.HasAny())
		suite.NoError(err)

		// This is testing that the Required Delivery Date is calculated correctly.
		// In order for the Required Delivery Date to be calculated, the following conditions must be true:
		// 1. The shipment is moving to the APPROVED status
		// 2. The shipment must already have the following fields present:
		// ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, DestinationAddress
		// 3. The shipment must not already have a Required Delivery Date
		// Note that MakeMTOShipment will automatically add a Required Delivery Date if the ScheduledPickupDate
		// is present, therefore we need to use MakeMTOShipmentMinimal and add the Pickup and Destination addresses
		estimatedWeight := unit.Pound(11000)
		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		shipmentHeavy := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)
		shipmentHeavyEtag := etag.GenerateEtag(shipmentHeavy.UpdatedAt)

		_, err = updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipmentHeavy.ID, status, nil, nil, shipmentHeavyEtag)
		suite.NoError(err)
		serviceItems := models.MTOServiceItems{}
		_ = suite.DB().All(&serviceItems)
		fetchedShipment := models.MTOShipment{}
		err = suite.DB().Find(&fetchedShipment, shipmentHeavy.ID)
		suite.NoError(err)
		// We also should have a required delivery date
		suite.NotNil(fetchedShipment.RequiredDeliveryDate)
	})

	suite.Run("Test that correct addresses are being used to calculate required delivery date", func() {
		setupTestData()
		appCtx := suite.AppContextForTest()

		ghcDomesticTransitTime0LbsUpper := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     10001,
			WeightLbsUpper:     0,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime0LbsUpper)
		suite.Assert().False(verrs.HasAny())
		suite.NoError(err)

		factory.FetchReServiceByCode(appCtx.DB(), models.ReServiceCodeDNPK)

		// This is testing that the Required Delivery Date is calculated correctly.
		// In order for the Required Delivery Date to be calculated, the following conditions must be true:
		// 1. The shipment is moving to the APPROVED status
		// 2. The shipment must already have the following fields present:
		// MTOShipmentTypeHHG: ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, DestinationAddress
		// MTOShipmentTypeHHGIntoNTSDom: ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, StorageFacility
		// MTOShipmentTypeHHGOutOfNTSDom: ScheduledPickupDate, NTSRecordedWeight, StorageFacility, DestinationAddress
		// 3. The shipment must not already have a Required Delivery Date
		// Note that MakeMTOShipment will automatically add a Required Delivery Date if the ScheduledPickupDate
		// is present, therefore we need to use MakeMTOShipmentMinimal and add the Pickup and Destination addresses
		estimatedWeight := unit.Pound(11000)

		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)

		hhgShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		ntsShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
		}, nil)

		ntsrShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:        models.MTOShipmentTypeHHGOutOfNTSDom,
					ScheduledPickupDate: &testdatagen.DateInsidePeakRateCycle,
					NTSRecordedWeight:   &estimatedWeight,
					Status:              models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		testCases := []struct {
			shipment            models.MTOShipment
			pickupLocation      *models.Address
			destinationLocation *models.Address
		}{
			{hhgShipment, hhgShipment.PickupAddress, hhgShipment.DestinationAddress},
			{ntsShipment, ntsShipment.PickupAddress, &ntsShipment.StorageFacility.Address},
			{ntsrShipment, &ntsrShipment.StorageFacility.Address, ntsrShipment.DestinationAddress},
		}

		for _, testCase := range testCases {
			shipmentEtag := etag.GenerateEtag(testCase.shipment.UpdatedAt)
			_, err = updater.UpdateMTOShipmentStatus(appCtx, testCase.shipment.ID, status, nil, nil, shipmentEtag)
			suite.NoError(err)

			fetchedShipment := models.MTOShipment{}
			err = suite.DB().Find(&fetchedShipment, testCase.shipment.ID)
			suite.NoError(err)
			// We also should have a required delivery date
			suite.NotNil(fetchedShipment.RequiredDeliveryDate)
			// Check that TransitDistance was called with the correct addresses
			suite.Equal(testCase.pickupLocation.PostalCode, TransitDistancePickupArg)
			suite.Equal(testCase.destinationLocation.PostalCode, TransitDistanceDestinationArg)
		}
	})

	suite.Run("Cannot set SUBMITTED status on shipment via UpdateMTOShipmentStatus", func() {
		setupTestData()

		// The only time a shipment gets set to the SUBMITTED status is when it is created, whether by the customer
		// or the Prime. This happens in the internal and prime API in the CreateMTOShipmentHandler. In that case,
		// the handlers will call ShipmentRouter.Submit().
		eTag = etag.GenerateEtag(draftShipment.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), draftShipment.ID, "SUBMITTED", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)

		err = suite.DB().Find(&draftShipment, draftShipment.ID)

		suite.NoError(err)
		suite.EqualValues(models.MTOShipmentStatusDraft, draftShipment.Status)
	})

	suite.Run("Rejecting a shipment in SUBMITTED status with a rejection reason should return no error", func() {
		setupTestData()

		eTag = etag.GenerateEtag(shipment2.UpdatedAt)
		rejectionReason := "Rejection reason"
		returnedShipment, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment2.ID, "REJECTED", &rejectionReason, nil, eTag)

		suite.NoError(err)
		suite.NotNil(returnedShipment)

		err = suite.DB().Find(&shipment2, shipment2.ID)

		suite.NoError(err)
		suite.EqualValues(models.MTOShipmentStatusRejected, shipment2.Status)
		suite.Equal(&rejectionReason, shipment2.RejectionReason)
	})

	suite.Run("Rejecting a shipment with no rejection reason returns an InvalidInputError", func() {
		setupTestData()

		eTag = etag.GenerateEtag(shipment3.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment3.ID, "REJECTED", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("Rejecting a shipment in APPROVED status returns a ConflictStatusError", func() {
		setupTestData()

		eTag = etag.GenerateEtag(approvedShipment.UpdatedAt)
		rejectionReason := "Rejection reason"
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), approvedShipment.ID, "REJECTED", &rejectionReason, nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Approving a shipment in REJECTED status returns a ConflictStatusError", func() {
		setupTestData()

		eTag = etag.GenerateEtag(rejectedShipment.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), rejectedShipment.ID, "APPROVED", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		setupTestData()

		staleETag := etag.GenerateEtag(time.Now())

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment4.ID, "APPROVED", nil, nil, staleETag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Passing in an invalid status returns a ConflictStatus error", func() {
		setupTestData()

		eTag = etag.GenerateEtag(shipment4.UpdatedAt)

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment4.ID, "invalid", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a bad shipment id returns a Not Found error", func() {
		setupTestData()

		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), badShipmentID, "APPROVED", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Changing to APPROVED status records approved_date", func() {
		setupTestData()

		shipment5 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		eTag = etag.GenerateEtag(shipment5.UpdatedAt)

		suite.Nil(shipment5.ApprovedDate)

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment5.ID, models.MTOShipmentStatusApproved, nil, nil, eTag)

		suite.NoError(err)
		suite.NoError(suite.DB().Find(&shipment5, shipment5.ID))
		suite.Equal(models.MTOShipmentStatusApproved, shipment5.Status)
		suite.NotNil(shipment5.ApprovedDate)
	})

	suite.Run("Changing to a non-APPROVED status does not record approved_date", func() {
		setupTestData()

		shipment6 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		eTag = etag.GenerateEtag(shipment6.UpdatedAt)
		rejectionReason := "reason"

		suite.Nil(shipment6.ApprovedDate)

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment6.ID, models.MTOShipmentStatusRejected, &rejectionReason, nil, eTag)

		suite.NoError(err)
		suite.NoError(suite.DB().Find(&shipment6, shipment6.ID))
		suite.Equal(models.MTOShipmentStatusRejected, shipment6.Status)
		suite.Nil(shipment6.ApprovedDate)
	})

	suite.Run("When move is not yet approved, cannot approve shipment", func() {
		setupTestData()

		submittedMTO := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		mtoShipment := submittedMTO.MTOShipments[0]
		eTag = etag.GenerateEtag(mtoShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), mtoShipment.ID, models.MTOShipmentStatusApproved, nil, nil, eTag)
		suite.NoError(suite.DB().Find(&mtoShipment, mtoShipment.ID))

		suite.Nil(updatedShipment)
		suite.Equal(models.MTOShipmentStatusSubmitted, mtoShipment.Status)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(
			err.Error(),
			fmt.Sprintf(
				"Cannot approve a shipment if the move status isn't %s or %s, or if it isn't a PPM shipment with a move status of %s. The current status for the move with ID %s is %s",
				models.MoveStatusAPPROVED,
				models.MoveStatusAPPROVALSREQUESTED,
				models.MoveStatusNeedsServiceCounseling,
				submittedMTO.ID,
				submittedMTO.Status,
			),
		)
	})

	suite.Run("An approved shipment can change to CANCELLATION_REQUESTED", func() {
		setupTestData()

		approvedShipment2 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)
		eTag = etag.GenerateEtag(approvedShipment2.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), approvedShipment2.ID, models.MTOShipmentStatusCancellationRequested, nil, nil, eTag)
		suite.NoError(suite.DB().Find(&approvedShipment2, approvedShipment2.ID))

		suite.NoError(err)
		suite.NotNil(updatedShipment)
		suite.Equal(models.MTOShipmentStatusCancellationRequested, updatedShipment.Status)
		suite.Equal(models.MTOShipmentStatusCancellationRequested, approvedShipment2.Status)
	})

	suite.Run("A CANCELLATION_REQUESTED shipment can change to CANCELED", func() {
		setupTestData()

		cancellationRequestedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusCancellationRequested,
				},
			},
		}, nil)
		eTag = etag.GenerateEtag(cancellationRequestedShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), cancellationRequestedShipment.ID, models.MTOShipmentStatusCanceled, nil, nil, eTag)
		suite.NoError(suite.DB().Find(&cancellationRequestedShipment, cancellationRequestedShipment.ID))

		suite.NoError(err)
		suite.NotNil(updatedShipment)
		suite.Equal(models.MTOShipmentStatusCanceled, updatedShipment.Status)
		suite.Equal(models.MTOShipmentStatusCanceled, cancellationRequestedShipment.Status)
	})

	suite.Run("An APPROVED shipment CANNOT change to CANCELED - ERROR", func() {
		setupTestData()

		eTag = etag.GenerateEtag(approvedShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), approvedShipment.ID, models.MTOShipmentStatusCanceled, nil, nil, eTag)
		suite.NoError(suite.DB().Find(&approvedShipment, approvedShipment.ID))

		suite.Error(err)
		suite.Nil(updatedShipment)
		suite.IsType(ConflictStatusError{}, err)
		suite.Equal(models.MTOShipmentStatusApproved, approvedShipment.Status)
	})

	suite.Run("An APPROVED shipment CAN change to Diversion Requested", func() {
		setupTestData()

		shipmentToDivert := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)
		eTag = etag.GenerateEtag(shipmentToDivert.UpdatedAt)

		diversionReason := "Test reason"
		_, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), shipmentToDivert.ID, models.MTOShipmentStatusDiversionRequested, nil, &diversionReason, eTag)
		suite.NoError(suite.DB().Find(&shipmentToDivert, shipmentToDivert.ID))

		suite.NoError(err)
		suite.Equal(models.MTOShipmentStatusDiversionRequested, shipmentToDivert.Status)
	})

	suite.Run("A diversion or diverted shipment can change to APPROVED", func() {
		setupTestData()
		diversionReason := "Test reason"

		// a diversion or diverted shipment is when the PRIME sets the diversion field to true
		// the status must also be in diversion requested status to be approvable as well
		diversionRequestedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:          models.MTOShipmentStatusDiversionRequested,
					Diversion:       true,
					DiversionReason: &diversionReason,
				},
			},
		}, nil)
		eTag = etag.GenerateEtag(diversionRequestedShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), diversionRequestedShipment.ID, models.MTOShipmentStatusApproved, nil, nil, eTag)

		suite.NoError(err)
		suite.NotNil(updatedShipment)
		suite.Equal(models.MTOShipmentStatusApproved, updatedShipment.Status)

		var shipmentServiceItems models.MTOServiceItems
		err = suite.DB().Where("mto_shipment_id = $1", updatedShipment.ID).All(&shipmentServiceItems)
		suite.NoError(err)
		suite.Len(shipmentServiceItems, 0, "should not have created shipment level service items for diversion shipment after approving")
	})
}

func (suite *MTOShipmentServiceSuite) TestMTOShipmentsMTOAvailableToPrime() {
	now := time.Now()
	hide := false
	var primeShipment models.MTOShipment
	var nonPrimeShipment models.MTOShipment
	var hiddenPrimeShipment models.MTOShipment

	setupTestData := func() {
		primeShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
		}, nil)
		nonPrimeShipment = factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		hiddenPrimeShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Show:               &hide,
				},
			},
		}, nil)
	}

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester())
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
	mockSender := setUpMockNotificationSender()
	addressUpdater := address.NewAddressUpdater()
	addressCreator := address.NewAddressCreator()

	updater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

	suite.Run("Shipment exists and is available to Prime - success", func() {
		setupTestData()

		isAvailable, err := updater.MTOShipmentsMTOAvailableToPrime(suite.AppContextForTest(), primeShipment.ID)
		suite.True(isAvailable)
		suite.NoError(err)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Shipment exists but is not available to Prime - failure", func() {
		setupTestData()

		isAvailable, err := updater.MTOShipmentsMTOAvailableToPrime(suite.AppContextForTest(), nonPrimeShipment.ID)
		suite.False(isAvailable)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), nonPrimeShipment.ID.String())

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Shipment exists, is available, but move is disabled - failure", func() {
		setupTestData()

		isAvailable, err := updater.MTOShipmentsMTOAvailableToPrime(suite.AppContextForTest(), hiddenPrimeShipment.ID)
		suite.False(isAvailable)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), hiddenPrimeShipment.ID.String())

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Shipment does not exist - failure", func() {
		setupTestData()

		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		isAvailable, err := updater.MTOShipmentsMTOAvailableToPrime(suite.AppContextForTest(), badUUID)
		suite.False(isAvailable)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), badUUID.String())

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateShipmentEstimatedWeightMoveExcessWeight() {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester())
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
	mockSender := setUpMockNotificationSender()
	addressUpdater := address.NewAddressUpdater()
	addressCreator := address.NewAddressCreator()
	mtoShipmentUpdaterPrime := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

	suite.Run("Updates to estimated weight change max billable weight", func() {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)

		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		suite.Equal(8000, *primeShipment.MoveTaskOrder.Orders.Entitlement.AuthorizedWeight())

		estimatedWeight := unit.Pound(1234)
		primeShipment.Status = ""
		primeShipment.PrimeEstimatedWeight = &estimatedWeight

		session := auth.Session{}
		_, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt), "test")
		suite.NoError(err)

		err = suite.DB().Reload(primeShipment.MoveTaskOrder.Orders.Entitlement)
		suite.NoError(err)

		estimatedWeight110 := int(math.Round(float64(*primeShipment.PrimeEstimatedWeight) * 1.10))
		suite.Equal(estimatedWeight110, *primeShipment.MoveTaskOrder.Orders.Entitlement.AuthorizedWeight())
	})

	suite.Run("Updating the shipment estimated weight will flag excess weight on the move and transitions move status", func() {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)

		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)
		estimatedWeight := unit.Pound(7200)
		// there is a validator check about updating the status
		primeShipment.Status = ""
		primeShipment.PrimeEstimatedWeight = &estimatedWeight

		suite.Nil(primeShipment.MoveTaskOrder.ExcessWeightQualifiedAt)
		suite.Equal(models.MoveStatusAPPROVED, primeShipment.MoveTaskOrder.Status)

		session := auth.Session{}
		_, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt), "test")
		suite.NoError(err)

		err = suite.DB().Reload(&primeShipment.MoveTaskOrder)
		suite.NoError(err)

		suite.NotNil(primeShipment.MoveTaskOrder.ExcessWeightQualifiedAt)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, primeShipment.MoveTaskOrder.Status)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Skips calling check excess weight if estimated weight was not provided in request", func() {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		addressUpdater := address.NewAddressUpdater()

		mockedUpdater := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
		}, nil)
		// there is a validator check about updating the status
		primeShipment.Status = ""
		actualWeight := unit.Pound(7200)
		primeShipment.PrimeActualWeight = &actualWeight

		moveWeights.On("CheckAutoReweigh", mock.AnythingOfType("*appcontext.appContext"), primeShipment.MoveTaskOrderID, mock.AnythingOfType("*models.MTOShipment")).Return(models.MTOShipments{}, nil)

		suite.Nil(primeShipment.MoveTaskOrder.ExcessWeightQualifiedAt)

		session := auth.Session{}
		_, err := mockedUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt), "test")
		suite.NoError(err)

		moveWeights.AssertNotCalled(suite.T(), "CheckExcessWeight")

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Skips calling check excess weight if the updated estimated weight matches the db value", func() {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		addressUpdater := address.NewAddressUpdater()

		mockedUpdater := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(7200)
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
		}, nil)
		// there is a validator check about updating the status
		primeShipment.Status = ""
		primeShipment.PrimeEstimatedWeight = &estimatedWeight

		suite.Nil(primeShipment.MoveTaskOrder.ExcessWeightQualifiedAt)

		session := auth.Session{}
		_, err := mockedUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt), "test")
		suite.Error(err)
		suite.Contains(err.Error(), "cannot be updated after initial estimation")

		moveWeights.AssertNotCalled(suite.T(), "CheckExcessWeight")

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateShipmentActualWeightAutoReweigh() {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester())
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
	mockSender := setUpMockNotificationSender()
	addressUpdater := address.NewAddressUpdater()
	addressCreator := address.NewAddressCreator()
	mtoShipmentUpdaterPrime := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

	suite.Run("Updating the shipment actual weight within weight allowance creates reweigh requests for", func() {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)

		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)
		actualWeight := unit.Pound(7200)
		// there is a validator check about updating the status
		primeShipment.Status = ""
		primeShipment.PrimeActualWeight = &actualWeight

		session := auth.Session{}
		_, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt), "test")
		suite.NoError(err)

		err = suite.DB().Eager("Reweigh").Reload(&primeShipment)
		suite.NoError(err)

		suite.NotNil(primeShipment.Reweigh)
		suite.Equal(primeShipment.ID.String(), primeShipment.Reweigh.ShipmentID.String())
		suite.NotNil(primeShipment.Reweigh.RequestedAt)
		suite.Equal(models.ReweighRequesterSystem, primeShipment.Reweigh.RequestedBy)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Skips calling check auto reweigh if actual weight was not provided in request", func() {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		addressUpdater := address.NewAddressUpdater()

		mockedUpdater := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
		}, nil)
		// there is a validator check about updating the status
		primeShipment.Status = ""
		estimatedWeight := unit.Pound(7200)
		primeShipment.PrimeEstimatedWeight = &estimatedWeight

		moveWeights.On("CheckExcessWeight", mock.AnythingOfType("*appcontext.appContext"), primeShipment.MoveTaskOrderID, mock.AnythingOfType("models.MTOShipment")).Return(&primeShipment.MoveTaskOrder, nil, nil)

		session := auth.Session{}
		_, err := mockedUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt), "test")
		suite.NoError(err)

		moveWeights.AssertNotCalled(suite.T(), "CheckAutoReweigh")

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.Run("Skips calling check auto reweigh if the updated actual weight matches the db value", func() {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		addressUpdater := address.NewAddressUpdater()

		mockedUpdater := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		actualWeight := unit.Pound(7200)
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
					PrimeActualWeight:   &actualWeight,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
		}, nil)
		// there is a validator check about updating the status
		primeShipment.Status = ""
		primeShipment.PrimeActualWeight = &actualWeight

		session := auth.Session{}
		_, err := mockedUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt), "test")
		suite.NoError(err)

		moveWeights.AssertNotCalled(suite.T(), "CheckAutoReweigh")

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateShipmentNullableFields() {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	moveRouter := moveservices.NewMoveRouter()
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)

	suite.Run("tacType and sacType are set to null when empty string is passed in", func() {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		addressUpdater := address.NewAddressUpdater()
		addressCreator := address.NewAddressCreator()
		mockedUpdater := NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		ntsLOAType := models.LOATypeNTS
		ntsMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
					TACType:      &ntsLOAType,
					SACType:      &ntsLOAType,
				},
			},
		}, nil)

		nullLOAType := models.LOAType("")
		requestedUpdate := &models.MTOShipment{
			ID:      ntsMove.MTOShipments[0].ID,
			TACType: &nullLOAType,
			SACType: &nullLOAType,
		}

		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		session.Roles = append(session.Roles, too.User.Roles...)
		_, err := mockedUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), requestedUpdate, etag.GenerateEtag(ntsMove.MTOShipments[0].UpdatedAt), "test")
		suite.NoError(err)
		suite.Equal(nil, nil)
		suite.Equal(nil, nil)
	})

	suite.Run("tacType and sacType are updated when passed in", func() {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()

		addressUpdater := address.NewAddressUpdater()
		addressCreator := address.NewAddressCreator()
		mockedUpdater := NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		ntsLOAType := models.LOATypeNTS
		hhgLOAType := models.LOATypeHHG

		ntsMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
					TACType:      &ntsLOAType,
					SACType:      &ntsLOAType,
				},
			},
		}, nil)
		shipment := ntsMove.MTOShipments[0]

		requestedUpdate := &models.MTOShipment{
			ID:      shipment.ID,
			TACType: &hhgLOAType,
		}

		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		session.Roles = append(session.Roles, too.User.Roles...)
		updatedMtoShipment, err := mockedUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), requestedUpdate, etag.GenerateEtag(shipment.UpdatedAt), "test")
		suite.NoError(err)
		suite.Equal(*requestedUpdate.TACType, *updatedMtoShipment.TACType)
		suite.Equal(*shipment.SACType, *updatedMtoShipment.SACType)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateStatusServiceItems() {

	expectedReServiceCodes := []models.ReServiceCode{
		models.ReServiceCodeDLH,
		models.ReServiceCodeDSH,
		models.ReServiceCodeFSC,
		models.ReServiceCodeDOP,
		models.ReServiceCodeDDP,
		models.ReServiceCodeDPK,
		models.ReServiceCodeDUPK,
	}

	var pickupAddress models.Address
	var longhaulDestinationAddress models.Address
	var shorthaulDestinationAddress models.Address
	var mto models.Move

	setupTestData := func() {
		for i := range expectedReServiceCodes {
			factory.FetchReServiceByCode(suite.DB(), expectedReServiceCodes[i])
		}

		pickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "7 Q St",
					City:           "Twentynine Palms",
					State:          "CA",
					PostalCode:     "92277",
				},
			},
		}, nil)

		longhaulDestinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "278 E Maple Drive",
					City:           "San Diego",
					State:          "CA",
					PostalCode:     "92114",
				},
			},
		}, nil)

		shorthaulDestinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "448 Washington Boulevard NE",
					City:           "Winterhaven",
					State:          "CA",
					PostalCode:     "92283",
				},
			},
		}, nil)

		mto = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)
	}

	builder := query.NewQueryBuilder()
	moveRouter := moveservices.NewMoveRouter()
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
	updater := NewMTOShipmentStatusUpdater(builder, siCreator, planner)

	suite.Run("Shipments with different origin/destination ZIP3 have longhaul service item", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    longhaulDestinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
					Status:       models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		appCtx := suite.AppContextForTest()
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, eTag)
		suite.NoError(err)

		serviceItems := models.MTOServiceItems{}
		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", updatedShipment.ID).All(&serviceItems)
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeDLH, serviceItems[0].ReService.Code)
	})

	suite.Run("Shipments with same origin/destination ZIP3 have shorthaul service item", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    shorthaulDestinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
					Status:       models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		appCtx := suite.AppContextForTest()
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, eTag)
		suite.NoError(err)

		serviceItems := models.MTOServiceItems{}
		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", updatedShipment.ID).All(&serviceItems)
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeDSH, serviceItems[0].ReService.Code)
	})
}

func (suite *MTOShipmentServiceSuite) TestConusInternationalServiceItems() {

	expectedWithConusReServiceCodes := []models.ReServiceCode{
		models.ReServiceCodeISLH,
		models.ReServiceCodePODFSC,
		models.ReServiceCodeIHPK,
		models.ReServiceCodeIHUPK,
	}

	var conusPickupAddress models.Address
	var conusDestinationAddress models.Address
	var mto models.Move

	setupTestData := func() {
		for i := range expectedWithConusReServiceCodes {
			factory.FetchReServiceByCode(suite.DB(), expectedWithConusReServiceCodes[i])
		}

		conus := false
		conusPickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "7 Q St",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       &conus,
				},
			},
		}, nil)

		conusDestinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "278 E Maple Drive",
					City:           "San Diego",
					State:          "CA",
					PostalCode:     "92114",
					IsOconus:       &conus,
				},
			},
		}, nil)

		mto = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)
	}

	builder := query.NewQueryBuilder()
	moveRouter := moveservices.NewMoveRouter()
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
	updater := NewMTOShipmentStatusUpdater(builder, siCreator, planner)

	suite.Run("Shipments without conus pickup/destination locations", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    conusPickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    conusDestinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
					Status:       models.MTOShipmentStatusSubmitted,
					MarketCode:   models.MarketCodeInternational,
				},
			},
		}, nil)

		appCtx := suite.AppContextForTest()
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, eTag)
		suite.NoError(err)

		serviceItems := models.MTOServiceItems{}
		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", updatedShipment.ID).All(&serviceItems)
		suite.NoError(err)

		for i := 0; i < len(expectedWithConusReServiceCodes); i++ {
			suite.Equal(expectedWithConusReServiceCodes[i], serviceItems[i].ReService.Code)
		}
	})
}

func (suite *MTOShipmentServiceSuite) TestWithoutConusInternationalServiceItems() {

	expectedWithoutConusReServiceCodes := []models.ReServiceCode{
		models.ReServiceCodeISLH,
		models.ReServiceCodePOEFSC,
		models.ReServiceCodeIHPK,
		models.ReServiceCodeIHUPK,
	}

	var NoConusPickupAddress models.Address
	var NoConusDestinationAddress models.Address
	var mto models.Move

	setupTestData := func() {
		for i := range expectedWithoutConusReServiceCodes {
			factory.FetchReServiceByCode(suite.DB(), expectedWithoutConusReServiceCodes[i])
		}

		isOconus := true
		NoConusPickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "7 Q St",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       &isOconus,
				},
			},
		}, nil)

		NoConusDestinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "278 E Maple Drive",
					City:           "San Diego",
					State:          "CA",
					PostalCode:     "92114",
					IsOconus:       &isOconus,
				},
			},
		}, nil)

		mto = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)
	}

	builder := query.NewQueryBuilder()
	moveRouter := moveservices.NewMoveRouter()
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
	updater := NewMTOShipmentStatusUpdater(builder, siCreator, planner)

	suite.Run("Shipments without conus pickup/destination locations", func() {
		setupTestData()

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    NoConusPickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    NoConusDestinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
					Status:       models.MTOShipmentStatusSubmitted,
					MarketCode:   models.MarketCodeInternational,
				},
			},
		}, nil)

		appCtx := suite.AppContextForTest()
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, eTag)
		suite.NoError(err)

		serviceItems := models.MTOServiceItems{}
		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", updatedShipment.ID).All(&serviceItems)
		suite.NoError(err)

		for i := 0; i < len(expectedWithoutConusReServiceCodes); i++ {
			suite.Equal(expectedWithoutConusReServiceCodes[i], serviceItems[i].ReService.Code)
		}
	})
}
