package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) Test_CheckServiceFee() {
	move := models.Move{
		Locator: "ABC123",
	}
	t := suite.T()
	engine := NewRateEngine(move)

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
		ServicesSchedule:   1,
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.MustSave(&serviceArea)

	feeAndRate, err := engine.serviceFeeCents(suite.AppContextForTest(), unit.CWT(50), "395", testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		suite.Fail("failed to calculate service fee: %s", err)
	}

	expectedFee := unit.Cents(17500)
	if feeAndRate.Fee != expectedFee {
		t.Errorf("wrong service fee: expected %d, got %d", expectedFee, feeAndRate.Fee)
	}

	expectedRate := unit.Cents(350).ToMillicents()
	if feeAndRate.Rate != expectedRate {
		t.Errorf("wrong service rate: expected %d, got %d", expectedRate, feeAndRate.Rate)
	}
}

func (suite *RateEngineSuite) Test_CheckFullPack() {
	move := models.Move{
		Locator: "ABC123",
	}
	t := suite.T()
	engine := NewRateEngine(move)

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
		ServicesSchedule:   1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.MustSave(&serviceArea)

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     0,
		WeightLbsUpper:     16001,
		RateCents:          5429,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&fullPackRate)

	feeAndRate, err := engine.fullPackCents(suite.AppContextForTest(), unit.CWT(50), "395", testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		suite.Fail("failed to calculate full pack fee: %s", err)
	}

	expectedFee := unit.Cents(271450)
	if feeAndRate.Fee != expectedFee {
		t.Errorf("wrong full pack fee: expected %d, got %d", expectedFee, feeAndRate.Fee)
	}
	expectedRate := unit.Cents(5429).ToMillicents()
	if feeAndRate.Rate != expectedRate {
		t.Errorf("wrong full pack rate: expected %d, got %d", expectedRate, feeAndRate.Rate)
	}
}

func (suite *RateEngineSuite) Test_CheckFullUnpack() {
	move := models.Move{
		Locator: "ABC123",
	}
	t := suite.T()
	engine := NewRateEngine(move)

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
		ServicesSchedule:   1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.MustSave(&serviceArea)

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

	feeAndRate, err := engine.fullUnpackCents(suite.AppContextForTest(), unit.CWT(50), "395", testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		suite.Fail("failed to calculate full unpack fee: %s", err)
	}

	expected := unit.Cents(27145)
	if feeAndRate.Fee != expected {
		t.Errorf("wrong full unpack fee: expected %d, got %d", expected, feeAndRate.Fee)
	}

	expectedRate := unit.Millicents(542900)
	if feeAndRate.Rate != expectedRate {
		t.Errorf("wrong full unpack rate: expected %d, got %d", expectedRate, feeAndRate.Rate)
	}
}

