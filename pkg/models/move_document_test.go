package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicMoveDocumentInstantiation() {
	moveDoc := &models.MoveDocument{}

	expErrors := map[string][]string{
		"document_id":        {"DocumentID can not be blank."},
		"move_id":            {"MoveID can not be blank."},
		"move_document_type": {"MoveDocumentType can not be blank."},
		"status":             {"Status can not be blank."},
		"title":              {"Title can not be blank."},
	}

	suite.verifyValidationErrors(moveDoc, expErrors)
}

func (suite *ModelSuite) TestFetchApprovedMovingExpenseDocuments() {
	// When: There is a move, ppm, move document and 2 expense docs
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

	// User is authorized to fetch move doc
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
	}

	status := models.MoveDocumentStatusOK
	moveDocs, err := models.FetchMoveDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID, &status, models.MoveDocumentTypeEXPENSE, false)

	if suite.NoError(err) {
		suite.Equal(2, len(moveDocs))
		for _, moveDoc := range moveDocs {
			suite.Equal(moveDoc.MoveDocumentType, models.MoveDocumentTypeEXPENSE)
			suite.Equal(moveDoc.Status, models.MoveDocumentStatusOK)
			suite.Equal(moveDoc.MoveID, ppm.Move.ID)
			suite.Equal((&moveDoc.MovingExpenseDocument.RequestedAmountCents).Int(), 2589)
		}
	}

	allMoveDocs, err2 := models.FetchMoveDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID, &status, models.MoveDocumentTypeEXPENSE, true)
	if suite.NoError(err2) {
		suite.Equal(3, len(allMoveDocs))
	}
	// When: the user is not authorized to fetch movedocs
	session.UserID = uuid.Must(uuid.NewV4())
	session.ServiceMemberID = uuid.Must(uuid.NewV4())
	_, err = models.FetchMoveDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID, &status, models.MoveDocumentTypeEXPENSE, false)
	if suite.Error(err) {
		suite.Equal(models.ErrFetchForbidden, err)
	}

	// When: the logged in user is an office user
	officeUser := testdatagen.MakeStubbedOfficeUser(suite.DB())
	session.UserID = *officeUser.UserID
	session.OfficeUserID = officeUser.ID
	session.ApplicationName = auth.OfficeApp

	moveDocsOffice, err := models.FetchMoveDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID, &status, models.MoveDocumentTypeEXPENSE, false)
	if suite.NoError(err) {
		for _, moveDoc := range moveDocsOffice {
			suite.Equal(moveDoc.MoveID, ppm.Move.ID)
			suite.Equal(moveDoc.Status, models.MoveDocumentStatusOK)
			suite.Equal(moveDoc.MoveDocumentType, models.MoveDocumentTypeEXPENSE)
		}
	}
}

func (suite *ModelSuite) TestFetchMovingExpenseDocumentsStorageExpense() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	sm := ppm.Move.Orders.ServiceMember
	start := time.Date(2016, 01, 01, 0, 0, 0, 0, time.UTC)
	end := time.Date(2016, 01, 16, 0, 0, 0, 0, time.UTC)
	storageExpense := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   "OK",
			MoveDocumentType:         "EXPENSE",
		},
		MovingExpenseDocument: models.MovingExpenseDocument{
			MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
			RequestedAmountCents: 100,
			PaymentMethod:        "GTCC",
			ReceiptMissing:       false,
			StorageStartDate:     &start,
			StorageEndDate:       &end,
		},
	}
	testdatagen.MakeMovingExpenseDocument(suite.DB(), storageExpense)
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
	}

	expenses, err := models.FetchMoveDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID, nil, models.MoveDocumentTypeEXPENSE, false)

	if suite.NoError(err) {
		suite.Equal(1, len(expenses))
		for _, moveDoc := range expenses {
			suite.Equal(moveDoc.MoveID, ppm.Move.ID)
			suite.Equal(*moveDoc.PersonallyProcuredMoveID, ppm.ID)
			suite.Equal(moveDoc.MovingExpenseDocument.StorageStartDate.UTC(), start)
			suite.Equal(moveDoc.MovingExpenseDocument.StorageEndDate.UTC(), end)
		}
	}
}

func (suite *ModelSuite) TestFetchMovingExpenseDocuments() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	sm := ppm.Move.Orders.ServiceMember
	awaitingReview := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   models.MoveDocumentStatusAWAITINGREVIEW,
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	}
	status := models.MoveDocumentStatusOK
	ok := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   status,
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	}
	testdatagen.MakeMovingExpenseDocument(suite.DB(), awaitingReview)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), ok)
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
	}

	allExpenses, err := models.FetchMoveDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID, nil, models.MoveDocumentTypeEXPENSE, false)
	if suite.NoError(err) {
		suite.Equal(2, len(allExpenses))
		for _, moveDoc := range allExpenses {
			suite.Equal(moveDoc.MoveID, ppm.Move.ID)
			suite.Equal(*moveDoc.PersonallyProcuredMoveID, ppm.ID)
		}
	}

	approvedExpenses, err := models.FetchMoveDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID, &status, models.MoveDocumentTypeEXPENSE, false)
	if suite.NoError(err) {
		suite.Equal(1, len(approvedExpenses))
		for _, moveDoc := range approvedExpenses {
			suite.Equal(moveDoc.MoveID, ppm.Move.ID)
			suite.Equal(*moveDoc.PersonallyProcuredMoveID, ppm.ID)
		}
	}

}

