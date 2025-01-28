package paperwork

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaperworkServiceSuite) TestUserUploadToPDFConverter() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)

	suite.FatalNoError(uploaderErr)

	uploadToPDFConverter := NewUserUploadToPDFConverter(userUploader)

	suite.Run("Returns an error if there is an issue with downloading an upload", func() {
		failingFakeS3 := storageTest.NewFakeS3Storage(false)

		failingUserUploader, failingUploaderErr := uploader.NewUserUploader(failingFakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(failingUploaderErr)

		failingUploadToPDFConverter := NewUserUploadToPDFConverter(failingUserUploader)

		appCtx := suite.AppContextForTest()

		badUserUpload := factory.BuildUserUpload(nil, []factory.Customization{
			{
				Model: models.UserUpload{
					ID: uuid.Must(uuid.NewV4()),
				},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: failingUserUploader,
					AppContext:   appCtx,
				},
			},
		}, nil)

		convertedFiles, err := failingUploadToPDFConverter.ConvertUserUploadsToPDF(appCtx, models.UserUploads{badUserUpload})

		if suite.Nil(convertedFiles) && suite.Error(err) {
			suite.Contains(err.Error(), fmt.Sprintf("failed to convert file %s (UserUpload ID: %d) to PDF", badUserUpload.Upload.Filename, badUserUpload.ID))

			suite.Contains(err.Error(), "failed to fetch file")
		}
	})

	suite.Run("Returns an upload stream as is if it is already a PDF", func() {
		expectedPDF := factory.FixtureOpen("test.pdf")

		defer expectedPDF.Close()

		expectedBytes, readExpectedErr := io.ReadAll(expectedPDF)

		suite.FatalNoError(readExpectedErr)

		_, seekErr := expectedPDF.Seek(0, io.SeekStart)

		suite.FatalNoError(seekErr)

		appCtx := suite.AppContextForTest()

		userUpload := factory.BuildUserUpload(appCtx.DB(), []factory.Customization{
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   appCtx,
					File:         expectedPDF,
				},
			},
		}, nil)

		convertedFiles, err := uploadToPDFConverter.ConvertUserUploadsToPDF(appCtx, models.UserUploads{userUpload})

		if suite.NoError(err) && suite.Len(convertedFiles, 1) {
			// The way this is tested, by reading the stream, also serves to let us know that the stream is open,
			// which is what we want since the caller will want to have access to it.
			actualBytes, readConvertedErr := io.ReadAll(convertedFiles[0].PDFStream)

			suite.NoError(readConvertedErr)

			suite.Equal(expectedBytes, actualBytes)

			// We also want to make sure that the original stream is closed, since we don't want to leave
			// the file open and we won't be using it since we have a PDF stream.
			originalStreamBytes, readOriginalStreamErr := io.ReadAll(convertedFiles[0].OriginalUploadStream)

			suite.Equal([]byte{}, originalStreamBytes)
			suite.ErrorContains(readOriginalStreamErr, "File is closed")
		}
	})

	suite.Run("Returns an error if one of the files fails to convert", func() {
		mockGotenbergServer := suite.setUpMockGotenbergServer(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		defer mockGotenbergServer.Close()

		appCtx := suite.AppContextForTest()

		userUpload1 := factory.BuildUserUpload(appCtx.DB(), []factory.Customization{
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   appCtx,
					File:         factory.FixtureOpen("test.pdf"), // PDF so we'll skip the gotenberg call
				},
			},
		}, nil)

		userUpload2 := factory.BuildUserUpload(appCtx.DB(), []factory.Customization{
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   appCtx,
					File:         factory.FixtureOpen("test.png"), // PNG so we'll hit the gotenberg call
				},
			},
		}, nil)

		convertedFiles, err := uploadToPDFConverter.ConvertUserUploadsToPDF(appCtx, models.UserUploads{userUpload1, userUpload2})

		if suite.Nil(convertedFiles) && suite.Error(err) {
			suite.Contains(err.Error(), fmt.Sprintf("failed to convert file %s (UserUpload ID: %d) to PDF", userUpload2.Upload.Filename, userUpload2.ID))

			// This is just to ensure we're bubbling up the error from the inner convert function and is tied to the
			// mock Gotenberg server returning a 404.
			suite.Contains(err.Error(), "404 Not Found")
		}
	})

	suite.Run("Can successfully convert multiple user uploads to PDF", func() {
		// We'll set up different PDFs to be returned by our mock gotenberg server to make it easier to test that
		// we are getting multiple files converted and passed back.
		expectedPDF1 := factory.FixtureOpen("empty-weight-ticket.pdf")

		defer expectedPDF1.Close()

		expectedBytes1, readExpectedErr1 := io.ReadAll(expectedPDF1)

		suite.FatalNoError(readExpectedErr1)

		_, seekErr1 := expectedPDF1.Seek(0, io.SeekStart)

		suite.FatalNoError(seekErr1)

		expectedPDF2 := factory.FixtureOpen("full-weight-ticket.pdf")

		defer expectedPDF2.Close()

		expectedBytes2, readExpectedErr2 := io.ReadAll(expectedPDF2)

		suite.FatalNoError(readExpectedErr2)

		_, seekErr2 := expectedPDF2.Seek(0, io.SeekStart)

		suite.FatalNoError(seekErr2)

		expectedFiles := []struct {
			pdf   io.Reader
			bytes []byte
		}{
			{expectedPDF1, expectedBytes1},
			{expectedPDF2, expectedBytes2},
		}

		timesGotenbergServerCalled := 0

		mockGotenbergServer := suite.setUpMockGotenbergServer(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := io.Copy(w, expectedFiles[timesGotenbergServerCalled].pdf)

			suite.FatalNoError(err)

			timesGotenbergServerCalled++
		})

		defer mockGotenbergServer.Close()

		appCtx := suite.AppContextForTest()

		// The actual files don't matter beyond them not being PDF files, so that we know the gotenberg code was called.
		userUpload1 := factory.BuildUserUpload(appCtx.DB(), []factory.Customization{
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   appCtx,
					File:         factory.FixtureOpen("empty-weight-ticket.png"),
				},
			},
		}, nil)

		userUpload2 := factory.BuildUserUpload(appCtx.DB(), []factory.Customization{
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   appCtx,
					File:         factory.FixtureOpen("full-weight-ticket.png"),
				},
			},
		}, nil)

		convertedFiles, err := uploadToPDFConverter.ConvertUserUploadsToPDF(appCtx, models.UserUploads{userUpload1, userUpload2})

		if suite.NoError(err) && suite.Len(convertedFiles, len(expectedFiles)) {
			for i, convertedFile := range convertedFiles {
				if suite.NotNil(convertedFile.PDFStream) {
					// The way this is tested, by reading the stream, also serves to let us know that the stream is
					// open, which is what we want since the caller will want to have access to it.
					actualBytes, readConvertedErr := io.ReadAll(convertedFile.PDFStream)

					suite.NoError(readConvertedErr)

					suite.Equal(expectedFiles[i].bytes, actualBytes)

					// We also want to make sure that the original stream is closed, since we don't want to leave
					// the file open and we won't be using it since we have a PDF stream.
					originalStreamBytes, readOriginalStreamErr := io.ReadAll(convertedFile.OriginalUploadStream)

					suite.Equal([]byte{}, originalStreamBytes)
					suite.ErrorContains(readOriginalStreamErr, "File is closed")
				}
			}
		}
	})
}
