package mtoshipment

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestMTOShipmentUpdater() {
	oldMTOShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
	mtoShipmentUpdater := NewMTOShipmentUpdater(suite.DB())

	requestedPickupDate := strfmt.Date(*oldMTOShipment.RequestedPickupDate)
	scheduledPickupDate := strfmt.Date(time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC))
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
		ID:                       strfmt.UUID(oldMTOShipment.ID.String()),
		DestinationAddress:       &destinationAddress,
		PickupAddress:            &pickupAddress,
		RequestedPickupDate:      &requestedPickupDate,
		ScheduledPickupDate:      &scheduledPickupDate,
		ShipmentType:             "INTERNATIONAL_UB",
		SecondaryPickupAddress:   &secondaryPickupAddress,
		SecondaryDeliveryAddress: &secondaryDeliveryAddress,
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
		unmodifiedSince := oldMTOShipment.UpdatedAt
		fmt.Println(unmodifiedSince)

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
	})
}