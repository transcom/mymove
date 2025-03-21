package ghcapi

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	uploadop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services/upload"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

const FixturePDF = "test.pdf"

func createPrereqs(suite *HandlerSuite, fixtureFile string) (models.Document, uploadop.CreateUploadParams) {
	document := factory.BuildDocument(suite.DB(), nil, nil)

	params := uploadop.NewCreateUploadParams()
	params.DocumentID = handlers.FmtUUID(document.ID)
	params.File = suite.Fixture(fixtureFile)

	return document, params
}

func makeRequest(suite *HandlerSuite, params uploadop.CreateUploadParams, serviceMember models.ServiceMember, fakeS3 *storageTest.FakeS3Storage) middleware.Responder {
	req := &http.Request{}
	req = suite.AuthenticateRequest(req, serviceMember)

	params.HTTPRequest = req

	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	handler := CreateUploadHandler{handlerConfig}
	response := handler.Handle(params)

	return response
}

func (suite *HandlerSuite) TestCreateUploadsHandlerSuccess() {
	t := suite.T()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	document, params := createPrereqs(suite, FixturePDF)

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

	expectedChecksum := "w7rJQqzlaazDW+mxTU9Q40Qchr3DW7FPQD7f8Js2J88="
	if upload.Checksum != expectedChecksum {
		t.Errorf("Did not calculate the correct MD5: expected %s, got %s", expectedChecksum, upload.Checksum)
	}
}

func (suite *HandlerSuite) TestCreateUploadsHandlerFailsWithWrongUser() {
	t := suite.T()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	_, params := createPrereqs(suite, FixturePDF)

	// Create a user that is not associated with the move
	otherUser := factory.BuildServiceMember(suite.DB(), nil, nil)

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
	document, params := createPrereqs(suite, FixturePDF)

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
	document, params := createPrereqs(suite, FixturePDF)

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
	document, params := createPrereqs(suite, FixturePDF)

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

func (suite *HandlerSuite) TestGetUploadStatusHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)
	localReceiver := notifications.StubNotificationReceiver{}

	orders := factory.BuildOrder(suite.DB(), nil, nil)
	uploadUser1 := factory.BuildUserUpload(suite.DB(), []factory.Customization{
		{
			Model:    orders.UploadedOrders,
			LinkOnly: true,
		},
		{
			Model: models.Upload{
				Filename:    "FileName",
				Bytes:       int64(15),
				ContentType: uploader.FileTypePDF,
			},
		},
	}, nil)

	file := suite.Fixture(FixturePDF)
	_, err := fakeS3.Store(uploadUser1.Upload.StorageKey, file.Data, "somehash", nil)
	suite.NoError(err)

	params := uploadop.NewGetUploadStatusParams()
	params.UploadID = strfmt.UUID(uploadUser1.Upload.ID.String())

	req := &http.Request{}
	req = suite.AuthenticateRequest(req, uploadUser1.Document.ServiceMember)
	params.HTTPRequest = req

	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	handlerConfig.SetNotificationReceiver(localReceiver)
	uploadInformationFetcher := upload.NewUploadInformationFetcher()
	handler := GetUploadStatusHandler{handlerConfig, uploadInformationFetcher}

	response := handler.Handle(params)
	_, ok := response.(*CustomGetUploadStatusResponse)
	suite.True(ok)

	queriedUpload := models.Upload{}
	err = suite.DB().Find(&queriedUpload, uploadUser1.Upload.ID)
	suite.NoError(err)
}

func (suite *HandlerSuite) TestGetUploadStatusHandlerFailure() {
	suite.Run("Error on no match for uploadId", func() {
		orders := factory.BuildOrder(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		uploadUUID := uuid.Must(uuid.NewV4())

		params := uploadop.NewGetUploadStatusParams()
		params.UploadID = strfmt.UUID(uploadUUID.String())

		req := &http.Request{}
		req = suite.AuthenticateRequest(req, orders.ServiceMember)
		params.HTTPRequest = req

		fakeS3 := storageTest.NewFakeS3Storage(true)
		localReceiver := notifications.StubNotificationReceiver{}

		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handlerConfig.SetNotificationReceiver(localReceiver)
		uploadInformationFetcher := upload.NewUploadInformationFetcher()
		handler := GetUploadStatusHandler{handlerConfig, uploadInformationFetcher}

		response := handler.Handle(params)
		_, ok := response.(*uploadop.GetUploadStatusNotFound)
		suite.True(ok)

		queriedUpload := models.Upload{}
		err := suite.DB().Find(&queriedUpload, uploadUUID)
		suite.Error(err)
	})

	suite.Run("Error when attempting access to another service member's upload", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		localReceiver := notifications.StubNotificationReceiver{}

		otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		orders := factory.BuildOrder(suite.DB(), nil, nil)
		uploadUser1 := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    orders.UploadedOrders,
				LinkOnly: true,
			},
			{
				Model: models.Upload{
					Filename:    "FileName",
					Bytes:       int64(15),
					ContentType: uploader.FileTypePDF,
				},
			},
		}, nil)

		file := suite.Fixture(FixturePDF)
		_, err := fakeS3.Store(uploadUser1.Upload.StorageKey, file.Data, "somehash", nil)
		suite.NoError(err)

		params := uploadop.NewGetUploadStatusParams()
		params.UploadID = strfmt.UUID(uploadUser1.Upload.ID.String())

		req := &http.Request{}
		req = suite.AuthenticateRequest(req, otherServiceMember)
		params.HTTPRequest = req

		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handlerConfig.SetNotificationReceiver(localReceiver)
		uploadInformationFetcher := upload.NewUploadInformationFetcher()
		handler := GetUploadStatusHandler{handlerConfig, uploadInformationFetcher}

		response := handler.Handle(params)
		_, ok := response.(*uploadop.GetUploadStatusForbidden)
		suite.True(ok)

		queriedUpload := models.Upload{}
		err = suite.DB().Find(&queriedUpload, uploadUser1.Upload.ID)
		suite.NoError(err)
	})
}
