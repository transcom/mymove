package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) SetupTest() {
	dbConnection.TruncateAll()
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func setupDBConnection() {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	dbConnection, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	Init(dbConnection)
}

func TestMain(m *testing.M) {
	setupDBConnection()

	os.Exit(m.Run())
}
