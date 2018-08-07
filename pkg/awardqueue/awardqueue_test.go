package awardqueue

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *AwardQueueSuite) Test_CheckAllTSPsBlackedOut() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "5")
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	blackoutStartDate := testdatagen.DateInsidePeakRateCycle
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
	testdatagen.MakeBlackoutDate(suite.db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp.ID,
			StartBlackoutDate:               blackoutStartDate,
			EndBlackoutDate:                 blackoutEndDate,
			TrafficDistributionListID:       &tdl.ID,
		},
	})

	pickupDate := blackoutStartDate.Add(time.Hour)
	deliverDate := blackoutStartDate.Add(time.Hour * 24 * 60)
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC

	shipment, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliverDate, tdl, sourceGBLOC, &market, nil)

	// Create a ShipmentWithOffer to feed the award queue
	shipmentWithOffer := models.ShipmentWithOffer{
		ID: shipment.ID,
		TrafficDistributionListID:       &tdl.ID,
		RequestedPickupDate:             &pickupDate,
		PickupDate:                      &pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        &testdatagen.DateInsidePerformancePeriod,
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

func (suite *AwardQueueSuite) Test_CheckShipmentDuringBlackOut() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "5")

	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC

	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	blackoutStartDate := testdatagen.DateInsidePeakRateCycle
	blackoutEndDate := blackoutStartDate.AddDate(0, 1, 0)
	testdatagen.MakeBlackoutDate(suite.db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp.ID,
			StartBlackoutDate:               blackoutStartDate,
			EndBlackoutDate:                 blackoutEndDate,
			TrafficDistributionListID:       &tdl.ID,
			SourceGBLOC:                     &sourceGBLOC,
			Market:                          &market,
		},
	})

	blackoutPickupDate := blackoutStartDate.AddDate(0, 0, 1)
	blackoutDeliverDate := blackoutStartDate.AddDate(0, 0, 5)
	blackoutShipment, _ := testdatagen.MakeShipment(suite.db, blackoutPickupDate, blackoutPickupDate, blackoutDeliverDate, tdl, sourceGBLOC, &market, nil)

	pickupDate := blackoutEndDate.AddDate(0, 0, 1)
	deliverDate := blackoutEndDate.AddDate(0, 0, 2)
	shipment, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliverDate, tdl, sourceGBLOC, &market, nil)

	// Run the Award Queue
	queue.assignShipments()

	shipmentOffer := models.ShipmentOffer{}
	query := suite.db.Where("shipment_id = $1", shipment.ID)
	if err := query.First(&shipmentOffer); err != nil {
		t.Errorf("Couldn't find shipment offer with shipment_ID: %v\n", shipment.ID)
	}

	blackoutShipmentOffer := models.ShipmentOffer{}
	blackoutQuery := suite.db.Where("shipment_id = $1", blackoutShipment.ID)
	if err := blackoutQuery.First(&blackoutShipmentOffer); err != nil {
		t.Errorf("Couldn't find shipment offer: %v", blackoutShipment.ID)
	}

	if shipmentOffer.AdministrativeShipment != false {
		t.Errorf("Shipment Awards erroneously assigned administrative status.")
	}

	if blackoutShipmentOffer.AdministrativeShipment != true {
		t.Errorf("Shipment Awards erroneously not assigned administrative status.")
	}

	suite.verifyOfferCount(tsp, 2)
}

