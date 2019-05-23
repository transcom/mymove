package storageintransit

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestDeliverStorageInTransits() {
	shipment, sit, _ := setupStorageInTransitServiceTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	deliverer := NewStorageInTransitsDeliverer(suite.DB())
	sit.Status = models.StorageInTransitStatusINSIT
	_, _ = suite.DB().ValidateAndSave(&sit)

	shipment.ActualDeliveryDate = &testdatagen.DateInsidePeakRateCycle
	_, _ = suite.DB().ValidateAndSave(&shipment)

	// Happy path
	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// change the status to in_sit.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)
	actualStorageInTransit, verrs, err := deliverer.DeliverStorageInTransits(
		shipment.ID,
		tspUser.TransportationServiceProviderID)

	suite.NoError(err)
	suite.Equal(false, verrs.HasAny())
	suite.Equal(models.StorageInTransitStatusDELIVERED, actualStorageInTransit[0].Status)
	suite.NotNil(actualStorageInTransit[0].OutDate)

	updatedShipment, err := models.FetchShipmentByTSP(suite.DB(), tspUser.TransportationServiceProviderID, shipment.ID)
	suite.Equal(updatedShipment.ActualDeliveryDate, actualStorageInTransit[0].OutDate)
}