func (suite *RateEngineSuite) Test_SITCharge() {
	var cwtAtMin, cwtBelowMin unit.CWT
	var expectedAtMinBase, expectedAtMin, expectedBelowMinBase, expectedBelowMin SITComputation

	move := models.Move{
		Locator: "ABC123",
	}
	engine := NewRateEngine(move)
	zip3 := "395"
	daysInSIT := 4

	suite.PreloadData(func() {

		z := models.Tariff400ngZip3{
			Zip3:          zip3,
			BasepointCity: "Saucier",
			State:         "MS",
			ServiceArea:   "428",
			RateArea:      "US48",
			Region:        "11",
		}
		suite.MustSave(&z)

		sit185ARate := unit.Cents(2324)
		sit185BRate := unit.Cents(431)
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
		suite.MustSave(&sa)

		sit210ARateCentsAtMin := unit.Cents(57600)
		itemRate210AAtMin := models.Tariff400ngItemRate{
			Code:               "210A",
			Schedule:           &sa.SITPDSchedule,
			WeightLbsLower:     unit.Pound(1000),
			WeightLbsUpper:     unit.Pound(1099),
			RateCents:          sit210ARateCentsAtMin,
			EffectiveDateLower: testdatagen.PeakRateCycleStart,
			EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		}
		suite.MustSave(&itemRate210AAtMin)

		sit225ARateCentsAtMin := unit.Cents(9900)
		itemRate225AAtMin := models.Tariff400ngItemRate{
			Code:               "225A",
			Schedule:           &sa.ServicesSchedule,
			WeightLbsLower:     itemRate210AAtMin.WeightLbsLower,
			WeightLbsUpper:     itemRate210AAtMin.WeightLbsUpper,
			RateCents:          sit225ARateCentsAtMin,
			EffectiveDateLower: testdatagen.PeakRateCycleStart,
			EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		}
		suite.MustSave(&itemRate225AAtMin)

		sit210ARateCentsBelowMin := unit.Cents(42100)
		itemRate210ABelowMin := models.Tariff400ngItemRate{
			Code:               "210A",
			Schedule:           &sa.SITPDSchedule,
			WeightLbsLower:     unit.Pound(0),
			WeightLbsUpper:     unit.Pound(999),
			RateCents:          sit210ARateCentsBelowMin,
			EffectiveDateLower: testdatagen.PeakRateCycleStart,
			EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		}
		suite.MustSave(&itemRate210ABelowMin)

		sit225ARateCentsBelowMin := unit.Cents(7700)
		itemRate225ABelowMin := models.Tariff400ngItemRate{
			Code:               "225A",
			Schedule:           &sa.ServicesSchedule,
			WeightLbsLower:     itemRate210ABelowMin.WeightLbsLower,
			WeightLbsUpper:     itemRate210ABelowMin.WeightLbsUpper,
			RateCents:          sit225ARateCentsBelowMin,
			EffectiveDateLower: testdatagen.PeakRateCycleStart,
			EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		}
		suite.MustSave(&itemRate225ABelowMin)

		cwtAtMin = unit.CWT(10)

		expectedAtMinBase.SITPart = sit185ARate.Multiply(cwtAtMin.Int())
		expectedAtMinBase.SITPart = expectedAtMinBase.SITPart.AddCents(sit185BRate.Multiply(daysInSIT - 1).Multiply(cwtAtMin.Int()))
		expectedAtMinBase.NonDiscountedTotal = expectedAtMinBase.SITPart.AddCents(expectedAtMinBase.LinehaulPart)
		expectedAtMin.SITPart = expectedAtMinBase.SITPart.AddCents(sit210ARateCentsAtMin)
		expectedAtMin.LinehaulPart = sit225ARateCentsAtMin
		expectedAtMin.NonDiscountedTotal = expectedAtMin.SITPart.AddCents(expectedAtMin.LinehaulPart)

		cwtBelowMin = unit.CWT(5)
		expectedBelowMinBase.SITPart = sit185ARate.Multiply(cwtBelowMin.Int())
		expectedBelowMinBase.SITPart = expectedBelowMinBase.SITPart.AddCents(sit185BRate.Multiply(daysInSIT - 1).Multiply(cwtBelowMin.Int()))
		expectedBelowMinBase.NonDiscountedTotal = expectedBelowMinBase.SITPart.AddCents(expectedBelowMinBase.LinehaulPart)
		expectedBelowMin.SITPart = expectedBelowMinBase.SITPart.AddCents(sit210ARateCentsBelowMin)
		expectedBelowMin.LinehaulPart = sit225ARateCentsBelowMin
		expectedBelowMin.NonDiscountedTotal = expectedBelowMin.SITPart.AddCents(expectedBelowMin.LinehaulPart)

	})

	// TODO: HHG SIT formula will be changing in future story to add in 225A/225B/225C (based on mileage).
	//   Current test just expecting baseline 185A and 185B charges.
	var testCases = []struct {
		description string
		cwt         unit.CWT
		isPPM       bool
		expected    SITComputation
	}{
		{"PPM at minimum weight", cwtAtMin, true, expectedAtMin},
		{"HHG at minimum weight", cwtAtMin, false, expectedAtMinBase},
		{"PPM below minimum weight", cwtBelowMin, true, expectedBelowMin},
		{"HHG below minimum weight", cwtBelowMin, false, expectedAtMinBase},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.description, func() {
			charge, err := engine.SitCharge(suite.AppContextForTest(), testCase.cwt, daysInSIT, zip3, testdatagen.DateInsidePeakRateCycle, testCase.isPPM)
			if err != nil {
				suite.Fail("error calculating SIT charge: %s", err)
			}
			suite.Equal(testCase.expected, charge)
		})
	}

	// Test zero days in SIT.
	suite.Run("zero days in SIT", func() {
		charge, err := engine.SitCharge(suite.AppContextForTest(), cwtAtMin, 0, zip3, testdatagen.DateInsidePeakRateCycle, true)
		suite.NoError(err)
		suite.Equal(SITComputation{}, charge)
	})

	// Test negative days in SIT.
	suite.Run("negative days in SIT", func() {
		_, err := engine.SitCharge(suite.AppContextForTest(), cwtAtMin, -1, zip3, testdatagen.DateInsidePeakRateCycle, true)
		suite.Error(err)
	})
}

