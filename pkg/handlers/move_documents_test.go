package handlers

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func createMoveDocumentSetup(suite *HandlerSuite) (CreateMoveDocumentHandler, movedocop.CreateMoveDocumentParams, internalmessages.CreateMoveDocumentPayload, models.PersonallyProcuredMove) {
	ppm := testdatagen.MakeDefaultPPM(suite.db)
	move := ppm.Move
	sm := move.Orders.ServiceMember

	upload := testdatagen.MakeUpload(suite.db, testdatagen.Assertions{
		Upload: models.Upload{
			UploaderID: sm.UserID,
		},
	})
	upload.DocumentID = nil
	suite.mustSave(&upload)
	uploadIds := []strfmt.UUID{*fmtUUID(upload.ID)}

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	payload := internalmessages.CreateMoveDocumentPayload{
		UploadIds:                uploadIds,
		PersonallyProcuredMoveID: fmtUUID(ppm.ID),
		MoveDocumentType:         internalmessages.MoveDocumentTypeOTHER,
		Title:                    fmtString("awesome_document.pdf"),
		Notes:                    fmtString("Some notes here"),
		Status:                   internalmessages.MoveDocumentStatusAWAITINGREVIEW,
	}

	params := movedocop.CreateMoveDocumentParams{
		HTTPRequest:               request,
		CreateMoveDocumentPayload: &payload,
		MoveID: strfmt.UUID(move.ID.String()),
	}

	context := NewHandlerContext(suite.db, suite.logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateMoveDocumentHandler(context)

	return handler, params, payload, ppm
}

func (suite *HandlerSuite) TestCreateMoveDocumentHandler() {
	handler, params, payload, _ := createMoveDocumentSetup(suite)

	response := handler.Handle(params)
	// assert we got back the 201 response
	suite.isNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateMoveDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)

	// Make sure the Upload was associated to the new document
	createdDocumentID := createdPayload.Document.ID
	var fetchedUpload models.Upload
	suite.db.Find(&fetchedUpload, payload.UploadIds[0])
	suite.Equal(createdDocumentID.String(), fetchedUpload.DocumentID.String())
}

func (suite *HandlerSuite) TestWrongUserCreateMoveDocumentHandler() {
	handler, params, _, _ := createMoveDocumentSetup(suite)

	// Next try the wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.db)
	params.HTTPRequest = suite.authenticateRequest(params.HTTPRequest, wrongUser)

	response := handler.Handle(params)
	suite.checkResponseForbidden(response)
}

func (suite *HandlerSuite) TestBadMoveCreateMoveDocumentHandler() {
	handler, params, _, _ := createMoveDocumentSetup(suite)

	params.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	response := handler.Handle(params)
	suite.checkResponseNotFound(response)
}

func (suite *HandlerSuite) TestWrongPPMCreateMoveDocumentHandler() {
	handler, params, payload, ppm := createMoveDocumentSetup(suite)

	otherPpm := testdatagen.MakePPM(suite.db, testdatagen.Assertions{
		Order: ppm.Move.Orders,
	})

	payload.PersonallyProcuredMoveID = fmtUUID(otherPpm.ID)
	params.CreateMoveDocumentPayload = &payload
	response := handler.Handle(params)
	suite.IsType(&movedocop.CreateMoveDocumentBadRequest{}, response)
}

func (suite *HandlerSuite) TestIndexMoveDocumentsHandler() {
	ppm := testdatagen.MakeDefaultPPM(suite.db)
	move := ppm.Move
	sm := move.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID: move.ID,
			Move:   move,
			PersonallyProcuredMoveID: &ppm.ID,
		},
	})

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	indexMoveDocParams := movedocop.IndexMoveDocumentsParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	context := NewHandlerContext(suite.db, suite.logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := IndexMoveDocumentsHandler(context)
	response := handler.Handle(indexMoveDocParams)

	// assert we got back the 201 response
	indexResponse := response.(*movedocop.IndexMoveDocumentsOK)
	indexPayload := indexResponse.Payload
	suite.NotNil(indexPayload)

	for _, moveDoc := range indexPayload {
		suite.Require().Equal(*moveDoc.ID, strfmt.UUID(moveDocument.ID.String()), "expected move ids to match")
		suite.Require().Equal(*moveDoc.PersonallyProcuredMoveID, strfmt.UUID(ppm.ID.String()), "expected ppm ids to match")
	}

	// Next try the wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.db)
	request = suite.authenticateRequest(request, wrongUser)
	indexMoveDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(indexMoveDocParams)
	suite.checkResponseForbidden(badUserResponse)

	// Now try a bad move
	indexMoveDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(indexMoveDocParams)
	suite.checkResponseNotFound(badMoveResponse)
}

func (suite *HandlerSuite) TestUpdateMoveDocumentHandler() {
	// When: there is a move and move document
	move := testdatagen.MakeDefaultMove(suite.db)
	sm := move.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID: move.ID,
			Move:   move,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	})
	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	// And: the title and status are updated
	updateMoveDocPayload := internalmessages.UpdateMoveDocumentPayload{
		Title:  fmtString("super_awesome.pdf"),
		Notes:  fmtString("This document is super awesome."),
		Status: internalmessages.MoveDocumentStatusOK,
	}

	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		HTTPRequest:        request,
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}

	handler := UpdateMoveDocumentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(updateMoveDocParams)

	// Then: we expect to get back a 200 response
	updateResponse := response.(*movedocop.UpdateMoveDocumentOK)
	updatePayload := updateResponse.Payload
	suite.NotNil(updatePayload)

	suite.Require().Equal(*updatePayload.ID, strfmt.UUID(moveDocument.ID.String()), "expected move doc ids to match")

	// And: the new data to be there
	suite.Require().Equal(*updatePayload.Title, "super_awesome.pdf")
	suite.Require().Equal(*updatePayload.Notes, "This document is super awesome.")

}
