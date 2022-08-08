package internalapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMovingExpenseDocumentHandler() {

	move := testdatagen.MakeDefaultMove(suite.DB())
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
	movingExpenseType := internalmessages.MovingExpenseTypeWEIGHINGFEES
	newMovingExpenseDocPayload := internalmessages.CreateMovingExpenseDocumentPayload{
		UploadIds:            uploadIds,
		MoveDocumentType:     &moveDocumentType,
		Title:                handlers.FmtString("awesome_document.pdf"),
		Notes:                handlers.FmtString("Some notes here"),
		MovingExpenseType:    &movingExpenseType,
		PaymentMethod:        handlers.FmtString("GTCC"),
		RequestedAmountCents: handlers.FmtInt64(2589),
	}

	newMovingExpenseDocParams := movedocop.CreateMovingExpenseDocumentParams{
		HTTPRequest:                        request,
		CreateMovingExpenseDocumentPayload: &newMovingExpenseDocPayload,
		MoveID:                             strfmt.UUID(move.ID.String()),
	}

	handlerConfig := suite.HandlerConfig()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := CreateMovingExpenseDocumentHandler{handlerConfig}
	response := handler.Handle(newMovingExpenseDocParams)
	// assert we got back the 201 response
	suite.IsNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateMovingExpenseDocumentOK)
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

	// Check that the status is correct
	suite.Equal(*createdPayload.Status, internalmessages.MoveDocumentStatusAWAITINGREVIEW)

	// Next try the wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	request = suite.AuthenticateRequest(request, wrongUser)
	newMovingExpenseDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(newMovingExpenseDocParams)
	suite.CheckResponseForbidden(badUserResponse)

	// Now try a bad move
	newMovingExpenseDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newMovingExpenseDocParams)
	suite.CheckResponseNotFound(badMoveResponse)

}

func (suite *HandlerSuite) TestCreateMovingExpenseDocumentHandlerReceiptMissingNoUploads() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	sm := move.Orders.ServiceMember
	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)
	moveDocumentType := internalmessages.MoveDocumentTypeOTHER
	movingExpenseType := internalmessages.MovingExpenseTypeWEIGHINGFEES
	newMovingExpenseDocPayload := internalmessages.CreateMovingExpenseDocumentPayload{
		MoveDocumentType:     &moveDocumentType,
		Title:                handlers.FmtString("awesome_document.pdf"),
		Notes:                handlers.FmtString("Some notes here"),
		MovingExpenseType:    &movingExpenseType,
		PaymentMethod:        handlers.FmtString("GTCC"),
		ReceiptMissing:       true,
		RequestedAmountCents: handlers.FmtInt64(2589),
	}
	newMovingExpenseDocParams := movedocop.CreateMovingExpenseDocumentParams{
		HTTPRequest:                        request,
		CreateMovingExpenseDocumentPayload: &newMovingExpenseDocPayload,
		MoveID:                             strfmt.UUID(move.ID.String()),
	}
	handlerConfig := suite.HandlerConfig()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := CreateMovingExpenseDocumentHandler{handlerConfig}

	response := handler.Handle(newMovingExpenseDocParams)

	suite.IsNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateMovingExpenseDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)
	suite.Equal(*createdPayload.Status, internalmessages.MoveDocumentStatusAWAITINGREVIEW)
}

func (suite *HandlerSuite) TestCreateMovingExpenseDocumentHandlerNoUploadsAndNotMissingReceipt() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	sm := move.Orders.ServiceMember
	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)
	moveDocumentType := internalmessages.MoveDocumentTypeOTHER
	movingExpenseType := internalmessages.MovingExpenseTypeWEIGHINGFEES
	newMovingExpenseDocPayload := internalmessages.CreateMovingExpenseDocumentPayload{
		MoveDocumentType:     &moveDocumentType,
		Title:                handlers.FmtString("awesome_document.pdf"),
		Notes:                handlers.FmtString("Some notes here"),
		MovingExpenseType:    &movingExpenseType,
		PaymentMethod:        handlers.FmtString("GTCC"),
		ReceiptMissing:       false,
		RequestedAmountCents: handlers.FmtInt64(2589),
	}
	newMovingExpenseDocParams := movedocop.CreateMovingExpenseDocumentParams{
		HTTPRequest:                        request,
		CreateMovingExpenseDocumentPayload: &newMovingExpenseDocPayload,
		MoveID:                             strfmt.UUID(move.ID.String()),
	}
	handlerConfig := suite.HandlerConfig()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := CreateMovingExpenseDocumentHandler{handlerConfig}

	response := handler.Handle(newMovingExpenseDocParams)

	// Submitting no uploads w/o selecting ReceiptMissing is an error
	suite.Assertions.IsType(&movedocop.CreateMovingExpenseDocumentBadRequest{}, response)
}

func (suite *HandlerSuite) TestCreateMovingExpenseDocumentHandlerStorageExpense() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	sm := move.Orders.ServiceMember
	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)
	moveDocumentType := internalmessages.MoveDocumentTypeOTHER
	movingExpenseType := internalmessages.MovingExpenseTypeSTORAGE
	newMovingExpenseDocPayload := internalmessages.CreateMovingExpenseDocumentPayload{
		MoveDocumentType:     &moveDocumentType,
		Title:                handlers.FmtString("awesome_document.pdf"),
		Notes:                handlers.FmtString("Some notes here"),
		MovingExpenseType:    &movingExpenseType,
		PaymentMethod:        handlers.FmtString("GTCC"),
		ReceiptMissing:       true,
		RequestedAmountCents: handlers.FmtInt64(200),
		StorageStartDate:     handlers.FmtDate(time.Date(2016, 01, 01, 0, 0, 0, 0, time.UTC)),
		StorageEndDate:       handlers.FmtDate(time.Date(2016, 01, 16, 0, 0, 0, 0, time.UTC)),
	}
	newMovingExpenseDocParams := movedocop.CreateMovingExpenseDocumentParams{
		HTTPRequest:                        request,
		CreateMovingExpenseDocumentPayload: &newMovingExpenseDocPayload,
		MoveID:                             strfmt.UUID(move.ID.String()),
	}
	handlerConfig := suite.HandlerConfig()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := CreateMovingExpenseDocumentHandler{handlerConfig}

	response := handler.Handle(newMovingExpenseDocParams)

	suite.Assertions.IsType(&movedocop.CreateMovingExpenseDocumentOK{}, response)
	createdResponse := response.(*movedocop.CreateMovingExpenseDocumentOK)
	movingExpense := models.MovingExpenseDocument{}
	err := suite.DB().Where("move_document_id = ?", createdResponse.Payload.ID).First(&movingExpense)
	suite.NoError(err)
	suite.Equal(movingExpense.StorageEndDate.UTC(), (time.Time)(*newMovingExpenseDocPayload.StorageEndDate))
	suite.Equal(movingExpense.StorageStartDate.UTC(), (time.Time)(*newMovingExpenseDocPayload.StorageStartDate))
}
