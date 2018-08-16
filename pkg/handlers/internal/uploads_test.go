package internal

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func createPrereqs(suite *utils.HandlerSuite) (models.Document, uploadop.CreateUploadParams) {
	document := testdatagen.MakeDefaultDocument(suite.db)

	params := uploadop.NewCreateUploadParams()
	params.DocumentID = utils.FmtUUID(document.ID)
	params.File = suite.fixture("test.pdf")

	return document, params
}

func makeRequest(suite *utils.HandlerSuite, params uploadop.CreateUploadParams, serviceMember models.ServiceMember, fakeS3 *storageTest.FakeS3Storage) middleware.Responder {
	req := &http.Request{}
	req = suite.authenticateRequest(req, serviceMember)

	params.HTTPRequest = req

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	handler := CreateUploadHandler(context)
	response := handler.Handle(params)

	return response
}

func (suite *utils.HandlerSuite) TestCreateUploadsHandlerSuccess() {
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
}

func (suite *utils.HandlerSuite) TestCreateUploadsHandlerFailsWithWrongUser() {
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

func (suite *utils.HandlerSuite) TestCreateUploadsHandlerFailsWithMissingDoc() {
	t := suite.T()

	fakeS3 := storageTest.NewFakeS3Storage(true)
	document, params := createPrereqs(suite)

	// Make a document ID that is not actually associated with a document
	params.DocumentID = utils.FmtUUID(uuid.Must(uuid.NewV4()))

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

func (suite *utils.HandlerSuite) TestCreateUploadsHandlerFailsWithZeroLengthFile() {
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

func (suite *utils.HandlerSuite) TestCreateUploadsHandlerFailure() {
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

func (suite *utils.HandlerSuite) TestDeleteUploadHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	upload := testdatagen.MakeDefaultUpload(suite.db)

	file := suite.fixture("test.pdf")
	fakeS3.Store(upload.StorageKey, file.Data, "somehash")

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

func (suite *utils.HandlerSuite) TestDeleteUploadsHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	upload1 := testdatagen.MakeDefaultUpload(suite.db)

	upload2Assertions := testdatagen.Assertions{
		Upload: models.Upload{
			Document:   upload1.Document,
			DocumentID: &upload1.Document.ID,
		},
	}
	upload2 := testdatagen.MakeUpload(suite.db, upload2Assertions)

	file := suite.fixture("test.pdf")
	fakeS3.Store(upload1.StorageKey, file.Data, "somehash")
	fakeS3.Store(upload2.StorageKey, file.Data, "somehash")

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
