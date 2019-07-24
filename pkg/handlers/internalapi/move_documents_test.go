package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMoveDocumentHandler() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := move.Orders.ServiceMember

	upload := testdatagen.MakeUpload(suite.DB(), testdatagen.Assertions{
		Upload: models.Upload{
			UploaderID: sm.UserID,
		},
	})
	upload.DocumentID = nil
	suite.MustSave(&upload)
	uploadIds := []strfmt.UUID{*handlers.FmtUUID(upload.ID)}

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)

	newMoveDocPayload := internalmessages.CreateGenericMoveDocumentPayload{
		UploadIds:                uploadIds,
		PersonallyProcuredMoveID: handlers.FmtUUID(ppm.ID),
		MoveDocumentType:         internalmessages.MoveDocumentTypeOTHER,
		Title:                    handlers.FmtString("awesome_document.pdf"),
		Notes:                    handlers.FmtString("Some notes here"),
	}

	newMoveDocParams := movedocop.CreateGenericMoveDocumentParams{
		HTTPRequest:                      request,
		CreateGenericMoveDocumentPayload: &newMoveDocPayload,
		MoveID:                           strfmt.UUID(move.ID.String()),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateGenericMoveDocumentHandler{context}
	response := handler.Handle(newMoveDocParams)
	// assert we got back the 201 response
	suite.IsNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateGenericMoveDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)

	// Make sure the Upload was associated to the new document
	createdDocumentID := createdPayload.Document.ID
	var fetchedUpload models.Upload
	suite.DB().Find(&fetchedUpload, upload.ID)
	suite.Equal(createdDocumentID.String(), fetchedUpload.DocumentID.String())

	// Next try the wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	request = suite.AuthenticateRequest(request, wrongUser)
	newMoveDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(newMoveDocParams)
	suite.CheckResponseForbidden(badUserResponse)

	// Now try a bad move
	newMoveDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newMoveDocParams)
	suite.CheckResponseNotFound(badMoveResponse)
}

func (suite *HandlerSuite) TestIndexMoveDocumentsHandler() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := move.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.DB(), testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   move.ID,
			Move:                     move,
			PersonallyProcuredMoveID: &ppm.ID,
		},
	})

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)

	indexMoveDocParams := movedocop.IndexMoveDocumentsParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := IndexMoveDocumentsHandler{context}
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
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	request = suite.AuthenticateRequest(request, wrongUser)
	indexMoveDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(indexMoveDocParams)
	suite.CheckResponseForbidden(badUserResponse)

	// Now try a bad move
	indexMoveDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(indexMoveDocParams)
	suite.CheckResponseNotFound(badMoveResponse)
}

func (suite *HandlerSuite) TestIndexWeightTicketSetDocumentsHandlerNoMissingFields() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := move.Orders.ServiceMember

	moveDoc := testdatagen.MakeMoveDocument(suite.DB(),
		testdatagen.Assertions{
			MoveDocument: models.MoveDocument{
				MoveID:                   move.ID,
				Move:                     move,
				PersonallyProcuredMoveID: &ppm.ID,
				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
			},
		})
	emptyWeight := unit.Pound(1000)
	fullWeight := unit.Pound(2500)
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDoc.ID,
		MoveDocument:             moveDoc,
		EmptyWeight:              &emptyWeight,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight,
		FullWeightTicketMissing:  false,
		VehicleNickname:          "My Car",
		VehicleOptions:           "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)

	indexMoveDocParams := movedocop.IndexMoveDocumentsParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := IndexMoveDocumentsHandler{context}
	response := handler.Handle(indexMoveDocParams)

	// assert we got back the 201 response
	indexResponse := response.(*movedocop.IndexMoveDocumentsOK)
	indexPayload := indexResponse.Payload
	suite.NotNil(indexPayload)
	for _, moveDoc := range indexPayload {
		suite.Require().Equal(*moveDoc.ID, strfmt.UUID(weightTicketSetDocument.MoveDocument.ID.String()), "expected move ids to match")
		suite.Require().Equal(*moveDoc.PersonallyProcuredMoveID, strfmt.UUID(ppm.ID.String()), "expected ppm ids to match")
		suite.Require().Equal(*moveDoc.EmptyWeight, int64(*weightTicketSetDocument.EmptyWeight), "expected empty weight to match")
		suite.Require().Equal(*moveDoc.EmptyWeightTicketMissing, weightTicketSetDocument.EmptyWeightTicketMissing, "expected empty weight ticket missing to match")
		suite.Require().Equal(*moveDoc.FullWeight, int64(*weightTicketSetDocument.FullWeight), "expected empty weight to match")
		suite.Require().Equal(*moveDoc.FullWeightTicketMissing, weightTicketSetDocument.FullWeightTicketMissing, "expected full weight ticket missing to match")
		suite.Require().Equal(moveDoc.WeightTicketDate.String(), strfmt.Date(*weightTicketSetDocument.WeightTicketDate).String(), "expected weight ticket date to match")
		suite.Require().Equal(*moveDoc.TrailerOwnershipMissing, weightTicketSetDocument.TrailerOwnershipMissing, "expected trailer ownership missing to match")
		suite.Require().Equal(moveDoc.VehicleOptions, weightTicketSetDocument.VehicleOptions, "expected vehicle options to match")
		suite.Require().Equal(moveDoc.VehicleNickname, weightTicketSetDocument.VehicleNickname, "expected vehicle nickname to match")
	}
}

