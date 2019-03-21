package models_test

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/dates"
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
	calendar := dates.NewUSCalendar()
	weekendDate := dates.NextNonWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 25, 0, 0, 0, 0, time.UTC))

	shipment := &Shipment{
		EstimatedPackDays:           &packDays,
		EstimatedTransitDays:        &transitDays,
		WeightEstimate:              &weightEstimate,
		ProgearWeightEstimate:       &progearWeightEstimate,
		SpouseProgearWeightEstimate: &spouseProgearWeightEstimate,
		RequestedPickupDate:         &weekendDate,
		OriginalDeliveryDate:        &weekendDate,
		OriginalPackDate:            &weekendDate,
		PmSurveyPlannedPackDate:     &weekendDate,
		PmSurveyPlannedPickupDate:   &weekendDate,
		PmSurveyPlannedDeliveryDate: &weekendDate,
		ActualPackDate:              &weekendDate,
		ActualPickupDate:            &weekendDate,
		ActualDeliveryDate:          &weekendDate,
	}

	stringDate := weekendDate.Format("2006-01-02 15:04:05 -0700 UTC")
	expErrors := map[string][]string{
		"move_id":                         []string{"move_id can not be blank."},
		"status":                          []string{"status can not be blank."},
		"estimated_pack_days":             []string{"-2 is less than or equal to zero."},
		"estimated_transit_days":          []string{"0 is less than or equal to zero."},
		"weight_estimate":                 []string{"-3 is less than zero."},
		"progear_weight_estimate":         []string{"-12 is less than zero."},
		"spouse_progear_weight_estimate":  []string{"-9 is less than zero."},
		"requested_pickup_date":           []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"original_delivery_date":          []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"original_pack_date":              []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"pm_survey_planned_pack_date":     []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"pm_survey_planned_pickup_date":   []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"pm_survey_planned_delivery_date": []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"actual_pack_date":                []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"actual_pickup_date":              []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"actual_delivery_date":            []string{fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
	}

	suite.verifyValidationErrors(shipment, expErrors)
}

// Test_FetchUnofferedShipments tests that a shipment is returned when we fetch shipments with offers.
func (suite *ModelSuite) Test_FetchUnofferedShipments() {
	t := suite.T()
	calendar := dates.NewUSCalendar()
	pickupDate := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 28, 0, 0, 0, 0, time.UTC))
	deliveryDate := dates.NextWorkday(*calendar, pickupDate)
	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	market := "dHHG"
	sourceGBLOC := "KKFA"
	destinationGBLOC := "HAFC"

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			ActualPickupDate:        &pickupDate,
			ActualDeliveryDate:      &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			DestinationGBLOC:        &destinationGBLOC,
			Market:                  &market,
			Status:                  ShipmentStatusSUBMITTED,
		},
	})

	shipment2 := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
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
	tspp := testdatagen.MakeDefaultTSPPerformance(suite.DB())
	CreateShipmentOffer(suite.DB(), shipment.ID, tspp.TransportationServiceProviderID, tspp.ID, false)
	shipments, err := FetchUnofferedShipments(suite.DB())

	// Expect only unassigned shipment returned
	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if len(shipments) != 1 {
		t.Errorf("Returned too many shipments. Expected %v, got %v", shipment2.ID, shipments)
	}
}

// TestShipmentStateMachine takes the shipment through valid state transitions
func (suite *ModelSuite) TestShipmentStateMachine() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
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

	// Can approve shipment (HHG)
	err = shipment.Approve()
	suite.Nil(err)
	suite.Equal(ShipmentStatusAPPROVED, shipment.Status, "expected Approved")

	shipDate := time.Now()

	// Can pack shipment
	err = shipment.Pack(shipDate)
	suite.Nil(err)
	suite.Equal(ShipmentStatusAPPROVED, shipment.Status, "expected Approved")
	suite.Equal(*shipment.ActualPackDate, shipDate, "expected Actual Pack Date to be set")

	// Can transport shipment
	err = shipment.Transport(shipDate)
	suite.Nil(err)
	suite.Equal(ShipmentStatusINTRANSIT, shipment.Status, "expected In Transit")
	suite.Equal(*shipment.ActualPickupDate, shipDate, "expected Actual Pickup Date to be set")

	// Can deliver shipment
	err = shipment.Deliver(shipDate)
	suite.Nil(err)
	suite.Equal(ShipmentStatusDELIVERED, shipment.Status, "expected Delivered")
	suite.Equal(*shipment.ActualDeliveryDate, shipDate, "expected Actual Delivery Date to be set")

	// Can complete shipment
	err = shipment.Complete()
	suite.Nil(err)
	suite.Equal(ShipmentStatusCOMPLETED, shipment.Status, "expected Completed")
}

