package mtoshipment

import (
	"errors"
	"fmt"
	"math"
	"strings"
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
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mockservices "github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
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
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	waf := entitlements.NewWeightAllotmentFetcher()
	mockSender := setUpMockNotificationSender()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()

	mtoShipmentUpdaterOffice := NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)
	mtoShipmentUpdaterCustomer := NewCustomerMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)
	mtoShipmentUpdaterPrime := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)
	scheduledPickupDate := now.Add(time.Hour * 24 * 3)
	firstAvailableDeliveryDate := now.Add(time.Hour * 24 * 4)
	actualPickupDate := now.Add(time.Hour * 24 * 3)
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

	suite.PreloadData(func() {
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
					County:         models.StringPointer("POLK"),
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
			ShipmentType:               models.MTOShipmentTypeHHG,
			PrimeActualWeight:          &primeActualWeight,
			PrimeEstimatedWeight:       &primeEstimatedWeight,
			FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
			Status:                     oldMTOShipment.Status,
			ActualPickupDate:           &actualPickupDate,
			ApprovedDate:               &firstAvailableDeliveryDate,
			MarketCode:                 oldMTOShipment.MarketCode,
		}

		primeEstimatedWeight = unit.Pound(9000)
	})

	setupAllAddressData := func() {
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
					County:         models.StringPointer("POLK"),
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
	}

	suite.Run("Etag is stale", func() {
		eTag := etag.GenerateEtag(time.Now())

		var testScheduledPickupDate time.Time
		mtoShipment.ScheduledPickupDate = &testScheduledPickupDate

		session := auth.Session{}
		_, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &mtoShipment, eTag, "test")
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.Run("404 Not Found Error - shipment can only be created for service member associated with the current session", func() {
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
		eTag := etag.GenerateEtag(oldMTOShipment.UpdatedAt)
		var testScheduledPickupDate time.Time
		mtoShipment.ScheduledPickupDate = &testScheduledPickupDate

		session := auth.Session{}
		updatedMTOShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &mtoShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.Equal(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeHHG)

		suite.Equal(updatedMTOShipment.PickupAddressID, oldMTOShipment.PickupAddressID)

		suite.Equal(updatedMTOShipment.PrimeActualWeight, &primeActualWeight)
		suite.True(actualPickupDate.Equal(*updatedMTOShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*updatedMTOShipment.FirstAvailableDeliveryDate))
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.Run("Updater can handle optional queries set as nil", func() {
		var testScheduledPickupDate time.Time

		oldMTOShipment2 := factory.BuildMTOShipment(suite.DB(), nil, nil)
		mtoShipment2 := models.MTOShipment{
			ID:                  oldMTOShipment2.ID,
			ShipmentType:        models.MTOShipmentTypeHHG,
			ScheduledPickupDate: &testScheduledPickupDate,
		}

		eTag := etag.GenerateEtag(oldMTOShipment2.UpdatedAt)
		session := auth.Session{}
		updatedMTOShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &mtoShipment2, eTag, "test")

		suite.Require().NoError(err)
		suite.Equal(updatedMTOShipment.ID, oldMTOShipment2.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment2.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeHHG)
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(suite.T(), "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.Run("Successfully remove a secondary pickup address", func() {
		secondaryPickupAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

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
		secondaryDeliveryAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
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
		secondaryPickupAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		tertiaryPickupAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
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
		tertiaryDeliveryAddress = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
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
		// Ensure we can update every address field on the shipment
		// Create an mtoShipment to update that has every address populated
		setupAllAddressData()

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

	suite.Run("Successful update on international shipment with estimated weight results in the update of estimated pricing for basic service items", func() {
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"50314",
			"99505",
		).Return(1000, nil)
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"97220",
			"99505",
		).Return(1000, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		pickupUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "50314")
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "Tester Address",
					City:               "Des Moines",
					State:              "IA",
					PostalCode:         "50314",
					IsOconus:           models.BoolPointer(false),
					UsPostRegionCityID: &pickupUSPRC.ID,
				},
			},
		}, nil)

		destUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99505")
		suite.FatalNoError(err)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "Anchorage",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		pickupDate := now.AddDate(0, 0, 10)
		requestedPickup := time.Now()
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: nil,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &pickupDate,
					RequestedPickupDate:  &requestedPickup,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeISLH,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHUPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)
		portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
			{
				Model:    portLocation,
				LinkOnly: true,
				Type:     &factory.PortLocations.PortOfDebarkation,
			},
		}, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		session := auth.Session{}
		_, err = mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")
		suite.NoError(err)

		// checking the service item data
		var serviceItems []models.MTOServiceItem
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)

		suite.Equal(4, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			// because the estimated weight is provided & POEFSC has a port location, estimated pricing should be updated
			suite.NotNil(serviceItems[i].PricingEstimate)
		}
	})

	suite.Run("Successful update on international shipment with estimated weight results in the update of estimated pricing for basic service items except for port fuel surcharge", func() {
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"50314",
			"99505",
		).Return(1000, nil)
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"50314",
			"97220",
		).Return(1000, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		pickupUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "50314")
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "Tester Address",
					City:               "Des Moines",
					State:              "IA",
					PostalCode:         "50314",
					IsOconus:           models.BoolPointer(false),
					UsPostRegionCityID: &pickupUSPRC.ID,
				},
			},
		}, nil)

		destUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99505")
		suite.FatalNoError(err)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "Anchorage",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		pickupDate := now.AddDate(0, 0, 10)
		requestedPickup := time.Now()
		dbShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: nil,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &pickupDate,
					RequestedPickupDate:  &requestedPickup,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    dbShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeISLH,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    dbShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    dbShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHUPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)

		// this will not have a port location and pricing shouldn't be updated
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    dbShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePODFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(dbShipment.UpdatedAt)

		shipment := models.MTOShipment{
			ID:                   dbShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		session := auth.Session{}
		_, err = mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &shipment, eTag, "test")
		suite.NoError(err)

		// checking the service item data
		var serviceItems []models.MTOServiceItem
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", dbShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)

		suite.Equal(4, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			if serviceItems[i].ReService.Code != models.ReServiceCodePODFSC {
				suite.NotNil(serviceItems[i].PricingEstimate)
			} else if serviceItems[i].ReService.Code == models.ReServiceCodePODFSC {
				suite.Nil(serviceItems[i].PricingEstimate)
			}
		}
	})

	suite.Run("Successful update to a minimal MTO shipment", func() {
		// Minimal MTO Shipment has no associated addresses created by default.
		// Part of this test ensures that if an address doesn't exist on a shipment,
		// the updater can successfully create it.
		setupAllAddressData()
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		scheduledDeliveryDate := now.Add(time.Hour * 24 * 4)
		actualDeliveryDate := now.Add(time.Hour * 24 * 4)
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

	suite.Run("Returns error if updated UB shipment addresses are both CONUS", func() {
		conusAddress := factory.BuildAddress(suite.DB(), nil, nil)

		oconusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Cordova",
					State:          "AK",
					PostalCode:     "99677",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		// UB shipment with an OCONUS pickup & a CONUS destination
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeUnaccompaniedBaggage,
				},
			},
			{
				Model:    oconusAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    conusAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		// updating pickup to be CONUS which should return an error because
		// UBs must have an OCONUS address
		updatedShipment := models.MTOShipment{
			ID:            oldShipment.ID,
			PickupAddress: &conusAddress,
		}

		session := auth.Session{}
		newShipment, err := mtoShipmentUpdaterCustomer.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Error(err)
		suite.Nil(newShipment)

		var invalidErr apperror.InvalidInputError
		suite.True(errors.As(err, &invalidErr), "error should be of type InvalidInputError")

		if invalidErr.ValidationErrors == nil {
			suite.Fail("ValidationErrors is nil")
		}

		ve := *invalidErr.ValidationErrors
		fieldErrors, exists := ve.Errors["UB shipment error"]
		suite.True(exists, "expected validation error for 'UB shipment error'")
		joinedErrors := strings.Join(fieldErrors, " ")
		suite.Contains(joinedErrors, "At least one address for a UB shipment must be OCONUS")

	})

	suite.Run("Updating a shipment does not nullify ApprovedDate", func() {
		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		setupAllAddressData()
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
		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		setupAllAddressData()
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
		defaultRole, err := too.User.Roles.Default()
		suite.FatalNoError(err)
		session.ActiveRole = *defaultRole
		newShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.Equal(destinationType, *newShipment.DestinationType)
	})

	suite.Run("Successfully update MTO Agents", func() {
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
		defaultRole, err := too.User.Roles.Default()
		suite.FatalNoError(err)
		session.ActiveRole = *defaultRole
		updatedMTOShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.StorageFacility)
	})

	suite.Run("Successfully edit storage facility on shipment", func() {
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
		defaultRole, err := too.User.Roles.Default()
		suite.FatalNoError(err)
		session.ActiveRole = *defaultRole
		updatedShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &newShipment, eTag, "test")
		suite.Require().NoError(err)
		suite.NotEqual(uuid.Nil, updatedShipment.ID)
		suite.Equal(&newEmail, updatedShipment.StorageFacility.Email)
		suite.Equal(newStorageFacilityAddress.StreetAddress1, updatedShipment.StorageFacility.Address.StreetAddress1)
	})

	suite.Run("Successfully update NTS previously recorded weight to shipment", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		ntsRecorededWeight := unit.Pound(980)
		updatedShipment := models.MTOShipment{
			ShipmentType:      models.MTOShipmentTypeHHGOutOfNTS,
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
		defaultRole, err := too.User.Roles.Default()
		suite.FatalNoError(err)
		session.ActiveRole = *defaultRole
		updatedMTOShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(ntsRecorededWeight, *updatedMTOShipment.NTSRecordedWeight)

	})

	suite.Run("Unable to update NTS previously recorded weight due to shipment type", func() {
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
		defaultRole, err := too.User.Roles.Default()
		suite.FatalNoError(err)
		session.ActiveRole = *defaultRole
		updatedMTOShipment, err := mtoShipmentUpdaterOffice.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

		suite.Require().Error(err)
		suite.Nil(updatedMTOShipment)
		suite.Equal("Invalid input found while updating the shipment", err.Error())

		var invalidErr apperror.InvalidInputError
		suite.True(errors.As(err, &invalidErr), "error should be of type InvalidInputError")

		if invalidErr.ValidationErrors == nil {
			suite.Fail("ValidationErrors is nil")
		}

		ve := *invalidErr.ValidationErrors
		fieldErrors, exists := ve.Errors["NTSRecordedWeight error"]
		suite.True(exists, "expected validation error for 'NTSRecordedWeight error'")
		joinedErrors := strings.Join(fieldErrors, " ")
		suite.Contains(joinedErrors, "field NTSRecordedWeight cannot be set for shipment type HHG")
	})

	suite.Run("Successfully divert a shipment and transition statuses", func() {
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
		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		setupAllAddressData()
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
		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		setupAllAddressData()
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
		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		setupAllAddressData()
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

	suite.Run("Successful Office/TOO UpdateShipment - CONUS Pickup, OCONUS Destination - mileage is recalculated and pricing estimates refreshed for International FSC SIT service items", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER1",
					City:           "Anchorage1",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		pickupDate := now.AddDate(0, 0, 10)
		requestedPickup := time.Now()
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: nil,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &pickupDate,
					RequestedPickupDate:  &requestedPickup,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// setup IOSFSC service item with SITOriginHHGOriginalAddress
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOSFSC,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)

		// setup IDSFSC service item with SITDestinationOriginalAddress
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDSFSC,
				},
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.SITDestinationOriginalAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
				LinkOnly: true,
			},
		}, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		var serviceItems []models.MTOServiceItem
		// verify pre-update mto service items for both origin/destination FSC SITs have not been set
		err := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		// expecting only IOSFSC and IDSFSC created for tests
		suite.Equal(2, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.Nil(serviceItems[i].PricingEstimate)
			suite.True(serviceItems[i].SITDeliveryMiles == (*int)(nil))
		}

		// As TOO
		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
			ActiveRole:      too.User.Roles[0],
		}
		expectedMileage := 314
		plannerSITFSC := &mocks.Planner{}
		// expecting 50314/50314 for IOSFSC mileage lookup for source, destination
		plannerSITFSC.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			// 99505/99505, 50314/50314
			mock.MatchedBy(func(source string) bool {
				return source == "50314" || source == "99505"
			}),
			mock.MatchedBy(func(destination string) bool {
				return destination == "50314" || destination == "99505"
			}),
		).Return(expectedMileage, nil)

		mtoShipmentUpdater := NewOfficeMTOShipmentUpdater(builder, fetcher, plannerSITFSC, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		_, err = mtoShipmentUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")
		suite.NoError(err)

		// verify post-update mto service items for both origin/destination FSC SITs have been set.
		// if set we know stored procedure update_service_item_pricing was executed sucessfully
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(2, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOSFSC || serviceItems[i].ReService.Code == models.ReServiceCodeIDSFSC)

			if serviceItems[i].ReService.Code == models.ReServiceCodeIOSFSC {
				suite.NotNil(*serviceItems[i].PricingEstimate)
				suite.Equal(*serviceItems[i].SITDeliveryMiles, expectedMileage)
			}
			// verify IDSFSC SIT with OCONUS destination does not calculate pricing resulting in 0.
			if serviceItems[i].ReService.Code == models.ReServiceCodeIDSFSC {
				suite.Equal(*serviceItems[i].SITDeliveryMiles, expectedMileage)
				suite.Equal(*serviceItems[i].PricingEstimate, unit.Cents(0))
			}
		}
	})

	suite.Run("Successful Office/TOO UpdateShipment - OCONUS Pickup, CONUS Destination - mileage is recalculated and pricing estimates refreshed for International FSC SIT service items", func() {
		setupAllAddressData()
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER1",
					City:           "Anchorage1",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		pickupDate := now.AddDate(0, 0, 10)
		requestedPickup := time.Now()
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: nil,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &pickupDate,
					RequestedPickupDate:  &requestedPickup,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// setup IOSFSC service item with SITOriginHHGOriginalAddress
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOSFSC,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)

		// setup IDSFSC service item with SITDestinationOriginalAddress
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDSFSC,
				},
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.SITDestinationOriginalAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
				LinkOnly: true,
			},
		}, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		var serviceItems []models.MTOServiceItem
		// verify pre-update mto service items for both origin/destination FSC SITs have not been set
		err := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		// expecting only IOSFSC and IDSFSC created for tests
		suite.Equal(2, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.Nil(serviceItems[i].PricingEstimate)
			suite.True(serviceItems[i].SITDeliveryMiles == (*int)(nil))
		}

		// As TOO
		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
			ActiveRole:      too.User.Roles[0],
		}
		expectedMileage := 314
		plannerSITFSC := &mocks.Planner{}
		// expecting 99505/99505, 50314/50314 for IOSFSC mileage lookup for source, destination
		plannerSITFSC.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(source string) bool {
				return source == "50314" || source == "99505"
			}),
			mock.MatchedBy(func(destination string) bool {
				return destination == "50314" || destination == "99505"
			}),
		).Return(expectedMileage, nil)

		mtoShipmentUpdater := NewOfficeMTOShipmentUpdater(builder, fetcher, plannerSITFSC, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		_, err = mtoShipmentUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")
		suite.NoError(err)

		// verify post-update mto service items for both origin/destination FSC SITs have been set.
		// if set we know stored procedure update_service_item_pricing was executed sucessfully
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(2, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOSFSC || serviceItems[i].ReService.Code == models.ReServiceCodeIDSFSC)

			if serviceItems[i].ReService.Code == models.ReServiceCodeIDSFSC {
				suite.NotNil(*serviceItems[i].PricingEstimate)
				suite.Equal(*serviceItems[i].SITDeliveryMiles, expectedMileage)
			}
			// verify IOSFSC SIT with OCONUS destination does not calculate mileage and pricing resulting in 0 for both.
			if serviceItems[i].ReService.Code == models.ReServiceCodeIOSFSC {
				suite.Equal(*serviceItems[i].SITDeliveryMiles, expectedMileage)
				suite.Equal(*serviceItems[i].PricingEstimate, unit.Cents(0))
			}
		}
	})

	suite.Run("Successful Office/TOO UpdateShipment - Pricing estimates calculated for Intl First Day SIT Service Items (IOFSIT, IDFSIT)", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "450 Street Dr",
					City:           "Charleston",
					State:          "SC",
					PostalCode:     "29404",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JB Snowtown",
					City:           "Juneau",
					State:          "AK",
					PostalCode:     "99801",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		pickupDate := now.AddDate(0, 0, 10)
		requestedPickup := time.Now()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: nil,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &pickupDate,
					RequestedPickupDate:  &requestedPickup,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// setup IOFSIT service item with SITOriginHHGOriginalAddress
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOFSIT,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)

		// setup IDFSIT service item
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDFSIT,
				},
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.SITDestinationOriginalAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
				LinkOnly: true,
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedShipment := models.MTOShipment{
			ID:                   shipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		var serviceItems []models.MTOServiceItem
		// verify pre-update mto service items for both origin/destination First Day SITs have not been set
		err := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		// expecting only IOFSIT and IDFSIT created for tests
		suite.Equal(2, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.Nil(serviceItems[i].PricingEstimate)
		}

		// As TOO
		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
			ActiveRole:      too.User.Roles[0],
		}
		plannerSITFSC := &mocks.Planner{}
		plannerSITFSC.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1, nil)

		mtoShipmentUpdater := NewOfficeMTOShipmentUpdater(builder, fetcher, plannerSITFSC, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		_, err = mtoShipmentUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")
		suite.NoError(err)

		// verify post-update mto service items for both origin/destination First Day SITs have been set.
		// if set we know stored procedure update_service_item_pricing was executed sucessfully
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(2, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOFSIT || serviceItems[i].ReService.Code == models.ReServiceCodeIDFSIT)
			suite.True(*serviceItems[i].PricingEstimate > 0)
		}
	})

	suite.Run("Successful Office/TOO UpdateShipment - no sit departure date - pricing estimates refreshed using MAX days for International Additional Days SIT service item", func() {
		parameter := models.ApplicationParameters{
			ParameterName:  models.StringPointer("maxSitDaysAllowance"),
			ParameterValue: models.StringPointer("90"),
		}
		suite.MustCreate(&parameter)

		testdatagen.FetchOrMakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate: testdatagen.ContractStartDate,
					EndDate:   testdatagen.ContractEndDate,
				},
			})

		usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99801")
		suite.NotNil(usprc)
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &usprc.ID,
					City:               usprc.USPostRegionCityNm,
					State:              usprc.State,
					PostalCode:         usprc.UsprZipID,
				},
			},
		}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		pickupDate := now.AddDate(0, 0, 10)
		requestedPickup := time.Now()
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: nil,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &pickupDate,
					RequestedPickupDate:  &requestedPickup,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    oldShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOASIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
					SITEntryDate:    &nowDate,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		// As TOO
		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
			ActiveRole:      too.User.Roles[0],
		}

		var serviceItems []models.MTOServiceItem
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOASIT)
			suite.NotNil(serviceItems[i].SITEntryDate)
			suite.Nil(serviceItems[i].SITDepartureDate)
			suite.Nil(serviceItems[i].PricingEstimate)
		}

		mtoShipmentUpdater := NewOfficeMTOShipmentUpdater(builder, fetcher, &mocks.Planner{}, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		updateShipment2, err := mtoShipmentUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")
		suite.NoError(err)

		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOASIT)
			suite.NotNil(serviceItems[i].SITEntryDate)
			suite.Nil(serviceItems[i].SITDepartureDate)
			suite.NotNil(*serviceItems[i].PricingEstimate)
		}

		var pricingEstimateWithMaxSitDays *unit.Cents
		// Set SIT Departure date
		serviceItems[0].SITDepartureDate = models.TimePointer(serviceItems[0].SITEntryDate.Add(time.Hour * 48))
		err = suite.AppContextForTest().DB().Update(&serviceItems[0])
		suite.NoError(err)
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOASIT)
			suite.NotNil(serviceItems[i].SITEntryDate)
			suite.NotNil(serviceItems[i].SITDepartureDate)
			suite.NotNil(*serviceItems[i].PricingEstimate)
			pricingEstimateWithMaxSitDays = serviceItems[i].PricingEstimate
		}

		eTag = etag.GenerateEtag(updateShipment2.UpdatedAt)
		updatedShipment = models.MTOShipment{
			ID:                   updateShipment2.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		var pricingEstimateWithOutMaxSitDays *unit.Cents
		_, err = mtoShipmentUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")
		suite.NoError(err)
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", oldShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOASIT)
			suite.NotNil(serviceItems[i].SITEntryDate)
			suite.NotNil(serviceItems[i].SITDepartureDate)
			suite.NotNil(*serviceItems[i].PricingEstimate)
			pricingEstimateWithOutMaxSitDays = serviceItems[i].PricingEstimate
		}

		// verify pricing is larger for smaller sit in days calculation versus one with default of 89
		suite.True(pricingEstimateWithMaxSitDays.Int() > pricingEstimateWithOutMaxSitDays.Int())
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

	suite.PreloadData(func() {
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
	})

	builder := query.NewQueryBuilder()
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
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

		appCtx := suite.AppContextForTest()
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.UpdatedAt)
		serviceItems := models.MTOServiceItems{}

		preApprovalTime := time.Now()
		approvedShipment, err := updater.UpdateMTOShipmentStatus(appCtx, shipmentForAutoApprove.ID, status, nil, nil, shipmentForAutoApproveEtag)
		suite.NoError(err)
		// Let's make sure the status is approved
		suite.Equal(models.MTOShipmentStatusApproved, approvedShipment.Status)

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
		// MTOShipmentTypeHHGIntoNTS: ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, StorageFacility
		// MTOShipmentTypeHHGOutOfNTS: ScheduledPickupDate, NTSRecordedWeight, StorageFacility, DestinationAddress
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
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTS,
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
					ShipmentType:        models.MTOShipmentTypeHHGOutOfNTS,
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

	suite.Run("Test that we are properly adding days to Alaska shipments", func() {
		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().Add(time.Hour * 12),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
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

		conusAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		zone1Address := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddressAKZone1})
		zone2Address := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddressAKZone2})
		zone3Address := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddressAKZone3})
		zone4Address := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddressAKZone4})
		zone5Address := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddressAKZone5})

		estimatedWeight := unit.Pound(11000)

		testCases10Days := []struct {
			pickupLocation      models.Address
			destinationLocation models.Address
		}{
			{conusAddress, zone1Address},
			{conusAddress, zone2Address},
			{zone1Address, conusAddress},
			{zone2Address, conusAddress},
		}
		// adding 22 days; ghcDomesticTransitTime0LbsUpper.MaxDaysTransitTime is 12, plus 10 for Zones 1 and 2
		rdd10DaysDate := testdatagen.DateInsidePeakRateCycle.AddDate(0, 0, 22)
		for _, testCase := range testCases10Days {
			shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
				{
					Model:    move,
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
					Model:    testCase.pickupLocation,
					Type:     &factory.Addresses.PickupAddress,
					LinkOnly: true,
				},
				{
					Model:    testCase.destinationLocation,
					Type:     &factory.Addresses.DeliveryAddress,
					LinkOnly: true,
				},
			}, nil)
			shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
			_, err = updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, status, nil, nil, shipmentEtag)
			suite.NoError(err)

			fetchedShipment := models.MTOShipment{}
			err = suite.DB().Find(&fetchedShipment, shipment.ID)
			suite.NoError(err)
			suite.NotNil(fetchedShipment.RequiredDeliveryDate)
			suite.Equal(rdd10DaysDate.Format(time.RFC3339), fetchedShipment.RequiredDeliveryDate.Format(time.RFC3339))
		}

		testCases20Days := []struct {
			pickupLocation      models.Address
			destinationLocation models.Address
		}{
			{conusAddress, zone3Address},
			{conusAddress, zone4Address},
			{zone3Address, conusAddress},
			{zone4Address, conusAddress},
		}
		// adding 32 days; ghcDomesticTransitTime0LbsUpper.MaxDaysTransitTime is 12, plus 20 for Zones 3 and 4
		rdd20DaysDate := testdatagen.DateInsidePeakRateCycle.AddDate(0, 0, 32)
		for _, testCase := range testCases20Days {
			shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
				{
					Model:    move,
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
					Model:    testCase.pickupLocation,
					Type:     &factory.Addresses.PickupAddress,
					LinkOnly: true,
				},
				{
					Model:    testCase.destinationLocation,
					Type:     &factory.Addresses.DeliveryAddress,
					LinkOnly: true,
				},
			}, nil)
			shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
			_, err = updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, status, nil, nil, shipmentEtag)
			suite.NoError(err)

			fetchedShipment := models.MTOShipment{}
			err = suite.DB().Find(&fetchedShipment, shipment.ID)
			suite.NoError(err)
			suite.NotNil(fetchedShipment.RequiredDeliveryDate)
			suite.Equal(rdd20DaysDate.Format(time.RFC3339), fetchedShipment.RequiredDeliveryDate.Format(time.RFC3339))
		}
		testCases60Days := []struct {
			pickupLocation      models.Address
			destinationLocation models.Address
		}{
			{conusAddress, zone5Address},
			{zone5Address, conusAddress},
		}

		// adding 72 days; ghcDomesticTransitTime0LbsUpper.MaxDaysTransitTime is 12, plus 60 for Zone 5 HHG
		rdd60DaysDate := testdatagen.DateInsidePeakRateCycle.AddDate(0, 0, 72)
		for _, testCase := range testCases60Days {
			shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
				{
					Model:    move,
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
					Model:    testCase.pickupLocation,
					Type:     &factory.Addresses.PickupAddress,
					LinkOnly: true,
				},
				{
					Model:    testCase.destinationLocation,
					Type:     &factory.Addresses.DeliveryAddress,
					LinkOnly: true,
				},
			}, nil)
			shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
			_, err = updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, status, nil, nil, shipmentEtag)
			suite.NoError(err)

			fetchedShipment := models.MTOShipment{}
			err = suite.DB().Find(&fetchedShipment, shipment.ID)
			suite.NoError(err)
			suite.NotNil(fetchedShipment.RequiredDeliveryDate)
			suite.Equal(rdd60DaysDate.Format(time.RFC3339), fetchedShipment.RequiredDeliveryDate.Format(time.RFC3339))
		}

		conusAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					City:           "Charlotte",
					State:          "NC",
					PostalCode:     "28290",
					IsOconus:       models.BoolPointer(false),
				},
			}}, nil)
		zone5Address = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Cordova",
					State:          "AK",
					PostalCode:     "99677",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		testCases60Days = []struct {
			pickupLocation      models.Address
			destinationLocation models.Address
		}{
			{conusAddress, zone5Address},
			{zone5Address, conusAddress},
		}

		for _, testCase := range testCases60Days {
			shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model: models.MTOShipment{
						ShipmentType:         models.MTOShipmentTypeUnaccompaniedBaggage,
						ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
						PrimeEstimatedWeight: &estimatedWeight,
						Status:               models.MTOShipmentStatusSubmitted,
					},
				},
				{
					Model:    testCase.pickupLocation,
					Type:     &factory.Addresses.PickupAddress,
					LinkOnly: true,
				},
				{
					Model:    testCase.destinationLocation,
					Type:     &factory.Addresses.DeliveryAddress,
					LinkOnly: true,
				},
			}, nil)
			// adding 42 days; ghcDomesticTransitTime0LbsUpper.MaxDaysTransitTime is 12, plus 30 for Zone 5 UB
			pickUpDate := shipment.ScheduledPickupDate
			rdd60DaysDateUB := pickUpDate.AddDate(0, 0, 27)
			shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
			_, err = updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, status, nil, nil, shipmentEtag)
			suite.NoError(err)

			fetchedShipment := models.MTOShipment{}
			err = suite.DB().Find(&fetchedShipment, shipment.ID)
			suite.NoError(err)
			suite.NotNil(fetchedShipment.RequiredDeliveryDate)
			suite.Equal(rdd60DaysDateUB.Format(time.RFC3339), fetchedShipment.RequiredDeliveryDate.Format(time.RFC3339))
		}
	})

	suite.Run("Update RDD on UB Shipment on status change", func() {
		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		appCtx := suite.AppContextForTest()

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

		conusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					City:           "Charlotte",
					State:          "NC",
					PostalCode:     "28290",
					IsOconus:       models.BoolPointer(false),
				},
			}}, nil)
		zone5Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Cordova",
					State:          "AK",
					PostalCode:     "99677",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)
		estimatedWeight := unit.Pound(4000)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeUnaccompaniedBaggage,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    conusAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    zone5Address,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		mtoShipment, err := updater.UpdateMTOShipmentStatus(appCtx, shipment.ID, status, nil, nil, shipmentEtag)
		suite.NoError(err)
		suite.NotNil(mtoShipment.RequiredDeliveryDate)
		suite.False(mtoShipment.RequiredDeliveryDate.IsZero())
	})

	suite.Run("Cannot set SUBMITTED status on shipment via UpdateMTOShipmentStatus", func() {
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
		eTag = etag.GenerateEtag(shipment3.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment3.ID, "REJECTED", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("Rejecting a shipment in APPROVED status returns a ConflictStatusError", func() {
		eTag = etag.GenerateEtag(approvedShipment.UpdatedAt)
		rejectionReason := "Rejection reason"
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), approvedShipment.ID, "REJECTED", &rejectionReason, nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Approving a shipment in REJECTED status returns a ConflictStatusError", func() {
		eTag = etag.GenerateEtag(rejectedShipment.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), rejectedShipment.ID, "APPROVED", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		staleETag := etag.GenerateEtag(time.Now())

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment4.ID, "APPROVED", nil, nil, staleETag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Passing in an invalid status returns a ConflictStatus error", func() {
		eTag = etag.GenerateEtag(shipment4.UpdatedAt)

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment4.ID, "invalid", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a bad shipment id returns a Not Found error", func() {
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), badShipmentID, "APPROVED", nil, nil, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Changing to APPROVED status records approved_date", func() {
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

		approvedShipment, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment5.ID, models.MTOShipmentStatusApproved, nil, nil, eTag)
		suite.NoError(err)
		suite.Equal(models.MTOShipmentStatusApproved, approvedShipment.Status)
		suite.NotNil(approvedShipment.ApprovedDate)
	})

	suite.Run("Changing to a non-APPROVED status does not record approved_date", func() {
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
		eTag = etag.GenerateEtag(approvedShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), approvedShipment.ID, models.MTOShipmentStatusCanceled, nil, nil, eTag)
		suite.NoError(suite.DB().Find(&approvedShipment, approvedShipment.ID))

		suite.Error(err)
		suite.Nil(updatedShipment)
		suite.IsType(ConflictStatusError{}, err)
		suite.Equal(models.MTOShipmentStatusApproved, approvedShipment.Status)
	})

	suite.Run("An APPROVALS_REQUESTED shipment CANNOT change to CANCELED - ERROR", func() {
		testShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitApprovalsRequestedShipment})
		eTag = etag.GenerateEtag(testShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), testShipment.ID, models.MTOShipmentStatusCanceled, nil, nil, eTag)
		suite.NoError(suite.DB().Find(&testShipment, testShipment.ID))

		suite.Error(err)
		suite.Nil(updatedShipment)
		suite.IsType(ConflictStatusError{}, err)
		suite.Equal(models.MTOShipmentStatusApprovalsRequested, testShipment.Status)
	})

	suite.Run("An APPROVED shipment CAN change to Diversion Requested", func() {
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
	waf := entitlements.NewWeightAllotmentFetcher()

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
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	mockSender := setUpMockNotificationSender()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
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
	waf := entitlements.NewWeightAllotmentFetcher()

	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	mockSender := setUpMockNotificationSender()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
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

		moveWeights.On("CheckAutoReweigh", mock.AnythingOfType("*appcontext.appContext"), primeShipment.MoveTaskOrderID, mock.AnythingOfType("*models.MTOShipment")).Return(nil)

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
	waf := entitlements.NewWeightAllotmentFetcher()

	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	mockSender := setUpMockNotificationSender()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
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
		moveWeights.On("CheckAutoReweigh",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil)
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

		moveWeights.On("CheckExcessWeight",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.MTOShipment"),
		).Return(&primeShipment.MoveTaskOrder, nil, nil)

		// there is a validator check about updating the status
		primeShipment.Status = ""

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
		weight := unit.Pound(7200)
		oldPrimeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					PrimeActualWeight:    &weight,
					PrimeEstimatedWeight: &weight,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
		}, nil)

		moveWeights.On("CheckExcessWeight",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.MTOShipment"),
		).Return(&oldPrimeShipment.MoveTaskOrder, nil, nil)

		newPrimeShipment := models.MTOShipment{
			ID:                oldPrimeShipment.ID,
			PrimeActualWeight: &weight,
		}

		eTag := etag.GenerateEtag(oldPrimeShipment.UpdatedAt)

		session := auth.Session{}
		_, err := mockedUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &newPrimeShipment, eTag, "test")
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
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)

	suite.Run("tacType and sacType are set to null when empty string is passed in", func() {
		moveWeights := &mockservices.MoveWeights{}
		moveWeights.On("CheckAutoReweigh",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil)

		mockSender := setUpMockNotificationSender()
		addressUpdater := address.NewAddressUpdater()
		addressCreator := address.NewAddressCreator()
		mockedUpdater := NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		ntsLOAType := models.LOATypeNTS
		ntsMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGIntoNTS,
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
		defaultRole, err := too.User.Roles.Default()
		suite.FatalNoError(err)
		session.ActiveRole = *defaultRole
		_, err = mockedUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), requestedUpdate, etag.GenerateEtag(ntsMove.MTOShipments[0].UpdatedAt), "test")
		suite.NoError(err)
		suite.Equal(nil, nil)
		suite.Equal(nil, nil)
	})

	suite.Run("tacType and sacType are updated when passed in", func() {
		moveWeights := &mockservices.MoveWeights{}
		moveWeights.On("CheckAutoReweigh",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil)
		mockSender := setUpMockNotificationSender()

		addressUpdater := address.NewAddressUpdater()
		addressCreator := address.NewAddressCreator()
		mockedUpdater := NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		ntsLOAType := models.LOATypeNTS
		hhgLOAType := models.LOATypeHHG

		ntsMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGIntoNTS,
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
		defaultRole, err := too.User.Roles.Default()
		suite.FatalNoError(err)
		session.ActiveRole = *defaultRole
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
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
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

		foundDLH := false
		for _, serviceItem := range serviceItems {
			if serviceItem.ReService.Code == models.ReServiceCodeDLH {
				foundDLH = true
				break
			}
		}

		// at least one service item should have the DLH code
		suite.True(foundDLH, "Expected to find at least one service item with ReService code DLH")
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

		isTestMatch := false
		for _, serviceItem := range serviceItems {
			if serviceItem.ReService.Code == models.ReServiceCodeDSH {
				isTestMatch = true
			}
		}
		suite.True(isTestMatch)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateDomesticServiceItems() {

	expectedReServiceCodes := []models.ReServiceCode{
		models.ReServiceCodeDLH,
		models.ReServiceCodeFSC,
		models.ReServiceCodeDOP,
		models.ReServiceCodeDDP,
		models.ReServiceCodeDNPK,
	}

	var pickupAddress models.Address
	var storageFacility models.StorageFacility
	var mto models.Move

	setupTestData := func() {
		pickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Test Street 1",
					City:           "Des moines",
					State:          "IA",
					PostalCode:     "50309",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		storageFacility = factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Test Street Adress 2",
					City:           "Des moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
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
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
		false,
	).Return(400, nil)
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
	updater := NewMTOShipmentStatusUpdater(builder, siCreator, planner)

	suite.Run("Preapproved service items successfully added to domestic nts shipments", func() {
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
				Model:    storageFacility,
				Type:     &factory.StorageFacility,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGIntoNTS,
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

		actualReServiceCodes := []models.ReServiceCode{}
		for _, item := range serviceItems {
			actualReServiceCodes = append(actualReServiceCodes, item.ReService.Code)
		}

		suite.ElementsMatch(expectedReServiceCodes, actualReServiceCodes)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateRequiredDeliveryDateUpdate() {

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	waf := entitlements.NewWeightAllotmentFetcher()
	mockSender := setUpMockNotificationSender()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()

	suite.Run("should update requiredDeliveryDate when scheduledPickupDate is updated", func() {
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(500, nil)
		mtoShipmentUpdaterPrime := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		appCtx := suite.AppContextForTest()

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

		conusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					City:           "Charlotte",
					State:          "NC",
					PostalCode:     "28290",
					IsOconus:       models.BoolPointer(false),
				},
			}}, nil)
		zone5Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Cordova",
					State:          "AK",
					PostalCode:     "99677",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)
		estimatedWeight := unit.Pound(4000)
		oldUbShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeUnaccompaniedBaggage,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusApproved,
					PrimeActualWeight:    &estimatedWeight,
				},
			},
			{
				Model:    conusAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    zone5Address,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		suite.Nil(oldUbShipment.RequiredDeliveryDate)

		pickUpDate := time.Now()
		expectedRequiredDeiliveryDate := pickUpDate.AddDate(0, 0, 27)
		newUbShipment := models.MTOShipment{
			ID:                  oldUbShipment.ID,
			ShipmentType:        models.MTOShipmentTypeUnaccompaniedBaggage,
			ScheduledPickupDate: &pickUpDate,
		}

		eTag := etag.GenerateEtag(oldUbShipment.UpdatedAt)
		updatedMTOShipment, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(appCtx, &newUbShipment, eTag, "test")

		suite.Nil(err)
		suite.NotNil(updatedMTOShipment)
		suite.NotNil(updatedMTOShipment.RequiredDeliveryDate)
		suite.False(updatedMTOShipment.RequiredDeliveryDate.IsZero())
		suite.Equal(expectedRequiredDeiliveryDate.Day(), updatedMTOShipment.RequiredDeliveryDate.Day())
		suite.Equal(expectedRequiredDeiliveryDate.Month(), updatedMTOShipment.RequiredDeliveryDate.Month())
		suite.Equal(expectedRequiredDeiliveryDate.Year(), updatedMTOShipment.RequiredDeliveryDate.Year())
	})

	suite.Run("AK -> AK - should update requiredDeliveryDate when scheduledPickupDate is updated", func() {
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(500, nil)
		mtoShipmentUpdaterPrime := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		appCtx := suite.AppContextForTest()

		zone1Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "24850 Gratiot Dr,",
					City:           "JBER",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)
		zone4Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Cordova",
					State:          "AK",
					PostalCode:     "99677",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)
		estimatedWeight := unit.Pound(400)

		oldUbShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeUnaccompaniedBaggage,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusApproved,
					PrimeActualWeight:    &estimatedWeight,
				},
			},
			{
				Model:    zone1Address,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    zone4Address,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)
		suite.Nil(oldUbShipment.RequiredDeliveryDate)

		pickUpDate := time.Now()
		expectedRequiredDeiliveryDateUB := pickUpDate.AddDate(0, 0, 30)
		newUbShipment := models.MTOShipment{
			ID:                  oldUbShipment.ID,
			ShipmentType:        models.MTOShipmentTypeUnaccompaniedBaggage,
			ScheduledPickupDate: &pickUpDate,
		}

		eTag := etag.GenerateEtag(oldUbShipment.UpdatedAt)
		updatedUBMTOShipment, err := mtoShipmentUpdaterPrime.UpdateMTOShipment(appCtx, &newUbShipment, eTag, "test")

		suite.Nil(err)
		suite.NotNil(updatedUBMTOShipment)
		suite.NotNil(updatedUBMTOShipment.RequiredDeliveryDate)
		suite.False(updatedUBMTOShipment.RequiredDeliveryDate.IsZero())
		suite.Equal(expectedRequiredDeiliveryDateUB.Day(), updatedUBMTOShipment.RequiredDeliveryDate.Day())
		suite.Equal(expectedRequiredDeiliveryDateUB.Month(), updatedUBMTOShipment.RequiredDeliveryDate.Month())
		suite.Equal(expectedRequiredDeiliveryDateUB.Year(), updatedUBMTOShipment.RequiredDeliveryDate.Year())
	})
}

