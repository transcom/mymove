package payloads

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// UserModel represents the user
// This does not copy over session IDs to the model
func UserModel(user *adminmessages.UserUpdate, id uuid.UUID, userOriginalActive bool) (*models.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user payload is nil")
	}
	model := &models.User{
		ID: uuid.FromStringOrNil(id.String()),
	}

	if user.OktaEmail != nil {
		model.OktaEmail = *user.OktaEmail
	}

	if user.Active == nil { // active status was nil in payload
		model.Active = userOriginalActive
	} else { // active status was provided in payload
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

func OfficeUserModelFromUpdate(payload *adminmessages.OfficeUserUpdate, officeUser *models.OfficeUser) *models.OfficeUser {
	if payload == nil || officeUser == nil {
		return officeUser
	}
	if payload.Email != nil {
		officeUser.Email = *payload.Email
	}

	if payload.FirstName != nil {
		officeUser.FirstName = *payload.FirstName
	}

	if payload.MiddleInitials != nil {
		officeUser.MiddleInitials = payload.MiddleInitials
	}

	if payload.LastName != nil {
		officeUser.LastName = *payload.LastName
	}

	if payload.Telephone != nil {
		officeUser.Telephone = *payload.Telephone
	}

	if payload.Active != nil {
		officeUser.Active = *payload.Active
	}
	return officeUser
}
