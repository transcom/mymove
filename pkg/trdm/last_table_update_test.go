package trdm_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/parser/loa"
	"github.com/transcom/mymove/pkg/parser/tac"
	"github.com/transcom/mymove/pkg/trdm"
)

const (
	gatewayURL string = "https://test.gateway.url.amazon.com"
)

func (suite *TRDMSuite) TestTGETLOADataOutOfDate() {
	newTime := time.Now().Add(1 * time.Hour)

	// We will have no LOA data yet, forcing us to immediately be out of date
	exists, _, err := trdm.TGETLOADataOutOfDate(suite.AppContextForTest(), newTime)
	suite.NoError(err)
	suite.True(exists)

	// Put a LOA inside the DB
	loa := factory.BuildDefaultLineOfAccounting(suite.DB())
	err = suite.DB().Update(&loa)
	suite.NoError(err)

	// Ensure that it finds the new LOA
	exists, _, err = trdm.TGETLOADataOutOfDate(suite.AppContextForTest(), newTime)
	suite.NoError(err)
	suite.True(exists)
}

func (suite *TRDMSuite) TestTGETTACDataOutOfDate() {
	newTime := time.Now().Add(1 * time.Hour)

	// We will have no TAC data yet, forcing us to immediately be out of date
	exists, _, err := trdm.TGETTACDataOutOfDate(suite.AppContextForTest(), newTime)
	suite.NoError(err)
	suite.True(exists)

	// Put a TAC inside the DB
	tac := factory.BuildDefaultTransportationAccountingCode(suite.DB())
	err = suite.DB().Update(&tac)
	suite.NoError(err)

	// Ensure that it finds the new TAC
	exists, _, err = trdm.TGETTACDataOutOfDate(suite.AppContextForTest(), newTime)
	suite.NoError(err)
	suite.True(exists)
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

// This is a rather large test. It will test lastTableUpdate triggering the need to retrieve
// additional TGET data for both TAC and LOA tables. It will then test a full
// successful TRDM cron flow
func (suite *TRDMSuite) TestSuccessfulTRDMFlowTACsAndLOAs() {
	// Factory customizations for TAC generation
	tacCustoms := []factory.Customization{
		{
			Model: models.TransportationAccountingCode{
				CreatedAt: time.Now().Add(time.Hour * 24 * 365 * -1),
				UpdatedAt: time.Now().Add(time.Hour * 24 * 365 * -1),
			},
		},
	}

	// Factory customizations for LOA generation
	loaCustoms := []factory.Customization{
		{
			Model: models.LineOfAccounting{
				CreatedAt: time.Now().Add(time.Hour * 24 * 365 * -1),
				UpdatedAt: time.Now().Add(time.Hour * 24 * 365 * -1),
			},
		},
	}

	// Create multiple old TACs in the DB
	outdatedTACCodes := make([]models.TransportationAccountingCode, 4)
	for i := range outdatedTACCodes {
		outdatedTACCodes[i] = factory.BuildTransportationAccountingCode(suite.DB(), tacCustoms, nil)
	}

	// Create multiple old LOAs in the DB
	outdatedLOACodes := make([]models.LineOfAccounting, 4)
	for i := range outdatedLOACodes {
		outdatedLOACodes[i] = factory.BuildLineOfAccounting(suite.DB(), loaCustoms, nil)
	}

	// Now pull the fixtures for testing
	tacFile, err := os.ReadFile("../parser/tac/fixtures/Transportation Account.txt")
	suite.NoError(err)
	loaFile, err := os.ReadFile("../parser/loa/fixtures/Line Of Accounting.txt")
	suite.NoError(err)

	// Mock creds provider
	mockProvider := &mockAssumeRoleProvider{creds: suite.creds}

	// Mock response from LastTableUpdate
	// We'll use the same response for both
	lastUpdateResponse := models.LastTableUpdateResponse{
		StatusCode: trdm.SuccessfulStatusCode,
		DateTime:   time.Now(),
		LastUpdate: time.Now(),
	}
	lastTableUpdateResponseBody, err := json.Marshal(lastUpdateResponse)
	suite.NoError(err)

	// Mock response from GetTable TAC
	getTableResponseTAC := models.GetTableResponse{
		RowCount:   1,
		StatusCode: trdm.SuccessfulStatusCode,
		DateTime:   time.Now(),
		Attachment: tacFile, // We're going to return our fixtureTACs to be added to the DB
	}

	getTableResponseBodyTAC, err := json.Marshal(getTableResponseTAC)
	suite.NoError(err)

	// Mock response from GetTable LOA
	getTableResponseLOA := models.GetTableResponse{
		RowCount:   1,
		StatusCode: trdm.SuccessfulStatusCode,
		DateTime:   time.Now(),
		Attachment: loaFile, // We're going to return our fixtureTACs to be added to the DB
	}

	getTableResponseBodyLOA, err := json.Marshal(getTableResponseLOA)
	suite.NoError(err)

	// Mock client and simulate the 2 different getTable responses and the reused lastTableUpdate response
	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			switch req.URL.String() {
			// Simulate lastTablUpdate response
			case gatewayURL + trdm.LastTableUpdateEndpoint:
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(lastTableUpdateResponseBody)),
				}, nil
				// Simulate getTable response
			case gatewayURL + trdm.GetTableEndpoint:
				// Parse the incoming request to determine if it's
				// trying to get LOA or TAC data from getTable
				var getTableRequest models.GetTableRequest
				// getTableErr used to avoid err shadowing
				reqBody, getTableErr := io.ReadAll(req.Body)
				suite.NoError(getTableErr)
				getTableErr = json.Unmarshal(reqBody, &getTableRequest)
				suite.NoError(getTableErr)

				// Give two different responses for LOA and TAC
				switch getTableRequest.PhysicalName {
				case trdm.LineOfAccounting:
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(getTableResponseBodyLOA)),
					}, nil
				case trdm.TransportationAccountingCode:
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(getTableResponseBodyTAC)),
					}, nil
					// If unable to determine LOA or TAC for getTable then error
				default:
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       nil,
					}, nil
				}

			default:
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       nil,
				}, nil
			}
		},
	}

	// Set the configuration for the test
	suite.viper.Set(trdm.TrdmIamRoleFlag, "mockRole")
	suite.viper.Set(trdm.TrdmGatewayRegionFlag, "us-gov-west-1") // TODO: Possibly switch to var that itself is pulled from viper
	suite.viper.Set(trdm.GatewayURLFlag, gatewayURL)

	// Now our flow will run and it will detect we have outdated TACs and LOAs.
	// When it finds these outdated TACs/LOAs it should trigger the `getTable` part of the flow
	// and attempt to store the new data found
	err = trdm.BeginTGETFlow(suite.viper, suite.AppContextForTest(), mockProvider, mockHTTPClient)
	suite.NoError(err)

	// Pull the TAC and LOA entries, they should have both the old and new ones
	var allTACs []models.TransportationAccountingCode
	err = suite.DB().All(&allTACs)
	suite.NoError(err)
	var allLOAs []models.LineOfAccounting
	err = suite.DB().All(&allLOAs)
	suite.NoError(err)

	/*
	**** Parse our fixture files ****
	 */

	// TAC
	reader := bytes.NewReader(tacFile)
	expectedTACCodes, err := tac.Parse(reader)
	suite.NoError(err)
	// Consolidate duplicates
	expectedTACCodes = tac.ConsolidateDuplicateTACsDesiredFromTRDM(expectedTACCodes)

	// LOA
	reader = bytes.NewReader(loaFile)
	expectedLOACodes, err := loa.Parse(reader)
	suite.NoError(err)

	// Now lets see if our new TACs and LOAs were stored
	suite.Equal(len(allTACs), len(outdatedTACCodes)+len(expectedTACCodes))
	// Also include the TAC length in the LOA calculation.
	// On factory build TAC, a LOA is generated alongside it.
	suite.Equal(len(allLOAs), len(outdatedLOACodes)+len(expectedLOACodes)+len(outdatedTACCodes))
}
func (suite *TRDMSuite) TestFetchWeeksOfMissingTime() {
	// These errs can be "_" because no matter what it will error out after FetchWeeksOfMissingData is called
	// Additionally, we are providing correct parameters that will never error
	ourLastUpdate, _ := time.Parse("Jan 02, 2006", "Aug 01, 2023")
	trdmLastUpdate, _ := time.Parse("Jan 02, 2006", "Aug 14, 2023")
	missingWeeks, err := trdm.FetchWeeksOfMissingTime(ourLastUpdate, trdmLastUpdate)
	suite.NoError(err)

	suite.Equal(len(missingWeeks), 2) // Assert 2 weeks are missing from the provided dates

	trdmLastUpdate, _ = time.Parse("Jan 02, 2006", "Aug 15, 2023") // 2 weeks and 1 day after our last update

	missingWeeks, err = trdm.FetchWeeksOfMissingTime(ourLastUpdate, trdmLastUpdate)
	suite.NoError(err)

	suite.Equal(len(missingWeeks), 3)
}

func (suite *TRDMSuite) TestIncorrectParametersForWeeksOfMissingTime() {
	// These errs can be "_" because no matter what it will error out after FetchWeeksOfMissingData is called
	// Additionally, we are providing correct parameters that will never error
	ourLastUpdate, _ := time.Parse("Jan 02, 2006", "Aug 02, 2023")
	trdmLastUpdate, _ := time.Parse("Jan 02, 2006", "Aug 01, 2023") // trdmLastUpdate is before our last update
	_, err := trdm.FetchWeeksOfMissingTime(ourLastUpdate, trdmLastUpdate)
	suite.Error(err)
}