func (suite *MTOShipmentServiceSuite) TestCalculateRequiredDeliveryDate() {
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(500, nil)

	suite.Run("CONUS -> AK - errors when rate area for the pickup address is not found", func() {
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

		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		conusAddressId := uuid.Must(uuid.NewV4())
		conusAddress := models.Address{
			ID:             conusAddressId,
			StreetAddress1: "1 some street",
			City:           "Charlotte",
			State:          "NC",
			PostalCode:     "28290",
			IsOconus:       models.BoolPointer(false),
		}

		zone5Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Cordova",
					State:          "AK",
					PostalCode:     "99677",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		RDD, err := CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, conusAddress, zone5Address, time.Now(), models.IntPointer(500), move.ID, models.MTOShipmentTypeUnaccompaniedBaggage)
		suite.NotNil(err)
		suite.Nil(RDD)
		suite.Equal(fmt.Sprintf("error fetching pickup rate area id for address ID: %s", conusAddress.ID), err.Error())
	})

	suite.Run("CONUS -> AK - errors when rate area for the destination address is not found", func() {
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

		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		conusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					City:           "Charlotte",
					State:          "NC",
					PostalCode:     "28290",
					IsOconus:       models.BoolPointer(false),
				},
			}}, nil)

		zone5AddressId := uuid.Must(uuid.NewV4())
		zone5Address := models.Address{
			ID:             zone5AddressId,
			StreetAddress1: "1 some street",
			StreetAddress2: models.StringPointer("P.O. Box 1234"),
			StreetAddress3: models.StringPointer("c/o Another Person"),
			City:           "Cordova",
			State:          "AK",
			PostalCode:     "99677",
			IsOconus:       models.BoolPointer(true),
		}

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		RDD, err := CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, conusAddress, zone5Address, time.Now(), models.IntPointer(500), move.ID, models.MTOShipmentTypeUnaccompaniedBaggage)
		suite.NotNil(err)
		suite.Nil(RDD)
		suite.Equal(fmt.Sprintf("error fetching destination rate area id for address ID: %s", zone5Address.ID), err.Error())
	})

	suite.Run("AK -> CONUS - correctly errors when weight is missing for HHG with one CONUS address", func() {
		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		conusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1 some street",
					City:           "Charlotte",
					State:          "NC",
					PostalCode:     "28290",
				},
			}}, nil)

		zone2Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "200 Sitka Spruce St",
					City:           "Eielson AFB",
					State:          "AK",
					PostalCode:     "99702",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		RDD, err := CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, zone2Address, conusAddress, time.Now(), nil, move.ID, models.MTOShipmentTypeHHG)
		suite.NotNil(err)
		suite.Nil(RDD)
		suite.Equal(fmt.Sprintf("unable to calculate domestic transit time due to missing weight for move ID: %s", move.ID), err.Error())
	})

	suite.Run("AK -> AK - successfully calculates RDD for intra-Alaska shipments", func() {
		pickupDate := time.Now()
		IntraAlaskaUBTransitTime := 30
		AlaskaZone1ToZone4HHGTransitTime := 15
		AlaskaZone2ToZone4HHGTransitTime := 60

		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		zone1Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "100 Mountain Hemlock St",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99504",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		zone2Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "200 Sitka Spruce St",
					City:           "FORT WAINWRIGHT",
					State:          "AK",
					PostalCode:     "99703",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		zone4Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "400 Tamarack St",
					City:           "Cordova",
					State:          "AK",
					PostalCode:     "99677",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		// Zone 1 -> Zone 4 UB
		RDD, err := CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, zone1Address, zone4Address, time.Now(), models.IntPointer(500), move.ID, models.MTOShipmentTypeUnaccompaniedBaggage)
		suite.Nil(err)
		suite.NotNil(RDD)
		suite.Equal(pickupDate.AddDate(0, 0, IntraAlaskaUBTransitTime).Day(), RDD.Day())

		// Zone 1 -> Zone 4 HHG
		RDD, err = CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, zone1Address, zone4Address, time.Now(), models.IntPointer(500), move.ID, models.MTOShipmentTypeHHG)
		suite.Nil(err)
		suite.NotNil(RDD)
		suite.Equal(pickupDate.AddDate(0, 0, AlaskaZone1ToZone4HHGTransitTime).Day(), RDD.Day())

		// Zone 1 -> Zone 4 NTS-Release
		RDD, err = CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, zone1Address, zone4Address, time.Now(), models.IntPointer(500), move.ID, models.MTOShipmentTypeHHGOutOfNTS)
		suite.Nil(err)
		suite.NotNil(RDD)
		suite.Equal(pickupDate.AddDate(0, 0, AlaskaZone1ToZone4HHGTransitTime).Day(), RDD.Day())

		// Zone 2 -> Zone 4 UB - Also tests that a nil weight is accepted for UBs
		RDD, err = CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, zone2Address, zone4Address, time.Now(), nil, move.ID, models.MTOShipmentTypeUnaccompaniedBaggage)
		suite.Nil(err)
		suite.NotNil(RDD)
		suite.Equal(pickupDate.AddDate(0, 0, IntraAlaskaUBTransitTime).Day(), RDD.Day())

		// Zone 2 -> Zone 4 HHG - Also tests that a nil weight is accepted for Intra-Alaska HHG
		RDD, err = CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, zone2Address, zone4Address, time.Now(), nil, move.ID, models.MTOShipmentTypeHHG)
		suite.Nil(err)
		suite.NotNil(RDD)
		suite.Equal(pickupDate.AddDate(0, 0, AlaskaZone2ToZone4HHGTransitTime).Day(), RDD.Day())
	})

	suite.Run("AK -> AK - correctly errors when rate area for pickup address is not found", func() {
		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now().AddDate(0, 0, -1),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		unsavedZone1AddressId := uuid.Must(uuid.NewV4())
		unsavedZone1Address := models.Address{
			ID:             unsavedZone1AddressId,
			StreetAddress1: "101 Mountain Hemlock St",
			City:           "JB Elmendorf-Richardson",
			State:          "AK",
			PostalCode:     "99506",
			IsOconus:       models.BoolPointer(true),
		}

		savedZone2Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "200 Sitka Spruce St",
					City:           "Eielson AFB",
					State:          "AK",
					PostalCode:     "99702",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		RDD, err := CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, unsavedZone1Address, savedZone2Address, time.Now(), models.IntPointer(500), move.ID, models.MTOShipmentTypeUnaccompaniedBaggage)
		suite.NotNil(err)
		suite.Nil(RDD)
		suite.Equal(fmt.Sprintf("error fetching pickup rate area id for address ID: %s", unsavedZone1Address.ID), err.Error())
	})

	suite.Run("AK -> AK - correctly errors when rate area for destination address is not found", func() {
		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now().AddDate(0, 0, -1),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		unsavedZone1AddressId := uuid.Must(uuid.NewV4())
		unsavedZone1Address := models.Address{
			ID:             unsavedZone1AddressId,
			StreetAddress1: "101 Mountain Hemlock St",
			City:           "JB Elmendorf-Richardson",
			State:          "AK",
			PostalCode:     "99506",
			IsOconus:       models.BoolPointer(true),
		}

		savedZone2Address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "200 Sitka Spruce St",
					City:           "Eielson AFB",
					State:          "AK",
					PostalCode:     "99702",
					IsOconus:       models.BoolPointer(true),
				},
			}}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		RDD, err := CalculateRequiredDeliveryDate(suite.AppContextForTest(), planner, savedZone2Address, unsavedZone1Address, time.Now(), models.IntPointer(500), move.ID, models.MTOShipmentTypeUnaccompaniedBaggage)
		suite.NotNil(err)
		suite.Nil(RDD)
		suite.Equal(fmt.Sprintf("error fetching destination rate area id for address ID: %s", unsavedZone1Address.ID), err.Error())
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateShipmentBasicServiceItemEstimatePricing() {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	setupOconusToConusNtsShipment := func(estimatedWeight *unit.Pound) (models.StorageFacility, models.Address, models.Address, models.MTOShipment) {
		storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.StorageFacility{
					FacilityName: *models.StringPointer("Test Storage Name"),
					Email:        models.StringPointer("old@email.com"),
					LotNumber:    models.StringPointer("Test lot number"),
					Phone:        models.StringPointer("555-555-5555"),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99507",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "148 S East St",
					City:           "Miami",
					State:          "FL",
					PostalCode:     "94535",
				},
			},
		}, nil)
		pickupDate := now.AddDate(0, 0, 10)
		shipment := factory.BuildNTSShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:             models.MoveStatusAPPROVED,
					AvailableToPrimeAt: &now,
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
			{
				Model: models.MTOShipment{
					MarketCode:           models.MarketCodeInternational,
					Status:               models.MTOShipmentStatusSubmitted,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTS,
					PrimeEstimatedWeight: estimatedWeight,
					ScheduledPickupDate:  &pickupDate,
					RequestedPickupDate:  &tomorrow,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)
		return storageFacility, pickupAddress, destinationAddress, shipment
	}

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(1000, nil)
	mockSender := setUpMockNotificationSender()
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	waf := entitlements.NewWeightAllotmentFetcher()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()
	mtoShipmentUpdaterPrime := NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

	// Start as nil when approved, then update with actual values once an estimated weight is set
	suite.Run("NTS OCONUS to CONUS INPK estimate prices on update", func() {
		var subtestData approveShipmentSubtestData

		subtestData.appCtx = suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		appCtx := subtestData.appCtx
		subtestData.planner = &mocks.Planner{}
		planner := subtestData.planner
		router := NewShipmentRouter()
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		subtestData.shipmentApprover = NewShipmentApprover(router, siCreator, subtestData.planner, subtestData.moveWeights, subtestData.mtoUpdater, moveRouter)
		approver := subtestData.shipmentApprover
		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		// Setup shipment with no estimated weight
		_, _, _, shipment := setupOconusToConusNtsShipment(nil)
		contract, err := models.FetchContractForMove(suite.AppContextForTest(), shipment.MoveTaskOrderID)
		suite.FatalNoError(err)

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(500, nil)

		// Approve the shipment to trigger the estimate pricing proc on INPK
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		_, approverErr := approver.ApproveShipment(appCtx, shipment.ID, shipmentEtag)
		suite.FatalNoError(approverErr)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)

		// Assert basic service items nil
		// Basic = created immediately on their own, not requested
		// Accessorial = created at a later date if requested
		expectedServiceItems := map[models.ReServiceCode]*unit.Cents{
			models.ReServiceCodeISLH:   nil,
			models.ReServiceCodePODFSC: nil,
			models.ReServiceCodeINPK:   nil,
		}
		suite.Equal(len(expectedServiceItems), len(serviceItems))

		// Look for INPK and assert its expected price matches the actual price the proc sets
		var foundINPKAfterApproval bool
		for _, serviceItem := range serviceItems {
			actualReServiceCode := serviceItem.ReService.Code
			suite.Contains(expectedServiceItems, actualReServiceCode, "Unexpected service code found: %s", actualReServiceCode)

			expectedPrice, found := expectedServiceItems[actualReServiceCode]
			suite.True(found, "Expected price for service code %s not found", actualReServiceCode)
			if actualReServiceCode == models.ReServiceCodeINPK {
				foundINPKAfterApproval = true
				suite.Nil(expectedPrice, "Expected price should be nil for service code %s", actualReServiceCode)
				suite.Nil(serviceItem.PricingEstimate, "Pricing estimate should be nil for service code %s", actualReServiceCode)
			}
		}
		suite.FatalTrue(foundINPKAfterApproval)

		// Now let's update it so it gets a non-nil value
		// Get the eTAG
		var oldUpdatedShipment = models.MTOShipment{
			ID: shipment.ID,
		}
		err = suite.DB().First(&oldUpdatedShipment)
		suite.FatalNoError(err)
		eTag := etag.GenerateEtag(oldUpdatedShipment.UpdatedAt)

		session := auth.Session{}
		primeEstimatedWeight := unit.Pound(1234)
		newUpdatedShipment := models.MTOShipment{
			ID:                   shipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		_, err = mtoShipmentUpdaterPrime.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &newUpdatedShipment, eTag, "test")
		suite.FatalNoError(err)

		// Get the contract escalation factor
		var escalationFactor float64
		err = suite.DB().RawQuery(`
			SELECT calculate_escalation_factor($1, $2)
		`, contract.ID, shipment.RequestedPickupDate).First(&escalationFactor)
		suite.FatalNoError(err)

		// Verify our non-truncated escalation factor db value is as expected
		// this also tests the calculate_escalation_factor proc
		// This information was pulled from the migration scripts (Or just run db fresh and perform the lookups
		// manually, whichever is your cup of tea)
		suite.Equal(escalationFactor, 1.11)

		// Fetch the INPK market factor from the DB
		inpkReService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeINPK)
		ntsMarketFactor, err := models.FetchMarketFactor(suite.AppContextForTest(), contract.ID, inpkReService.ID, "O")
		suite.FatalNoError(err)

		// Assert basic service items
		// Basic = created immediately on their own, not requested
		// Accessorial = created at a later date if requested
		expectedServiceItems = map[models.ReServiceCode]*unit.Cents{
			// Not testing ISLH or PODFSC so their prices will be nil
			models.ReServiceCodeISLH:   nil,
			models.ReServiceCodePODFSC: nil,
			// Remember that we pass in IHPK base price, not INPK base price. INPK doesn't have a base price
			// because it uses IHPK for iHHG -> iNTS packing
			models.ReServiceCodeINPK: func() *unit.Cents {
				ihpkService, err := models.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIHPK)
				suite.FatalNoError(err)

				ihpkRIOP, err := models.FetchReIntlOtherPrice(suite.DB(), *shipment.PickupAddressID, ihpkService.ID, contract.ID, shipment.RequestedPickupDate)
				suite.FatalNoError(err)
				suite.NotEmpty(ihpkRIOP)

				return models.CentPointer(computeINPKExpectedPriceCents(ihpkRIOP.PerUnitCents.Int(), escalationFactor, ntsMarketFactor, primeEstimatedWeight.Int()))
			}(),
		}
		// Get updated service items
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.FatalNoError(err)
		suite.Equal(len(expectedServiceItems), len(serviceItems))

		// Look for INPK and assert its expected price matches the actual price the proc sets
		var foundINPK bool
		for _, serviceItem := range serviceItems {
			actualReServiceCode := serviceItem.ReService.Code
			suite.Contains(expectedServiceItems, actualReServiceCode, "Unexpected service code found: %s", actualReServiceCode)

			expectedPrice, found := expectedServiceItems[actualReServiceCode]
			suite.True(found, "Expected price for service code %s not found", actualReServiceCode)
			if actualReServiceCode == models.ReServiceCodeINPK {
				foundINPK = true
				suite.FatalNotNil(expectedPrice)
				suite.FatalNotNil(serviceItem.PricingEstimate)
				suite.Equal(*expectedPrice, *serviceItem.PricingEstimate, "Pricing estimate mismatch for service code %s", actualReServiceCode)

			}
		}
		suite.FatalTrue(foundINPK)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateRequestedPickupDate() {

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	waf := entitlements.NewWeightAllotmentFetcher()
	mockSender := setUpMockNotificationSender()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()

	createSubtestData := func() (services.MTOShipmentUpdater, models.Move) {
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(500, nil)
		mtoShipmentUpdaterPrime := NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator, addressUpdater, addressCreator)

		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().AddDate(1, 0, 0),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		return mtoShipmentUpdaterPrime, move
	}

	suite.Run("RequestedPickupDate validation check - must be in the future for shipment types other than PPM", func() {
		shipmentUpdater, move := createSubtestData()

		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		tomorrow := now.AddDate(0, 0, 1)

		testCases := []struct {
			input        *time.Time
			shipmentType models.MTOShipmentType
			shouldError  bool
		}{
			// HHG
			{nil, models.MTOShipmentTypeHHG, false},
			{&time.Time{}, models.MTOShipmentTypeHHG, false},
			{&yesterday, models.MTOShipmentTypeHHG, true},
			{&now, models.MTOShipmentTypeHHG, true},
			{&tomorrow, models.MTOShipmentTypeHHG, false},
			// NTS
			{nil, models.MTOShipmentTypeHHGIntoNTS, false},
			{&time.Time{}, models.MTOShipmentTypeHHGIntoNTS, false},
			{&yesterday, models.MTOShipmentTypeHHGIntoNTS, true},
			{&now, models.MTOShipmentTypeHHGIntoNTS, true},
			{&tomorrow, models.MTOShipmentTypeHHGIntoNTS, false},
			// NTSR
			{nil, models.MTOShipmentTypeHHGOutOfNTS, false},
			{&time.Time{}, models.MTOShipmentTypeHHGOutOfNTS, false},
			{&yesterday, models.MTOShipmentTypeHHGOutOfNTS, true},
			{&now, models.MTOShipmentTypeHHGOutOfNTS, true},
			{&tomorrow, models.MTOShipmentTypeHHGOutOfNTS, false},
			// BOAT HAUL AWAY
			{nil, models.MTOShipmentTypeBoatHaulAway, false},
			{&time.Time{}, models.MTOShipmentTypeBoatHaulAway, false},
			{&yesterday, models.MTOShipmentTypeBoatHaulAway, true},
			{&now, models.MTOShipmentTypeBoatHaulAway, true},
			{&tomorrow, models.MTOShipmentTypeBoatHaulAway, false},
			// BOAT TOW AWAY
			{nil, models.MTOShipmentTypeBoatTowAway, false},
			{&time.Time{}, models.MTOShipmentTypeBoatTowAway, false},
			{&yesterday, models.MTOShipmentTypeBoatTowAway, true},
			{&now, models.MTOShipmentTypeBoatTowAway, true},
			{&tomorrow, models.MTOShipmentTypeBoatTowAway, false},
			// MOBILE HOME
			{nil, models.MTOShipmentTypeMobileHome, false},
			{&time.Time{}, models.MTOShipmentTypeMobileHome, false},
			{&yesterday, models.MTOShipmentTypeMobileHome, true},
			{&now, models.MTOShipmentTypeMobileHome, true},
			{&tomorrow, models.MTOShipmentTypeMobileHome, false},
			// UB
			{nil, models.MTOShipmentTypeUnaccompaniedBaggage, false},
			{&time.Time{}, models.MTOShipmentTypeUnaccompaniedBaggage, false},
			{&yesterday, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&now, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&tomorrow, models.MTOShipmentTypeUnaccompaniedBaggage, false},
			// PPM - should always pass validation
			{nil, models.MTOShipmentTypePPM, false},
			{&time.Time{}, models.MTOShipmentTypePPM, false},
			{&yesterday, models.MTOShipmentTypePPM, false},
			{&now, models.MTOShipmentTypePPM, false},
			{&tomorrow, models.MTOShipmentTypePPM, false},
		}

		for _, testCase := range testCases {
			// Default is HHG, but we set it explicitly below via the test cases
			var oldShipment models.MTOShipment
			if testCase.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
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
				moveForPrime := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
				oldShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    moveForPrime,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:         testCase.shipmentType,
							RequestedPickupDate:  &tomorrow,
							ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
							PrimeEstimatedWeight: models.PoundPointer(unit.Pound(4000)),
							Status:               models.MTOShipmentStatusSubmitted,
						},
					},
				}, nil)
			} else {
				oldShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
					{
						Model:    move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:        testCase.shipmentType,
							Status:              models.MTOShipmentStatusSubmitted,
							RequestedPickupDate: &tomorrow,
						},
					},
				}, nil)
			}

			updatedShipment := models.MTOShipment{
				ID:                  oldShipment.ID,
				RequestedPickupDate: testCase.input,
			}

			eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
			too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
			session := auth.Session{
				ApplicationName: auth.OfficeApp,
				UserID:          *too.UserID,
				OfficeUserID:    too.ID,
				ActiveRole:      too.User.Roles[0],
			}

			shipment, err := shipmentUpdater.UpdateMTOShipment(suite.AppContextWithSessionForTest(&session), &updatedShipment, eTag, "test")

			testCaseInputString := ""
			if testCase.input == nil {
				testCaseInputString = "nil"
			} else {
				testCaseInputString = (*testCase.input).String()
			}

			if testCase.shouldError {
				suite.Nil(shipment, "Should error for %s | %s", testCase.shipmentType, testCaseInputString)
				suite.Error(err)
				if testCase.input != nil && !(*testCase.input).IsZero() {
					suite.Equal("RequestedPickupDate must be greater than or equal to tomorrow's date.", err.Error())
				} else {
					suite.Contains(err.Error(), fmt.Sprintf("RequestedPickupDate is required to create or modify %s %s shipment", GetAorAnByShipmentType(testCase.shipmentType), testCase.shipmentType))
				}
			} else {
				suite.NoError(err, "Should not error for %s | %s", testCase.shipmentType, testCaseInputString)
				suite.NotNil(shipment)
			}
		}
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateSITServiceItemsSITIfPostalCodeChanged() {

	setupData := func(isPickupAddressTest bool, isOConus bool) (models.MTOShipment, models.Address, models.Address) {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		isPickupAddressOconus := false
		isDestinatonaAddressOconus := false

		if isPickupAddressTest {
			isPickupAddressOconus = isOConus
		} else {
			isDestinatonaAddressOconus = isOConus
		}

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(isPickupAddressOconus),
				},
			},
		}, nil)

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER1",
					City:           "Anchorage1",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(isDestinatonaAddressOconus),
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
					ShipmentType:       models.MTOShipmentTypeHHG,
					UsesExternalVendor: true,
					Status:             models.MTOShipmentStatusApproved,
					MarketCode:         models.MarketCodeInternational,
				},
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
		}, nil)

		customization := make([]factory.Customization, 0)
		customization = append(customization,
			factory.Customization{
				Model:    move,
				LinkOnly: true,
			},
			factory.Customization{
				Model:    shipment,
				LinkOnly: true,
			},
			factory.Customization{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			})
		if isPickupAddressTest {
			customization = append(customization,
				factory.Customization{
					Model: models.ReService{
						Code: models.ReServiceCodeIOSFSC,
					},
				},
				factory.Customization{
					Model:    pickupAddress,
					Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
					LinkOnly: true,
				},
				factory.Customization{
					Model:    pickupAddress,
					Type:     &factory.Addresses.SITOriginHHGActualAddress,
					LinkOnly: true,
				},
			)
		} else {
			customization = append(customization,
				factory.Customization{
					Model: models.ReService{
						Code: models.ReServiceCodeIDSFSC,
					},
				},
				factory.Customization{
					Model:    destinationAddress,
					Type:     &factory.Addresses.SITDestinationOriginalAddress,
					LinkOnly: true,
				},
				factory.Customization{
					Model:    destinationAddress,
					Type:     &factory.Addresses.SITDestinationFinalAddress,
					LinkOnly: true,
				})
		}

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), customization, nil)

		shipment.MTOServiceItems = append(shipment.MTOServiceItems, serviceItem)

		return shipment, pickupAddress, destinationAddress
	}

	suite.Run("IOSFSC - success", func() {
		shipment, pickupAddress, _ := setupData(true, false)

		expectedMileage := 23
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			pickupAddress.PostalCode,
			pickupAddress.PostalCode,
			mock.Anything,
		).Return(expectedMileage, nil)

		var serviceItems []models.MTOServiceItem
		err := suite.AppContextForTest().DB().EagerPreload("ReService", "SITOriginHHGOriginalAddress", "SITOriginHHGActualAddress").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].SITDeliveryMiles == (*int)(nil))
			suite.Equal(serviceItems[i].SITOriginHHGOriginalAddress.PostalCode, pickupAddress.PostalCode)
			suite.Equal(serviceItems[i].SITOriginHHGActualAddress.PostalCode, pickupAddress.PostalCode)
		}

		addressCreator := address.NewAddressCreator()
		err = UpdateSITServiceItemsSITIfPostalCodeChanged(planner, suite.AppContextForTest(), addressCreator, &shipment)
		suite.Nil(err)

		err = suite.AppContextForTest().DB().EagerPreload("ReService", "SITOriginHHGOriginalAddress", "SITOriginHHGActualAddress").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOSFSC)
			suite.Equal(*serviceItems[i].SITDeliveryMiles, expectedMileage)
			suite.Equal(serviceItems[i].SITOriginHHGOriginalAddress.PostalCode, pickupAddress.PostalCode)
			suite.Equal(serviceItems[i].SITOriginHHGActualAddress.PostalCode, pickupAddress.PostalCode)
		}
	})

	suite.Run("IDSFSC - success", func() {
		shipment, _, destinationAddress := setupData(false, false)

		expectedMileage := 23
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			destinationAddress.PostalCode,
			destinationAddress.PostalCode,
			mock.Anything,
		).Return(expectedMileage, nil)

		var serviceItems []models.MTOServiceItem
		err := suite.AppContextForTest().DB().EagerPreload("ReService", "SITDestinationOriginalAddress", "SITDestinationFinalAddress").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].SITDeliveryMiles == (*int)(nil))
			suite.Equal(serviceItems[i].SITDestinationOriginalAddress.PostalCode, destinationAddress.PostalCode)
			suite.Equal(serviceItems[i].SITDestinationFinalAddress.PostalCode, destinationAddress.PostalCode)
		}

		addressCreator := address.NewAddressCreator()
		err = UpdateSITServiceItemsSITIfPostalCodeChanged(planner, suite.AppContextForTest(), addressCreator, &shipment)
		suite.Nil(err)

		err = suite.AppContextForTest().DB().EagerPreload("ReService", "SITDestinationOriginalAddress", "SITDestinationFinalAddress").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIDSFSC)
			suite.Equal(*serviceItems[i].SITDeliveryMiles, expectedMileage)
			suite.Equal(serviceItems[i].SITDestinationOriginalAddress.PostalCode, destinationAddress.PostalCode)
			suite.Equal(serviceItems[i].SITDestinationFinalAddress.PostalCode, destinationAddress.PostalCode)
		}
	})
}
