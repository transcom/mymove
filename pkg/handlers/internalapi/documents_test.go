package internalapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestCreateDocumentsHandler() {
	t := suite.T()

	serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

	params := documentop.NewCreateDocumentParams()
	params.DocumentPayload = &internalmessages.PostDocumentPayload{
		ServiceMemberID: *handlers.FmtUUID(serviceMember.ID),
	}

	req := &http.Request{}
	req = suite.AuthenticateRequest(req, serviceMember)
	params.HTTPRequest = req

	handler := CreateDocumentHandler{suite.NewHandlerConfig()}
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

	if len(documentPayload.Uploads) != 0 {
		t.Errorf("wrong number of uploads, expected 0, got %d", len(documentPayload.Uploads))
	}

	document := models.Document{}
	err := suite.DB().Find(&document, documentPayload.ID)
	if err != nil {
		t.Errorf("Couldn't find expected document.")
	}
}

func (suite *HandlerSuite) TestShowDocumentHandler() {
	t := suite.T()

	userUpload := factory.BuildUserUpload(suite.DB(), nil, nil)

	documentID := userUpload.DocumentID
	var document models.Document

	err := suite.DB().Eager("ServiceMember.User").Find(&document, documentID)
	if err != nil {
		t.Fatalf("could not load document: %s", err)
	}

	params := documentop.NewShowDocumentParams()
	params.DocumentID = strfmt.UUID(documentID.String())

	req := &http.Request{}
	req = suite.AuthenticateRequest(req, document.ServiceMember)
	params.HTTPRequest = req

	handlerConfig := suite.NewHandlerConfig()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := ShowDocumentHandler{handlerConfig}
	response := handler.Handle(params)

	showResponse, ok := response.(*documentop.ShowDocumentOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
	documentPayload := showResponse.Payload

	// Double quote the filename to be able to handle filenames with commas in them
	quotedFilename := strconv.Quote(userUpload.Upload.Filename)

	responseDocumentUUID := documentPayload.ID.String()
	if responseDocumentUUID != documentID.String() {
		t.Errorf("wrong document uuid, expected %v, got %v", documentID, responseDocumentUUID)
	}

	if len(documentPayload.Uploads) != 1 {
		t.Errorf("wrong number of uploads, expected 1, got %d", len(documentPayload.Uploads))
	}

	uploadPayload := documentPayload.Uploads[0]
	values := url.Values{}
	values.Add("response-content-type", uploader.FileTypePDF)
	values.Add("response-content-disposition", "attachment; filename="+quotedFilename)
	values.Add("signed", "test")
	expectedURL := fmt.Sprintf("https://example.com/dir/%s?", userUpload.Upload.StorageKey) + values.Encode()
	if (uploadPayload.URL).String() != expectedURL {
		t.Errorf("wrong URL for upload, expected %s, got %s", expectedURL, uploadPayload.URL)
	}
}
