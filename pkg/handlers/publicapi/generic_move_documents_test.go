package publicapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	movedocop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) testCreateGenericMoveHandler(moveType models.SelectedMoveType) *models.MoveDocument {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	shipment := shipments[0]
	tspUser := tspUsers[0]

	upload := testdatagen.MakeUpload(suite.DB(), testdatagen.Assertions{
		Upload: models.Upload{
			UploaderID: *tspUser.UserID,
		},
	})
	upload.DocumentID = nil
	suite.MustSave(&upload)
	uploadIds := []strfmt.UUID{*handlers.FmtUUID(upload.ID)}

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateTspRequest(request, tspUser)

	newMoveDocPayload := apimessages.CreateGenericMoveDocumentPayload{
		UploadIds:        uploadIds,
		MoveDocumentType: apimessages.MoveDocumentTypeOTHER,
		Title:            handlers.FmtString("awesome_document.pdf"),
		Notes:            handlers.FmtString("Some notes here"),
	}

	newMoveDocParams := movedocop.CreateGenericMoveDocumentParams{
		HTTPRequest:                      request,
		CreateGenericMoveDocumentPayload: &newMoveDocPayload,
		ShipmentID:                       strfmt.UUID(shipment.ID.String()),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateGenericMoveDocumentHandler{context}
	response := handler.Handle(newMoveDocParams)

	// Assert we got back the 201 response
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
	wrongUser := testdatagen.MakeTspUser(suite.DB(), testdatagen.Assertions{
		TspUser: models.TspUser{
			Email: "unauthorized@example.com",
		},
		User: models.User{
			LoginGovEmail: "unauthorized@example.com",
		},
	})

	request = suite.AuthenticateTspRequest(request, wrongUser)
	newMoveDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(newMoveDocParams)
	suite.CheckResponseForbidden(badUserResponse)

	//Now try a bad shipment
	newMoveDocParams.ShipmentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newMoveDocParams)
	suite.CheckResponseForbidden(badMoveResponse)

	var moveDocument models.MoveDocument
	suite.Nil(suite.DB().Find(&moveDocument, createdPayload.ID), "could not load MoveDocument")
	return &moveDocument
}

func (suite *HandlerSuite) TestCreateGenericMoveDocumentHandlerHHG() {
	moveDocument := suite.testCreateGenericMoveHandler(models.SelectedMoveTypeHHG)
	suite.Nil(moveDocument.PersonallyProcuredMoveID, "moveDocument.PersonallyProcuredMoveID was not nil")
	suite.NotNil(moveDocument.ShipmentID, "moveDocument.ShipmentID was nil")
}

func (suite *HandlerSuite) TestCreateGenericMoveDocumentHandlerHHGPPM() {
	moveDocument := suite.testCreateGenericMoveHandler(models.SelectedMoveTypeHHGPPM)
	suite.Nil(moveDocument.PersonallyProcuredMoveID, "moveDocument.PersonallyProcuredMoveID was not nil")
	suite.NotNil(moveDocument.ShipmentID, "moveDocument.ShipmentID was nil")
}
