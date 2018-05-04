package handlers

import (
	"context"
	"net/http"

	"github.com/satori/go.uuid"

	authcontext "github.com/transcom/mymove/pkg/auth"
	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateDocumentsHandler() {
	t := suite.T()

	serviceMember, err := testdatagen.MakeServiceMember(suite.db)
	if err != nil {
		t.Fatalf("could not create serviceMember: %s", err)
	}

	userID := serviceMember.UserID

	params := documentop.NewCreateDocumentParams()
	params.DocumentPayload = &internalmessages.PostDocumentPayload{Name: "test document", ServiceMemberID: *fmtUUID(serviceMember.ID)}

	ctx := authcontext.PopulateAuthContext(context.Background(), userID, "fake token")
	params.HTTPRequest = (&http.Request{}).WithContext(ctx)

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

	if uuid.Must(uuid.FromString(documentPayload.ServiceMemberID.String())) == uuid.Nil {
		t.Errorf("got empty serviceMember uuid")
	}

	if documentPayload.Name == nil {
		t.Errorf("got nil document name")
	} else if *documentPayload.Name != "test document" {
		t.Errorf("wrong document name, expected %s, got %s", "test document", *documentPayload.Name)
	}

	if len(documentPayload.Uploads) != 0 {
		t.Errorf("wrong number of uploads, expected 0, got %d", len(documentPayload.Uploads))
	}

	document := models.Document{}
	err = suite.db.Find(&document, documentPayload.ID)
	if err != nil {
		t.Errorf("Couldn't find expected document.")
	}
}
