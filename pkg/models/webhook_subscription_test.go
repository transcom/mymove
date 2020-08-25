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
		"status":        {"Status is not in the list [ACTIVE, DISABLED]."},
		"callback_url":  {"CallbackURL can not be blank."},
	}

	verrs, err := webhookSubscription.Validate(suite.DB())
	if err != nil || verrs != nil {
		print(err)
	}

	suite.verifyValidationErrors(webhookSubscription, expectedErrors)
}

func (suite *ModelSuite) TestWebhookSubscription_Instantiation() {
	t := suite.T()
	webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())

	verrs, err := suite.DB().ValidateAndSave(&webhookSubscription)

	if err != nil {
		t.Fatalf("could not save WebhookSubscription: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}
}