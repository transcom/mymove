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

func (suite *RateEngineSuite) Test_CheckDetermineMileage() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage, err := engine.determineMileage("10024", "18209")
	if err != nil {
		t.Error("Unable to determine mileage: ", err)
	}
	if mileage != 1000 {
		t.Errorf("Determined mileage incorrectly. Expected 1000 got %d", mileage)
	}
}

func (suite *RateEngineSuite) Test_CheckBaseLinehaul() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage := 3200
	weight := 4000

	blh, _ := engine.baseLinehaul(mileage, weight)

	if blh != 12800000 {
		t.Errorf("CWT should have been 12800000 but is %d.", blh)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulFactors() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	linehaulFactor, err := engine.linehaulFactors(6000, "18209")
	if err != nil {
		t.Error("Unable to determine linehaulFactor: ", err)
	}
	if linehaulFactor != 30.6 {
		t.Errorf("Determined linehaul factor incorrectly. Expected 30.6 got %f", linehaulFactor)
	}
}

func (suite *RateEngineSuite) Test_CheckShorthaulCharge() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage := 799
	cwt := 40

	shc, _ := engine.shorthaulCharge(mileage, cwt)

	if shc != 31960 {
		t.Errorf("Shorthaul charge should have been 31960 but is %f.", shc)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulChargeTotal() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	linehaulChargeTotal, err := engine.linehaulChargeTotal("10024", "94103")
	if err != nil {
		t.Error("Unable to determine linehaulChargeTotal: ", err)
	}
	if linehaulChargeTotal != 1180012.036 {
		t.Errorf("Determined linehaul factor incorrectly. Expected 1180012.036 got %f", linehaulChargeTotal)
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
