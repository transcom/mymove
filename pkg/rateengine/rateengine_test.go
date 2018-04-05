package rateengine

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

func (suite *RateEngineSuite) Test_CheckDetermineCWT() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	weight := 2500
	cwt := engine.determineCWT(weight)

	if cwt != 25 {
		t.Errorf("CWT should have been 25 but is %d.", cwt)
	}
}

type RateEngineSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *RateEngineSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *RateEngineSuite) mustSave(model interface{}) {
	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		suite.T().Fatalf("errors encountered saving %v: %v", model, verrs)
	}
}

func TestRateEngineSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &RateEngineSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
