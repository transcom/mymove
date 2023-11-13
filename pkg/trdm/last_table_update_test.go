package trdm_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/trdm"
)

func (suite *TRDMSuite) TestFetchLOARecordsByTime() {
	// Get initial TAC codes count
	initialCodes, err := trdm.FetchLOARecordsByTime(suite.AppContextForTest(), time.Now())
	suite.NoError(err)
	intialLoaCodesLength := len(initialCodes)

	// Creates a test TAC code record in the DB
	factory.BuildDefaultLineOfAccounting(suite.DB())

	// Fetch All TAC Records
	// !! A second time.Now() statement is intentional based on the SQL query.
	codes, err := trdm.FetchLOARecordsByTime(suite.AppContextForTest(), time.Now())
	suite.NoError(err)
	// Compare new TAC Code count to initial count
	finalCodesLength := len(codes)

	suite.NoError(err)
	suite.NotEqual(finalCodesLength, intialLoaCodesLength)
}

func (suite *TRDMSuite) TestFetchTACRecordsByTime() {
	// Get initial TAC codes count
	initialCodes, err := trdm.FetchTACRecordsByTime(suite.AppContextForTest(), time.Now())
	suite.NoError(err)
	initialTacCodeLength := len(initialCodes)

	// Creates a test TAC code record in the DB
	factory.BuildDefaultTransportationAccountingCode(suite.DB())

	// Fetch All TAC Records
	codes, err := trdm.FetchTACRecordsByTime(suite.AppContextForTest(), time.Now())
	suite.NoError(err)

	// Compare new TAC Code count to initial count
	finalCodesLength := len(codes)

	suite.NotEqual(finalCodesLength, initialTacCodeLength)
}

// Mock provider to inject our own dummy creds
type mockAssumeRoleProvider struct {
	creds aws.Credentials
}

// Implement retrieve to fully mock, see aws go sdk v2
// sts assume role provider retrieve
// Throw away var because it does get used in real code
// Our test overrides the need for it though.
func (m *mockAssumeRoleProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	return m.creds, nil
}

// Mock HTTP client for test injection, we're providing an HTTPClient type now in the
// non test functions of the package
type MockHTTPClient struct {
	// DoFunc allows us to implement our own responses from when
	// executing client.Do()
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func (suite *TRDMSuite) TestLastTableUpdate() {
	// Setup mock creds
	mockCreds := aws.Credentials{
		AccessKeyID:     "mockAccessKeyID",
		SecretAccessKey: "mockSecretAccessKey",
		SessionToken:    "mockSessionToken",
		Source:          "mockProvider",
	}
	mockProvider := &mockAssumeRoleProvider{creds: mockCreds}

	lastUpdateResponse := models.LastTableUpdateResponse{
		StatusCode: trdm.SuccessfulStatusCode,
		DateTime:   time.Now(),
		LastUpdate: time.Now().Add(-24 * time.Hour),
	}

	responseBody, err := json.Marshal(lastUpdateResponse)
	suite.NoError(err)

	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(responseBody)),
			}, nil
		},
	}
	// Set the configuration for the test
	suite.viper.Set(trdm.TrdmIamRoleFlag, "mockRole")
	suite.viper.Set(trdm.TrdmGatewayRegionFlag, "us-gov-west-1") // TODO: Possibly switch to var that itself is pulled from viper
	suite.viper.Set(trdm.GatewayURLFlag, "https://test.gateway.url.amazon.com")

	err = trdm.LastTableUpdate(suite.viper, suite.AppContextForTest(), mockProvider, mockHTTPClient)
	suite.NoError(err)
}

func (suite *TRDMSuite) TestFetchAllTACRecords() {
	// Get initial TAC codes count
	initialCodes, err := trdm.FetchAllTACRecords(suite.AppContextForTest())
	initialTacCodeLength := len(initialCodes)
	suite.NoError(err)

	// Creates a test TAC code record in the DB
	factory.BuildFullTransportationAccountingCode(suite.DB())

	// Fetch All TAC Records
	codes, err := trdm.FetchAllTACRecords(suite.AppContextForTest())

	// Compare new TAC Code count to initial count
	finalCodesLength := len(codes)

	suite.NoError(err)
	suite.NotEqual(finalCodesLength, initialTacCodeLength)
}
