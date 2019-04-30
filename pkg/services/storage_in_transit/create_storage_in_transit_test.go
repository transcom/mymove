package storageintransit

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestCreateStorageInTransit() {
	shipment, sit, user := setupStorageInTransitServiceTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser.ID,
	}

	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// Use these to create a SIT for them.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)

	payload := payloadForStorageInTransitModel(&sit)

	// Happy path. This should succeed.
	creator := NewStorageInTransitCreator(suite.DB())
	actualStorageInTransit, verrs, err := creator.CreateStorageInTransit(*payload, shipment.ID, &session)

	suite.False(verrs.HasAny())
	suite.NoError(err)
	storageInTransitCompare(suite, *actualStorageInTransit, sit)

	// This should fail with a forbidden error if we try with an office user.
	session = auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}
	// Happy path. This should succeed.
	actualStorageInTransit, verrs, err = creator.CreateStorageInTransit(*payload, shipment.ID, &session)
	suite.Error(err, "FETCH_FORBIDDEN")

	// This should fail if the tsp user doesn't 'own' the shipment associated with the SIT.
	assertions = testdatagen.Assertions{
		TspUser: models.TspUser{
			Email: "unused_test_email_for_sit@sit.com",
		},
	}
	tspUser2 := testdatagen.MakeTspUser(suite.DB(), assertions)
	session = auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser2.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser2.ID,
	}

	actualStorageInTransit, verrs, err = creator.CreateStorageInTransit(*payload, shipment.ID, &session)
	suite.Error(err, "USER_UNAUTHORIZED")
}
