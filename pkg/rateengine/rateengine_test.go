package rateengine

import (
	"log"
	"testing"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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

func (suite *RateEngineSuite) Test_CheckPPMTotal() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	defaultRateDateLower := time.Date(2017, 5, 15, 0, 0, 0, 0, time.UTC)
	defaultRateDateUpper := time.Date(2018, 5, 15, 0, 0, 0, 0, time.UTC)

	originZip3 := models.Tariff400ngZip3{
		Zip3:          395,
		BasepointCity: "Saucier",
		State:         "MS",
		ServiceArea:   428,
		RateArea:      "48",
		Region:        11,
	}
	suite.mustSave(&originZip3)

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        428,
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		ServicesSchedule:   1,
		EffectiveDateLower: defaultRateDateLower,
		EffectiveDateUpper: defaultRateDateUpper,
	}
	suite.mustSave(&originServiceArea)

	destinationZip3 := models.Tariff400ngZip3{
		Zip3:          336,
		BasepointCity: "Tampa",
		State:         "FL",
		ServiceArea:   197,
		RateArea:      "4964400",
		Region:        13,
	}
	suite.mustSave(&destinationZip3)

	destinationServiceArea := models.Tariff400ngServiceArea{
		Name:               "Tampa, FL",
		ServiceArea:        197,
		LinehaulFactor:     69,
		ServiceChargeCents: 663,
		ServicesSchedule:   1,
		EffectiveDateLower: defaultRateDateLower,
		EffectiveDateUpper: defaultRateDateUpper,
	}
	suite.mustSave(&destinationServiceArea)

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     0,
		WeightLbsUpper:     16001,
		RateCents:          5429,
		EffectiveDateLower: defaultRateDateLower,
		EffectiveDateUpper: defaultRateDateUpper,
	}
	suite.mustSave(&fullPackRate)

	fullUnpackRate := models.Tariff400ngFullUnpackRate{
		Schedule:           1,
		RateMillicents:     542900,
		EffectiveDateLower: defaultRateDateLower,
		EffectiveDateUpper: defaultRateDateUpper,
	}
	suite.mustSave(&fullUnpackRate)

	newBaseLinehaul := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1,
		DistanceMilesUpper: 10000,
		WeightLbsLower:     1000,
		WeightLbsUpper:     4000,
		RateCents:          20000,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.RateEngineDate,
		EffectiveDateUpper: testdatagen.RateEngineDate,
	}

	_, err := suite.db.ValidateAndSave(&newBaseLinehaul)
	if err != nil {
		t.Errorf("The newBaselineHaul didn't save.")
	}

	// 139698 +20000
	fee, err := engine.computePPM(2000, 395, 336, testdatagen.RateEngineDate, .40)

	if err != nil {
		t.Fatalf("failed to calculate ppm charge: %s", err)
	}

	expected := 61642
	if fee != expected {
		t.Errorf("wrong PPM charge total: expected %d, got %d", expected, fee)
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
