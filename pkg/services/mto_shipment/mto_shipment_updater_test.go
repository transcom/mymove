package mtoshipment

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestMTOShipmentUpdater() {
	oldMTOShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
	mtoShipmentUpdater := NewMTOShipmentUpdater(suite.DB())
	requestedPickupDate := strfmt.Date(*oldMTOShipment.RequestedPickupDate)
	scheduledPickupDate := strfmt.Date(time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC))
	firstAvailableDeliveryDate := strfmt.Date(time.Date(2019, time.March, 10, 0, 0, 0, 0, time.UTC))
	pickupAddress := primemessages.Address{
		City:           &oldMTOShipment.PickupAddress.City,
		Country:        oldMTOShipment.PickupAddress.Country,
		ID:             strfmt.UUID(oldMTOShipment.PickupAddress.ID.String()),
		PostalCode:     &oldMTOShipment.PickupAddress.PostalCode,
		State:          &oldMTOShipment.PickupAddress.State,
		StreetAddress1: &oldMTOShipment.PickupAddress.StreetAddress1,
		StreetAddress2: oldMTOShipment.PickupAddress.StreetAddress2,
		StreetAddress3: oldMTOShipment.PickupAddress.StreetAddress3,
	}

	destinationAddress := primemessages.Address{
		City:           &oldMTOShipment.DestinationAddress.City,
		Country:        oldMTOShipment.DestinationAddress.Country,
		ID:             strfmt.UUID(oldMTOShipment.DestinationAddress.ID.String()),
		PostalCode:     &oldMTOShipment.DestinationAddress.PostalCode,
		State:          &oldMTOShipment.DestinationAddress.State,
		StreetAddress1: &oldMTOShipment.DestinationAddress.StreetAddress1,
		StreetAddress2: oldMTOShipment.DestinationAddress.StreetAddress2,
		StreetAddress3: oldMTOShipment.DestinationAddress.StreetAddress3,
	}

	secondaryPickupAddressModel := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})

	secondaryPickupAddress := primemessages.Address{
		City:           &secondaryPickupAddressModel.City,
		Country:        secondaryPickupAddressModel.Country,
		ID:             strfmt.UUID(secondaryPickupAddressModel.ID.String()),
		PostalCode:     &secondaryPickupAddressModel.PostalCode,
		State:          &secondaryPickupAddressModel.State,
		StreetAddress1: &secondaryPickupAddressModel.StreetAddress1,
		StreetAddress2: secondaryPickupAddressModel.StreetAddress2,
		StreetAddress3: secondaryPickupAddressModel.StreetAddress3,
	}

	secondaryDeliveryAddressModel := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})

	secondaryDeliveryAddress := primemessages.Address{
		City:           &secondaryDeliveryAddressModel.City,
		Country:        secondaryDeliveryAddressModel.Country,
		ID:             strfmt.UUID(secondaryDeliveryAddressModel.ID.String()),
		PostalCode:     &secondaryDeliveryAddressModel.PostalCode,
		State:          &secondaryDeliveryAddressModel.State,
		StreetAddress1: &secondaryDeliveryAddressModel.StreetAddress1,
		StreetAddress2: secondaryDeliveryAddressModel.StreetAddress2,
		StreetAddress3: secondaryDeliveryAddressModel.StreetAddress3,
	}

	payload := primemessages.MTOShipment{
		ID:                         strfmt.UUID(oldMTOShipment.ID.String()),
		DestinationAddress:         &destinationAddress,
		PickupAddress:              &pickupAddress,
		RequestedPickupDate:        requestedPickupDate,
		ScheduledPickupDate:        scheduledPickupDate,
		ShipmentType:               "INTERNATIONAL_UB",
		SecondaryPickupAddress:     &secondaryPickupAddress,
		SecondaryDeliveryAddress:   &secondaryDeliveryAddress,
		PrimeActualWeight:          123,
		FirstAvailableDeliveryDate: firstAvailableDeliveryDate,
	}

	suite.T().Run("If-Unmodified-Since is not equal to the updated_at date", func(t *testing.T) {
		unmodifiedSince := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload,
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		_, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.Error(err)
		suite.IsType(ErrPreconditionFailed{}, err)
	})

	suite.T().Run("If-Unmodified-Since is equal to the updated_at date", func(t *testing.T) {
		weight := unit.Pound(123)
		actualWeight := &weight
		unmodifiedSince := oldMTOShipment.UpdatedAt

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload,
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)

		suite.NotZero(updatedMTOShipment.PickupAddress.ID, pickupAddress.ID)
		suite.Equal(updatedMTOShipment.PickupAddress.StreetAddress1, *pickupAddress.StreetAddress1)
		suite.NotZero(updatedMTOShipment.DestinationAddress.ID, destinationAddress.ID)
		suite.Equal(updatedMTOShipment.DestinationAddress.StreetAddress1, *destinationAddress.StreetAddress1)

		suite.NotZero(updatedMTOShipment.SecondaryPickupAddress.ID, secondaryPickupAddress.ID)
		suite.Equal(updatedMTOShipment.SecondaryPickupAddress.StreetAddress1, *secondaryPickupAddress.StreetAddress1)
		suite.NotZero(updatedMTOShipment.SecondaryDeliveryAddress.ID, secondaryDeliveryAddress.ID)
		suite.Equal(updatedMTOShipment.SecondaryDeliveryAddress.StreetAddress1, *secondaryDeliveryAddress.StreetAddress1)
		suite.Equal(updatedMTOShipment.PrimeActualWeight, *&actualWeight)
		suite.True(time.Date(2019, time.March, 10, 0, 0, 0, 0, time.UTC).Equal(*updatedMTOShipment.FirstAvailableDeliveryDate))
	})

	payload2 := primemessages.MTOShipment{
		ID:                       strfmt.UUID(oldMTOShipment.ID.String()),
		DestinationAddress:       &destinationAddress,
		PickupAddress:            &pickupAddress,
		RequestedPickupDate:      requestedPickupDate,
		ScheduledPickupDate:      scheduledPickupDate,
		ShipmentType:             "INTERNATIONAL_UB",
		SecondaryPickupAddress:   nil,
		SecondaryDeliveryAddress: nil,
		PrimeActualWeight:        123,
	}

	suite.T().Run("Updater can handle optional queries set as nil", func(t *testing.T) {
		suite.DB().Find(&oldMTOShipment, oldMTOShipment.ID)
		weight := unit.Pound(123)
		actualWeight := &weight
		unmodifiedSince := oldMTOShipment.UpdatedAt

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload2,
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)

		suite.NotZero(updatedMTOShipment.PickupAddress.ID, pickupAddress.ID)
		suite.Equal(updatedMTOShipment.PickupAddress.StreetAddress1, *pickupAddress.StreetAddress1)
		suite.NotZero(updatedMTOShipment.DestinationAddress.ID, destinationAddress.ID)
		suite.Equal(updatedMTOShipment.DestinationAddress.StreetAddress1, *destinationAddress.StreetAddress1)

		suite.NotZero(updatedMTOShipment.SecondaryPickupAddress.ID, secondaryPickupAddress.ID)
		suite.Equal(updatedMTOShipment.SecondaryPickupAddress.StreetAddress1, *secondaryPickupAddress.StreetAddress1)
		suite.NotZero(updatedMTOShipment.SecondaryDeliveryAddress.ID, secondaryDeliveryAddress.ID)
		suite.Equal(updatedMTOShipment.SecondaryDeliveryAddress.StreetAddress1, *secondaryDeliveryAddress.StreetAddress1)
		suite.Equal(updatedMTOShipment.PrimeActualWeight, *&actualWeight)
	})

	now := time.Now()

	suite.T().Run("Failed case if not both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		eightDaysFromNow := now.AddDate(0, 0, 8)
		threeDaysBefore := now.AddDate(0, 0, -3)
		oldMTOShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &eightDaysFromNow,
				ApprovedDate:        &threeDaysBefore,
			},
		})
		payload3 := payload
		payload3.ID = strfmt.UUID(oldMTOShipment2.ID.String())
		unmodifiedSince := oldMTOShipment2.UpdatedAt
		payload3.PrimeEstimatedWeight = 4500

		//remove this when remove required fields
		scheduledPickupDate = strfmt.Date(eightDaysFromNow)
		payload3.ScheduledPickupDate = scheduledPickupDate

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload3,
			MoveTaskOrderID:   strfmt.UUID(oldMTOShipment2.MoveTaskOrderID.String()),
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		_, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.Error(err)
	})

	suite.T().Run("Successful case if both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		tenDaysFromNow := now.AddDate(0, 0, 11)
		oldMTOShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &tenDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		payload3 := payload
		payload3.ID = strfmt.UUID(oldMTOShipment2.ID.String())
		unmodifiedSince := oldMTOShipment2.UpdatedAt
		payload3.PrimeEstimatedWeight = 4500

		//remove this when remove required fields
		scheduledPickupDate = strfmt.Date(tenDaysFromNow)
		payload3.ScheduledPickupDate = scheduledPickupDate

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload3,
			MoveTaskOrderID:   strfmt.UUID(oldMTOShipment2.MoveTaskOrderID.String()),
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	})

	suite.T().Run("Failed case if approved date is 3-9 days from scheduled move date but estimated weight recorded date isn't at least 3 days prior to scheduled move date", func(t *testing.T) {
		twoDaysFromNow := now.AddDate(0, 0, 2)
		twoDaysBefore := now.AddDate(0, 0, -2)
		oldMTOShipment4 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &twoDaysFromNow,
				ApprovedDate:        &twoDaysBefore,
			},
		})
		payload4 := payload
		payload4.ID = strfmt.UUID(oldMTOShipment4.ID.String())
		unmodifiedSince := oldMTOShipment4.UpdatedAt
		payload4.PrimeEstimatedWeight = 4500

		//remove this when remove required fields
		scheduledPickupDate = strfmt.Date(twoDaysFromNow)
		payload4.ScheduledPickupDate = scheduledPickupDate

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload4,
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		_, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.Error(err)
	})

	suite.T().Run("Successful case if approved date is 3-9 days from scheduled move date and estimated weight recorded date is at least 3 days prior to scheduled move date", func(t *testing.T) {
		sixDaysFromNow := now.AddDate(0, 0, 6)
		twoDaysBefore := now.AddDate(0, 0, -2)
		oldMTOShipment3 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &sixDaysFromNow,
				ApprovedDate:        &twoDaysBefore,
			},
		})
		payload4 := payload
		payload4.ID = strfmt.UUID(oldMTOShipment3.ID.String())
		unmodifiedSince := oldMTOShipment3.UpdatedAt
		payload4.PrimeEstimatedWeight = 4500

		//remove this when remove required fields
		scheduledPickupDate = strfmt.Date(sixDaysFromNow)
		payload4.ScheduledPickupDate = scheduledPickupDate

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload4,
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	})

	suite.T().Run("Failed case if approved date is less than 3 days from scheduled move date but estimated weight recorded date isn't at least 1 day prior to scheduled move date", func(t *testing.T) {
		oneDayFromNow := now.AddDate(0, 0, 1)
		oneDayBefore := now.AddDate(0, 0, -1)
		oldMTOShipment4 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &oneDayFromNow,
				ApprovedDate:        &oneDayBefore,
			},
		})
		payload5 := payload
		payload5.ID = strfmt.UUID(oldMTOShipment4.ID.String())
		unmodifiedSince := oldMTOShipment4.UpdatedAt
		payload5.PrimeEstimatedWeight = 4500

		//remove this when remove required fields
		scheduledPickupDate = strfmt.Date(oneDayFromNow)
		payload5.ScheduledPickupDate = scheduledPickupDate

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload5,
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		_, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.Error(err)
	})

	suite.T().Run("Successful case if approved date is less than 3 days from scheduled move date and estimated weight recorded date is at least 1 day prior to scheduled move date", func(t *testing.T) {
		twoDaysFromNow := now.AddDate(0, 0, 2)
		oldMTOShipment4 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &twoDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		payload5 := payload
		payload5.ID = strfmt.UUID(oldMTOShipment4.ID.String())
		unmodifiedSince := oldMTOShipment4.UpdatedAt
		payload5.PrimeEstimatedWeight = 4500

		//remove this when remove required fields
		scheduledPickupDate = strfmt.Date(twoDaysFromNow)
		payload5.ScheduledPickupDate = scheduledPickupDate

		params := mtoshipmentops.UpdateMTOShipmentParams{
			Body:              &payload5,
			IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		}
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(params)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	})
}
