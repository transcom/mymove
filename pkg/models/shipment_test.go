package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ShipmentValidations() {
	shipment := &Shipment{}

	expErrors := map[string][]string{
		"traffic_distribution_list_id": []string{"traffic_distribution_list_id can not be blank."},
	}

	suite.verifyValidationErrors(shipment, expErrors)
}
