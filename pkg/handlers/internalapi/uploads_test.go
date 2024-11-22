// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
// RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
// RA: in which this would be considered a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package internalapi

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"golang.org/x/text/encoding/charmap"

	"github.com/transcom/mymove/pkg/factory"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services/upload"
	weightticketparser "github.com/transcom/mymove/pkg/services/weight_ticket_parser"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

const FixturePDF = "test.pdf"
const FixturePNG = "test.png"
const FixtureJPG = "test.jpg"
const FixtureTXT = "test.txt"
const FixtureXLS = "Weight Estimator.xls"
const FixtureXLSX = "Weight Estimator.xlsx"
const WeightEstimatorFullXLSX = "Weight Estimator Full.xlsx"
const WeightEstimatorPrefix = "Weight Estimator Full"
const FixtureEmpty = "empty.pdf"
const FixtureScreenshot = "Screenshot 2024-10-10 at 10.46.48 AM.png"

func createPrereqs(suite *HandlerSuite, fixtureFile string) (models.Document, uploadop.CreateUploadParams) {
	document := factory.BuildDocument(suite.DB(), nil, nil)

	params := uploadop.NewCreateUploadParams()
	params.DocumentID = handlers.FmtUUID(document.ID)
	params.File = suite.Fixture(fixtureFile)

	return document, params
}

func createPPMPrereqs(suite *HandlerSuite, fixtureFile string, weightReceipt bool) (models.Document, ppmop.CreatePPMUploadParams) {
	ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

	weightTicket := factory.BuildWeightTicket(suite.DB(), []factory.Customization{
		{
			Model:    ppmShipment,
			LinkOnly: true,
		},
	}, nil)

	params := ppmop.NewCreatePPMUploadParams()
	params.DocumentID = strfmt.UUID(weightTicket.EmptyDocumentID.String())
	params.PpmShipmentID = strfmt.UUID(ppmShipment.ID.String())
	params.File = suite.Fixture(fixtureFile)
	params.WeightReceipt = weightReceipt

	return weightTicket.EmptyDocument, params
}

func createPPMProgearPrereqs(suite *HandlerSuite, fixtureFile string) (models.Document, ppmop.CreatePPMUploadParams) {
	ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

	proGear := factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
		{
			Model:    ppmShipment,
			LinkOnly: true,
		},
	}, nil)

	params := ppmop.NewCreatePPMUploadParams()
	params.DocumentID = strfmt.UUID(proGear.DocumentID.String())
	params.PpmShipmentID = strfmt.UUID(ppmShipment.ID.String())
	params.File = suite.Fixture(fixtureFile)

	return proGear.Document, params
}

func createPPMExpensePrereqs(suite *HandlerSuite, fixtureFile string) (models.Document, ppmop.CreatePPMUploadParams) {
	ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

	movingExpense := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
		{
			Model:    ppmShipment,
			LinkOnly: true,
		},
	}, nil)

	params := ppmop.NewCreatePPMUploadParams()
	params.DocumentID = strfmt.UUID(movingExpense.DocumentID.String())
	params.PpmShipmentID = strfmt.UUID(ppmShipment.ID.String())
	params.File = suite.Fixture(fixtureFile)

	return movingExpense.Document, params
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