func (suite *RateEngineSuite) Test_CheckNonLinehaulChargeTotal() {
	move := models.Move{
		Locator: "ABC123",
	}
	t := suite.T()
	engine := NewRateEngine(move)

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
		SIT185ARateCents:   unit.Cents(1750),
		SIT185BRateCents:   unit.Cents(50),
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

	cost, err := engine.nonLinehaulChargeComputation(
		suite.AppContextForTest(), unit.Pound(2000), "39503", "33607", testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		suite.Fail("failed to calculate non linehaul charge: %s", err)
	}

	// origin service fee:  7000
	// dest. service fee:  13260
	// pack fee:          108580
	// unpack fee:         10858
	expected := unit.Cents(139698)
	totalFee := cost.OriginService.Fee + cost.DestinationService.Fee + cost.Pack.Fee + cost.Unpack.Fee
	if totalFee != expected {
		t.Errorf("wrong non-linehaul charge total: expected %d, got %d", expected, totalFee)
	}
}

func (suite *RateEngineSuite) Test_SitComputationApplyDiscount() {
	var discountTestCases = []struct {
		description      string
		sitPart          unit.Cents
		linehaulPart     unit.Cents
		sitDiscount      unit.DiscountRate
		linehaulDiscount unit.DiscountRate
		expected         unit.Cents
	}{
		{
			"all values",
			unit.Cents(57500),
			unit.Cents(45555),
			unit.DiscountRate(0.5),
			unit.DiscountRate(0.75),
			unit.Cents(40139),
		},
		{
			"all zeros",
			unit.Cents(0),
			unit.Cents(0),
			unit.DiscountRate(0),
			unit.DiscountRate(0),
			unit.Cents(0),
		},
		{
			"no discount",
			unit.Cents(57500),
			unit.Cents(45555),
			unit.DiscountRate(0),
			unit.DiscountRate(0),
			unit.Cents(103055),
		},
		{
			"full discount",
			unit.Cents(57500),
			unit.Cents(45555),
			unit.DiscountRate(1),
			unit.DiscountRate(1),
			unit.Cents(0),
		},
	}

	for _, testCase := range discountTestCases {
		suite.Run(testCase.description, func() {
			sitComputation := SITComputation{
				SITPart:            testCase.sitPart,
				LinehaulPart:       testCase.linehaulPart,
				NonDiscountedTotal: testCase.sitPart.AddCents(testCase.linehaulPart),
			}

			sitDiscounted := sitComputation.ApplyDiscount(testCase.linehaulDiscount, testCase.sitDiscount)
			suite.Equal(testCase.expected, sitDiscounted)
		})
	}
}
