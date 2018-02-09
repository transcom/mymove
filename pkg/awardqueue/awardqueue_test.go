package awardqueue

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
)

func TestFindAllUnawardedShipments(t *testing.T) {
	_, err := findAllUnawardedShipments()

	if err != nil {
		t.Fatal("Unable to find shipments: ", err)
	}
}

func setupDBConnection() {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	conn, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	dbConnection = conn
}

func TestMain(m *testing.M) {
	setupDBConnection()

	os.Exit(m.Run())
}
