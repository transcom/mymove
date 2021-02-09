package webhooksubscription

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestCreateWebhookSubscription() {
	queryBuilder := query.NewQueryBuilder(suite.DB())
	subscriber := testdatagen.MakeContractor(suite.DB(), testdatagen.Assertions{})

	webhookSubscriptionInfo := models.WebhookSubscription{
		SubscriberID: subscriber.ID,
		Status:       models.WebhookSubscriptionStatusActive,
		EventKey:     "PaymentRequest.Update",
		CallbackURL:  "/my/callback/url",
	}

	// Happy path
	suite.T().Run("If the subscription is created successfully it should be returned", func(t *testing.T) {
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", subscriber.ID)}

		creator := NewWebhookSubscriptionCreator(suite.DB(), queryBuilder)
		webhookSubscription, verrs, err := creator.CreateWebhookSubscription(&webhookSubscriptionInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(webhookSubscription.ID)
		suite.NotNil(webhookSubscription.Severity)
		suite.Equal(webhookSubscriptionInfo.Status, webhookSubscription.Status)
	})

	// Bad subscriber ID
	suite.T().Run("If we are provided a organization that doesn't exist, the create should fail", func(t *testing.T) {
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", "b9c41d03-c730-4580-bd37-9ccf4845af6c")}

		creator := NewWebhookSubscriptionCreator(suite.DB(), queryBuilder)
		_, _, err := creator.CreateWebhookSubscription(&webhookSubscriptionInfo, filter)
		suite.Error(err)
		suite.Contains(err.Error(), "not found while looking for SubscriberID")
	})
}
