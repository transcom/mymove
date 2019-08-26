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
	suite.NoError(err)
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
		{
			MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
			RequestedAmountCents: unit.Cents(100000),
			PaymentMethod:        "GTCC",
		},
		{
			MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
			RequestedAmountCents: unit.Cents(20000),
			PaymentMethod:        "GTCC",
		},
		{
			MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
			RequestedAmountCents: unit.Cents(10000),
			PaymentMethod:        "OTHER",
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
	suite.Equal("$100.00", sswPage2.TotalMemberPaidSIT)
	suite.Equal("$1,200.00", sswPage2.TotalGTCCPaidSIT)
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

func (suite *ModelSuite) TestFormatSSWGetEntitlement() {
	spouseHasProGear := true
	hasDependants := true
	allotment := models.GetWeightAllotment(models.ServiceMemberRankE1)
	totalEntitlement, err := models.GetEntitlement(models.ServiceMemberRankE1, hasDependants, spouseHasProGear)
	suite.NoError(err)
	sswEntitlement := models.SSWGetEntitlement(models.ServiceMemberRankE1, hasDependants, spouseHasProGear)

	suite.Equal(unit.Pound(totalEntitlement), sswEntitlement.TotalWeight)
	suite.Equal(unit.Pound(allotment.TotalWeightSelfPlusDependents), sswEntitlement.Entitlement)
	suite.Equal(unit.Pound(allotment.ProGearWeightSpouse), sswEntitlement.SpouseProGear)
	suite.Equal(unit.Pound(allotment.ProGearWeight), sswEntitlement.ProGear)
}

func (suite *ModelSuite) TestFormatSSWGetEntitlementNoDependants() {
	spouseHasProGear := false
	hasDependants := false
	allotment := models.GetWeightAllotment(models.ServiceMemberRankE1)
	totalEntitlement, err := models.GetEntitlement(models.ServiceMemberRankE1, hasDependants, spouseHasProGear)
	suite.NoError(err)
	sswEntitlement := models.SSWGetEntitlement(models.ServiceMemberRankE1, hasDependants, spouseHasProGear)

	suite.Equal(unit.Pound(totalEntitlement), sswEntitlement.TotalWeight)
	suite.Equal(unit.Pound(allotment.TotalWeightSelf), sswEntitlement.Entitlement)
	suite.Equal(unit.Pound(allotment.ProGearWeight), sswEntitlement.ProGear)
	suite.Equal(unit.Pound(0), sswEntitlement.SpouseProGear)
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
	singlePPM := models.PersonallyProcuredMoves{models.PersonallyProcuredMove{}}
	multiplePPMs := models.PersonallyProcuredMoves{models.PersonallyProcuredMove{}, models.PersonallyProcuredMove{}}

	multiplePPMsFormatted := models.FormatAllShipments(multiplePPMs)
	singlePPMFormatted := models.FormatAllShipments(singlePPM)

	// testing single shipment moves
	suite.Equal("01 - PPM", singlePPMFormatted.ShipmentNumberAndTypes)

	// testing multiple ppm moves
	suite.Equal("01 - PPM\n\n02 - PPM", multiplePPMsFormatted.ShipmentNumberAndTypes)
}

func (suite *ModelSuite) TestFormatAllSITExpenses() {
	startdate1 := time.Date(2019, 5, 12, 0, 0, 0, 0, time.UTC)
	endDate1 := time.Date(2019, 5, 15, 0, 0, 0, 0, time.UTC)
	startdate2 := time.Date(2019, 5, 15, 0, 0, 0, 0, time.UTC)
	endDate2 := time.Date(2019, 5, 20, 0, 0, 0, 0, time.UTC)
	sitExpenses := models.MovingExpenseDocuments{
		{
			MovingExpenseType: models.MovingExpenseTypeSTORAGE,
			StorageStartDate:  &startdate1,
			StorageEndDate:    &endDate1,
		},
		{
			MovingExpenseType: models.MovingExpenseTypeSTORAGE,
			StorageStartDate:  &startdate2,
			StorageEndDate:    &endDate2,
		},
	}

	formattedSitExpenses := models.FormatAllSITExpenses(sitExpenses)

	suite.Equal("01 - PPM\n\n02 - PPM", formattedSitExpenses.NumberAndTypes)
	suite.Equal("12-May-2019\n\n15-May-2019", formattedSitExpenses.EntryDates)
	suite.Equal("15-May-2019\n\n20-May-2019", formattedSitExpenses.EndDates)
	suite.Equal("3\n\n5", formattedSitExpenses.DaysInStorage)
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
	pounds := unit.Pound(1000)
	ppm := models.PersonallyProcuredMove{NetWeight: &pounds}
	noWtg := models.PersonallyProcuredMove{NetWeight: nil}

	suite.Equal("1,000 lbs - FINAL", models.FormatPPMWeight(ppm))
	suite.Equal("", models.FormatPPMWeight(noWtg))
}

func (suite *ModelSuite) TestCalculatePPMEntitlementNoHHGPPMLessThanMaxEntitlement() {
	ppmWeight := unit.Pound(900)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement, err := models.CalculateRemainingPPMEntitlement(move, totalEntitlement)
	suite.NoError(err)

	suite.Equal(unit.Pound(ppmWeight), ppmRemainingEntitlement)
}

func (suite *ModelSuite) TestCalculatePPMEntitlementNoHHGPPMGreaterThanMaxEntitlement() {
	ppmWeight := unit.Pound(1100)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement, err := models.CalculateRemainingPPMEntitlement(move, totalEntitlement)
	suite.NoError(err)

	suite.Equal(totalEntitlement, ppmRemainingEntitlement)
}

func (suite *ModelSuite) TestFormatSignature() {
	signatureDate := time.Date(2019, time.January, 26, 14, 40, 0, 0, time.UTC)
	sm := models.ServiceMember{
		FirstName: models.StringPointer("John"),
		LastName:  models.StringPointer("Smith"),
	}
	signature := models.SignedCertification{
		Date: signatureDate,
	}
	sswfd := models.ShipmentSummaryFormData{
		ServiceMember:       sm,
		SignedCertification: signature,
	}

	formattedSignature := models.FormatSignature(sswfd)

	suite.Equal("John Smith electronically signed on 26 Jan 2019 at 2:40pm", formattedSignature)
}
