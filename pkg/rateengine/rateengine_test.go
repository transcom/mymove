package rateengine

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) setupRateEngineTest() {
	originZip3 := models.Tariff400ngZip3{
		Zip3:          "395",
		BasepointCity: "Saucier",
		State:         "MS",
		ServiceArea:   "428",
		RateArea:      "US48",
		Region:        "11",
	}
	suite.MustSave(&originZip3)
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
	suite.MustSave(&originServiceArea)
	destinationZip3 := models.Tariff400ngZip3{
		Zip3:          "336",
		BasepointCity: "Tampa",
		State:         "FL",
		ServiceArea:   "197",
		RateArea:      "US4964400",
		Region:        "13",
	}
	suite.MustSave(&destinationZip3)
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
	suite.MustSave(&destinationServiceArea)
	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     0,
		WeightLbsUpper:     16001,
		RateCents:          5429,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&fullPackRate)
	fullUnpackRate := models.Tariff400ngFullUnpackRate{
		Schedule:           1,
		RateMillicents:     542900,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&fullUnpackRate)
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
	suite.MustSave(&newBaseLinehaul)
	shorthaul := models.Tariff400ngShorthaulRate{
		CwtMilesLower:      1,
		CwtMilesUpper:      50000,
		RateCents:          5656,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&shorthaul)

}

func (suite *RateEngineSuite) computePPMIncludingLHRates(originZip string, destinationZip string, weight unit.Pound, logger *zap.Logger, planner route.Planner) (CostComputation, error) {
	suite.setupRateEngineTest()
	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US48",
			DestinationRegion: "13",
			CodeOfService:     "2",
		},
	})
	tsp := testdatagen.MakeDefaultTSP(suite.DB())
	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50.0),
	}
	suite.MustSave(&tspPerformance)
	lhDiscount, _, err := models.PPMDiscountFetch(suite.DB(),
		logger,
		originZip,
		destinationZip, testdatagen.RateEngineDate,
	)
	suite.Require().Nil(err)
	engine := NewRateEngine(suite.DB(), logger)
	cost, err := engine.ComputePPM(
		weight,
		originZip,
		destinationZip,
		1044,
		testdatagen.RateEngineDate,
		0,
		lhDiscount,
		0,
	)
	suite.Require().Nil(err)
	suite.Require().True(cost.GCC > 0)
	return cost, err
}

func (suite *RateEngineSuite) Test_CheckPPMTotal() {
	suite.setupRateEngineTest()
	t := suite.T()

	engine := NewRateEngine(suite.DB(), suite.logger)

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	// 139698 +20000
	cost, err := engine.ComputePPM(2000, "39574", "33633", 1234, testdatagen.RateEngineDate,
		1, unit.DiscountRate(.6), unit.DiscountRate(.5))

	if err != nil {
		t.Fatalf("failed to calculate ppm charge: %s", err)
	}

	expected := unit.Cents(64887)
	if cost.GCC != expected {
		t.Errorf("wrong GCC: expected %d, got %d", expected, cost.GCC)
	}
}

func (suite *RateEngineSuite) TestComputePPMWithLHDiscount() {
	logger, _ := zap.NewDevelopment()
	planner := route.NewTestingPlanner(1234)
	originZip := "39574"
	destinationZip := "33633"
	weight := unit.Pound(2000)
	cost, err := suite.computePPMIncludingLHRates(originZip, destinationZip, weight, logger, planner)

	engine := NewRateEngine(suite.DB(), logger)
	ppmCost, err := engine.ComputePPMIncludingLHDiscount(
		weight,
		originZip,
		destinationZip,
		1044,
		testdatagen.RateEngineDate,
		0,
		0,
	)
	suite.Require().Nil(err)

	suite.True(ppmCost.GCC > 0)
	suite.Equal(ppmCost, cost)
}

type RateEngineSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *RateEngineSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestRateEngineSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger, _ := zap.NewDevelopment()

	hs := &RateEngineSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}