func (suite *ModelSuite) TestSetBookDateWhenSubmitted() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	// There is not a way to set a field to nil using testdatagen.Assertions
	shipment.BookDate = nil
	suite.MustSave(&shipment)
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
	tspUsers, shipments, shipmentOffers, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	shipmentOffer := shipmentOffers[0]

	suite.Equal(ShipmentStatusAWARDED, shipment.Status, "expected Awarded")
	suite.Nil(shipmentOffer.Accepted)
	suite.Nil(shipmentOffer.RejectionReason)

	newShipment, newShipmentOffer, _, err := AcceptShipmentForTSP(suite.DB(), tspUser.TransportationServiceProviderID, shipment.ID)
	suite.NoError(err)

	suite.Equal(ShipmentStatusACCEPTED, newShipment.Status, "expected Accepted")
	suite.True(*newShipmentOffer.Accepted)
	suite.Nil(newShipmentOffer.RejectionReason)
	suite.NotEqual(shipment.Move.Orders.NewDutyStation.Address.ID, newShipment.DestinationAddressOnAcceptance.ID)
	suite.Equal(shipment.Move.Orders.NewDutyStation.Address.City, newShipment.DestinationAddressOnAcceptance.City)
}

func createAndAcceptShipmentWithDeliveryAddress(suite *ModelSuite, hasDeliveryAddress bool) (Shipment, Shipment, error) {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []ShipmentStatus{ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	unitedStates := "United States"

	addressAssertions := testdatagen.Assertions{
		Address: Address{
			StreetAddress1: "Fort Gordon",
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "30813",
			Country:        &unitedStates,
		},
	}

	deliveryAddress := testdatagen.MakeAddress3(suite.DB(), addressAssertions)
	shipment.HasDeliveryAddress = hasDeliveryAddress
	shipment.DeliveryAddress = &deliveryAddress
	shipment.DeliveryAddressID = &deliveryAddress.ID
	suite.DB().ValidateAndSave(&shipment)

	newShipment, _, _, err := AcceptShipmentForTSP(suite.DB(), tspUser.TransportationServiceProviderID, shipment.ID)

	return shipment, *newShipment, err
}

// TestAcceptShipmentForTSPWithDeliveryAddress tests that delivery address is used for a shipment when TSP accepts
// a offer and delivery address is available instead of duty station
func (suite *ModelSuite) TestAcceptShipmentForTSPWithDeliveryAddress() {
	hasDeliveryAddress := true
	shipment, newShipment, err := createAndAcceptShipmentWithDeliveryAddress(suite, hasDeliveryAddress)
	suite.NoError(err)
	suite.Equal(shipment.DeliveryAddress.City, newShipment.DestinationAddressOnAcceptance.City)
}

// TestAcceptShipmentForTSPWithDeliveryAddress tests that delivery address is used for a shipment when TSP accepts
// a offer and delivery address is available instead of duty station
func (suite *ModelSuite) TestAcceptShipmentForTSPWithDeliveryAddressHasDeliveryAddressFalse() {
	hasDeliveryAddress := false
	shipment, newShipment, err := createAndAcceptShipmentWithDeliveryAddress(suite, hasDeliveryAddress)
	suite.NoError(err)
	suite.Equal(shipment.Move.Orders.NewDutyStation.Address.City, newShipment.DestinationAddressOnAcceptance.City)
}

// TestCurrentTransportationServiceProviderID tests that a shipment returns the proper current tsp id
func (suite *ModelSuite) TestCurrentTransportationServiceProviderID() {
	tsp := testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{})
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	var emptyUUID uuid.UUID

	suite.Equal(shipment.CurrentTransportationServiceProviderID(), emptyUUID)

	testdatagen.MakeShipmentOffer(suite.DB(), testdatagen.Assertions{
		ShipmentOffer: ShipmentOffer{
			TransportationServiceProviderID: tsp.ID,
			ShipmentID:                      shipment.ID,
		},
	})

	// CurrentTransportationServiceProviderID looks at the shipment offers on a shipment
	// Since it doesn't re-fetch the shipment, if the offers have changed
	// We need to re-fetch the shipment to reload the offers
	reloadShipment, err := FetchShipmentByTSP(suite.DB(), tsp.ID, shipment.ID)
	suite.Nil(err)
	suite.Equal(tsp.ID, reloadShipment.CurrentTransportationServiceProviderID(), "expected ids to be equal")
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
		shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
			Shipment: Shipment{
				SourceGBLOC: &d[0],
			},
		})
		err := shipment.AssignGBLNumber(suite.DB())
		suite.NoError(err)
		suite.NotNil(shipment.GBLNumber)
		suite.Equal(*shipment.GBLNumber, d[1])
	}
}

