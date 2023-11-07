package trdm_test

import (
	"crypto/tls"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/trdm"
)

// const (
// 	physicalName = "fakePhysicalName"
// )

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

func (suite *TRDMSuite) TestLastTableUpdate() {
	err := trdm.LastTableUpdate(suite.viper, &tls.Config{MinVersion: tls.VersionTLS13}, suite.AppContextForTest())
	suite.NoError(err)
}
