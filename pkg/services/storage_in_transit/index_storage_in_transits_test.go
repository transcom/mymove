package storageintransit

import (
	"github.com/transcom/mymove/pkg/auth"
)

func (suite *StorageInTransitServiceSuite) TestIndexStorageInTransits() {
	shipment, _, user := setupStorageInTransitServiceTest(suite)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}
	// Happy path. This should succeed.
	indexer := NewStorageInTransitIndexer(suite.DB())
	actualStorageInTransits, err := indexer.IndexStorageInTransits(shipment.ID, &session)
	suite.NoError(err)
	suite.Equal(2, len(actualStorageInTransits))

	// Let's make sure this fails for a servicemember user
	session = auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		ServiceMemberID: user.ID,
	}

	_, err = indexer.IndexStorageInTransits(shipment.ID, &session)
	suite.Error(err, "FETCH_FORBIDDEN")
}
