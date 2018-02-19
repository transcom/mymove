package awardqueue

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

var db *pop.Connection

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
func TestAwardShipment(t *testing.T) {
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
		TransportationServiceProviderID: &tsp.ID,
		AdministrativeShipment:          newBoolPtr(false),
	}

	// Run the Award Queue
	award, err := AttemptShipmentAward(pas)

	// See if shipment was awarded
	if err != nil {
		t.Fatalf("Shipment award expected no errors, received: %v", err)
	}
	if award == nil {
		t.Fatal("ShipmentAward was not found.")
	}
}

func setupDBConnection() {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	conn, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	db = conn
}

func TestMain(m *testing.M) {
	setupDBConnection()

	os.Exit(m.Run())
}
