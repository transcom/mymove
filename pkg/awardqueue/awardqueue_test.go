package awardqueue

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/markbates/pop"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AwardQueueSuite) Test_CheckAllTSPsBlackedOut() {
	t := suite.T()
	queue := NewAwardQueue(suite.db)

	tsp, err := testdatagen.MakeTSP(suite.db, "A Very Excellent TSP", "XYZA")
	tdl, err := testdatagen.MakeTDL(suite.db, "Oklahoma", "62240", "5")
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0)
	blackoutStartDate := time.Now()
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
	testdatagen.MakeBlackoutDate(suite.db, tsp, blackoutStartDate, blackoutEndDate, &tdl, nil, nil, nil, nil)

	pickupDate := blackoutStartDate.Add(time.Hour)
	deliverDate := blackoutStartDate.Add(time.Hour * 24 * 60)
	shipment, _ := testdatagen.MakeShipment(suite.db, pickupDate, deliverDate, tdl)

	// Create a PossiblyAwardedShipment to feed the award queue
	pas := models.PossiblyAwardedShipment{
		ID: shipment.ID,
		TrafficDistributionListID:       tdl.ID,
		PickupDate:                      pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		AwardDate:                       testdatagen.DateInsidePerformancePeriod,
	}

	// Run the Award Queue
	award, err := queue.attemptShipmentAward(pas)

	expectedError := "Could not find a TSP without blackout dates"
	// See if shipment was awarded
	if err == nil || award != nil {
		t.Errorf("Shipment was awarded to a blacked out TSP!")
	} else if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Did not receive proper error message. Expected '%s', got '%s' instead.", expectedError, err)
	}
}

func (suite *AwardQueueSuite) Test_CheckShipmentDuringBlackOut() {
	t := suite.T()
	queue := NewAwardQueue(suite.db)

	tsp, _ := testdatagen.MakeTSP(suite.db, "A Very Excellent TSP", "XYZA")
	tdl, _ := testdatagen.MakeTDL(suite.db, "Oklahoma", "62240", "5")
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0)
	blackoutStartDate := time.Now().AddDate(1, 0, 0)
	blackoutEndDate := blackoutStartDate.AddDate(0, 1, 0)
	testdatagen.MakeBlackoutDate(suite.db, tsp, blackoutStartDate, blackoutEndDate, &tdl, nil, nil, nil, nil)
	pickupDate := blackoutStartDate.AddDate(0, 0, 1)
	deliverDate := blackoutStartDate.AddDate(0, 0, 5)

	// Create a shipment within blackout dates and one not within blackout dates
	blackoutShipment, _ := testdatagen.MakeShipment(suite.db, pickupDate, deliverDate, tdl)
	shipment, _ := testdatagen.MakeShipment(suite.db, time.Now(), time.Now().AddDate(0, 0, 1), tdl)

	// Run the Award Queue
	queue.assignUnawardedShipments()

	shipmentAward := models.ShipmentAward{}
	query := suite.db.Where("shipment_id = $1", shipment.ID)
	if err := query.First(&shipmentAward); err != nil {
		t.Errorf("Couldn't find shipment award with shipment_ID: %v\n", shipment.ID)
	}

	blackoutShipmentAward := models.ShipmentAward{}
	blackoutQuery := suite.db.Where("shipment_id = $1", blackoutShipment.ID)
	if err := blackoutQuery.First(&blackoutShipmentAward); err != nil {
		t.Errorf("Couldn't find shipment award: %v", blackoutShipment.ID)
	}

	if shipmentAward.AdministrativeShipment != false || blackoutShipmentAward.AdministrativeShipment != true {
		t.Errorf("Shipment Awards erroneously assigned administrative status.")
	}

	suite.verifyAwardCount(tsp, 2)
}

