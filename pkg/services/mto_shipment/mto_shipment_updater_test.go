package mtoshipment

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestMTOShipmentUpdater() {
	oldMTOShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
	builder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(builder)
	mtoShipmentUpdater := NewMTOShipmentUpdater(suite.DB(), builder, fetcher)

	requestedPickupDate := *oldMTOShipment.RequestedPickupDate
	scheduledPickupDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	firstAvailableDeliveryDate := time.Date(2019, time.March, 10, 0, 0, 0, 0, time.UTC)
	secondaryPickupAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
	secondaryDeliveryAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})
	primeActualWeight := unit.Pound(1234)

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
		SecondaryPickupAddress:     &secondaryPickupAddress,
		SecondaryDeliveryAddress:   &secondaryDeliveryAddress,
		PrimeActualWeight:          &primeActualWeight,
		FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
		Status:                     oldMTOShipment.Status,
	}

	suite.T().Run("Etag is stale", func(t *testing.T) {
		eTag := base64.StdEncoding.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
		_, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment, eTag)
		suite.Error(err)
		suite.IsType(ErrPreconditionFailed{}, err)
	})

	suite.T().Run("If-Unmodified-Since is equal to the updated_at date", func(t *testing.T) {
		eTag := base64.StdEncoding.EncodeToString([]byte(oldMTOShipment.UpdatedAt.Format(time.RFC3339Nano)))
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment, eTag)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)

		suite.NotZero(updatedMTOShipment.PickupAddressID, oldMTOShipment.PickupAddressID)

		suite.NotZero(updatedMTOShipment.SecondaryPickupAddressID, secondaryPickupAddress.ID)
		suite.NotZero(updatedMTOShipment.SecondaryDeliveryAddressID, secondaryDeliveryAddress.ID)
		suite.Equal(updatedMTOShipment.PrimeActualWeight, &primeActualWeight)
		suite.True(firstAvailableDeliveryDate.Equal(*updatedMTOShipment.FirstAvailableDeliveryDate))
	})

	oldMTOShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
	mtoShipment2 := models.MTOShipment{
		ID:           oldMTOShipment2.ID,
		ShipmentType: "INTERNATIONAL_UB",
	}

	suite.T().Run("Updater can handle optional queries set as nil", func(t *testing.T) {
		eTag := base64.StdEncoding.EncodeToString([]byte(oldMTOShipment2.UpdatedAt.Format(time.RFC3339Nano)))

		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment2, eTag)
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
		eTag := base64.StdEncoding.EncodeToString([]byte(oldShipment.UpdatedAt.Format(time.RFC3339Nano)))
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
		suite.Error(err)
	})
	//
	suite.T().Run("Successful case if both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		tenDaysFromNow := now.AddDate(0, 0, 11)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &tenDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		eTag := base64.StdEncoding.EncodeToString([]byte(oldShipment.UpdatedAt.Format(time.RFC3339Nano)))
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
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
		eTag := base64.StdEncoding.EncodeToString([]byte(oldShipment.UpdatedAt.Format(time.RFC3339Nano)))
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
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
		eTag := base64.StdEncoding.EncodeToString([]byte(oldShipment.UpdatedAt.Format(time.RFC3339Nano)))
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
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
		eTag := base64.StdEncoding.EncodeToString([]byte(oldShipment.UpdatedAt.Format(time.RFC3339Nano)))
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}

		_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
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
		eTag := base64.StdEncoding.EncodeToString([]byte(oldShipment.UpdatedAt.Format(time.RFC3339Nano)))
		updatedShipment := models.MTOShipment{
			ID:                   oldShipment.ID,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		}
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	})
}
