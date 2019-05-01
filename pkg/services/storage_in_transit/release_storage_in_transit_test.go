package storageintransit

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestReleaseStorageInTransit() {
	shipment, sit, user := setupStorageInTransitServiceTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser.ID,
	}
	payload := apimessages.StorageInTransitReleasePayload{
		ReleasedOn: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
	}

	releaser := NewStorageInTransitInReleaser(suite.DB())
	sit.Status = models.StorageInTransitStatusINSIT
	_, _ = suite.DB().ValidateAndSave(&sit)

	// Should fail if TSP doesn't 'own' the storage in transit
	_, _, err := releaser.ReleaseStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.Error(err, "WRITE_CONFLICT")

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

	actualStorageInTransit, verrs, err := releaser.ReleaseStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.NoError(err)
	suite.Equal(false, verrs.HasAny())
	suite.Equal(models.StorageInTransitStatusRELEASED, actualStorageInTransit.Status)

	// It should also work if we're coming back from delivered status
	sit.Status = models.StorageInTransitStatusDELIVERED
	_, _ = suite.DB().ValidateAndSave(&sit)

	actualStorageInTransit, verrs, err = releaser.ReleaseStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.NoError(err)
	suite.Equal(false, verrs.HasAny())
	suite.Equal(models.StorageInTransitStatusRELEASED, actualStorageInTransit.Status)

	// It should fail if we're skipping a step (going from requested to released)
	sit.Status = models.StorageInTransitStatusREQUESTED
	_, _ = suite.DB().ValidateAndSave(&sit)

	actualStorageInTransit, verrs, err = releaser.ReleaseStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.Error(err, "WRITE_CONFLICT")

	// Should fail for an office user
	session = auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}

	sit.Status = models.StorageInTransitStatusINSIT
	_, _ = suite.DB().ValidateAndSave(&sit)

	actualStorageInTransit, verrs, err = releaser.ReleaseStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.Error(err, "FETCH_FORBIDDEN")

}
