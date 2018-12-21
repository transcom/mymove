package dpsapi

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/handlers"
	"go.uber.org/zap"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	handlers.BaseTestSuite
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &HandlerSuite{}
	hs.SetTestDB(db)
	hs.SetTestLogger(logger)

	suite.Run(t, hs)
}
