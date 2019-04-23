package storageintransit

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type StorageInTransitServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *StorageInTransitServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestStorageInTransitSuite(t *testing.T) {
	hs := &StorageInTransitServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
}

func setupStorageInTransitServiceTest(suite *StorageInTransitServiceSuite) (shipment models.Shipment, sit models.StorageInTransit, user models.OfficeUser) {
	shipment = testdatagen.MakeDefaultShipment(suite.DB())
	user = testdatagen.MakeDefaultOfficeUser(suite.DB())

	assertions := testdatagen.Assertions{
		StorageInTransit: models.StorageInTransit{
			Location:           models.StorageInTransitLocationORIGIN,
			ShipmentID:         shipment.ID,
			EstimatedStartDate: testdatagen.DateInsidePeakRateCycle,
		},
	}
	testdatagen.MakeStorageInTransit(suite.DB(), assertions)
	sit = testdatagen.MakeStorageInTransit(suite.DB(), assertions)

	return shipment, sit, user
}

func storageInTransitCompare(suite *StorageInTransitServiceSuite, expected models.StorageInTransit, actual models.StorageInTransit) {
	suite.Equal(expected.WarehouseEmail, actual.WarehouseEmail)
	suite.Equal(expected.Notes, actual.Notes)
	suite.Equal(expected.WarehouseID, actual.WarehouseID)
	suite.Equal(expected.Location, actual.Location)
	suite.Equal(expected.WarehouseName, actual.WarehouseName)
	suite.Equal(expected.WarehousePhone, actual.WarehousePhone)
	suite.True(expected.EstimatedStartDate.Equal(actual.EstimatedStartDate))
	suite.Equal(expected.Status, actual.Status)
}
