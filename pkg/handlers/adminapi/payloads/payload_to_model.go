package payloads

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// UserModel represents the user
// This does not copy over session IDs to the model
func UserModel(user *adminmessages.UserUpdatePayload, id uuid.UUID) (*models.User, *validate.Errors) {
	verrs := validate.NewErrors()

	if user == nil {
		verrs.Add("User", "payload is nil") // does this make sense
		return nil, verrs
	}
	model := &models.User{
		ID: uuid.FromStringOrNil(id.String()),
	}

	if user.Active != nil {
		model.Active = *user.Active
	}

	return model, nil
}

// WebhookSubscriptionModel converts a webhook subscription payload to a model
func WebhookSubscriptionModel(sub *adminmessages.WebhookSubscription) *models.WebhookSubscription {
	model := &models.WebhookSubscription{
		ID: uuid.FromStringOrNil(sub.ID.String()),
	}

	if sub.Severity != nil {
		model.Severity = int(*sub.Severity)
	}

	if sub.CallbackURL != nil {
		model.CallbackURL = *sub.CallbackURL
	}

	if sub.EventKey != nil {
		model.EventKey = *sub.EventKey
	}

	if sub.Status != nil {
		model.Status = models.WebhookSubscriptionStatus(*sub.Status)
	}

	if sub.SubscriberID != nil {
		model.SubscriberID = uuid.FromStringOrNil(sub.SubscriberID.String())
	}

	return model
}

// WebhookSubscriptionModelFromCreate converts a payload for creating a webhook subscription to a model
func WebhookSubscriptionModelFromCreate(sub *adminmessages.CreateWebhookSubscription) *models.WebhookSubscription {
	model := &models.WebhookSubscription{
		// EventKey and CallbackURL are required fields in the YAML, so we don't have to worry about the potential
		// nil dereference errors here:
		EventKey:     *sub.EventKey,
		CallbackURL:  *sub.CallbackURL,
		SubscriberID: uuid.FromStringOrNil(sub.SubscriberID.String()),
	}
	if sub.Status != nil {
		model.Status = models.WebhookSubscriptionStatus(*sub.Status)
	}
	return model
}
