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

func (suite *HandlerSuite) TestCreateMoveDocumentHandler() {
	move := testdatagen.MakeDefaultMove(suite.db)
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

	newMoveDocPayload := internalmessages.CreateMoveDocumentPayload{
		UploadIds:        uploadIds,
		MoveDocumentType: internalmessages.MoveDocumentTypeOTHER,
		Title:            fmtString("awesome_document.pdf"),
		Notes:            fmtString("Some notes here"),
		Status:           internalmessages.MoveDocumentStatusAWAITINGREVIEW,
	}

	newMoveDocParams := movedocop.CreateMoveDocumentParams{
		HTTPRequest:               request,
		CreateMoveDocumentPayload: &newMoveDocPayload,
		MoveID: strfmt.UUID(move.ID.String()),
	}

	context := NewHandlerContext(suite.db, suite.logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateMoveDocumentHandler(context)
	response := handler.Handle(newMoveDocParams)
	// assert we got back the 201 response
	suite.isNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateMoveDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)

	// Make sure the Upload was associated to the new document
	createdDocumentID := createdPayload.Document.ID
	var fetchedUpload models.Upload
	suite.db.Find(&fetchedUpload, upload.ID)
	suite.Equal(createdDocumentID.String(), fetchedUpload.DocumentID.String())

	// Next try the wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.db)
	request = suite.authenticateRequest(request, wrongUser)
	newMoveDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(newMoveDocParams)
	suite.checkResponseForbidden(badUserResponse)

	// Now try a bad move
	newMoveDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newMoveDocParams)
	suite.checkResponseNotFound(badMoveResponse)
}

func (suite *HandlerSuite) TestIndexMoveDocumentsHandler() {
	move1 := testdatagen.MakeDefaultMove(suite.db)
	sm := move1.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID: move1.ID,
			Move:   move1,
		},
	})

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	indexMoveDocParams := movedocop.IndexMoveDocumentsParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move1.ID.String()),
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
		Title: fmtString("super_awesome.pdf"),
		Notes: fmtString("This document is super awesome."),
	}

	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		HTTPRequest:        request,
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveID:             strfmt.UUID(move.ID.String()),
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

	// When: The wrong move id is passed in, expect 404
	updateMoveDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(updateMoveDocParams)

	_, ok := badMoveResponse.(*movedocop.UpdateMoveDocumentBadRequest)
	if !ok {
		suite.T().Fatalf("Request failed: %#v", badMoveResponse)
	}
}
