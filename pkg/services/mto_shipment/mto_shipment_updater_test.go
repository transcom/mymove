//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package mtoshipment

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	notificationMocks "github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/fetch"

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
	oldMTOShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(500, nil)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(NewShipmentReweighRequester())

	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
	mockSender := setUpMockNotificationSender()
	mtoShipmentUpdater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator)

	requestedPickupDate := *oldMTOShipment.RequestedPickupDate
	scheduledPickupDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	firstAvailableDeliveryDate := time.Date(2019, time.March, 10, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(2020, time.June, 8, 0, 0, 0, 0, time.UTC)
	secondaryPickupAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
	secondaryDeliveryAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})
	primeActualWeight := unit.Pound(1234)
	primeEstimatedWeight := unit.Pound(1234)
	newDestinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "987 Other Avenue",
			StreetAddress2: swag.String("P.O. Box 1234"),
			StreetAddress3: swag.String("c/o Another Person"),
			City:           "Des Moines",
			State:          "IA",
			PostalCode:     "50309",
			Country:        swag.String("US"),
		},
	})

	newPickupAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "987 Over There Avenue",
			StreetAddress2: swag.String("P.O. Box 1234"),
			StreetAddress3: swag.String("c/o Another Person"),
			City:           "Houston",
			State:          "TX",
			PostalCode:     "77083",
			Country:        swag.String("US"),
		},
	})

	mtoShipment := models.MTOShipment{
		ID:                         oldMTOShipment.ID,
		MoveTaskOrderID:            oldMTOShipment.MoveTaskOrderID,
		DestinationAddress:         oldMTOShipment.DestinationAddress,
		DestinationAddressID:       oldMTOShipment.DestinationAddressID,
		PickupAddress:              oldMTOShipment.PickupAddress,
		PickupAddressID:            oldMTOShipment.PickupAddressID,
		RequestedPickupDate:        &requestedPickupDate,
		ScheduledPickupDate:        &scheduledPickupDate,
		ShipmentType:               "INTERNATIONAL_UB",
		PrimeActualWeight:          &primeActualWeight,
		PrimeEstimatedWeight:       &primeEstimatedWeight,
		FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
		Status:                     oldMTOShipment.Status,
		ActualPickupDate:           &actualPickupDate,
		ApprovedDate:               &firstAvailableDeliveryDate,
	}

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	//now := time.Now()
	primeEstimatedWeight = unit.Pound(4500)

	suite.T().Run("Can retrieve existing shipment", func(t *testing.T) {

		existingShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})

		var reServiceDomCrating models.ReService
		if err := suite.DB().Where("code = $1", models.ReServiceCodeDCRT).First(&reServiceDomCrating); err != nil {
			// Something is truncating this when all server tests run, but we need this ReService value to exist
			reServiceDomCrating = testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code:      models.ReServiceCodeDCRT,
					Name:      "test",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			})
		}

		mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				MoveTaskOrderID: existingShipment.MoveTaskOrderID,
				MTOShipmentID:   &existingShipment.ID,
			},
			ReService: reServiceDomCrating,
		})

		item := testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
			MTOServiceItemDimension: models.MTOServiceItemDimension{
				Type:      models.DimensionTypeItem,
				Length:    1000,
				Height:    1000,
				Width:     1000,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			MTOServiceItem: mtoServiceItem1,
		})

		crate := testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
			MTOServiceItemDimension: models.MTOServiceItemDimension{
				MTOServiceItemID: mtoServiceItem1.ID,
				Type:             models.DimensionTypeCrate,
				Length:           2000,
				Height:           2000,
				Width:            2000,
				CreatedAt:        time.Time{},
				UpdatedAt:        time.Time{},
			},
		})

		shipment, err := mtoShipmentUpdater.RetrieveMTOShipment(suite.AppContextForTest(), existingShipment.ID)

		suite.NoError(err)

		suite.Equal(existingShipment.ID, shipment.ID)
		suite.Equal(existingShipment.CreatedAt.UTC(), shipment.CreatedAt.UTC())
		suite.Equal(existingShipment.ShipmentType, shipment.ShipmentType)
		suite.Equal(existingShipment.UpdatedAt.UTC(), shipment.UpdatedAt.UTC())

		suite.Require().Equal(1, len(shipment.MTOServiceItems))
		suite.Require().Equal(2, len(shipment.MTOServiceItems[0].Dimensions))
		for _, s := range shipment.MTOServiceItems[0].Dimensions {
			if s.Type == models.DimensionTypeCrate {
				suite.Equal(crate.Height, s.Height)
			} else {
				suite.Equal(item.Height, s.Height)
			}
		}

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	servicesCounselor := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{})

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *servicesCounselor.UserID,
		OfficeUserID:    servicesCounselor.ID,
	}
	session.Roles = append(session.Roles, servicesCounselor.User.Roles...)

	var statusTests = []struct {
		name      string
		status    models.MTOShipmentStatus
		updatable bool
	}{
		{"Draft isn't updatable", models.MTOShipmentStatusDraft, false},
		{"Submitted is updatable", models.MTOShipmentStatusSubmitted, true},
		{"Approved isn't updatable", models.MTOShipmentStatusApproved, false},
	}

	for _, tt := range statusTests {
		suite.T().Run(fmt.Sprintf("Updatable status returned as expected: %v", tt.name), func(t *testing.T) {
			shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					Status: tt.status,
				},
			})

			updatable, err := mtoShipmentUpdater.CheckIfMTOShipmentCanBeUpdated(suite.AppContextForTest(), &shipment, &session)

			suite.NoError(err)

			suite.Equal(tt.updatable, updatable,
				"Expected updatable to be %v when status is %v. Got %v", tt.updatable, tt.status, updatable)

			// Verify that shipment recalculate was handled correctly
			mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
		})
	}

	suite.T().Run("Etag is stale", func(t *testing.T) {
		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), &mtoShipment, eTag)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.T().Run("If-Unmodified-Since is equal to the updated_at date", func(t *testing.T) {
		eTag := etag.GenerateEtag(oldMTOShipment.UpdatedAt)
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), &mtoShipment, eTag)

		suite.Require().NoError(err)
		suite.Equal(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)

		suite.Equal(updatedMTOShipment.PickupAddressID, oldMTOShipment.PickupAddressID)

		suite.Equal(updatedMTOShipment.PrimeActualWeight, &primeActualWeight)
		suite.True(actualPickupDate.Equal(*updatedMTOShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*updatedMTOShipment.FirstAvailableDeliveryDate))

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	oldMTOShipment2 := testdatagen.MakeDefaultMTOShipment(suite.DB())
	mtoShipment2 := models.MTOShipment{
		ID:           oldMTOShipment2.ID,
		ShipmentType: "INTERNATIONAL_UB",
	}

	suite.T().Run("Updater can handle optional queries set as nil", func(t *testing.T) {
		eTag := etag.GenerateEtag(oldMTOShipment2.UpdatedAt)
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), &mtoShipment2, eTag)

		suite.Require().NoError(err)
		suite.Equal(updatedMTOShipment.ID, oldMTOShipment2.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment2.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	suite.T().Run("Successful update to all address fields", func(t *testing.T) {
		// Ensure we can update every address field on the shipment
		// Create an mtoShipment to update that has every address populated
		oldMTOShipment3 := testdatagen.MakeDefaultMTOShipment(suite.DB())

		eTag := etag.GenerateEtag(oldMTOShipment3.UpdatedAt)

		updatedShipment := &models.MTOShipment{
			ID:                         oldMTOShipment3.ID,
			DestinationAddress:         &newDestinationAddress,
			DestinationAddressID:       &newDestinationAddress.ID,
			PickupAddress:              &newPickupAddress,
			PickupAddressID:            &newPickupAddress.ID,
			SecondaryPickupAddress:     &secondaryPickupAddress,
			SecondaryPickupAddressID:   &secondaryDeliveryAddress.ID,
			SecondaryDeliveryAddress:   &secondaryDeliveryAddress,
			SecondaryDeliveryAddressID: &secondaryDeliveryAddress.ID,
		}

		updatedShipment, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), updatedShipment, eTag)

		suite.Require().NoError(err)
		suite.Equal(newDestinationAddress.ID, *updatedShipment.DestinationAddressID)
		suite.Equal(newDestinationAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(newPickupAddress.ID, *updatedShipment.PickupAddressID)
		suite.Equal(newPickupAddress.StreetAddress1, updatedShipment.PickupAddress.StreetAddress1)
		suite.Equal(secondaryPickupAddress.ID, *updatedShipment.SecondaryPickupAddressID)
		suite.Equal(secondaryPickupAddress.StreetAddress1, updatedShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(secondaryDeliveryAddress.ID, *updatedShipment.SecondaryDeliveryAddressID)
		suite.Equal(secondaryDeliveryAddress.StreetAddress1, updatedShipment.SecondaryDeliveryAddress.StreetAddress1)
		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))

	})

	suite.T().Run("Successful update to a minimal MTO shipment", func(t *testing.T) {
		// Minimal MTO Shipment has no associated addresses created by default.
		// Part of this test ensures that if an address doesn't exist on a shipment,
		// the updater can successfully create it.
		oldShipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := time.Date(2019, time.March, 15, 0, 0, 0, 0, time.UTC)
		scheduledPickupDate := time.Date(2019, time.March, 17, 0, 0, 0, 0, time.UTC)
		requestedDeliveryDate := time.Date(2019, time.March, 30, 0, 0, 0, 0, time.UTC)
		primeEstimatedWeightRecordedDate := time.Date(2019, time.March, 12, 0, 0, 0, 0, time.UTC)
		customerRemarks := "I have a grandfather clock"
		counselorRemarks := "Counselor approved"
		updatedShipment := models.MTOShipment{
			ID:                               oldShipment.ID,
			DestinationAddress:               &newDestinationAddress,
			DestinationAddressID:             &newDestinationAddress.ID,
			PickupAddress:                    &newPickupAddress,
			PickupAddressID:                  &newPickupAddress.ID,
			SecondaryPickupAddress:           &secondaryPickupAddress,
			SecondaryDeliveryAddress:         &secondaryDeliveryAddress,
			RequestedPickupDate:              &requestedPickupDate,
			ScheduledPickupDate:              &scheduledPickupDate,
			RequestedDeliveryDate:            &requestedDeliveryDate,
			ActualPickupDate:                 &actualPickupDate,
			PrimeActualWeight:                &primeActualWeight,
			PrimeEstimatedWeight:             &primeEstimatedWeight,
			FirstAvailableDeliveryDate:       &firstAvailableDeliveryDate,
			PrimeEstimatedWeightRecordedDate: &primeEstimatedWeightRecordedDate,
			Status:                           models.MTOShipmentStatusSubmitted,
			CustomerRemarks:                  &customerRemarks,
			CounselorRemarks:                 &counselorRemarks,
		}

		newShipment, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), &updatedShipment, eTag)

		suite.Require().NoError(err)
		suite.True(requestedPickupDate.Equal(*newShipment.RequestedPickupDate))
		suite.True(scheduledPickupDate.Equal(*newShipment.ScheduledPickupDate))
		suite.True(requestedDeliveryDate.Equal(*newShipment.RequestedDeliveryDate))
		suite.True(actualPickupDate.Equal(*newShipment.ActualPickupDate))
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

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Updating a shipment does not nullify ApprovedDate", func(t *testing.T) {
		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		oldShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		})

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := time.Date(2019, time.March, 15, 0, 0, 0, 0, time.UTC)
		requestedDeliveryDate := time.Date(2019, time.March, 30, 0, 0, 0, 0, time.UTC)
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
			RequestedPickupDate:      &requestedPickupDate,
			RequestedDeliveryDate:    &requestedDeliveryDate,
			CustomerRemarks:          &customerRemarks,
			CounselorRemarks:         &counselorRemarks,
		}

		newShipment, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), &updatedShipment, eTag)

		suite.Require().NoError(err)
		suite.NotEmpty(newShipment.ApprovedDate)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Successfully update MTO Agents", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		mtoAgent1 := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOAgent: models.MTOAgent{
				MTOShipment:   shipment,
				MTOShipmentID: shipment.ID,
				FirstName:     swag.String("Test"),
				LastName:      swag.String("Agent"),
				Email:         swag.String("test@test.email.com"),
				MTOAgentType:  models.MTOAgentReleasing,
			},
		})
		mtoAgent2 := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOAgent: models.MTOAgent{
				MTOShipment:   shipment,
				MTOShipmentID: shipment.ID,
				FirstName:     swag.String("Test2"),
				LastName:      swag.String("Agent2"),
				Email:         swag.String("test2@test.email.com"),
				MTOAgentType:  models.MTOAgentReceiving,
			},
		})
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

		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), &updatedShipment, eTag)

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(phone, *updatedMTOShipment.MTOAgents[0].Phone)
		suite.Equal(newFirstName, *updatedMTOShipment.MTOAgents[0].FirstName)
		suite.Equal(email, *updatedMTOShipment.MTOAgents[1].Email)
		suite.Equal(newLastName, *updatedMTOShipment.MTOAgents[1].LastName)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Successfully add new MTO Agent and edit another", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		existingAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOAgent: models.MTOAgent{
				MTOShipment:   shipment,
				MTOShipmentID: shipment.ID,
				FirstName:     swag.String("Test"),
				LastName:      swag.String("Agent"),
				Email:         swag.String("test@test.email.com"),
				MTOAgentType:  models.MTOAgentReleasing,
			},
		})

		mtoAgentToCreate := models.MTOAgent{
			MTOShipment:   shipment,
			MTOShipmentID: shipment.ID,
			FirstName:     swag.String("Ima"),
			LastName:      swag.String("Newagent"),
			Email:         swag.String("test2@test.email.com"),
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

		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), &updatedShipment, eTag)

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(phone, *updatedMTOShipment.MTOAgents[0].Phone)
		suite.Equal(*mtoAgentToCreate.FirstName, *updatedMTOShipment.MTOAgents[1].FirstName)
		suite.Equal(*mtoAgentToCreate.LastName, *updatedMTOShipment.MTOAgents[1].LastName)
		suite.Equal(*mtoAgentToCreate.Email, *updatedMTOShipment.MTOAgents[1].Email)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Successfully add storage facility to shipment", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		storageFacility := testdatagen.MakeStorageFacility(suite.DB(), testdatagen.Assertions{})

		updatedShipment := models.MTOShipment{
			ID:              shipment.ID,
			StorageFacility: &storageFacility,
		}
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipmentOffice(suite.AppContextForTest(), &updatedShipment, eTag)

		suite.Require().NoError(err)
		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.StorageFacility)
	})

	suite.T().Run("Successfully edit storage facility on shipment", func(t *testing.T) {
		email := "old@email.com"
		storageFacility := testdatagen.MakeStorageFacility(suite.DB(), testdatagen.Assertions{
			StorageFacility: models.StorageFacility{
				Address: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						StreetAddress1: "1234 Over Here Street",
						City:           "Houston",
						State:          "TX",
						PostalCode:     "77083",
						Country:        swag.String("US"),
					},
				}),
				Email: &email,
			},
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				StorageFacility: &storageFacility,
			},
		})

		newStorageFacilityAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "987 Over There Avenue",
				City:           "Houston",
				State:          "TX",
				PostalCode:     "77083",
				Country:        swag.String("US"),
			},
		})

		newEmail := "old@email.com"
		newStorageFacility := testdatagen.MakeStorageFacility(suite.DB(), testdatagen.Assertions{
			StorageFacility: models.StorageFacility{
				Address: newStorageFacilityAddress,
				Email:   &newEmail,
			},
		})

		newShipment := models.MTOShipment{
			ID:              shipment.ID,
			StorageFacility: &newStorageFacility,
		}

		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		updatedShipment, err := mtoShipmentUpdater.UpdateMTOShipmentOffice(suite.AppContextForTest(), &newShipment, eTag)
		suite.Require().NoError(err)
		suite.NotZero(newShipment.ID, oldMTOShipment.ID)
		suite.Equal(&newEmail, updatedShipment.StorageFacility.Email)
		suite.Equal(newStorageFacilityAddress.StreetAddress1, updatedShipment.StorageFacility.Address.StreetAddress1)
	})

	suite.T().Run("Successfully divert a shipment and transition statuses", func(t *testing.T) {
		// A diverted shipment should transition to the SUBMITTED status.
		// If the move it is connected to is APPROVED, that move should transition to APPROVALS REQUESTED
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				MoveTaskOrder: move,
				Status:        models.MTOShipmentStatusApproved,
				Diversion:     false,
			},
		})
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		shipmentInput := models.MTOShipment{
			ID:        shipment.ID,
			Diversion: true,
		}

		updatedShipment, err := mtoShipmentUpdater.UpdateMTOShipmentCustomer(suite.AppContextForTest(), &shipmentInput, eTag)

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
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"))
	})

	// Test UpdateMTOShipmentPrime
	// TODO: Add more tests, such as making sure this function fails if the
	// move is not available to the prime.
	suite.T().Run("Updating a shipment does not nullify ApprovedDate", func(t *testing.T) {
		// This test was added because of a bug that nullified the ApprovedDate
		// when ScheduledPickupDate was included in the payload. See PR #6919.
		// ApprovedDate affects shipment diversions, so we want to make sure it
		// never gets nullified, regardless of which fields are being updated.
		move := testdatagen.MakeAvailableMove(suite.DB())
		oldShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
			Move: move,
		})

		suite.NotNil(oldShipment.ApprovedDate)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		requestedPickupDate := time.Date(2019, time.March, 15, 0, 0, 0, 0, time.UTC)
		scheduledPickupDate := time.Date(2019, time.March, 17, 0, 0, 0, 0, time.UTC)
		requestedDeliveryDate := time.Date(2019, time.March, 30, 0, 0, 0, 0, time.UTC)
		updatedShipment := models.MTOShipment{
			ID:                         oldShipment.ID,
			DestinationAddress:         &newDestinationAddress,
			DestinationAddressID:       &newDestinationAddress.ID,
			PickupAddress:              &newPickupAddress,
			PickupAddressID:            &newPickupAddress.ID,
			SecondaryPickupAddress:     &secondaryPickupAddress,
			SecondaryDeliveryAddress:   &secondaryDeliveryAddress,
			RequestedPickupDate:        &requestedPickupDate,
			ScheduledPickupDate:        &scheduledPickupDate,
			RequestedDeliveryDate:      &requestedDeliveryDate,
			ActualPickupDate:           &actualPickupDate,
			PrimeActualWeight:          &primeActualWeight,
			PrimeEstimatedWeight:       &primeEstimatedWeight,
			FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
		}

		newShipment, err := mtoShipmentUpdater.UpdateMTOShipmentPrime(suite.AppContextForTest(), &updatedShipment, eTag)

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

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})
}

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentStatus() {
	mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{Status: models.MoveStatusAPPROVED}})
	estimatedWeight := unit.Pound(2000)
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
			PrimeEstimatedWeight: &estimatedWeight,
			Status:               models.MTOShipmentStatusSubmitted,
		},
	})
	draftShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusDraft,
		},
	})
	shipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	shipment3 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	shipment4 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	shipmentForAutoApprove := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
	})
	rejectionReason := "exotic animals are banned"
	rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status:          models.MTOShipmentStatusRejected,
			RejectionReason: &rejectionReason,
		},
	})
	shipment.Status = models.MTOShipmentStatusSubmitted
	eTag := etag.GenerateEtag(shipment.UpdatedAt)
	status := models.MTOShipmentStatusApproved

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	// Let's also create a transit time object with a zero upper bound for weight (this can happen in the table).
	ghcDomesticTransitTime0LbsUpper := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     10001,
		WeightLbsUpper:     0,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime0LbsUpper)

	builder := query.NewQueryBuilder()
	moveRouter := moveservices.NewMoveRouter()
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(500, nil)
	updater := NewMTOShipmentStatusUpdater(builder, siCreator, planner)

	suite.T().Run("If the mtoShipment is approved successfully it should create approved mtoServiceItems", func(t *testing.T) {
		appCtx := suite.AppContextForTest()
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		serviceItems := models.MTOServiceItems{}
		var expectedReServiceCodes []models.ReServiceCode
		expectedReServiceCodes = append(expectedReServiceCodes,
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
		)

		var reServiceCode models.ReService
		if err := appCtx.DB().Where("code = $1", expectedReServiceCodes[0]).First(&reServiceCode); err != nil {
			// Something is truncating these when all server tests run, but we need some values for reServices
			for _, serviceCode := range expectedReServiceCodes {
				testdatagen.MakeReService(appCtx.DB(), testdatagen.Assertions{
					ReService: models.ReService{
						Code:      serviceCode,
						Name:      "test",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				})
			}
		}

		preApprovalTime := time.Now()
		_, err := updater.UpdateMTOShipmentStatus(appCtx, shipmentForAutoApprove.ID, status, nil, shipmentForAutoApproveEtag)
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
		// service items created as a side effect to see if they are approved.
		missingReServiceCodes := expectedReServiceCodes
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

	suite.T().Run("If we act on a shipment with a weight that has a 0 upper weight it should still work", func(t *testing.T) {
		// This is testing that the Required Delivery Date is calculated correctly.
		// In order for the Required Delivery Date to be calculated, the following conditions must be true:
		// 1. The shipment is moving to the APPROVED status
		// 2. The shipment must already have the following fields present:
		// ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, DestinationAddress
		// 3. The shipment must not already have a Required Delivery Date
		// Note that MakeMTOShipment will automatically add a Required Delivery Date if the ScheduledPickupDate
		// is present, therefore we need to use MakeMTOShipmentMinimal and add the Pickup and Destination addresses
		estimatedWeight := unit.Pound(11000)
		destinationAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
		pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
		shipmentHeavy := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
				ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
				PrimeEstimatedWeight: &estimatedWeight,
				Status:               models.MTOShipmentStatusSubmitted,
				DestinationAddress:   &destinationAddress,
				DestinationAddressID: &destinationAddress.ID,
				PickupAddress:        &pickupAddress,
				PickupAddressID:      &pickupAddress.ID,
			},
		})
		shipmentHeavyEtag := etag.GenerateEtag(shipmentHeavy.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipmentHeavy.ID, status, nil, shipmentHeavyEtag)
		suite.NoError(err)
		serviceItems := models.MTOServiceItems{}
		_ = suite.DB().All(&serviceItems)
		fetchedShipment := models.MTOShipment{}
		err = suite.DB().Find(&fetchedShipment, shipmentHeavy.ID)
		suite.NoError(err)
		// We also should have a required delivery date
		suite.NotNil(fetchedShipment.RequiredDeliveryDate)
	})

	suite.T().Run("Cannot set SUBMITTED status on shipment via UpdateMTOShipmentStatus", func(t *testing.T) {
		// The only time a shipment gets set to the SUBMITTED status is when it is created, whether by the customer
		// or the Prime. This happens in the internal and prime API in the CreateMTOShipmentHandler. In that case,
		// the handlers will call ShipmentRouter.Submit().
		eTag = etag.GenerateEtag(draftShipment.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), draftShipment.ID, "SUBMITTED", nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)

		err = suite.DB().Find(&draftShipment, draftShipment.ID)

		suite.NoError(err)
		suite.EqualValues(models.MTOShipmentStatusDraft, draftShipment.Status)
	})

	suite.T().Run("Rejecting a shipment in SUBMITTED status with a rejection reason should return no error", func(t *testing.T) {
		eTag = etag.GenerateEtag(shipment2.UpdatedAt)
		rejectionReason := "Rejection reason"
		returnedShipment, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment2.ID, "REJECTED", &rejectionReason, eTag)

		suite.NoError(err)
		suite.NotNil(returnedShipment)

		err = suite.DB().Find(&shipment2, shipment2.ID)

		suite.NoError(err)
		suite.EqualValues(models.MTOShipmentStatusRejected, shipment2.Status)
		suite.Equal(&rejectionReason, shipment2.RejectionReason)
	})

	suite.T().Run("Rejecting a shipment with no rejection reason returns an InvalidInputError", func(t *testing.T) {
		eTag = etag.GenerateEtag(shipment3.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment3.ID, "REJECTED", nil, eTag)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.T().Run("Rejecting a shipment in APPROVED status returns a ConflictStatusError", func(t *testing.T) {
		eTag = etag.GenerateEtag(approvedShipment.UpdatedAt)
		rejectionReason := "Rejection reason"
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), approvedShipment.ID, "REJECTED", &rejectionReason, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.T().Run("Approving a shipment in REJECTED status returns a ConflictStatusError", func(t *testing.T) {
		eTag = etag.GenerateEtag(rejectedShipment.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), rejectedShipment.ID, "APPROVED", nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.T().Run("Passing in a stale identifier returns a PreconditionFailedError", func(t *testing.T) {
		staleETag := etag.GenerateEtag(time.Now())

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment4.ID, "APPROVED", nil, staleETag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.T().Run("Passing in an invalid status returns a ConflictStatus error", func(t *testing.T) {
		eTag = etag.GenerateEtag(shipment4.UpdatedAt)

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment4.ID, "invalid", nil, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.T().Run("Passing in a bad shipment id returns a Not Found error", func(t *testing.T) {
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), badShipmentID, "APPROVED", nil, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Changing to APPROVED status records approved_date", func(t *testing.T) {
		shipment5 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		eTag = etag.GenerateEtag(shipment5.UpdatedAt)

		suite.Nil(shipment5.ApprovedDate)

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment5.ID, models.MTOShipmentStatusApproved, nil, eTag)

		suite.NoError(err)
		suite.DB().Find(&shipment5, shipment5.ID)
		suite.Equal(models.MTOShipmentStatusApproved, shipment5.Status)
		suite.NotNil(shipment5.ApprovedDate)
	})

	suite.T().Run("Changing to a non-APPROVED status does not record approved_date", func(t *testing.T) {
		shipment6 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		eTag = etag.GenerateEtag(shipment6.UpdatedAt)
		rejectionReason := "reason"

		suite.Nil(shipment6.ApprovedDate)

		_, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), shipment6.ID, models.MTOShipmentStatusRejected, &rejectionReason, eTag)

		suite.NoError(err)
		suite.DB().Find(&shipment6, shipment6.ID)
		suite.Equal(models.MTOShipmentStatusRejected, shipment6.Status)
		suite.Nil(shipment6.ApprovedDate)
	})

	suite.T().Run("When move is not yet approved, cannot approve shipment", func(t *testing.T) {
		submittedMTO := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
		mtoShipment := submittedMTO.MTOShipments[0]
		eTag = etag.GenerateEtag(mtoShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(suite.AppContextForTest(), mtoShipment.ID, models.MTOShipmentStatusApproved, nil, eTag)
		suite.DB().Find(&mtoShipment, mtoShipment.ID)

		suite.Nil(updatedShipment)
		suite.Equal(models.MTOShipmentStatusSubmitted, mtoShipment.Status)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "Cannot approve a shipment if the move isn't approved.")
	})

	suite.T().Run("An approved shipment can change to CANCELLATION_REQUESTED", func(t *testing.T) {
		approvedShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		})
		eTag = etag.GenerateEtag(approvedShipment2.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), approvedShipment2.ID, models.MTOShipmentStatusCancellationRequested, nil, eTag)
		suite.DB().Find(&approvedShipment2, approvedShipment2.ID)

		suite.NoError(err)
		suite.NotNil(updatedShipment)
		suite.Equal(models.MTOShipmentStatusCancellationRequested, updatedShipment.Status)
		suite.Equal(models.MTOShipmentStatusCancellationRequested, approvedShipment2.Status)
	})

	suite.T().Run("A CANCELLATION_REQUESTED shipment can change to CANCELED", func(t *testing.T) {
		cancellationRequestedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusCancellationRequested,
			},
		})
		eTag = etag.GenerateEtag(cancellationRequestedShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), cancellationRequestedShipment.ID, models.MTOShipmentStatusCanceled, nil, eTag)
		suite.DB().Find(&cancellationRequestedShipment, cancellationRequestedShipment.ID)

		suite.NoError(err)
		suite.NotNil(updatedShipment)
		suite.Equal(models.MTOShipmentStatusCanceled, updatedShipment.Status)
		suite.Equal(models.MTOShipmentStatusCanceled, cancellationRequestedShipment.Status)
	})

	suite.T().Run("An APPROVED shipment CANNOT change to CANCELED - ERROR", func(t *testing.T) {
		eTag = etag.GenerateEtag(approvedShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), approvedShipment.ID, models.MTOShipmentStatusCanceled, nil, eTag)
		suite.DB().Find(&approvedShipment, approvedShipment.ID)

		suite.Error(err)
		suite.Nil(updatedShipment)
		suite.IsType(ConflictStatusError{}, err)
		suite.Equal(models.MTOShipmentStatusApproved, approvedShipment.Status)
	})

	suite.T().Run("An APPROVED shipment CAN change to Diversion Requested", func(t *testing.T) {
		shipmentToDivert := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		})
		eTag = etag.GenerateEtag(shipmentToDivert.UpdatedAt)

		_, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), shipmentToDivert.ID, models.MTOShipmentStatusDiversionRequested, nil, eTag)
		suite.DB().Find(&shipmentToDivert, shipmentToDivert.ID)

		suite.NoError(err)
		suite.Equal(models.MTOShipmentStatusDiversionRequested, shipmentToDivert.Status)
	})

	suite.T().Run("A diversion or diverted shipment can change to APPROVED", func(t *testing.T) {
		// a diversion or diverted shipment is when the PRIME sets the diversion field to true
		// the status must also be in diversion requested status to be approvable as well
		diversionRequestedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
			MTOShipment: models.MTOShipment{
				Status:    models.MTOShipmentStatusDiversionRequested,
				Diversion: true,
			},
		})
		eTag = etag.GenerateEtag(diversionRequestedShipment.UpdatedAt)

		updatedShipment, err := updater.UpdateMTOShipmentStatus(
			suite.AppContextForTest(), diversionRequestedShipment.ID, models.MTOShipmentStatusApproved, nil, eTag)

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
	primeShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
	})
	nonPrimeShipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
	hiddenPrimeShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
			Show:               &hide,
		},
	})

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
	updater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator)

	suite.T().Run("Shipment exists and is available to Prime - success", func(t *testing.T) {
		isAvailable, err := updater.MTOShipmentsMTOAvailableToPrime(suite.AppContextForTest(), primeShipment.ID)
		suite.True(isAvailable)
		suite.NoError(err)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Shipment exists but is not available to Prime - failure", func(t *testing.T) {
		isAvailable, err := updater.MTOShipmentsMTOAvailableToPrime(suite.AppContextForTest(), nonPrimeShipment.ID)
		suite.False(isAvailable)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), nonPrimeShipment.ID.String())

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Shipment exists, is available, but move is disabled - failure", func(t *testing.T) {
		isAvailable, err := updater.MTOShipmentsMTOAvailableToPrime(suite.AppContextForTest(), hiddenPrimeShipment.ID)
		suite.False(isAvailable)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), hiddenPrimeShipment.ID.String())

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Shipment does not exist - failure", func(t *testing.T) {
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		isAvailable, err := updater.MTOShipmentsMTOAvailableToPrime(suite.AppContextForTest(), badUUID)
		suite.False(isAvailable)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), badUUID.String())

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
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
	updater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator)

	suite.T().Run("Updating the shipment estimated weight will flag excess weight on the move and transitions move status", func(t *testing.T) {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)

		primeShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: models.Move{
				AvailableToPrimeAt: &now,
				Status:             models.MoveStatusAPPROVED,
			},
		})
		estimatedWeight := unit.Pound(7200)
		// there is a validator check about updating the status
		primeShipment.Status = ""
		primeShipment.PrimeEstimatedWeight = &estimatedWeight

		suite.Nil(primeShipment.MoveTaskOrder.ExcessWeightQualifiedAt)
		suite.Equal(models.MoveStatusAPPROVED, primeShipment.MoveTaskOrder.Status)

		_, err := updater.UpdateMTOShipmentPrime(suite.AppContextForTest(), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt))
		suite.NoError(err)

		err = suite.DB().Reload(&primeShipment.MoveTaskOrder)
		suite.NoError(err)

		suite.NotNil(primeShipment.MoveTaskOrder.ExcessWeightQualifiedAt)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, primeShipment.MoveTaskOrder.Status)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Skips calling check excess weight if estimated weight was not provided in request", func(t *testing.T) {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		mockedUpdater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		primeShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
		})
		// there is a validator check about updating the status
		primeShipment.Status = ""
		actualWeight := unit.Pound(7200)
		primeShipment.PrimeActualWeight = &actualWeight

		moveWeights.On("CheckAutoReweigh", mock.AnythingOfType("*appcontext.appContext"), primeShipment.MoveTaskOrderID, mock.AnythingOfType("*models.MTOShipment")).Return(models.MTOShipments{}, nil)

		suite.Nil(primeShipment.MoveTaskOrder.ExcessWeightQualifiedAt)

		_, err := mockedUpdater.UpdateMTOShipmentPrime(suite.AppContextForTest(), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt))
		suite.NoError(err)

		moveWeights.AssertNotCalled(t, "CheckExcessWeight")

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Skips calling check excess weight if the updated estimated weight matches the db value", func(t *testing.T) {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		mockedUpdater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(7200)
		primeShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:               models.MTOShipmentStatusApproved,
				ApprovedDate:         &now,
				ScheduledPickupDate:  &pickupDate,
				PrimeEstimatedWeight: &estimatedWeight,
			},
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
		})
		// there is a validator check about updating the status
		primeShipment.Status = ""
		primeShipment.PrimeEstimatedWeight = &estimatedWeight

		suite.Nil(primeShipment.MoveTaskOrder.ExcessWeightQualifiedAt)

		_, err := mockedUpdater.UpdateMTOShipmentPrime(suite.AppContextForTest(), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt))
		suite.NoError(err)

		moveWeights.AssertNotCalled(t, "CheckExcessWeight")

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
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
	updater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator)

	suite.T().Run("Updating the shipment actual weight within weight allowance creates reweigh requests for", func(t *testing.T) {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)

		primeShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: models.Move{
				AvailableToPrimeAt: &now,
				Status:             models.MoveStatusAPPROVED,
			},
		})
		actualWeight := unit.Pound(7200)
		// there is a validator check about updating the status
		primeShipment.Status = ""
		primeShipment.PrimeActualWeight = &actualWeight

		_, err := updater.UpdateMTOShipmentPrime(suite.AppContextForTest(), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt))
		suite.NoError(err)

		err = suite.DB().Eager("Reweigh").Reload(&primeShipment)
		suite.NoError(err)

		suite.NotNil(primeShipment.Reweigh)
		suite.Equal(primeShipment.ID.String(), primeShipment.Reweigh.ShipmentID.String())
		suite.NotNil(primeShipment.Reweigh.RequestedAt)
		suite.Equal(models.ReweighRequesterSystem, primeShipment.Reweigh.RequestedBy)

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Skips calling check auto reweigh if actual weight was not provided in request", func(t *testing.T) {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		mockedUpdater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		primeShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
		})
		// there is a validator check about updating the status
		primeShipment.Status = ""
		estimatedWeight := unit.Pound(7200)
		primeShipment.PrimeEstimatedWeight = &estimatedWeight

		moveWeights.On("CheckExcessWeight", mock.AnythingOfType("*appcontext.appContext"), primeShipment.MoveTaskOrderID, mock.AnythingOfType("models.MTOShipment")).Return(&primeShipment.MoveTaskOrder, nil, nil)

		_, err := mockedUpdater.UpdateMTOShipmentPrime(suite.AppContextForTest(), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt))
		suite.NoError(err)

		moveWeights.AssertNotCalled(t, "CheckAutoReweigh")

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})

	suite.T().Run("Skips calling check auto reweigh if the updated actual weight matches the db value", func(t *testing.T) {
		moveWeights := &mockservices.MoveWeights{}
		mockSender := setUpMockNotificationSender()
		mockedUpdater := NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, &mockShipmentRecalculator)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		actualWeight := unit.Pound(7200)
		primeShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
				PrimeActualWeight:   &actualWeight,
			},
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
		})
		// there is a validator check about updating the status
		primeShipment.Status = ""
		primeShipment.PrimeActualWeight = &actualWeight

		_, err := mockedUpdater.UpdateMTOShipmentPrime(suite.AppContextForTest(), &primeShipment, etag.GenerateEtag(primeShipment.UpdatedAt))
		suite.NoError(err)

		moveWeights.AssertNotCalled(t, "CheckAutoReweigh")

		// Verify that shipment recalculate was handled correctly
		mockShipmentRecalculator.AssertNotCalled(t, "ShipmentRecalculatePaymentRequest", mock.Anything, mock.Anything)
	})
}
