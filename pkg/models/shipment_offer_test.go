package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ShipmentOfferValidations() {
	sa := &ShipmentOffer{}

	var expErrors = map[string][]string{
		"shipment_id":                                    {"ShipmentID can not be blank."},
		"transportation_service_provider_id":             {"TransportationServiceProviderID can not be blank."},
		"transportation_service_provider_performance_id": {"TransportationServiceProviderPerformanceID can not be blank."},
	}

	suite.verifyValidationErrors(sa, expErrors)
}
