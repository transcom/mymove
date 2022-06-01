//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package internalapi

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func createPrereqs(suite *HandlerSuite) (models.Document, uploadop.CreateUploadParams) {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	params := uploadop.NewCreateUploadParams()
	params.DocumentID = handlers.FmtUUID(document.ID)
	params.File = suite.Fixture("test.pdf")

	return document, params
}

func makeRequest(suite *HandlerSuite, params uploadop.CreateUploadParams, serviceMember models.ServiceMember, fakeS3 *storageTest.FakeS3Storage) middleware.Responder {
	req := &http.Request{}
	req = suite.AuthenticateRequest(req, serviceMember)

	params.HTTPRequest = req

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handlerConfig.SetFileStorer(fakeS3)
	handler := CreateUploadHandler{handlerConfig}
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
	err := suite.DB().Find(&upload, uploadPayload.ID)
	if err != nil {
		t.Fatalf("Couldn't find expected upload.")
	}

	expectedChecksum := "nOE6HwzyE4VEDXn67ULeeA=="
	if upload.Checksum != expectedChecksum {
		t.Errorf("Did not calculate the correct MD5: expected %s, got %s", expectedChecksum, upload.Checksum)
	}
}

func (suite *HandlerSuite) TestCreateUploadsHandlerFailsWithWrongUser() {
	t := suite.T()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	_, params := createPrereqs(suite)

	// Create a user that is not associated with the move
	otherUser := testdatagen.MakeDefaultServiceMember(suite.DB())

	response := makeRequest(suite, params, otherUser, fakeS3)
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.Code)

	count, err := suite.DB().Count(&models.Upload{})

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
	params.DocumentID = handlers.FmtUUID(uuid.Must(uuid.NewV4()))

	response := makeRequest(suite, params, document.ServiceMember, fakeS3)
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.Code)

	count, err := suite.DB().Count(&models.Upload{})

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

	params.File = suite.Fixture("empty.pdf")

	response := makeRequest(suite, params, document.ServiceMember, fakeS3)
	suite.CheckResponseBadRequest(response)

	count, err := suite.DB().Count(&models.Upload{})

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

	currentCount, countErr := suite.DB().Count(&models.Upload{})
	suite.Nil(countErr)

	response := makeRequest(suite, params, document.ServiceMember, fakeS3)
	suite.CheckResponseInternalServerError(response)

	count, err := suite.DB().Count(&models.Upload{})
	if err != nil {
		t.Fatalf("Couldn't count uploads in database: %s", err)
	}

	if count != currentCount {
		t.Fatalf("Wrong number of uploads in database: expected %d, got %d", currentCount, count)
	}
}

func (suite *HandlerSuite) TestDeleteUploadHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	uploadUser := testdatagen.MakeDefaultUserUpload(suite.DB())
	suite.Nil(uploadUser.Upload.DeletedAt)

	file := suite.Fixture("test.pdf")
	fakeS3.Store(uploadUser.Upload.StorageKey, file.Data, "somehash", nil)

	params := uploadop.NewDeleteUploadParams()
	params.UploadID = strfmt.UUID(uploadUser.Upload.ID.String())

	req := &http.Request{}
	req = suite.AuthenticateRequest(req, uploadUser.Document.ServiceMember)
	params.HTTPRequest = req

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handlerConfig.SetFileStorer(fakeS3)
	handler := DeleteUploadHandler{handlerConfig}
	response := handler.Handle(params)

	_, ok := response.(*uploadop.DeleteUploadNoContent)
	suite.True(ok)

	queriedUpload := models.Upload{}
	err := suite.DB().Find(&queriedUpload, uploadUser.Upload.ID)
	suite.Nil(err)
	suite.NotNil(queriedUpload.DeletedAt)
}

func (suite *HandlerSuite) TestDeleteUploadsHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	uploadUser1 := testdatagen.MakeDefaultUserUpload(suite.DB())
	suite.Nil(uploadUser1.Upload.DeletedAt)

	uploadUser2Assertions := testdatagen.Assertions{
		UserUpload: models.UserUpload{
			Document:   uploadUser1.Document,
			DocumentID: &uploadUser1.Document.ID,
		},
	}
	uploadUser2 := testdatagen.MakeUserUpload(suite.DB(), uploadUser2Assertions)

	file := suite.Fixture("test.pdf")
	fakeS3.Store(uploadUser1.Upload.StorageKey, file.Data, "somehash", nil)
	fakeS3.Store(uploadUser2.Upload.StorageKey, file.Data, "somehash", nil)

	params := uploadop.NewDeleteUploadsParams()
	params.UploadIds = []strfmt.UUID{
		strfmt.UUID(uploadUser1.Upload.ID.String()),
		strfmt.UUID(uploadUser2.Upload.ID.String()),
	}

	req := &http.Request{}
	req = suite.AuthenticateRequest(req, uploadUser1.Document.ServiceMember)
	params.HTTPRequest = req

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handlerConfig.SetFileStorer(fakeS3)
	handler := DeleteUploadsHandler{handlerConfig}
	response := handler.Handle(params)

	_, ok := response.(*uploadop.DeleteUploadsNoContent)
	suite.True(ok)

	queriedUpload := models.Upload{}
	err := suite.DB().Find(&queriedUpload, uploadUser1.Upload.ID)
	suite.Nil(err)
	suite.NotNil(queriedUpload.DeletedAt)
}

func (suite *HandlerSuite) TestDeleteUploadHandlerFailure() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	uploadUser := testdatagen.MakeDefaultUserUpload(suite.DB())
	suite.Nil(uploadUser.Upload.DeletedAt)

	file := suite.Fixture("test.pdf")
	fakeS3.Store(uploadUser.Upload.StorageKey, file.Data, "somehash", nil)

	params := uploadop.NewDeleteUploadParams()
	params.UploadID = strfmt.UUID(uploadUser.Upload.ID.String())

	req := &http.Request{}
	req = suite.AuthenticateRequest(req, uploadUser.Document.ServiceMember)
	params.HTTPRequest = req

	fakeS3Failure := storageTest.NewFakeS3Storage(false)

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handlerConfig.SetFileStorer(fakeS3Failure)
	handler := DeleteUploadHandler{handlerConfig}
	response := handler.Handle(params)

	_, ok := response.(*uploadop.DeleteUploadNoContent)
	suite.False(ok)

	queriedUpload := models.Upload{}
	err := suite.DB().Find(&queriedUpload, uploadUser.Upload.ID)
	suite.Nil(err)
	suite.Nil(queriedUpload.DeletedAt)
}
