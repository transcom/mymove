package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
	"go.uber.org/zap"
)

var testDbConnection *pop.Connection
var testLogger *zap.Logger

func setupDependencies() {

	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	dbConnection, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	testDbConnection = dbConnection
	testLogger, err = zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	Init(dbConnection)
}

func TestMain(m *testing.M) {
	setupDependencies()

	os.Exit(m.Run())
}
