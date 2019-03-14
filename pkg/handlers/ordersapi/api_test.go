package ordersapi

import (
	"log"
	"net/http"
	"testing"

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
	suite.DB().TruncateAll()
}

// AfterTest completes tests by trying to close open files
func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.TestFilesToClose() {
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
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(logger, notifications.NewStubNotificationSender("milmovelocal", logger)),
	}

	suite.Run(t, hs)
	// clean up whatever the last test added to the DB
	hs.DB().TruncateAll()
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
