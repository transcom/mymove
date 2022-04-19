package webhooksubscription

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestCreateWebhookSubscription() {
	queryBuilder := query.NewQueryBuilder()
	subscriber := testdatagen.MakeContractor(suite.DB(), testdatagen.Assertions{})

	webhookSubscriptionInfo := models.WebhookSubscription{
		SubscriberID: subscriber.ID,
		Status:       models.WebhookSubscriptionStatusActive,
		EventKey:     "PaymentRequest.Update",
		CallbackURL:  "/my/callback/url",
	}

	// Happy path
	suite.T().Run("If the subscription is created successfully it should be returned", func(t *testing.T) {
		creator := NewWebhookSubscriptionCreator(queryBuilder)
		webhookSubscription, verrs, err := creator.CreateWebhookSubscription(suite.AppContextForTest(), &webhookSubscriptionInfo)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(webhookSubscription.ID)
		suite.NotNil(webhookSubscription.Severity)
		suite.Equal(webhookSubscriptionInfo.Status, webhookSubscription.Status)
	})

	// Bad subscriber ID
	suite.T().Run("If we are provided an organization that doesn't exist, the create should fail", func(t *testing.T) {
		creator := NewWebhookSubscriptionCreator(queryBuilder)
		invalidSubscription := models.WebhookSubscription{
			SubscriberID: uuid.Must(uuid.FromString("b9c41d03-c730-4580-bd37-9ccf4845af6c")),
			Status:       models.WebhookSubscriptionStatusActive,
			EventKey:     "PaymentRequest.Update",
			CallbackURL:  "",
		}
		_, _, err := creator.CreateWebhookSubscription(suite.AppContextForTest(), &invalidSubscription)
		suite.Error(err)
		suite.Contains(err.Error(), "not found while looking for SubscriberID")
	})
}
