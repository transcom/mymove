package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestWebhookSubscription_NotNullConstraint() {
	webhookSubscription := &models.WebhookSubscription{}

	expectedErrors := map[string][]string{
		"subscriber_id": {"SubscriberID can not be blank."},
		"event_key":     {"EventKey can not be blank."},
		"status":        {"Status is not in the list [ACTIVE, DISABLED, FAILING]."},
		"callback_url":  {"CallbackURL can not be blank."},
	}

	verrs, err := webhookSubscription.Validate(suite.DB())
	if err != nil || verrs != nil {
		print(err)
	}

	suite.verifyValidationErrors(webhookSubscription, expectedErrors)
}

func (suite *ModelSuite) TestWebhookSubscription_Instantiation() {

	suite.Run("Default subscription", func() {
		webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())

		verrs, err := suite.DB().ValidateAndSave(&webhookSubscription)

		// Check that there were no errors
		suite.Nil(err, "could not save WebhookSubscription: %v", err)
		suite.Zero(verrs.Count(), "did not expect validation errors: %v", verrs)
		// Check default severity is set
		suite.Equal(0, webhookSubscription.Severity)
	})

	suite.Run("Updated subscription severity", func() {
		webhookSubscription := testdatagen.MakeWebhookSubscription(suite.DB(), testdatagen.Assertions{
			WebhookSubscription: models.WebhookSubscription{
				Severity: 2,
			},
		})
		verrs, err := suite.DB().ValidateAndSave(&webhookSubscription)

		// Check that there were no errors
		suite.Nil(err, "could not save WebhookSubscription: %v", err)
		suite.Zero(verrs.Count(), "did not expect validation errors: %v", verrs)
		// Check non default severity is set
		suite.Equal(2, webhookSubscription.Severity)
	})

}
