package mtoshipment

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestMTOShipmentUpdater() {
	oldMTOShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
	mtoShipmentUpdater := NewMTOShipmentUpdater(suite.DB())

	requestedPickupDate := *oldMTOShipment.RequestedPickupDate
	scheduledPickupDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	firstAvailableDeliveryDate := time.Date(2019, time.March, 10, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(2020, time.June, 8, 0, 0, 0, 0, time.UTC)

	secondaryPickupAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
	secondaryDeliveryAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})
	primeActualWeight := unit.Pound(1234)

	mtoShipment := models.MTOShipment{
		ID:                         oldMTOShipment.ID,
		DestinationAddress:         oldMTOShipment.DestinationAddress,
		PickupAddress:              oldMTOShipment.PickupAddress,
		RequestedPickupDate:        &requestedPickupDate,
		ScheduledPickupDate:        &scheduledPickupDate,
		ShipmentType:               "INTERNATIONAL_UB",
		SecondaryPickupAddress:     &secondaryPickupAddress,
		SecondaryDeliveryAddress:   &secondaryDeliveryAddress,
		PrimeActualWeight:          &primeActualWeight,
		FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
		ActualPickupDate:           &actualPickupDate,
	}

	suite.T().Run("If-Unmodified-Since is not equal to the updated_at date", func(t *testing.T) {
		unmodifiedSince := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)

		_, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment, unmodifiedSince)
		suite.Error(err)
		suite.IsType(ErrPreconditionFailed{}, err)
	})

	suite.T().Run("If-Unmodified-Since is equal to the updated_at date", func(t *testing.T) {
		unmodifiedSince := oldMTOShipment.UpdatedAt
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment, unmodifiedSince)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)

		suite.NotZero(updatedMTOShipment.PickupAddress.ID, oldMTOShipment.PickupAddress.ID)

		suite.NotZero(updatedMTOShipment.SecondaryPickupAddress.ID, secondaryPickupAddress.ID)
		suite.NotZero(updatedMTOShipment.SecondaryDeliveryAddress.ID, secondaryDeliveryAddress.ID)
		suite.Equal(updatedMTOShipment.PrimeActualWeight, &primeActualWeight)
		suite.True(actualPickupDate.Equal(*updatedMTOShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*updatedMTOShipment.FirstAvailableDeliveryDate))
	})

	oldMTOShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
	mtoShipment2 := models.MTOShipment{
		ID:           oldMTOShipment2.ID,
		ShipmentType: "INTERNATIONAL_UB",
	}

	suite.T().Run("Updater can handle optional queries set as nil", func(t *testing.T) {
		unmodifiedSince := oldMTOShipment2.UpdatedAt

		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment2, unmodifiedSince)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment2.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)
		suite.Nil(updatedMTOShipment.PrimeEstimatedWeight)
	})

	now := time.Now()
	primeEstimatedWeight := unit.Pound(4500)

	suite.T().Run("Failed case if not both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		eightDaysFromNow := now.AddDate(0, 0, 8)
		threeDaysBefore := now.AddDate(0, 0, -3)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &eightDaysFromNow,
				ApprovedDate:        &threeDaysBefore,
			},
		})
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		unmodifiedSince := oldShipment.UpdatedAt

		_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, unmodifiedSince)
		suite.Error(err)
	})

	suite.T().Run("Successful case if both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		tenDaysFromNow := now.AddDate(0, 0, 11)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &tenDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		unmodifiedSince := oldShipment.UpdatedAt
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, unmodifiedSince)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	})

	suite.T().Run("Failed case if approved date is 3-9 days from scheduled move date but estimated weight recorded date isn't at least 3 days prior to scheduled move date", func(t *testing.T) {
		twoDaysFromNow := now.AddDate(0, 0, 2)
		twoDaysBefore := now.AddDate(0, 0, -2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &twoDaysFromNow,
				ApprovedDate:        &twoDaysBefore,
			},
		})
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		unmodifiedSince := oldShipment.UpdatedAt

		_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, unmodifiedSince)
		suite.Error(err)
	})

	suite.T().Run("Successful case if approved date is 3-9 days from scheduled move date and estimated weight recorded date is at least 3 days prior to scheduled move date", func(t *testing.T) {
		sixDaysFromNow := now.AddDate(0, 0, 6)
		twoDaysBefore := now.AddDate(0, 0, -2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &sixDaysFromNow,
				ApprovedDate:        &twoDaysBefore,
			},
		})
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		unmodifiedSince := oldShipment.UpdatedAt
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, unmodifiedSince)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	})

	suite.T().Run("Failed case if approved date is less than 3 days from scheduled move date but estimated weight recorded date isn't at least 1 day prior to scheduled move date", func(t *testing.T) {
		oneDayFromNow := now.AddDate(0, 0, 1)
		oneDayBefore := now.AddDate(0, 0, -1)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &oneDayFromNow,
				ApprovedDate:        &oneDayBefore,
			},
		})
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		unmodifiedSince := oldShipment.UpdatedAt

		_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, unmodifiedSince)
		suite.Error(err)
	})

	suite.T().Run("Successful case if approved date is less than 3 days from scheduled move date and estimated weight recorded date is at least 1 day prior to scheduled move date", func(t *testing.T) {
		twoDaysFromNow := now.AddDate(0, 0, 2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &twoDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		unmodifiedSince := oldShipment.UpdatedAt
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, unmodifiedSince)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	})
}
