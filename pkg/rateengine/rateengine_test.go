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
	originZip503 := models.Tariff400ngZip3{
		Zip3:          "503",
		BasepointCity: "Des Moines",
		State:         "IA",
		ServiceArea:   "296",
		RateArea:      "US53",
		Region:        "7",
	}
	suite.MustSave(&originZip503)
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
	originServiceArea503 := models.Tariff400ngServiceArea{
		Name:               "Des Moines, IA",
		ServiceArea:        "296",
		LinehaulFactor:     unit.Cents(263),
		ServiceChargeCents: unit.Cents(489),
		ServicesSchedule:   3,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1447),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      3,
	}
	suite.MustSave(&originServiceArea503)
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
	itemRate210A := models.Tariff400ngItemRate{
		Code:               "210A",
		Schedule:           &destinationServiceArea.SITPDSchedule,
		WeightLbsLower:     newBaseLinehaul.WeightLbsLower,
		WeightLbsUpper:     newBaseLinehaul.WeightLbsUpper,
		RateCents:          unit.Cents(57600),
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&itemRate210A)
	itemRate225A := models.Tariff400ngItemRate{
		Code:               "225A",
		Schedule:           &destinationServiceArea.ServicesSchedule,
		WeightLbsLower:     newBaseLinehaul.WeightLbsLower,
		WeightLbsUpper:     newBaseLinehaul.WeightLbsUpper,
		RateCents:          unit.Cents(9900),
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&itemRate225A)
	fullPackRate1 := models.Tariff400ngFullPackRate{
		Schedule:           3,
		WeightLbsLower:     0,
		WeightLbsUpper:     16001,
		RateCents:          6130,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&fullPackRate1)
	tdl := models.TrafficDistributionList{
		SourceRateArea:    "US48",
		DestinationRegion: "13",
		CodeOfService:     "2",
	}
	suite.MustSave(&tdl)
	tdl1 := models.TrafficDistributionList{
		SourceRateArea:    "US53",
		DestinationRegion: "13",
		CodeOfService:     "2",
	}
	suite.MustSave(&tdl1)
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
	tspPerformance1 := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl1.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50.0),
	}
	suite.MustSave(&tspPerformance1)
}

func (suite *RateEngineSuite) computePPMIncludingLHRates(originZip string, destinationZip string, weight unit.Pound, logger Logger, planner route.Planner) (CostComputation, error) {
	move := models.Move{
		Locator: "ABC123",
	}
	lhDiscount, sitDiscount, err := models.PPMDiscountFetch(suite.DB(),
		logger,
		move,
		originZip,
		destinationZip,
		testdatagen.RateEngineDate,
	)
	suite.Require().Nil(err)
	engine := NewRateEngine(suite.DB(), logger, move)
	cost, err := engine.computePPM(
		weight,
		originZip,
		destinationZip,
		1044,
		testdatagen.RateEngineDate,
		0,
		lhDiscount,
		sitDiscount,
	)
	suite.Require().Nil(err)
	suite.Require().True(cost.GCC > 0)
	return cost, err
}

func (suite *RateEngineSuite) Test_CheckPPMTotal() {
	move := models.Move{
		Locator: "ABC123",
	}
	suite.setupRateEngineTest()
	t := suite.T()

	engine := NewRateEngine(suite.DB(), suite.logger, move)

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	// 139698 +20000
	cost, err := engine.computePPM(2000, "39574", "33633", 1234, testdatagen.RateEngineDate,
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
	move := models.Move{
		Locator: "ABC123",
	}
	suite.setupRateEngineTest()
	logger, _ := zap.NewDevelopment()
	planner := route.NewTestingPlanner(1234)
	originZip := "39574"
	destinationZip := "33633"
	weight := unit.Pound(2000)
	cost, err := suite.computePPMIncludingLHRates(originZip, destinationZip, weight, logger, planner)
	suite.Require().Nil(err)

	engine := NewRateEngine(suite.DB(), logger, move)
	ppmCost, err := engine.computePPMIncludingLHDiscount(
		weight,
		originZip,
		destinationZip,
		1044,
		testdatagen.RateEngineDate,
		0,
	)
	suite.Require().Nil(err)

	suite.True(ppmCost.GCC > 0)
	suite.Equal(ppmCost, cost)
}

func (suite *RateEngineSuite) TestComputeLowestCostPPMMove() {
	move := models.Move{
		Locator: "ABC123",
	}
	suite.setupRateEngineTest()
	logger, _ := zap.NewDevelopment()
	planner := route.NewTestingPlanner(1234)
	originZip := "39574"
	originDutyStationZip := "50309"
	destinationZip := "33633"
	distanceMilesFromOriginPickupZip := 1044
	distanceMilesFromDutyStationZip := 3300
	weight := unit.Pound(2000)
	engine := NewRateEngine(suite.DB(), logger, move)

	suite.Run("TestComputeLowestCostPPMMove when pickup zip results in lower GCC", func() {
		ppmCostWithPickupZip, err := suite.computePPMIncludingLHRates(
			originZip,
			destinationZip,
			weight,
			logger,
			planner,
		)
		suite.NoError(err)

		ppmCostWithDutyStationZip, err := suite.computePPMIncludingLHRates(
			originDutyStationZip,
			destinationZip,
			weight,
			logger,
			planner,
		)
		suite.NoError(err)

		cost, err := engine.ComputeLowestCostPPMMove(
			weight,
			originZip,
			originDutyStationZip,
			destinationZip,
			distanceMilesFromOriginPickupZip,
			distanceMilesFromDutyStationZip,
			testdatagen.RateEngineDate,
			0,
		)
		suite.NoError(err)

		suite.True(cost.GCC > 0)
		suite.True(ppmCostWithPickupZip.GCC > 0)
		suite.True(ppmCostWithDutyStationZip.GCC > 0)
		suite.True(ppmCostWithPickupZip.GCC < ppmCostWithDutyStationZip.GCC)
		suite.Equal(cost, ppmCostWithPickupZip)
	})

	suite.Run("TestComputeLowestCostPPMMove when duty station results in lower GCC", func() {
		originZip := "50309"
		originDutyStationZip := "39574"
		distanceMilesFromOriginPickupZip := 3300
		distanceMilesFromDutyStationZip := 1044

		ppmCostWithPickupZip, err := suite.computePPMIncludingLHRates(
			originZip,
			destinationZip,
			weight,
			logger,
			planner,
		)
		suite.NoError(err)

		ppmCostWithDutyStationZip, err := suite.computePPMIncludingLHRates(
			originDutyStationZip,
			destinationZip,
			weight,
			logger,
			planner,
		)
		suite.NoError(err)

		cost, err := engine.ComputeLowestCostPPMMove(
			weight,
			originZip,
			originDutyStationZip,
			destinationZip,
			distanceMilesFromOriginPickupZip,
			distanceMilesFromDutyStationZip,
			testdatagen.RateEngineDate,
			0,
		)
		suite.NoError(err)

		suite.True(cost.GCC > 0)
		suite.True(ppmCostWithPickupZip.GCC > 0)
		suite.True(ppmCostWithDutyStationZip.GCC > 0)
		suite.True(ppmCostWithPickupZip.GCC > ppmCostWithDutyStationZip.GCC)
		suite.Equal(cost, ppmCostWithDutyStationZip)
	})
}

type RateEngineSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *RateEngineSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestRateEngineSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger, _ := zap.NewDevelopment()

	hs := &RateEngineSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
