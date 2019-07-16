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
	logger Logger
}

func (suite *GexSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestGexSuite(t *testing.T) {

	hs := &GexSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("gex")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
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
