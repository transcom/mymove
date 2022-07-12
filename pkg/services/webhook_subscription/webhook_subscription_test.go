package webhooksubscription

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type WebhookSubscriptionServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestWebhookSubscriptionSuite(t *testing.T) {

	ts := &WebhookSubscriptionServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
