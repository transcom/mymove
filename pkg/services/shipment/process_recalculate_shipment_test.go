package shipment

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

// Deliver marks the Shipment request as Delivered. Must be IN TRANSIT state.
func helperForceDeliver(s *models.Shipment, actualDeliveryDate time.Time) error {
	s.Status = models.ShipmentStatusDELIVERED
	s.ActualDeliveryDate = &actualDeliveryDate
	pickup := actualDeliveryDate.AddDate(0, 0, -2)
	s.ActualPickupDate = &pickup
	pack := actualDeliveryDate.AddDate(0, 0, -3)
	s.ActualPackDate = &pack

	return nil
}

func (suite *ProcessRecalculateShipmentSuite) helperCreateShipmentAndPlanner() (*models.Shipment, *route.Planner) {
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

	// Make sure there's a FuelEIADieselPrice
	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	planner := route.NewTestingPlanner(1100)
	//engine := rateengine.NewRateEngine(suite.DB(), suite.logger, planner)

	return &shipment, &planner
}

func (suite *ProcessRecalculateShipmentSuite) TestProcessRecalculateShipmentCall() {

	shipment, planner := suite.helperCreateShipmentAndPlanner()
	shipmentID := shipment.ID

	shipmentLineItems, err := models.FetchLineItemsByShipmentID(suite.DB(), &shipmentID)
	suite.Nil(err)

	shipment2, err := invoice.FetchShipmentForInvoice{DB: suite.DB()}.Call(shipmentID)
	if err != nil {
		suite.logger.Error("Error fetching Shipment for re-pricing line items for shipment", zap.Error(err))
	}
	shipment = &shipment2

	//
	// TEST: No date range records (return false)
	//

	dates, err := models.FetchShipmentRecalculateDates(suite.DB())
	suite.Empty(dates)

	update, err := ProcessRecalculateShipment{
		DB:     suite.DB(),
		Logger: suite.logger,
	}.Call(shipment, shipmentLineItems, *planner)
	// TEST Validation: No date range records (return false)
	suite.Nil(err)
	suite.Equal(false, update)

	id, err := uuid.NewV4()
	// Create date range
	recalculateRange := models.ShipmentRecalculate{
		ID:                    id,
		ShipmentUpdatedAfter:  time.Date(1970, time.January, 01, 0, 0, 0, 0, time.UTC),
		ShipmentUpdatedBefore: time.Now(),
		Active:                true,
	}
	suite.MustCreate(suite.DB(), &recalculateRange)
	fetchRecalculateRange, err := models.FetchShipmentRecalculateDates(suite.DB())
	recalculateRange = *fetchRecalculateRange
	suite.Nil(err)

	//
	// TEST: shipment is not in DELIVERED or COMPLETED state (return false)
	//

	suite.NotEqual(models.ShipmentStatusDELIVERED, shipment.Status)
	suite.NotEqual(models.ShipmentStatusCOMPLETED, shipment.Status)

	update, err = ProcessRecalculateShipment{
		DB:     suite.DB(),
		Logger: suite.logger,
	}.Call(shipment, shipmentLineItems, *planner)
	// TEST Validation: shipment is not in DELIVERED or COMPLETED state (return false)
	suite.Nil(err)
	suite.Equal(false, update)

	// Shipment is delivered
	// Bypassing the shipment.Deliver(testdatagen.DateInsidePeakRateCycle)
	helperForceDeliver(shipment, testdatagen.DateInsidePeakRateCycle.AddDate(0, 0, 2))

	suite.MustSave(shipment)
	suite.Equal(models.ShipmentStatusDELIVERED, shipment.Status, "expected Delivered")

	//
	// TEST: shipment after date range (return false)
	//

	shipment.CreatedAt = recalculateRange.ShipmentUpdatedBefore.AddDate(0, 0, 2)
	suite.MustSave(shipment)

	update, err = ProcessRecalculateShipment{
		DB:     suite.DB(),
		Logger: suite.logger,
	}.Call(shipment, shipmentLineItems, *planner)
	// TEST Validation: shipment after date range (return false)
	suite.Nil(err)
	suite.Equal(false, update)

	//
	// TEST: Shipment missing base line item or line item was updated in date range (return true)
	//

	// Move recalculate date range to a date in the past
	recalculateRange.ShipmentUpdatedBefore = recalculateRange.ShipmentUpdatedBefore.AddDate(-5, 0, 0)
	suite.MustSave(&recalculateRange)

	// Move Shipment.CreatedAt to date within range
	shipment.CreatedAt = recalculateRange.ShipmentUpdatedBefore.AddDate(0, 0, -1)
	suite.MustSave(shipment)

	update, err = ProcessRecalculateShipment{
		DB:     suite.DB(),
		Logger: suite.logger,
	}.Call(shipment, shipmentLineItems, *planner)
	// TEST Validation: Shipment missing base line item or line item was updated in date range (return true)
	suite.Nil(err)
	suite.Equal(true, update)

	shipmentLineItems, err = models.FetchLineItemsByShipmentID(suite.DB(), &shipmentID)
	suite.Nil(err)

	shipment2, err = invoice.FetchShipmentForInvoice{DB: suite.DB()}.Call(shipmentID)
	if err != nil {
		suite.logger.Error("Error fetching Shipment for re-pricing line items for shipment", zap.Error(err))
	}
	shipment = &shipment2

	// Move recalculate date range back to a date earlier than the newly generated line items
	// Shipment line items are all present and all updated
	recalculateRange.ShipmentUpdatedBefore = time.Now().AddDate(0, -1, 0)
	suite.MustSave(&recalculateRange)

	fetchRecalculateRange, err = models.FetchShipmentRecalculateDates(suite.DB())
	suite.Nil(err)
	suite.NotNil(fetchRecalculateRange)

	// Move Shipment.CreatedAt to date within range
	shipment.CreatedAt = testdatagen.DateInsidePeakRateCycle.AddDate(0, 0, 2)
	suite.MustSave(shipment)

	// TEST: Do not recalculate shipment all line items are up to date (return false)
	update, err = ProcessRecalculateShipment{
		DB:     suite.DB(),
		Logger: suite.logger,
	}.Call(shipment, shipmentLineItems, *planner)
	// TEST Validation:  Do not recalculate shipment all line items are up to date (return false)
	suite.Nil(err)
	suite.Equal(false, update)
}

type ProcessRecalculateShipmentSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *ProcessRecalculateShipmentSuite) SetupTest() {
	suite.DB().TruncateAll()
}
func TestProcessRecalculateShipmentSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &ProcessRecalculateShipmentSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}
