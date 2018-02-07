package models

import (
	"fmt"
	"testing"

	"github.com/markbates/pop"
	"github.com/satori/go.uuid"
)

var db *pop.Connection

func Test_AwardedShipment(t *testing.T) {
	awardedShipment := AwardedShipment{
		ShipmentID:                      uuid.Must(uuid.NewV4()),
		TransportationServiceProviderID: uuid.Must(uuid.NewV4()),
		AdministrativeShipment:          false,
	}

	db.ValidateAndSave(&awardedShipment)
	fmt.Printf(">1> %s\n", awardedShipment)
	/*
		fmt.Printf(">2> %s\n", awardedShipmentSaved)
		if err != nil {
			t.Fatal("Didn't write it to the db")
		}
		if awardedShipment.ID == uuid.Nil {
			t.Error("didn't get an ID back")
		}

		if awardedShipment.CreatedAt.IsZero() {
			t.Error("wasn't assigned a created_at time")
		}
	*/
}

func Test_AwardedShipmentValidations(t *testing.T) {
	as := &AwardedShipment{}
	verrs, err := db.ValidateAndSave(as)
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
