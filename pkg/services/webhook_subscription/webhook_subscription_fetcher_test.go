package webhooksubscription

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestWebhookSubscriptionFetcher() {
	builder := query.NewQueryBuilder()
	fetcher := NewWebhookSubscriptionFetcher(builder)

	suite.Run("Get a webhook subscription successfully", func() {
		webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		webhookSubscriptionID := webhookSubscription.ID
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", webhookSubscription.ID.String())}

		webhookSubscription, err := fetcher.FetchWebhookSubscription(suite.AppContextForTest(), filters)

		suite.NoError(err)
		suite.Equal(webhookSubscriptionID, webhookSubscription.ID)
	})

	suite.Run("Failure to fetch - return empty webhookSubscription and error", func() {
		fakeID, err := uuid.NewV4()
		suite.NoError(err)

		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", fakeID.String())}

		fakeWebhookSubscription, err := fetcher.FetchWebhookSubscription(suite.AppContextForTest(), filters)

		suite.Error(err)
		suite.Equal(models.WebhookSubscription{}, fakeWebhookSubscription)
	})

}
