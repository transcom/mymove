package paperwork

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *PaperworkServiceSuite) TestPDFMerger() {
	merger := NewPDFMerger()

	suite.Run("Returns an error if there is an issue reading one of the PDF streams", func() {
		pdf := factory.FixtureOpen("test.pdf")

		// closing the file should lead to an error because we can't read it to send it as part of a request
		suite.FatalNoError(pdf.Close())

		mergedPDF, err := merger.MergePDFs(suite.AppContextForTest(), []io.ReadCloser{pdf})

		if suite.Nil(mergedPDF) && suite.Error(err) {
			suite.ErrorContains(err, "failed to copy PDF stream 0 to request")

			suite.ErrorContains(err, "file already closed")
		}
	})

	suite.Run("Returns an error if there is an error creating the request", func() {
		// We already have a protocol set in `.envrc`, but for the sake of this test, we want to blank it out to get an
		// error with creating a request.
		if originalProtocol := os.Getenv(GotenbergProtocol); originalProtocol != "" {
			os.Unsetenv(GotenbergProtocol)

			// We want to restore the original value after the test ends.
			defer func() {
				os.Setenv(GotenbergProtocol, originalProtocol)
			}()
		}

		pdf := factory.FixtureOpen("test.pdf")

		defer pdf.Close()

		mergedPDF, err := merger.MergePDFs(suite.AppContextForTest(), []io.ReadCloser{pdf})

		if suite.Nil(mergedPDF) && suite.Error(err) {
			suite.ErrorContains(err, "failed to create request to merge PDFs")

			suite.ErrorContains(err, "missing protocol scheme")
		}
	})

	suite.Run("Returns an error if there is an error with making the http request", func() {
		// We should already have these set because of `.envrc`, but just in case, we'll want these set for this test.
		os.Setenv(GotenbergProtocol, "http")
		os.Setenv(GotenbergHost, "localhost")
		os.Setenv(GotenbergPort, "2000")

		pdf := factory.FixtureOpen("test.pdf")

		defer pdf.Close()

		mergedPDF, err := merger.MergePDFs(suite.AppContextForTest(), []io.ReadCloser{pdf})

		if suite.Nil(mergedPDF) && suite.Error(err) {
			suite.ErrorContains(err, "failed to merge PDFs")

			suite.ErrorContains(err, "connection refused")
		}
	})

	suite.Run(fmt.Sprintf("Returns an error if the response code isn't %d", http.StatusOK), func() {
		mockGotenbergServer := suite.setUpMockGotenbergServer(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		defer mockGotenbergServer.Close()

		pdf := factory.FixtureOpen("test.pdf")

		defer pdf.Close()

		mergedPDF, err := merger.MergePDFs(suite.AppContextForTest(), []io.ReadCloser{pdf})

		if suite.Nil(mergedPDF) && suite.Error(err) {
			suite.ErrorContains(err, "failed to merge PDFs")

			suite.ErrorContains(err, fmt.Sprintf("bad status | code: %d", http.StatusNotFound))
		}
	})

	suite.Run("Returns a merged PDF if there are no errors", func() {
		expectedPDF := factory.FixtureOpen("test.pdf")

		expectedBytes, readExpectedErr := io.ReadAll(expectedPDF)

		suite.FatalNoError(readExpectedErr)

		_, seekErr := expectedPDF.Seek(0, 0)

		suite.FatalNoError(seekErr)

		mockGotenbergServer := suite.setUpMockGotenbergServer(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := io.Copy(w, expectedPDF)

			suite.FatalNoError(err)
		})

		defer mockGotenbergServer.Close()

		pdfsToMerge := []io.ReadCloser{
			factory.FixtureOpen("empty-weight-ticket.pdf"),
			factory.FixtureOpen("full-weight-ticket.pdf"),
		}

		mergedPDF, err := merger.MergePDFs(suite.AppContextForTest(), pdfsToMerge)

		if suite.NoError(err) && suite.NotNil(mergedPDF) {
			mergedBytes, readMergedErr := io.ReadAll(mergedPDF)

			suite.NoError(readMergedErr)

			// We want to make sure that the merged PDF is the same as the expected PDF.
			suite.Equal(expectedBytes, mergedBytes)
		}
	})
}
