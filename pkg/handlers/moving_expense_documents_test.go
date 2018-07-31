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

func (suite *HandlerSuite) TestCreateMovingExpenseDocumentHandler() {

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

	method := internalmessages.MethodOfReceiptMILPAY
	newReimbursementPayload := internalmessages.CreateReimbursement{
		MethodOfReceipt: &method,
		RequestedAmount: fmtInt64(12000),
	}

	newMovingExpenseDocPayload := internalmessages.CreateMovingExpenseDocumentPayload{
		UploadIds:         uploadIds,
		MoveDocumentType:  internalmessages.MoveDocumentTypeOTHER,
		Title:             fmtString("awesome_document.pdf"),
		Notes:             fmtString("Some notes here"),
		MovingExpenseType: internalmessages.MovingExpenseTypeWEIGHINGFEES,
		Reimbursement:     &newReimbursementPayload,
	}

	newMovingExpenseDocParams := movedocop.CreateMovingExpenseDocumentParams{
		HTTPRequest:                        request,
		CreateMovingExpenseDocumentPayload: &newMovingExpenseDocPayload,
		MoveID: strfmt.UUID(move.ID.String()),
	}

	context := NewHandlerContext(suite.db, suite.logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateMovingExpenseDocumentHandler(context)
	response := handler.Handle(newMovingExpenseDocParams)
	// assert we got back the 201 response
	suite.isNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateMovingExpenseDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)

	// Make sure the Upload was associated to the new document
	createdDocumentID := createdPayload.Document.ID
	var fetchedUpload models.Upload
	suite.db.Find(&fetchedUpload, upload.ID)
	suite.Equal(createdDocumentID.String(), fetchedUpload.DocumentID.String())

	// Check that the status is correct
	suite.Equal(createdPayload.Status, internalmessages.MoveDocumentStatusAWAITINGREVIEW)

	// Check that the reimbursment has the right status
	suite.Equal(*createdPayload.Reimbursement.Status, internalmessages.ReimbursementStatusREQUESTED)

	// Next try the wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.db)
	request = suite.authenticateRequest(request, wrongUser)
	newMovingExpenseDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(newMovingExpenseDocParams)
	suite.checkResponseForbidden(badUserResponse)

	// Now try a bad move
	newMovingExpenseDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newMovingExpenseDocParams)
	suite.checkResponseNotFound(badMoveResponse)

}
