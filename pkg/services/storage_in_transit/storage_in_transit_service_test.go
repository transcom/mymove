package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"

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

func payloadForStorageInTransitModel(s *models.StorageInTransit) *apimessages.StorageInTransit {
	if s == nil {
		return nil
	}

	location := string(s.Location)
	status := string(s.Status)

	return &apimessages.StorageInTransit{
		ID:                  *handlers.FmtUUID(s.ID),
		ShipmentID:          *handlers.FmtUUID(s.ShipmentID),
		EstimatedStartDate:  handlers.FmtDate(s.EstimatedStartDate),
		Notes:               handlers.FmtStringPtr(s.Notes),
		WarehouseAddress:    payloadForAddressModel(&s.WarehouseAddress),
		WarehouseEmail:      handlers.FmtStringPtr(s.WarehouseEmail),
		WarehouseID:         handlers.FmtString(s.WarehouseID),
		WarehouseName:       handlers.FmtString(s.WarehouseName),
		WarehousePhone:      handlers.FmtStringPtr(s.WarehousePhone),
		Location:            &location,
		Status:              *handlers.FmtString(status),
		AuthorizationNotes:  handlers.FmtStringPtr(s.AuthorizationNotes),
		AuthorizedStartDate: handlers.FmtDatePtr(s.AuthorizedStartDate),
		ActualStartDate:     handlers.FmtDatePtr(s.ActualStartDate),
		OutDate:             handlers.FmtDatePtr(s.OutDate),
	}
}
func TestStorageInTransitServiceSuite(t *testing.T) {
	hs := &StorageInTransitServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, hs)
}
