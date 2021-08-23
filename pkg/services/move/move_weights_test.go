package move

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MoveServiceSuite) TestExcessWeight() {
	moveWeights := NewMoveWeights()

	suite.Run("qualifies move for excess weight when an approved shipment estimated weight is updated within threshold", func() {
		// The default weight allotment for this move is 8000 and the threshold is 90% of that
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: approvedMove,
		})

		estimatedWeight := unit.Pound(7200)
		approvedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.TestAppContext(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.NotNil(updatedMove.ExcessWeightQualifiedAt)

		// refetch the move from the database not just the return value
		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.NotNil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("does not flag move for excess weight when an approved shipment estimated weight is lower than threshold", func() {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: approvedMove,
		})

		estimatedWeight := unit.Pound(7199)
		approvedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.TestAppContext(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("qualifies move for excess weight when the sum of approved shipments is updated within threshold", func() {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(3600)

		testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:               models.MTOShipmentStatusApproved,
				ApprovedDate:         &now,
				ScheduledPickupDate:  &pickupDate,
				PrimeEstimatedWeight: &estimatedWeight,
			},
			Move: approvedMove,
		})
		approvedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: approvedMove,
		})

		approvedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.TestAppContext(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.NotNil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.NotNil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("does not flag move for excess weight when the sum of non-approved shipments meets the threshold", func() {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(3600)

		testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:               models.MTOShipmentStatusCanceled,
				ApprovedDate:         &now,
				ScheduledPickupDate:  &pickupDate,
				PrimeEstimatedWeight: &estimatedWeight,
			},
			Move: approvedMove,
		})
		approvedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: approvedMove,
		})

		approvedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.TestAppContext(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("does not flag move for excess weight when updated shipment status is not approved", func() {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)

		unapprovedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusSubmitted,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
				Diversion:           true,
			},
			Move: approvedMove,
		})

		estimatedWeight := unit.Pound(7200)
		unapprovedShipment.PrimeEstimatedWeight = &estimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.TestAppContext(), approvedMove.ID, unapprovedShipment)

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
		approvedMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt:      &now,
				Status:                  models.MoveStatusAPPROVED,
				ExcessWeightQualifiedAt: &now,
			},
		})

		pickupDate := now.AddDate(0, 0, 10)
		estimatedWeight := unit.Pound(7200)
		approvedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:               models.MTOShipmentStatusApproved,
				ApprovedDate:         &now,
				ScheduledPickupDate:  &pickupDate,
				PrimeEstimatedWeight: &estimatedWeight,
			},
			Move: approvedMove,
		})

		updatedEstimatedWeight := unit.Pound(7199)
		approvedShipment.PrimeEstimatedWeight = &updatedEstimatedWeight
		updatedMove, verrs, err := moveWeights.CheckExcessWeight(suite.TestAppContext(), approvedMove.ID, approvedShipment)

		suite.Nil(verrs)
		suite.NoError(err)

		suite.NotNil(approvedMove.ExcessWeightQualifiedAt)
		suite.Nil(updatedMove.ExcessWeightQualifiedAt)

		err = suite.DB().Reload(&approvedMove)
		suite.NoError(err)
		suite.Nil(approvedMove.ExcessWeightQualifiedAt)
	})

	suite.Run("returns error if orders grade is unset to lookup weight allowance", func() {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		approvedMove.Orders.Grade = nil

		err := suite.DB().Save(&approvedMove.Orders)
		suite.NoError(err)

		_, verrs, err := moveWeights.CheckExcessWeight(suite.TestAppContext(), approvedMove.ID, models.MTOShipment{})
		suite.Nil(verrs)
		suite.EqualError(err, "could not determine excess weight entitlement without grade")
	})

	suite.Run("returns error if dependents authorized is unset to lookup weight allowance", func() {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		approvedMove.Orders.Entitlement.DependentsAuthorized = nil

		err := suite.DB().Save(approvedMove.Orders.Entitlement)
		suite.NoError(err)

		_, verrs, err := moveWeights.CheckExcessWeight(suite.TestAppContext(), approvedMove.ID, models.MTOShipment{})
		suite.Nil(verrs)
		suite.EqualError(err, "could not determine excess weight entitlement without dependents authorization value")
	})
}
