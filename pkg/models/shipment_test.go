package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_ShipmentValidations() {
	packDays := int64(-2)
	transitDays := int64(0)
	var weightEstimate unit.Pound = -3
	var progearWeightEstimate unit.Pound = -12
	var spouseProgearWeightEstimate unit.Pound = -9

	shipment := &Shipment{
		EstimatedPackDays:           &packDays,
		EstimatedTransitDays:        &transitDays,
		WeightEstimate:              &weightEstimate,
		ProgearWeightEstimate:       &progearWeightEstimate,
		SpouseProgearWeightEstimate: &spouseProgearWeightEstimate,
	}

	expErrors := map[string][]string{
		"move_id":                        []string{"move_id can not be blank."},
		"status":                         []string{"status can not be blank."},
		"estimated_pack_days":            []string{"-2 is less than or equal to zero."},
		"estimated_transit_days":         []string{"0 is less than or equal to zero."},
		"weight_estimate":                []string{"-3 is less than or equal to zero."},
		"progear_weight_estimate":        []string{"-12 is less than or equal to zero."},
		"spouse_progear_weight_estimate": []string{"-9 is less than or equal to zero."},
	}

	suite.verifyValidationErrors(shipment, expErrors)
}

// Test_FetchAllShipments tests that a shipment is returned when we fetch shipments with their offers.
func (suite *ModelSuite) Test_FetchAllShipments() {
	t := suite.T()
	pickupDate := time.Now()
	deliveryDate := time.Now().AddDate(0, 0, 1)
	tdl := testdatagen.MakeDefaultTDL(suite.db)
	market := "dHHG"
	sourceGBLOC := "OHAI"

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			PickupDate:              &pickupDate,
			DeliveryDate:            &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})

	shipment2 := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			PickupDate:              &pickupDate,
			DeliveryDate:            &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})
	tsp := testdatagen.MakeDefaultTSP(suite.db)
	CreateShipmentOffer(suite.db, shipment.ID, tsp.ID, false)
	shipments, err := FetchShipments(suite.db, false)

	// Expect both shipments returned
	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if shipments[0].ID != shipment.ID || shipments[1].ID != shipment2.ID {
		t.Errorf("Failed to return correct shipments. Expected shipments %v and %v, got %v and %v",
			shipment.ID, shipment2.ID, shipments[0].ID, shipments[1].ID)
	}
}

// Test_FetchUnassignedShipments tests that a shipment is returned when we fetch shipments with offers.
func (suite *ModelSuite) Test_FetchUnassignedShipments() {
	t := suite.T()
	pickupDate := time.Now()
	deliveryDate := time.Now().AddDate(0, 0, 1)
	tdl := testdatagen.MakeDefaultTDL(suite.db)
	market := "dHHG"
	sourceGBLOC := "OHAI"

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			PickupDate:              &pickupDate,
			DeliveryDate:            &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})

	shipment2 := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			PickupDate:              &pickupDate,
			DeliveryDate:            &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})
	tsp := testdatagen.MakeDefaultTSP(suite.db)
	CreateShipmentOffer(suite.db, shipment.ID, tsp.ID, false)
	shipments, err := FetchShipments(suite.db, true)

	// Expect only unassigned shipment returned
	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if len(shipments) != 1 {
		t.Errorf("Returned too many shipments. Expected %v, got %v", shipment2.ID, shipments)
	}
}

func equalShipmentsSlice(a []ShipmentWithOffer, b []ShipmentWithOffer) bool {
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
