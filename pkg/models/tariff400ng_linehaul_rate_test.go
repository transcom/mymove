package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
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
