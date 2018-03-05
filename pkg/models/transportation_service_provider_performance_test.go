package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_BestValueScoreValidations() {
	tspPerformance := &TransportationServiceProviderPerformance{BestValueScore: 101}

	var expErrors = map[string][]string{
		"best_value_score": []string{"101 is not less than 101."},
	}

	suite.verifyValidationErrors(tspPerformance, expErrors)

	tspPerformance = &TransportationServiceProviderPerformance{BestValueScore: -1}

	expErrors = map[string][]string{
		"best_value_score": []string{"-1 is not greater than -1."},
	}

	suite.verifyValidationErrors(tspPerformance, expErrors)
}
