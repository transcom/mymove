package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_ShipmentValidations() {
	packDays := -2
	transitDays := 0
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
		"traffic_distribution_list_id":   []string{"traffic_distribution_list_id can not be blank."},
		"source_gbloc":                   []string{"source_gbloc can not be blank."},
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
	now := time.Now()
	tdl, _ := testdatagen.MakeTDL(
		suite.db,
		testdatagen.DefaultSrcRateArea,
		testdatagen.DefaultDstRegion,
		testdatagen.DefaultCOS)
	market := "dHHG"
	sourceGBLOC := "OHAI"
	shipment, _ := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl, sourceGBLOC, &market)
	shipment2, _ := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl, sourceGBLOC, &market)
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
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
	now := time.Now()
	tdl, _ := testdatagen.MakeTDL(
		suite.db,
		testdatagen.DefaultSrcRateArea,
		testdatagen.DefaultDstRegion,
		testdatagen.DefaultCOS)
	market := "dHHG"
	sourceGBLOC := "OHAI"
	shipment, _ := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl, sourceGBLOC, &market)
	shipment2, _ := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl, sourceGBLOC, &market)
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
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
