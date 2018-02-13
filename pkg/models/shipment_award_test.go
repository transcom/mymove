package models

import (
	"testing"

	"github.com/satori/go.uuid"
)

func Test_ShipmentAward(t *testing.T) {
	shipmentAward := ShipmentAward{
		ShipmentID:                      uuid.Must(uuid.NewV4()),
		TransportationServiceProviderID: uuid.Must(uuid.NewV4()),
		AdministrativeShipment:          false,
	}

	dbConnection.Create(&shipmentAward)

	// TODO: Complete tests here. As written, this will fail because we're
	// making up UUIDs for foreign keys and that will violate the foreign
	// key constraint! Hoist with his own petard.

	/*
		if err != nil {
			t.Fatal("Didn't write it to the db: ", err)
		}

		if shipmentAward.ID == uuid.Nil {
			t.Error("didn't get an ID back")
		}

		if shipmentAward.CreatedAt.IsZero() {
			t.Error("wasn't assigned a created_at time")
		}
	*/
}

func Test_ShipmentAwardValidations(t *testing.T) {
	as := &ShipmentAward{}
	verrs, err := dbConnection.ValidateAndSave(as)
	if err != nil {
		t.Error(err)
	}

	if verrs.Count() != 2 {
		t.Errorf("expected %d errors, got %d", 2, verrs.Count())
	}

	shipmentErrs := verrs.Get("shipment_id")
	expected := []string{"ShipmentID can not be blank."}
	if !equalSlice(shipmentErrs, expected) {
		t.Errorf("expected errors on %s to be %v, got %v", "ShipmentID", expected, shipmentErrs)
	}

	tspErrs := verrs.Get("transportation_service_provider_id")
	expected = []string{"TransportationServiceProviderID can not be blank."}
	if !equalSlice(tspErrs, expected) {
		t.Errorf("expected errors on %s to be %v, got %v", "TransportationServiceProviderID", expected, tspErrs)
	}
}

func equalSlice(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