func makePPMRequest(suite *HandlerSuite, params ppmop.CreatePPMUploadParams, serviceMember models.ServiceMember, fakeS3 *storageTest.FakeS3Storage) middleware.Responder {
	req := &http.Request{}
	req = suite.AuthenticateRequest(req, serviceMember)

	params.HTTPRequest = req

	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	userUploader, err := uploader.NewUserUploader(handlerConfig.FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)
	suite.FatalNoError(err)

	pdfGenerator, err := paperworkgenerator.NewGenerator(userUploader.Uploader())
	suite.FatalNoError(err)

	parserComputer := weightticketparser.NewWeightTicketComputer()
	weightGenerator, err := weightticketparser.NewWeightTicketParserGenerator(pdfGenerator)
	suite.FatalNoError(err)

	handler := CreatePPMUploadHandler{handlerConfig, weightGenerator, parserComputer, userUploader}
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

func (suite *HandlerSuite) TestDeleteUploadHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)
	setupTestData := func() (models.UserUpload, models.Move) {
		move := factory.BuildMove(suite.DB(), nil, nil)
		uploadUser := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    move.Orders.UploadedOrders,
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

		return uploadUser, move
	}

	//when Move is in draft, upload can be deleted
	suite.Run("delete upload from DRAFT Move", func() {
		uploadUser, move := setupTestData()
		suite.Equal(move.Status, models.MoveStatusDRAFT)

		suite.Nil(uploadUser.Upload.DeletedAt)

		file := suite.Fixture(FixturePDF)
		fakeS3.Store(uploadUser.Upload.StorageKey, file.Data, "somehash", nil)

		params := uploadop.NewDeleteUploadParams()
		params.UploadID = strfmt.UUID(uploadUser.Upload.ID.String())

		req := &http.Request{}
		req = suite.AuthenticateRequest(req, uploadUser.Document.ServiceMember)
		params.HTTPRequest = req

		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		uploadInformationFetcher := upload.NewUploadInformationFetcher()
		fmt.Print(uploadInformationFetcher)
		handler := DeleteUploadHandler{handlerConfig, uploadInformationFetcher}
		response := handler.Handle(params)

		_, ok := response.(*uploadop.DeleteUploadNoContent)
		suite.True(ok)

		queriedUpload := models.Upload{}
		err := suite.DB().Find(&queriedUpload, uploadUser.Upload.ID)
		suite.Nil(err)
		suite.NotNil(queriedUpload.DeletedAt)
	})

	suite.Run("cannot delete upload once Move is out of DRAFT", func() {
		uploadUser, _ := setupTestData()
		move := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil)
		suite.Equal(move.Status, models.MoveStatusNeedsServiceCounseling)

		suite.Nil(uploadUser.Upload.DeletedAt)

		file := suite.Fixture(FixturePDF)
		fakeS3.Store(uploadUser.Upload.StorageKey, file.Data, "somehash", nil)

		params := uploadop.NewDeleteUploadParams()

		req := &http.Request{}
		req = suite.AuthenticateRequest(req, uploadUser.Document.ServiceMember)
		params.HTTPRequest = req

		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		uploadInformationFetcher := upload.NewUploadInformationFetcher()
		fmt.Print(uploadInformationFetcher)
		handler := DeleteUploadHandler{handlerConfig, uploadInformationFetcher}
		response := handler.Handle(params)

		_, ok := response.(*uploadop.DeleteUploadNoContent)
		suite.False(ok)

		queriedUpload := models.Upload{}
		err := suite.DB().Find(&queriedUpload, uploadUser.Upload.ID)
		suite.Nil(err)
		suite.Nil(queriedUpload.DeletedAt)
	})

}

func (suite *HandlerSuite) TestDeleteUploadsHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	uploadUser1 := factory.BuildUserUpload(suite.DB(), nil, nil)
	suite.Nil(uploadUser1.Upload.DeletedAt)

	uploadUser2Customization := []factory.Customization{
		{
			Model:    uploadUser1.Document,
			LinkOnly: true,
		},
	}
	uploadUser2 := factory.BuildUserUpload(suite.DB(), uploadUser2Customization, nil)

	file := suite.Fixture(FixturePDF)
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

	handlerConfig := suite.HandlerConfig()
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

