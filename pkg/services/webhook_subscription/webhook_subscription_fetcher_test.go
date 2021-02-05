package webhooksubscription

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestWebhookSubscriptionFetcher() {
	builder := query.NewQueryBuilder(suite.DB())
	fetcher := NewWebhookSubscriptionFetcher(builder)

	webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	webhookSubscriptionID := webhookSubscription.ID

	suite.T().Run("Get a webhook subscription successfully", func(t *testing.T) {
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", webhookSubscription.ID.String())}

		webhookSubscription, err := fetcher.FetchWebhookSubscription(filters)

		suite.NoError(err)
		suite.Equal(webhookSubscriptionID, webhookSubscription.ID)
	})

	suite.T().Run("Fetch error returns nil", func(t *testing.T) {
		fakeID, err := uuid.NewV4()
		suite.NoError(err)

		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", fakeID.String())}

		fakeWebhookSubscription, err := fetcher.FetchWebhookSubscription(filters)

		suite.Error(err)
		suite.Equal(models.WebhookSubscription{}, *fakeWebhookSubscription)
	})

}
