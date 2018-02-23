package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var testDbConnection *pop.Connection
var testLogger *zap.Logger

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) SetupTest() {
	testDbConnection.TruncateAll()
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func setupDependencies() {

	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	dbConnection, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}
	testDbConnection = dbConnection
	Init(dbConnection)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	testLogger = logger
}

func TestMain(m *testing.M) {
	setupDependencies()

	os.Exit(m.Run())
}
