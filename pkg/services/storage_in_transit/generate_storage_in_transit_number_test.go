package storageintransit

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestGenerateStorageInTransitNumber() {
	var storageInTransitNumberTestCases = []struct {
		name                           string
		placeInSitTime                 time.Time
		startSequenceNumber            int // <= 0 indicates to not reset.
		expectedStorageInTransitNumber string
	}{
		{
			"first SIT number for a day",
			time.Date(2019, 4, 22, 0, 0, 0, 0, time.UTC),
			0,
			"191120001",
		},
		{
			"second SIT number for a day",
			time.Date(2019, 4, 22, 0, 0, 0, 0, time.UTC),
			0,
			"191120002",
		},
		{
			"same year, different day",
			time.Date(2019, 4, 23, 0, 0, 0, 0, time.UTC),
			0,
			"191130001",
		},
		{
			"different year, same day",
			time.Date(2018, 4, 22, 0, 0, 0, 0, time.UTC),
			0,
			"181120001",
		},
		{
			"max 4-digit sequence number",
			time.Date(2019, 4, 24, 0, 0, 0, 0, time.UTC),
			9999,
			"1911410000",
		},
	}

	storageInTransitNumberGenerator := NewStorageInTransitNumberGenerator(suite.DB())

	for _, testCase := range storageInTransitNumberTestCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			// Reset sequence number if needed.
			if testCase.startSequenceNumber > 0 {
				utcTime := testCase.placeInSitTime.UTC()
				err := testdatagen.SetStorageInTransitSequenceNumber(suite.DB(), utcTime.Year(), utcTime.YearDay(), testCase.startSequenceNumber)
				suite.NoError(err)
			}

			storageInTransitNumber, err := storageInTransitNumberGenerator.GenerateStorageInTransitNumber(testCase.placeInSitTime)
			suite.NoError(err)

			suite.Equal(testCase.expectedStorageInTransitNumber, storageInTransitNumber)
		})
	}
}
