package invoice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type GexSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func TestGexSuite(t *testing.T) {

	ts := &GexSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("gex")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GexSuite) TestSendToGexHTTP_Call() {
	ediString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	resp, err := NewGexSenderHTTP(mockServer.URL, "", false, nil, "", "").
		SendToGex(ediString, "test_transaction")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}
	expectedStatus := http.StatusOK
	suite.Equal(expectedStatus, resp.StatusCode)

	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	resp, err = NewGexSenderHTTP(mockServer.URL, "", false, nil, "", "").
		SendToGex(ediString, "test_transaction")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}

	expectedStatus = http.StatusInternalServerError
	suite.Equal(expectedStatus, resp.StatusCode)
}

func (suite *GexSuite) TestSendToGexHTTP_QueryParams() {
	ediString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	resp, err := NewGexSenderHTTP(mockServer.URL, "test_channel", false, nil, "", "").
		SendToGex(ediString, "test_filename")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}
	expectedStatus := http.StatusOK
	suite.Equal(expectedStatus, resp.StatusCode)

	// Make sure we sent a request with correct channel and filename query parameters
	suite.Contains(resp.Request.URL.RawQuery, "channel=test_channel")
	suite.Contains(resp.Request.URL.RawQuery, "fname=test_filename")
}
