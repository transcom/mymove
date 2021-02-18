package webhooksubscription

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestWebhookSubscriptionUpdater() {
	builder := query.NewQueryBuilder(suite.DB())
	updater := NewWebhookSubscriptionUpdater(builder)

	webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	webhookSubscriptionID := webhookSubscription.ID

	eTag := etag.GenerateEtag(webhookSubscription.UpdatedAt)

	suite.T().Run("Etag is stale", func(t *testing.T) {
		eTag = etag.GenerateEtag(time.Now())
		_, err := updater.UpdateWebhookSubscription(&webhookSubscription, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Gets a webhook subscription successfully", func(t *testing.T) {
		// filters := []services.QueryFilter{query.NewQueryFilter("id", "=", webhookSubscription.ID.String())}
		// var foundWebhookSubscription models.WebhookSubscription
		// var eTag = "eTag"
		eTag = etag.GenerateEtag(webhookSubscription.UpdatedAt)
		foundWebhookSubscription, err := updater.UpdateWebhookSubscription(&webhookSubscription, eTag)

		suite.NoError(err)
		suite.Equal(webhookSubscriptionID, foundWebhookSubscription.ID)
	})

	// suite.T().Run("Fails to update - return empty webhookSubscription and error", func(t *testing.T) {
	// 	// fakeID, err := uuid.NewV4()
	// 	// suite.NoError(err)

	// 	// filters := []services.QueryFilter{query.NewQueryFilter("id", "=", fakeID.String())}
	// 	// var foundWebhookSubscription models.WebhookSubscription

	// 	fakeWebhookSubscription, err := updater.UpdateWebhookSubscription(&webhookSubscription, eTag)

	// 	suite.Error(err)
	// 	suite.Equal(models.WebhookSubscription{}, fakeWebhookSubscription)
	// })

}
