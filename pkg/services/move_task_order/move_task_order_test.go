package movetaskorder

import (
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	expectedMTO := serviceItem.MoveTaskOrder
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())
	actualMTO, err := mtoFetcher.FetchMoveTaskOrder(expectedMTO.ID)

	suite.NoError(err)
	suite.Equal(expectedMTO.CustomerID, actualMTO.CustomerID)
	suite.NotZero(actualMTO.Customer)
	suite.Equal(expectedMTO.DestinationAddressID, actualMTO.DestinationAddressID)
	suite.NotZero(actualMTO.DestinationAddress)
	suite.Equal(expectedMTO.DestinationDutyStationID, actualMTO.DestinationDutyStationID)
	suite.NotZero(actualMTO.DestinationDutyStation)
	suite.Equal(expectedMTO.MoveID, actualMTO.MoveID)
	suite.NotZero(actualMTO.Move)
	suite.Equal(expectedMTO.NTSEntitlement, actualMTO.NTSEntitlement)
	suite.Equal(expectedMTO.OriginDutyStationID, actualMTO.OriginDutyStationID)
	suite.NotZero(actualMTO.OriginDutyStation)
	suite.Equal(expectedMTO.POVEntitlement, actualMTO.POVEntitlement)
	suite.Equal(expectedMTO.PickupAddressID, actualMTO.PickupAddressID)
	suite.NotZero(actualMTO.PickupAddress)
	suite.Equal(expectedMTO.RequestedPickupDates.UTC(), actualMTO.RequestedPickupDates.UTC())
	suite.Equal(expectedMTO.SitEntitlement, actualMTO.SitEntitlement)
	suite.Len(actualMTO.ServiceItems, 1)
	suite.Equal(serviceItem.ID, actualMTO.ServiceItems[0].ID)
	suite.Equal(expectedMTO.Status, actualMTO.Status)
	suite.Equal(expectedMTO.WeightEntitlement, actualMTO.WeightEntitlement)

}