func (suite *AwardQueueSuite) Test_ShipmentWithinBlackoutDates() {
	t := suite.T()
	queue := NewAwardQueue(suite.db)
	// Creates a TSP and TDL with a blackout date connected to both.
	testTSP1, _ := testdatagen.MakeTSP(suite.db, "A Very Excellent TSP", "XYZA")
	testTDL, _ := testdatagen.MakeTDL(suite.db, "Oklahoma", "62240", "5")
	testStartDate := time.Now()
	testEndDate := testStartDate.Add(time.Hour * 24 * 2)
	testdatagen.MakeBlackoutDate(suite.db, testTSP1, testStartDate, testEndDate, &testTDL, nil, nil, nil, nil)

	// Two pickup times to check with ShipmentWithinBlackoutDates
	testPickupDateBetween := testStartDate.Add(time.Hour * 24)
	testPickupDateAfter := testEndDate.Add(time.Hour * 24 * 5)

	// One TSP with no blackout dates
	testTSP2, _ := testdatagen.MakeTSP(suite.db, "A Spotless TSP", "PORK")

	// Checks a date that falls within the blackout date range; returns true.
	test1, err := queue.ShipmentWithinBlackoutDates(testTSP1.ID, testPickupDateBetween)

	if err != nil {
		t.Fatal(err)
	} else if !test1 {
		t.Errorf("Expected true, got false instead.")
	}

	// Checks a date that falls after the blackout date range; returns false.
	test2, err := queue.ShipmentWithinBlackoutDates(testTSP1.ID, testPickupDateAfter)

	if err != nil {
		t.Fatal(err)
	} else if test2 {
		t.Errorf("Expected false, got true instead.")
	}

	// Checks a TSP with no blackout dates and returns false.
	test3, err := queue.ShipmentWithinBlackoutDates(testTSP2.ID, testPickupDateAfter)

	if err != nil {
		t.Fatal(err)
	} else if test3 {
		t.Errorf("Expected false, got true instead.")
	}
}

func (suite *AwardQueueSuite) Test_FindAllUnawardedShipments() {
	t := suite.T()
	queue := NewAwardQueue(suite.db)
	_, err := queue.findAllUnawardedShipments()

	if err != nil {
		t.Error("Unable to find shipments: ", err)
	}
}

// Test that we can create a shipment that should be awarded, and that
// it actually gets awarded.
func (suite *AwardQueueSuite) Test_AwardSingleShipment() {
	t := suite.T()
	queue := NewAwardQueue(suite.db)

	// Make a shipment
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(suite.db, time.Now(), time.Now(), tdl)

	// Make a TSP to handle it
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test Shipper", "TEST")
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0)

	// Create a PossiblyAwardedShipment to feed the award queue
	pas := models.PossiblyAwardedShipment{
		ID: shipment.ID,
		TrafficDistributionListID:       tdl.ID,
		PickupDate:                      time.Now(),
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		AwardDate:                       shipment.AwardDate,
	}

	// Run the Award Queue
	award, err := queue.attemptShipmentAward(pas)

	// See if shipment was awarded
	if err != nil {
		t.Errorf("Shipment award expected no errors, received: %v", err)
	} else if award == nil {
		t.Error("ShipmentAward was not found.")
	}
}

// Test that we can create a shipment that should NOT be awarded because it is not in a TDL
// with any TSPs, and that it doens't get awarded.
func (suite *AwardQueueSuite) Test_FailAwardingSingleShipment() {
	t := suite.T()
	queue := NewAwardQueue(suite.db)

	// Make a shipment in a new TDL, which inherently has no TSPs
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(suite.db, time.Now(), time.Now(), tdl)

	// Create a PossiblyAwardedShipment to feed the award queue
	pas := models.PossiblyAwardedShipment{
		ID: shipment.ID,
		TrafficDistributionListID:       tdl.ID,
		PickupDate:                      time.Now(),
		TransportationServiceProviderID: nil,
		AdministrativeShipment:          swag.Bool(false),
	}

	// Run the Award Queue
	award, err := queue.attemptShipmentAward(pas)

	// See if shipment was awarded
	if err == nil {
		t.Errorf("Shipment award expected an error, received none.")
	}
	if award != nil {
		t.Errorf("Wrong return value, expected nil, got %v", award)
	}
}