func (suite *ModelSuite) TestFetchMovingExpenseDocumentsAuth() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	sm := ppm.Move.Orders.ServiceMember
	officeUser := testdatagen.MakeStubbedOfficeUser(suite.DB())
	authorizedSession := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
	}
	officeSession := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
		ServiceMemberID: sm.ID,
	}
	unauthorizedSession := &auth.Session{
		ApplicationName: auth.MilApp,
	}

	_, err1 := models.FetchMoveDocuments(suite.DB(), authorizedSession, ppm.Move.PersonallyProcuredMoves[0].ID, nil, models.MoveDocumentTypeEXPENSE, false)
	_, err2 := models.FetchMoveDocuments(suite.DB(), officeSession, ppm.Move.PersonallyProcuredMoves[0].ID, nil, models.MoveDocumentTypeEXPENSE, false)
	_, err3 := models.FetchMoveDocuments(suite.DB(), unauthorizedSession, ppm.Move.PersonallyProcuredMoves[0].ID, nil, models.MoveDocumentTypeEXPENSE, false)

	suite.Nil(err1)
	suite.Nil(err2)
	suite.Equal(models.ErrFetchForbidden, err3)
}

func (suite *ModelSuite) TestFetchMoveDocument() {
	// When: there is a move and move document
	move := testdatagen.MakeDefaultMove(suite.DB())
	sm := move.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.DB(), testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID: move.ID,
			Move:   move,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	})
	// User is authorized to fetch move doc
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
	}

	moveDoc, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	if suite.NoError(err) {
		suite.Equal(moveDocument.MoveID, moveDoc.MoveID)
	}

	// When: the user is not authorized to fetch movedoc
	session.UserID = uuid.Must(uuid.NewV4())
	session.ServiceMemberID = uuid.Must(uuid.NewV4())
	_, err = models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	if suite.Error(err) {
		suite.Equal(models.ErrFetchForbidden, err)
	}

	// When: the logged in user is an office user
	officeUser := testdatagen.MakeStubbedOfficeUser(suite.DB())
	session.UserID = *officeUser.UserID
	session.OfficeUserID = officeUser.ID
	session.ApplicationName = auth.OfficeApp

	// Then: move document is returned
	moveDocOfficeUser, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	if suite.NoError(err) {
		suite.Equal(moveDocOfficeUser.MoveID, moveDoc.MoveID)
	}
}

func (suite *ModelSuite) TestMoveDocumentStatuses() {
	// When: there is a move and move document
	move := testdatagen.MakeDefaultMove(suite.DB())
	sm := move.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.DB(), testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID: move.ID,
			Move:   move,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	})

	suite.Equal(moveDocument.Status, models.MoveDocumentStatusAWAITINGREVIEW)

	err := moveDocument.Approve()
	suite.NoError(err)

	err = moveDocument.Reject()
	suite.NoError(err)

	err = moveDocument.Approve()
	suite.NoError(err)

	err = moveDocument.Approve()
	suite.Error(err)

	// JUST for testing, resetting Status by hand.
	moveDocument.Status = models.MoveDocumentStatusAWAITINGREVIEW

	err = moveDocument.AttemptTransition(models.MoveDocumentStatusOK)
	suite.NoError(err)
	suite.Equal(moveDocument.Status, models.MoveDocumentStatusOK)

	err = moveDocument.AttemptTransition(models.MoveDocumentStatusHASISSUE)
	suite.NoError(err)
	suite.Equal(moveDocument.Status, models.MoveDocumentStatusHASISSUE)

	err = moveDocument.AttemptTransition(models.MoveDocumentStatusOK)
	suite.NoError(err)
	suite.Equal(moveDocument.Status, models.MoveDocumentStatusOK)

	err = moveDocument.AttemptTransition(models.MoveDocumentStatusOK)
	suite.NoError(err)
	suite.Equal(moveDocument.Status, models.MoveDocumentStatusOK)

}

func (suite *ModelSuite) TestDeleteMoveDocument() {
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

	expenseDoc := testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions)
	moveDocument := expenseDoc.MoveDocument
	suite.Nil(expenseDoc.DeletedAt)
	suite.Nil(moveDocument.DeletedAt)

	err := models.DeleteMoveDocument(suite.DB(), &moveDocument)

	if suite.NoError(err) {
		suite.NotNil(moveDocument.DeletedAt)
		suite.NotNil(moveDocument.MovingExpenseDocument.DeletedAt)
	}
}