// TestCreateShipmentLineItem tests that a shipment line item is created correctly
func (suite *ModelSuite) TestCreateShipmentLineItem() {
	acc := testdatagen.MakeDefaultTariff400ngItem(suite.DB())
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	q1 := unit.BaseQuantityFromInt(5)
	notes := "It's a giant moose head named Fred he seemed rather pleasant"
	baseParams := BaseShipmentLineItemParams{
		Tariff400ngItemID: acc.ID,
		Quantity1:         &q1,
		Location:          "ORIGIN",
		Notes:             &notes,
	}
	additionalParams := AdditionalShipmentLineItemParams{}
	shipmentLineItem, verrs, err := shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantityFromInt(5), shipmentLineItem.Quantity1)
		suite.Equal(acc.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
	}
}

// TestCreateShipmentLineItemCode105BAndE tests that 105B/E line items are created correctly
func (suite *ModelSuite) TestCreateShipmentLineItemCode105BAndE() {
	acc105B := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "105B",
		},
	})

	acc105E := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "105E",
		},
	})
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	notes := "It's a giant moose head named Fred he seemed rather pleasant"
	baseParams := BaseShipmentLineItemParams{
		Tariff400ngItemID:   acc105B.ID,
		Tariff400ngItemCode: acc105B.Code,
		Location:            "ORIGIN",
		Notes:               &notes,
	}
	additionalParams := AdditionalShipmentLineItemParams{
		ItemDimensions: &AdditionalLineItemDimensions{
			Length: 10000,
			Width:  10000,
			Height: 10000,
		},
		CrateDimensions: &AdditionalLineItemDimensions{
			Length: 10000,
			Width:  10000,
			Height: 10000,
		},
	}
	// Create 105B preapproval
	shipmentLineItem, verrs, err := shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	// 10x10x10 cubic inches is roughly 0.5787 cubic feet.
	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantity(5787), shipmentLineItem.Quantity1)
		suite.Equal(acc105B.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
		suite.NotZero(shipmentLineItem.ItemDimensions.ID)
		suite.NotZero(shipmentLineItem.CrateDimensions.ID)
	}

	//Create 105E preapproval
	baseParams.Tariff400ngItemID = acc105E.ID
	baseParams.Tariff400ngItemCode = acc105E.Code
	shipmentLineItem, verrs, err = shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantity(5787), shipmentLineItem.Quantity1)
		suite.Equal(acc105E.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
		suite.NotZero(shipmentLineItem.ItemDimensions.ID)
		suite.NotZero(shipmentLineItem.CrateDimensions.ID)
	}

	//Create 105E preapproval with base quantity
	baseParams.Tariff400ngItemID = acc105E.ID
	baseParams.Tariff400ngItemCode = acc105E.Code
	var q1 unit.BaseQuantity = 1000
	baseParams.Quantity1 = &q1
	additionalParams = AdditionalShipmentLineItemParams{}
	shipmentLineItem, verrs, err = shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantity(1000), shipmentLineItem.Quantity1)
		suite.Equal(acc105E.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
		suite.Zero(shipmentLineItem.ItemDimensionsID)
		suite.Zero(shipmentLineItem.CrateDimensionsID)
	}
}

