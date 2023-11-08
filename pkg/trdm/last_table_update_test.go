package trdm_test

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/transcom/mymove/pkg/factory"
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

func (suite *TRDMSuite) TestLastTableUpdate() {
	// Setup mock creds
	mockCreds := aws.Credentials{
		AccessKeyID:     "mockAccessKeyID",
		SecretAccessKey: "mockSecretAccessKey",
		SessionToken:    "mockSessionToken",
		Source:          "mockProvider",
	}
	mockProvider := &mockAssumeRoleProvider{creds: mockCreds}
	// Set the configuration for the test
	suite.viper.Set(trdm.TrdmIamRoleFlag, "mockRole")
	suite.viper.Set(trdm.TrdmGatewayRegionFlag, "us-gov-west-1") // TODO: Possibly switch to var that itself is pulled from viper
	suite.viper.Set(trdm.GatewayURLFlag, "https://test.gateway.url.amazon.com")

	err := trdm.LastTableUpdate(suite.viper, &tls.Config{MinVersion: tls.VersionTLS13}, suite.AppContextForTest(), mockProvider)
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
