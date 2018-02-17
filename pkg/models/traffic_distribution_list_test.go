package models

import "testing"

func Test_TrafficDistributionList(t *testing.T) {
	tdl := &TrafficDistributionList{}

	var expErrors = map[string][]string{
		"source_rate_area":   []string{"SourceRateArea can not be blank."},
		"destination_region": []string{"DestinationRegion can not be blank."},
		"code_of_service":    []string{"CodeOfService can not be blank."},
	}

	verifyValidationErrors(tdl, expErrors, t)
}