// TestCreateShipmentLineItemCode35A tests that 35A line items are created correctly
func (suite *ModelSuite) TestCreateShipmentLineItemCode35A() {
	acc35A := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "35A",
		},
	})

	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	desc := "This is a description"
	reas := "This is the reason"
	estAmt := unit.Cents(1234)
	actAmt := unit.Cents(1000)
	baseParams := BaseShipmentLineItemParams{
		Tariff400ngItemID:   acc35A.ID,
		Tariff400ngItemCode: acc35A.Code,
		Location:            "ORIGIN",
	}
	additionalParams := AdditionalShipmentLineItemParams{
		Description:         &desc,
		Reason:              &reas,
		EstimateAmountCents: &estAmt,
		ActualAmountCents:   &actAmt,
	}

	// Create 35A preapproval
	shipmentLineItem, verrs, err := shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(acc35A.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
		suite.Equal(desc, *shipmentLineItem.Description)
		suite.Equal(reas, *shipmentLineItem.Reason)
		suite.Equal(estAmt, *shipmentLineItem.EstimateAmountCents)
		suite.Equal(actAmt, *shipmentLineItem.ActualAmountCents)
		suite.Equal(unit.BaseQuantity(100000), shipmentLineItem.Quantity1)
	}
}

// TestCreateShipmentLineItemCode226A tests that 226A line items are created correctly
func (suite *ModelSuite) TestCreateShipmentLineItemCode226A() {
	acc226A := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "226A",
		},
	})

	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	desc := "This is a description"
	reas := "This is the reason"
	actAmt := unit.Cents(1000)
	baseParams := BaseShipmentLineItemParams{
		Tariff400ngItemID:   acc226A.ID,
		Tariff400ngItemCode: acc226A.Code,
		Location:            "ORIGIN",
	}
	additionalParams := AdditionalShipmentLineItemParams{
		Description:       &desc,
		Reason:            &reas,
		ActualAmountCents: &actAmt,
	}

	// Create 226A preapproval
	shipmentLineItem, verrs, err := shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantityFromCents(actAmt), shipmentLineItem.Quantity1)
		suite.Equal(acc226A.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
		suite.Equal(desc, *shipmentLineItem.Description)
		suite.Equal(reas, *shipmentLineItem.Reason)
		suite.Equal(actAmt, *shipmentLineItem.ActualAmountCents)
	}
}

// TestUpdateShipmentLineItem tests that line items are updated correctly
func (suite *ModelSuite) TestUpdateShipmentLineItem() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	notes := "It's a giant moose head named Fred he seemed rather pleasant"
	description := "This is a description."
	acc4A := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "4A",
		},
	})
	lineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: ShipmentLineItem{
			Tariff400ngItem:   acc4A,
			Tariff400ngItemID: acc4A.ID,
			Quantity1:         unit.BaseQuantityFromInt(1234),
			Location:          "ORIGIN",
			Notes:             notes,
			Description:       &description,
		},
	})

	updateNotes := "Updated notes"
	baseParams := BaseShipmentLineItemParams{
		Quantity1:           &lineItem.Quantity1,
		Tariff400ngItemID:   lineItem.Tariff400ngItemID,
		Tariff400ngItemCode: lineItem.Tariff400ngItem.Code,
		Location:            string(lineItem.Location),
		Notes:               &updateNotes,
	}
	additionalParams := AdditionalShipmentLineItemParams{Description: lineItem.Description}

	// Create 105B preapproval
	verrs, err := shipment.UpdateShipmentLineItem(suite.DB(),
		baseParams, additionalParams, &lineItem)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantityFromInt(1234), lineItem.Quantity1)
		suite.Equal(*baseParams.Notes, lineItem.Notes)
	}
}

