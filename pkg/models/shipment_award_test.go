package models

import "testing"

func Test_ShipmentAwardValidations(t *testing.T) {
	sa := &ShipmentAward{}

	var expErrors = map[string][]string{
		"shipment_id":                        []string{"ShipmentID can not be blank."},
		"transportation_service_provider_id": []string{"TransportationServiceProviderID can not be blank."},
	}

	verifyValidationErrors(sa, expErrors, t)
}
