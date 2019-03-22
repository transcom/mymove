package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchDataShipmentSummaryWorksheet() {
	moveID, _ := uuid.NewV4()
	serviceMemberID, _ := uuid.NewV4()
	//advanceID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	fortGordon := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	rank := models.ServiceMemberRankE9
	moveType := models.SelectedMoveTypeHHGPPM

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:               moveID,
			SelectedMoveType: &moveType,
		},
		Order: models.Order{
			OrdersType:       ordersType,
			NewDutyStationID: fortGordon.ID,
		},
		ServiceMember: models.ServiceMember{
			ID:            serviceMemberID,
			DutyStationID: &yuma.ID,
			Rank:          &rank,
		},
	})

	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID:              move.ID,
			NetWeight:           models.Int64Pointer(10000),
			HasRequestedAdvance: true,
			AdvanceID:           &advance.ID,
			Advance:             &advance,
		},
	})
	// Only concerned w/ approved advances for ssw
	ppm.Move.PersonallyProcuredMoves[0].Advance.Request()
	ppm.Move.PersonallyProcuredMoves[0].Advance.Approve()
	// Save advance in reimbursements table by saving ppm
	models.SavePersonallyProcuredMove(suite.DB(), &ppm)
	movedocuments := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   "OK",
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: serviceMemberID,
			ServiceMember:   move.Orders.ServiceMember,
		},
	}
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments)
	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			ServiceMemberID: serviceMemberID,
			MoveID:          move.ID,
		},
	})
	session := auth.Session{
		UserID:          move.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}
	ppm.Move.Submit()
	ppm.Move.Approve()
	// This is the same PPM model as ppm, but this is the one that will be saved by SaveMoveDependencies
	ppm.Move.PersonallyProcuredMoves[0].Submit()
	ppm.Move.PersonallyProcuredMoves[0].Approve()
	ppm.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(suite.DB(), &ppm.Move)
	ssd, err := models.FetchDataShipmentSummaryWorksheetFormData(suite.DB(), &session, moveID)

	suite.NoError(err)
	suite.Equal(move.Orders.ID, ssd.Order.ID)
	suite.Require().Len(ssd.Shipments, 1)
	suite.Equal(shipment.ID, ssd.Shipments[0].ID)
	suite.Require().Len(ssd.PersonallyProcuredMoves, 1)
	suite.Equal(ppm.ID, ssd.PersonallyProcuredMoves[0].ID)
	suite.Equal(serviceMemberID, ssd.ServiceMember.ID)
	suite.Equal(yuma.ID, ssd.CurrentDutyStation.ID)
	suite.Equal(yuma.Address.ID, ssd.CurrentDutyStation.Address.ID)
	suite.Equal(fortGordon.ID, ssd.NewDutyStation.ID)
	suite.Equal(fortGordon.Address.ID, ssd.NewDutyStation.Address.ID)
	rankWtgAllotment := models.GetWeightAllotment(rank)
	suite.Equal(rankWtgAllotment, ssd.WeightAllotment)
	suite.Require().NotNil(ssd.ServiceMember.Rank)
	totalEntitlement, err := models.GetEntitlement(*ssd.ServiceMember.Rank, ssd.Order.HasDependents, ssd.Order.SpouseHasProGear)
	suite.Require().Nil(err)
	suite.Equal(unit.Pound(totalEntitlement), ssd.TotalWeightAllotment)
	suite.Require().Len(ssd.MovingExpenseDocuments, 2)
	suite.NotNil(ssd.MovingExpenseDocuments[0].ID)
	suite.NotNil(ssd.MovingExpenseDocuments[1].ID)
	suite.Equal(ppm.NetWeight, ssd.PersonallyProcuredMoves[0].NetWeight)
	suite.Require().NotNil(ssd.PersonallyProcuredMoves[0].Advance)
	suite.Equal(ppm.Advance.ID, ssd.PersonallyProcuredMoves[0].Advance.ID)
	suite.Equal(unit.Cents(1000), ssd.PersonallyProcuredMoves[0].Advance.RequestedAmount)
}

