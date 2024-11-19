package move

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MoveServiceSuite) TestExcessWeight() {
	moveWeights := NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

	suite.Run("qualifies move for excess weight when an approved shipment estimated weight is updated within threshold", func() {
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

		estimatedWeight := unit.Pound(7200)
		approvedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.NotNil(updatedMove.ExcessWeightQualifiedAt)

		// refetch the move from the database not just the return value
		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.NotNil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("qualifies move for excess weight when an approved UB shipment estimated weight is updated within threshold", func() {
		dependentsAuthorized := true
		dependentsUnderTwelve := 1
		dependentsOverTwelve := 2
		ubAllowance := 2000
		primeEstimatedWeight := 1833

		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Entitlement{
					ID:                      uuid.Must(uuid.NewV4()),
					ProGearWeight:           200,
					UBAllowance:             &ubAllowance,
					DependentsAuthorized:    &dependentsAuthorized,
					DependentsUnderTwelve:   &dependentsUnderTwelve,
					DependentsTwelveAndOver: &dependentsOverTwelve,
				},
			},
		}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					ScheduledPickupDate:  &pickupDate,
					ShipmentType:         models.MTOShipmentTypeUnaccompaniedBaggage,
					PrimeEstimatedWeight: (*unit.Pound)(&primeEstimatedWeight),
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), move.ID, approvedShipment)

		suite.NotNil(updatedMove)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(move.ExcessWeightQualifiedAt)
		suite.NotNil(updatedMove.ExcessWeightQualifiedAt)

		// refetch the move from the database not just the return value
		err = suite.DB().Reload(&move)
		suite.NoError(err)
		suite.NotNil(move.ExcessWeightQualifiedAt)
	})

	suite.Run("does not flag move for excess weight when an approved shipment estimated weight is lower than threshold", func() {
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

		estimatedWeight := unit.Pound(7199)
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
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
		}, nil)

		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(7200)
		approvedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
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

		updatedEstimatedWeight := unit.Pound(7199)
		approvedShipment.PrimeEstimatedWeight = &updatedEstimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.AppContextForTest(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.NotNil(approvedMove.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
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
	moveWeights := NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

	suite.Run("requests reweigh on shipment if the acutal weight is 90% of the weight allowance", func() {
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
		autoReweighShipments, err := moveWeights.CheckAutoReweigh(session, approvedMove.ID, &approvedShipment)

		suite.NoError(err)

		suite.NotNil(approvedShipment.Reweigh)
		suite.Equal(approvedShipment.ID.String(), approvedShipment.Reweigh.ShipmentID.String())
		suite.Equal(models.ReweighRequesterSystem, approvedShipment.Reweigh.RequestedBy)
		suite.NotNil(approvedShipment.Reweigh.RequestedAt)
		suite.NotNil(autoReweighShipments)
	})

	suite.Run("does not request reweigh on shipments when below 90% of weight allowance threshold", func() {
		mockedReweighRequestor := mocks.ShipmentReweighRequester{}
		mockedWeightService := NewMoveWeights(&mockedReweighRequestor)

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

		actualWeight := unit.Pound(7199)
		approvedShipment.PrimeEstimatedWeight = &actualWeight
		_, err := mockedWeightService.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &approvedShipment)

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
		autoReweighShipments, err := moveWeights.CheckAutoReweigh(session, approvedMove.ID, &approvedShipment)

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
		suite.Equal(len(autoReweighShipments), 2)
	})

	suite.Run("does not request reweigh when shipments aren't in approved statuses", func() {
		mockedReweighRequestor := mocks.ShipmentReweighRequester{}
		mockedWeightService := NewMoveWeights(&mockedReweighRequestor)

		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		actualWeight := unit.Pound(3600)

		existingShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusCanceled,
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
		_, err := mockedWeightService.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &approvedShipment)

		suite.NoError(err)

		err = suite.DB().Eager("Reweigh").Reload(&existingShipment)
		suite.NoError(err)
		suite.Equal(uuid.Nil, existingShipment.Reweigh.ID)
		suite.Equal(uuid.Nil, approvedShipment.Reweigh.ID)
		mockedReweighRequestor.AssertNotCalled(suite.T(), "RequestShipmentReweigh")
	})

	suite.Run("uses lower reweigh weight on shipments that already have reweighs", func() {
		mockedReweighRequestor := mocks.ShipmentReweighRequester{}
		mockedWeightService := NewMoveWeights(&mockedReweighRequestor)
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		actualWeight := unit.Pound(2400)
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

		reweighedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
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

		approvedShipment.PrimeActualWeight = &actualWeight
		_, err := mockedWeightService.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &approvedShipment)

		suite.NoError(err)

		err = suite.DB().Eager("Reweigh").Reload(&existingShipment)
		suite.NoError(err)
		suite.Equal(uuid.Nil, existingShipment.Reweigh.ID)
		suite.Equal(uuid.Nil, approvedShipment.Reweigh.ID)
		mockedReweighRequestor.AssertNotCalled(suite.T(), "RequestShipmentReweigh")
	})

	suite.Run("returns error if orders grade is unset to lookup weight allowance", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		approvedMove.Orders.Grade = nil

		err := suite.DB().Save(&approvedMove.Orders)
		suite.NoError(err)

		_, err = moveWeights.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &models.MTOShipment{})
		suite.EqualError(err, "could not determine excess weight entitlement without grade")
	})

	suite.Run("returns error if dependents authorized is unset to lookup weight allowance", func() {
		approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		approvedMove.Orders.Entitlement.DependentsAuthorized = nil

		err := suite.DB().Save(approvedMove.Orders.Entitlement)
		suite.NoError(err)

		_, err = moveWeights.CheckAutoReweigh(suite.AppContextForTest(), approvedMove.ID, &models.MTOShipment{})
		suite.EqualError(err, "could not determine excess weight entitlement without dependents authorization value")
	})
}
