package rateengine

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	// . "github.com/transcom/mymove/pkg/rateengine"
)

func (suite *RateEngineSuite) Test_determineMileage() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage, err := engine.determineMileage("10024", "18209")
	if err != nil {
		t.Error("Unable to determine mileage: ", err)
	}
	if mileage != 1000 {
		t.Errorf("Determined mileage incorrectly. Expected 1000 got %v", mileage)
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

func (suite *RateEngineSuite) Test_determineCWT() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	weight := 2500
	cwt := engine.determineCWT(weight)

	if cwt != 25 {
		t.Errorf("CWT should have been 25 but is %d.", cwt)
	}
}

// // Determine the Base Linehaul (BLH)
// func (re *RateEngine) determineBaseLinehaul(mileage int, weight int) (base_linehaul_charge int, error) {
//  // TODO (Rebecca): This will come from a fetch
//  base_linehaul_charge := mileage * weight
//  // TODO (Rebecca): make a proper error
//  err := "whoops"
//  return base_linehaul_charge, err
// }

func (suite *RateEngineSuite) Test_CheckDetermineBaseLinehaul() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage := 3200
	weight := 4000

	blh, _ := engine.determineBaseLinehaul(mileage, weight)

	if blh != 12800000 {
		t.Errorf("CWT should have been 12800000 but is %d.", blh)
	}
}
