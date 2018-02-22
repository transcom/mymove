package models

import "testing"

func Test_BestValueScoreValidations(t *testing.T) {
	bvs := &BestValueScore{Score: 101}

	var expErrors = map[string][]string{
		"score": []string{"101 is not less than 101."},
	}

	verifyValidationErrors(bvs, expErrors, t)

	bvs = &BestValueScore{Score: -1}

	expErrors = map[string][]string{
		"score": []string{"-1 is not greater than -1."},
	}

	verifyValidationErrors(bvs, expErrors, t)
}
