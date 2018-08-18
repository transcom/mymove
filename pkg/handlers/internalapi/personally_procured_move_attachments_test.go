package internalapi

import (
	"net/http/httptest"
	"os"
	"regexp"

	"github.com/pjdufour-truss/pdfcpu/pkg/api"
	"github.com/pjdufour-truss/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/afero"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) assertPDFPageCount(count int, file afero.File, storer storage.FileStorer) {
	pdfConfig := pdfcpu.NewInMemoryConfiguration()
	pdfConfig.FileSystem = storer.FileSystem()

	ctx, err := api.Read(file.Name(), pdfConfig)
	suite.NoError(err)

	err = pdfcpu.ValidateXRefTable(ctx.XRefTable)
	suite.NoError(err)

	suite.Equal(2, ctx.PageCount)
}

func (suite *HandlerSuite) createHandlerContext() HandlerContext {
	context := NewHandlerContext(suite.db, suite.logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)

	return context
}

func (suite *HandlerSuite) TestCreatePPMAttachmentsHandler() {
	uploadKeyRe := regexp.MustCompile(`(user/.+/uploads/.+)\?`)

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)
	ppm := testdatagen.MakeDefaultPPM(suite.db)
	expDoc := testdatagen.MakeMovingExpenseDocument(suite.db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			PersonallyProcuredMoveID: &ppm.ID,
			Status: models.MoveDocumentStatusOK,
		},
	})
	// Doc with an unapproved status
	testdatagen.MakeMovingExpenseDocument(suite.db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			PersonallyProcuredMoveID: &ppm.ID,
			Status: models.MoveDocumentStatusHASISSUE,
		},
	})

	// Context gives us our file storer and filesystem
	context := suite.createHandlerContext()

	// Open our test file
	f, err := os.Open("../fixtures/test.pdf")
	suite.NoError(err)

	// Backfill the uploaded orders file in filesystem
	uploadedOrdersUpload := ppm.Move.Orders.UploadedOrders.Uploads[0]
	_, err = context.storage.Store(uploadedOrdersUpload.StorageKey, f, uploadedOrdersUpload.Checksum)
	suite.NoError(err)

	// Create upload for expense document model
	loader := uploader.NewUploader(suite.db, suite.logger, context.storage)
	loader.CreateUpload(&expDoc.MoveDocument.DocumentID, *officeUser.UserID, f)

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateOfficeRequest(request, officeUser)

	params := ppmop.CreatePPMAttachmentsParams{
		PersonallyProcuredMoveID: *handlers.FmtUUID(ppm.ID),
		HTTPRequest:              request,
	}

	handler := CreatePersonallyProcuredMoveAttachmentsHandler(context)
	response := handler.Handle(params)
	// assert we got back the 201 response
	suite.isNotErrResponse(response)
	createdResponse := response.(*ppmop.CreatePPMAttachmentsOK)
	createdPDFPayload := createdResponse.Payload
	suite.NotNil(createdPDFPayload.URL)

	// Extract upload key from returned URL
	attachmentsURL := string(*createdPDFPayload.URL)
	uploadKey := uploadKeyRe.FindStringSubmatch(attachmentsURL)[1]

	merged, err := context.storage.Fetch(uploadKey)
	suite.NoError(err)
	mergedFile := merged.(afero.File)

	suite.assertPDFPageCount(2, mergedFile, context.storage)
}
