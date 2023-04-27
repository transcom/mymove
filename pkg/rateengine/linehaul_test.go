package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) BeforeTest(suiteName string, testName string) {
	suite.T().Skip()
}

func (suite *RateEngineSuite) Test_CheckBaseLinehaul() {
	move := models.Move{
		Locator: "ABC123",
	}

	t := suite.T()
	engine := NewRateEngine(move)

	expected := unit.Cents(128000)

	mileage := 3200
	weight := unit.Pound(3900)
	date := testdatagen.DateInsidePeakRateCycle

	blh, err := engine.baseLinehaul(suite.AppContextForTest(), mileage, weight, date)
	if blh != expected {
		t.Errorf("BaseLinehaulCents should have been %d but is %d.", expected, blh)
	}
	if err != nil {
		t.Errorf("Encountered error trying to get baseLinehaul: %v", err)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulFactors() {
	move := models.Move{
		Locator: "ABC123",
	}
	t := suite.T()
	engine := NewRateEngine(move)

	linehaulFactor, err := engine.linehaulFactors(suite.AppContextForTest(), 60, "395", testdatagen.RateEngineDate)
	if err != nil {
		t.Error("Unable to determine linehaulFactor: ", err)
	}
	expected := unit.Cents(3420)
	if linehaulFactor != expected {
		t.Errorf("Determined linehaul factor incorrectly. Expected %d, got %d", expected, linehaulFactor)
	}
}

func (suite *RateEngineSuite) Test_CheckShorthaulCharge() {
	move := models.Move{
		Locator: "ABC123",
	}
	t := suite.T()
	engine := NewRateEngine(move)
	mileage := 799
	cwt := unit.CWT(40)
	rate := unit.Cents(5656)

	shc, _ := engine.shorthaulCharge(suite.AppContextForTest(), mileage, cwt, testdatagen.DateInsidePeakRateCycle)
	if shc != rate {
		t.Errorf("Shorthaul charge should have been %d, but is %d.", rate, shc)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulChargeTotal() {
	move := models.Move{
		Locator: "ABC123",
	}
	t := suite.T()
	engine := NewRateEngine(move)
	weight := unit.Pound(2000)
	expected := unit.Cents(11462)
	zip5Austin := "78717"
	zip5SanFrancisco := "94103"
	distanceMiles := 1234

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	cost, err := engine.linehaulChargeComputation(
		suite.AppContextForTest(), weight, zip5Austin, zip5SanFrancisco, distanceMiles, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Error("Unable to determine linehaulChargeTotal: ", err)
	}
	if cost.LinehaulChargeTotal != expected {
		t.Errorf("Determined linehaul factor incorrectly. Expected %d, got %d", expected, cost.LinehaulChargeTotal)
	}
}

func (suite *RateEngineSuite) Test_CheckFuelSurchargeComputation() {
	move := models.Move{
		Locator: "ABC123",
	}
	engine := NewRateEngine(move)

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	fuelSurcharge, err := engine.fuelSurchargeComputation(suite.AppContextForTest(), unit.Cents(12000), testdatagen.NonPeakRateCycleEnd)

	suite.NoError(err)
	suite.Equal(unit.Cents(720), fuelSurcharge.Fee)
}
