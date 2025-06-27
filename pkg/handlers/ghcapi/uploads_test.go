package ghcapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"golang.org/x/text/encoding/charmap"

	"github.com/transcom/mymove/pkg/factory"
	ppmop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	uploadop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
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
const WeightEstimatorXlsxFail = "Weight Estimator Expect Failed Upload.xlsx"
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

	handlerConfig := suite.NewHandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)

	handler := CreateUploadHandler{handlerConfig}
	response := handler.Handle(params)

	return response
}

func makePPMRequest(suite *HandlerSuite, params ppmop.CreatePPMUploadParams, officeUser models.OfficeUser, fakeS3 *storageTest.FakeS3Storage) middleware.Responder {
	req := &http.Request{}
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params.HTTPRequest = req

	handlerConfig := suite.NewHandlerConfig()
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

func (suite *HandlerSuite) TestCreatePPMUploadsHandlerSuccess() {
	suite.Run("uploads .xls file", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureXLS, false)

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		// Double quote the filename to be able to handle filenames with commas in them
		quotedFilename := strconv.Quote(upload.Filename)

		suite.NoError(err)
		suite.Equal("V/Q6K9rVdEPVzgKbh5cn2x4Oci4XDaG4fcG04R41Iz4=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixtureXLS, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypeExcel, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(officeUser.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypeExcel))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+quotedFilename))
	})

	suite.Run("uploads .xlsx file", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureXLSX, false)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		// Double quote the filename to be able to handle filenames with commas in them
		quotedFilename := strconv.Quote(upload.Filename)

		suite.NoError(err)
		suite.Equal("eRZ1Cr3Ms0692k03ftoEdqXpvd/CHcbxmhEGEQBYVdY=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixtureXLSX, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypeExcelXLSX, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(officeUser.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypeExcelXLSX))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+quotedFilename))
	})

	suite.Run("uploads weight estimator .xlsx file (full weight)", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, WeightEstimatorFullXLSX, true)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		suite.NoError(err)

		// Double quote the filename to be able to handle filenames with commas in them
		quotedFilename := strconv.Quote(upload.Filename)

		// uploaded xlsx document should now be converted to a pdf so we check for pdf instead of xlsx
		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Contains(createdResponse.Payload.Filename, WeightEstimatorPrefix)
		suite.Equal(uploader.FileTypePDF, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(officeUser.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypePDF))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+quotedFilename))
	})

	suite.Run("uploads file for a progear document", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMProgearPrereqs(suite, FixturePNG)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		// Double quote the filename to be able to handle filenames with commas in them
		quotedFilename := strconv.Quote(upload.Filename)

		suite.NoError(err)
		suite.Equal("/io1MRhLi2BFk9eF+lH1Ax+hyH+bPhlEK7A9/bqWlPY=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixturePNG, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypePNG, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(officeUser.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypePNG))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+quotedFilename))
	})

	suite.Run("uploads file for an expense document", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMExpensePrereqs(suite, FixtureJPG)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		// Double quote the filename to be able to handle filenames with commas in them
		quotedFilename := strconv.Quote(upload.Filename)

		suite.NoError(err)
		suite.Equal("ibKT78j4CJecDXC6CbGISkqWFG5eSjCjlZJHlaFRho4=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixtureJPG, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypeJPEG, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(officeUser.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypeJPEG))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+quotedFilename))
	})

	suite.Run("uploads file with filename characters not supported by ISO8859_1", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMExpensePrereqs(suite, FixtureScreenshot)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadCreated{}, response)

		createdResponse, _ := response.(*ppmop.CreatePPMUploadCreated)

		upload := models.Upload{}
		err := suite.DB().Find(&upload, createdResponse.Payload.ID)

		// Double quote the filename to be able to handle filenames with commas in them
		quotedFilename := strconv.Quote(upload.Filename)

		filenameBuffer := make([]byte, 0)
		for _, r := range quotedFilename {
			if encodedRune, ok := charmap.ISO8859_1.EncodeRune(r); ok {
				filenameBuffer = append(filenameBuffer, encodedRune)
			}
		}

		suite.NoError(err)
		suite.Equal("/io1MRhLi2BFk9eF+lH1Ax+hyH+bPhlEK7A9/bqWlPY=", upload.Checksum)

		suite.NotEmpty(createdResponse.Payload.ID)
		suite.Equal(FixtureScreenshot, createdResponse.Payload.Filename)
		suite.Equal(uploader.FileTypePNG, createdResponse.Payload.ContentType)
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(officeUser.UserID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(upload.ID.String()))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape(uploader.FileTypePNG))
		suite.NotContains(createdResponse.Payload.URL, url.QueryEscape(upload.Filename))
		suite.Contains(createdResponse.Payload.URL, url.QueryEscape("attachment; filename="+string(filenameBuffer)))
	})
}

func (suite *HandlerSuite) TestCreatePPMUploadsHandlerFailure() {
	suite.Run("documentId does not exist", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureXLS, false)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		params.DocumentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadNotFound{}, response)
		notFoundResponse, _ := response.(*ppmop.CreatePPMUploadNotFound)

		suite.Equal(fmt.Sprintf("documentId %q was not found", params.DocumentID), *notFoundResponse.Payload.Message)
	})

	suite.Run("documentId is not associated with the PPM shipment", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureXLS, false)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		document := factory.BuildDocument(suite.DB(), nil, nil)
		params.DocumentID = strfmt.UUID(document.ID.String())

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadNotFound{}, response)
		notFoundResponse, _ := response.(*ppmop.CreatePPMUploadNotFound)

		suite.Equal(fmt.Sprintf("documentId %q was not found for this shipment", params.DocumentID), *notFoundResponse.Payload.Message)
	})

	suite.Run("ppmShipmentId does not exist", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureXLS, false)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		params.PpmShipmentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadNotFound{}, response)
		notFoundResponse, _ := response.(*ppmop.CreatePPMUploadNotFound)

		suite.Equal(fmt.Sprintf("documentId %q was not found for this shipment", params.DocumentID), *notFoundResponse.Payload.Message)
	})

	suite.Run("unsupported content type upload", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureTXT, false)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadUnprocessableEntity{}, response)
		invalidContentTypeResponse, _ := response.(*ppmop.CreatePPMUploadUnprocessableEntity)

		unsupportedErr := uploader.NewErrUnsupportedContentType(uploader.FileTypeTextUTF8, uploader.AllowedTypesPPMDocuments)
		suite.Equal(unsupportedErr.Error(), *invalidContentTypeResponse.Payload.Detail)
	})

	suite.Run("empty file upload", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, FixtureEmpty, false)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.CheckResponseBadRequest(response)

		badResponseErr := response.(*handlers.ErrResponse)
		suite.Equal("File has length of 0", badResponseErr.Err.Error())
	})

	suite.Run("Non-weight Estimator FIle Submitted for upload", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		_, params := createPPMPrereqs(suite, WeightEstimatorXlsxFail, true)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		response := makePPMRequest(suite, params, officeUser, fakeS3)

		suite.IsType(&ppmop.CreatePPMUploadForbidden{}, response)
		incorrectXlsxResponse, _ := response.(*ppmop.CreatePPMUploadForbidden)
		suite.Equal("The uploaded .xlsx file does not match the expected weight estimator file format.", *incorrectXlsxResponse.Payload.Message)
	})
}
