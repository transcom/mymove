package storageintransit

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestDeliverStorageInTransits() {
	shipment, sit, _ := setupStorageInTransitServiceTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser.ID,
	}

	deliverer := NewStorageInTransitsDeliverer(suite.DB())
	sit.Status = models.StorageInTransitStatusINSIT
	_, _ = suite.DB().ValidateAndSave(&sit)

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

	actualStorageInTransit, verrs, err := deliverer.DeliverStorageInTransits(shipment.ID, &session)
	suite.NoError(err)
	suite.Equal(false, verrs.HasAny())
	suite.Equal(models.StorageInTransitStatusDELIVERED, actualStorageInTransit[0].Status)
	suite.Equal(shipment.ActualDeliveryDate, actualStorageInTransit[0].OutDate)

	// It should not work if we're coming back from released status
	sit.Status = models.StorageInTransitStatusRELEASED
	_, _ = suite.DB().ValidateAndSave(&sit)
}
