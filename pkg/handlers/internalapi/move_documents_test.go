package internalapi

import (
	"net/http"
	"net/http/httptest"

	"github.com/gobuffalo/validate/v3"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

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

	userUpload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
		UserUpload: models.UserUpload{
			UploaderID: sm.UserID,
		},
	})
	userUpload.DocumentID = nil
	suite.MustSave(&userUpload)
	uploadIds := []strfmt.UUID{*handlers.FmtUUID(userUpload.Upload.ID)}

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)

	moveDocumentType := internalmessages.MoveDocumentTypeOTHER
	newMoveDocPayload := internalmessages.CreateGenericMoveDocumentPayload{
		UploadIds:                uploadIds,
		PersonallyProcuredMoveID: handlers.FmtUUID(ppm.ID),
		MoveDocumentType:         &moveDocumentType,
		Title:                    handlers.FmtString("awesome_document.pdf"),
		Notes:                    handlers.FmtString("Some notes here"),
	}

	newMoveDocParams := movedocop.CreateGenericMoveDocumentParams{
		HTTPRequest:                      request,
		CreateGenericMoveDocumentPayload: &newMoveDocPayload,
		MoveID:                           strfmt.UUID(move.ID.String()),
	}

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := CreateGenericMoveDocumentHandler{handlerConfig}
	response := handler.Handle(newMoveDocParams)
	// assert we got back the 201 response
	suite.IsNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateGenericMoveDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)

	// Make sure the UserUpload was associated to the new document
	createdDocumentID := createdPayload.Document.ID
	var fetchedUpload models.UserUpload
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
	//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
	//RA: in a unit test, then there is no risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	suite.DB().Find(&fetchedUpload, userUpload.ID)
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

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := IndexMoveDocumentsHandler{handlerConfig}
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

	vehicleNickname := "My Car"
	emptyWeight := unit.Pound(1000)
	fullWeight := unit.Pound(2500)
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDoc.ID,
		MoveDocument:             moveDoc,
		EmptyWeight:              &emptyWeight,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight,
		FullWeightTicketMissing:  false,
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      "CAR",
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

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := IndexMoveDocumentsHandler{handlerConfig}
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
		suite.Require().Equal(*moveDoc.WeightTicketSetType, internalmessages.WeightTicketSetType(weightTicketSetDocument.WeightTicketSetType), "expected vehicle options to match")
		suite.Require().Equal(moveDoc.VehicleNickname, weightTicketSetDocument.VehicleNickname, "expected vehicle nickname to match")
	}
}

func (suite *HandlerSuite) TestIndexWeightTicketSetDocumentsHandlerMissingFields() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := move.Orders.ServiceMember
	vehicleNickname := "My Car"

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
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      "CAR",
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

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := IndexMoveDocumentsHandler{handlerConfig}
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
		suite.Require().Equal(*moveDoc.WeightTicketSetType, internalmessages.WeightTicketSetType(weightTicketSetDocument.WeightTicketSetType), "expected vehicle options to match")
		suite.Require().Equal(moveDoc.VehicleNickname, weightTicketSetDocument.VehicleNickname, "expected vehicle nickname to match")
	}
}

func (suite *HandlerSuite) TestUpdateMoveDocumentHandler() {
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

	status := internalmessages.MoveDocumentStatusOK
	moveDocumentType := internalmessages.MoveDocumentTypeSHIPMENTSUMMARY
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
		MoveID:           handlers.FmtUUID(move.ID),
		Title:            handlers.FmtString(moveDocument.Title),
		Notes:            moveDocument.Notes,
		Status:           &status,
		MoveDocumentType: &moveDocumentType,
	}

	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		HTTPRequest:        request,
		UpdateMoveDocument: updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}

	moveDocumentUpdateHandler := &mocks.MoveDocumentUpdater{}

	handler := UpdateMoveDocumentHandler{
		handlers.NewHandlerConfig(suite.DB(),
			suite.Logger()),
		moveDocumentUpdateHandler,
	}

	// happy path
	returnedMoveDocument := models.MoveDocument{ID: moveDocument.ID}
	moveDocumentUpdateHandler.On("Update",
		mock.AnythingOfType("*appcontext.appContext"),
		updateMoveDocPayload,
		moveDocument.ID,
	).Return(&returnedMoveDocument, validate.NewErrors(), nil).Once()

	response := handler.Handle(updateMoveDocParams)

	suite.Assertions.IsType(&movedocop.UpdateMoveDocumentOK{}, response)
	responsePayload := response.(*movedocop.UpdateMoveDocumentOK).Payload
	suite.Equal(moveDocument.ID.String(), responsePayload.ID.String())

	// error scenario
	expectedError := models.ErrFetchForbidden
	returnedMoveDocument = models.MoveDocument{ID: moveDocument.ID}
	moveDocumentUpdateHandler.On("Update",
		mock.AnythingOfType("*appcontext.appContext"),
		updateMoveDocPayload,
		moveDocument.ID,
	).Return(&returnedMoveDocument, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(updateMoveDocParams)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}

func (suite *HandlerSuite) TestDeleteMoveDocumentHandler() {
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
	request := httptest.NewRequest("DELETE", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)

	deleteMoveDocParams := movedocop.DeleteMoveDocumentParams{
		HTTPRequest:    request,
		MoveDocumentID: strfmt.UUID(moveDocument.ID.String()),
	}

	handler := DeleteMoveDocumentHandler{
		handlers.NewHandlerConfig(suite.DB(),
			suite.Logger()),
	}

	response := handler.Handle(deleteMoveDocParams)
	errorResponse := response.(*handlers.ErrResponse)
	suite.Equal(http.StatusInternalServerError, errorResponse.Code)
	suite.Equal("Can only delete weight ticket set and expense documents", errorResponse.Err.Error())

	moveDocument2 := testdatagen.MakeMoveDocument(suite.DB(), testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   move.ID,
			Move:                     move,
			MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
			PersonallyProcuredMoveID: &ppm.ID,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	})

	vehicleNickname := "My Car"
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDocument2.ID,
		MoveDocument:             moveDocument2,
		EmptyWeight:              nil,
		EmptyWeightTicketMissing: true,
		FullWeight:               nil,
		FullWeightTicketMissing:  true,
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      "CAR",
		WeightTicketDate:         nil,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	request2 := httptest.NewRequest("DELETE", "/fake/path", nil)
	request2 = suite.AuthenticateRequest(request2, sm)

	deleteMoveDocParams2 := movedocop.DeleteMoveDocumentParams{
		HTTPRequest:    request2,
		MoveDocumentID: strfmt.UUID(moveDocument2.ID.String()),
	}

	successResponse := handler.Handle(deleteMoveDocParams2)

	suite.Assertions.IsType(&movedocop.DeleteMoveDocumentNoContent{}, successResponse)
}