func (suite *ModelSuite) TestFetchMovingExpensesShipmentSummaryWorksheetNoPPM() {
	moveID, _ := uuid.NewV4()
	serviceMemberID, _ := uuid.NewV4()
	moveType := models.SelectedMoveTypeHHG

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:               moveID,
			SelectedMoveType: &moveType,
		},
	})
	session := auth.Session{
		UserID:          move.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}

	movingExpenses, err := models.FetchMovingExpensesShipmentSummaryWorksheet(move, suite.DB(), &session)

	suite.Len(movingExpenses, 0)
	suite.Nil(err)
}

func (suite *ModelSuite) TestFormatValuesShipmentSummaryWorksheetFormPage1() {
	yuma := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	fortGordon := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	wtgEntitlements := models.WeightAllotment{
		TotalWeightSelf:               13000,
		TotalWeightSelfPlusDependents: 15000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}
	serviceMemberID, _ := uuid.NewV4()
	serviceBranch := models.AffiliationAIRFORCE
	rank := models.ServiceMemberRankE9
	serviceMember := models.ServiceMember{
		ID:            serviceMemberID,
		FirstName:     models.StringPointer("Marcus"),
		MiddleName:    models.StringPointer("Joseph"),
		LastName:      models.StringPointer("Jenkins"),
		Suffix:        models.StringPointer("Jr."),
		Telephone:     models.StringPointer("444-555-8888"),
		PersonalEmail: models.StringPointer("michael+ppm-expansion_1@truss.works"),
		Edipi:         models.StringPointer("1234567890"),
		Affiliation:   &serviceBranch,
		Rank:          &rank,
		DutyStationID: &yuma.ID,
	}

	orderIssueDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	order := models.Order{
		IssueDate:           orderIssueDate,
		OrdersType:          internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:        models.StringPointer("012345"),
		NewDutyStationID:    fortGordon.ID,
		OrdersIssuingAgency: models.StringPointer(string(serviceBranch)),
		TAC:                 models.StringPointer("NTA4"),
		SAC:                 models.StringPointer("SAC"),
		HasDependents:       true,
		SpouseHasProGear:    true,
	}
	pickupDate := time.Date(2019, time.January, 11, 0, 0, 0, 0, time.UTC)
	weight := unit.Pound(5000)
	shipments := []models.Shipment{
		{
			ActualPickupDate: &pickupDate,
			NetWeight:        &weight,
			Status:           models.ShipmentStatusDELIVERED,
		},
	}
	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)
	personallyProcuredMoves := []models.PersonallyProcuredMove{
		{
			OriginalMoveDate: &pickupDate,
			Status:           models.PPMStatusPAYMENTREQUESTED,
			NetWeight:        models.Int64Pointer(4000),
			Advance:          &advance,
		},
	}
	ssd := models.ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   order,
		CurrentDutyStation:      yuma,
		NewDutyStation:          fortGordon,
		PPMRemainingEntitlement: 3000,
		WeightAllotment:         wtgEntitlements,
		TotalWeightAllotment:    17500,
		Shipments:               shipments,
		PreparationDate:         time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
		PersonallyProcuredMoves: personallyProcuredMoves,
		Obligations: models.Obligations{
			MaxObligation:    models.Obligation{Gcc: unit.Cents(600000), SIT: unit.Cents(53000)},
			ActualObligation: models.Obligation{Gcc: unit.Cents(500000), SIT: unit.Cents(30000)},
		},
	}
	sswPage1 := models.FormatValuesShipmentSummaryWorksheetFormPage1(ssd)

	suite.Equal("01-Jan-2019", sswPage1.PreparationDate)

	suite.Equal("Jenkins Jr., Marcus Joseph", sswPage1.ServiceMemberName)
	suite.Equal("E-9", sswPage1.RankGrade)
	suite.Equal("Air Force", sswPage1.ServiceBranch)
	suite.Equal("90 days per each shipment", sswPage1.MaxSITStorageEntitlement)
	suite.Equal("Yuma AFB, IA 50309", sswPage1.AuthorizedOrigin)
	suite.Equal("Fort Gordon, GA 30813", sswPage1.AuthorizedDestination)
	suite.Equal("NO", sswPage1.POVAuthorized)
	suite.Equal("444-555-8888", sswPage1.PreferredPhoneNumber)
	suite.Equal("michael+ppm-expansion_1@truss.works", sswPage1.PreferredEmail)
	suite.Equal("1234567890", sswPage1.DODId)

	suite.Equal("Air Force", sswPage1.IssuingBranchOrAgency)
	suite.Equal("21-Dec-2018", sswPage1.OrdersIssueDate)
	suite.Equal("PCS/012345", sswPage1.OrdersTypeAndOrdersNumber)
	suite.Equal("NTA4", sswPage1.TAC)
	suite.Equal("SAC", sswPage1.SAC)

	suite.Equal("Fort Gordon, GA 30813", sswPage1.NewDutyAssignment)

	suite.Equal("15,000", sswPage1.WeightAllotment)
	suite.Equal("2,000", sswPage1.WeightAllotmentProgear)
	suite.Equal("500", sswPage1.WeightAllotmentProgearSpouse)
	suite.Equal("17,500", sswPage1.TotalWeightAllotment)

	suite.Equal("01 - HHG (GBL)\n\n02 - PPM", sswPage1.ShipmentNumberAndTypes)
	suite.Equal("11-Jan-2019\n\n11-Jan-2019", sswPage1.ShipmentPickUpDates)
	suite.Equal("5,000 lbs - FINAL\n\n4,000 lbs - FINAL", sswPage1.ShipmentWeights)
	suite.Equal("Delivered\n\nAt destination", sswPage1.ShipmentCurrentShipmentStatuses)

	suite.Equal("17,500", sswPage1.TotalWeightAllotmentRepeat)
	suite.Equal("$6,000.00", sswPage1.MaxObligationGCC100)
	suite.Equal("$5,700.00", sswPage1.MaxObligationGCC95)
	suite.Equal("$530.00", sswPage1.MaxObligationSIT)
	suite.Equal("$3,600.00", sswPage1.MaxObligationGCCMaxAdvance)

	suite.Equal("3,000", sswPage1.PPMRemainingEntitlement)
	suite.Equal("$5,000.00", sswPage1.ActualObligationGCC100)
	suite.Equal("$4,750.00", sswPage1.ActualObligationGCC95)
	suite.Equal("$300.00", sswPage1.ActualObligationSIT)
	suite.Equal("$10.00", sswPage1.ActualObligationAdvance)
}

