package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateGenericMoveDocumentHandler() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	sm := move.Orders.ServiceMember

	userUpload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
		UserUpload: models.UserUpload{
			UploaderID: sm.UserID,
		},
	})
	userUpload.DocumentID = nil
	suite.MustSave(&userUpload)
	uploadIds := []strfmt.UUID{*handlers.FmtUUID(userUpload.UploadID)}

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, sm)

	newMoveDocPayload := internalmessages.CreateGenericMoveDocumentPayload{
		UploadIds:        uploadIds,
		MoveDocumentType: internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeOTHER),
		Title:            handlers.FmtString("awesome_document.pdf"),
		Notes:            handlers.FmtString("Some notes here"),
	}

	newMoveDocParams := movedocop.CreateGenericMoveDocumentParams{
		HTTPRequest:                      request,
		CreateGenericMoveDocumentPayload: &newMoveDocPayload,
		MoveID:                           strfmt.UUID(move.ID.String()),
	}

	handlerConfig := suite.HandlerConfig()
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
