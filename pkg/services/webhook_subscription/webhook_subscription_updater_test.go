package webhooksubscription

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestWebhookSubscriptionUpdater() {
	builder := query.NewQueryBuilder(suite.DB())
	updater := NewWebhookSubscriptionUpdater(builder)

	webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	webhookSubscriptionID := webhookSubscription.ID

	// TODO: Add eTags back in once eTag functionality is created
	// eTag := etag.GenerateEtag(webhookSubscription.UpdatedAt)

	// suite.T().Run("Etag is stale", func(t *testing.T) {
	// 	eTag = etag.GenerateEtag(time.Now())
	// 	_, err := updater.UpdateWebhookSubscription(&webhookSubscription, eTag)

	// 	suite.Error(err)
	// 	suite.IsType(services.PreconditionFailedError{}, err)
	// })

	suite.T().Run("Gets a webhook subscription successfully", func(t *testing.T) {
		// eTag = etag.GenerateEtag(webhookSubscription.UpdatedAt)
		foundWebhookSubscription, err := updater.UpdateWebhookSubscription(&webhookSubscription)

		suite.NoError(err)
		suite.Equal(webhookSubscriptionID, foundWebhookSubscription.ID)
	})

	suite.T().Run("Fails to find correct webhookSuscription - return empty webhookSubscription and error", func(t *testing.T) {
		fakeID, _ := uuid.NewV4()
		webhookSubscription.ID = fakeID

		_, err := updater.UpdateWebhookSubscription(&webhookSubscription)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Fails to update - return empty webhookSubscription and error", func(t *testing.T) {
		badWebhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		fakeID, _ := uuid.NewV4()
		badWebhookSubscription.SubscriberID = fakeID

		_, err := updater.UpdateWebhookSubscription(&badWebhookSubscription)

		suite.Error(err)
	})

}
