package models

import "testing"

func Test_BestValueScoreValidations(t *testing.T) {
	tspPerformance := &TransportationServiceProviderPerformance{BestValueScore: 101}

	var expErrors = map[string][]string{
		"best_value_score": []string{"101 is not less than 101."},
	}

	verifyValidationErrors(tspPerformance, expErrors, t)

	tspPerformance = &TransportationServiceProviderPerformance{BestValueScore: -1}

	expErrors = map[string][]string{
		"best_value_score": []string{"-1 is not greater than -1."},
	}

	verifyValidationErrors(tspPerformance, expErrors, t)
}