func (suite *AwardQueueSuite) Test_AwardAssignUnawardedShipmentsSingleTSP() {
	t := suite.T()
	queue := NewAwardQueue(suite.db)

	shipmentsToMake := 10

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")

	// Make a few shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(suite.db, time.Now(), time.Now(), tdl)
	}

	// Make a TSP in the same TDL to handle these shipments
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test Shipper", "TEST")

	// ... and give this TSP a performance record
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0)

	// Run the Award Queue
	queue.assignUnawardedShipments()

	// Count the number of shipments awarded to our TSP
	query := suite.db.Where("transportation_service_provider_id = $1", tsp.ID)
	awards := []models.ShipmentAward{}
	count, err := query.Count(&awards)

	if err != nil {
		t.Errorf("Error counting shipment awards: %v", err)
	}
	if count != shipmentsToMake {
		t.Errorf("Not all ShipmentAwards found. Expected %d found %d", shipmentsToMake, count)
	}
}

func (suite *AwardQueueSuite) Test_AwardAssignUnawardedShipmentsToMultipleTSPs() {
	suite.db.TruncateAll()

	queue := NewAwardQueue(suite.db)

	shipmentsToMake := 17

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(suite.db, time.Now(), time.Now(), tdl)
	}

	// Make TSPs in the same TDL to handle these shipments
	tsp1, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(suite.db, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(suite.db, "Test TSP 3", "TSP3")
	tsp4, _ := testdatagen.MakeTSP(suite.db, "Test TSP 4", "TSP4")
	tsp5, _ := testdatagen.MakeTSP(suite.db, "Test TSP 5", "TSP5")

	// TSPs should be orderd by award_count first, then BVS.
	testdatagen.MakeTSPPerformance(suite.db, tsp1, tdl, swag.Int(1), mps+5, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp2, tdl, swag.Int(1), mps+4, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp3, tdl, swag.Int(2), mps+2, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp4, tdl, swag.Int(3), mps+3, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp5, tdl, swag.Int(4), mps+1, 0)

	// Run the Award Queue
	queue.assignUnawardedShipments()

	suite.verifyAwardCount(tsp1, 6)
	suite.verifyAwardCount(tsp2, 5)
	suite.verifyAwardCount(tsp3, 3)
	suite.verifyAwardCount(tsp4, 2)
	suite.verifyAwardCount(tsp5, 1)
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
	queue := NewAwardQueue(suite.db)
	tspsToMake := 5

	tdl, err := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	if err != nil {
		t.Errorf("Failed to create TDL: %v", err)
	}

	for i := 0; i < tspsToMake; i++ {
		tsp, _ := testdatagen.MakeTSP(suite.db, "Test Shipper", "TEST")
		score := mps + i + 1
		testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, nil, score, 0)
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

func (suite *AwardQueueSuite) verifyAwardCount(tsp models.TransportationServiceProvider, expectedCount int) {
	t := suite.T()
	t.Helper()

	// TODO is there a more concise way to do this?
	query := suite.db.Where("transportation_service_provider_id = $1", tsp.ID)
	awards := []models.ShipmentAward{}
	count, err := query.Count(&awards)

	if err != nil {
		t.Fatalf("Error counting shipment awards: %v", err)
	}
	if count != expectedCount {
		t.Errorf("Wrong number of ShipmentAwards found: expected %d, got %d", expectedCount, count)
	}

	var tspPerformance models.TransportationServiceProviderPerformance
	if err := query.First(&tspPerformance); err != nil {
		t.Errorf("No TSP Performance record found with id %s", tsp.ID)
	}
	if expectedCount != tspPerformance.AwardCount {
		t.Errorf("Wrong AwardCount for TSP: expected %d, got %d", expectedCount, tspPerformance.AwardCount)
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
	db *pop.Connection
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

	hs := &AwardQueueSuite{db: db}
	suite.Run(t, hs)
}
