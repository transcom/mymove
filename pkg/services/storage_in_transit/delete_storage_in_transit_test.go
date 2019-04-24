package storageintransit

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestDeleteStorageInTransit() {
	shipment, sit, user := setupStorageInTransitServiceTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}
	// Office user can't delete. This should fail.
	deleter := NewStorageInTransitDeleter(suite.DB())
	err := deleter.DeleteStorageInTransit(shipment.ID, sit.ID, &session)
	suite.Error(err, "FETCH_FORBIDDEN")

	// If a TSP doesn't 'own' the storage in transit, it should fail.
	session = auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser.ID,
	}
	deleter = NewStorageInTransitDeleter(suite.DB())
	err = deleter.DeleteStorageInTransit(shipment.ID, sit.ID, &session)
	suite.Error(err, "FETCH_FORBIDDEN")

	// Happy path
	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// Use these to create a SIT for them.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)

	err = deleter.DeleteStorageInTransit(shipment.ID, sit.ID, &session)
	suite.NoError(err)
	deletedStorageInTransit, err := models.FetchStorageInTransitByID(suite.DB(), sit.ID)
	suite.Nil(deletedStorageInTransit)

}
