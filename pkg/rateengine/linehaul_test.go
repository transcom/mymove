package rateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *RateEngineSuite) Test_CheckDetermineMileage() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger, suite.date)
	mileage, err := engine.determineMileage("10024", "18209")
	if err != nil {
		t.Error("Unable to determine mileage: ", err)
	}
	expected := 1000
	if mileage != expected {
		t.Errorf("Determined mileage incorrectly. Expected %d, got %d", expected, mileage)
	}
}

func (suite *RateEngineSuite) Test_CheckBaseLinehaul() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger, suite.date)
	mileage := 3200
	cwt := 40

	blh, _ := engine.baseLinehaul(mileage, cwt)
	expected := 128000
	if blh != 128000 {
		t.Errorf("CWT should have been %d but is %d.", expected, blh)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulFactors() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger, suite.date)

	// Load fake data
	defaultRateDateLower := time.Date(2017, 5, 15, 0, 0, 0, 0, time.UTC)
	defaultRateDateUpper := time.Date(2018, 5, 15, 0, 0, 0, 0, time.UTC)

	originZip3 := models.Tariff400ngZip3{
		Zip3:          395,
		BasepointCity: "Saucier",
		State:         "MS",
		ServiceArea:   428,
		RateArea:      48,
		Region:        11,
	}
	suite.mustSave(&originZip3)

	serviceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        428,
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		EffectiveDateLower: defaultRateDateLower,
		EffectiveDateUpper: defaultRateDateUpper,
	}
	suite.mustSave(&serviceArea)
	fmt.Print(serviceArea)
	fmt.Print(originZip3)

	linehaulFactor, err := engine.linehaulFactors(60, 395)
	if err != nil {
		t.Error("Unable to determine linehaulFactor: ", err)
	}
	expected := 3060
	if linehaulFactor != expected {
		t.Errorf("Determined linehaul factor incorrectly. Expected %d, got %d", expected, linehaulFactor)
	}
}

func (suite *RateEngineSuite) Test_CheckShorthaulCharge() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger, suite.date)
	mileage := 799
	cwt := 40

	shc, _ := engine.shorthaulCharge(mileage, cwt)
	expected := 31960
	if shc != expected {
		t.Errorf("Shorthaul charge should have been %d, but is %d.", expected, shc)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulChargeTotal() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger, suite.date)
	linehaulChargeTotal, err := engine.linehaulChargeTotal("10024", "94103")
	if err != nil {
		t.Error("Unable to determine linehaulChargeTotal: ", err)
	}
	expected := 13003
	if linehaulChargeTotal != expected {
		t.Errorf("Determined linehaul factor incorrectly. Expected %d, got %d", expected, linehaulChargeTotal)
	}
}