func (suite *AwardQueueSuite) Test_ShipmentWithinBlackoutDates() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)
	// Creates a TSP and TDL with a blackout date connected to both.
	testTSP1, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	testTDL, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "5")

	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	testStartDate := testdatagen.DateInsidePeakRateCycle
	testEndDate := testStartDate.Add(time.Hour * 24 * 2)
	testdatagen.MakeBlackoutDate(suite.db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: testTSP1.ID,
			StartBlackoutDate:               testStartDate,
			EndBlackoutDate:                 testEndDate,
			TrafficDistributionListID:       &testTDL.ID,
		},
	})

	// Two pickup times to check with ShipmentWithinBlackoutDates
	testPickupDateBetween := testStartDate.Add(time.Hour * 24)
	testPickupDateAfter := testEndDate.Add(time.Hour * 24 * 5)

	// Two shipments using these pickup dates to provide to ShipmentWithinBlackoutDates
	testShipmentBetween, _ := testdatagen.MakeShipment(suite.db, testPickupDateBetween, testStartDate, testEndDate, testTDL, sourceGBLOC, &market, nil)
	testShipmentAfter, _ := testdatagen.MakeShipment(suite.db, testPickupDateAfter, testStartDate, testEndDate, testTDL, sourceGBLOC, &market, nil)

	// One TSP with no blackout dates
	testTSP2, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())

	// Two shipments with offers, using the shipments above
	shipmentWithOfferBetween := models.ShipmentWithOffer{
		ID: testShipmentBetween.ID,
		TrafficDistributionListID:       &testTDL.ID,
		PickupDate:                      &testPickupDateBetween,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        &testdatagen.DateInsidePerformancePeriod,
	}

	shipmentWithOfferAfter := models.ShipmentWithOffer{
		ID: testShipmentAfter.ID,
		TrafficDistributionListID:       &testTDL.ID,
		PickupDate:                      &testPickupDateAfter,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        &testdatagen.DateInsidePerformancePeriod,
	}

	// Checks a date that falls within the blackout date range; returns true.
	test1, err := queue.ShipmentWithinBlackoutDates(testTSP1.ID, shipmentWithOfferBetween)

	if err != nil {
		t.Fatal(err)
	} else if !test1 {
		t.Errorf("Expected true, got false instead.")
	}

	// Checks a date that falls after the blackout date range; returns false.
	test2, err := queue.ShipmentWithinBlackoutDates(testTSP1.ID, shipmentWithOfferAfter)

	if err != nil {
		t.Fatal(err)
	} else if test2 {
		t.Errorf("Expected false, got true instead.")
	}

	// Checks a TSP with no blackout dates and returns false.
	test3, err := queue.ShipmentWithinBlackoutDates(testTSP2.ID, shipmentWithOfferAfter)

	if err != nil {
		t.Fatal(err)
	} else if test3 {
		t.Errorf("Expected false, got true instead.")
	}
}

func (suite *AwardQueueSuite) Test_FindAllUnassignedShipments() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)
	_, err := queue.findAllUnassignedShipments()

	if err != nil {
		t.Error("Unable to find shipments: ", err)
	}
}

// Test that we can create a shipment that should be offered, and that
// it actually gets offered.
func (suite *AwardQueueSuite) Test_OfferSingleShipment() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	// Make a shipment
	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "2")
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliverDate := testdatagen.DateInsidePeakRateCycle

	shipment, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliverDate, tdl, sourceGBLOC, &market, nil)

	// Make a TSP to handle it
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	// Create a ShipmentWithOffer to feed the award queue
	shipmentWithOffer := models.ShipmentWithOffer{
		ID: shipment.ID,
		TrafficDistributionListID:       &tdl.ID,
		PickupDate:                      &pickupDate,
		RequestedPickupDate:             &pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        shipment.BookDate,
	}

	// Run the Award Queue
	offer, err := queue.attemptShipmentOffer(shipmentWithOffer)

	// See if shipment was offered
	if err != nil {
		t.Errorf("Shipment offer expected no errors, received: %v", err)
	} else if offer == nil {
		t.Error("ShipmentOffer was not found.")
	}
}

// Test that we can create a shipment that should NOT be offered because it is not in a TDL
// with any TSPs, and that it doens't get offered.
func (suite *AwardQueueSuite) Test_FailOfferingSingleShipment() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	// Make a shipment in a new TDL, which inherently has no TSPs
	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "2")
	market := "dHHG"
	sourceGBLOC := "OHAI"
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliverDate := testdatagen.DateInsidePeakRateCycle

	shipment, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliverDate, tdl, sourceGBLOC, &market, nil)

	// Create a ShipmentWithOffer to feed the award queue
	shipmentWithOffer := models.ShipmentWithOffer{
		ID: shipment.ID,
		TrafficDistributionListID:       &tdl.ID,
		PickupDate:                      &pickupDate,
		RequestedPickupDate:             &pickupDate,
		BookDate:                        &pickupDate,
		TransportationServiceProviderID: nil,
		AdministrativeShipment:          swag.Bool(false),
	}

	// Run the Award Queue
	offer, err := queue.attemptShipmentOffer(shipmentWithOffer)

	// See if shipment was offered
	if err == nil {
		t.Errorf("Shipment offer expected an error, received none.")
	}
	if offer != nil {
		t.Errorf("Wrong return value, expected nil, got %v", offer)
	}
}

