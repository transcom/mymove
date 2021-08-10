package webhooksubscription

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestWebhookSubscriptionFetcher() {
	builder := query.NewQueryBuilder()
	fetcher := NewWebhookSubscriptionFetcher(builder)

	webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	webhookSubscriptionID := webhookSubscription.ID

	suite.T().Run("Get a webhook subscription successfully", func(t *testing.T) {
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", webhookSubscription.ID.String())}

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		webhookSubscription, err := fetcher.FetchWebhookSubscription(appCfg, filters)

		suite.NoError(err)
		suite.Equal(webhookSubscriptionID, webhookSubscription.ID)
	})

	suite.T().Run("Failure to fetch - return empty webhookSubscription and error", func(t *testing.T) {
		fakeID, err := uuid.NewV4()
		suite.NoError(err)

		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", fakeID.String())}

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		fakeWebhookSubscription, err := fetcher.FetchWebhookSubscription(appCfg, filters)

		suite.Error(err)
		suite.Equal(models.WebhookSubscription{}, fakeWebhookSubscription)
	})

}
