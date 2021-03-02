package webhooksubscription

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestWebhookSubscriptionUpdater() {
	builder := query.NewQueryBuilder(suite.DB())
	updater := NewWebhookSubscriptionUpdater(builder)

	// Create a webhook subscription
	origSub := testdatagen.MakeDefaultWebhookSubscription(suite.DB())

	suite.T().Run("Updates a webhook subscription successfully", func(t *testing.T) {

		newSub := models.WebhookSubscription{
			ID:          origSub.ID,
			CallbackURL: "/this/is/changed",
			Severity:    2,
			EventKey:    "Change.The.Event",
		}
		sev := int64(newSub.Severity)
		updatedSub, err := updater.UpdateWebhookSubscription(&newSub, &sev)

		suite.NoError(err)
		suite.Equal(newSub.CallbackURL, updatedSub.CallbackURL)
		suite.Equal(newSub.Severity, updatedSub.Severity)
		suite.Equal(newSub.EventKey, updatedSub.EventKey)
		suite.Equal(origSub.ID, updatedSub.ID)
		suite.Equal(origSub.Status, updatedSub.Status)
	})

	suite.T().Run("Fails to find correct webhookSubscription - return empty webhookSubscription and error", func(t *testing.T) {
		fakeID, _ := uuid.NewV4()

		newSub := models.WebhookSubscription{
			ID:          fakeID,
			CallbackURL: "/this/is/changed/again"}

		updatedSub, err := updater.UpdateWebhookSubscription(&newSub, nil)

		suite.Equal(models.RecordNotFoundErrorString, err.Error())
		suite.Nil(updatedSub)
	})

	suite.T().Run("Fails to update - return empty webhookSubscription and error", func(t *testing.T) {
		badWebhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		fakeID, _ := uuid.NewV4()

		newSub := models.WebhookSubscription{
			ID:           badWebhookSubscription.ID,
			SubscriberID: fakeID,
			CallbackURL:  "/this/is/changed/again",
		}

		updatedSub, err := updater.UpdateWebhookSubscription(&newSub, nil)

		suite.Error(err)
		suite.Nil(updatedSub)
	})
}
