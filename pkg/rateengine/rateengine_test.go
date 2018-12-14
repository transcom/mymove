package rateengine

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) Test_CheckPPMTotal() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger, suite.planner)
	originZip3 := models.Tariff400ngZip3{
		Zip3:          "395",
		BasepointCity: "Saucier",
		State:         "MS",
		ServiceArea:   "428",
		RateArea:      "US48",
		Region:        "11",
	}
	suite.mustSave(&originZip3)

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        "428",
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		ServicesSchedule:   1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.mustSave(&originServiceArea)

	destinationZip3 := models.Tariff400ngZip3{
		Zip3:          "336",
		BasepointCity: "Tampa",
		State:         "FL",
		ServiceArea:   "197",
		RateArea:      "US4964400",
		Region:        "13",
	}
	suite.mustSave(&destinationZip3)

	destinationServiceArea := models.Tariff400ngServiceArea{
		Name:               "Tampa, FL",
		ServiceArea:        "197",
		LinehaulFactor:     69,
		ServiceChargeCents: 663,
		ServicesSchedule:   1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(5550),
		SIT185BRateCents:   unit.Cents(222),
		SITPDSchedule:      1,
	}
	suite.mustSave(&destinationServiceArea)

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     0,
		WeightLbsUpper:     16001,
		RateCents:          5429,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.mustSave(&fullPackRate)

	fullUnpackRate := models.Tariff400ngFullUnpackRate{
		Schedule:           1,
		RateMillicents:     542900,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.mustSave(&fullUnpackRate)

	newBaseLinehaul := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1,
		DistanceMilesUpper: 10000,
		WeightLbsLower:     1000,
		WeightLbsUpper:     4000,
		RateCents:          20000,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.mustSave(&newBaseLinehaul)

	shorthaul := models.Tariff400ngShorthaulRate{
		CwtMilesLower:      1,
		CwtMilesUpper:      50000,
		RateCents:          5656,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.mustSave(&shorthaul)

	// 139698 +20000
	cost, err := engine.ComputePPM(2000, "39574", "33633", testdatagen.RateEngineDate,
		1, unit.DiscountRate(.6), unit.DiscountRate(.5))

	if err != nil {
		t.Fatalf("failed to calculate ppm charge: %s", err)
	}

	expected := unit.Cents(64887)
	if cost.GCC != expected {
		t.Errorf("wrong GCC: expected %d, got %d", expected, cost.GCC)
	}
}

type RateEngineSuite struct {
	testingsuite.LocalTestSuite
	db      *pop.Connection
	logger  *zap.Logger
	planner route.Planner
}

func (suite *RateEngineSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *RateEngineSuite) mustSave(model interface{}) {
	t := suite.T()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		t.Fatalf("error: %s", err)
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		t.Fatalf("errors encountered saving %v: %v", model, verrs)
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
	logger, _ := zap.NewDevelopment()
	planner := route.NewTestingPlanner(1234)

	hs := &RateEngineSuite{db: db, logger: logger, planner: planner}
	suite.Run(t, hs)
}
