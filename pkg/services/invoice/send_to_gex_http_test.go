package invoice

import (
	"github.com/transcom/mymove/pkg/testingsuite"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type GexSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *GexSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestGexSuite(t *testing.T) {

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &GexSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, hs)
}

func (suite *GexSuite) TestSendToGexHTTP_Call() {
	ediString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	resp, err := NewGexSenderHTTP(mockServer.URL, false, nil, "", "").
		SendToGex(ediString, "test_transaction")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}
	expectedStatus := http.StatusOK
	suite.Equal(expectedStatus, resp.StatusCode)

	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	resp, err = NewGexSenderHTTP(mockServer.URL, false, nil, "", "").
		SendToGex(ediString, "test_transaction")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}

	expectedStatus = http.StatusInternalServerError
	suite.Equal(expectedStatus, resp.StatusCode)
}
