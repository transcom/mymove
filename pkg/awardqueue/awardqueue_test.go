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

func TestFindAllUnawardedShipments(t *testing.T) {
	_, err := findAllUnawardedShipments()

	if err != nil {
		t.Fatal("Unable to find shipments: ", err)
	}
}

// Test that we can create a shipment that should be awarded, and that
// it actually gets awarded.
func TestAwardSingleShipment(t *testing.T) {
	// Make a shipment
	tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(db, time.Now(), time.Now(), tdl)

	// Make a TSP to handle it
	tsp, _ := testdatagen.MakeTSP(db, "Test Shipper", "TEST")
	testdatagen.MakeTSPPerformance(db, tsp, tdl, nil, mps+1, 0)

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
	award, err := AttemptShipmentAward(pas)

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
	// Make a shipment in a new TDL, which inherently has no TSPs
	tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(db, time.Now(), time.Now(), tdl)

	// Create a PossiblyAwardedShipment to feed the award queue
	pas := models.PossiblyAwardedShipment{
		ID: shipment.ID,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: nil,
		AdministrativeShipment:          swag.Bool(false),
	}

	// Run the Award Queue
	award, err := AttemptShipmentAward(pas)

	// See if shipment was awarded
	if err == nil {
		t.Errorf("Shipment award expected an error, received none.")
	}
	if award != nil {
		t.Error("ShipmentAward was created, expected 'nil'.")
	}
}

func TestAwardQueueEndToEnd(t *testing.T) {
	shipmentsToMake := 10

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")

	// Make a few shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, time.Now(), time.Now(), tdl)
	}

	// Make a TSP in the same TDL to handle these shipments
	tsp, _ := testdatagen.MakeTSP(db, "Test Shipper", "TEST")

	// ... and give this TSP a performance record
	testdatagen.MakeTSPPerformance(db, tsp, tdl, nil, mps+1, 0)

	// Run the Award Queue
	Run(db)

	// Count the number of shipments awarded to our TSP
	query := db.Where("transportation_service_provider_id = $1", tsp.ID)
	awards := []models.ShipmentAward{}
	count, err := query.Count(&awards)

	if err != nil {
		t.Errorf("Error counting shipment awards: %v", err)
	}
	if count != shipmentsToMake {
		t.Errorf("Not all ShipmentAwards found. Expected %d found %d", shipmentsToMake, count)
	}
}

