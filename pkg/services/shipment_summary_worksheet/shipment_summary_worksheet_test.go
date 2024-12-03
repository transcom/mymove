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
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services/mocks"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFetchDataShipmentSummaryWorksheet() {
	//advanceID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	grade := models.ServiceMemberGradeE9
	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	SSWPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)

	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
				Grade:      &grade,
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
			Model: models.SignedCertification{},
		},
	}, nil)

	var expenseAmount unit.Cents = 1000.00
	var currentExpenseType = models.MovingExpenseReceiptTypeOther
	paidGTCC := true
	movingExpense := models.MovingExpense{
		Amount:            &expenseAmount,
		MovingExpenseType: &currentExpenseType,
		PaidWithGTCC:      &paidGTCC,
	}

	factory.AddMovingExpenseToPPMShipment(suite.DB(), &ppmShipment, nil, &movingExpense)

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
	gradeWtgAllotment := models.GetWeightAllotment(grade)
	suite.Equal(unit.Pound(gradeWtgAllotment.TotalWeightSelf), ssd.WeightAllotment.Entitlement)
	suite.Equal(unit.Pound(gradeWtgAllotment.ProGearWeight), ssd.WeightAllotment.ProGear)
	suite.Equal(unit.Pound(500), ssd.WeightAllotment.SpouseProGear)
	suite.Require().NotNil(ssd.Order.Grade)
	weightAllotment := models.GetWeightAllotment(*ssd.Order.Grade)
	// E_9 rank, no dependents, with spouse pro-gear
	totalWeight := weightAllotment.TotalWeightSelf + weightAllotment.ProGearWeight + weightAllotment.ProGearWeightSpouse
	suite.Require().Nil(err)
	suite.Equal(unit.Pound(totalWeight), ssd.WeightAllotment.TotalWeight)
	suite.Equal(ppmShipment.EstimatedWeight, ssd.PPMShipments[0].EstimatedWeight)
	suite.Require().NotNil(ssd.PPMShipments[0].AdvanceAmountRequested)
	suite.Equal(ppmShipment.AdvanceAmountRequested, ssd.PPMShipments[0].AdvanceAmountRequested)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFetchDataShipmentSummaryWorksheetWithErrorNoMove() {
	//advanceID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	grade := models.ServiceMemberGradeE9
	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	SSWPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)

	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
				Grade:      &grade,
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
	suite.Nil(emptySSD)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatEMPLID() {
	edipi := "12345567890"
	affiliation := models.AffiliationCOASTGUARD
	emplid := "9999999"
	serviceMember := models.ServiceMember{
		ID:          uuid.Must(uuid.NewV4()),
		Edipi:       &edipi,
		Affiliation: &affiliation,
		Emplid:      &emplid,
	}

	result, err := formatEmplid(serviceMember)

	suite.Equal("EMPLID: 9999999", *result)
	suite.NoError(err)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatEMPLIDNotCG() {
	edipi := "12345567890"
	affiliation := models.AffiliationARMY
	emplid := "9999999"
	serviceMember := models.ServiceMember{
		ID:          uuid.Must(uuid.NewV4()),
		Edipi:       &edipi,
		Affiliation: &affiliation,
		Emplid:      &emplid,
	}

	result, err := formatEmplid(serviceMember)

	suite.Equal("12345567890", *result)
	suite.NoError(err)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatEMPLIDNull() {
	edipi := "12345567890"
	affiliation := models.AffiliationARMY
	serviceMember := models.ServiceMember{
		ID:          uuid.Must(uuid.NewV4()),
		Edipi:       &edipi,
		Affiliation: &affiliation,
	}

	result, err := formatEmplid(serviceMember)

	suite.Equal("12345567890", *result)
	suite.NoError(err)
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
	grade := models.ServiceMemberGradeE9
	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	SSWPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)

	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
				Grade:      &grade,
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
	gradeWtgAllotment := models.GetWeightAllotment(grade)
	suite.Equal(unit.Pound(gradeWtgAllotment.TotalWeightSelf), ssd.WeightAllotment.Entitlement)
	suite.Equal(unit.Pound(gradeWtgAllotment.ProGearWeight), ssd.WeightAllotment.ProGear)
	suite.Equal(unit.Pound(500), ssd.WeightAllotment.SpouseProGear)
	suite.Require().NotNil(ssd.Order.Grade)
	weightAllotment := models.GetWeightAllotment(*ssd.Order.Grade)
	// E_9 rank, no dependents, with spouse pro-gear
	totalWeight := weightAllotment.TotalWeightSelf + weightAllotment.ProGearWeight + weightAllotment.ProGearWeightSpouse
	suite.Equal(unit.Pound(totalWeight), ssd.WeightAllotment.TotalWeight)
	suite.Equal(ppmShipment.EstimatedWeight, ssd.PPMShipments[0].EstimatedWeight)
	suite.Require().NotNil(ssd.PPMShipments[0].AdvanceAmountRequested)
	suite.Equal(ppmShipment.AdvanceAmountRequested, ssd.PPMShipments[0].AdvanceAmountRequested)
	suite.Require().Len(ssd.MovingExpenses, 0)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatValuesShipmentSummaryWorksheetFormPage1() {
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	wtgEntitlements := models.SSWMaxWeightEntitlement{
		Entitlement:   15000,
		ProGear:       2000,
		SpouseProGear: 500,
		TotalWeight:   17500,
	}

	serviceMemberID, _ := uuid.NewV4()
	serviceBranch := models.AffiliationAIRFORCE
	grade := models.ServiceMemberGradeE9
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
		Grade:             &grade,
	}
	expectedPickupDate := time.Date(2019, time.January, 11, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(2019, time.February, 11, 0, 0, 0, 0, time.UTC)
	netWeight := unit.Pound(4000)
	cents := unit.Cents(1000)
	locator := "ABCDEF-01"
	estIncentive := unit.Cents(1000000)
	maxIncentive := unit.Cents(2000000)
	PPMShipments := models.PPMShipment{
		ExpectedDepartureDate:  expectedPickupDate,
		ActualMoveDate:         &actualPickupDate,
		Status:                 models.PPMShipmentStatusWaitingOnCustomer,
		EstimatedWeight:        &netWeight,
		AdvanceAmountRequested: &cents,
		EstimatedIncentive:     &estIncentive,
		MaxIncentive:           &maxIncentive,
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
		IsActualExpenseReimbursement: models.BoolPointer(true),
	}
	ssd := models.ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   order,
		CurrentDutyLocation:     yuma,
		NewDutyLocation:         fortGordon,
		PPMRemainingEntitlement: 3000,
		WeightAllotment:         wtgEntitlements,
		PreparationDate:         time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
		PPMShipment:             PPMShipments,
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	sswPage1, err := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage1(ssd, false)
	suite.NoError(err)
	suite.Equal(FormatDate(time.Now()), sswPage1.PreparationDate1)

	suite.Equal("Jenkins Jr., Marcus Joseph", sswPage1.ServiceMemberName)
	suite.Equal("E-9", sswPage1.RankGrade)
	suite.Equal("Air Force", sswPage1.ServiceBranch)
	suite.Equal("00 Days in SIT", sswPage1.MaxSITStorageEntitlement)
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
	suite.Equal("2,000", sswPage1.WeightAllotmentProGear)
	suite.Equal("500", sswPage1.WeightAllotmentProgearSpouse)
	suite.Equal("17,500", sswPage1.TotalWeightAllotment)

	suite.Equal(locator+" PPM", sswPage1.ShipmentNumberAndTypes)
	suite.Equal("11-Jan-2019", sswPage1.ShipmentPickUpDates)
	suite.Equal("4,000 lbs - Estimated", sswPage1.ShipmentWeights)
	suite.Equal("Waiting On Customer", sswPage1.ShipmentCurrentShipmentStatuses)
	suite.Equal("17,500", sswPage1.TotalWeightAllotmentRepeat)
	suite.Equal("15,000 lbs; $10,000.00", sswPage1.MaxObligationGCC100)
	suite.True(sswPage1.IsActualExpenseReimbursement)
	suite.Equal("Actual Expense Reimbursement", sswPage1.GCCIsActualExpenseReimbursement)

	// quick test when there is no PPM actual move date
	PPMShipmentWithoutActualMoveDate := models.PPMShipment{
		Status:                 models.PPMShipmentStatusWaitingOnCustomer,
		EstimatedWeight:        &netWeight,
		AdvanceAmountRequested: &cents,
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}

	ssdWithoutPPMActualMoveDate := models.ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   order,
		CurrentDutyLocation:     yuma,
		NewDutyLocation:         fortGordon,
		PPMRemainingEntitlement: 3000,
		WeightAllotment:         wtgEntitlements,
		PreparationDate:         time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
		PPMShipment:             PPMShipmentWithoutActualMoveDate,
	}
	sswPage1NoActualMoveDate, err := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage1(ssdWithoutPPMActualMoveDate, false)
	suite.NoError(err)
	suite.Equal("N/A", sswPage1NoActualMoveDate.ShipmentPickUpDates)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatValuesShipmentSummaryWorksheetFormPage2() {
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	orderIssueDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	locator := "ABCDEF-01"
	shipment := models.PPMShipment{
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
		IsActualExpenseReimbursement: models.BoolPointer(true),
	}

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
	paidWithGTCCFalse := false
	paidWithGTCCTrue := true
	tollExpense := models.MovingExpenseReceiptTypeTolls
	oilExpense := models.MovingExpenseReceiptTypeOil
	amount := unit.Cents(10000)
	movingExpenses := models.MovingExpenses{
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCFalse,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCFalse,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCTrue,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCFalse,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCTrue,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCTrue,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCFalse,
		},
	}

	ssd := models.ShipmentSummaryFormData{
		Order:          order,
		MovingExpenses: movingExpenses,
		PPMShipment:    shipment,
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	expensesMap := SubTotalExpenses(ssd.MovingExpenses)
	sswPage2, err := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage2(ssd, false, expensesMap)
	suite.NoError(err)
	suite.Equal("$200.00", sswPage2.TollsGTCCPaid)
	suite.Equal("$200.00", sswPage2.TollsMemberPaid)
	suite.Equal("$200.00", sswPage2.OilMemberPaid)
	suite.Equal("$100.00", sswPage2.OilGTCCPaid)
	suite.Equal("$300.00", sswPage2.TotalGTCCPaid)
	suite.Equal("$400.00", sswPage2.TotalMemberPaid)
	suite.Equal("NTA4", sswPage2.TAC)
	suite.Equal("SAC", sswPage2.SAC)
	suite.Equal("Actual Expense Reimbursement", sswPage2.IncentiveIsActualExpenseReimbursement)
	suite.Equal(`This PPM is being processed at actual expense reimbursement for valid expenses not to exceed the
		government constructed cost (GCC).`, sswPage2.HeaderIsActualExpenseReimbursement)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatValuesShipmentSummaryWorksheetFormPage2ExcludeRejectedOrExcludedExpensesFromTotal() {
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	orderIssueDate := time.Date(2018, time.December, 23, 0, 0, 0, 0, time.UTC)
	locator := "ABCDEF-01"
	singlePPM := models.PPMShipment{
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}
	order := models.Order{
		IssueDate:         orderIssueDate,
		OrdersType:        internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:      models.StringPointer("012346"),
		NewDutyLocationID: fortGordon.ID,
		TAC:               models.StringPointer("NTA4"),
		SAC:               models.StringPointer("SAC"),
		HasDependents:     true,
		SpouseHasProGear:  true,
	}
	paidWithGTCCFalse := false
	paidWithGTCCTrue := true
	tollExpense := models.MovingExpenseReceiptTypeTolls
	oilExpense := models.MovingExpenseReceiptTypeOil
	approvedStatus := models.PPMDocumentStatusApproved
	excludedStatus := models.PPMDocumentStatusExcluded
	rejectedStatus := models.PPMDocumentStatusRejected
	amount := unit.Cents(10000)
	smallerAmount := unit.Cents(5000)
	movingExpenses := models.MovingExpenses{
		// APPROVED
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCFalse,
			Status:            &approvedStatus,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &smallerAmount,
			PaidWithGTCC:      &paidWithGTCCTrue,
			Status:            &approvedStatus,
		},
		// EXCLUDED
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCTrue,
			Status:            &excludedStatus,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCFalse,
			Status:            &excludedStatus,
		},
		{
			MovingExpenseType: &oilExpense,
			Amount:            &smallerAmount,
			PaidWithGTCC:      &paidWithGTCCFalse,
			Status:            &excludedStatus,
		},
		// REJECTED
		{
			MovingExpenseType: &oilExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCFalse,
			Status:            &rejectedStatus,
		},
		{
			MovingExpenseType: &tollExpense,
			Amount:            &amount,
			PaidWithGTCC:      &paidWithGTCCTrue,
			Status:            &rejectedStatus,
		},
	}

	ssd := models.ShipmentSummaryFormData{
		Order:          order,
		MovingExpenses: movingExpenses,
		PPMShipment:    singlePPM,
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	expensesMap := SubTotalExpenses(ssd.MovingExpenses)
	sswPage2, err := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage2(ssd, false, expensesMap)
	suite.NoError(err)
	suite.Equal("$0.00", sswPage2.TollsGTCCPaid)
	suite.Equal("$100.00", sswPage2.TollsMemberPaid)
	suite.Equal("$0.00", sswPage2.OilMemberPaid)
	suite.Equal("$50.00", sswPage2.OilGTCCPaid)
	suite.Equal("$50.00", sswPage2.TotalGTCCPaid)
	suite.Equal("$100.00", sswPage2.TotalMemberPaid)
	suite.Equal("NTA4", sswPage2.TAC)
	suite.Equal("SAC", sswPage2.SAC)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatValuesShipmentSummaryWorksheetFormPage3() {
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	wtgEntitlements := models.SSWMaxWeightEntitlement{}
	serviceMember := models.ServiceMember{}
	order := models.Order{}
	expectedPickupDate := time.Date(2019, time.January, 11, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(2019, time.February, 11, 0, 0, 0, 0, time.UTC)
	netWeight := unit.Pound(4000)
	cents := unit.Cents(1000)
	locator := "ABCDEF-01"
	move := factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)
	PPMShipment := models.PPMShipment{
		ID:                     move.MTOShipments[0].PPMShipment.ID,
		ExpectedDepartureDate:  expectedPickupDate,
		ActualMoveDate:         &actualPickupDate,
		Status:                 models.PPMShipmentStatusWaitingOnCustomer,
		EstimatedWeight:        &netWeight,
		AdvanceAmountRequested: &cents,
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}
	ssd := models.ShipmentSummaryFormData{
		AllShipments:            move.MTOShipments,
		ServiceMember:           serviceMember,
		Order:                   order,
		CurrentDutyLocation:     yuma,
		NewDutyLocation:         fortGordon,
		PPMRemainingEntitlement: 3000,
		WeightAllotment:         wtgEntitlements,
		PreparationDate:         time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
		PPMShipment:             PPMShipment,
	}
	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	sswPage3, err := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage3(ssd, false)
	suite.NoError(err)
	suite.Equal(FormatDate(time.Now()), sswPage3.PreparationDate3)
	suite.Equal(make(map[string]string), sswPage3.AddShipments)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatAdditionalHHG() {
	page3Map := make(map[string]string)
	i := 1
	hhg := factory.BuildMTOShipment(suite.DB(), nil, nil)
	locator := "ABCDEF"
	hhg.ShipmentLocator = &locator

	page3Map, err := formatAdditionalHHG(page3Map, i, hhg)
	suite.NoError(err)
	suite.Equal(*hhg.ShipmentLocator+" HHG", page3Map["AddShipmentNumberAndTypes1"])
	suite.Equal("16-Mar-2020 Actual", page3Map["AddShipmentPickUpDates1"])
	suite.Equal("980 Actual", page3Map["AddShipmentWeights1"])
	suite.Equal(FormatEnum(string(hhg.Status), ""), page3Map["AddShipmentStatus1"])
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestMemberPaidRemainingPPMEntitlementFormatValuesShipmentSummaryWorksheetFormPage2() {
	storageExpense := models.MovingExpenseReceiptTypeStorage
	amount := unit.Cents(10000)
	movingExpenses := models.MovingExpenses{
		{
			MovingExpenseType:      &storageExpense,
			Amount:                 &amount,
			PaidWithGTCC:           models.BoolPointer(false),
			SITReimburseableAmount: models.CentPointer(unit.Cents(100)),
		},
	}

	locator := "ABCDEF-01"
	id := uuid.Must(uuid.NewV4())
	PPMShipments := []models.PPMShipment{
		{
			FinalIncentive:        models.CentPointer(unit.Cents(500)),
			AdvanceAmountReceived: models.CentPointer(unit.Cents(200)),
			ID:                    id,
			Shipment: models.MTOShipment{
				ShipmentLocator: &locator,
			},
		},
	}

	signedCertType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
	cert := models.SignedCertification{
		CertificationType: &signedCertType,
		CertificationText: "APPROVED",
		Signature:         "Firstname Lastname",
		UpdatedAt:         time.Now(),
		PpmID:             models.UUIDPointer(PPMShipments[0].ID),
	}
	var certs []*models.SignedCertification
	certs = append(certs, &cert)

	ssd := models.ShipmentSummaryFormData{
		MovingExpenses:       movingExpenses,
		PPMShipment:          PPMShipments[0],
		SignedCertifications: certs,
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	expensesMap := SubTotalExpenses(ssd.MovingExpenses)
	sswPage2, _ := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage2(ssd, true, expensesMap)
	suite.Equal("$4.00", sswPage2.PPMRemainingEntitlement)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestAOAPacketPPMEntitlementFormatValuesShipmentSummaryWorksheetFormPage2() {
	storageExpense := models.MovingExpenseReceiptTypeStorage
	amount := unit.Cents(10000)
	movingExpenses := models.MovingExpenses{
		{
			MovingExpenseType:      &storageExpense,
			Amount:                 &amount,
			PaidWithGTCC:           models.BoolPointer(false),
			SITReimburseableAmount: models.CentPointer(unit.Cents(100)),
		},
	}

	locator := "ABCDEF-01"

	PPMShipments := []models.PPMShipment{
		{
			FinalIncentive:        models.CentPointer(unit.Cents(500)),
			AdvanceAmountReceived: models.CentPointer(unit.Cents(200)),
			Shipment: models.MTOShipment{
				ShipmentLocator: &locator,
			},
		},
	}

	ssd := models.ShipmentSummaryFormData{
		MovingExpenses: movingExpenses,
		PPMShipment:    PPMShipments[0],
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	expensesMap := SubTotalExpenses(ssd.MovingExpenses)
	sswPage2, _ := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage2(ssd, false, expensesMap)
	suite.Equal("N/A", sswPage2.PPMRemainingEntitlement)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestNullCheckForFinalIncentiveAndAOAPPMEntitlementFormatValuesShipmentSummaryWorksheetFormPage2() {
	storageExpense := models.MovingExpenseReceiptTypeStorage
	amount := unit.Cents(10000)
	movingExpenses := models.MovingExpenses{
		{
			MovingExpenseType:      &storageExpense,
			Amount:                 &amount,
			PaidWithGTCC:           models.BoolPointer(false),
			SITReimburseableAmount: models.CentPointer(unit.Cents(100)),
		},
	}

	locator := "ABCDEF-01"
	id := uuid.Must(uuid.NewV4())
	PPMShipments := []models.PPMShipment{
		{

			ID: id,
			Shipment: models.MTOShipment{
				ShipmentLocator: &locator,
			},
		},
	}

	signedCertType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
	cert := models.SignedCertification{
		CertificationType: &signedCertType,
		CertificationText: "APPROVED",
		Signature:         "Firstname Lastname",
		UpdatedAt:         time.Now(),
		PpmID:             models.UUIDPointer(PPMShipments[0].ID),
	}
	var certs []*models.SignedCertification
	certs = append(certs, &cert)

	ssd := models.ShipmentSummaryFormData{
		MovingExpenses:       movingExpenses,
		PPMShipment:          PPMShipments[0],
		SignedCertifications: certs,
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	expensesMap := SubTotalExpenses(ssd.MovingExpenses)
	sswPage2, _ := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage2(ssd, true, expensesMap)
	suite.Equal("$1.00", sswPage2.PPMRemainingEntitlement)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestGTCCPaidRemainingPPMEntitlementFormatValuesShipmentSummaryWorksheetFormPage2() {
	storageExpense := models.MovingExpenseReceiptTypeStorage
	amount := unit.Cents(10000)
	movingExpenses := models.MovingExpenses{
		{
			MovingExpenseType:      &storageExpense,
			Amount:                 &amount,
			PaidWithGTCC:           models.BoolPointer(true),
			SITReimburseableAmount: models.CentPointer(unit.Cents(200)),
		},
	}

	locator := "ABCDEF-01"
	id := uuid.Must(uuid.NewV4())
	PPMShipments := []models.PPMShipment{
		{
			FinalIncentive:        models.CentPointer(unit.Cents(600)),
			AdvanceAmountReceived: models.CentPointer(unit.Cents(100)),
			ID:                    id,
			Shipment: models.MTOShipment{
				ShipmentLocator: &locator,
			},
		},
	}

	signedCertType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
	cert := models.SignedCertification{
		CertificationType: &signedCertType,
		CertificationText: "APPROVED",
		Signature:         "Firstname Lastname",
		UpdatedAt:         time.Now(),
		PpmID:             models.UUIDPointer(PPMShipments[0].ID),
	}
	var certs []*models.SignedCertification
	certs = append(certs, &cert)

	ssd := models.ShipmentSummaryFormData{
		MovingExpenses:       movingExpenses,
		PPMShipment:          PPMShipments[0],
		SignedCertifications: certs,
	}

	expensesMap := SubTotalExpenses(ssd.MovingExpenses)

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	sswPage2, _ := sswPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage2(ssd, true, expensesMap)
	suite.Equal("$105.00", sswPage2.PPMRemainingEntitlement)
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
				"TotalMemberPaid": 500,
			},
		},
	}

	for _, testCase := range testCases {
		actual := SubTotalExpenses(testCase.input)
		suite.Equal(testCase.expected, actual)
	}

}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatSSWGetEntitlement() {
	spouseHasProGear := true
	hasDependants := true
	allotment := models.GetWeightAllotment(models.ServiceMemberGradeE1)
	expectedTotalWeight := allotment.TotalWeightSelfPlusDependents + allotment.ProGearWeight + allotment.ProGearWeightSpouse
	sswEntitlement := SSWGetEntitlement(models.ServiceMemberGradeE1, hasDependants, spouseHasProGear)

	suite.Equal(unit.Pound(expectedTotalWeight), sswEntitlement.TotalWeight)
	suite.Equal(unit.Pound(allotment.TotalWeightSelfPlusDependents), sswEntitlement.Entitlement)
	suite.Equal(unit.Pound(allotment.ProGearWeightSpouse), sswEntitlement.SpouseProGear)
	suite.Equal(unit.Pound(allotment.ProGearWeight), sswEntitlement.ProGear)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatSSWGetEntitlementNoDependants() {
	spouseHasProGear := false
	hasDependants := false
	allotment := models.GetWeightAllotment(models.ServiceMemberGradeE1)
	expectedTotalWeight := allotment.TotalWeightSelf + allotment.ProGearWeight + allotment.ProGearWeightSpouse
	sswEntitlement := SSWGetEntitlement(models.ServiceMemberGradeE1, hasDependants, spouseHasProGear)

	suite.Equal(unit.Pound(expectedTotalWeight), sswEntitlement.TotalWeight)
	suite.Equal(unit.Pound(allotment.TotalWeightSelf), sswEntitlement.Entitlement)
	suite.Equal(unit.Pound(allotment.ProGearWeight), sswEntitlement.ProGear)
	suite.Equal(unit.Pound(500), sswEntitlement.SpouseProGear)
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
	e9 := models.ServiceMemberGradeE9
	multipleGrades := models.ServiceMemberGradeO1ACADEMYGRADUATE

	suite.Equal("E-9", FormatGrade(&e9))
	suite.Equal("O-1 or Service Academy Graduate", FormatGrade(&multipleGrades))
}

// This is the test
func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatShipmentNumberAndType() {
	locator := "ABCDEF-01"
	singlePPM := models.PPMShipment{
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}

	wtgEntitlements := models.SSWMaxWeightEntitlement{
		Entitlement:   15000,
		ProGear:       2000,
		SpouseProGear: 500,
		TotalWeight:   17500,
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	singlePPMFormatted := sswPPMComputer.FormatShipment(singlePPM, wtgEntitlements, false)

	// testing single shipment moves
	suite.Equal("ABCDEF-01 PPM", singlePPMFormatted.ShipmentNumberAndTypes)
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

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatPPMWeightEstimated() {
	pounds := unit.Pound(1000)
	ppm := models.PPMShipment{EstimatedWeight: &pounds}
	noWtg := models.PPMShipment{EstimatedWeight: nil}

	suite.Equal("1,000 lbs - Estimated", FormatPPMWeightEstimated(ppm))
	suite.Equal("", FormatPPMWeightEstimated(noWtg))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatPPMWeightFinal() {
	pounds := unit.Pound(1000)

	suite.Equal("1,000 lbs - Actual", FormatPPMWeightFinal(pounds))
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatAOASignedCertifications() {
	var err error
	move := factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)
	testDate := time.Now() // due to using updatedAt, time.Now() needs to be used to test cert times and dates
	aoaCertifications := Certifications{
		CustomerField: "",
		OfficeField:   "AOA: Firstname Lastname",
		DateField:     "AOA: " + FormatDate(testDate),
	}
	sswCertifications := Certifications{
		CustomerField: "",
		OfficeField:   "AOA: Firstname Lastname\nSSW: Firstname Lastname",
		DateField:     "AOA: " + FormatDate(testDate) + "\nSSW: " + FormatDate(testDate),
	}
	prepAOADate := FormatDate(testDate)
	prepSSWDate := FormatDate(testDate)

	signedCertType := models.SignedCertificationTypePreCloseoutReviewedPPMPAYMENT
	aoaSignedCertification := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &signedCertType,
				CertificationText: "APPROVED",
				Signature:         "Firstname Lastname",
				UpdatedAt:         testDate,
				PpmID:             models.UUIDPointer(move.MTOShipments[0].PPMShipment.ID),
			},
		},
	}, nil)
	var certs []*models.SignedCertification
	certs = append(certs, &aoaSignedCertification)

	formattedSignature := formatSignedCertifications(certs, move.MTOShipments[0].PPMShipment.ID, false)
	formattedDate := formatAOADate(certs, move.MTOShipments[0].PPMShipment.ID)
	suite.Equal(prepAOADate, formattedDate)
	suite.Equal(aoaCertifications, formattedSignature)

	signedCertType = models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
	ppmPaymentsignedCertification := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &signedCertType,
				CertificationText: "APPROVED",
				Signature:         "Firstname Lastname",
				UpdatedAt:         testDate,
				PpmID:             models.UUIDPointer(move.MTOShipments[0].PPMShipment.ID),
			},
		},
	}, nil)
	certs = append(certs, &ppmPaymentsignedCertification)

	formattedSignature = formatSignedCertifications(certs, move.MTOShipments[0].PPMShipment.ID, true)
	formattedDate, err = formatSSWDate(certs, move.MTOShipments[0].PPMShipment.ID)
	suite.NoError(err)
	suite.Equal(prepSSWDate, formattedDate)
	suite.Equal(sswCertifications, formattedSignature)

}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatSSWSignedCertifications() {
	var err error
	move := factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)
	testDate := time.Now() // due to using updatedAt, time.Now() needs to be used to test cert times and dates
	sswCertifications := Certifications{
		CustomerField: "",
		OfficeField:   "AOA: \nSSW: Firstname Lastname",
		DateField:     "AOA: " + "\nSSW: " + FormatDate(testDate),
	}
	prepSSWDate := FormatDate(testDate)

	var certs []*models.SignedCertification

	signedCertType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
	ppmPaymentsignedCertification := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &signedCertType,
				CertificationText: "APPROVED",
				Signature:         "Firstname Lastname",
				UpdatedAt:         testDate,
				PpmID:             models.UUIDPointer(move.MTOShipments[0].PPMShipment.ID),
			},
		},
	}, nil)
	certs = append(certs, &ppmPaymentsignedCertification)

	formattedSignature := formatSignedCertifications(certs, move.MTOShipments[0].PPMShipment.ID, true)
	formattedDate, err := formatSSWDate(certs, move.MTOShipments[0].PPMShipment.ID)
	suite.NoError(err)
	suite.Equal(prepSSWDate, formattedDate)
	suite.Equal(sswCertifications, formattedSignature)

}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatAddress() {
	// Test case 1: Valid W2 address
	validAddress := &models.Address{
		StreetAddress1: "123 Main St",
		City:           "Cityville",
		State:          "ST",
		PostalCode:     "12345",
	}

	expectedValidResult := "123 Main St,  Cityville ST 12345"

	resultValid := FormatAddress(validAddress)

	suite.Equal(expectedValidResult, resultValid)

	// Test case 2: Valid W2 Address with country
	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	validAddress2 := &models.Address{
		StreetAddress1: "123 Main St",
		City:           "Cityville",
		State:          "ST",
		PostalCode:     "12345",
		Country:        &country,
	}

	expectedValidResult2 := "123 Main St,  Cityville ST US12345"

	resultValid2 := FormatAddress(validAddress2)

	suite.Equal(expectedValidResult2, resultValid2)

	// Test case 2: Nil W2 address
	nilAddress := (*models.Address)(nil)

	expectedNilResult := "W2 Address not found"

	resultNil := FormatAddress(nilAddress)

	suite.Equal(expectedNilResult, resultNil)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestNilOrValue() {
	// Test case 1: Non-nil pointer
	validPointer := "ValidValue"
	validResult := nilOrValue(&validPointer)
	expectedValidResult := "ValidValue"

	if validResult != expectedValidResult {
		suite.Equal(expectedValidResult, validResult)
	}

	// Test case 2: Nil pointer
	nilPointer := (*string)(nil)
	nilResult := nilOrValue(nilPointer)
	expectedNilResult := ""

	if nilResult != expectedNilResult {
		suite.Equal(expectedNilResult, nilResult)
	}
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestMergeTextFields() {
	// Test case 1: Non-empty input slices
	fields1 := []textField{
		{Pages: []int{1, 2}, ID: "1", Name: "Field1", Value: "Value1", Multiline: false, Locked: true},
		{Pages: []int{3, 4}, ID: "2", Name: "Field2", Value: "Value2", Multiline: true, Locked: false},
	}

	fields2 := []textField{
		{Pages: []int{5, 6}, ID: "3", Name: "Field3", Value: "Value3", Multiline: true, Locked: false},
		{Pages: []int{7, 8}, ID: "4", Name: "Field4", Value: "Value4", Multiline: false, Locked: true},
	}

	fields3 := []textField{
		{Pages: []int{9, 10}, ID: "5", Name: "Field5", Value: "Value5", Multiline: true, Locked: false},
		{Pages: []int{11, 12}, ID: "6", Name: "Field6", Value: "Value6", Multiline: false, Locked: true},
	}

	mergedResult := mergeTextFields(fields1, fields2, fields3)

	expectedMergedResult := []textField{
		{Pages: []int{1, 2}, ID: "1", Name: "Field1", Value: "Value1", Multiline: false, Locked: true},
		{Pages: []int{3, 4}, ID: "2", Name: "Field2", Value: "Value2", Multiline: true, Locked: false},
		{Pages: []int{5, 6}, ID: "3", Name: "Field3", Value: "Value3", Multiline: true, Locked: false},
		{Pages: []int{7, 8}, ID: "4", Name: "Field4", Value: "Value4", Multiline: false, Locked: true},
		{Pages: []int{9, 10}, ID: "5", Name: "Field5", Value: "Value5", Multiline: true, Locked: false},
		{Pages: []int{11, 12}, ID: "6", Name: "Field6", Value: "Value6", Multiline: false, Locked: true},
	}

	suite.Equal(mergedResult, expectedMergedResult)

	// Test case 2: Empty input slices
	emptyResult := mergeTextFields([]textField{}, []textField{}, []textField{})
	expectedEmptyResult := []textField{}

	suite.Equal(emptyResult, expectedEmptyResult)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestCreateTextFields() {
	// Test case 1: Non-empty input
	type TestData struct {
		Field1 string
		Field2 int
		Field3 bool
		Field4 map[string]string
	}

	field4 := make(map[string]string)
	field4["Field4"] = "Value 4"
	testData := TestData{"Value1", 42, true, field4}
	pages := []int{1, 2}

	result := createTextFields(testData, pages...)

	expectedResult := []textField{
		{Pages: pages, ID: "1", Name: "Field1", Value: "Value1", Multiline: true, Locked: false},
		{Pages: pages, ID: "2", Name: "Field2", Value: "42", Multiline: true, Locked: false},
		{Pages: pages, ID: "3", Name: "Field3", Value: "true", Multiline: true, Locked: false},
		{Pages: pages, ID: "4", Name: "Field4", Value: "Value 4", Multiline: true, Locked: false},
	}

	suite.Equal(result, expectedResult)

	// Test case 2: Empty input
	emptyResult := createTextFields(struct{}{})

	suite.Nil(emptyResult)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFillSSWPDFForm() {
	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.FatalNoError(uploaderErr)
	generator, err := paperworkgenerator.NewGenerator(userUploader.Uploader())
	suite.FatalNil(err)

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	ppmGenerator, err := NewSSWPPMGenerator(generator)
	suite.FatalNoError(err)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	grade := models.ServiceMemberGradeE9
	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
				Grade:      &grade,
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
			Model: models.SignedCertification{
				UpdatedAt: time.Now(),
			},
		},
	}, nil)

	storageExpenseType := models.MovingExpenseReceiptTypeStorage
	movingExpense := models.MovingExpense{
		MovingExpenseType: &storageExpenseType,
		Amount:            models.CentPointer(unit.Cents(67899)),
		SITStartDate:      models.TimePointer(time.Now()),
		SITEndDate:        models.TimePointer(time.Now()),
	}

	factory.AddMovingExpenseToPPMShipment(suite.DB(), &ppmShipment, nil, &movingExpense)

	ppmShipmentID := ppmShipment.ID

	serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID

	session := auth.Session{
		UserID:          ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}

	models.SaveMoveDependencies(suite.DB(), &ppmShipment.Shipment.MoveTaskOrder)

	ssd, err := sswPPMComputer.FetchDataShipmentSummaryWorksheetFormData(suite.AppContextForTest(), &session, ppmShipmentID)
	suite.NoError(err)
	page1Data, page2Data, Page3Data, err := sswPPMComputer.FormatValuesShipmentSummaryWorksheet(*ssd, false)
	suite.NoError(err)
	test, info, err := ppmGenerator.FillSSWPDFForm(page1Data, page2Data, Page3Data)
	suite.NoError(err)
	println(test.Name())           // ensures was generated with temp filesystem
	suite.Equal(info.PageCount, 3) // ensures PDF is not corrupted
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestActualExpenseReimbursementCalculations() {

	// Helper function to format disbursement field for equal checks
	expectedDisbursementString := func(expectedGTCC int, expectedMember int) string {
		return "GTCC: " + FormatDollars((models.CentPointer(unit.Cents(expectedGTCC)).ToMillicents().ToDollarFloat())) + "\nMember: " + FormatDollars(models.CentPointer(unit.Cents(expectedMember)).ToMillicents().ToDollarFloat())
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.FatalNoError(uploaderErr)
	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)
	generator, err := paperworkgenerator.NewGenerator(userUploader.Uploader())
	suite.NoError(err)
	ppmGenerator, err := NewSSWPPMGenerator(generator)
	suite.NoError(err)
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	orderIssueDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	locator := "ABCDEF-01"

	shipment := models.PPMShipment{
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
		IsActualExpenseReimbursement: models.BoolPointer(true),
		FinalIncentive:               models.CentPointer(20000),
	}

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
	storageExpense := models.MovingExpenseReceiptTypeStorage
	contractedExpense := models.MovingExpenseReceiptTypeContractedExpense
	movingExpenses := models.MovingExpenses{
		{
			MovingExpenseType: &contractedExpense,
			PaidWithGTCC:      models.BoolPointer(false),
		},
		{
			MovingExpenseType: &contractedExpense,
			PaidWithGTCC:      models.BoolPointer(true),
		},
		{
			MovingExpenseType: &storageExpense,
			PaidWithGTCC:      models.BoolPointer(false),
		},
		{
			MovingExpenseType: &storageExpense,
			PaidWithGTCC:      models.BoolPointer(true),
		},
	}

	/**
		Expenses map:
			- movingExpenses[0] == Total Member Expenses
			- movingExpenses[1] == Total GTCC Expenses
			- movingExpenses[2] == Member SIT Expenses
			- movingExpenses[3] == GTCC SIT Expenses
	**/
	const (
		MemberTotalExpenses = 0
		GTCCTotalExpenses   = 1
		MemberSITExpenses   = 2
		GTCCSITExpenses     = 3
	)

	signedCertType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
	cert := models.SignedCertification{
		CertificationType: &signedCertType,
		CertificationText: "APPROVED",
		Signature:         "Firstname Lastname",
		UpdatedAt:         time.Now(),
		PpmID:             models.UUIDPointer(shipment.ID),
	}
	var certs []*models.SignedCertification
	certs = append(certs, &cert)

	ssd := models.ShipmentSummaryFormData{
		Order:                        order,
		MovingExpenses:               movingExpenses,
		PPMShipment:                  shipment,
		SignedCertifications:         certs,
		IsActualExpenseReimbursement: true,
	}

	// Final Incentive == 100% GCC

	/**
		Test case 1: GTCC is greater or equal to GCC

		Expected outcome: 	GTCC disbursement == 100% GCC (20000 cents)
							Member disbursement == 0
	**/
	movingExpenses[MemberTotalExpenses].Amount = models.CentPointer(50000)
	movingExpenses[GTCCTotalExpenses].Amount = models.CentPointer(50000)
	movingExpenses[MemberSITExpenses].Amount = models.CentPointer(50000)
	movingExpenses[MemberSITExpenses].SITReimburseableAmount = models.CentPointer(50000)
	movingExpenses[GTCCSITExpenses].Amount = models.CentPointer(50000)

	page1Data, page2Data, Page3Data, err := sswPPMComputer.FormatValuesShipmentSummaryWorksheet(ssd, true)
	suite.NoError(err)
	suite.Equal(expectedDisbursementString(20000, 0), page2Data.Disbursement)
	suite.Equal("$0.00", page2Data.PPMRemainingEntitlement) // Check that pre-tax remaining incentive has been set to 0

	// Usual test checks to ensure PDF was generated properly
	test, info, err := ppmGenerator.FillSSWPDFForm(page1Data, page2Data, Page3Data)
	suite.NoError(err)
	println(test.Name())           // ensures was generated with temp filesystem
	suite.Equal(info.PageCount, 3) // ensures PDF is not corrupted

	// Also test for AOA instead of payment packet
	page1Data, page2Data, Page3Data, err = sswPPMComputer.FormatValuesShipmentSummaryWorksheet(ssd, false)
	suite.NoError(err)
	suite.Equal(expectedDisbursementString(20000, 0), page2Data.Disbursement)
	suite.Equal("$0.00", page2Data.PPMRemainingEntitlement)

	// Check PDF generation again
	test, info, err = ppmGenerator.FillSSWPDFForm(page1Data, page2Data, Page3Data)
	suite.NoError(err)
	println(test.Name())
	suite.Equal(info.PageCount, 3)

	/**
		Test case 2: GTCC is less than GCC, and total member expenses (incl. SIT) exceed amount left over from GCC - GTCC

		Expected outcome: 	GTCC disbursement == GTCC paid expenses + GTCC paid SIT
							Member disbursement == 100% GCC - GTCC
	**/
	movingExpenses[MemberTotalExpenses].Amount = models.CentPointer(50000)
	movingExpenses[GTCCTotalExpenses].Amount = models.CentPointer(10000)
	movingExpenses[MemberSITExpenses].Amount = models.CentPointer(5000)
	movingExpenses[MemberSITExpenses].SITReimburseableAmount = models.CentPointer(5000)
	movingExpenses[GTCCSITExpenses].Amount = models.CentPointer(1500)

	page1Data, page2Data, Page3Data, err = sswPPMComputer.FormatValuesShipmentSummaryWorksheet(ssd, true)
	suite.NoError(err)
	suite.Equal(expectedDisbursementString(11500, 8500), page2Data.Disbursement)
	suite.Equal("$0.00", page2Data.PPMRemainingEntitlement)

	test, info, err = ppmGenerator.FillSSWPDFForm(page1Data, page2Data, Page3Data)
	suite.NoError(err)
	println(test.Name())
	suite.Equal(info.PageCount, 3)

	page1Data, page2Data, Page3Data, err = sswPPMComputer.FormatValuesShipmentSummaryWorksheet(ssd, false)
	suite.NoError(err)
	suite.Equal(expectedDisbursementString(11500, 8500), page2Data.Disbursement)
	suite.Equal("$0.00", page2Data.PPMRemainingEntitlement)

	test, info, err = ppmGenerator.FillSSWPDFForm(page1Data, page2Data, Page3Data)
	suite.NoError(err)
	println(test.Name())
	suite.Equal(info.PageCount, 3)

	/**
		Test case 3: GTCC is less than GCC, and total member expenses (incl. SIT) are lower than amount left over from GCC - GTCC

		Expected outcome: 	GTCC disbursement == GTCC paid expenses + GTCC paid SIT
							Member disbursement == Total Member Expenses + Member paid SIT
	**/
	movingExpenses[MemberTotalExpenses].Amount = models.CentPointer(1000)
	movingExpenses[GTCCTotalExpenses].Amount = models.CentPointer(10000)
	movingExpenses[MemberSITExpenses].Amount = models.CentPointer(2000)
	movingExpenses[MemberSITExpenses].SITReimburseableAmount = models.CentPointer(2000)
	movingExpenses[GTCCSITExpenses].Amount = models.CentPointer(1500)

	page1Data, page2Data, Page3Data, err = sswPPMComputer.FormatValuesShipmentSummaryWorksheet(ssd, true)
	suite.NoError(err)
	suite.Equal(expectedDisbursementString(11500, 3000), page2Data.Disbursement)
	suite.Equal("$0.00", page2Data.PPMRemainingEntitlement)

	test, info, err = ppmGenerator.FillSSWPDFForm(page1Data, page2Data, Page3Data)
	suite.NoError(err)
	println(test.Name())
	suite.Equal(info.PageCount, 3)

	page1Data, page2Data, Page3Data, err = sswPPMComputer.FormatValuesShipmentSummaryWorksheet(ssd, false)
	suite.NoError(err)
	suite.Equal(expectedDisbursementString(11500, 3000), page2Data.Disbursement)
	suite.Equal("$0.00", page2Data.PPMRemainingEntitlement)

	test, info, err = ppmGenerator.FillSSWPDFForm(page1Data, page2Data, Page3Data)
	suite.NoError(err)
	println(test.Name())
	suite.Equal(info.PageCount, 3)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatMaxAdvance() {
	cents := unit.Cents(1000)
	tests := []struct {
		name               string
		estimatedIncentive *unit.Cents
		expectedResult     string
	}{
		{
			name:               "Valid estimated incentive",
			estimatedIncentive: &cents,
			expectedResult:     "$6.00",
		},
		{
			name:               "Nil estimated incentive",
			estimatedIncentive: nil,
			expectedResult:     "No Incentive Found",
		},
	}

	for _, tt := range tests {
		result := formatMaxAdvance(tt.estimatedIncentive)
		suite.Equal(tt.expectedResult, result)
	}

}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatShipment() {
	exampleValue1 := unit.Cents(5000)
	exampleValue2 := unit.Cents(3000)
	exampleValue3 := unit.Cents(1000)
	maxIncentive := unit.Cents(1000)
	exampleValue4 := models.PPMAdvanceStatusReceived
	exampleValue5 := true
	locator := "ABCDEF-01"

	wtgEntitlements := models.SSWMaxWeightEntitlement{
		Entitlement:   15000,
		ProGear:       2000,
		SpouseProGear: 500,
		TotalWeight:   17500,
	}

	tests := []struct {
		name           string
		shipment       models.PPMShipment
		expectedResult models.WorkSheetShipment
		entitlements   models.SSWMaxWeightEntitlement
	}{
		{
			name: "All fields present",
			shipment: models.PPMShipment{
				FinalIncentive:        &exampleValue1, // Example value
				EstimatedIncentive:    &exampleValue2, // Example value
				MaxIncentive:          &maxIncentive,
				AdvanceAmountReceived: &exampleValue3, // Example value
				AdvanceStatus:         &exampleValue4,
				HasRequestedAdvance:   &exampleValue5,
				Shipment: models.MTOShipment{
					ShipmentLocator: &locator,
				},
			},
			expectedResult: models.WorkSheetShipment{
				FinalIncentive:         "$50.00",                     // Example expected result
				MaxIncentive:           "$500.00",                    // Example expected result
				MaxAdvance:             "$18.00",                     // Assuming formatMaxAdvance correctly formats
				EstimatedIncentive:     "$30.00",                     // Example expected result
				AdvanceAmountReceived:  "$10.00 Requested, Received", // Example expected result
				ShipmentNumberAndTypes: locator,
			},
			entitlements: wtgEntitlements,
		},
		{
			name: "Final Incentive nil",
			shipment: models.PPMShipment{
				FinalIncentive:        nil,
				EstimatedIncentive:    &exampleValue2, // Example value
				MaxIncentive:          &maxIncentive,
				AdvanceAmountReceived: &exampleValue3, // Example value
				AdvanceStatus:         &exampleValue4,
				HasRequestedAdvance:   &exampleValue5,
				Shipment: models.MTOShipment{
					ShipmentLocator: &locator,
				},
			},
			expectedResult: models.WorkSheetShipment{
				FinalIncentive:         "No final incentive.",
				MaxIncentive:           "$500.00",
				MaxAdvance:             "$18.00",                     // Assuming formatMaxAdvance correctly formats
				EstimatedIncentive:     "$30.00",                     // Example expected result
				AdvanceAmountReceived:  "$10.00 Requested, Received", // Example expected result
				ShipmentNumberAndTypes: locator,
			},
			entitlements: wtgEntitlements,
		},
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	sswPPMComputer := NewSSWPPMComputer(mockPPMCloseoutFetcher)

	for _, tt := range tests {
		result := sswPPMComputer.FormatShipment(tt.shipment, tt.entitlements, false)

		suite.Equal(tt.expectedResult.FinalIncentive, result.FinalIncentive)
		suite.Equal(tt.expectedResult.MaxAdvance, result.MaxAdvance)
		suite.Equal(tt.expectedResult.EstimatedIncentive, result.EstimatedIncentive)
		suite.Equal(tt.expectedResult.AdvanceAmountReceived, result.AdvanceAmountReceived)
		suite.Equal(tt.expectedResult.ShipmentNumberAndTypes+" PPM", result.ShipmentNumberAndTypes)
	}
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatAdditionalShipments() {
	locator := "ABCDEF-01"
	now := time.Now()

	ppm := models.PPMShipment{
		ID: uuid.Must(uuid.NewV4()),
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}

	ppm2 := models.PPMShipment{
		ID:     uuid.Must(uuid.NewV4()),
		Status: models.PPMShipmentStatusSubmitted,
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
		ActualMoveDate: &now,
		WeightTickets: models.WeightTickets{
			models.WeightTicket{
				AdjustedNetWeight: models.PoundPointer(1200),
			},
			models.WeightTicket{
				AdjustedNetWeight: models.PoundPointer(1200),
			},
			models.WeightTicket{
				EmptyWeight: models.PoundPointer(3000),
				FullWeight:  models.PoundPointer(4200),
			},
		},
	}

	ppm3 := models.PPMShipment{
		ID:     uuid.Must(uuid.NewV4()),
		Status: models.PPMShipmentStatusSubmitted,
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
		ActualMoveDate:  &now,
		EstimatedWeight: models.PoundPointer(25),
	}

	ppm4 := models.PPMShipment{
		ID:     uuid.Must(uuid.NewV4()),
		Status: models.PPMShipmentStatusSubmitted,
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
		ExpectedDepartureDate: now,
		EstimatedWeight:       models.PoundPointer(25),
	}

	primeActualWeight := unit.Pound(1234)
	primeEstimatedWeight := unit.Pound(1234)

	shipments := []models.MTOShipment{
		{
			ShipmentType:    models.MTOShipmentTypePPM,
			ShipmentLocator: &locator,
			PPMShipment:     &ppm2,
			Status:          models.MTOShipmentStatusSubmitted,
		},
		{
			ShipmentType:    models.MTOShipmentTypePPM,
			ShipmentLocator: &locator,
			PPMShipment:     &ppm3,
			Status:          models.MTOShipmentStatusSubmitted,
		},
		{
			ShipmentType:    models.MTOShipmentTypePPM,
			ShipmentLocator: &locator,
			PPMShipment:     &ppm4,
			Status:          models.MTOShipmentStatusSubmitted,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
			ShipmentLocator:      &locator,
			RequestedPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeActualWeight:    &primeActualWeight,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			ShipmentLocator:      &locator,
			RequestedPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeActualWeight:    &primeActualWeight,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypePPM,
			ShipmentLocator:      &locator,
			RequestedPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeActualWeight:    &primeActualWeight,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypePPM,
			ShipmentLocator:      &locator,
			RequestedPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeActualWeight:    &primeActualWeight,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypeMobileHome,
			ShipmentLocator:      &locator,
			RequestedPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeActualWeight:    &primeActualWeight,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypeBoatHaulAway,
			ShipmentLocator:      &locator,
			RequestedPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeActualWeight:    &primeActualWeight,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypeBoatTowAway,
			ShipmentLocator:      &locator,
			RequestedPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeActualWeight:    &primeActualWeight,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypePPM,
			ShipmentLocator:      &locator,
			RequestedPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypePPM,
			ShipmentLocator:      &locator,
			ActualPickupDate:     &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:          &ppm2,
			ShipmentType:         models.MTOShipmentTypePPM,
			ShipmentLocator:      &locator,
			ScheduledPickupDate:  &now,
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &primeEstimatedWeight,
		},
		{
			PPMShipment:     &ppm2,
			ShipmentType:    models.MTOShipmentTypePPM,
			ShipmentLocator: &locator,
			Status:          models.MTOShipmentStatusSubmitted,
		},
	}

	ssd := models.ShipmentSummaryFormData{
		PPMShipment:  ppm,
		AllShipments: shipments,
	}

	results, _ := formatAdditionalShipments(ssd)
	suite.Equal(len(results), 56) // # of shipments multiply by 4

	expectedMapKeys := [4]string{"AddShipmentNumberAndTypes", "AddShipmentPickUpDates", "AddShipmentWeights", "AddShipmentStatus"}

	for indexShipment, shipment := range shipments {
		for index, key := range expectedMapKeys {
			value, contains := results[fmt.Sprintf("%s%d", key, indexShipment+1)]
			suite.True(contains)
			// verify AddShipmentNumberAndTypes
			if index == 0 {
				if shipment.ShipmentType == models.MTOShipmentTypePPM {
					suite.Equal(fmt.Sprintf("%s %s", locator, string(shipment.ShipmentType)), value)
				} else if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
					suite.Equal(fmt.Sprintf("%s %s", locator, "NTS Release"), value)
				} else if shipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTSDom {
					suite.Equal(fmt.Sprintf("%s %s", locator, "NTS"), value)
				} else if shipment.ShipmentType == models.MTOShipmentTypeMobileHome {
					suite.Equal(fmt.Sprintf("%s %s", locator, "Mobile Home"), value)
				} else if shipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway {
					suite.Equal(fmt.Sprintf("%s %s", locator, "Boat Haul"), value)
				} else if shipment.ShipmentType == models.MTOShipmentTypeBoatTowAway {
					suite.Equal(fmt.Sprintf("%s %s", locator, "Boat Tow"), value)
				} else if shipment.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
					suite.Equal(fmt.Sprintf("%s %s", locator, "UB"), value)
				} else {
					suite.Fail(fmt.Sprintf("unaccounted type: %s", string(shipment.ShipmentType)))
				}
			}
			// verify AddShipmentPickUpDates
			if index == 1 {
				if shipment.ShipmentType == models.MTOShipmentTypePPM {
					if shipment.PPMShipment.ActualMoveDate != nil {
						suite.Equal(fmt.Sprintf("%s %s", FormatDate(now), "Actual"), value)
					} else {
						suite.Equal(fmt.Sprintf("%s %s", FormatDate(now), "Expected"), value)
					}
				} else {
					if shipment.RequestedPickupDate != nil {
						suite.Equal(fmt.Sprintf("%s %s", FormatDate(now), "Requested"), value)
					} else if shipment.ActualPickupDate != nil {
						suite.Equal(fmt.Sprintf("%s %s", FormatDate(now), "Actual"), value)
					} else if shipment.ScheduledPickupDate != nil {
						suite.Equal(fmt.Sprintf("%s %s", FormatDate(now), "Scheduled"), value)
					} else {
						suite.Equal(" - ", value)
					}
				}
			}
			// verify AddShipmentWeights
			if index == 2 {
				if shipment.ShipmentType == models.MTOShipmentTypePPM {
					if shipment.PPMShipment.EstimatedWeight != nil {
						suite.Equal("25 lbs - Estimated", value)
					} else {
						suite.Equal("3,600 lbs - Actual", value)
					}
				} else {
					if shipment.PrimeActualWeight != nil {
						suite.Equal("1,234 Actual", value)
					} else if shipment.PrimeEstimatedWeight != nil {
						suite.Equal("1,234 Estimated", value)
					} else {
						suite.Equal(" - ", value)
					}
				}
			}
			// verify AddShipmentStatus
			if index == 3 {
				suite.Equal(FormatEnum(string(shipment.Status), ""), value)
			}
		}
	}
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestTooManyShipmentsErrorFormatAdditionalShipments() {
	locator := "ABCDEF-01"

	ppm := models.PPMShipment{
		ID: uuid.Must(uuid.NewV4()),
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}

	ppm2 := models.PPMShipment{
		ID: uuid.Must(uuid.NewV4()),
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}

	var shipments []models.MTOShipment
	i := 0
	// build 18 shipments to exceed limit
	for next := true; next; next = i < 18 {
		shipments = append(shipments, models.MTOShipment{
			ShipmentType:    models.MTOShipmentTypePPM,
			ShipmentLocator: &locator,
			PPMShipment:     &ppm2,
			Status:          models.MTOShipmentStatusSubmitted,
		})
		i++
	}

	ssd := models.ShipmentSummaryFormData{
		PPMShipment:  ppm,
		AllShipments: shipments,
	}

	_, err := formatAdditionalShipments(ssd)
	suite.NotNil(err)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestMissingShipmentLocatorErrorFormatAdditionalShipments() {
	locator := "ABCDEF-01"

	ppm := models.PPMShipment{
		ID: uuid.Must(uuid.NewV4()),
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}

	ppm2 := models.PPMShipment{
		ID: uuid.Must(uuid.NewV4()),
		Shipment: models.MTOShipment{
			ShipmentLocator: &locator,
		},
	}

	shipments := []models.MTOShipment{
		{
			ShipmentType: models.MTOShipmentTypePPM,
			PPMShipment:  &ppm2,
			Status:       models.MTOShipmentStatusSubmitted,
			//No -- ShipmentLocator: &locator,
		},
	}

	ssd := models.ShipmentSummaryFormData{
		PPMShipment:  ppm,
		AllShipments: shipments,
	}

	_, err := formatAdditionalShipments(ssd)
	suite.NotNil(err)
}

func (suite *ShipmentSummaryWorksheetServiceSuite) TestFormatDisbursement() {
	expensesMap := make(map[string]float64)

	// Test case 1: GTCC calculation B is less than GTCC calculation A
	// Additionally, Member should not be less than 0
	expectedResult := "GTCC: " + FormatDollars(100.00) + "\nMember: " + FormatDollars(0)
	expensesMap["TotalGTCCPaid"] = 200.00
	expensesMap["StorageGTCCPaid"] = 300.00
	ppmRemainingEntitlement := 60.00
	expensesMap["StorageMemberPaid"] = 40.00
	result := formatDisbursement(expensesMap, ppmRemainingEntitlement)
	suite.Equal(result, expectedResult)

	// Test case 2: GTCC calculation A is less than GTCC calculation B
	expectedResult = "GTCC: " + FormatDollars(100.00) + "\nMember: " + FormatDollars(400.00)
	expensesMap = make(map[string]float64)
	expensesMap["TotalGTCCPaid"] = 60.00
	expensesMap["StorageGTCCPaid"] = 40.00
	ppmRemainingEntitlement = 300.00
	expensesMap["StorageMemberPaid"] = 200.00
	result = formatDisbursement(expensesMap, ppmRemainingEntitlement)
	suite.Equal(result, expectedResult)

	// Test case 3: GTCC calculation is less than 0
	expectedResult = "GTCC: " + FormatDollars(0) + "\nMember: " + FormatDollars(-250.00)
	expensesMap = make(map[string]float64)
	expensesMap["TotalGTCCPaid"] = 0
	expensesMap["StorageGTCCPaid"] = 0
	ppmRemainingEntitlement = -300.00
	expensesMap["StorageMemberPaid"] = 50.00
	result = formatDisbursement(expensesMap, ppmRemainingEntitlement)
	suite.Equal(result, expectedResult)
}