// TestUpdateShipmentLineItemCode105BAndE tests that 105B/E line items are updated correctly
func (suite *ModelSuite) TestUpdateShipmentLineItemCode105BAndE() {
	acc105B := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "105B",
		},
	})

	acc105E := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "105E",
		},
	})
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	notes := "It's a giant moose head named Fred he seemed rather pleasant"
	description := "This is a description."
	item := testdatagen.MakeDefaultShipmentLineItemDimensions(suite.DB())
	crate := testdatagen.MakeDefaultShipmentLineItemDimensions(suite.DB())
	lineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: ShipmentLineItem{
			Tariff400ngItemID: acc105B.ID,
			Location:          "ORIGIN",
			Notes:             notes,
			Description:       &description,
			ItemDimensionsID:  &item.ID,
			ItemDimensions:    item,
			CrateDimensionsID: &crate.ID,
			CrateDimensions:   crate,
		},
	})

	updateNotes := "Updated notes"
	baseParams := BaseShipmentLineItemParams{
		Tariff400ngItemID:   lineItem.Tariff400ngItemID,
		Tariff400ngItemCode: lineItem.Tariff400ngItem.Code,
		Location:            string(lineItem.Location),
		Notes:               &updateNotes,
	}
	additionalParams := AdditionalShipmentLineItemParams{
		ItemDimensions: &AdditionalLineItemDimensions{
			Length: unit.ThousandthInches(20000),
			Width:  unit.ThousandthInches(20000),
			Height: unit.ThousandthInches(20000),
		},
		CrateDimensions: &AdditionalLineItemDimensions{
			Length: unit.ThousandthInches(20000),
			Width:  unit.ThousandthInches(20000),
			Height: unit.ThousandthInches(20000),
		},
		Description: lineItem.Description,
	}

	// Create 105B preapproval
	verrs, err := shipment.UpdateShipmentLineItem(suite.DB(),
		baseParams, additionalParams, &lineItem)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantity(46296), lineItem.Quantity1)
		suite.Equal(acc105B.ID.String(), lineItem.Tariff400ngItem.ID.String())
		suite.Equal(*baseParams.Notes, lineItem.Notes)
		suite.Equal(*additionalParams.Description, *lineItem.Description)

		suite.NotZero(lineItem.ItemDimensions.ID)
		suite.Equal(additionalParams.ItemDimensions.Length, lineItem.ItemDimensions.Length)
		suite.Equal(additionalParams.ItemDimensions.Width, lineItem.ItemDimensions.Width)
		suite.Equal(additionalParams.ItemDimensions.Height, lineItem.ItemDimensions.Height)

		suite.NotZero(lineItem.CrateDimensions.ID)
		suite.Equal(additionalParams.CrateDimensions.Height, lineItem.CrateDimensions.Height)
		suite.Equal(additionalParams.CrateDimensions.Width, lineItem.CrateDimensions.Width)
		suite.Equal(additionalParams.CrateDimensions.Height, lineItem.CrateDimensions.Height)
	}

	//Update to 105E preapproval
	baseParams.Tariff400ngItemID = acc105E.ID
	baseParams.Tariff400ngItemCode = acc105E.Code
	verrs, err = shipment.UpdateShipmentLineItem(suite.DB(),
		baseParams, additionalParams, &lineItem)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantity(46296), lineItem.Quantity1)
		suite.Equal(acc105E.ID.String(), lineItem.Tariff400ngItem.ID.String())
		suite.NotZero(lineItem.ItemDimensions.ID)
		suite.Equal(additionalParams.ItemDimensions.Length, lineItem.ItemDimensions.Length)
		suite.Equal(additionalParams.ItemDimensions.Width, lineItem.ItemDimensions.Width)
		suite.Equal(additionalParams.ItemDimensions.Height, lineItem.ItemDimensions.Height)
		suite.NotZero(lineItem.CrateDimensions.ID)
		suite.Equal(additionalParams.CrateDimensions.Height, lineItem.CrateDimensions.Height)
		suite.Equal(additionalParams.CrateDimensions.Width, lineItem.CrateDimensions.Width)
		suite.Equal(additionalParams.CrateDimensions.Height, lineItem.CrateDimensions.Height)
	}
}

// TestUpdateShipmentLineItemCode35A tests that 35A line items are updated correctly
func (suite *ModelSuite) TestUpdateShipmentLineItemCode35A() {
	acc35A := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "35A",
		},
	})

	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	desc := "This is a description"
	reas := "This is the reason"
	notes := "Notes"
	loc := ShipmentLineItemLocationORIGIN
	estAmt := unit.Cents(1000)
	actAmt := unit.Cents(1000)
	lineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: ShipmentLineItem{
			Tariff400ngItemID:   acc35A.ID,
			Location:            loc,
			Notes:               notes,
			Description:         &desc,
			Reason:              &reas,
			EstimateAmountCents: &estAmt,
			ActualAmountCents:   &actAmt,
		},
	})

	// Update values
	baseParams := BaseShipmentLineItemParams{
		Tariff400ngItemID:   acc35A.ID,
		Tariff400ngItemCode: acc35A.Code,
		Location:            "ORIGIN",
	}
	desc = "updated description"
	reas = "updated reason"
	estAmt = unit.Cents(2000)
	actAmt = unit.Cents(1500)
	additionalParams := AdditionalShipmentLineItemParams{
		Description:         &desc,
		Reason:              &reas,
		EstimateAmountCents: &estAmt,
		ActualAmountCents:   &actAmt,
	}

	verrs, err := shipment.UpdateShipmentLineItem(suite.DB(),
		baseParams, additionalParams, &lineItem)
	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantity(150000), lineItem.Quantity1)
		suite.Equal(acc35A.ID.String(), lineItem.Tariff400ngItem.ID.String())
		suite.Equal(desc, *lineItem.Description)
		suite.Equal(reas, *lineItem.Reason)
		suite.Equal(estAmt, *lineItem.EstimateAmountCents)
		suite.Equal(actAmt, *lineItem.ActualAmountCents)
	}
}