func (suite *AwardQueueSuite) TestAssignShipmentsSingleTSP() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	shipmentsToMake := 10

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "2")

	// Shipment details
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle.Add(time.Hour)

	// Make a few shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliveryDate, tdl, sourceGBLOC, &market, nil)
	}

	// Make a TSP in the same TDL to handle these shipments
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())

	// ... and give this TSP a performance record
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	// Run the Award Queue
	queue.assignShipments()

	// Count the number of shipments offered to our TSP
	query := suite.db.Where("transportation_service_provider_id = $1", tsp.ID)
	offers := []models.ShipmentOffer{}
	count, err := query.Count(&offers)

	if err != nil {
		t.Errorf("Error counting shipment offers: %v", err)
	}
	if count != shipmentsToMake {
		t.Errorf("Not all ShipmentOffers found. Expected %d found %d", shipmentsToMake, count)
	}
}

func (suite *AwardQueueSuite) TestAssignShipmentsToMultipleTSPs() {
	suite.db.TruncateAll()

	queue := NewAwardQueue(suite.db, suite.logger)

	shipmentsToMake := 17

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "2")

	// Shipment details
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle.Add(time.Hour)

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliveryDate, tdl, sourceGBLOC, &market, nil)
	}

	// Make TSPs in the same TDL to handle these shipments
	tsp1, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	tsp2, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	tsp3, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	tsp4, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	tsp5, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())

	// TSPs should be orderd by offer_count first, then BVS.
	testdatagen.MakeTSPPerformance(suite.db, tsp1, tdl, swag.Int(1), mps+5, 0, .4, .4)
	testdatagen.MakeTSPPerformance(suite.db, tsp2, tdl, swag.Int(1), mps+4, 0, .3, .3)
	testdatagen.MakeTSPPerformance(suite.db, tsp3, tdl, swag.Int(2), mps+2, 0, .2, .2)
	testdatagen.MakeTSPPerformance(suite.db, tsp4, tdl, swag.Int(3), mps+3, 0, .1, .1)
	testdatagen.MakeTSPPerformance(suite.db, tsp5, tdl, swag.Int(4), mps+1, 0, .6, .6)

	// Run the Award Queue
	queue.assignShipments()

	suite.verifyOfferCount(tsp1, 6)
	suite.verifyOfferCount(tsp2, 5)
	suite.verifyOfferCount(tsp3, 3)
	suite.verifyOfferCount(tsp4, 2)
	suite.verifyOfferCount(tsp5, 1)
}

func (suite *AwardQueueSuite) Test_GetTSPsPerBandWithRemainder() {
	t := suite.T()
	// Check bands should expect differing num of TSPs when not divisible by 4
	// Remaining TSPs should be divided among bands in descending order
	tspPerBandList := getTSPsPerBand(10)
	expectedBandList := []int{3, 3, 2, 2}
	if !equalSlice(tspPerBandList, expectedBandList) {
		t.Errorf("Failed to correctly divide TSP counts. Expected to find %d, found %d", expectedBandList, tspPerBandList)
	}
}

func (suite *AwardQueueSuite) Test_GetTSPsPerBandNoRemainder() {
	t := suite.T()
	// Check bands should expect correct num of TSPs when num of TSPs is divisible by 4
	tspPerBandList := getTSPsPerBand(8)
	expectedBandList := []int{2, 2, 2, 2}
	if !equalSlice(tspPerBandList, expectedBandList) {
		t.Errorf("Failed to correctly divide TSP counts. Expected to find %d, found %d", expectedBandList, tspPerBandList)
	}
}

func (suite *AwardQueueSuite) Test_AssignTSPsToBands() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)
	tspsToMake := 5

	tdl, err := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "2")
	if err != nil {
		t.Errorf("Failed to create TDL: %v", err)
	}

	for i := 0; i < tspsToMake; i++ {
		tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
		score := float64(mps + i + 1)

		rate := unit.NewDiscountRateFromPercent(45.3)
		testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, nil, score, 0, rate, rate)
	}

	err = queue.assignPerformanceBands()

	if err != nil {
		t.Errorf("Failed to assign to performance bands: %v", err)
	}

	perfs, err := models.FetchTSPPerformanceForQualityBandAssignment(suite.db, tdl.ID, mps)
	if err != nil {
		t.Errorf("Failed to fetch TSPPerformances: %v", err)
	}

	expectedBands := []int{1, 1, 2, 3, 4}

	for i, perf := range perfs {
		band := expectedBands[i]
		if perf.QualityBand == nil {
			t.Errorf("No quality band assigned for Peformance #%v, got nil", perf.ID)
		} else if (*perf.QualityBand) != band {
			t.Errorf("Wrong quality band: expected %v, got %v", band, *perf.QualityBand)
		}
	}
}

