package movetaskorder

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {
	expectedMTO := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	expectedEntitlement := testdatagen.MakeEntitlement(suite.DB(), testdatagen.Assertions{
		GHCEntitlement: models.GHCEntitlement{
			MoveTaskOrderID: expectedMTO.ID,
		},
	})
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{
		ServiceItem: models.ServiceItem{
			MoveTaskOrderID: expectedMTO.ID,
		},
	})
	mtoFetcher := NewMoveTaskOrderFetcher(suite.DB())

	actualMTO, err := mtoFetcher.FetchMoveTaskOrder(expectedMTO.ID)

	suite.NoError(err)
	suite.NotZero(actualMTO.Customer)
	suite.Equal(expectedMTO.CustomerID, actualMTO.CustomerID)
	suite.Equal(expectedMTO.CustomerRemarks, actualMTO.CustomerRemarks)
	suite.Equal(expectedMTO.DestinationAddressID, actualMTO.DestinationAddressID)
	suite.NotZero(actualMTO.DestinationAddress)
	suite.Equal(expectedMTO.DestinationDutyStationID, actualMTO.DestinationDutyStationID)
	suite.NotZero(actualMTO.DestinationDutyStation)
	suite.NotZero(expectedEntitlement.ID, actualMTO.Entitlements.ID)
	suite.Equal(expectedMTO.MoveID, actualMTO.MoveID)
	suite.NotZero(actualMTO.Move)
	suite.Equal(expectedMTO.OriginDutyStationID, actualMTO.OriginDutyStationID)
	suite.NotZero(actualMTO.OriginDutyStation)
	suite.Equal(expectedMTO.PickupAddressID, actualMTO.PickupAddressID)
	suite.NotZero(actualMTO.PickupAddress)
	suite.Equal(expectedMTO.RequestedPickupDate.UTC(), actualMTO.RequestedPickupDate.UTC())
	suite.Len(actualMTO.ServiceItems, 1)
	suite.Equal(serviceItem.ID, actualMTO.ServiceItems[0].ID)
	suite.Equal(expectedMTO.Status, actualMTO.Status)

}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderStatusUpdater() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	originalMTO := serviceItem.MoveTaskOrder
	// check not equal to what asserting against below
	suite.NotEqual(originalMTO.Status, models.MoveTaskOrderStatusDraft)
	mtoStatusUpdater := NewMoveTaskOrderStatusUpdater(suite.DB())

	updatedMTO, err := mtoStatusUpdater.UpdateMoveTaskOrderStatus(originalMTO.ID, models.MoveTaskOrderStatusDraft)

	suite.NoError(err)
	suite.Equal(models.MoveTaskOrderStatusDraft, updatedMTO.Status)
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderStatusUpdaterEmptyStatus() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	originalMTO := serviceItem.MoveTaskOrder
	// check not equal to what asserting against below
	suite.NotEqual(originalMTO.Status, models.MoveTaskOrderStatusDraft)
	mtoStatusUpdater := NewMoveTaskOrderStatusUpdater(suite.DB())

	_, err := mtoStatusUpdater.UpdateMoveTaskOrderStatus(originalMTO.ID, "")

	suite.Error(err)
}
