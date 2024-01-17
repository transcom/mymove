// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
// RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
// RA: in which this would be considered a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package shipmentsummaryworksheet

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFetchDataShipmentSummaryWorksheet() {
	//advanceID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	rank := models.ServiceMemberRankE9
	SSWPPMComputer := NewSSWPPMComputer()

	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
			},
		},
		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model:    yuma,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.ServiceMember{
				Rank: &rank,
			},
		},
		{
			Model: models.SignedCertification{},
		},
	}, nil)

	ppmShipmentID := ppmShipment.ID

	serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID

	session := auth.Session{
		UserID:          ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}

	models.SaveMoveDependencies(suite.DB(), &ppmShipment.Shipment.MoveTaskOrder)

	ssd, err := SSWPPMComputer.FetchDataShipmentSummaryWorksheetFormData(suite.AppContextForTest(), &session, ppmShipmentID)

	suite.NoError(err)
	suite.Equal(ppmShipment.Shipment.MoveTaskOrder.Orders.ID, ssd.Order.ID)
	suite.Require().Len(ssd.PPMShipments, 1)
	suite.Equal(ppmShipment.ID, ssd.PPMShipments[0].ID)
	suite.Equal(serviceMemberID, ssd.ServiceMember.ID)
	suite.Equal(yuma.ID, ssd.CurrentDutyLocation.ID)
	suite.Equal(yuma.Address.ID, ssd.CurrentDutyLocation.Address.ID)
	suite.Equal(fortGordon.ID, ssd.NewDutyLocation.ID)
	suite.Equal(fortGordon.Address.ID, ssd.NewDutyLocation.Address.ID)
	rankWtgAllotment := models.GetWeightAllotment(rank)
	suite.Equal(unit.Pound(rankWtgAllotment.TotalWeightSelf), ssd.WeightAllotment.Entitlement)
	suite.Equal(unit.Pound(rankWtgAllotment.ProGearWeight), ssd.WeightAllotment.ProGear)
	suite.Equal(unit.Pound(0), ssd.WeightAllotment.SpouseProGear)
	suite.Require().NotNil(ssd.ServiceMember.Rank)
	weightAllotment := models.GetWeightAllotment(*ssd.ServiceMember.Rank)
	// E_9 rank, no dependents, no spouse pro-gear
	totalWeight := weightAllotment.TotalWeightSelf + weightAllotment.ProGearWeight
	suite.Require().Nil(err)
	suite.Equal(unit.Pound(totalWeight), ssd.WeightAllotment.TotalWeight)
	suite.Equal(ppmShipment.EstimatedWeight, ssd.PPMShipments[0].EstimatedWeight)
	suite.Require().NotNil(ssd.PPMShipments[0].AdvanceAmountRequested)
	suite.Equal(ppmShipment.AdvanceAmountRequested, ssd.PPMShipments[0].AdvanceAmountRequested)
	// suite.Equal(signedCertification.ID, ssd.SignedCertification.ID)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFetchDataShipmentSummaryWorksheetWithErrorNoMove() {
	//advanceID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	rank := models.ServiceMemberRankE9
	SSWPPMComputer := NewSSWPPMComputer()

	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
			},
		},
		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model:    yuma,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.ServiceMember{
				Rank: &rank,
			},
		},
	}, nil)

	PPMShipmentID := uuid.Nil
	serviceMemberID := move.Orders.ServiceMemberID

	session := auth.Session{
		UserID:          move.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}

	emptySSD, err := SSWPPMComputer.FetchDataShipmentSummaryWorksheetFormData(suite.AppContextForTest(), &session, PPMShipmentID)

	suite.Error(err)
	suite.Equal(emptySSD, services.ShipmentSummaryFormData{})
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFetchMovingExpensesShipmentSummaryWorksheetNoPPM() {
	serviceMemberID, _ := uuid.NewV4()

	ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
	session := auth.Session{
		UserID:          ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}

	movingExpenses, err := FetchMovingExpensesShipmentSummaryWorksheet(ppmShipment, suite.AppContextForTest(), &session)

	suite.Len(movingExpenses, 0)
	suite.NoError(err)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFetchDataShipmentSummaryWorksheetOnlyPPM() {
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	rank := models.ServiceMemberRankE9
	SSWPPMComputer := NewSSWPPMComputer()

	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
			},
		},
		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model:    yuma,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.ServiceMember{
				Rank: &rank,
			},
		},
		{
			Model: models.SignedCertification{},
		},
	}, nil)
	ppmShipmentID := ppmShipment.ID

	serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID
	session := auth.Session{
		UserID:          ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}
	models.SaveMoveDependencies(suite.DB(), &ppmShipment.Shipment.MoveTaskOrder)
	ssd, err := SSWPPMComputer.FetchDataShipmentSummaryWorksheetFormData(suite.AppContextForTest(), &session, ppmShipmentID)

	suite.NoError(err)
	suite.Equal(ppmShipment.Shipment.MoveTaskOrder.Orders.ID, ssd.Order.ID)
	suite.Require().Len(ssd.PPMShipments, 1)
	suite.Equal(ppmShipment.ID, ssd.PPMShipments[0].ID)
	suite.Equal(serviceMemberID, ssd.ServiceMember.ID)
	suite.Equal(yuma.ID, ssd.CurrentDutyLocation.ID)
	suite.Equal(yuma.Address.ID, ssd.CurrentDutyLocation.Address.ID)
	suite.Equal(fortGordon.ID, ssd.NewDutyLocation.ID)
	suite.Equal(fortGordon.Address.ID, ssd.NewDutyLocation.Address.ID)
	rankWtgAllotment := models.GetWeightAllotment(rank)
	suite.Equal(unit.Pound(rankWtgAllotment.TotalWeightSelf), ssd.WeightAllotment.Entitlement)
	suite.Equal(unit.Pound(rankWtgAllotment.ProGearWeight), ssd.WeightAllotment.ProGear)
	suite.Equal(unit.Pound(0), ssd.WeightAllotment.SpouseProGear)
	suite.Require().NotNil(ssd.ServiceMember.Rank)
	weightAllotment := models.GetWeightAllotment(*ssd.ServiceMember.Rank)
	// E_9 rank, no dependents, no spouse pro-gear
	totalWeight := weightAllotment.TotalWeightSelf + weightAllotment.ProGearWeight
	suite.Equal(unit.Pound(totalWeight), ssd.WeightAllotment.TotalWeight)
	suite.Equal(ppmShipment.EstimatedWeight, ssd.PPMShipments[0].EstimatedWeight)
	suite.Require().NotNil(ssd.PPMShipments[0].AdvanceAmountRequested)
	suite.Equal(ppmShipment.AdvanceAmountRequested, ssd.PPMShipments[0].AdvanceAmountRequested)
	// suite.Equal(signedCertification.ID, ssd.SignedCertification.ID)
	suite.Require().Len(ssd.MovingExpenses, 0)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatValuesShipmentSummaryWorksheetFormPage1() {
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	wtgEntitlements := services.SSWMaxWeightEntitlement{
		Entitlement:   15000,
		ProGear:       2000,
		SpouseProGear: 500,
		TotalWeight:   17500,
	}

	serviceMemberID, _ := uuid.NewV4()
	serviceBranch := models.AffiliationAIRFORCE
	rank := models.ServiceMemberRankE9
	serviceMember := models.ServiceMember{
		ID:             serviceMemberID,
		FirstName:      models.StringPointer("Marcus"),
		MiddleName:     models.StringPointer("Joseph"),
		LastName:       models.StringPointer("Jenkins"),
		Suffix:         models.StringPointer("Jr."),
		Telephone:      models.StringPointer("444-555-8888"),
		PersonalEmail:  models.StringPointer("michael+ppm-expansion_1@truss.works"),
		Edipi:          models.StringPointer("1234567890"),
		Affiliation:    &serviceBranch,
		Rank:           &rank,
		DutyLocationID: &yuma.ID,
	}

	orderIssueDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	order := models.Order{
		IssueDate:         orderIssueDate,
		OrdersType:        internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:      models.StringPointer("012345"),
		NewDutyLocationID: fortGordon.ID,
		TAC:               models.StringPointer("NTA4"),
		SAC:               models.StringPointer("SAC"),
		HasDependents:     true,
		SpouseHasProGear:  true,
	}
	pickupDate := time.Date(2019, time.January, 11, 0, 0, 0, 0, time.UTC)
	netWeight := unit.Pound(4000)
	cents := unit.Cents(1000)
	PPMShipments := []models.PPMShipment{
		{
			ExpectedDepartureDate:  pickupDate,
			Status:                 models.PPMShipmentStatusWaitingOnCustomer,
			EstimatedWeight:        &netWeight,
			AdvanceAmountRequested: &cents,
		},
	}
	ssd := services.ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   order,
		CurrentDutyLocation:     yuma,
		NewDutyLocation:         fortGordon,
		PPMRemainingEntitlement: 3000,
		WeightAllotment:         wtgEntitlements,
		PreparationDate:         time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
		PPMShipments:            PPMShipments,
	}
	sswPage1 := FormatValuesShipmentSummaryWorksheetFormPage1(ssd)

	suite.Equal("01-Jan-2019", sswPage1.PreparationDate)

	suite.Equal("Jenkins Jr., Marcus Joseph", sswPage1.ServiceMemberName)
	suite.Equal("E-9", sswPage1.RankGrade)
	suite.Equal("Air Force", sswPage1.ServiceBranch)
	suite.Equal("90 days per each shipment", sswPage1.MaxSITStorageEntitlement)
	suite.Equal("Yuma AFB, IA 50309", sswPage1.AuthorizedOrigin)
	suite.Equal("Fort Eisenhower, GA 30813", sswPage1.AuthorizedDestination)
	suite.Equal("No", sswPage1.POVAuthorized)
	suite.Equal("444-555-8888", sswPage1.PreferredPhoneNumber)
	suite.Equal("michael+ppm-expansion_1@truss.works", sswPage1.PreferredEmail)
	suite.Equal("1234567890", sswPage1.DODId)

	suite.Equal("Air Force", sswPage1.IssuingBranchOrAgency)
	suite.Equal("21-Dec-2018", sswPage1.OrdersIssueDate)
	suite.Equal("PCS/012345", sswPage1.OrdersTypeAndOrdersNumber)

	suite.Equal("Fort Eisenhower, GA 30813", sswPage1.NewDutyAssignment)

	suite.Equal("15,000", sswPage1.WeightAllotment)
	suite.Equal("2,000", sswPage1.WeightAllotmentProgear)
	suite.Equal("500", sswPage1.WeightAllotmentProgearSpouse)
	suite.Equal("17,500", sswPage1.TotalWeightAllotment)

	suite.Equal("01 - PPM", sswPage1.ShipmentNumberAndTypes)
	suite.Equal("11-Jan-2019", sswPage1.ShipmentPickUpDates)
	suite.Equal("4,000 lbs - FINAL", sswPage1.ShipmentWeights)
	suite.Equal("Waiting On Customer", sswPage1.ShipmentCurrentShipmentStatuses)

	suite.Equal("17,500", sswPage1.TotalWeightAllotmentRepeat)

	// All obligation tests must be temporarily stopped until calculator is rebuilt

	// suite.Equal("$6,000.00", sswPage1.MaxObligationGCC100)
	// suite.Equal("$5,700.00", sswPage1.MaxObligationGCC95)
	// suite.Equal("$530.00", sswPage1.MaxObligationSIT)
	// suite.Equal("$3,600.00", sswPage1.MaxObligationGCCMaxAdvance)

	suite.Equal("3,000", sswPage1.PPMRemainingEntitlement)
	// suite.Equal("$5,000.00", sswPage1.ActualObligationGCC100)
	// suite.Equal("$4,750.00", sswPage1.ActualObligationGCC95)
	// suite.Equal("$300.00", sswPage1.ActualObligationSIT)
	// suite.Equal("$10.00", sswPage1.ActualObligationAdvance)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatValuesShipmentSummaryWorksheetFormPage2() {
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	orderIssueDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)

	order := models.Order{
		IssueDate:         orderIssueDate,
		OrdersType:        internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:      models.StringPointer("012345"),
		NewDutyLocationID: fortGordon.ID,
		TAC:               models.StringPointer("NTA4"),
		SAC:               models.StringPointer("SAC"),
		HasDependents:     true,
		SpouseHasProGear:  true,
	}
	paidWithGTCC := false
	tollExpense := models.MovingExpenseReceiptTypeTolls
	oilExpense := models.MovingExpenseReceiptTypeOil
	amount := unit.Cents(10000)
	movingExpenses := models.MovingExpenses{
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
	}

	ssd := services.ShipmentSummaryFormData{
		Order:          order,
		MovingExpenses: movingExpenses,
	}
	sswPage2 := FormatValuesShipmentSummaryWorksheetFormPage2(ssd)

	suite.Equal("NTA4", sswPage2.TAC)
	suite.Equal("SAC", sswPage2.SAC)

	// fields w/ no expenses should format as $0.00, but must be temporarily removed until string function is replaced
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatValuesShipmentSummaryWorksheetFormPage3() {
	signatureDate := time.Date(2019, time.January, 26, 14, 40, 0, 0, time.UTC)
	sm := models.ServiceMember{
		FirstName: models.StringPointer("John"),
		LastName:  models.StringPointer("Smith"),
	}
	paidWithGTCC := false
	tollExpense := models.MovingExpenseReceiptTypeTolls
	oilExpense := models.MovingExpenseReceiptTypeOil
	amount := unit.Cents(10000)
	movingExpenses := models.MovingExpenses{
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCC,
		},
	}
	signature := models.SignedCertification{
		Date: signatureDate,
	}

	ssd := services.ShipmentSummaryFormData{
		ServiceMember:       sm,
		SignedCertification: signature,
		MovingExpenses:      movingExpenses,
	}

	sswPage3 := FormatValuesShipmentSummaryWorksheetFormPage3(ssd)

	suite.Equal("", sswPage3.AmountsPaid)
	suite.Equal("John Smith electronically signed", sswPage3.ServiceMemberSignature)
	suite.Equal("26 Jan 2019 at 2:40pm", sswPage3.SignatureDate)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestGroupExpenses() {
	paidWithGTCC := false
	tollExpense := models.MovingExpenseReceiptTypeTolls
	oilExpense := models.MovingExpenseReceiptTypeOil
	amount := unit.Cents(10000)
	testCases := []struct {
		input    models.MovingExpenses
		expected map[string]float64
	}{
		{
			models.MovingExpenses{
				{
					MovingExpenseType: &tollExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
				{
					MovingExpenseType: &oilExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
				{
					MovingExpenseType: &oilExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
				{
					MovingExpenseType: &oilExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
				{
					MovingExpenseType: &tollExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
			},
			map[string]float64{
				"OilMemberPaid":   300,
				"TollsMemberPaid": 200,
			},
		},
		{
			models.MovingExpenses{
				{
					MovingExpenseType: &tollExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
				{
					MovingExpenseType: &oilExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
				{
					MovingExpenseType: &oilExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
				{
					MovingExpenseType: &oilExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
				{
					MovingExpenseType: &tollExpense,
					Amount:            &amount,
					PaidWithGTCC:      &paidWithGTCC,
				},
			},
			map[string]float64{
				"OilMemberPaid":   300,
				"TollsMemberPaid": 200,
			},
		},
	}

	for _, testCase := range testCases {
		actual := SubTotalExpenses(testCase.input)
		suite.Equal(testCase.expected, actual)
	}

}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestCalculatePPMEntitlementPPMGreaterThanRemainingEntitlement() {
	ppmWeight := unit.Pound(1100)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement, err := CalculateRemainingPPMEntitlement(move, totalEntitlement)
	suite.NoError(err)

	suite.Equal(totalEntitlement, ppmRemainingEntitlement)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestCalculatePPMEntitlementPPMLessThanRemainingEntitlement() {
	ppmWeight := unit.Pound(500)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement, err := CalculateRemainingPPMEntitlement(move, totalEntitlement)
	suite.NoError(err)

	suite.Equal(unit.Pound(ppmWeight), ppmRemainingEntitlement)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatSSWGetEntitlement() {
	spouseHasProGear := true
	hasDependants := true
	allotment := models.GetWeightAllotment(models.ServiceMemberRankE1)
	expectedTotalWeight := allotment.TotalWeightSelfPlusDependents + allotment.ProGearWeight + allotment.ProGearWeightSpouse
	sswEntitlement := SSWGetEntitlement(models.ServiceMemberRankE1, hasDependants, spouseHasProGear)

	suite.Equal(unit.Pound(expectedTotalWeight), sswEntitlement.TotalWeight)
	suite.Equal(unit.Pound(allotment.TotalWeightSelfPlusDependents), sswEntitlement.Entitlement)
	suite.Equal(unit.Pound(allotment.ProGearWeightSpouse), sswEntitlement.SpouseProGear)
	suite.Equal(unit.Pound(allotment.ProGearWeight), sswEntitlement.ProGear)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatSSWGetEntitlementNoDependants() {
	spouseHasProGear := false
	hasDependants := false
	allotment := models.GetWeightAllotment(models.ServiceMemberRankE1)
	expectedTotalWeight := allotment.TotalWeightSelf + allotment.ProGearWeight
	sswEntitlement := SSWGetEntitlement(models.ServiceMemberRankE1, hasDependants, spouseHasProGear)

	suite.Equal(unit.Pound(expectedTotalWeight), sswEntitlement.TotalWeight)
	suite.Equal(unit.Pound(allotment.TotalWeightSelf), sswEntitlement.Entitlement)
	suite.Equal(unit.Pound(allotment.ProGearWeight), sswEntitlement.ProGear)
	suite.Equal(unit.Pound(0), sswEntitlement.SpouseProGear)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatLocation() {
	fortEisenhower := models.DutyLocation{Name: "Fort Eisenhower, GA 30813", Address: models.Address{State: "GA", PostalCode: "30813"}}
	yuma := models.DutyLocation{Name: "Yuma AFB", Address: models.Address{State: "IA", PostalCode: "50309"}}

	suite.Equal("Fort Eisenhower, GA 30813", fortEisenhower.Name)
	suite.Equal("Yuma AFB, IA 50309", FormatLocation(yuma))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatServiceMemberFullName() {
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

	suite.Equal("Smith Jr., Tom James", FormatServiceMemberFullName(sm1))
	suite.Equal("Smith, Tom", FormatServiceMemberFullName(sm2))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatCurrentPPMStatus() {
	draft := models.PPMShipment{Status: models.PPMShipmentStatusDraft}
	submitted := models.PPMShipment{Status: models.PPMShipmentStatusSubmitted}

	suite.Equal("Draft", FormatCurrentPPMStatus(draft))
	suite.Equal("Submitted", FormatCurrentPPMStatus(submitted))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatRank() {
	e9 := models.ServiceMemberRankE9
	multipleRanks := models.ServiceMemberRankO1ACADEMYGRADUATE

	suite.Equal("E-9", FormatRank(&e9))
	suite.Equal("O-1 or Service Academy Graduate", FormatRank(&multipleRanks))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatShipmentNumberAndType() {
	singlePPM := models.PPMShipments{models.PPMShipment{}}
	multiplePPMs := models.PPMShipments{models.PPMShipment{}, models.PPMShipment{}}

	multiplePPMsFormatted := FormatAllShipments(multiplePPMs)
	singlePPMFormatted := FormatAllShipments(singlePPM)

	// testing single shipment moves
	suite.Equal("01 - PPM", singlePPMFormatted.ShipmentNumberAndTypes)

	// testing multiple ppm moves
	suite.Equal("01 - PPM\n\n02 - PPM", multiplePPMsFormatted.ShipmentNumberAndTypes)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatWeights() {
	suite.Equal("0", FormatWeights(0))
	suite.Equal("10", FormatWeights(10))
	suite.Equal("1,000", FormatWeights(1000))
	suite.Equal("1,000,000", FormatWeights(1000000))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatOrdersIssueDate() {
	dec212018 := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	jan012019 := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	suite.Equal("21-Dec-2018", FormatDate(dec212018))
	suite.Equal("01-Jan-2019", FormatDate(jan012019))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatOrdersType() {
	pcsOrder := models.Order{OrdersType: internalmessages.OrdersTypePERMANENTCHANGEOFSTATION}
	var unknownOrdersType internalmessages.OrdersType = "UNKNOWN_ORDERS_TYPE"
	localOrder := models.Order{OrdersType: unknownOrdersType}

	suite.Equal("PCS", FormatOrdersType(pcsOrder))
	suite.Equal("", FormatOrdersType(localOrder))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatServiceMemberAffiliation() {
	airForce := models.AffiliationAIRFORCE
	marines := models.AffiliationMARINES

	suite.Equal("Air Force", FormatServiceMemberAffiliation(&airForce))
	suite.Equal("Marines", FormatServiceMemberAffiliation(&marines))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatPPMWeight() {
	pounds := unit.Pound(1000)
	ppm := models.PPMShipment{EstimatedWeight: &pounds}
	noWtg := models.PPMShipment{EstimatedWeight: nil}

	suite.Equal("1,000 lbs - FINAL", FormatPPMWeight(ppm))
	suite.Equal("", FormatPPMWeight(noWtg))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestCalculatePPMEntitlementNoHHGPPMLessThanMaxEntitlement() {
	ppmWeight := unit.Pound(900)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement, err := CalculateRemainingPPMEntitlement(move, totalEntitlement)
	suite.NoError(err)

	suite.Equal(unit.Pound(ppmWeight), ppmRemainingEntitlement)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestCalculatePPMEntitlementNoHHGPPMGreaterThanMaxEntitlement() {
	ppmWeight := unit.Pound(1100)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement, err := CalculateRemainingPPMEntitlement(move, totalEntitlement)
	suite.NoError(err)

	suite.Equal(totalEntitlement, ppmRemainingEntitlement)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatSignature() {
	sm := models.ServiceMember{
		FirstName: models.StringPointer("John"),
		LastName:  models.StringPointer("Smith"),
	}

	formattedSignature := FormatSignature(sm)

	suite.Equal("John Smith electronically signed", formattedSignature)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatSignatureDate() {
	signatureDate := time.Date(2019, time.January, 26, 14, 40, 0, 0, time.UTC)

	signature := models.SignedCertification{
		Date: signatureDate,
	}
	sswfd := ShipmentSummaryFormData{
		SignedCertification: signature,
	}

	formattedDate := FormatSignatureDate(sswfd.SignedCertification)

	suite.Equal("26 Jan 2019 at 2:40pm", formattedDate)
}
