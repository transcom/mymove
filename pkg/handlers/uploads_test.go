package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func createPrereqs(suite *HandlerSuite) (models.Document, uploadop.CreateUploadParams) {
	document := testdatagen.MakeDefaultDocument(suite.db)

	params := uploadop.NewCreateUploadParams()
	params.DocumentID = strfmt.UUID(document.ID.String())
	params.File = suite.fixture("test.pdf")

	return document, params
}

func makeRequest(suite *HandlerSuite, params uploadop.CreateUploadParams, serviceMember models.ServiceMember, fakeS3 *storageTest.FakeS3Storage) middleware.Responder {
	req := &http.Request{}
	req = suite.authenticateRequest(req, serviceMember)

	params.HTTPRequest = req

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	handler := CreateUploadHandler(context)
	response := handler.Handle(params)

	return response
}

func (suite *HandlerSuite) TestCreateUploadsHandlerSuccess() {
	t := suite.T()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	document, params := createPrereqs(suite)

	response := makeRequest(suite, params, document.ServiceMember, fakeS3)
	createdResponse, ok := response.(*uploadop.CreateUploadCreated)
	if !ok {
		t.Fatalf("Wrong response type. Expected CreateUploadCreated, got %T", response)
	}

	uploadPayload := createdResponse.Payload
	upload := models.Upload{}
	err := suite.db.Find(&upload, uploadPayload.ID)
	if err != nil {
		t.Fatalf("Couldn't find expected upload.")
	}

	expectedChecksum := "nOE6HwzyE4VEDXn67ULeeA=="
	if upload.Checksum != expectedChecksum {
		t.Errorf("Did not calculate the correct MD5: expected %s, got %s", expectedChecksum, upload.Checksum)
	}

	if len(fakeS3.PutFiles) != 1 {
		t.Errorf("Wrong number of putFiles: expected 1, got %d", len(fakeS3.PutFiles))
	}

	key := fmt.Sprintf("documents/%s/uploads/%s", document.ID, upload.ID)
	if fakeS3.PutFiles[0].Key != key {
		t.Errorf("Wrong key name: expected %s, got %s", key, fakeS3.PutFiles[0].Key)
	}

	pos, err := fakeS3.PutFiles[0].Body.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Fatalf("Could't check position in uploaded file: %s", err)
	}

	if pos != 0 {
		t.Errorf("Wrong file position: expected 0, got %d", pos)
	}

	// TODO verify Body
}

func (suite *HandlerSuite) TestCreateUploadsHandlerFailsWithWrongUser() {
	t := suite.T()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	_, params := createPrereqs(suite)

	// Create a user that is not associated with the move
	otherUser := testdatagen.MakeDefaultServiceMember(suite.db)

	response := makeRequest(suite, params, otherUser, fakeS3)
	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.code)

	count, err := suite.db.Count(&models.Upload{})

	if err != nil {
		t.Fatalf("Couldn't count uploads in database: %s", err)
	}

	if count != 0 {
		t.Fatalf("Wrong number of uploads in database: expected 0, got %d", count)
	}
}

func (suite *HandlerSuite) TestCreateUploadsHandlerFailsWithMissingDoc() {
	t := suite.T()

	fakeS3 := storageTest.NewFakeS3Storage(true)
	document, params := createPrereqs(suite)

	// Make a document ID that is not actually associated with a document
	params.DocumentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())

	response := makeRequest(suite, params, document.ServiceMember, fakeS3)
	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.code)

	count, err := suite.db.Count(&models.Upload{})

	if err != nil {
		t.Fatalf("Couldn't count uploads in database: %s", err)
	}

	if count != 0 {
		t.Fatalf("Wrong number of uploads in database: expected 0, got %d", count)
	}
}

func (suite *HandlerSuite) TestCreateUploadsHandlerFailsWithZeroLengthFile() {
	t := suite.T()

	fakeS3 := storageTest.NewFakeS3Storage(true)
	document, params := createPrereqs(suite)

	params.File = suite.fixture("empty.pdf")

	response := makeRequest(suite, params, document.ServiceMember, fakeS3)
	_, ok := response.(*uploadop.CreateUploadBadRequest)
	if !ok {
		t.Fatalf("Wrong response type. Expected CreateUploadBadRequest, got %T", response)
	}

	count, err := suite.db.Count(&models.Upload{})

	if err != nil {
		t.Fatalf("Couldn't count uploads in database: %s", err)
	}

	if count != 0 {
		t.Fatalf("Wrong number of uploads in database: expected 0, got %d", count)
	}
}

func (suite *HandlerSuite) TestCreateUploadsHandlerFailure() {
	t := suite.T()
	fakeS3 := storageTest.NewFakeS3Storage(false)
	document, params := createPrereqs(suite)

	response := makeRequest(suite, params, document.ServiceMember, fakeS3)
	_, ok := response.(*uploadop.CreateUploadInternalServerError)
	if !ok {
		t.Fatalf("Wrong response type. Expected CreateUploadInternalServerError, got %T", response)
	}

	count, err := suite.db.Count(&models.Upload{})

	if err != nil {
		t.Fatalf("Couldn't count uploads in database: %s", err)
	}

	if count != 0 {
		t.Fatalf("Wrong number of uploads in database: expected 0, got %d", count)
	}
}

func (suite *HandlerSuite) TestDeleteUploadHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	upload := testdatagen.MakeDefaultUpload(suite.db)

	file := suite.fixture("test.pdf")
	key := fakeS3.Key("documents", upload.DocumentID.String(), "uploads", upload.ID.String())
	fakeS3.Store(key, file.Data, "somehash")

	params := uploadop.NewDeleteUploadParams()
	params.UploadID = strfmt.UUID(upload.ID.String())

	req := &http.Request{}
	req = suite.authenticateRequest(req, upload.Document.ServiceMember)
	params.HTTPRequest = req

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	handler := DeleteUploadHandler(context)
	response := handler.Handle(params)

	_, ok := response.(*uploadop.DeleteUploadNoContent)
	suite.True(ok)

	queriedUpload := models.Upload{}
	err := suite.db.Find(&queriedUpload, upload.ID)
	suite.NotNil(err)
}

func (suite *HandlerSuite) TestDeleteUploadsHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	upload1 := testdatagen.MakeDefaultUpload(suite.db)

	upload2Assertions := testdatagen.Assertions{
		Upload: models.Upload{
			Document:   upload1.Document,
			DocumentID: upload1.Document.ID,
		},
	}
	upload2 := testdatagen.MakeUpload(suite.db, upload2Assertions)

	file := suite.fixture("test.pdf")
	key1 := fakeS3.Key("documents", upload1.DocumentID.String(), "uploads", upload1.ID.String())
	key2 := fakeS3.Key("documents", upload2.DocumentID.String(), "uploads", upload2.ID.String())
	fakeS3.Store(key1, file.Data, "somehash")
	fakeS3.Store(key2, file.Data, "somehash")

	params := uploadop.NewDeleteUploadsParams()
	params.UploadIds = []strfmt.UUID{
		strfmt.UUID(upload1.ID.String()),
		strfmt.UUID(upload2.ID.String()),
	}

	req := &http.Request{}
	req = suite.authenticateRequest(req, upload1.Document.ServiceMember)
	params.HTTPRequest = req

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	handler := DeleteUploadsHandler(context)
	response := handler.Handle(params)

	_, ok := response.(*uploadop.DeleteUploadsNoContent)
	suite.True(ok)

	queriedUpload := models.Upload{}
	err := suite.db.Find(&queriedUpload, upload1.ID)
	suite.NotNil(err)
}
