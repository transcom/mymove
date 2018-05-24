package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) Test_CheckServiceFee() {
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
	suite.mustSave(&serviceArea)

	fee, err := engine.serviceFeeCents(unit.CWT(50), "395", testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Fatalf("failed to calculate service fee: %s", err)
	}

	expected := unit.Cents(17500)
	if fee != expected {
		t.Errorf("wrong service fee: expected %d, got %d", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_CheckFullPack() {
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

	serviceArea := models.Tariff400ngServiceArea{
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
	suite.mustSave(&serviceArea)

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     0,
		WeightLbsUpper:     16001,
		RateCents:          5429,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.mustSave(&fullPackRate)

	fee, err := engine.fullPackCents(unit.CWT(50), "395", testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Fatalf("failed to calculate full pack fee: %s", err)
	}

	expected := unit.Cents(271450)
	if fee != expected {
		t.Errorf("wrong full pack fee: expected %d, got %d", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_CheckFullUnpack() {
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

	serviceArea := models.Tariff400ngServiceArea{
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
	suite.mustSave(&serviceArea)

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

	fee, err := engine.fullUnpackCents(unit.CWT(50), "395", testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Fatalf("failed to calculate full unpack fee: %s", err)
	}

	expected := unit.Cents(27145)
	if fee != expected {
		t.Errorf("wrong full unpack fee: expected %d, got %d", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_SITCharge() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger, suite.planner)

	cwt := unit.CWT(10)
	daysInSIT := 4
	zip3 := "395"
	sit185ARate := unit.Cents(2324)
	sit185BRate := unit.Cents(431)

	z := models.Tariff400ngZip3{
		Zip3:          zip3,
		BasepointCity: "Saucier",
		State:         "MS",
		ServiceArea:   "428",
		RateArea:      "US48",
		Region:        "11",
	}
	suite.mustSave(&z)

	sa := models.Tariff400ngServiceArea{
		Name:               "Tampa, FL",
		ServiceArea:        "428",
		LinehaulFactor:     69,
		ServiceChargeCents: 663,
		ServicesSchedule:   1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   sit185ARate,
		SIT185BRateCents:   sit185BRate,
		SITPDSchedule:      1,
	}
	suite.mustSave(&sa)

	// Test PPM SIT charges
	charge, err := engine.SitCharge(cwt, daysInSIT, zip3, testdatagen.DateInsidePeakRateCycle, true)
	if err != nil {
		t.Fatalf("error calculating SIT charge: %s", err)
	}
	expected := sit185BRate.Multiply(daysInSIT).Multiply(cwt.Int())
	if charge != expected {
		t.Errorf("wrong PPM SIT charge total: expected %d, got %d", expected, charge)
	}

	// Test HHG SIT charges
	charge, err = engine.SitCharge(cwt, daysInSIT, zip3, testdatagen.DateInsidePeakRateCycle, false)
	if err != nil {
		t.Fatalf("error calculating SIT charge: %s", err)
	}
	expectedFirstDay := sit185ARate.Multiply(cwt.Int()).Int()
	expectedAddtlDay := sit185BRate.Multiply(daysInSIT - 1).Multiply(cwt.Int()).Int()
	expected = unit.Cents(expectedFirstDay + expectedAddtlDay)
	if charge != expected {
		t.Errorf("wrong HHG SIT charge total: expected %d, got %d", expected, charge)
	}
}

func (suite *RateEngineSuite) Test_CheckNonLinehaulChargeTotal() {
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
		SIT185ARateCents:   unit.Cents(1750),
		SIT185BRateCents:   unit.Cents(50),
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

	cost, err := engine.nonLinehaulChargeComputation(
		unit.Pound(2000), "39503", "33607", testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Fatalf("failed to calculate non linehaul charge: %s", err)
	}

	// origin service fee:  7000
	// dest. service fee:  13260
	// pack fee:          108580
	// unpack fee:         10858
	expected := unit.Cents(139698)
	totalFee := cost.OriginServiceFee + cost.DestinationServiceFee + cost.PackFee + cost.UnpackFee
	if totalFee != expected {
		t.Errorf("wrong non-linehaul charge total: expected %d, got %d", expected, totalFee)
	}
}
