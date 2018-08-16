package internal

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

func (suite *HandlerSuite) TestCreateMovingExpenseDocumentHandler() {

	move := testdatagen.MakeDefaultMove(suite.parent.Db)
	sm := move.Orders.ServiceMember

	upload := testdatagen.MakeUpload(suite.parent.Db, testdatagen.Assertions{
		Upload: models.Upload{
			UploaderID: sm.UserID,
		},
	})
	upload.DocumentID = nil
	suite.parent.MustSave(&upload)
	uploadIds := []strfmt.UUID{*utils.FmtUUID(upload.ID)}

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.parent.AuthenticateRequest(request, sm)

	newMovingExpenseDocPayload := internalmessages.CreateMovingExpenseDocumentPayload{
		UploadIds:            uploadIds,
		MoveDocumentType:     internalmessages.MoveDocumentTypeOTHER,
		Title:                utils.FmtString("awesome_document.pdf"),
		Notes:                utils.FmtString("Some notes here"),
		MovingExpenseType:    internalmessages.MovingExpenseTypeWEIGHINGFEES,
		PaymentMethod:        utils.FmtString("GTCC"),
		RequestedAmountCents: utils.FmtInt64(2589),
	}

	newMovingExpenseDocParams := movedocop.CreateMovingExpenseDocumentParams{
		HTTPRequest:                        request,
		CreateMovingExpenseDocumentPayload: &newMovingExpenseDocPayload,
		MoveID: strfmt.UUID(move.ID.String()),
	}

	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateMovingExpenseDocumentHandler(context)
	response := handler.Handle(newMovingExpenseDocParams)
	// assert we got back the 201 response
	suite.parent.IsNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateMovingExpenseDocumentOK)
	createdPayload := createdResponse.Payload
	suite.parent.NotNil(createdPayload.ID)

	// Make sure the Upload was associated to the new document
	createdDocumentID := createdPayload.Document.ID
	var fetchedUpload models.Upload
	suite.parent.Db.Find(&fetchedUpload, upload.ID)
	suite.parent.Equal(createdDocumentID.String(), fetchedUpload.DocumentID.String())

	// Check that the status is correct
	suite.parent.Equal(createdPayload.Status, internalmessages.MoveDocumentStatusAWAITINGREVIEW)

	// Next try the wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.parent.Db)
	request = suite.parent.AuthenticateRequest(request, wrongUser)
	newMovingExpenseDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(newMovingExpenseDocParams)
	suite.parent.CheckResponseForbidden(badUserResponse)

	// Now try a bad move
	newMovingExpenseDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newMovingExpenseDocParams)
	suite.parent.CheckResponseNotFound(badMoveResponse)

}
