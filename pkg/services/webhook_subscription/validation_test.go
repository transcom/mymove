package webhooksubscription

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WebhookSubscriptionServiceSuite) TestWebhookSubscriptionValidation() {

}

func (suite *WebhookSubscriptionServiceSuite) TestWebhookSubscriptionRules() {
	builder := query.NewQueryBuilder()

	subscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	invalidSubscription := models.WebhookSubscription{
		SubscriberID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111")),
		Status:       models.WebhookSubscriptionStatusActive,
		EventKey:     "PaymentRequest.Update",
		CallbackURL:  "/my/callback/url",
	}
	suite.Run("checkSubscriberExists", func() {
		suite.Run("success", func() {
			err := checkSubscriberExists(builder).Validate(suite.AppContextForTest(), subscription)
			suite.Require().NoError(err)
		})
		suite.Run("failure", func() {
			err := checkSubscriberExists(builder).Validate(suite.AppContextForTest(), invalidSubscription)
			suite.Require().Error(err)
		})
	})

}
