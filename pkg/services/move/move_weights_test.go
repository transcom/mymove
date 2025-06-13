package move

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	notificationMocks "github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/services/entitlements"
	mocks "github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/reweigh"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func setUpMockNotificationSender() notifications.NotificationSender {
	// The UserUpdater needs a NotificationSender for sending user activity emails to system admins.
	// This function allows us to set up a fresh mock for each test so we can check the number of calls it has.
	mockSender := notificationMocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
	).Return(nil)

	return &mockSender
}

func (suite *MoveServiceSuite) TestExcessWeight() {
	waf := entitlements.NewWeightAllotmentFetcher()
	mockSender := setUpMockNotificationSender()

	moveWeights := NewMoveWeights(mtoshipment.NewShipmentReweighRequester(mockSender), waf)

	suite.Run("qualifies move for excess weight when an approved shipment estimated weight is updated within threshold", func() {
		// The default weight allotment for this move is 8000 and the threshold is 90% of that
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedHHGShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)
		approvedUbShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
					ShipmentType:        models.MTOShipmentTypeUnaccompaniedBaggage,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		estimatedWeight := unit.Pound(7200)
		approvedHHGShipment.PrimeEstimatedWeight = &estimatedWeight
		approvedUbShipment.PrimeEstimatedWeight = &estimatedWeight
		_, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, approvedHHGShipment)
		suite.Nil(verrs)
		suite.NoError(err)
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, approvedUbShipment)
		suite.Nil(verrs)
		suite.NoError(err)

		// Move has nil excess weight risks before checking for excess weight
		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.Nil(approvedMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)
		// Move has not nil excess weight risks after checking for excess weight
		suite.NotNil(updatedMove.ExcessWeightQualifiedAt)
		suite.NotNil(updatedMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)

		// refetch the move from the database not just the return value
		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		// Ensure it saved to db
		suite.NotNil(approvedMove.ExcessWeightQualifiedAt)
		suite.NotNil(approvedMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)
	})

	suite.Run("does not flag move for excess weight when an approved shipment estimated weight is lower than threshold", func() {
		// Create a move with an oconus duty location so it qualifies for UB allowance
		// The allowance based on these params should be 500 ub
		oconusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus: models.BoolPointer(true),
				},
			},
		}, nil)
		oconusDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model:    oconusAddress,
				LinkOnly: true,
			},
		}, nil)
		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    oconusDutyLocation,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    order,
				LinkOnly: true,
			}}, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)

		approvedHHGShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		approvedUbShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
					ShipmentType:        models.MTOShipmentTypeUnaccompaniedBaggage,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		estimatedHHGWeight := unit.Pound(7199)
		estimatedUBWeight := unit.Pound(250)
		approvedHHGShipment.PrimeEstimatedWeight = &estimatedHHGWeight
		approvedUbShipment.PrimeEstimatedWeight = &estimatedUBWeight
		_, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), move.ID, approvedHHGShipment)
		suite.Nil(verrs)
		suite.NoError(err)
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), move.ID, approvedUbShipment)
		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(move.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)
		suite.Nil(move.ExcessUnaccompaniedBaggageWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)

		err = suite.DB().Reload(&move)
		suite.NoError(err)
		suite.Nil(move.ExcessWeightQualifiedAt)
		suite.Nil(move.ExcessUnaccompaniedBaggageWeightQualifiedAt)
	})

	suite.Run("qualifies move for excess weight when the sum of approved shipments is updated within threshold", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(3600)

		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		approvedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.NotNil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.NotNil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("does not flag move for excess weight when the sum of non-approved shipments meets the threshold", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(3600)

		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusCanceled,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		approvedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("does not flag move for excess weight when updated shipment status is not approved", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)

		unapprovedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
					Diversion:           true,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		estimatedWeight := unit.Pound(7200)
		unapprovedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, unapprovedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("removes excess weight qualification when estimated weight drops below previously met threshold", func() {
		now := time.Now()
		estimatedUbWeight := unit.Pound(250)
		estimatedWeight := unit.Pound(7199 - estimatedUbWeight)

		// Add an OCONUS address so it qualifies for UB allowance
		// The allowance based on these params should be 500 ub
		oconusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus: models.BoolPointer(true),
				},
			},
		}, nil)
		oconusDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model:    oconusAddress,
				LinkOnly: true,
			},
		}, nil)
		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    oconusDutyLocation,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)
		// By default have excess weight turned on, we want to simulate it resetting
		initialMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt:                     &now,
					ExcessUnaccompaniedBaggageWeightQualifiedAt: &now,
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			}}, nil)

		pickupDate := now.AddDate(0, 0, 10)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model:    initialMove,
				LinkOnly: true,
			},
		}, nil)
		approvedUbShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					ShipmentType:         models.MTOShipmentTypeUnaccompaniedBaggage,
					PrimeEstimatedWeight: &estimatedUbWeight,
				},
			},
			{
				Model:    initialMove,
				LinkOnly: true,
			},
		}, nil)

		// We defaulted to excess amounts
		suite.NotNil(initialMove.ExcessWeightQualifiedAt)
		suite.NotNil(initialMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)

		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), initialMove.ID, approvedUbShipment)
		suite.Nil(verrs)
		suite.NoError(err)

		// The shipments we created will not qualify for risk of excess
		// This means that after we CheckExcessWeight again, the
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)

		// Reload our original move that had excess weight qualified at present
		// and now make sure it is nil
		err = suite.DB().Reload(&initialMove)
		suite.NoError(err)
		suite.Nil(initialMove.ExcessWeightQualifiedAt)
		suite.Nil(initialMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)
	})

	suite.Run("returns error if orders grade is unset to lookup weight allowance", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		approvedMove.Orders.Grade = nil

		err := suite.DB().Save(&approvedMove.Orders)
		suite.NoError(err)

		_, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, models.MTOShipment{})
		suite.Nil(verrs)
		suite.EqualError(err, "could not determine excess weight entitlement without grade")
	})

	suite.Run("returns error if dependents authorized is unset to lookup weight allowance", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		approvedMove.Orders.Entitlement.DependentsAuthorized = nil

		err := suite.DB().Save(approvedMove.Orders.Entitlement)
		suite.NoError(err)

		_, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, models.MTOShipment{})
		suite.Nil(verrs)
		suite.EqualError(err, "could not determine excess weight entitlement without dependents authorization value")
	})

	suite.Run("qualifies move for excess weight when an approved shipment with PPM weights is greater than threshold", func() {
		// The default weight allotment for this move is 8000 and the threshold is 90% of that
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		//Default estimatedWeight for ppm is 4000
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		primeEstimatedWeight := unit.Pound(3200)
		approvedShipment.PrimeEstimatedWeight = &primeEstimatedWeight
		approvedShipment.PPMShipment = &ppmShipment
		//When accounting for PPM weight, the sum should exceed the 90% threshold
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.NotNil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.NotNil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("does not flag move for excess weight when an approved shipment with PPM weights is below the threshold", func() {
		// The default weight allotment for this move is 8000 and the threshold is 90% of that
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		//Default estimatedWeight for ppm is 4000
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		primeEstimatedWeight := unit.Pound(3199)
		approvedShipment.PrimeEstimatedWeight = &primeEstimatedWeight
		//When accounting for PPM weight, the sum should NOT exceed the 90% threshold
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
	})
}

func (suite *MoveServiceSuite) TestAutoReweigh() {
	waf := entitlements.NewWeightAllotmentFetcher()
	mockSender := setUpMockNotificationSender()
	moveWeights := NewMoveWeights(mtoshipment.NewShipmentReweighRequester(mockSender), waf)

	suite.Run("requests reweigh on shipment if the actual weight is 90% of the weight allowance", func() {
		// The default weight allotment for this move is 8000 and the threshold is 90% of that
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		actualWeight := unit.Pound(7200)
		approvedShipment.PrimeActualWeight = &actualWeight
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		err := moveWeights.CheckAutoReweigh(session, approvedMove.ID, &approvedShipment)
		suite.NoError(err)
		err = suite.DB().Eager("Reweigh").Reload(&approvedShipment)
		suite.NoError(err)

		suite.NotNil(approvedShipment.Reweigh)
		suite.Equal(approvedShipment.ID.String(), approvedShipment.Reweigh.ShipmentID.String())
		suite.Equal(models.ReweighRequesterSystem, approvedShipment.Reweigh.RequestedBy)
		suite.NotNil(approvedShipment.Reweigh.RequestedAt)
	})

	suite.Run("does not request reweigh on shipments when below 90% of weight allowance threshold", func() {
		mockedReweighRequestor := mocks.ShipmentReweighRequester{}
		mockedWeightService := NewMoveWeights(&mockedReweighRequestor, waf)

		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		weight := unit.Pound(7199)
		approvedShipment.PrimeActualWeight = &weight
		approvedShipment.PrimeEstimatedWeight = &weight
		err := mockedWeightService.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &approvedShipment)
		suite.NoError(err)
		err = suite.DB().Eager("Reweigh").Reload(&approvedShipment)
		suite.NoError(err)

		suite.Equal(uuid.Nil, approvedShipment.Reweigh.ID)
		mockedReweighRequestor.AssertNotCalled(suite.T(), "RequestShipmentReweigh")
	})

	suite.Run("requests reweigh on existing shipments in addition to the one being updated", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		actualWeight := unit.Pound(3600)

		existingShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
					PrimeActualWeight:   &actualWeight,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		approvedShipment.PrimeActualWeight = &actualWeight
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		err := moveWeights.CheckAutoReweigh(session, approvedMove.ID, &approvedShipment)
		suite.NoError(err)
		err = suite.DB().Eager("Reweigh").Reload(&approvedShipment)
		suite.NoError(err)

		suite.NotNil(approvedShipment.Reweigh)
		suite.NotEqual(uuid.Nil, approvedShipment.Reweigh.ID)
		suite.Equal(approvedShipment.ID, approvedShipment.Reweigh.ShipmentID)
		suite.Equal(models.ReweighRequesterSystem, approvedShipment.Reweigh.RequestedBy)

		err = suite.DB().Eager("Reweigh").Reload(&existingShipment)
		suite.NoError(err)

		suite.NotNil(existingShipment.Reweigh)
		suite.NotEqual(uuid.Nil, existingShipment.Reweigh.ID)
		suite.Equal(existingShipment.ID, existingShipment.Reweigh.ShipmentID)
		suite.Equal(models.ReweighRequesterSystem, existingShipment.Reweigh.RequestedBy)
	})

	suite.Run("does not request reweigh when shipments aren't in approved statuses", func() {
		mockedReweighRequestor := mocks.ShipmentReweighRequester{}
		mockedWeightService := NewMoveWeights(&mockedReweighRequestor, waf)

		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		weight := unit.Pound(3600)

		existingShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusCanceled,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					PrimeEstimatedWeight: &weight,
					PrimeActualWeight:    &weight,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		approvedShipment.PrimeEstimatedWeight = &weight
		approvedShipment.PrimeActualWeight = &weight
		err := mockedWeightService.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &approvedShipment)
		suite.NoError(err)
		err = suite.DB().Eager("Reweigh").Reload(&approvedShipment)
		suite.NoError(err)

		err = suite.DB().Eager("Reweigh").Reload(&existingShipment)
		suite.NoError(err)
		suite.Equal(uuid.Nil, existingShipment.Reweigh.ID)
		suite.Equal(uuid.Nil, approvedShipment.Reweigh.ID)
		mockedReweighRequestor.AssertNotCalled(suite.T(), "RequestShipmentReweigh")
	})

	suite.Run("requests reweighs for all shipments if any other shipments have active reweighs", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(2400)
		existingShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		reweighedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)
		reweighWeight := unit.Pound(2399)
		testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			Reweigh: models.Reweigh{
				Weight: &reweighWeight,
			},
			MTOShipment: reweighedShipment,
		})

		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		approvedShipment.PrimeEstimatedWeight = &estimatedWeight

		err := suite.DB().Eager("Reweigh").Reload(&approvedShipment)
		suite.NoError(err)

		err = suite.DB().Eager("Reweigh").Reload(&existingShipment)
		suite.NoError(err)

		// Test that entire function returns without error
		reweighErr := moveWeights.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &approvedShipment)
		suite.NoError(reweighErr)

	})

	suite.Run("returns error if orders grade is unset to lookup weight allowance", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		approvedMove.Orders.Grade = nil

		err := suite.DB().Save(&approvedMove.Orders)
		suite.NoError(err)

		err = moveWeights.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &models.MTOShipment{})
		suite.EqualError(err, "could not determine excess weight entitlement without grade")
	})

	suite.Run("returns error if dependents authorized is unset to lookup weight allowance", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		approvedMove.Orders.Entitlement.DependentsAuthorized = nil

		err := suite.DB().Save(approvedMove.Orders.Entitlement)
		suite.NoError(err)

		err = moveWeights.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &models.MTOShipment{})
		suite.EqualError(err, "could not determine excess weight entitlement without dependents authorization value")
	})

	suite.Run("returns error if can't find move when checking for auto-reweigh", func() {
		randomID, err := uuid.NewV4()
		suite.NoError(err)
		err = moveWeights.CheckAutoReweigh(suite.AppContextForTest(), randomID, &models.MTOShipment{})
		suite.EqualError(err, apperror.NewNotFoundError(randomID, "looking for Move").Error())
	})

	suite.Run("returns error if shipment returns nil when checking for auto-reweigh", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		err := suite.DB().Save(approvedMove.Orders.Entitlement)
		suite.NoError(err)

		err = moveWeights.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, nil)
		suite.EqualError(err, apperror.NewBadDataError("received a nil MTO shipment, an MTO shipment must be supplied for checking reweighs").Error())
	})

	suite.Run("will correctly factor in dependents authorized when requesting a reweigh", func() {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		actualWeight := unit.Pound(5000)

		move1 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment1 := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    move1,
				LinkOnly: true,
			},
		}, nil)
		move1.Orders.Entitlement.DependentsAuthorized = models.BoolPointer(false)
		err := suite.DB().Save(move1.Orders.Entitlement)
		suite.NoError(err)
		shipment1.PrimeActualWeight = &actualWeight

		move2 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment2 := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    move2,
				LinkOnly: true,
			},
		}, nil)
		move2.Orders.Entitlement.DependentsAuthorized = models.BoolPointer(true)
		err = suite.DB().Save(move2.Orders.Entitlement)
		suite.NoError(err)
		shipment2.PrimeActualWeight = &actualWeight

		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		// Test that correct number of shipments are returned by GetAutoReweighShipments
		var move models.Move
		err = suite.DB().EagerPreload("MTOShipments", "Orders", "Orders.Entitlement", "MTOShipments.Reweigh", "MTOShipments.ShipmentType", "MTOShipments.Status", "MTOShipments.DeletedAt", "MTOShipments.PrimeActualWeight", "MTOShipments.PrimeEstimatedWeight").Find(&move, move1.ID)
		suite.NoError(err)
		autoReweighShipments, err := moveWeights.GetAutoReweighShipments(session, &move, &shipment1)
		suite.NoError(err)
		suite.Equal(1, len(*autoReweighShipments))

		err = suite.DB().EagerPreload("MTOShipments", "Orders", "Orders.Entitlement", "MTOShipments.Reweigh", "MTOShipments.ShipmentType", "MTOShipments.Status", "MTOShipments.DeletedAt", "MTOShipments.PrimeActualWeight", "MTOShipments.PrimeEstimatedWeight").Find(&move, move2.ID)
		suite.NoError(err)
		autoReweighShipments, err = moveWeights.GetAutoReweighShipments(session, &move, &shipment2)
		suite.NoError(err)
		suite.Equal(0, len(*autoReweighShipments))

		// Test that entire function executes without error
		err = moveWeights.CheckAutoReweigh(session, move1.ID, &shipment1)
		suite.NoError(err)

		err = moveWeights.CheckAutoReweigh(session, move2.ID, &shipment2)
		suite.NoError(err)
		suite.Equal(0, len(*autoReweighShipments))
	})

	suite.Run("doesn't request a reweigh if one already exists for a shipment", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		now := time.Now()

		pickupDate := now.AddDate(0, 0, 10)
		actualWeight := unit.Pound(3600)

		existingShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
					PrimeActualWeight:   &actualWeight,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ApprovedDate:        &now,
					ScheduledPickupDate: &pickupDate,
				},
			},
			{
				Model:    approvedMove,
				LinkOnly: true,
			},
		}, nil)

		approvedShipment.PrimeActualWeight = &actualWeight

		// Create an existing reweigh
		approvedShipmentReweighModel := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterPrime,
			ShipmentID:  approvedShipment.ID,
		}

		existingShipmentReweighModel := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterPrime,
			ShipmentID:  existingShipment.ID,
		}

		reweighCreator := reweigh.NewReweighCreator()
		approvedShipmentReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), approvedShipmentReweighModel)
		suite.NoError(err)
		existingShipmentReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), existingShipmentReweighModel)
		suite.NoError(err)
		suite.NotNil(approvedShipmentReweigh)
		suite.NotNil(existingShipmentReweigh)

		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		err = moveWeights.CheckAutoReweigh(session, approvedMove.ID, &approvedShipment)
		suite.NoError(err)

		autoReweighShipments, err := moveWeights.GetAutoReweighShipments(suite.AppContextForTest(), &approvedMove, &approvedShipment)
		suite.NoError(err)
		suite.Empty(autoReweighShipments)
	})
}
