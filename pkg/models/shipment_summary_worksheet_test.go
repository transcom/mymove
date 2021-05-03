//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
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
	netWeight := unit.Pound(10000)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID:              move.ID,
			NetWeight:           &netWeight,
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
			Status:                   models.MoveDocumentStatusOK,
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: serviceMemberID,
			ServiceMember:   move.Orders.ServiceMember,
		},
	}
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments)

	session := auth.Session{
		UserID:          move.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}
	ppm.Move.Submit()
	ppm.Move.Approve()
	// This is the same PPM model as ppm, but this is the one that will be saved by SaveMoveDependencies
	ppm.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppm.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(suite.DB(), &ppm.Move)
	certificationType := models.SignedCertificationTypePPMPAYMENT
	signedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   moveID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &certificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	ssd, err := models.FetchDataShipmentSummaryWorksheetFormData(suite.DB(), &session, moveID)

	suite.NoError(err)
	suite.Equal(move.Orders.ID, ssd.Order.ID)
	suite.Require().Len(ssd.PersonallyProcuredMoves, 1)
	suite.Equal(ppm.ID, ssd.PersonallyProcuredMoves[0].ID)
	suite.Equal(serviceMemberID, ssd.ServiceMember.ID)
	suite.Equal(yuma.ID, ssd.CurrentDutyStation.ID)
	suite.Equal(yuma.Address.ID, ssd.CurrentDutyStation.Address.ID)
	suite.Equal(fortGordon.ID, ssd.NewDutyStation.ID)
	suite.Equal(fortGordon.Address.ID, ssd.NewDutyStation.Address.ID)
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
	suite.Require().Len(ssd.MovingExpenseDocuments, 2)
	suite.NotNil(ssd.MovingExpenseDocuments[0].ID)
	suite.NotNil(ssd.MovingExpenseDocuments[1].ID)
	suite.Equal(ppm.NetWeight, ssd.PersonallyProcuredMoves[0].NetWeight)
	suite.Require().NotNil(ssd.PersonallyProcuredMoves[0].Advance)
	suite.Equal(ppm.Advance.ID, ssd.PersonallyProcuredMoves[0].Advance.ID)
	suite.Equal(unit.Cents(1000), ssd.PersonallyProcuredMoves[0].Advance.RequestedAmount)
	suite.Equal(signedCertification.ID, ssd.SignedCertification.ID)
}

func (suite *ModelSuite) TestFetchDataShipmentSummaryWorksheetOnlyPPM() {
	moveID, _ := uuid.NewV4()
	serviceMemberID, _ := uuid.NewV4()
	//advanceID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	fortGordon := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	rank := models.ServiceMemberRankE9
	moveType := models.SelectedMoveTypePPM

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
	netWeight := unit.Pound(10000)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID:              move.ID,
			NetWeight:           &netWeight,
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
			Status:                   models.MoveDocumentStatusOK,
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: serviceMemberID,
			ServiceMember:   move.Orders.ServiceMember,
		},
	}
	movedocuments2 := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   models.MoveDocumentStatusEXCLUDEFROMCALCULATION,
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: serviceMemberID,
			ServiceMember:   move.Orders.ServiceMember,
		},
	}
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments2)
	session := auth.Session{
		UserID:          move.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}
	ppm.Move.Submit()
	ppm.Move.Approve()
	// This is the same PPM model as ppm, but this is the one that will be saved by SaveMoveDependencies
	ppm.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppm.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(suite.DB(), &ppm.Move)
	certificationType := models.SignedCertificationTypePPMPAYMENT
	signedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   moveID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &certificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	ssd, err := models.FetchDataShipmentSummaryWorksheetFormData(suite.DB(), &session, moveID)

	suite.NoError(err)
	suite.Equal(move.Orders.ID, ssd.Order.ID)
	suite.Require().Len(ssd.PersonallyProcuredMoves, 1)
	suite.Equal(ppm.ID, ssd.PersonallyProcuredMoves[0].ID)
	suite.Equal(serviceMemberID, ssd.ServiceMember.ID)
	suite.Equal(yuma.ID, ssd.CurrentDutyStation.ID)
	suite.Equal(yuma.Address.ID, ssd.CurrentDutyStation.Address.ID)
	suite.Equal(fortGordon.ID, ssd.NewDutyStation.ID)
	suite.Equal(fortGordon.Address.ID, ssd.NewDutyStation.Address.ID)
	rankWtgAllotment := models.GetWeightAllotment(rank)
	suite.Equal(unit.Pound(rankWtgAllotment.TotalWeightSelf), ssd.WeightAllotment.Entitlement)
	suite.Equal(unit.Pound(rankWtgAllotment.ProGearWeight), ssd.WeightAllotment.ProGear)
	suite.Equal(unit.Pound(0), ssd.WeightAllotment.SpouseProGear)
	suite.Require().NotNil(ssd.ServiceMember.Rank)
	weightAllotment := models.GetWeightAllotment(*ssd.ServiceMember.Rank)
	// E_9 rank, no dependents, no spouse pro-gear
	totalWeight := weightAllotment.TotalWeightSelf + weightAllotment.ProGearWeight
	suite.Equal(unit.Pound(totalWeight), ssd.WeightAllotment.TotalWeight)
	suite.Require().Len(ssd.MovingExpenseDocuments, 2)
	suite.NotNil(ssd.MovingExpenseDocuments[0].ID)
	suite.NotNil(ssd.MovingExpenseDocuments[1].ID)
	suite.Equal(ppm.NetWeight, ssd.PersonallyProcuredMoves[0].NetWeight)
	suite.Require().NotNil(ssd.PersonallyProcuredMoves[0].Advance)
	suite.Equal(ppm.Advance.ID, ssd.PersonallyProcuredMoves[0].Advance.ID)
	suite.Equal(unit.Cents(1000), ssd.PersonallyProcuredMoves[0].Advance.RequestedAmount)
	suite.Equal(signedCertification.ID, ssd.SignedCertification.ID)
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
	suite.NoError(err)
}

