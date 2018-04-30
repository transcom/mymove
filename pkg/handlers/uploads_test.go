package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	authcontext "github.com/transcom/mymove/pkg/auth"
	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type putFile struct {
	key      string
	body     io.ReadSeeker
	checksum string
}

type fakeS3Storage struct {
	putFiles    []putFile
	willSucceed bool
}

func (fake *fakeS3Storage) Key(args ...string) string {
	return path.Join(args...)
}

func (fake *fakeS3Storage) Store(key string, data io.ReadSeeker, md5 string) (*storage.StoreResult, error) {
	file := putFile{
		key:      key,
		body:     data,
		checksum: md5,
	}
	fake.putFiles = append(fake.putFiles, file)
	buf := []byte{}
	_, err := data.Read(buf)
	if err != nil {
		return nil, err
	}
	if fake.willSucceed {
		return &storage.StoreResult{}, nil
	}
	return nil, errors.New("failed to push")
}

func (fake *fakeS3Storage) PresignedURL(key string, contentType string) (string, error) {
	url := fmt.Sprintf("https://example.com/dir/%s?contentType=%%2F%s", key, contentType)
	return url, nil
}

func newFakeS3Storage(willSucceed bool) *fakeS3Storage {
	return &fakeS3Storage{
		willSucceed: willSucceed,
	}
}

func createPrereqs(suite *HandlerSuite) (models.Document, uploadop.CreateUploadParams) {
	t := suite.T()

	document, err := testdatagen.MakeDocument(suite.db, nil)
	if err != nil {
		t.Fatalf("could not create document: %s", err)
	}

	params := uploadop.NewCreateUploadParams()
	params.DocumentID = strfmt.UUID(document.ID.String())
	params.File = suite.fixture("test.pdf")

	return document, params
}

func makeRequest(suite *HandlerSuite, params uploadop.CreateUploadParams, userID uuid.UUID, fakeS3 *fakeS3Storage) middleware.Responder {
	ctx := authcontext.PopulateAuthContext(context.Background(), userID, "fake token")
	params.HTTPRequest = (&http.Request{}).WithContext(ctx)

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	handler := CreateUploadHandler(context)
	response := handler.Handle(params)

	return response
}

func (suite *HandlerSuite) TestCreateUploadsHandlerSuccess() {
	t := suite.T()
	fakeS3 := newFakeS3Storage(true)
	document, params := createPrereqs(suite)

	response := makeRequest(suite, params, document.UploaderID, fakeS3)
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

	if len(fakeS3.putFiles) != 1 {
		t.Errorf("Wrong number of putFiles: expected 1, got %d", len(fakeS3.putFiles))
	}

	key := fmt.Sprintf("/documents/%s/uploads/%s", document.ID, upload.ID)
	if fakeS3.putFiles[0].key != key {
		t.Errorf("Wrong key name: expected %s, got %s", key, fakeS3.putFiles[0].key)
	}

	pos, err := fakeS3.putFiles[0].body.Seek(0, io.SeekCurrent)
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
	fakeS3 := newFakeS3Storage(true)
	_, params := createPrereqs(suite)

	// Create a user that is not associated with the move
	otherUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&otherUser)

	response := makeRequest(suite, params, otherUser.ID, fakeS3)
	_, ok := response.(*uploadop.CreateUploadForbidden)
	if !ok {
		t.Fatalf("Wrong response type. Expected CreateUploadForbidden, got %T", response)
	}

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

	fakeS3 := newFakeS3Storage(true)
	document, params := createPrereqs(suite)

	// Make a document ID that is not actually associated with a document
	params.DocumentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())

	response := makeRequest(suite, params, document.UploaderID, fakeS3)
	_, ok := response.(*uploadop.CreateUploadNotFound)
	if !ok {
		t.Fatalf("Wrong response type. Expected CreateUploadNotFound, got %T", response)
	}

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

	fakeS3 := newFakeS3Storage(true)
	document, params := createPrereqs(suite)

	params.File = suite.fixture("empty.pdf")

	response := makeRequest(suite, params, document.UploaderID, fakeS3)
	_, ok := response.(*uploadop.CreateUploadBadRequest)
	if !ok {
		t.Fatalf("Wrong response type. Expected CreateUploadNotFound, got %T", response)
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
	fakeS3 := newFakeS3Storage(false)
	document, params := createPrereqs(suite)

	response := makeRequest(suite, params, document.UploaderID, fakeS3)
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
