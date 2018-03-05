package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_TrafficDistributionList() {
	tdl := &TrafficDistributionList{}

	var expErrors = map[string][]string{
		"source_rate_area":   []string{"SourceRateArea can not be blank."},
		"destination_region": []string{"DestinationRegion can not be blank."},
		"code_of_service":    []string{"CodeOfService can not be blank."},
	}

	suite.verifyValidationErrors(tdl, expErrors)
}
