package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SpaHandlerSuite struct {
	*testingsuite.PopTestSuite
	logger *zap.Logger
	mfs    afero.HttpFs
}

func setupMockFileSystem() *afero.HttpFs {
	afs := afero.NewMemMapFs()

	errMkdir := afs.MkdirAll("test", 0755)

	if errMkdir != nil {
		log.Panic(errMkdir)
	}

	errWriteFile := afero.WriteFile(afs, "/test/a", []byte("file a"), 0644)

	if errWriteFile != nil {
		log.Panic(errWriteFile)
	}

	errWriteFile = afero.WriteFile(afs, "/test/index.html", []byte("index html file"), 0644)

	if errWriteFile != nil {
		log.Panic(errWriteFile)
	}

	errMkdir = afs.MkdirAll("/test/noIndexDir", 0755)

	if errMkdir != nil {
		log.Panic(errMkdir)
	}

	errWriteFile = afero.WriteFile(afs, "/test/noIndexDir/b", []byte("file b"), 0644)

	if errWriteFile != nil {
		log.Panic(errWriteFile)
	}

	ahttpFs := afero.NewHttpFs(afs)
	return ahttpFs
}

func TestSpaHandlerSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	mfs := setupMockFileSystem()

	hs := &SpaHandlerSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       logger,
		mfs:          *mfs,
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

type testCase struct {
	name               string
	request            string
	expectedStatusCode int
	expectedBody       string
}

func (suite *SpaHandlerSuite) TestSpaHandlerServeHttp() {
	cases := []testCase{
		{
			name:               "A directory without a trailing slash and that has an index.html",
			request:            "test",
			expectedStatusCode: http.StatusMovedPermanently,
			expectedBody:       "",
		},
		{
			name:               "A directory with a trailing slash and that has an index.html",
			request:            "test/",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "index html file",
		},
		{
			name:               "A directory without a trailing slash and that does not have an index.html",
			request:            "test/noIndexDir",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "404 page not found\n",
		},
		{
			name:               "A directory with a trailing slash and that does not have an index.html",
			request:            "test/noIndexDir/",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "404 page not found\n",
		},
		{
			name:               "A file that exists in a directory that does have an index.html",
			request:            "test/a",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "file a",
		},
		{
			name:               "A file that exists in a directory that does not have an index.html",
			request:            "test/noIndexDir/b",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "file b",
		},
		{
			name:               "A file that does not exist",
			request:            "test/c",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "404 page not found\n",
		},
	}

	cfs := NewCustomFileSystem(
		suite.mfs,
		"index.html",
		suite.logger,
	)

	for _, testCase := range cases {
		suite.T().Run(testCase.name, func(t *testing.T) {

			req, err := http.NewRequest("GET", testCase.request, nil)
			suite.NoError(err)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			sh := NewSpaHandler(cfs)
			sh.ServeHTTP(rr, req)

			suite.Equal(testCase.expectedStatusCode, rr.Code, "Status codes did not match when retreiving %v for request %v: expected %v, got %v", testCase.name, testCase.request, testCase.expectedStatusCode, rr.Code)

			// Check the response body is what we expect.
			suite.Equal(testCase.expectedBody, rr.Body.String(), "Handler returned unexpected body when retrieving %v: expected %v, got %v", testCase.name, testCase.expectedBody, rr.Body.String())
		})
	}
}
