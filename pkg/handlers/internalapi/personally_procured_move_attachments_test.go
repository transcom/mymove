package internalapi

import (
	"net/http/httptest"
	"os"
	"regexp"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
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
	pdfConfig := pdfcpu.NewDefaultConfiguration()

	f, err := storer.FileSystem().Open(file.Name())
	suite.NoError(err)
	ctx, err := api.ReadContext(f, pdfConfig)
	suite.NoError(err)

	err = validate.XRefTable(ctx.XRefTable)
	suite.NoError(err)

	suite.Equal(count, ctx.PageCount)
}

func (suite *HandlerSuite) createHandlerContext() handlers.HandlerContext {
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)

	return context
}

func (suite *HandlerSuite) TestCreatePPMAttachmentsHandlerTests() {
	tests := []struct {
		name          string
		pdfName       string
		expectedPages int
	}{
		//We reuse the pdf twice for the merged pdf so expect new pdf to be 2x the length or original
		{name: "vanilla pdf", pdfName: "../../testdatagen/testdata/test.pdf", expectedPages: 2},
		//problem pdf is specific pdf that was causing decoding errors previously
		{name: "problem pdf", pdfName: "../../testdatagen/testdata/orders.pdf", expectedPages: 4},
	}
	uploadKeyRe := regexp.MustCompile(`(user/.+/uploads/.+)\?`)
	for _, test := range tests {
		suite.Run(test.name, func() {
			officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
			// Context gives us our file storer and filesystem
			context := suite.createHandlerContext()

			ppm := testdatagen.MakeDefaultPPM(suite.DB())
			expDoc := testdatagen.MakeMovingExpenseDocument(suite.DB(), testdatagen.Assertions{
				MoveDocument: models.MoveDocument{
					PersonallyProcuredMoveID: &ppm.ID,
					Status:                   models.MoveDocumentStatusOK,
				},
			})
			// Doc with an unapproved status
			testdatagen.MakeMovingExpenseDocument(suite.DB(), testdatagen.Assertions{
				MoveDocument: models.MoveDocument{
					PersonallyProcuredMoveID: &ppm.ID,
					Status:                   models.MoveDocumentStatusHASISSUE,
				},
			})
			// Open our test file
			f, err := os.Open(test.pdfName)
			suite.NoError(err)
			// Backfill the uploaded orders file in filesystem
			uploadedOrdersUpload := ppm.Move.Orders.UploadedOrders.UserUploads[0].Upload
			_, err = context.FileStorer().Store(uploadedOrdersUpload.StorageKey, f, uploadedOrdersUpload.Checksum, nil)
			suite.NoError(err)

			// Create upload for expense document model
			userUploader, err := uploader.NewUserUploader(suite.DB(), suite.TestLogger(), context.FileStorer(), 100*uploader.MB)
			suite.NoError(err)
			//RA Summary: gosec - errcheck - Unchecked return value
			//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
			//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
			//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
			//RA: in a unit test, then there is no risk
			//RA Developer Status: Mitigated
			//RA Validator Status: Mitigated
			//RA Modified Severity: N/A
			// nolint:errcheck
			userUploader.CreateUserUploadForDocument(&expDoc.MoveDocument.DocumentID, *officeUser.UserID, uploader.File{File: f}, uploader.AllowedTypesServiceMember)

			request := httptest.NewRequest("POST", "/fake/path", nil)
			request = suite.AuthenticateOfficeRequest(request, officeUser)

			docTypesToFetch := []string{"WEIGHT_TICKET", "EXPENSE", "OTHER", "STORAGE_EXPENSE"}
			params := ppmop.CreatePPMAttachmentsParams{
				PersonallyProcuredMoveID: *handlers.FmtUUID(ppm.ID),
				HTTPRequest:              request,
				DocTypes:                 docTypesToFetch,
			}

			handler := CreatePersonallyProcuredMoveAttachmentsHandler{context}
			response := handler.Handle(params)
			// assert we got back the 201 response
			suite.IsNotErrResponse(response)
			createdResponse := response.(*ppmop.CreatePPMAttachmentsOK)
			createdPDFPayload := createdResponse.Payload
			suite.NotNil(createdPDFPayload.URL)

			// Extract upload key from returned URL
			attachmentsURL := string(*createdPDFPayload.URL)
			uploadKey := uploadKeyRe.FindStringSubmatch(attachmentsURL)[1]

			merged, err := context.FileStorer().Fetch(uploadKey)
			suite.NoError(err)
			mergedFile := merged.(afero.File)

			suite.assertPDFPageCount(test.expectedPages, mergedFile, context.FileStorer())
		})
	}
}