// TestCreateShipmentLineItemCode105BAndEMissingDimensions tests that missing dimensions for 105B/E throws error
func (suite *ModelSuite) TestCreateShipmentLineItemCode105BAndEMissingDimensions() {
	acc105B := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "105B",
		},
	})

	acc105E := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "105E",
		},
	})
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	notes := "It's a giant moose head named Fred he seemed rather pleasant"
	baseParams := BaseShipmentLineItemParams{
		Tariff400ngItemID:   acc105B.ID,
		Tariff400ngItemCode: acc105B.Code,
		Location:            "ORIGIN",
		Notes:               &notes,
	}
	additionalParams := AdditionalShipmentLineItemParams{}
	// Try create 105B preapproval
	_, _, err := shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	suite.Error(err)

	// Try create 105E preapproval
	baseParams.Tariff400ngItemID = acc105E.ID
	baseParams.Tariff400ngItemCode = acc105E.Code
	_, _, err = shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	suite.Error(err)
}

// TestSaveShipmentAndPricingInfo tests that a shipment and line items can be saved
func (suite *ModelSuite) TestSaveShipmentAndPricingInfo() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	distance := testdatagen.MakeDefaultDistanceCalculation(suite.DB())

	var lineItems []ShipmentLineItem
	codes := []string{"LHS", "135A", "135B", "105A", "105C"}
	for _, code := range codes {
		item := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
			Tariff400ngItem: Tariff400ngItem{
				Code: code,
			},
		})
		lineItem := ShipmentLineItem{
			ShipmentID:        shipment.ID,
			Tariff400ngItemID: item.ID,
			Tariff400ngItem:   item,
			Location:          ShipmentLineItemLocationDESTINATION,
			Status:            ShipmentLineItemStatusAPPROVED,
		}
		lineItems = append(lineItems, lineItem)
	}

	verrs, err := shipment.SaveShipmentAndPricingInfo(suite.DB(), lineItems, []ShipmentLineItem{}, distance)
	suite.NoError(err)
	suite.NoVerrs(verrs)
}

// TestSaveShipmentAndPricingInfoDisallowDuplicates tests that duplicate baseline charges with the same
// tariff 400ng codes cannot be saved.
func (suite *ModelSuite) TestSaveShipmentAndPricingInfoDisallowBaselineDuplicates() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	var lineItems []ShipmentLineItem

	distance := testdatagen.MakeDefaultDistanceCalculation(suite.DB())

	item := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "LHS",
		},
	})
	testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: ShipmentLineItem{
			Tariff400ngItem:   item,
			ShipmentID:        shipment.ID,
			Tariff400ngItemID: item.ID,
			Shipment:          shipment,
		},
	})
	lineItem := ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Tariff400ngItemID: item.ID,
		Tariff400ngItem:   item,
	}
	lineItems = append(lineItems, lineItem)
	verrs, err := shipment.SaveShipmentAndPricingInfo(suite.DB(), lineItems, []ShipmentLineItem{}, distance)

	suite.Error(err)
	suite.NoVerrs(verrs)
}

