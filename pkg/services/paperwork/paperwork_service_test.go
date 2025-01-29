package paperwork

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/uploader"
)

type PaperworkServiceSuite struct {
	*testingsuite.PopTestSuite
	userUploader *uploader.UserUploader
	filesToClose []afero.File
}

func TestPaperworkServiceSuite(t *testing.T) {

	storer := storageTest.NewFakeS3Storage(true)

	newUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		log.Panic(err)
	}
	hs := &PaperworkServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		userUploader: newUploader,
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

// setUpMockGotenbergServer sets up a mock Gotenberg server and sets the corresponding env vars to make it easier to
// test functions that hit gotenberg endpoints. This asks the caller to pass in a handler function so that the caller
// can decide how the mock server should respond to their requests (e.g. setting a 404 status code or returning a
// specific file). The caller will need to close the mock server when they are done with it. The easiest way to ensure
// it happens even if the test fails is to defer the call to mockGotenbergServer.Close() right after calling this
// function.
// E.g.
//
//	mockGotenbergServer := suite.setUpMockGotenbergServer(func(w http.ResponseWriter, r *http.Request) {
//	    w.WriteHeader(http.StatusNotFound)
//	})
//
//	defer mockGotenbergServer.Close()
func (suite *PaperworkServiceSuite) setUpMockGotenbergServer(handlerFunc http.HandlerFunc) *httptest.Server {
	mockGotenbergServer := httptest.NewServer(http.HandlerFunc(handlerFunc))

	// The mock server sets its own url so we'll want to break it down and set the corresponding env vars
	url, urlParseErr := url.ParseRequestURI(mockGotenbergServer.URL)

	suite.FatalNoError(urlParseErr)

	os.Setenv(GotenbergProtocol, url.Scheme)
	os.Setenv(GotenbergHost, url.Hostname())
	os.Setenv(GotenbergPort, url.Port())

	return mockGotenbergServer
}

func (suite *PaperworkServiceSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		file.Close()
	}
}

func (suite *PaperworkServiceSuite) closeFile(file afero.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

func (suite *PaperworkServiceSuite) openLocalFile(path string, fs *afero.Afero) (afero.File, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

	outputFile, err := fs.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "error creating afero file")
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		return nil, errors.Wrap(err, "error copying over file contents")
	}

	suite.closeFile(outputFile)

	return outputFile, nil
}
