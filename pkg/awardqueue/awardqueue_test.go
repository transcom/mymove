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
	testdatagen.MakeBestValueScore(db, tsp, tdl, 10)

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

	// ... and give this TSP a BVS
	testdatagen.MakeBestValueScore(db, tsp, tdl, 10)

	// Run the Award Queue
	Run(db)

	// Count the number of shipments awarded to our TSP
	query := db.Where(fmt.Sprintf("transportation_service_provider_id = '%s'", tsp.ID))
	awards := []models.ShipmentAward{}
	count, err := query.Count(&awards)

	if err != nil {
		t.Errorf("Error counting shipment awards: %v", err)
	}
	if count != shipmentsToMake {
		t.Errorf("Not all ShipmentAwards found. Expected %d found %d", shipmentsToMake, count)
	}
}

// Ensure that if we create a TSP in a TDL, the function that finds it can
// indeed find it.
func Test_FetchTSPsInTDL(t *testing.T) {
	tdl, _ := testdatagen.MakeTDL(db, "source", "dest", "cos")
	tsp, _ := testdatagen.MakeTSP(db, "Test TSP", "TSP1")
	testdatagen.MakeBestValueScore(db, tsp, tdl, 10)

	tsps, err := models.FetchTSPsInTDLSortByAward(db, tdl.ID)

	if err != nil {
		t.Errorf("Failed to find TSP: %v", err)
	} else if len(tsps) != 1 {
		t.Errorf("Failed to find TSP. Expected to find 1, found %d", len(tsps))
	}
}

func Test_FetchTSPsInTDLByBVS(t *testing.T) {
	tdl, _ := testdatagen.MakeTDL(db, "source", "dest", "cos")
	tsp, _ := testdatagen.MakeTSP(db, "Test TSP", "TSP1")
	testdatagen.MakeBestValueScore(db, tsp, tdl, 10)

	tsps, err := models.FetchTSPsInTDLSortByBVS(db, tdl.ID)

	if err != nil {
		t.Errorf("Failed to find TSPs: %v", err)
	} else if len(tsps) != 1 {
		t.Errorf("Failed to find TSP. Expected to find 1, found %d", len(tsps))
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

	// Make 10 (not divisible by 4) TSPs in this TDL with BVSs
	for i := 0; i < tspsToMake; i++ {
		tsp, _ := testdatagen.MakeTSP(db, "Test Shipper", "TEST")
		testdatagen.MakeBestValueScore(db, tsp, tdl, 10)
	}
	// Fetch TSPs in TDL
	tspsbb, err := models.FetchTSPsInTDLSortByBVS(db, tdl.ID)
	qbs := assignTSPsToBands(tspsbb)
	if err != nil {
		t.Errorf("Failed to find TSPs: %v", err)
	}
	if len(qbs[0]) != 2 || len(qbs[2]) != 1 {
		t.Errorf("Failed to correctly add TSPs to quality bands.")
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
