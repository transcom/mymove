package shipment

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

func (suite *DeliverPriceShipmentSuite) TestUpdateInvoicesCall() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.FatalNoError(err)

	shipment := shipments[0]

	// And an unpriced, approved pre-approval
	testdatagen.MakeCompleteShipmentLineItem(suite.db, testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
			Status:     models.ShipmentLineItemStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	deliveryDate := testdatagen.DateInsidePerformancePeriod
	planner := route.NewTestingPlanner(1100)
	engine := rateengine.NewRateEngine(suite.db, suite.logger, planner)
	verrs, err := DeliverAndPriceShipment{
		DB:     suite.db,
		Engine: engine,
	}.Call(deliveryDate, &shipment)

	suite.FatalNoError(err)
	suite.FatalFalse(verrs.HasAny())

	suite.Equal(shipment.Status, models.ShipmentStatusDELIVERED)

	fetchedLineItems, err := models.FetchLineItemsByShipmentID(suite.db, &shipment.ID)
	suite.FatalNoError(err)
	// All items should be priced
	for _, item := range fetchedLineItems {
		suite.NotNil(item.AmountCents, item.Tariff400ngItem.Code)
	}
}

type DeliverPriceShipmentSuite struct {
	testingsuite.LocalTestSuite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *DeliverPriceShipmentSuite) SetupTest() {
	suite.db.TruncateAll()
}
func TestUpdateInvoiceSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &DeliverPriceShipmentSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
