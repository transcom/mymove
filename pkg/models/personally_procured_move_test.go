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

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestPPMValidation() {
	ppm := &models.PersonallyProcuredMove{}

	expErrors := map[string][]string{
		"status": {"Status can not be blank."},
	}

	suite.verifyValidationErrors(ppm, expErrors)
}

func (suite *ModelSuite) TestPPMAdvance() {

	move := testdatagen.MakeDefaultMove(suite.DB())
	serviceMember := move.Orders.ServiceMember

	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	advance.Request()
	models.SavePersonallyProcuredMove(suite.DB(), ppm)
	session := auth.Session{
		UserID:          serviceMember.User.ID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: serviceMember.ID,
	}
	fetchedPPM, err := models.FetchPersonallyProcuredMove(suite.DB(), &session, ppm.ID)
	suite.NoError(err)
	suite.Equal(fetchedPPM.Advance.Status, models.ReimbursementStatusREQUESTED, "expected Requested")
}

// TODO: Fix test now that we capture transaction error
/* func (suite *ModelSuite) TestPPMAdvanceNoGTCC() {
	move := testdatagen.MakeDefaultMove(suite.DB())

	advance := BuildDraftReimbursement(1000, MethodOfReceiptGTCC)

	_, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.NoError(err)
	suite.True(verrs.HasAny())
} */

func (suite *ModelSuite) TestPPMStateMachine() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = models.OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.MustSave(&orders)
	testdatagen.MakeDefaultContractor(suite.DB())

	selectedMoveType := models.SelectedMoveTypeHHGPPM

	moveOptions := models.MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	ppm.Status = models.PPMStatusSUBMITTED // NEVER do this outside of a test.

	// Can cancel ppm
	err = ppm.Cancel()
	suite.NoError(err)
	suite.Equal(models.PPMStatusCANCELED, ppm.Status, "expected Canceled")
}

func (suite *ModelSuite) TestFetchMoveDocumentsForTypes() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	sm := ppm.Move.Orders.ServiceMember

	assertions := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   "OK",
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	}

	testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions)

	deletedAt := time.Date(2019, 8, 7, 0, 0, 0, 0, time.UTC)
	deleteAssertions := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   "OK",
			MoveDocumentType:         "EXPENSE",
			DeletedAt:                &deletedAt,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
			DeletedAt:       &deletedAt,
		},
	}
	testdatagen.MakeMovingExpenseDocument(suite.DB(), deleteAssertions)

	docTypes := []string{"EXPENSE"}
	moveDocs, err := ppm.FetchMoveDocumentsForTypes(suite.DB(), docTypes)

	if suite.NoError(err) {
		suite.Equal(2, len(moveDocs))
	}

}

func (suite *ModelSuite) TestFetchPersonallyProcuredMoveByOrderID() {
	orderID := uuid.Must(uuid.NewV4())
	moveID, _ := uuid.FromString("7112b18b-7e03-4b28-adde-532b541bba8d")
	invalidID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

	order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID: orderID,
		},
	})
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       moveID,
			OrdersID: orderID,
			Orders:   order,
		},
	})

	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	advance.Request()
	models.SavePersonallyProcuredMove(suite.DB(), ppm)

	tests := []struct {
		lookupID  uuid.UUID
		resultID  uuid.UUID
		resultErr bool
	}{
		{lookupID: orderID, resultID: moveID, resultErr: false},
		{lookupID: invalidID, resultID: invalidID, resultErr: true},
	}

	for _, ts := range tests {
		ppm, err := models.FetchPersonallyProcuredMoveByOrderID(suite.DB(), ts.lookupID)
		if ts.resultErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
		suite.Equal(ppm.MoveID, ts.resultID, "Wrong moveID: %s", ts.lookupID)
	}
}
