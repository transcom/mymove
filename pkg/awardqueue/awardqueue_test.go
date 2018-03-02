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
	testdatagen.MakeTSPPerformance(testDB, tsp, tdl, nil, mps+1, 0)

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

func TestAwardQueueEndToEnd(t *testing.T) {
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
	testdatagen.MakeTSPPerformance(testDB, tsp, tdl, nil, mps+1, 0)

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