func (suite *ModelSuite) TestFormatValuesShipmentSummaryWorksheetFormPage1() {
	yuma := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	fortGordon := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	wtgEntitlements := models.SSWMaxWeightEntitlement{
		Entitlement:   15000,
		ProGear:       2000,
		SpouseProGear: 500,
		TotalWeight:   17500,
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
		IssueDate:        orderIssueDate,
		OrdersType:       internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:     models.StringPointer("012345"),
		NewDutyStationID: fortGordon.ID,
		TAC:              models.StringPointer("NTA4"),
		SAC:              models.StringPointer("SAC"),
		HasDependents:    true,
		SpouseHasProGear: true,
	}
	pickupDate := time.Date(2019, time.January, 11, 0, 0, 0, 0, time.UTC)
	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)
	netWeight := unit.Pound(4000)
	personallyProcuredMoves := []models.PersonallyProcuredMove{
		{
			OriginalMoveDate: &pickupDate,
			Status:           models.PPMStatusPAYMENTREQUESTED,
			NetWeight:        &netWeight,
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
		PreparationDate:         time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
		PersonallyProcuredMoves: personallyProcuredMoves,
		Obligations: models.Obligations{
			MaxObligation:              models.Obligation{Gcc: unit.Cents(600000), SIT: unit.Cents(53000)},
			ActualObligation:           models.Obligation{Gcc: unit.Cents(500000), SIT: unit.Cents(30000), Miles: unit.Miles(4050)},
			NonWinningMaxObligation:    models.Obligation{Gcc: unit.Cents(700000), SIT: unit.Cents(63000)},
			NonWinningActualObligation: models.Obligation{Gcc: unit.Cents(600000), SIT: unit.Cents(40000), Miles: unit.Miles(5050)},
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
	suite.Equal("No", sswPage1.POVAuthorized)
	suite.Equal("444-555-8888", sswPage1.PreferredPhoneNumber)
	suite.Equal("michael+ppm-expansion_1@truss.works", sswPage1.PreferredEmail)
	suite.Equal("1234567890", sswPage1.DODId)

	suite.Equal("Air Force", sswPage1.IssuingBranchOrAgency)
	suite.Equal("21-Dec-2018", sswPage1.OrdersIssueDate)
	suite.Equal("PCS/012345", sswPage1.OrdersTypeAndOrdersNumber)

	suite.Equal("Fort Gordon, GA 30813", sswPage1.NewDutyAssignment)

	suite.Equal("15,000", sswPage1.WeightAllotment)
	suite.Equal("2,000", sswPage1.WeightAllotmentProgear)
	suite.Equal("500", sswPage1.WeightAllotmentProgearSpouse)
	suite.Equal("17,500", sswPage1.TotalWeightAllotment)

	suite.Equal("01 - PPM", sswPage1.ShipmentNumberAndTypes)
	suite.Equal("11-Jan-2019", sswPage1.ShipmentPickUpDates)
	suite.Equal("4,000 lbs - FINAL", sswPage1.ShipmentWeights)
	suite.Equal("At destination", sswPage1.ShipmentCurrentShipmentStatuses)

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
	fortGordon := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	orderIssueDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)

	order := models.Order{
		IssueDate:        orderIssueDate,
		OrdersType:       internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:     models.StringPointer("012345"),
		NewDutyStationID: fortGordon.ID,
		TAC:              models.StringPointer("NTA4"),
		SAC:              models.StringPointer("SAC"),
		HasDependents:    true,
		SpouseHasProGear: true,
	}
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
		Order:                  order,
		MovingExpenseDocuments: movingExpenses,
	}
	sswPage2, _ := models.FormatValuesShipmentSummaryWorksheetFormPage2(ssd)

	suite.Equal("NTA4", sswPage2.TAC)
	suite.Equal("SAC", sswPage2.SAC)

	// fields w/ no expenses should format as $0.00
	suite.Equal("$0.00", sswPage2.RentalEquipmentGTCCPaid.String())
	suite.Equal("$0.00", sswPage2.PackingMaterialsGTCCPaid.String())

	suite.Equal("$300.00", sswPage2.ContractedExpenseGTCCPaid.String())
	suite.Equal("$300.00", sswPage2.TotalGTCCPaid.String())
	suite.Equal("$300.00", sswPage2.TotalGTCCPaidRepeated.String())

	suite.Equal("$100.00", sswPage2.TollsMemberPaid.String())
	suite.Equal("$100.00", sswPage2.GasMemberPaid.String())
	suite.Equal("$200.00", sswPage2.TotalMemberPaid.String())
	suite.Equal("$200.00", sswPage2.TotalMemberPaidRepeated.String())
	suite.Equal("$100.00", sswPage2.TotalMemberPaidSIT.String())
	suite.Equal("$1,200.00", sswPage2.TotalGTCCPaidSIT.String())
}

