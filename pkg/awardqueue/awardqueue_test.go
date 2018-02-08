package awardqueue

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
)

func TestFindAllShipments(t *testing.T) {
	err := findAllShipments()

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
