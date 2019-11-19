package movetaskorder

import (
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
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

	suite.Nil(originalMTO.AvailableToPrimeDate)
	// check not equal to what asserting against below
	suite.NotEqual(originalMTO.Status, models.MoveTaskOrderStatusApproved)
	mtoStatusUpdater := NewMoveTaskOrderStatusUpdater(suite.DB())

	updatedMTO, err := mtoStatusUpdater.UpdateMoveTaskOrderStatus(originalMTO.ID, models.MoveTaskOrderStatusApproved)

	suite.NoError(err)
	suite.Equal(models.MoveTaskOrderStatusApproved, updatedMTO.Status)
	// date should be filled when mto has been approved
	suite.NotNil(updatedMTO.AvailableToPrimeDate)
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderStatusUpdaterEmptyStatus() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	originalMTO := serviceItem.MoveTaskOrder
	// check not equal to what asserting against below
	suite.NotEqual(originalMTO.Status, models.MoveTaskOrderStatusSubmitted)
	mtoStatusUpdater := NewMoveTaskOrderStatusUpdater(suite.DB())

	_, err := mtoStatusUpdater.UpdateMoveTaskOrderStatus(originalMTO.ID, "")

	suite.Error(err)
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderActualWeightUpdater() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	originalMTO := serviceItem.MoveTaskOrder
	// check not equal to what asserting against below
	suite.Nil(originalMTO.PrimeActualWeight)
	mtoActualWeightUpdater := NewMoveTaskOrderActualWeightUpdater(suite.DB())

	newWeight := int64(566)
	updatedMTO, err := mtoActualWeightUpdater.UpdateMoveTaskOrderActualWeight(originalMTO.ID, newWeight)

	suite.NoError(err)
	suite.Equal(unit.Pound(newWeight), *updatedMTO.PrimeActualWeight)
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderPrimeEstimatedWeightUpdater() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	originalMTO := serviceItem.MoveTaskOrder
	// check not equal to what asserting against below
	suite.Nil(originalMTO.PrimeEstimatedWeight)
	suite.Nil(originalMTO.PrimeEstimatedWeightRecordedDate)
	mtoActualWeightUpdater := NewMoveTaskOrderEstimatedWeightUpdater(suite.DB())
	mtoActualWeightFetcher := NewMoveTaskOrderFetcher(suite.DB())

	newWeight := unit.Pound(1234)
	now := time.Now()
	updatedMTO, updateErr := mtoActualWeightUpdater.UpdatePrimeEstimatedWeight(originalMTO.ID, newWeight, now)
	suite.NoError(updateErr)
	suite.NotNil(updatedMTO)
	dbUpdatedMTO, fetchErr := mtoActualWeightFetcher.FetchMoveTaskOrder(updatedMTO.ID)
	suite.NoError(fetchErr)

	suite.Equal(newWeight, *updatedMTO.PrimeEstimatedWeight)
	suite.Equal(now.Format(time.RFC3339), updatedMTO.PrimeEstimatedWeightRecordedDate.Format(time.RFC3339))
	suite.Equal(newWeight, *dbUpdatedMTO.PrimeEstimatedWeight)
	suite.Equal(now.Format(time.RFC3339), dbUpdatedMTO.PrimeEstimatedWeightRecordedDate.Format(time.RFC3339))
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderPrimeEstimatedWeightUpdaterInvalidWeight() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	originalMTO := serviceItem.MoveTaskOrder
	// check not equal to what asserting against below
	suite.Nil(originalMTO.PrimeEstimatedWeight)
	suite.Nil(originalMTO.PrimeEstimatedWeightRecordedDate)
	mtoActualWeightUpdater := NewMoveTaskOrderEstimatedWeightUpdater(suite.DB())

	newWeight := unit.Pound(-1000)
	now := time.Now()
	_, updateErr := mtoActualWeightUpdater.UpdatePrimeEstimatedWeight(originalMTO.ID, newWeight, now)

	suite.Error(updateErr)
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderDestinationAddressUpdater() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	originalMTO := serviceItem.MoveTaskOrder
	// check not equal to what asserting against below
	address := testdatagen.MakeDefaultAddress(suite.DB())
	mtoActualWeightUpdater := NewMoveTaskOrderDestinationAddressUpdater(suite.DB())
	moveTaskOrderFetcher := NewMoveTaskOrderFetcher(suite.DB())

	updatedMTO, updateErr := mtoActualWeightUpdater.UpdateMoveTaskOrderDestinationAddress(originalMTO.ID, &address)
	suite.NoError(updateErr)
	suite.NotNil(updatedMTO)
	// CreatedAt, UpdatedAt will be different so just assert against string format
	suite.Equal(address.LineFormat(), updatedMTO.DestinationAddress.LineFormat())

	dbUpdatedMTO, fetchErr := moveTaskOrderFetcher.FetchMoveTaskOrder(updatedMTO.ID)
	suite.NoError(fetchErr)
	suite.Equal(address.LineFormat(), dbUpdatedMTO.DestinationAddress.LineFormat())
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderPrimePostCounselingUpdater() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	originalMTO := serviceItem.MoveTaskOrder
	suite.Nil(originalMTO.SubmittedCounselingInfoDate)

	// check not equal to what asserting against below
	address := testdatagen.MakeDefaultAddress(suite.DB())
	address2 := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
	mtoPostCounselingInformationUpdater := NewMoveTaskOrderPostCounselingInformationUpdater(suite.DB())
	moveTaskOrderFetcher := NewMoveTaskOrderFetcher(suite.DB())

	now := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	information := services.PostCounselingInformation{
		PPMIsIncluded:            true,
		ScheduledMoveDate:        now,
		SecondaryDeliveryAddress: &address,
		SecondaryPickupAddress:   &address2,
	}
	updatedMTO, updateErr := mtoPostCounselingInformationUpdater.UpdateMoveTaskOrderPostCounselingInformation(originalMTO.ID, information)
	suite.NoError(updateErr)
	suite.NotNil(updatedMTO)
	suite.Equal(information.ScheduledMoveDate, *updatedMTO.ScheduledMoveDate)
	suite.Equal(information.SecondaryDeliveryAddress, updatedMTO.SecondaryDeliveryAddress)
	suite.Equal(information.SecondaryPickupAddress, updatedMTO.SecondaryPickupAddress)
	suite.Equal(information.PPMIsIncluded, *updatedMTO.PpmIsIncluded)
	suite.NotNil(updatedMTO.SubmittedCounselingInfoDate)

	dbUpdatedMTO, fetchErr := moveTaskOrderFetcher.FetchMoveTaskOrder(updatedMTO.ID)
	suite.NoError(fetchErr)
	suite.Equal(information.ScheduledMoveDate.String(), (*dbUpdatedMTO.ScheduledMoveDate).UTC().String())
	suite.Equal(information.SecondaryDeliveryAddress.LineFormat(), dbUpdatedMTO.SecondaryDeliveryAddress.LineFormat())
	suite.Equal(information.SecondaryPickupAddress.LineFormat(), dbUpdatedMTO.SecondaryPickupAddress.LineFormat())
	suite.Equal(information.PPMIsIncluded, *dbUpdatedMTO.PpmIsIncluded)
}
