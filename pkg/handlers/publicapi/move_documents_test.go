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

func (suite *HandlerSuite) TestIndexMoveDocumentsHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	shipment := shipments[0]
	tspUser := tspUsers[0]
	move := shipment.Move

	moveDocument := testdatagen.MakeMoveDocument(suite.DB(), testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:     move.ID,
			Move:       move,
			ShipmentID: &shipment.ID,
		},
	})

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateTspRequest(request, tspUser)

	indexMoveDocParams := movedocop.IndexMoveDocumentsParams{
		HTTPRequest: request,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
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
		suite.Require().Equal(*moveDoc.ShipmentID, strfmt.UUID(shipment.ID.String()), "expected shipment ids to match")
	}

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
	indexMoveDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(indexMoveDocParams)
	suite.CheckResponseForbidden(badUserResponse)

	// Now try a bad shipment
	indexMoveDocParams.ShipmentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(indexMoveDocParams)
	suite.CheckResponseForbidden(badMoveResponse)
}

func (suite *HandlerSuite) TestUpdateMoveDocumentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	shipment := shipments[0]
	tspUser := tspUsers[0]
	move := shipment.Move
	sm := shipment.Move.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.DB(), testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:     move.ID,
			Move:       move,
			ShipmentID: &shipment.ID,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	})
	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateTspRequest(request, tspUser)

	// And: the title and status are updated
	updateMoveDocPayload := apimessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
		Title:            handlers.FmtString("super_awesome.pdf"),
		Notes:            handlers.FmtString("This document is super awesome."),
		Status:           apimessages.MoveDocumentStatusOK,
		MoveDocumentType: apimessages.MoveDocumentTypeOTHER,
	}

	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		HTTPRequest:        request,
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
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
	updateMoveDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(updateMoveDocParams)
	suite.CheckResponseForbidden(badUserResponse)

	// Now try a bad shipment
	updateMoveDocParams.ShipmentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(updateMoveDocParams)
	suite.CheckResponseForbidden(badMoveResponse)
}
