package handlers

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMoveDocumentHandler() {
	move := testdatagen.MakeDefaultMove(suite.db)
	sm := move.Orders.ServiceMember

	document := testdatagen.MakeDocument(suite.db, testdatagen.Assertions{
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
			Name:            "Move document test",
		},
	})
	docID := strfmt.UUID(document.ID.String())

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	newMoveDocPayload := internalmessages.CreateMoveDocumentPayload{
		DocumentID:       &docID,
		MoveDocumentType: internalmessages.MoveDocumentTypeOTHER,
		Notes:            fmtString("Some notes here"),
		Status:           internalmessages.MoveDocumentStatusAWAITINGREVIEW,
	}

	newMoveDocParams := moveop.CreateMoveDocumentParams{
		HTTPRequest:               request,
		CreateMoveDocumentPayload: &newMoveDocPayload,
		MoveID: strfmt.UUID(move.ID.String()),
	}

	handler := CreateMoveDocumentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(newMoveDocParams)
	// assert we got back the 201 response
	createdResponse := response.(*moveop.CreateMoveDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)

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
	move1 := testdatagen.MakeDefaultMove(suite.db)
	sm := move1.Orders.ServiceMember

	moveDocument := testdatagen.MakeMoveDocument(suite.db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID: move1.ID,
			Move:   move1,
		},
	})

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	indexMoveDocParams := moveop.IndexMoveDocumentsParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move1.ID.String()),
	}

	handler := IndexMoveDocumentsHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(indexMoveDocParams)

	// assert we got back the 201 response
	indexResponse := response.(*moveop.IndexMoveDocumentsOK)
	indexPayload := indexResponse.Payload
	suite.NotNil(indexPayload)

	for _, moveDoc := range indexPayload {
		suite.Require().Equal(*moveDoc.ID, strfmt.UUID(moveDocument.ID.String()), "expected move ids to match")
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
