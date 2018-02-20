package awardqueue

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// newBoolPtr creates a pointer to a boolean value. Apparently it's not
// straightforward to do this in Go otherwise.
// https://stackoverflow.com/questions/28817992/how-to-set-bool-pointer-to-true-in-struct-literal
func newBoolPtr(value bool) *bool {
	return &value
}

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
		AdministrativeShipment:          newBoolPtr(false),
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
		AdministrativeShipment:          newBoolPtr(false),
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