func (suite *HandlerSuite) TestDeleteUploadHandlerSuccessEvenWithS3Failure() {
	// uploader.DeleteUpload only performs soft deletes and does not use the Storer's Delete method,
	// therefore a failure in the S3 storer will still result in a successful soft delete.
	fakeS3 := storageTest.NewFakeS3Storage(true)
	setupTestData := func() (models.UserUpload, models.Move) {
		move := factory.BuildMove(suite.DB(), nil, nil)
		uploadUser := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    move.Orders.UploadedOrders,
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

		return uploadUser, move
	}

	uploadUser, move := setupTestData()
	suite.Nil(uploadUser.Upload.DeletedAt)
	suite.Equal(move.Status, models.MoveStatusDRAFT)

	file := suite.Fixture(FixturePDF)
	fakeS3.Store(uploadUser.Upload.StorageKey, file.Data, "somehash", nil)

	params := uploadop.NewDeleteUploadParams()
	params.UploadID = strfmt.UUID(uploadUser.Upload.ID.String())

	req := &http.Request{}
	req = suite.AuthenticateRequest(req, uploadUser.Document.ServiceMember)
	params.HTTPRequest = req

	fakeS3Failure := storageTest.NewFakeS3Storage(false)

	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3Failure)
	uploadInformationFetcher := upload.NewUploadInformationFetcher()
	handler := DeleteUploadHandler{handlerConfig, uploadInformationFetcher}
	response := handler.Handle(params)

	_, ok := response.(*uploadop.DeleteUploadNoContent)
	suite.True(ok)

	queriedUpload := models.Upload{}
	err := suite.DB().Find(&queriedUpload, uploadUser.Upload.ID)
	suite.Nil(err)
	suite.NotNil(queriedUpload.DeletedAt)
}

func (suite *HandlerSuite) TestGetUploadStatusHandlerSuccess() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	move := factory.BuildMove(suite.DB(), nil, nil)
	uploadUser1 := factory.BuildUserUpload(suite.DB(), []factory.Customization{
		{
			Model:    move.Orders.UploadedOrders,
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
	fakeS3.Store(uploadUser1.Upload.StorageKey, file.Data, "somehash", nil)

	params := uploadop.NewGetUploadStatusParams()
	params.UploadID = strfmt.UUID(uploadUser1.ID.String())

	req := &http.Request{}
	req = suite.AuthenticateRequest(req, uploadUser1.Document.ServiceMember)
	params.HTTPRequest = req

	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	uploadInformationFetcher := upload.NewUploadInformationFetcher()
	handler := GetUploadStatusHandler{handlerConfig, uploadInformationFetcher}

	response := handler.Handle(params)

	res, ok := response.(*uploadop.GetUploadStatusOK)
	suite.True(ok)

	queriedUpload := models.Upload{}
	err := suite.DB().Find(&queriedUpload, uploadUser1.Upload.ID)
	suite.Nil(err)
	suite.Equal("CLEAN", res.Payload)
}

func (suite *HandlerSuite) TestCreatePPMUploadsHandlerSuccess() {
	suite.Run("uploads .xls file", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMPrereqs(suite, FixtureXLS, false)

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		suite.NoError(err)
		suite.Equal("V/Q6K9rVdEPVzgKbh5cn2x4Oci4XDaG4fcG04R41Iz4=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixtureXLS, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypeExcel, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(document.ServiceMember.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypeExcel))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+upload.Filename))
	})

	suite.Run("uploads .xlsx file", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMPrereqs(suite, FixtureXLSX, false)

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		suite.NoError(err)
		suite.Equal("eRZ1Cr3Ms0692k03ftoEdqXpvd/CHcbxmhEGEQBYVdY=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixtureXLSX, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypeExcelXLSX, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(document.ServiceMember.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypeExcelXLSX))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+upload.Filename))
	})

	suite.Run("uploads weight estimator .xlsx file (full weight)", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMPrereqs(suite, WeightEstimatorFullXLSX, true)

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		suite.NoError(err)

		// uploaded xlsx document should now be converted to a pdf so we check for pdf instead of xlsx
		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Contains(createdResponse.Payload.Filename, WeightEstimatorPrefix)
		suite.Equal(uploader.FileTypePDF, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(document.ServiceMember.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypePDF))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+upload.Filename))
	})

	suite.Run("uploads file for a progear document", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMProgearPrereqs(suite, FixturePNG)

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		suite.NoError(err)
		suite.Equal("/io1MRhLi2BFk9eF+lH1Ax+hyH+bPhlEK7A9/bqWlPY=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixturePNG, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypePNG, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(document.ServiceMember.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypePNG))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+upload.Filename))
	})

	suite.Run("uploads file for an expense document", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMExpensePrereqs(suite, FixtureJPG)

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		suite.NoError(err)
		suite.Equal("ibKT78j4CJecDXC6CbGISkqWFG5eSjCjlZJHlaFRho4=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixtureJPG, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypeJPEG, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(document.ServiceMember.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypeJPEG))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+upload.Filename))
	})

	suite.Run("uploads file with filename characters not supported by ISO8859_1", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMExpensePrereqs(suite, FixtureScreenshot)

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		filenameBuffer := make([]byte, 0)
		for _, r := range upload.Filename {
			if encodedRune, ok := charmap.ISO8859_1.EncodeRune(r); ok {
				filenameBuffer = append(filenameBuffer, encodedRune)
			}
		}

		suite.NoError(err)
		suite.Equal("/io1MRhLi2BFk9eF+lH1Ax+hyH+bPhlEK7A9/bqWlPY=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixtureScreenshot, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypePNG, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(document.ServiceMember.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypePNG))
		suite.NotContains(createdResponse.Payload.URL, url.QueryEscape(upload.Filename))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+string(filenameBuffer)))
	})
}

