package mtoshipment

import (
	"github.com/go-openapi/strfmt"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/testdatagen"
	"time"
)

func (suite *MTOShipmentServiceSuite) TestMTOShipmentFetcher() {
	expectedMTOShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
	mtoShipmentFetcher := NewMTOShipmentFetcher(suite.DB())

	actualMTOShipment, err := mtoShipmentFetcher.FetchMTOShipment(expectedMTOShipment.ID)
	suite.NoError(err)

	suite.NotZero(expectedMTOShipment.ID, actualMTOShipment.ID)
	suite.Equal(expectedMTOShipment.MoveTaskOrder.ID, actualMTOShipment.MoveTaskOrder.ID)
	suite.NotNil(actualMTOShipment.MoveTaskOrder)
	suite.NotNil(expectedMTOShipment.ShipmentType)
	suite.NotNil(expectedMTOShipment.PickupAddress)
	suite.NotNil(expectedMTOShipment.DestinationAddress)
}

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

	payload := primemessages.MTOShipment{
		ID: strfmt.UUID(oldMTOShipment.ID.String()),
		DestinationAddress:      &destinationAddress,
		PickupAddress:           &pickupAddress,
		RequestedPickupDate:      &requestedPickupDate,
		ScheduledPickupDate:      &scheduledPickupDate,
		ShipmentType:             "INTERNATIONAL_UB",
	}

	updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC), &payload)
	//fmt.Println(updatedMTOShipment)
	suite.NoError(err)

	suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
	suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
}