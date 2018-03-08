package awardqueue

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

var testDB *pop.Connection

func TestFindAllUnawardedShipments(t *testing.T) {
	queue := NewAwardQueue(testDB)
	_, err := queue.findAllUnawardedShipments()

	if err != nil {
		t.Fatal("Unable to find shipments: ", err)
	}
}

// Test that we can create a shipment that should be awarded, and that
// it actually gets awarded.
func TestAwardSingleShipment(t *testing.T) {
	queue := NewAwardQueue(testDB)

	// Make a shipment
	tdl, _ := testdatagen.MakeTDL(testDB, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(testDB, time.Now(), time.Now(), tdl)

	// Make a TSP to handle it
	tsp, _ := testdatagen.MakeTSP(testDB, "Test Shipper", "TEST")
	testdatagen.MakeTSPPerformance(testDB, tsp, tdl, swag.Int(1), mps+1, 0)

	// Create a PossiblyAwardedShipment to feed the award queue
	pas := models.PossiblyAwardedShipment{
		ID: shipment.ID,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
	}

	// Run the Award Queue
	award, err := queue.attemptShipmentAward(pas)

	// See if shipment was awarded
	if err != nil {
		t.Errorf("Shipment award expected no errors, received: %v", err)
	}
	if award == nil {
		t.Error("ShipmentAward was not found.")
	}
}

// Test that we can create a shipment that should NOT be awarded because it is not in a TDL
// with any TSPs, and that it doens't get awarded.
func TestFailAwardingSingleShipment(t *testing.T) {
	queue := NewAwardQueue(testDB)

	// Make a shipment in a new TDL, which inherently has no TSPs
	tdl, _ := testdatagen.MakeTDL(testDB, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(testDB, time.Now(), time.Now(), tdl)

	// Create a PossiblyAwardedShipment to feed the award queue
	pas := models.PossiblyAwardedShipment{
		ID: shipment.ID,
		TrafficDistributionListID:       tdl.ID,
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
		t.Error("ShipmentAward was created, expected 'nil'.")
	}
}

func TestAwardAssignUnawardedShipmentsSingleTSP(t *testing.T) {
	queue := NewAwardQueue(testDB)

	shipmentsToMake := 10

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(testDB, "california", "90210", "2")

	// Make a few shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(testDB, time.Now(), time.Now(), tdl)
	}

	// Make a TSP in the same TDL to handle these shipments
	tsp, _ := testdatagen.MakeTSP(testDB, "Test Shipper", "TEST")

	// ... and give this TSP a performance record
	testdatagen.MakeTSPPerformance(testDB, tsp, tdl, swag.Int(1), mps+1, 0)

	// Run the Award Queue
	queue.assignUnawardedShipments()

	// Count the number of shipments awarded to our TSP
	query := testDB.Where("transportation_service_provider_id = $1", tsp.ID)
	awards := []models.ShipmentAward{}
	count, err := query.Count(&awards)

	if err != nil {
		t.Errorf("Error counting shipment awards: %v", err)
	}
	if count != shipmentsToMake {
		t.Errorf("Not all ShipmentAwards found. Expected %d found %d", shipmentsToMake, count)
	}
}

func TestAwardAssignUnawardedShipmentsToMultipleTSPs(t *testing.T) {
	testDB.TruncateAll()

	queue := NewAwardQueue(testDB)

	shipmentsToMake := 17

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(testDB, "california", "90210", "2")

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(testDB, time.Now(), time.Now(), tdl)
	}

	// Make TSPs in the same TDL to handle these shipments
	tsp1, _ := testdatagen.MakeTSP(testDB, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(testDB, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(testDB, "Test TSP 3", "TSP3")
	tsp4, _ := testdatagen.MakeTSP(testDB, "Test TSP 4", "TSP4")
	tsp5, _ := testdatagen.MakeTSP(testDB, "Test TSP 5", "TSP5")

	// TSPs should be orderd by award_count first, then BVS.
	testdatagen.MakeTSPPerformance(testDB, tsp1, tdl, swag.Int(1), mps+5, 0)
	testdatagen.MakeTSPPerformance(testDB, tsp2, tdl, swag.Int(1), mps+4, 0)
	testdatagen.MakeTSPPerformance(testDB, tsp3, tdl, swag.Int(2), mps+2, 0)
	testdatagen.MakeTSPPerformance(testDB, tsp4, tdl, swag.Int(3), mps+3, 0)
	testdatagen.MakeTSPPerformance(testDB, tsp5, tdl, swag.Int(4), mps+1, 0)

	// Run the Award Queue
	queue.assignUnawardedShipments()

	verifyAwardCount(t, tsp1, 6)
	verifyAwardCount(t, tsp2, 5)
	verifyAwardCount(t, tsp3, 3)
	verifyAwardCount(t, tsp4, 2)
	verifyAwardCount(t, tsp5, 1)
}

func verifyAwardCount(t *testing.T, tsp models.TransportationServiceProvider, expectedCount int) {
	t.Helper()

	// TODO is there a more concise way to do this?
	query := testDB.Where("transportation_service_provider_id = $1", tsp.ID)
	awards := []models.ShipmentAward{}
	count, err := query.Count(&awards)

	if err != nil {
		t.Fatalf("Error counting shipment awards: %v", err)
	}
	if count != expectedCount {
		t.Errorf("Wrong number of ShipmentAwards found: expected %d, got %d", expectedCount, count)
	}
}

func Test_getTSPsPerBandWithRemainder(t *testing.T) {
	// Check bands should expect differing num of TSPs when not divisible by 4
	// Remaining TSPs should be divided among bands in descending order
	tspPerBandList := getTSPsPerBand(10)
	expectedBandList := []int{3, 3, 2, 2}
	if !equalSlice(tspPerBandList, []int{3, 3, 2, 2}) {
		t.Errorf("Failed to correctly divide TSP counts. Expected to find %d, found %d", expectedBandList, tspPerBandList)
	}
}

func Test_getTSPsPerBandNoRemainder(t *testing.T) {
	// Check bands should expect correct num of TSPs when num of TSPs is divisible by 4
	tspPerBandList := getTSPsPerBand(8)
	expectedBandList := []int{2, 2, 2, 2}
	if !equalSlice(tspPerBandList, []int{2, 2, 2, 2}) {
		t.Errorf("Failed to correctly divide TSP counts. Expected to find %d, found %d", expectedBandList, tspPerBandList)
	}
}

func Test_assignTSPsToBands(t *testing.T) {
	pop.Debug = true
	queue := NewAwardQueue(testDB)
	tspsToMake := 5

	tdl, err := testdatagen.MakeTDL(testDB, "california", "90210", "2")
	if err != nil {
		t.Errorf("Failed to create TDL: %v", err)
	}

	for i := 0; i < tspsToMake; i++ {
		tsp, _ := testdatagen.MakeTSP(testDB, "Test Shipper", "TEST")
		score := mps + i + 1
		testdatagen.MakeTSPPerformance(testDB, tsp, tdl, nil, score, 0)
	}

	err = queue.assignPerformanceBands()

	if err != nil {
		t.Errorf("Failed to assign to performance bands: %v", err)
	}

	perfs, err := models.FetchTSPPerformanceForQualityBandAssignment(testDB, tdl.ID, mps)
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
	pop.Debug = false
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

func setupDBConnection() {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	conn, err := pop.Connect("test")
	if err != nil {
		fmt.Println(err)
	}

	testDB = conn
}

func TestMain(m *testing.M) {
	setupDBConnection()

	os.Exit(m.Run())
}