func (suite *ModelSuite) TestFormatValuesShipmentSummaryWorksheetFormPage3() {
	signatureDate := time.Date(2019, time.January, 26, 14, 40, 0, 0, time.UTC)
	sm := models.ServiceMember{
		FirstName: models.StringPointer("John"),
		LastName:  models.StringPointer("Smith"),
	}
	movingExpenses := models.MovingExpenseDocuments{
		{
			MovingExpenseType:    models.MovingExpenseTypeOTHER,
			RequestedAmountCents: unit.Cents(100000),
			PaymentMethod:        "GTCC",
		},
		{
			MovingExpenseType:    models.MovingExpenseTypeOTHER,
			RequestedAmountCents: unit.Cents(20000),
			PaymentMethod:        "GTCC",
		},
		{
			MovingExpenseType:    models.MovingExpenseTypeOTHER,
			RequestedAmountCents: unit.Cents(10000),
			PaymentMethod:        "OTHER",
		},
	}
	signature := models.SignedCertification{
		Date: signatureDate,
	}

	ssd := models.ShipmentSummaryFormData{
		ServiceMember:          sm,
		SignedCertification:    signature,
		MovingExpenseDocuments: movingExpenses,
	}

	sswPage3 := models.FormatValuesShipmentSummaryWorksheetFormPage3(ssd)

	suite.Equal("$1,000.00\n\n$200.00\n\n$100.00", sswPage3.AmountsPaid)
	suite.Equal("John Smith electronically signed", sswPage3.ServiceMemberSignature)
	suite.Equal("26 Jan 2019 at 2:40pm", sswPage3.SignatureDate)
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
				{
					MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
					RequestedAmountCents: unit.Cents(20000),
					PaymentMethod:        "GTCC",
				},
			},
			map[string]float64{
				"TollsGTCCPaid":    100,
				"GasMemberPaid":    300,
				"TotalMemberPaid":  300,
				"TotalGTCCPaid":    100,
				"TotalPaidNonSIT":  400,
				"StorageGTCCPaid":  200,
				"TotalGTCCPaidSIT": 200,
				"TotalPaidSIT":     200,
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
				"TotalPaidNonSIT":            800,
			},
		},
	}

	for _, testCase := range testCases {
		actual := models.SubTotalExpenses(testCase.input)
		suite.Equal(testCase.expected, actual)
	}

}

func (suite *ModelSuite) TestCalculatePPMEntitlementPPMGreaterThanRemainingEntitlement() {
	ppmWeight := unit.Pound(1100)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement, err := models.CalculateRemainingPPMEntitlement(move, totalEntitlement)
	suite.NoError(err)

	suite.Equal(totalEntitlement, ppmRemainingEntitlement)
}