func (suite *ModelSuite) TestFormatValuesShipmentSummaryWorksheetFormPage2() {
	movingExpenses := models.MovingExpenseDocuments{
		{
			MovingExpenseType:    "TOLLS",
			RequestedAmountCents: unit.Cents(10000),
			PaymentMethod:        "OTHER",
		},
		{
			MovingExpenseType:    "GAS",
			RequestedAmountCents: unit.Cents(10000),
			PaymentMethod:        "OTHER",
		},
		{
			MovingExpenseType:    "CONTRACTED_EXPENSE",
			RequestedAmountCents: unit.Cents(20000),
			PaymentMethod:        "GTCC",
		},
		{
			MovingExpenseType:    "CONTRACTED_EXPENSE",
			RequestedAmountCents: unit.Cents(10000),
			PaymentMethod:        "GTCC",
		},
	}

	ssd := models.ShipmentSummaryFormData{
		MovingExpenseDocuments: movingExpenses,
	}
	sswPage2, _ := models.FormatValuesShipmentSummaryWorksheetFormPage2(ssd)
	// fields w/ no expenses should format as $0.00
	suite.Equal("$0.00", sswPage2.RentalEquipmentGTCCPaid)
	suite.Equal("$0.00", sswPage2.PackingMaterialsGTCCPaid)

	suite.Equal("$300.00", sswPage2.ContractedExpenseGTCCPaid)
	suite.Equal("$300.00", sswPage2.TotalGTCCPaid)
	suite.Equal("$300.00", sswPage2.TotalGTCCPaidRepeated)

	suite.Equal("$100.00", sswPage2.TollsMemberPaid)
	suite.Equal("$100.00", sswPage2.GasMemberPaid)
	suite.Equal("$200.00", sswPage2.TotalMemberPaid)
	suite.Equal("$200.00", sswPage2.TotalMemberPaidRepeated)
}

