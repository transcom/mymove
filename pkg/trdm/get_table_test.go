package trdm_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/parser/loa"
	"github.com/transcom/mymove/pkg/trdm"
)

func (suite *TRDMSuite) TestGetTGETDataLOA() {
	getTableRequest := models.GetTableRequest{
		PhysicalName:                trdm.LineOfAccounting,
		ContentUpdatedSinceDateTime: time.Now().Add(time.Hour * 24 * 365 * 5),
		ReturnContent:               true,
	}

	// Mock a successful attachment return from GetTable
	// Use parser fixtures for LOA
	loaFile, err := os.ReadFile("../parser/loa/fixtures/Line Of Accounting.txt")
	suite.NoError(err)

	getTableResponse := models.GetTableResponse{
		RowCount:   2,
		StatusCode: trdm.SuccessfulStatusCode,
		DateTime:   time.Now(),
		Attachment: loaFile,
	}

	responseBody, err := json.Marshal(getTableResponse)
	suite.NoError(err)

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(responseBody)),
			}, nil
		}}

	mockProvider := &mockAssumeRoleProvider{creds: suite.creds}

	service := trdm.NewGatewayService(mockClient, suite.logger, "us-gov-west-1", "mockRole", "https://test.gateway.url.amazon.com", mockProvider)

	// GetTable for Line of Accounting, parse it, and then store it in the database
	err = trdm.GetTGETData(getTableRequest, *service, suite.AppContextForTest())
	suite.NoError(err)

	// Load our test file
	reader := bytes.NewReader(loaFile)
	expectedLOACodes, err := loa.Parse(reader)
	suite.NoError(err)

	// Load our loa codes from the DB
	var loaCodes []models.LineOfAccounting

	err = suite.AppContextForTest().DB().All(&loaCodes)
	suite.NoError(err)

	// Only compare len. Remember, the attachment does not
	// have the same automatically populated values we do in the DB
	// Such as UUID, Created, Updated, that sort. That's why we're
	// just checking to see if the values are now present, not 1:1 matches
	suite.Equal(len(expectedLOACodes), len(loaCodes))
}