func (suite *ModelSuite) TestCalculatePPMEntitlementPPMLessThanRemainingEntitlement() {
	ppmWeight := unit.Pound(500)
	totalEntitlement := unit.Pound(1000)
	move := models.Move{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{models.PersonallyProcuredMove{NetWeight: &ppmWeight}},
	}

	ppmRemainingEntitlement, err := models.CalculateRemainingPPMEntitlement(move, totalEntitlement)
	suite.NoError(err)

	suite.Equal(unit.Pound(ppmWeight), ppmRemainingEntitlement)
}

func (suite *ModelSuite) TestFormatSSWGetEntitlement() {
	spouseHasProGear := true
	hasDependants := true
	allotment := models.GetWeightAllotment(models.ServiceMemberRankE1)
	expectedTotalWeight := allotment.TotalWeightSelfPlusDependents + allotment.ProGearWeight + allotment.ProGearWeightSpouse
	sswEntitlement := models.SSWGetEntitlement(models.ServiceMemberRankE1, hasDependants, spouseHasProGear)

	suite.Equal(unit.Pound(expectedTotalWeight), sswEntitlement.TotalWeight)
	suite.Equal(unit.Pound(allotment.TotalWeightSelfPlusDependents), sswEntitlement.Entitlement)
	suite.Equal(unit.Pound(allotment.ProGearWeightSpouse), sswEntitlement.SpouseProGear)
	suite.Equal(unit.Pound(allotment.ProGearWeight), sswEntitlement.ProGear)
}

func (suite *ModelSuite) TestFormatSSWGetEntitlementNoDependants() {
	spouseHasProGear := false
	hasDependants := false
	allotment := models.GetWeightAllotment(models.ServiceMemberRankE1)
	expectedTotalWeight := allotment.TotalWeightSelf + allotment.ProGearWeight
	sswEntitlement := models.SSWGetEntitlement(models.ServiceMemberRankE1, hasDependants, spouseHasProGear)

	suite.Equal(unit.Pound(expectedTotalWeight), sswEntitlement.TotalWeight)
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
	localOrder := models.Order{OrdersType: unknownOrdersType}

	suite.Equal("PCS", models.FormatOrdersType(pcsOrder))
	suite.Equal("", models.FormatOrdersType(localOrder))
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

func (suite *ModelSuite) TestFormatOtherExpenses() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	sm := ppm.Move.Orders.ServiceMember

	assertions1 := testdatagen.Assertions{
		MovingExpenseDocument: models.MovingExpenseDocument{
			MovingExpenseType:    models.MovingExpenseTypeOTHER,
			RequestedAmountCents: unit.Cents(2589),
		},
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   "OK",
			MoveDocumentType:         "EXPENSE",
			Title:                    "The Bard",
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	}

	assertions2 := testdatagen.Assertions{
		MovingExpenseDocument: models.MovingExpenseDocument{
			MovingExpenseType:    models.MovingExpenseTypeOTHER,
			RequestedAmountCents: unit.Cents(1439),
		},
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   "OK",
			MoveDocumentType:         "EXPENSE",
			Title:                    "The Beedle",
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	}

	otherExpenseDocs := models.MovingExpenseDocuments{}
	otherExpenseDocs = append(otherExpenseDocs, testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions1))
	otherExpenseDocs = append(otherExpenseDocs, testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions2))

	formattedOtherExpenses := models.FormatOtherExpenses(otherExpenseDocs)

	suite.Equal("The Bard\n\nThe Beedle", formattedOtherExpenses.Descriptions)
	suite.Equal("$25.89\n\n$14.39", formattedOtherExpenses.AmountsPaid)
}

func (suite *ModelSuite) TestFormatSignature() {
	sm := models.ServiceMember{
		FirstName: models.StringPointer("John"),
		LastName:  models.StringPointer("Smith"),
	}

	formattedSignature := models.FormatSignature(sm)

	suite.Equal("John Smith electronically signed", formattedSignature)
}

func (suite *ModelSuite) TestFormatSignatureDate() {
	signatureDate := time.Date(2019, time.January, 26, 14, 40, 0, 0, time.UTC)

	signature := models.SignedCertification{
		Date: signatureDate,
	}
	sswfd := models.ShipmentSummaryFormData{
		SignedCertification: signature,
	}

	formattedDate := models.FormatSignatureDate(sswfd.SignedCertification)

	suite.Equal("26 Jan 2019 at 2:40pm", formattedDate)
}