func (suite *ModelSuite) TestGroupExpenses() {
	testCases := []struct {
		input    models.MovingExpenseDocuments
		expected map[string]float64
	}{
		{
			models.MovingExpenseDocuments{
				{
					MovingExpenseType:    "TOLLS",
					RequestedAmountCents: unit.Cents(10000),
					PaymentMethod:        "GTCC",
				},
				{
					MovingExpenseType:    "GAS",
					RequestedAmountCents: unit.Cents(10000),
					PaymentMethod:        "OTHER",
				},
				{
					MovingExpenseType:    "GAS",
					RequestedAmountCents: unit.Cents(20000),
					PaymentMethod:        "OTHER",
				},
			},
			map[string]float64{
				"TollsGTCCPaid":   100,
				"GasMemberPaid":   300,
				"TotalMemberPaid": 300,
				"TotalGTCCPaid":   100,
			},
		},
		{
			models.MovingExpenseDocuments{
				{
					MovingExpenseType:    "PACKING_MATERIALS",
					RequestedAmountCents: unit.Cents(10000),
					PaymentMethod:        "GTCC",
				},
				{
					MovingExpenseType:    "PACKING_MATERIALS",
					RequestedAmountCents: unit.Cents(10000),
					PaymentMethod:        "OTHER",
				},
				{
					MovingExpenseType:    "WEIGHING_FEES",
					RequestedAmountCents: unit.Cents(10000),
					PaymentMethod:        "GTCC",
				},
				{
					MovingExpenseType:    "WEIGHING_FEES",
					RequestedAmountCents: unit.Cents(10000),
					PaymentMethod:        "OTHER",
				},
				{
					MovingExpenseType:    "WEIGHING_FEES",
					RequestedAmountCents: unit.Cents(20000),
					PaymentMethod:        "GTCC",
				},
				{
					MovingExpenseType:    "GAS",
					RequestedAmountCents: unit.Cents(20000),
					PaymentMethod:        "OTHER",
				},
			},
			map[string]float64{
				"PackingMaterialsGTCCPaid":   100,
				"PackingMaterialsMemberPaid": 100,
				"WeighingFeesMemberPaid":     100,
				"WeighingFeesGTCCPaid":       300,
				"GasMemberPaid":              200,
				"TotalMemberPaid":            400,
				"TotalGTCCPaid":              400,
			},
		},
	}

	for _, testCase := range testCases {
		actual := models.SubTotalExpenses(testCase.input)
		suite.Equal(testCase.expected, actual)
	}

}

func (suite *ModelSuite) TestFormatWeightAllotment() {
	hasDependant := models.ShipmentSummaryFormData{
		Order: models.Order{HasDependents: true},
		WeightAllotment: models.WeightAllotment{
			TotalWeightSelf:               1000,
			TotalWeightSelfPlusDependents: 2000,
		},
	}
	noDependant := models.ShipmentSummaryFormData{
		WeightAllotment: models.WeightAllotment{
			TotalWeightSelf:               1000,
			TotalWeightSelfPlusDependents: 2000,
		},
	}

	suite.Equal("2,000", models.FormatWeightAllotment(hasDependant))
	suite.Equal("1,000", models.FormatWeightAllotment(noDependant))
}

func (suite *ModelSuite) TestFormatLocation() {
	fortGordon := models.DutyStation{Name: "Fort Gordon", Address: models.Address{State: "GA", PostalCode: "30813"}}
	yuma := models.DutyStation{Name: "Yuma AFB", Address: models.Address{State: "IA", PostalCode: "50309"}}

	suite.Equal("Fort Gordon, GA 30813", models.FormatLocation(fortGordon))
	suite.Equal("Yuma AFB, IA 50309", models.FormatLocation(yuma))
}

