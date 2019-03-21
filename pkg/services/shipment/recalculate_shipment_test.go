package shipment

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RecalculateShipmentSuite) helperDeliverAndPriceShipment() *models.Shipment {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
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

	// Make sure there's a FuelEIADieselPrice
	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	deliveryDate := testdatagen.DateInsidePerformancePeriod
	planner := route.NewTestingPlanner(1100)
	engine := rateengine.NewRateEngine(suite.DB(), suite.logger)
	verrs, err := DeliverAndPriceShipment{
		DB:      suite.DB(),
		Engine:  engine,
		Planner: planner,
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

	return &shipment
}

func (suite *RecalculateShipmentSuite) TestRecalculateShipmentCall() {

	shipment := suite.helperDeliverAndPriceShipment()
	shipmentID := shipment.ID

	fetchedLineItems, err := models.FetchLineItemsByShipmentID(suite.DB(), &shipmentID)
	suite.FatalNoError(err)
	// All items should be priced
	for _, item := range fetchedLineItems {
		suite.NotNil(item.AmountCents, item.Tariff400ngItem.Code)
	}

	// Verify that all base shipment line items are present
	allPresent := ProcessRecalculateShipment{}.hasAllBaseLineItems(fetchedLineItems)
	suite.Equal(true, allPresent)

	// Remove 1 base shipment line time
	// Removing Fuel Surcharge
	var fuelSurcharge models.ShipmentLineItem
	removedFuelSurcharge := false
	for _, item := range fetchedLineItems {
		if item.Tariff400ngItem.Code == "16A" {
			fuelSurcharge = item
			err := suite.DB().Destroy(&fuelSurcharge)
			if err != nil {
				suite.logger.Fatal("Error Removing Fuel Surcharge")
			}
			removedFuelSurcharge = true
		}
	}
	suite.Equal(true, removedFuelSurcharge)

	zeroCents := unit.Cents(0)
	zeroMillicents := unit.Millicents(0)

	// Set price of 1 base shipment line item to zero
	updatedUnpack := false
	for _, item := range fetchedLineItems {
		if item.Tariff400ngItem.Code == "105C" {
			item.AmountCents = &zeroCents
			item.AppliedRate = &zeroMillicents
			updatedUnpack = true
			suite.MustSave(&item)
			break
		}
	}
	suite.Equal(true, updatedUnpack)

	// Fetch shipment line items after saves
	fetchedLineItems, err = models.FetchLineItemsByShipmentID(suite.DB(), &shipmentID)
	suite.FatalNoError(err)

	// Verify base shipment line item is  zero
	for _, item := range fetchedLineItems {
		if item.Tariff400ngItem.Code == "4A" {
			suite.Equal(zeroCents, *item.AmountCents)
		}
	}

	shipment2, err := invoice.FetchShipmentForInvoice{DB: suite.DB()}.Call(shipmentID)
	if err != nil {
		suite.logger.Error("Error fetching Shipment for re-pricing line items for shipment", zap.Error(err))
	}
	shipment = &shipment2

	// Verify all base shipment line items are not present
	allPresent = ProcessRecalculateShipment{}.hasAllBaseLineItems(fetchedLineItems)
	suite.Equal(false, allPresent)

	// Re-calculate the Shipment!
	planner := route.NewTestingPlanner(1100)
	engine := rateengine.NewRateEngine(suite.DB(), suite.logger)
	verrs, err := RecalculateShipment{
		DB:      suite.DB(),
		Logger:  suite.logger,
		Engine:  engine,
		Planner: planner,
	}.Call(shipment)
	suite.Equal(false, verrs.HasAny())
	suite.Nil(err, "Failed to recalculate shipment")

	// Fetch shipment line items after recalculation
	fetchedLineItems, err = models.FetchLineItemsByShipmentID(suite.DB(), &shipmentID)
	suite.FatalNoError(err)

	// Verify all base shipment line items are present
	allPresent = ProcessRecalculateShipment{}.hasAllBaseLineItems(fetchedLineItems)
	suite.Equal(true, allPresent)

	// All items should be priced
	// Verify base shipment line item is not zero
	// Verify approved accessorial is not zero
	for _, item := range fetchedLineItems {
		suite.NotEqual(zeroCents, *item.AmountCents)
	}
}

type RecalculateShipmentSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *RecalculateShipmentSuite) SetupTest() {
	suite.DB().TruncateAll()
}
func TestRecalculateShipmentSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &RecalculateShipmentSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}
