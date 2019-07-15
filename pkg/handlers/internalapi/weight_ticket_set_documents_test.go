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

func (suite *HandlerSuite) TestCreateWeightTicketSetDocumentHandler() {

	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	sm := ppm.Move.Orders.ServiceMember

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

	newWeightTicketSetDocumentPayload := internalmessages.CreateWeightTicketDocumentsPayload{
		UploadIds:                uploadIds,
		EmptyWeight:              handlers.FmtInt64(1000),
		FullWeight:               handlers.FmtInt64(2000),
		EmptyWeightTicketMissing: handlers.FmtBool(false),
		FullWeightTicketMissing:  handlers.FmtBool(false),
		PersonallyProcuredMoveID: handlers.FmtUUID(ppm.ID),
		VehicleNickname:          handlers.FmtString("My car"),
		VehicleOptions:           handlers.FmtString("CAR"),
		WeightTicketDate:         handlers.FmtDate(testdatagen.NextValidMoveDate),
		TrailerOwnershipMissing:  handlers.FmtBool(false),
	}

	newWeightTicketSetDocParams := movedocop.CreateWeightTicketDocumentParams{
		HTTPRequest:                request,
		CreateWeightTicketDocument: &newWeightTicketSetDocumentPayload,
		MoveID:                     strfmt.UUID(ppm.MoveID.String()),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateWeightTicketSetDocumentHandler{context}
	response := handler.Handle(newWeightTicketSetDocParams)
	suite.IsNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateWeightTicketDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)

	// confirm Upload is associated to the new document
	createdDocumentID := createdPayload.Document.ID
	var fetchedUpload models.Upload
	suite.DB().Find(&fetchedUpload, upload.ID)
	suite.Equal(createdDocumentID.String(), fetchedUpload.DocumentID.String())
	suite.Equal(createdPayload.Status, internalmessages.MoveDocumentStatusAWAITINGREVIEW)

	var fetchedMoveDocument models.MoveDocument
	err := suite.DB().Q().Where("move_id = ?", ppm.MoveID).First(&fetchedMoveDocument)
	suite.NoError(err)
	var fetchedWeightTicket models.WeightTicketSetDocument
	err = suite.DB().Q().Where("move_document_id = ?", fetchedMoveDocument.ID).First(&fetchedWeightTicket)
	suite.NoError(err)

	// Wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	request = suite.AuthenticateRequest(request, wrongUser)
	newWeightTicketSetDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(newWeightTicketSetDocParams)
	suite.CheckResponseForbidden(badUserResponse)

	// Bad move
	newWeightTicketSetDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newWeightTicketSetDocParams)
	suite.CheckResponseNotFound(badMoveResponse)

}