func (suite *ModelSuite) TestFormatServiceMemberFullName() {
	sm1 := models.ServiceMember{
		Suffix:     models.StringPointer("Jr."),
		FirstName:  models.StringPointer("Tom"),
		MiddleName: models.StringPointer("James"),
		LastName:   models.StringPointer("Smith"),
	}
	sm2 := models.ServiceMember{
		FirstName: models.StringPointer("Tom"),
		LastName:  models.StringPointer("Smith"),
	}

	suite.Equal("Smith Jr., Tom James", models.FormatServiceMemberFullName(sm1))
	suite.Equal("Smith, Tom", models.FormatServiceMemberFullName(sm2))
}

func (suite *ModelSuite) TestFormatCurrentShipmentStatus() {
	completed := models.Shipment{Status: models.ShipmentStatusDELIVERED}
	inTransit := models.Shipment{Status: models.ShipmentStatusINTRANSIT}

	suite.Equal("Delivered", models.FormatCurrentShipmentStatus(completed))
	suite.Equal("In Transit", models.FormatCurrentShipmentStatus(inTransit))
}

func (suite *ModelSuite) TestFormatCurrentPPMStatus() {
	paymentRequested := models.PersonallyProcuredMove{Status: models.PPMStatusPAYMENTREQUESTED}
	completed := models.PersonallyProcuredMove{Status: models.PPMStatusCOMPLETED}

	suite.Equal("At destination", models.FormatCurrentPPMStatus(paymentRequested))
	suite.Equal("Completed", models.FormatCurrentPPMStatus(completed))
}

func (suite *ModelSuite) TestFormatRank() {
	e9 := models.ServiceMemberRankE9
	multipleRanks := models.ServiceMemberRankO1ACADEMYGRADUATE

	suite.Equal("E-9", models.FormatRank(&e9))
	suite.Equal("O-1/Service Academy Graduate", models.FormatRank(&multipleRanks))
}

func (suite *ModelSuite) TestFormatShipmentNumberAndType() {
	singleShipment := models.Shipments{models.Shipment{}}
	multipleShipments := models.Shipments{models.Shipment{}, models.Shipment{}}
	singlePPM := models.PersonallyProcuredMoves{models.PersonallyProcuredMove{}}
	multiplePPMs := models.PersonallyProcuredMoves{models.PersonallyProcuredMove{}, models.PersonallyProcuredMove{}}
	var blankHHGSlice []models.Shipment
	var blankPPMSlice []models.PersonallyProcuredMove

	multipleShipmentsFormatted := models.FormatAllShipments(blankPPMSlice, multipleShipments)
	multiplePPMsFormatted := models.FormatAllShipments(multiplePPMs, blankHHGSlice)
	varietyOfShipmentsFormatted := models.FormatAllShipments(multiplePPMs, singleShipment)

	// testing single shipment moves
	suite.Equal("01 - HHG (GBL)", models.FormatAllShipments(blankPPMSlice, singleShipment).ShipmentNumberAndTypes)
	suite.Equal("01 - PPM", models.FormatAllShipments(singlePPM, blankHHGSlice).ShipmentNumberAndTypes)

	// testing multiple shipment moves
	suite.Equal("01 - HHG (GBL)\n\n02 - HHG (GBL)", multipleShipmentsFormatted.ShipmentNumberAndTypes)

	// testing multiple ppm moves
	suite.Equal("01 - PPM\n\n02 - PPM", multiplePPMsFormatted.ShipmentNumberAndTypes)

	// testing a variety of shipments and ppms
	suite.Equal("01 - HHG (GBL)\n\n02 - PPM\n\n03 - PPM", varietyOfShipmentsFormatted.ShipmentNumberAndTypes)
}

func (suite *ModelSuite) TestFormatShipmentWeight() {
	pounds := unit.Pound(1000)
	shipment := models.Shipment{NetWeight: &pounds}

	suite.Equal("1,000 lbs - FINAL", models.FormatShipmentWeight(shipment))
}