func (suite *HandlerSuite) TestCreatePPMUploadsHandlerFailure() {
	suite.Run("documentId does not exist", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMPrereqs(suite, FixtureXLS, false)

		params.DocumentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadNotFound{}, response)
		notFoundResponse, _ := response.(*ppmop.CreatePPMUploadNotFound)

		suite.Equal(fmt.Sprintf("documentId %q was not found for this user", params.DocumentID), *notFoundResponse.Payload.Detail)
	})

	suite.Run("documentId is not associated with the PPM shipment", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureXLS, false)

		document := factory.BuildDocument(suite.DB(), nil, nil)
		params.DocumentID = strfmt.UUID(document.ID.String())

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadNotFound{}, response)
		notFoundResponse, _ := response.(*ppmop.CreatePPMUploadNotFound)

		suite.Equal(fmt.Sprintf("documentId %q was not found for this shipment", params.DocumentID), *notFoundResponse.Payload.Detail)
	})

	suite.Run("service member session does not match document creator", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureXLS, false)

		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		response := makePPMRequest(suite, params, serviceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadNotFound{}, response)
		notFoundResponse, _ := response.(*ppmop.CreatePPMUploadNotFound)

		suite.Equal(fmt.Sprintf("documentId %q was not found for this user", params.DocumentID), *notFoundResponse.Payload.Detail)
	})

	suite.Run("ppmShipmentId does not exist", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMPrereqs(suite, FixtureXLS, false)

		params.PpmShipmentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadNotFound{}, response)
		notFoundResponse, _ := response.(*ppmop.CreatePPMUploadNotFound)

		suite.Equal(fmt.Sprintf("documentId %q was not found for this shipment", params.DocumentID), *notFoundResponse.Payload.Detail)
	})

	suite.Run("unsupported content type upload", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMPrereqs(suite, FixtureTXT, false)

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadUnprocessableEntity{}, response)
		invalidContentTypeResponse, _ := response.(*ppmop.CreatePPMUploadUnprocessableEntity)

		unsupportedErr := uploader.NewErrUnsupportedContentType(uploader.FileTypeTextUTF8, uploader.AllowedTypesPPMDocuments)
		suite.Equal(unsupportedErr.Error(), *invalidContentTypeResponse.Payload.Detail)
	})

	suite.Run("empty file upload", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		document, params := createPPMPrereqs(suite, FixtureEmpty, false)

		response := makePPMRequest(suite, params, document.ServiceMember, fakeS3)

		suite.CheckResponseBadRequest(response)

		badResponseErr := response.(*handlers.ErrResponse)
		suite.Equal("File has length of 0", badResponseErr.Err.Error())
	})
}
