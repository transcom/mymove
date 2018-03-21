package handlers

import (
	"github.com/satori/go.uuid"

	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	// "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateDocumentsHandler() {
	t := suite.T()

	move, err := testdatagen.MakeMove(suite.db)
	if err != nil {
		t.Fatalf("could not create move: %s", err)
	}

	params := documentop.NewCreateDocumentParams()
	handler := CreateDocumentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	createdResponse, ok := response.(*documentop.CreateDocumentCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
	documentPayload := createdResponse.Payload

	if uuid.Must(uuid.FromString(documentPayload.ID.String())) == uuid.Nil {
		t.Errorf("got empty document uuid")
	}

	if documentPayload.Name == nil {
		t.Errorf("got nil document name")
	} else if *documentPayload.Name != "test document" {
		t.Errorf("wrong document name, expected %s, got %s")
	}

	if len(documentPayload.Uploads) != 0 {
		t.Errorf("wrong number of uploads, expected 0, got %d", len(documentPayload.Uploads))
	}
}
