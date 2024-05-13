package paperwork

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *PaperworkServiceSuite) TestFormGotenbergURL() {
	suite.Run("protocol, host, and port are based on env vars", func() {
		testCases := []struct {
			protocol string
			host     string
			port     string
			expected string
		}{
			{"http", "localhost", "3000", "http://localhost:3000/"},
			{"https", "my.move.mil", "2000", "https://my.move.mil:2000/"},
		}

		for _, testCase := range testCases {
			testCase := testCase

			suite.Run(fmt.Sprintf("protocol: %s, host: %s, port: %s", testCase.protocol, testCase.host, testCase.port), func() {
				os.Setenv(GotenbergProtocol, testCase.protocol)
				os.Setenv(GotenbergHost, testCase.host)
				os.Setenv(GotenbergPort, testCase.port)

				actual := formGotenbergURL("")

				suite.Equal(testCase.expected, actual)
			})
		}
	})

	suite.Run("can append an endpoint to the URL", func() {
		os.Setenv(GotenbergProtocol, "http")
		os.Setenv(GotenbergHost, "localhost")
		os.Setenv(GotenbergPort, "2000")

		endpoints := []string{
			"forms/libreoffice/convert",
			"forms/pdfengines/merge",
		}

		for _, endpoint := range endpoints {
			endpoint := endpoint

			suite.Run(fmt.Sprintf("endpoint: %s", endpoint), func() {
				actual := formGotenbergURL(endpoint)

				suite.True(strings.HasSuffix(actual, endpoint), fmt.Sprintf("expected %s to end with %s", actual, endpoint))
			})
		}
	})
}

func (suite *PaperworkServiceSuite) TestConvertFileToPDF() {
	suite.Run("Returns an error if there is an issue with the stream", func() {
		stream := factory.FixtureOpen("test.png")

		userUpload := factory.BuildUserUpload(nil, nil, nil)
		userUpload.ID = uuid.Must(uuid.NewV4())

		fileInfo := services.NewFileInfo(&userUpload, stream)

		// closing the file should make the conversion fail because it'll fail to read the file.
		closeErr := fileInfo.OriginalUploadStream.Close()

		suite.FatalNoError(closeErr)

		err := convertFileToPDF(suite.AppContextForTest(), fileInfo)

		if suite.Nil(fileInfo.PDFStream) && suite.Error(err) {
			suite.Contains(err.Error(), fmt.Sprintf("failed to convert file %s (UserUpload ID: %d) to PDF", fileInfo.UserUpload.Upload.Filename, fileInfo.UserUpload.ID))

			suite.Contains(err.Error(), "file already closed")
		}
	})

	suite.Run("Returns an error if there is an issue creating the request", func() {
		// We already have a protocol set in `.envrc`, but for the sake of this test, we want to blank it out to get an
		// error with creating a request.
		if originalProtocol := os.Getenv(GotenbergProtocol); originalProtocol != "" {
			os.Unsetenv(GotenbergProtocol)

			// We want to restore the original value after the test ends.
			defer func() {
				os.Setenv(GotenbergProtocol, originalProtocol)
			}()
		}

		stream := factory.FixtureOpen("test.png")

		defer stream.Close()

		userUpload := factory.BuildUserUpload(nil, nil, nil)
		userUpload.ID = uuid.Must(uuid.NewV4())

		fileInfo := services.NewFileInfo(&userUpload, stream)
		err := convertFileToPDF(suite.AppContextForTest(), fileInfo)

		if suite.Nil(fileInfo.PDFStream) && suite.Error(err) {
			suite.Contains(err.Error(), fmt.Sprintf("failed to convert file %s (UserUpload ID: %d) to PDF", fileInfo.UserUpload.Upload.Filename, fileInfo.UserUpload.ID))

			suite.Contains(err.Error(), "missing protocol scheme")
		}
	})

	suite.Run("Returns an error if there is an error with making the http request", func() {
		// We should already have these set because of `.envrc`, but just in case, we'll want these set for this test.
		os.Setenv(GotenbergProtocol, "http")
		os.Setenv(GotenbergHost, "localhost")
		os.Setenv(GotenbergPort, "2000")

		stream := factory.FixtureOpen("test.png")

		defer stream.Close()

		userUpload := factory.BuildUserUpload(nil, nil, nil)
		userUpload.ID = uuid.Must(uuid.NewV4())

		fileInfo := services.NewFileInfo(&userUpload, stream)

		err := convertFileToPDF(suite.AppContextForTest(), fileInfo)

		if suite.Nil(fileInfo.PDFStream) && suite.Error(err) {
			suite.Contains(err.Error(), fmt.Sprintf("failed to convert file %s (UserUpload ID: %d) to PDF", fileInfo.UserUpload.Upload.Filename, fileInfo.UserUpload.ID))

			suite.Contains(err.Error(), "connection refused")
		}
	})

	suite.Run(fmt.Sprintf("Returns an error if the response code isn't %d", http.StatusOK), func() {
		mockGotenbergServer := suite.setUpMockGotenbergServer(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		defer mockGotenbergServer.Close()

		stream := factory.FixtureOpen("test.png")

		defer stream.Close()

		userUpload := factory.BuildUserUpload(nil, nil, nil)
		userUpload.ID = uuid.Must(uuid.NewV4())

		fileInfo := services.NewFileInfo(&userUpload, stream)

		err := convertFileToPDF(suite.AppContextForTest(), fileInfo)

		if suite.Nil(fileInfo.PDFStream) && suite.Error(err) {
			suite.Contains(err.Error(), fmt.Sprintf("failed to convert file %s (UserUpload ID: %d) to PDF", fileInfo.UserUpload.Upload.Filename, fileInfo.UserUpload.ID))

			suite.Contains(err.Error(), fmt.Sprintf("bad status | code: %d", http.StatusNotFound))
		}
	})

	suite.Run("Set PDF stream on file info struct if the conversion is successful", func() {
		expectedPDF := factory.FixtureOpen("test.pdf")

		defer expectedPDF.Close()

		expectedBytes, readOriginalErr := io.ReadAll(expectedPDF)

		suite.FatalNoError(readOriginalErr)

		_, seekErr := expectedPDF.Seek(0, io.SeekStart)

		suite.FatalNoError(seekErr)

		mockGotenbergServer := suite.setUpMockGotenbergServer(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := io.Copy(w, expectedPDF)

			suite.FatalNoError(err)
		})

		defer mockGotenbergServer.Close()

		stream := factory.FixtureOpen("test.png")

		defer stream.Close()

		userUpload := factory.BuildUserUpload(nil, nil, nil)
		userUpload.ID = uuid.Must(uuid.NewV4())

		fileInfo := services.NewFileInfo(&userUpload, stream)

		err := convertFileToPDF(suite.AppContextForTest(), fileInfo)

		if suite.NotNil(fileInfo.PDFStream) && suite.NoError(err) {
			actualBytes, readConvertedErr := io.ReadAll(fileInfo.PDFStream)

			suite.NoError(readConvertedErr)

			suite.Equal(expectedBytes, actualBytes)
		}
	})
}