// Test_FetchTSPPerformanceForAwardQueue ensures that TSPs are returned in the expected
// order for the Award Queue operation.
func Test_FetchTSPPerformanceForAwardQueue(t *testing.T) {
	tdl, _ := testdatagen.MakeTDL(db, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(db, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(db, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(db, "Test TSP 3", "TSP2")
	// TSPs should be orderd by award_count first, then BVS.
	testdatagen.MakeTSPPerformance(db, tsp1, tdl, nil, mps+1, 0)
	testdatagen.MakeTSPPerformance(db, tsp2, tdl, nil, mps+3, 1)
	testdatagen.MakeTSPPerformance(db, tsp3, tdl, nil, mps+2, 1)

	tsps, err := models.FetchTSPPerformanceForAwardQueue(db, tdl.ID, mps)

	if err != nil {
		t.Errorf("Failed to find TSP: %v", err)
	} else if len(tsps) != 3 {
		t.Errorf("Failed to find TSPs. Expected to find 3, found %d", len(tsps))
	} else if tsps[0].TransportationServiceProviderID != tsp1.ID &&
		tsps[1].TransportationServiceProviderID != tsp2.ID &&
		tsps[2].TransportationServiceProviderID != tsp3.ID {

		t.Errorf("TSPs returned out of expected order.\n"+
			"\tExpected: [%s, %s, %s]\nFound:    [%s, %s, %s]",
			tsp1.ID, tsp2.ID, tsp3.ID,
			tsps[0].TransportationServiceProviderID,
			tsps[1].TransportationServiceProviderID,
			tsps[2].TransportationServiceProviderID)
	}
}

// Test_MinimumPerformanceScore ensures that TSPs whose BVS is below the MPS
// do not enter the Award Queue process.
func Test_MinimumPerformanceScore(t *testing.T) {
	tdl, _ := testdatagen.MakeTDL(db, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(db, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(db, "Test TSP 2", "TSP2")
	// Make 2 TSPs, one with a BVS above the MPS and one below the MPS.
	testdatagen.MakeTSPPerformance(db, tsp1, tdl, nil, mps+1, 0)
	testdatagen.MakeTSPPerformance(db, tsp2, tdl, nil, mps-1, 1)

	tsps, err := models.FetchTSPPerformanceForQualityBandAssignment(db, tdl.ID, mps)

	if err != nil {
		t.Errorf("Failed to find TSP: %v", err)
	} else if len(tsps) != 1 {
		t.Errorf("Failed to find TSPs. Expected to find 1, found %d", len(tsps))
	} else if tsps[0].TransportationServiceProviderID != tsp1.ID {
		t.Errorf("Incorrect TSP returned. Expected %s, received %s.",
			tsp1.ID,
			tsps[0].TransportationServiceProviderID)
	}
}

// Test_FetchTSPPerformanceForQualityBandAssignment ensures that TSPs are returned in the expected
// order for the division into quality bands.
func Test_FetchTSPPerformanceForQualityBandAssignment(t *testing.T) {
	tdl, _ := testdatagen.MakeTDL(db, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(db, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(db, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(db, "Test TSP 3", "TSP2")
	// What matter is the BVS score order; award_count has no influence.
	testdatagen.MakeTSPPerformance(db, tsp1, tdl, nil, 90, 0)
	testdatagen.MakeTSPPerformance(db, tsp2, tdl, nil, 50, 1)
	testdatagen.MakeTSPPerformance(db, tsp3, tdl, nil, 15, 1)

	tsps, err := models.FetchTSPPerformanceForQualityBandAssignment(db, tdl.ID, mps)

	if err != nil {
		t.Errorf("Failed to find TSP: %v", err)
	} else if len(tsps) != 3 {
		t.Errorf("Failed to find TSPs. Expected to find 3, found %d", len(tsps))
	} else if tsps[0].TransportationServiceProviderID != tsp1.ID &&
		tsps[1].TransportationServiceProviderID != tsp2.ID &&
		tsps[2].TransportationServiceProviderID != tsp3.ID {

		t.Errorf("TSPs returned out of expected order.\n"+
			"\tExpected: [%s, %s, %s]\nFound:    [%s, %s, %s]",
			tsp1.ID, tsp2.ID, tsp3.ID,
			tsps[0].TransportationServiceProviderID,
			tsps[1].TransportationServiceProviderID,
			tsps[2].TransportationServiceProviderID)
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
	tspsToMake := 5

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")

	// Make 5 (not divisible by 4) TSPs in this TDL with BVSs
	for i := 0; i < tspsToMake; i++ {
		tsp, _ := testdatagen.MakeTSP(db, "Test Shipper", "TEST")
		testdatagen.MakeTSPPerformance(db, tsp, tdl, nil, mps+1, 0)
	}
	// Fetch TSPs in TDL
	tspPerfs, err := models.FetchTSPPerformanceForQualityBandAssignment(db, tdl.ID, mps)
	qbs := assignTSPsToBands(tspPerfs)
	if err != nil {
		t.Errorf("Failed to find TSPs: %v", err)
	}
	if len(qbs[0]) != 2 || len(qbs[1]) != 1 {
		t.Errorf("Failed to correctly add TSPs to quality bands.")
	}
}

func Test_BVSWithLowMPS(t *testing.T) {
	tspsToMake := 5

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")

	// Make 5 (not divisible by 4) TSPs in this TDL with BVSs above MPS threshold
	for i := 0; i < tspsToMake; i++ {
		tsp, _ := testdatagen.MakeTSP(db, "Test Shipper", "TEST")
		testdatagen.MakeTSPPerformance(db, tsp, tdl, nil, 15, 0)
	}
	// Make 1 TSP in this TDL with BVS below the MPS threshold
	mpsTSP, _ := testdatagen.MakeTSP(db, "Low BVS Test Shipper", "TEST")
	testdatagen.MakeTSPPerformance(db, mpsTSP, tdl, nil, mps-1, 0)

	// Fetch TSPs in TDL
	tspsbb, err := models.FetchTSPPerformanceForQualityBandAssignment(db, tdl.ID, mps)

	// Then: Expect to find TSPs in TDL
	if err != nil {
		t.Errorf("Failed to find TSPs: %v", err)
	}
	// Then: Expect TSP with low BVS won't be in sorted TSP slice
	for _, tsp := range tspsbb {
		if tsp.ID == mpsTSP.ID {
			t.Errorf("TSP: %v with a BVS below MPS incorrectly included.", mpsTSP.ID)
		}
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

func setupDBConnection() {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	conn, err := pop.Connect("test")
	if err != nil {
		fmt.Println(err)
	}

	db = conn
}

func TestMain(m *testing.M) {
	setupDBConnection()

	os.Exit(m.Run())
}
