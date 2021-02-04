package webhooksubscription

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type WebhookSubscriptionServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestWebhookSubscriptionSuite(t *testing.T) {

	ts := &WebhookSubscriptionServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
