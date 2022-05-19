package invoice

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type GexSuite struct {
	testingsuite.PopTestSuite
}

func TestGexSuite(t *testing.T) {

	ts := &GexSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("gex"),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GexSuite) TestSendToGexHTTP_Call() {
	bodyString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	resp, err := NewGexSenderHTTP(mockServer.URL, false, nil, "", "").
		SendToGex(services.GEXChannelInvoice, bodyString, "test_transaction")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}
	expectedStatus := http.StatusOK
	suite.Equal(expectedStatus, resp.StatusCode)

	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	resp, err = NewGexSenderHTTP(mockServer.URL, false, nil, "", "").
		SendToGex(services.GEXChannelInvoice, bodyString, "test_transaction")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}

	expectedStatus = http.StatusInternalServerError
	suite.Equal(expectedStatus, resp.StatusCode)
}

func (suite *GexSuite) TestSendToGexHTTP_QueryParams() {
	bodyString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	resp, err := NewGexSenderHTTP(mockServer.URL, false, nil, "", "").
		SendToGex(services.GEXChannelInvoice, bodyString, "test_filename")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}
	expectedStatus := http.StatusOK
	suite.Equal(expectedStatus, resp.StatusCode)

	// Make sure we sent a request with correct channel and filename query parameters
	suite.Contains(resp.Request.URL.RawQuery, fmt.Sprintf("channel=%s", services.GEXChannelInvoice))
	suite.Contains(resp.Request.URL.RawQuery, "fname=test_filename")
}

func (suite *GexSuite) TestSendToGexHTTP_InvalidChannel() {
	bodyString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	var invalidChannel services.GEXChannel = "INVALID-CHANNEL"
	resp, err := NewGexSenderHTTP(mockServer.URL, false, nil, "", "").
		SendToGex(invalidChannel, bodyString, "test_filename")
	suite.Nil(resp)
	suite.NotNil(err)
	suite.Equal("Invalid channel type, expected [\"TRANSCOM-DPS-MILMOVE-CPS-IN-USBANK-RCOM\" \"TRANSCOM-DPS-MILMOVE-GHG-IN-IGC-RCOM\"]", err.Error())
}
