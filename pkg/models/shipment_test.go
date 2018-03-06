package models

import "testing"

func Test_Shipment(t *testing.T) {
	shipment := &Shipment{}

	expErrors := map[string][]string{
		"traffic_distribution_list_id": []string{"traffic_distribution_list_id can not be blank."},
		"gbloc": []string{"gbloc can not be blank."},
	}

	verifyValidationErrors(shipment, expErrors, t)
}