// Test_AwardTSPsInDifferentRateCycles ensures that TSPs that service different
// rate cycles get awarded shipments appropriately
func (suite *AwardQueueSuite) Test_AwardTSPsInDifferentRateCycles() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	twoMonths, _ := time.ParseDuration("2 months")
	twoMonthsLater := testdatagen.PerformancePeriodStart.Add(twoMonths)

	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, "2")

	// Make Peak TSP and Shipment
	tspPeak, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())

	tspPerfPeak := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TransportationServiceProviderID: tspPeak.ID,
		TrafficDistributionListID:       tdl.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  100,
		OfferCount:                      0,
	}
	_, err := suite.db.ValidateAndSave(&tspPerfPeak)
	if err != nil {
		t.Error(err)
	}

	shipmentPeak := models.Shipment{
		TrafficDistributionListID: &tdl.ID,
		PickupDate:                &testdatagen.DateInsidePeakRateCycle,
		RequestedPickupDate:       &testdatagen.DateInsidePeakRateCycle,
		DeliveryDate:              &twoMonthsLater,
		BookDate:                  &testdatagen.PerformancePeriodStart,
		SourceGBLOC:               &testdatagen.DefaultSrcGBLOC,
		Market:                    &testdatagen.DefaultMarket,
		MoveID:                    testdatagen.MakeDefaultMove(suite.db).ID,
		Status:                    "DEFAULT",
	}
	_, err = suite.db.ValidateAndSave(&shipmentPeak)
	if err != nil {
		t.Error(err)
	}

	// Make Non-Peak TSP and Shipment
	tspNonPeak, _ := testdatagen.MakeTSP(suite.db, "NPEK")
	tspPerfNonPeak := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.NonPeakRateCycleStart,
		RateCycleEnd:                    testdatagen.NonPeakRateCycleEnd,
		TransportationServiceProviderID: tspNonPeak.ID,
		TrafficDistributionListID:       tdl.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  100,
		OfferCount:                      0,
	}
	_, err = suite.db.ValidateAndSave(&tspPerfNonPeak)
	if err != nil {
		t.Error(err)
	}

	shipmentNonPeak := models.Shipment{
		TrafficDistributionListID: &tdl.ID,
		PickupDate:                &testdatagen.DateInsideNonPeakRateCycle,
		RequestedPickupDate:       &testdatagen.DateInsideNonPeakRateCycle,
		DeliveryDate:              &twoMonthsLater,
		BookDate:                  &testdatagen.PerformancePeriodStart,
		SourceGBLOC:               &testdatagen.DefaultSrcGBLOC,
		Market:                    &testdatagen.DefaultMarket,
		MoveID:                    testdatagen.MakeDefaultMove(suite.db).ID,
		Status:                    "DEFAULT",
	}
	_, err = suite.db.ValidateAndSave(&shipmentNonPeak)
	if err != nil {
		t.Error(err)
	}

	queue.assignShipments()

	suite.verifyOfferCount(tspPeak, 1)
	suite.verifyOfferCount(tspNonPeak, 1)
}

func (suite *AwardQueueSuite) verifyOfferCount(tsp models.TransportationServiceProvider, expectedCount int) {
	t := suite.T()
	t.Helper()

	query := suite.db.Where("transportation_service_provider_id = $1", tsp.ID)
	offers := []models.ShipmentOffer{}
	count, err := query.Count(&offers)

	if err != nil {
		t.Fatalf("Error counting shipment offers: %v", err)
	}
	if count != expectedCount {
		t.Errorf("Wrong number of ShipmentOffers found: expected %d, got %d", expectedCount, count)
	}

	var tspPerformance models.TransportationServiceProviderPerformance
	if err := query.First(&tspPerformance); err != nil {
		t.Errorf("No TSP Performance record found with id %s", tsp.ID)
	}
	if expectedCount != tspPerformance.OfferCount {
		t.Errorf("Wrong OfferCount for TSP: expected %d, got %d", expectedCount, tspPerformance.OfferCount)
	}
}

func equalSlice(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type AwardQueueSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *AwardQueueSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestAwardQueueSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &AwardQueueSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
