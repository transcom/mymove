package ordersapi

import (
	"log"
	"net/http"
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	handlers.BaseHandlerTestSuite
}

// AuthenticateClientCertRequest authenticates mutual TLS auth API users with the provided ClientCert object
func (suite *HandlerSuite) AuthenticateClientCertRequest(req *http.Request, cert *models.ClientCert) *http.Request {
	ctx := authentication.SetClientCertInRequestContext(req, cert)
	return req.WithContext(ctx)
}

// SetupTest sets up the test suite by preparing the DB
func (suite *HandlerSuite) SetupTest() {
	errTruncateAll := //RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return value in the file is used for test database teardown
		//RA: Given the database is being reset for unit test use, there are no unexpected states and conditions to account for
		//RA Developer Status: Mitigated
		//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
		//RA Validator: jneuner@mitre.org
		//RA Modified Severity:
		suite.DB().TruncateAll() // nolint:errcheck
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database: %#v", errTruncateAll)
	}
}

// AfterTest completes tests by trying to close open files
func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.TestFilesToClose() {
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used to close a local server connection to ensure a unit test server is not left running indefinitely
		//RA: Given the functions causing the lint errors are used to close a local server connection for testing purposes, it is not deemed a risk
		//RA Developer Status: Mitigated
		//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
		//RA Validator: jneuner@mitre.org
		//RA Modified Severity:
		file.Data.Close() // nolint:errcheck
	}
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &HandlerSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(logger, notifications.NewStubNotificationSender("milmovelocal", logger), testingsuite.CurrentPackage()),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func makeAllPowerfulClientCert() *models.ClientCert {
	return &models.ClientCert{
		AllowAirForceOrdersRead:     true,
		AllowAirForceOrdersWrite:    true,
		AllowArmyOrdersRead:         true,
		AllowArmyOrdersWrite:        true,
		AllowCoastGuardOrdersRead:   true,
		AllowCoastGuardOrdersWrite:  true,
		AllowDpsAuthAPI:             true,
		AllowMarineCorpsOrdersRead:  true,
		AllowMarineCorpsOrdersWrite: true,
		AllowNavyOrdersRead:         true,
		AllowNavyOrdersWrite:        true,
		AllowOrdersAPI:              true,
	}
}
