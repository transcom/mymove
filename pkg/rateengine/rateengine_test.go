package awardqueue

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

func (suite *RateEngineSuite) Test_CheckAllTSPsBlackedOut() {
	t := suite.T()
	queue := NewRateEngine(suite.db, suite.logger)

	tsp, err := testdatagen.MakeTSP(suite.db, "A Very Excellent TSP", "XYZA")
	tdl, err := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "5")
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0)

	blackoutStartDate := testdatagen.DateInsidePeakRateCycle
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
	testdatagen.MakeBlackoutDate(suite.db, tsp, blackoutStartDate, blackoutEndDate, &tdl, nil, nil)

	pickupDate := blackoutStartDate.Add(time.Hour)
	deliverDate := blackoutStartDate.Add(time.Hour * 24 * 60)
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC

	shipment, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliverDate, tdl, sourceGBLOC, &market)

	// Create a ShipmentWithOffer to feed the award queue
	shipmentWithOffer := models.ShipmentWithOffer{
		ID: shipment.ID,
		TrafficDistributionListID:       tdl.ID,
		RequestedPickupDate:             pickupDate,
		PickupDate:                      pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        testdatagen.DateInsidePerformancePeriod,
	}

	// Run the Award Queue
	offer, err := queue.attemptShipmentOffer(shipmentWithOffer)

	expectedError := "could not find a TSP without blackout dates"
	// See if shipment was offered
	if err == nil || offer != nil {
		t.Errorf("Shipment was offered to a blacked out TSP!")
	} else if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Did not receive proper error message. Expected '%s', got '%s' instead.", expectedError, err)
	}
}

type RateEngineSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *RateEngineSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestRateEngineSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &RateEngineSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