func (suite *HandlerSuite) TestIndexWeightTicketSetDocumentsHandlerMissingFields() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := move.Orders.ServiceMember

	moveDoc := testdatagen.MakeMoveDocument(suite.DB(),
		testdatagen.Assertions{
			MoveDocument: models.MoveDocument{
				MoveID:                   move.ID,
				Move:                     move,
				PersonallyProcuredMoveID: &ppm.ID,
				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
			},
		})
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDoc.ID,
		MoveDocument:             moveDoc,
		EmptyWeight:              nil,
		EmptyWeightTicketMissing: true,
		FullWeight:               nil,
		FullWeightTicketMissing:  true,
		VehicleNickname:          "My Car",
		VehicleOptions:           "CAR",
		WeightTicketDate:         nil,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)

	indexMoveDocParams := movedocop.IndexMoveDocumentsParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := IndexMoveDocumentsHandler{context}
	response := handler.Handle(indexMoveDocParams)

	// assert we got back the 201 response
	indexResponse := response.(*movedocop.IndexMoveDocumentsOK)
	indexPayload := indexResponse.Payload
	suite.NotNil(indexPayload)
	for _, moveDoc := range indexPayload {
		suite.Require().Equal(*moveDoc.ID, strfmt.UUID(weightTicketSetDocument.MoveDocument.ID.String()), "expected move ids to match")
		suite.Require().Equal(*moveDoc.PersonallyProcuredMoveID, strfmt.UUID(ppm.ID.String()), "expected ppm ids to match")
		suite.Require().Nil(moveDoc.EmptyWeight)
		suite.Require().Equal(*moveDoc.EmptyWeightTicketMissing, weightTicketSetDocument.EmptyWeightTicketMissing, "expected empty weight ticket missing to match")
		suite.Require().Nil(moveDoc.FullWeight)
		suite.Require().Equal(*moveDoc.FullWeightTicketMissing, weightTicketSetDocument.FullWeightTicketMissing, "expected full weight ticket missing to match")
		suite.Require().Nil(moveDoc.WeightTicketDate)
		suite.Require().Equal(*moveDoc.TrailerOwnershipMissing, weightTicketSetDocument.TrailerOwnershipMissing, "expected trailer ownership missing to match")
		suite.Require().Equal(moveDoc.VehicleOptions, weightTicketSetDocument.VehicleOptions, "expected vehicle options to match")
		suite.Require().Equal(moveDoc.VehicleNickname, weightTicketSetDocument.VehicleNickname, "expected vehicle nickname to match")
	}
}

func (suite *HandlerSuite) TestUpdateMoveDocumentHandler() {
	// When: there is a move and move document
	move := testdatagen.MakeDefaultMove(suite.DB())
	sm := move.Orders.ServiceMember
	moveDocument := testdatagen.MakeMoveDocumentWeightTicketSet(suite.DB(), testdatagen.Assertions{
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
	request = suite.AuthenticateRequest(request, sm)

	emptyWeight := (int64)(200)
	fullWeight := (int64)(500)

	// And: the title and status are updated
	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
		MoveID:           handlers.FmtUUID(move.ID),
		Title:            handlers.FmtString("super_awesome.pdf"),
		Notes:            handlers.FmtString("This document is super awesome."),
		Status:           internalmessages.MoveDocumentStatusOK,
		MoveDocumentType: internalmessages.MoveDocumentTypeWEIGHTTICKETSET,
		EmptyWeight:      &emptyWeight,
		FullWeight:       &fullWeight,
	}

	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		HTTPRequest:        request,
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}

	handler := UpdateMoveDocumentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(updateMoveDocParams)

	// Then: we expect to get back a 200 response
	suite.IsNotErrResponse(response)
	updateResponse := response.(*movedocop.UpdateMoveDocumentOK)
	updatePayload := updateResponse.Payload
	suite.NotNil(updatePayload)

	suite.Require().Equal(*updatePayload.ID, strfmt.UUID(moveDocument.ID.String()), "expected move doc ids to match")

	// And: the new data to be there
	suite.Require().Equal(*updatePayload.Title, "super_awesome.pdf")
	suite.Require().Equal(*updatePayload.Notes, "This document is super awesome.")
	suite.Require().Equal(*updatePayload.EmptyWeight, int64(200))
	suite.Require().Equal(*updatePayload.FullWeight, int64(500))
}
func (suite *HandlerSuite) TestApproveMoveDocumentHandler() {
	// When: there is a move and move document
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusPAYMENTREQUESTED,
		},
	})
	move := ppm.Move
	sm := move.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.DB(), testdatagen.Assertions{
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
	request = suite.AuthenticateRequest(request, sm)

	// And: the title and status are updated
	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
		MoveID:           handlers.FmtUUID(move.ID),
		Title:            handlers.FmtString(moveDocument.Title),
		Notes:            moveDocument.Notes,
		Status:           internalmessages.MoveDocumentStatusOK,
		MoveDocumentType: internalmessages.MoveDocumentTypeSHIPMENTSUMMARY,
	}

	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		HTTPRequest:        request,
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}

	handler := UpdateMoveDocumentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(updateMoveDocParams)

	// Then: we expect to get back a 200 response
	suite.IsNotErrResponse(response)
	updateResponse := response.(*movedocop.UpdateMoveDocumentOK)
	updatePayload := updateResponse.Payload
	suite.NotNil(updatePayload)

	suite.Require().Equal(*updatePayload.ID, strfmt.UUID(moveDocument.ID.String()), "expected move doc ids to match")

	// And: the new data to be there
	suite.Require().Equal(updatePayload.Status, internalmessages.MoveDocumentStatusOK)

	var ppms models.PersonallyProcuredMoves
	q := suite.DB().Where("move_id = ?", move.ID)
	q.All(&ppms)
	suite.Require().Equal(len(ppms), 1, "Should have a PPM!")
	reloadedPPM := ppms[0]
	suite.Require().Equal(string(models.PPMStatusCOMPLETED), string(reloadedPPM.Status))

}
