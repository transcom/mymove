package models_test

import (
	"testing"

	. "github.com/transcom/mymove/pkg/models"
)

func Test_TransportationServiceProvider(t *testing.T) {
	tsp := &TransportationServiceProvider{}

	var expErrors = map[string][]string{
		"name": []string{"Name can not be blank."},
		"standard_carrier_alpha_code": []string{"StandardCarrierAlphaCode can not be blank."},
	}

	verifyValidationErrors(tsp, expErrors, t)
}
