package storageintransit

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestPatchStorageInTransit() {
	shipment, sit, user := setupStorageInTransitServiceTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}

	sit.WarehouseID = "123456"
	sit.Notes = swag.String("Updated Note")
	sit.WarehouseEmail = swag.String("updated@email.com")
	payload := payloadForStorageInTransitModel(&sit)

	patcher := NewStorageInTransitPatcher(suite.DB())

	// Office happy path
	actualStorageInTransit, verrs, err := patcher.PatchStorageInTransit(*payload, shipment.ID, sit.ID, &session)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	storageInTransitCompare(suite, *actualStorageInTransit, sit)

	// Fail with service member user
	session = auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		ServiceMemberID: user.ID,
	}

	_, _, err = patcher.PatchStorageInTransit(*payload, shipment.ID, sit.ID, &session)
	suite.Error(err, "FETCH_FORBIDDEN")

	// Fail when tsp doesn't 'own' shipment
	session = auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser.ID,
	}

	_, _, err = patcher.PatchStorageInTransit(*payload, shipment.ID, sit.ID, &session)
	suite.Error(err, "FETCH_FORBIDDEN")

	// TSP Happy path
	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// change the status to in_sit.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)
	actualStorageInTransit, verrs, err = patcher.PatchStorageInTransit(*payload, shipment.ID, sit.ID, &session)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	storageInTransitCompare(suite, *actualStorageInTransit, sit)
}
