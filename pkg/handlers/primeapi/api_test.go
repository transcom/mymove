package primeapi

import (
	"log"
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	handlers.BaseHandlerTestSuite
}

// AfterTest completes tests by trying to close open files
func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.TestFilesToClose() {
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used to close a local server connection to ensure a unit test server is not left running indefinitely
		//RA: Given the functions causing the lint errors are used to close a local server connection for testing purposes, it is not deemed a risk
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: N/A
		// nolint:errcheck
		file.Data.Close()
	}
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &HandlerSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(logger, notifications.NewStubNotificationSender("milmovelocal", logger), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

// EqualAddress compares a model address against a payload address
func (suite *HandlerSuite) EqualAddress(expected models.Address, actual *primemessages.Address, checkID bool) {
	if checkID == true {
		suite.Equal(expected.ID.String(), actual.ID.String())
	}
	suite.Equal(&expected.StreetAddress1, actual.StreetAddress1)
	suite.Equal(expected.StreetAddress2, actual.StreetAddress2)
	suite.Equal(expected.StreetAddress3, actual.StreetAddress3)
	suite.Equal(&expected.City, actual.City)
	suite.Equal(&expected.State, actual.State)
	suite.Equal(&expected.PostalCode, actual.PostalCode)
	suite.Equal(expected.Country, actual.Country)
}

// EqualAddressPayload compares a payload address against a payload address
// If you don't want to compare IDs set checkID to false
func (suite *HandlerSuite) EqualAddressPayload(expected *primemessages.Address, actual *primemessages.Address, checkID bool) {
	if expected == nil || actual == nil {
		suite.Nil(expected)
		suite.Nil(actual)
	}
	if checkID == true {
		suite.Equal(expected.ID.String(), actual.ID.String())
	}
	suite.Equal(expected.StreetAddress1, actual.StreetAddress1)
	suite.Equal(expected.StreetAddress2, actual.StreetAddress2)
	suite.Equal(expected.StreetAddress3, actual.StreetAddress3)
	suite.Equal(expected.City, actual.City)
	suite.Equal(expected.State, actual.State)
	suite.Equal(expected.PostalCode, actual.PostalCode)
	suite.Equal(expected.Country, actual.Country)
}
