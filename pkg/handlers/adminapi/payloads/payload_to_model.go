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
		CallbackURL:  sub.CallbackURL,
		EventKey:     sub.EventKey,
		ID:           uuid.FromStringOrNil(sub.ID.String()),
		SubscriberID: uuid.FromStringOrNil(sub.SubscriberID.String()),
	}

	if sub.Severity != nil {
		model.Severity = int(*sub.Severity)
	}

	return model
}
