package models

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
)

var dbConnection *pop.Connection

func setupDBConnection() {
	configLocation := "../../config"
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
