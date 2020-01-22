package mtoshipment

import (
	"github.com/transcom/mymove/pkg/testdatagen"
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
