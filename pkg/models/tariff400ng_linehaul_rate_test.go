package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_LinehaulEffectiveDateValidation() {
	now := time.Now()

	validLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validLinehaulRate, expErrors)

	invalidLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"effective_date_upper": []string{"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidLinehaulRate, expErrors)
}

func (suite *ModelSuite) Test_LinehaulWeightValidation() {
	validLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validLinehaulRate, expErrors)

	invalidLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     200,
		WeightLbsUpper:     100,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors = map[string][]string{
		"weight_lbs_lower": []string{"200 is not less than 100."},
	}
	suite.verifyValidationErrors(&invalidLinehaulRate, expErrors)
}

func (suite *ModelSuite) Test_LinehaulRateValidation() {
	validLinehaulRate := Tariff400ngLinehaulRate{
		RateCents:          100,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validLinehaulRate, expErrors)

	invalidLinehaulRate := Tariff400ngLinehaulRate{
		RateCents:          -1,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors = map[string][]string{
		"rate_cents": []string{"-1 is not greater than -1."},
	}
	suite.verifyValidationErrors(&invalidLinehaulRate, expErrors)
}

func (suite *ModelSuite) Test_LinehaulDistanceValidation() {
	validLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validLinehaulRate, expErrors)

	invalidLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 200,
		DistanceMilesUpper: 100,
	}

	expErrors = map[string][]string{
		"distance_miles_lower": []string{"200 is not less than 100."},
	}
	suite.verifyValidationErrors(&invalidLinehaulRate, expErrors)
}

func (suite *ModelSuite) Test_LinehaulRateCreateAndSave() {
	t := suite.T()
	mySpecificRate := 474747

	newBaseLinehaul := Tariff400ngLinehaulRate{
		DistanceMilesLower: 3101,
		DistanceMilesUpper: 3300,
		WeightLbsLower: 3000,
		WeightLbsUpper: 4000,
		RateCents: mySpecificRate,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}

	suite.mustSave(&newBaseLinehaul)

	linehaulRate := 0

	sql := `SELECT
			rate_cents
		FROM
			tariff400ng_linehaul_rates
		WHERE
			rate_cents = $1;`

	err := suite.db.RawQuery(sql, mySpecificRate).First(&linehaulRate)

	if err != nil {
		t.Errorf("You got an error: %s\n", err)
	}
	if linehaulRate != mySpecificRate {
		t.Errorf("The record object didn't save!")
	}
}
