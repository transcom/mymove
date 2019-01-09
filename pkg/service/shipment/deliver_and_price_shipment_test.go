package shipment

import (
	"testing"

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
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.FatalNoError(err)

	shipment := shipments[0]

	// And an unpriced, approved pre-approval
	testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
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
	engine := rateengine.NewRateEngine(suite.DB(), suite.logger, planner)
	verrs, err := DeliverAndPriceShipment{
		DB:     suite.DB(),
		Engine: engine,
	}.Call(deliveryDate, &shipment)

	suite.FatalNoError(err)
	suite.FatalFalse(verrs.HasAny())

	suite.Equal(shipment.Status, models.ShipmentStatusDELIVERED)

	fetchedLineItems, err := models.FetchLineItemsByShipmentID(suite.DB(), &shipment.ID)
	suite.FatalNoError(err)
	// All items should be priced
	for _, item := range fetchedLineItems {
		suite.NotNil(item.AmountCents, item.Tariff400ngItem.Code)
	}
}

type DeliverPriceShipmentSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *DeliverPriceShipmentSuite) SetupTest() {
	suite.DB().TruncateAll()
}
func TestUpdateInvoiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &DeliverPriceShipmentSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}
