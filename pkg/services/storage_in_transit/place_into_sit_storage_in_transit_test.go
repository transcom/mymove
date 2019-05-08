package storageintransit

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestPlaceIntoSITStorageInTransit() {
	shipment, sit, user := setupStorageInTransitServiceTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser.ID,
	}

	sit.AuthorizedStartDate = &testdatagen.DateInsidePerformancePeriod
	suite.DB().Save(&sit)
	payload := apimessages.StorageInTransitInSitPayload{
		ActualStartDate: *handlers.FmtDate(testdatagen.DateInsidePerformancePeriod),
	}

	inSITPlacer := NewStorageInTransitInSITPlacer(suite.DB())

	// Happy path
	sit.Status = models.StorageInTransitStatusAPPROVED
	_, _ = suite.DB().ValidateAndSave(&sit)
	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// change the status to in_sit.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)

	actualStorageInTransit, verrs, err := inSITPlacer.PlaceIntoSITStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.Equal(models.StorageInTransitStatusINSIT, actualStorageInTransit.Status)

	// Shouldn't work with an office user
	session = auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}

	_, _, err = inSITPlacer.PlaceIntoSITStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.Error(err, "FETCH_FORBIDDEN")

	// Shouldn't work if status is not approved
	sit.Status = models.StorageInTransitStatusREQUESTED
	_, _ = suite.DB().ValidateAndSave(&sit)

	_, _, err = inSITPlacer.PlaceIntoSITStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.Error(err, "FETCH_FORBIDDEN")
}
