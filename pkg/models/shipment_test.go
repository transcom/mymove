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

// Test_FetchUnofferedShipments tests that a shipment is returned when we fetch shipments with offers.
func (suite *ModelSuite) Test_FetchUnofferedShipments() {
	t := suite.T()
	pickupDate := time.Now()
	deliveryDate := time.Now().AddDate(0, 0, 1)
	tdl := testdatagen.MakeDefaultTDL(suite.db)
	market := "dHHG"
	sourceGBLOC := "OHAI"

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			ActualPickupDate:        &pickupDate,
			ActualDeliveryDate:      &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
			Status:                  ShipmentStatusSUBMITTED,
		},
	})

	shipment2 := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			ActualPickupDate:        &pickupDate,
			ActualDeliveryDate:      &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
			Status:                  ShipmentStatusSUBMITTED,
		},
	})
	tspp := testdatagen.MakeDefaultTSPPerformance(suite.db)
	CreateShipmentOffer(suite.db, shipment.ID, tspp.TransportationServiceProviderID, tspp.ID, false)
	shipments, err := FetchUnofferedShipments(suite.db)

	// Expect only unassigned shipment returned
	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if len(shipments) != 1 {
		t.Errorf("Returned too many shipments. Expected %v, got %v", shipment2.ID, shipments)
	}
}

// TestShipmentStateMachine takes the shipment through valid state transitions
func (suite *ModelSuite) TestShipmentStateMachine() {
	shipment := testdatagen.MakeDefaultShipment(suite.db)
	suite.Equal(ShipmentStatusDRAFT, shipment.Status, "expected Draft")

	// Can submit shipment
	err := shipment.Submit()
	suite.Nil(err)
	suite.Equal(ShipmentStatusSUBMITTED, shipment.Status, "expected Submitted")

	// Can award shipment
	err = shipment.Award()
	suite.Nil(err)
	suite.Equal(ShipmentStatusAWARDED, shipment.Status, "expected Awarded")

	// Can accept shipment
	err = shipment.Accept()
	suite.Nil(err)
	suite.Equal(ShipmentStatusACCEPTED, shipment.Status, "expected Accepted")

	// Can approve shipment
	err = shipment.Approve()
	suite.Nil(err)
	suite.Equal(ShipmentStatusAPPROVED, shipment.Status, "expected Approved")

	shipDate := time.Now()

	// Can transport shipment
	err = shipment.Transport(shipDate)
	suite.Nil(err)
	suite.Equal(ShipmentStatusINTRANSIT, shipment.Status, "expected In Transit")

	// Can deliver shipment
	err = shipment.Deliver(shipDate)
	suite.Nil(err)
	suite.Equal(ShipmentStatusDELIVERED, shipment.Status, "expected Delivered")

	// Can complete shipment
	err = shipment.Complete()
	suite.Nil(err)
	suite.Equal(ShipmentStatusCOMPLETED, shipment.Status, "expected Completed")
}

func (suite *ModelSuite) TestSetBookDateWhenSubmitted() {
	shipment := testdatagen.MakeDefaultShipment(suite.db)

	// There is not a way to set a field to nil using testdatagen.Assertions
	shipment.BookDate = nil
	suite.mustSave(&shipment)
	suite.Nil(shipment.BookDate)

	// Can submit shipment
	err := shipment.Submit()
	suite.Nil(err)
	suite.NotNil(shipment.BookDate)
}

// TestAcceptShipmentForTSP tests that a shipment and shipment offer is correctly accepted
func (suite *ModelSuite) TestAcceptShipmentForTSP() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []ShipmentStatus{ShipmentStatusAWARDED}
	tspUsers, shipments, shipmentOffers, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	shipmentOffer := shipmentOffers[0]

	suite.Equal(ShipmentStatusAWARDED, shipment.Status, "expected Awarded")
	suite.Nil(shipmentOffer.Accepted)
	suite.Nil(shipmentOffer.RejectionReason)

	newShipment, newShipmentOffer, _, err := AcceptShipmentForTSP(suite.db, tspUser.TransportationServiceProviderID, shipment.ID)
	suite.NoError(err)

	suite.Equal(ShipmentStatusACCEPTED, newShipment.Status, "expected Awarded")
	suite.True(*newShipmentOffer.Accepted)
	suite.Nil(newShipmentOffer.RejectionReason)
}

// TestShipmentAssignGBLNumber tests that a GBL number is created correctly
func (suite *ModelSuite) TestShipmentAssignGBLNumber() {
	testData := [][]string{
		// {GBLOC, expected GBL number}
		{"GBO1", "GBO17000001"},
		{"GBO1", "GBO17000002"},
		{"GBO1", "GBO17000003"},
		// New GBLOC starts new sequence
		{"GBO2", "GBO27000001"},
		// Old sequence should still work
		{"GBO1", "GBO17000004"},
	}

	for _, d := range testData {
		shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
			Shipment: Shipment{
				SourceGBLOC: &d[0],
			},
		})
		err := shipment.AssignGBLNumber(suite.db)
		suite.NoError(err)
		suite.NotNil(shipment.GBLNumber)
		suite.Equal(*shipment.GBLNumber, d[1])
	}
}

// TestShipmentAssignGBLNumber tests that a GBL number is created correctly
func (suite *ModelSuite) TestCreateShipmentAccessorial() {
	acc := testdatagen.MakeDummyAccessorial(suite.db)
	shipment := testdatagen.MakeDefaultShipment(suite.db)

	q1 := int64(5)
	notes := "It's a giant moose head named Fred he seemed rather pleasant"
	shipmentAccessorial, verrs, err := shipment.CreateShipmentAccessorial(suite.db,
		acc.ID,
		&q1,
		nil,
		"O",
		&notes)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(5, shipmentAccessorial.Quantity1.ToUnitInt())
		suite.Equal(acc.ID.String(), shipmentAccessorial.Accessorial.ID.String())
	}
}
