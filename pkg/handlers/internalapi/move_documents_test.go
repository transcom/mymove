package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMoveDocumentHandler() {
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
	uploadIds := []strfmt.UUID{*utils.FmtUUID(upload.ID)}

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	newMoveDocPayload := internalmessages.CreateGenericMoveDocumentPayload{
		UploadIds:                uploadIds,
		PersonallyProcuredMoveID: utils.FmtUUID(ppm.ID),
		MoveDocumentType:         internalmessages.MoveDocumentTypeOTHER,
		Title:                    utils.FmtString("awesome_document.pdf"),
		Notes:                    utils.FmtString("Some notes here"),
	}

	newMoveDocParams := movedocop.CreateGenericMoveDocumentParams{
		HTTPRequest:                      request,
		CreateGenericMoveDocumentPayload: &newMoveDocPayload,
		MoveID: strfmt.UUID(move.ID.String()),
	}

	context := NewHandlerContext(suite.db, suite.logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateGenericMoveDocumentHandler(context)
	response := handler.Handle(newMoveDocParams)
	// assert we got back the 201 response
	suite.isNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateGenericMoveDocumentOK)
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
	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		ID:               utils.FmtUUID(moveDocument.ID),
		MoveID:           utils.FmtUUID(move.ID),
		Title:            utils.FmtString("super_awesome.pdf"),
		Notes:            utils.FmtString("This document is super awesome."),
		Status:           internalmessages.MoveDocumentStatusOK,
		MoveDocumentType: internalmessages.MoveDocumentTypeOTHER,
	}

	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		HTTPRequest:        request,
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}

	handler := UpdateMoveDocumentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(updateMoveDocParams)

	// Then: we expect to get back a 200 response
	suite.isNotErrResponse(response)
	updateResponse := response.(*movedocop.UpdateMoveDocumentOK)
	updatePayload := updateResponse.Payload
	suite.NotNil(updatePayload)

	suite.Require().Equal(*updatePayload.ID, strfmt.UUID(moveDocument.ID.String()), "expected move doc ids to match")

	// And: the new data to be there
	suite.Require().Equal(*updatePayload.Title, "super_awesome.pdf")
	suite.Require().Equal(*updatePayload.Notes, "This document is super awesome.")
}

func (suite *HandlerSuite) TestApproveMoveDocumentHandler() {
	// When: there is a move and move document
	ppm := testdatagen.MakePPM(suite.db, testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusPAYMENTREQUESTED,
		},
	})
	move := ppm.Move
	sm := move.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   move.ID,
			Move:                     move,
			MoveDocumentType:         models.MoveDocumentTypeSHIPMENTSUMMARY,
			PersonallyProcuredMoveID: &ppm.ID,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	})
	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	// And: the title and status are updated
	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		ID:               utils.FmtUUID(moveDocument.ID),
		MoveID:           utils.FmtUUID(move.ID),
		Title:            utils.FmtString(moveDocument.Title),
		Notes:            moveDocument.Notes,
		Status:           internalmessages.MoveDocumentStatusOK,
		MoveDocumentType: internalmessages.MoveDocumentTypeSHIPMENTSUMMARY,
	}

	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		HTTPRequest:        request,
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}

	handler := UpdateMoveDocumentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(updateMoveDocParams)

	// Then: we expect to get back a 200 response
	suite.isNotErrResponse(response)
	updateResponse := response.(*movedocop.UpdateMoveDocumentOK)
	updatePayload := updateResponse.Payload
	suite.NotNil(updatePayload)

	suite.Require().Equal(*updatePayload.ID, strfmt.UUID(moveDocument.ID.String()), "expected move doc ids to match")

	// And: the new data to be there
	suite.Require().Equal(updatePayload.Status, internalmessages.MoveDocumentStatusOK)

	var ppms models.PersonallyProcuredMoves
	q := suite.db.Where("move_id = ?", move.ID)
	q.All(&ppms)
	suite.Require().Equal(len(ppms), 1, "Should have a PPM!")
	reloadedPPM := ppms[0]
	suite.Require().Equal(string(models.PPMStatusCOMPLETED), string(reloadedPPM.Status))

}
