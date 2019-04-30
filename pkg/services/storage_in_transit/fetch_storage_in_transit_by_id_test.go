package storageintransit

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestFetchStorageInTransitByID() {
	shipment, sit, user := setupStorageInTransitServiceTest(suite)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}

	fetcher := NewStorageInTransitByIDFetcher(suite.DB())
	actualSIT, err := fetcher.FetchStorageInTransitByID(sit.ID, shipment.ID, &session)
	suite.NoError(err, "Error fetching SIT")

	storageInTransitCompare(suite, sit, *actualSIT)

	// Let's make sure it fails when a TSP who doesn't own the shipment tries to do a GET on this
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	session = auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		TspUserID:       user.ID,
	}
	actualSIT, err = fetcher.FetchStorageInTransitByID(sit.ID, shipment.ID, &session)
	suite.Error(models.ErrFetchForbidden)
}
