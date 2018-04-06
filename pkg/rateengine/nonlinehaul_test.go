package rateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *RateEngineSuite) Test_CheckServiceFee() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger, suite.date)

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

	serviceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        428,
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		EffectiveDateLower: defaultRateDateLower,
		EffectiveDateUpper: defaultRateDateUpper,
	}
	suite.mustSave(&serviceArea)

	fee, err := engine.serviceFeeCents(50, 395)
	if err != nil {
		t.Fatalf("failed to calculate service fee: %s", err)
	}

	expected := 17500
	if fee != expected {
		t.Errorf("wrong service fee: expected %d, got %d", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_CheckFullPack() {
	t := suite.T()
	t.Skip("Not yet implemented")

	engine := NewRateEngine(suite.db, suite.logger, suite.date)
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

	serviceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        428,
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		ServicesSchedule:   1,
		EffectiveDateLower: defaultRateDateLower,
		EffectiveDateUpper: defaultRateDateUpper,
	}
	suite.mustSave(&serviceArea)

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     0,
		WeightLbsUpper:     16001,
		RateCents:          5429,
		EffectiveDateLower: defaultRateDateLower,
		EffectiveDateUpper: defaultRateDateUpper,
	}
	suite.mustSave(&fullPackRate)

	fee, err := engine.fullPackCents(50, 395)
	if err != nil {
		t.Fatalf("failed to calculate full pack fee: %s", err)
	}

	expected := 1375
	if fee != expected {
		t.Errorf("wrong full pack fee: expected %d, got %d", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_CheckFullUnpack() {
	t := suite.T()
	t.Skip("Not yet implemented")

	engine := NewRateEngine(suite.db, suite.logger, suite.date)

	fee, err := engine.fullUnpackCents(25, 18209)
	if err != nil {
		t.Fatalf("failed to calculate full unpack fee: %s", err)
	}

	expected := 125
	if fee != expected {
		t.Errorf("wrong full pack fee: expected %d, got %d", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_CheckNonLinehaulChargeTotal() {
	t := suite.T()
	t.Skip("Not yet implemented")

	engine := NewRateEngine(suite.db, suite.logger, suite.date)

	fee, err := engine.nonLinehaulChargeTotalCents(10024, 18209, 0.5)

	if err != nil {
		t.Fatalf("failed to calculate non linehaul charge: %s", err)
	}
	// (155.2 + 155.2 + 2200 + 200) * .5
	expected := 1355
	if fee != expected {
		t.Errorf("wrong non-linehaul charge total: expected %d, got %d", expected, fee)
	}
}
