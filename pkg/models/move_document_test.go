package models_test

import (
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/auth"
	. "github.com/transcom/mymove/pkg/models"
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

func (suite *ModelSuite) TestFetchMoveDocumentsByTypeForShipment() {
	// When: There is a move, shipment, and move document of type GBL
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	sm := shipment.Move.Orders.ServiceMember

	assertions := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:           shipment.Move.ID,
			Move:             shipment.Move,
			ShipmentID:       &shipment.ID,
			Status:           "OK",
			MoveDocumentType: "GOV_BILL_OF_LADING",
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
		TransportationServiceProvider: models.TransportationServiceProvider{
			ID: tspUser.TransportationServiceProvider.ID,
		},
		Shipment: models.Shipment{
			ID: shipment.ID,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Shipment:                        shipment,
			ShipmentID:                      shipment.ID,
		},
	}

	testdatagen.MakeMoveDocument(suite.DB(), assertions)
	testdatagen.MakeMoveDocument(suite.DB(), assertions)
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)
	// When: the logged in user is a TSP user
	session := &auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		TspUserID:       tspUser.ID,
	}

	moveDocs, err := FetchMoveDocumentsByTypeForShipment(suite.DB(), session, MoveDocumentTypeGOVBILLOFLADING, shipment.ID)

	if suite.NoError(err) {
		suite.Equal(2, len(moveDocs))
		for _, moveDoc := range moveDocs {
			suite.Equal(moveDoc.MoveDocumentType, MoveDocumentTypeGOVBILLOFLADING)
			suite.Equal(moveDoc.Status, MoveDocumentStatusOK)
			suite.Equal(*moveDoc.ShipmentID, shipment.ID)
			suite.Equal(moveDoc.MoveID, shipment.Move.ID)
		}
	}

	// When: a document doesn't exist
	nonExistantDocs, err := FetchMoveDocumentsByTypeForShipment(suite.DB(), session, MoveDocumentTypeSHIPMENTSUMMARY, shipment.ID)
	// Then: No docs should be returned
	if suite.NoError(err) {
		suite.Equal(0, len(nonExistantDocs))
	}
	// When: a user without authority is logged in
	session = &auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          sm.UserID,
		TspUserID:       sm.ID,
	}
	_, err = FetchMoveDocumentsByTypeForShipment(suite.DB(), session, MoveDocumentTypeSHIPMENTSUMMARY, shipment.ID)
	// Then: FetchForbiddenError should be returned
	suite.Equal("FETCH_NOT_FOUND", err.Error())
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

	// User is authorized to fetch move doc
	session := &auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
	}

	moveDocs, err := FetchApprovedMovingExpenseDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID)

	if suite.NoError(err) {
		suite.Equal(2, len(moveDocs))
		for _, moveDoc := range moveDocs {
			suite.Equal(moveDoc.MoveDocumentType, MoveDocumentTypeEXPENSE)
			suite.Equal(moveDoc.Status, MoveDocumentStatusOK)
			suite.Equal(moveDoc.MoveID, ppm.Move.ID)
			suite.Equal((&moveDoc.MovingExpenseDocument.RequestedAmountCents).Int(), 2589)
		}
	}

	// When: the user is not authorized to fetch movedocs
	session.UserID = uuid.Must(uuid.NewV4())
	session.ServiceMemberID = uuid.Must(uuid.NewV4())
	_, err = FetchApprovedMovingExpenseDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID)
	if suite.Error(err) {
		suite.Equal(ErrFetchForbidden, err)
	}

	// When: the logged in user is an office user
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session.UserID = *officeUser.UserID
	session.OfficeUserID = officeUser.ID
	session.ApplicationName = auth.OfficeApp

	moveDocsOffice, err := FetchApprovedMovingExpenseDocuments(suite.DB(), session, ppm.Move.PersonallyProcuredMoves[0].ID)
	if suite.NoError(err) {
		for _, moveDoc := range moveDocsOffice {
			suite.Equal(moveDoc.MoveID, ppm.Move.ID)
			suite.Equal(moveDoc.Status, MoveDocumentStatusOK)
			suite.Equal(moveDoc.MoveDocumentType, MoveDocumentTypeEXPENSE)
		}
	}
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
		ApplicationName: auth.MyApp,
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
	}

	moveDoc, err := FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	if suite.NoError(err) {
		suite.Equal(moveDocument.MoveID, moveDoc.MoveID)
	}

	// When: the user is not authorized to fetch movedoc
	session.UserID = uuid.Must(uuid.NewV4())
	session.ServiceMemberID = uuid.Must(uuid.NewV4())
	_, err = FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	if suite.Error(err) {
		suite.Equal(ErrFetchForbidden, err)
	}

	// When: the logged in user is an office user
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session.UserID = *officeUser.UserID
	session.OfficeUserID = officeUser.ID
	session.ApplicationName = auth.OfficeApp

	// Then: move document is returned
	moveDocOfficeUser, err := FetchMoveDocument(suite.DB(), session, moveDocument.ID)
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

	suite.Equal(moveDocument.Status, MoveDocumentStatusAWAITINGREVIEW)

	err := moveDocument.Approve()
	suite.NoError(err)

	err = moveDocument.Reject()
	suite.NoError(err)

	err = moveDocument.Approve()
	suite.NoError(err)

	err = moveDocument.Approve()
	suite.Error(err)

	// JUST for testing, resetting Status by hand.
	moveDocument.Status = MoveDocumentStatusAWAITINGREVIEW

	err = moveDocument.AttemptTransition(MoveDocumentStatusOK)
	suite.NoError(err)
	suite.Equal(moveDocument.Status, MoveDocumentStatusOK)

	err = moveDocument.AttemptTransition(MoveDocumentStatusHASISSUE)
	suite.NoError(err)
	suite.Equal(moveDocument.Status, MoveDocumentStatusHASISSUE)

	err = moveDocument.AttemptTransition(MoveDocumentStatusOK)
	suite.NoError(err)
	suite.Equal(moveDocument.Status, MoveDocumentStatusOK)

	err = moveDocument.AttemptTransition(MoveDocumentStatusOK)
	suite.NoError(err)
	suite.Equal(moveDocument.Status, MoveDocumentStatusOK)

}