// TestSaveShipmentAndPricingInfoDisallowDuplicates tests that duplicate baseline charges with the same
// tariff 400ng codes cannot be saved.
func (suite *ModelSuite) TestSaveShipmentAndPricingInfoAllowOtherDuplicates() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	var lineItems []ShipmentLineItem

	distance := testdatagen.MakeDefaultDistanceCalculation(suite.DB())

	item := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: Tariff400ngItem{
			Code: "105B",
		},
	})
	testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: ShipmentLineItem{
			Tariff400ngItem:   item,
			ShipmentID:        shipment.ID,
			Tariff400ngItemID: item.ID,
			Shipment:          shipment,
		},
	})

	lineItem := ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Tariff400ngItemID: item.ID,
		Tariff400ngItem:   item,
		Location:          ShipmentLineItemLocationDESTINATION,
		Status:            ShipmentLineItemStatusAPPROVED,
	}
	lineItems = append(lineItems, lineItem)
	verrs, err := shipment.SaveShipmentAndPricingInfo(suite.DB(), []ShipmentLineItem{}, lineItems, distance)

	suite.NoError(err)
	suite.NoVerrs(verrs)
}

// TestSaveShipment tests that a shipment can be saved
func (suite *ModelSuite) TestSaveShipment() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	now := time.Now()
	shipment.PmSurveyCompletedAt = &now

	verrs, err := SaveShipment(suite.DB(), &shipment)

	suite.NoError(err)
	suite.NoVerrs(verrs)
}

// TestAcceptedShipmentOffer test that we can retrieve a valid accepted shipment offer
func (suite *ModelSuite) TestAcceptedShipmentOffer() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	suite.Equal(ShipmentStatusDRAFT, shipment.Status, "expected Draft")

	// Shipment does not have an accepted shipment offer
	noAcceptedShipmentOffer, err := shipment.AcceptedShipmentOffer()
	suite.Nil(err) // Shipment.Status does not require an accepted ShipmentOffer
	suite.Nil(noAcceptedShipmentOffer)

	shipmentOffer := testdatagen.MakeDefaultShipmentOffer(suite.DB())
	shipment.ShipmentOffers = append(shipment.ShipmentOffers, shipmentOffer)
	suite.Len(shipment.ShipmentOffers, 1)

	// Can submit shipment
	err = shipment.Submit()
	suite.Nil(err)
	suite.Equal(ShipmentStatusSUBMITTED, shipment.Status, "expected Submitted")

	// Can award shipment
	err = shipment.Award()
	suite.Nil(err)
	suite.Equal(ShipmentStatusAWARDED, shipment.Status, "expected Awarded")

	// ShipmentOffer has not been accepted yet
	// Shipment does not have an accepted shipment offer
	noAcceptedShipmentOffer, err = shipment.AcceptedShipmentOffer()
	suite.Nil(err) // Shipment.Status does not require an accepted ShipmentOffer
	suite.Nil(noAcceptedShipmentOffer)

	// Can accept shipment
	err = shipment.Accept()
	suite.Nil(err)
	suite.Equal(ShipmentStatusACCEPTED, shipment.Status, "expected Accepted")

	// ShipmentOffer has not been accepted yet
	// Shipment does not have an accepted shipment offer, but Shipment is in the Accepted state
	noAcceptedShipmentOffer, err = shipment.AcceptedShipmentOffer()
	suite.NotNil(err) // Shipment.Status requires an accepted ShipmentOffer
	suite.Nil(noAcceptedShipmentOffer)

	// Accept ShipmentOffer for the TSP
	err = shipment.ShipmentOffers[0].Accept()
	suite.Nil(err)
	suite.True(*shipment.ShipmentOffers[0].Accepted)
	suite.Nil(shipment.ShipmentOffers[0].RejectionReason)

	// Get accepted shipment offer from shipment
	acceptedShipmentOffer, err := shipment.AcceptedShipmentOffer()
	suite.Nil(err)
	suite.NotNil(acceptedShipmentOffer)

	// Test results of TSP for an accepted shipment offer
	// accepted shipment offer can't have empty or nil values for certain data
	scac, err := acceptedShipmentOffer.SCAC()
	suite.Nil(err)
	suite.NotEmpty(scac)
	supplierID, err := acceptedShipmentOffer.SupplierID()
	suite.Nil(err)
	suite.NotNil(supplierID)
	suite.NotEmpty(*supplierID)

	// Do TSPs have the same ID
	suite.NotEmpty(acceptedShipmentOffer.TransportationServiceProviderPerformance.TransportationServiceProvider.ID.String())
	suite.Equal(acceptedShipmentOffer.TransportationServiceProviderPerformance.TransportationServiceProvider.ID,
		shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.TransportationServiceProvider.ID)
}
