package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildWebhookNotification() {
	suite.Run("Successful creation of default WebhookNotification", func() {
		// Under test:      BuildWebhookNotification
		// Mocked:          None
		// Set up:          Create a webhook notification with no customizations or traits
		// Expected outcome:webhookNotification should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		webhookNotification := BuildWebhookNotification(suite.DB(), nil, nil)

		suite.NotNil(webhookNotification.MoveTaskOrderID)
		suite.Equal(models.WebhookNotificationPending, webhookNotification.Status)
		suite.Equal("Payment.Create", webhookNotification.EventKey)
	})

	suite.Run("Successful creation of customized WebhookNotification", func() {
		// Under test:      BuildWebhookNotification
		// Mocked:          None
		// Set up:          Create a webhook notification and pass custom fields
		// Expected outcome:webhookNotification should be created with custom values

		// SETUP
		customPayload := "{\"message\":\"This is an test notification #1\"}"

		customWebhookNotification := models.WebhookNotification{
			Status:  models.WebhookNotificationSent,
			Payload: customPayload,
		}

		// CALL FUNCTION UNDER TEST
		webhookNotification := BuildWebhookNotification(suite.DB(), []Customization{
			{
				Model: customWebhookNotification,
			},
		}, nil)

		suite.Equal(customWebhookNotification.Status, webhookNotification.Status)
		suite.Equal(customWebhookNotification.Payload, webhookNotification.Payload)
	})

	suite.Run("Successful return of linkOnly WebhookNotification", func() {
		// Under test:       BuildWebhookNotification
		// Set up:           Pass in a linkOnly webhookNotification
		// Expected outcome: No new WebhookNotification should be created.

		// Check num WebhookNotification records
		precount, err := suite.DB().Count(&models.WebhookNotification{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		webhookNotification := BuildWebhookNotification(suite.DB(), []Customization{
			{
				Model: models.WebhookNotification{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, webhookNotification.ID)
	})
}
