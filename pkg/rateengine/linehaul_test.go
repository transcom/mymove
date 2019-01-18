package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) Test_CheckDetermineMileage() {
	t := suite.T()
	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)
	mileage, err := engine.determineMileage("39574", "33633")
	if err != nil {
		t.Error("Unable to determine mileage: ", err)
	}
	expected := 1234
	if mileage != expected {
		t.Errorf("Determined mileage incorrectly. Expected %d, got %d", expected, mileage)
	}
}

func (suite *RateEngineSuite) Test_CheckBaseLinehaul() {
	t := suite.T()
	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)

	expected := unit.Cents(128000)

	newBaseLinehaul := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 3101,
		DistanceMilesUpper: 3300,
		WeightLbsLower:     3000,
		WeightLbsUpper:     4000,
		RateCents:          expected,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}

	otherBaseLinehaul := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 3401,
		DistanceMilesUpper: 3500,
		WeightLbsLower:     3000,
		WeightLbsUpper:     4000,
		RateCents:          158000,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.NonPeakRateCycleStart,
		EffectiveDateUpper: testdatagen.NonPeakRateCycleEnd,
	}

	suite.MustSave(&newBaseLinehaul)
	suite.MustSave(&otherBaseLinehaul)

	mileage := 3200
	weight := unit.Pound(3900)
	date := testdatagen.DateInsidePeakRateCycle

	blh, err := engine.baseLinehaul(mileage, weight, date)
	if blh != expected {
		t.Errorf("BaseLinehaulCents should have been %d but is %d.", expected, blh)
	}
	if err != nil {
		t.Errorf("Encountered error trying to get baseLinehaul: %v", err)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulFactors() {
	t := suite.T()
	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)

	// Load fake data
	originZip3 := models.Tariff400ngZip3{
		Zip3:          "395",
		BasepointCity: "Saucier",
		State:         "MS",
		ServiceArea:   "428",
		RateArea:      "US48",
		Region:        "11",
	}
	suite.MustSave(&originZip3)

	serviceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        "428",
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.MustSave(&serviceArea)

	linehaulFactor, err := engine.linehaulFactors(60, "395", testdatagen.RateEngineDate)
	if err != nil {
		t.Error("Unable to determine linehaulFactor: ", err)
	}
	expected := unit.Cents(3420)
	if linehaulFactor != expected {
		t.Errorf("Determined linehaul factor incorrectly. Expected %d, got %d", expected, linehaulFactor)
	}
}

func (suite *RateEngineSuite) Test_CheckShorthaulCharge() {
	t := suite.T()
	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)
	mileage := 799
	cwt := unit.CWT(40)
	rate := unit.Cents(5656)

	sh := models.Tariff400ngShorthaulRate{
		CwtMilesLower:      1,
		CwtMilesUpper:      50000,
		RateCents:          rate,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&sh)

	shc, _ := engine.shorthaulCharge(mileage, cwt, testdatagen.DateInsidePeakRateCycle)
	if shc != rate {
		t.Errorf("Shorthaul charge should have been %d, but is %d.", rate, shc)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulChargeTotal() {
	t := suite.T()
	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)
	weight := unit.Pound(2000)
	expected := unit.Cents(11462)
	zip3Austin := "787"
	zip5Austin := "78717"
	zip3SanFrancisco := "941"
	zip5SanFrancisco := "94103"

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	// $4642 is the 2018 baseline rate for a 1700 mile (Austin -> SF), 2000lb move
	newBaseLinehaul := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1,
		DistanceMilesUpper: 10000,
		WeightLbsLower:     1000,
		WeightLbsUpper:     4000,
		RateCents:          4642,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&newBaseLinehaul)

	// Create Service Area entries for Zip3s
	zipAustin := models.Tariff400ngZip3{
		Zip3:          zip3Austin,
		BasepointCity: "Austin",
		State:         "TX",
		ServiceArea:   "744",
		RateArea:      "US1",
		Region:        "1",
	}
	suite.MustSave(&zipAustin)

	zipSanFrancisco := models.Tariff400ngZip3{
		Zip3:          zip3SanFrancisco,
		BasepointCity: "San Francisco",
		State:         "CA",
		ServiceArea:   "81",
		RateArea:      "US2",
		Region:        "2",
	}
	suite.MustSave(&zipSanFrancisco)

	// Create fees for service areas
	sa1 := models.Tariff400ngServiceArea{
		Name:               "Austin",
		ServiceChargeCents: 100,
		ServiceArea:        "744",
		LinehaulFactor:     78,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.MustSave(&sa1)

	sa2 := models.Tariff400ngServiceArea{
		Name:               "SF",
		ServiceChargeCents: 200,
		ServiceArea:        "81",
		LinehaulFactor:     263,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.MustSave(&sa2)

	cost, err := engine.linehaulChargeComputation(
		weight, zip5Austin, zip5SanFrancisco, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Error("Unable to determine linehaulChargeTotal: ", err)
	}
	if cost.LinehaulChargeTotal != expected {
		t.Errorf("Determined linehaul factor incorrectly. Expected %d, got %d", expected, cost.LinehaulChargeTotal)
	}
}

func (suite *RateEngineSuite) Test_CheckFuelSurchargeComputation() {
	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	fuelSurcharge, err := engine.fuelSurchargeComputation(unit.Cents(12000), testdatagen.NonPeakRateCycleEnd)

	suite.Nil(err)
	suite.Equal(unit.Cents(720), fuelSurcharge.Fee)
}