func (suite *ModelSuite) TestFormatPickupDate() {
	hhgPickupDate := time.Date(2018, time.December, 1, 0, 0, 0, 0, time.UTC)
	shipment := models.Shipment{ActualPickupDate: &hhgPickupDate}

	suite.Equal("01-Dec-2018", models.FormatShipmentPickupDate(shipment))
}

func (suite *ModelSuite) TestFormatWeights() {
	suite.Equal("0", models.FormatWeights(0))
	suite.Equal("10", models.FormatWeights(10))
	suite.Equal("1,000", models.FormatWeights(1000))
	suite.Equal("1,000,000", models.FormatWeights(1000000))
}

func (suite *ModelSuite) TestFormatOrdersIssueDate() {
	dec212018 := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	jan012019 := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	suite.Equal("21-Dec-2018", models.FormatDate(dec212018))
	suite.Equal("01-Jan-2019", models.FormatDate(jan012019))
}

func (suite *ModelSuite) TestFormatOrdersType() {
	pcsOrder := models.Order{OrdersType: internalmessages.OrdersTypePERMANENTCHANGEOFSTATION}
	var unknownOrdersType internalmessages.OrdersType = "UNKNOWN_ORDERS_TYPE"
	localMoveOrder := models.Order{OrdersType: unknownOrdersType}

	suite.Equal("PCS", models.FormatOrdersType(pcsOrder))
	suite.Equal("", models.FormatOrdersType(localMoveOrder))
}

func (suite *ModelSuite) TestFormatServiceMemberAffiliation() {
	airForce := models.AffiliationAIRFORCE
	marines := models.AffiliationMARINES

	suite.Equal("Air Force", models.FormatServiceMemberAffiliation(&airForce))
	suite.Equal("Marines", models.FormatServiceMemberAffiliation(&marines))
}

func (suite *ModelSuite) TestFormatPPMWeight() {
	var pounds int64 = 1000
	ppm := models.PersonallyProcuredMove{NetWeight: &pounds}
	noWtg := models.PersonallyProcuredMove{NetWeight: nil}

	suite.Equal("1,000 lbs - FINAL", models.FormatPPMWeight(ppm))
	suite.Equal("", models.FormatPPMWeight(noWtg))
}

func (suite *ModelSuite) TestCalculatePPMEntitlementNoHHGPPMLessThanMaxEntitlement() {
	ppmWeight := int64(900)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement := models.CalculateRemainingPPMEntitlement(move, totalEntitlement)

	suite.Equal(unit.Pound(ppmWeight), ppmRemainingEntitlement)
}

func (suite *ModelSuite) TestCalculatePPMEntitlementNoHHGPPMGreaterThanMaxEntitlement() {
	ppmWeight := int64(1100)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement := models.CalculateRemainingPPMEntitlement(move, totalEntitlement)

	suite.Equal(totalEntitlement, ppmRemainingEntitlement)
}

func (suite *ModelSuite) TestCalculatePPMEntitlementPPMGreaterThanRemainingEntitlement() {
	ppmWeight := int64(1100)
	totalEntitlement := unit.Pound(1000)
	hhg := unit.Pound(100)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
		Shipments:               models.Shipments{models.Shipment{NetWeight: &hhg}},
	}

	ppmRemainingEntitlement := models.CalculateRemainingPPMEntitlement(move, totalEntitlement)

	suite.Equal(totalEntitlement-hhg, ppmRemainingEntitlement)
}

func (suite *ModelSuite) TestCalculatePPMEntitlementPPMLessThanRemainingEntitlement() {
	ppmWeight := int64(500)
	totalEntitlement := unit.Pound(1000)
	hhg := unit.Pound(100)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
		Shipments:               models.Shipments{models.Shipment{NetWeight: &hhg}},
	}

	ppmRemainingEntitlement := models.CalculateRemainingPPMEntitlement(move, totalEntitlement)

	suite.Equal(unit.Pound(ppmWeight), ppmRemainingEntitlement)
}
