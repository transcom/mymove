package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/transcom/mymove/cmd/webhook-client/utils"
	"github.com/transcom/mymove/pkg/models"
)

// Message is the body of our request
type Message struct {
	ID              uuid.UUID       `json:"id"`
	EventName       string          `json:"eventName"`
	TriggeredAt     strfmt.DateTime `json:"triggeredAt"`
	ObjectType      string          `json:"objectType"`
	UpdatedObjectID uuid.UUID       `json:"updatedObjectID"`
	Object          string          `json:"object"`
}

// Engine encapsulates the services used by the webhook notification engine
type Engine struct {
	DB                  *pop.Connection
	Logger              utils.Logger
	Client              utils.WebhookClientPoster
	Cmd                 *cobra.Command
	PeriodInSeconds     int
	MaxImmediateRetries int
}

// processNotifications reads all the notifications and all the subscriptions and processes them one by one
func (eng *Engine) processNotifications(notifications []models.WebhookNotification, subscriptions []models.WebhookSubscription) {
	for _, notif := range notifications {
		notif := notif

		// search for subscription
		foundSub := false
		for _, sub := range subscriptions {
			if sub.EventKey == notif.EventKey {
				foundSub = true
				// If found, send  to subscription
				// #nosec G601 TODO needs review
				err := eng.sendOneNotification(&notif, &sub)
				if err != nil {
					eng.Logger.Error("Webhook Notification send failed", zap.Error(err))
					return
				}
				continue
			}
		}
		if foundSub == false {
			//If no subscription was found, update notification status to skipped.
			eng.Logger.Debug("No subscription found for notification event, skipping.", zap.String("eventKey", notif.EventKey))
			notif.Status = models.WebhookNotificationSkipped
			err := eng.updateNotification(&notif)
			if err != nil {
				eng.Logger.Error("Notification update failed", zap.Error(err))
			}
		}

	}
}

// updateNotification is a helper function to write the notification to the db.
func (eng *Engine) updateNotification(notif *models.WebhookNotification) error {
	// Update notification (status is updated by sendOneNotification)
	verrs, err := eng.DB.ValidateAndUpdate(notif)
	if verrs != nil && verrs.HasAny() {
		err = errors.New(verrs.Error())
		eng.Logger.Error(err.Error())
		return err
	}
	if err != nil {
		eng.Logger.Error(err.Error())
		return err
	}
	return nil
}

// sendOneNotification sends the notification and marks the notification as sent or failed
func (eng *Engine) sendOneNotification(notif *models.WebhookNotification, sub *models.WebhookSubscription) error {
	logger := eng.Logger

	// Construct notification to send
	// message := fmt.Sprintf("There was a %s notification, id: %s, moveTaskOrderID: %s",
	// 	string(notif.EventKey), notif.ID.String(), notif.MoveTaskOrderID.String())
	message := &Message{
		ID:          notif.ID,
		EventName:   notif.EventKey,
		TriggeredAt: strfmt.DateTime(notif.UpdatedAt),
	}
	if notif.Payload != nil {
		message.Object = *notif.Payload
	}
	if notif.ObjectID != nil {
		message.UpdatedObjectID = *notif.ObjectID
	}
	json, err := json.Marshal(message)
	if err != nil {
		notif.Status = models.WebhookNotificationFailed
		eng.updateNotification(notif)
		logger.Error("Error creating payload:", zap.Error(err))
		return err
	}

	// Try MaxImmediateRetries times to send
	try := 0
	for try = 0; try < eng.MaxImmediateRetries; try++ {
		// Post the notification
		url := sub.CallbackURL
		resp, body, err2 := eng.Client.Post(json, url)

		if notif.Status == models.WebhookNotificationPending {
			notif.FirstAttemptedAt = time.Now()
			// Not writing to db, but should be written within
			// this function.
		}

		if err2 == nil && resp.StatusCode == 200 {
			// Update notification
			notif.Status = models.WebhookNotificationSent
			eng.updateNotification(notif)
			logger.Info("Notification successfully sent:",
				zap.String("Status", resp.Status),
				zap.String("EventName", message.EventName),
				zap.String("NotificationID", message.ID.String()),
				zap.String("UpdatedObjectID", message.UpdatedObjectID.String()),
			)
			break // send was successful
		}
		// If there was an error sending, log error and continue
		if err2 != nil {
			logger.Error("Failed to send, error sending webhook:", zap.Error(err2),
				zap.String("notificationID", notif.ID.String()),
				zap.Int("Retry #", try))
			continue
		}
		// If there was an error response from server, log error and continue
		if resp.StatusCode != 200 {
			errmsg := fmt.Sprintf("Failed to send. Response Status: %s. Body: %s", resp.Status, string(body))
			logger.Error("Received error on sending webhook", zap.String("Error", errmsg),
				zap.String("notificationID", notif.ID.String()),
				zap.Int("Retry #", try))
		}
	}

	// Update Notification with failing if appropriate
	if try == eng.MaxImmediateRetries {
		notif.Status = models.WebhookNotificationFailing
		eng.updateNotification(notif)

		errmsg := fmt.Sprintf("Failed to send notification ID: %s after %d immediate retries", notif.ID, try)
		err = errors.New(errmsg)
		return err
	}

	return nil

}

// Run runs the engine once. This happens periodically and is called by Start()
// It collects all the pending notifications and active subscriptions in the db
// and starts processing them.
// If a new notification or subscription were to be adding during the course of one run
// by the Milmove server, it would only be processed on the next call of run().
func (eng *Engine) run() error {

	logger := eng.Logger
	// Read all notifications
	notifications := []models.WebhookNotification{}
	err := eng.DB.Order("created_at asc").Where("status = ? OR status = ?", models.WebhookNotificationPending, models.WebhookNotificationFailing).All(&notifications)

	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	}
	logger.Debug("Notification Check:", zap.Int("Num notifications found", len(notifications)))

	// If none, return
	if len(notifications) == 0 {
		return nil
	}

	// If there are notifications, get subscriptions
	subscriptions := []models.WebhookSubscription{}
	err = eng.DB.Where("status = ? OR status = ?", models.WebhookSubscriptionStatusActive, models.WebhookSubscriptionStatusFailing).All(&subscriptions)

	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	}
	logger.Debug("Subscription Check!", zap.Int("Num subscriptions found", len(subscriptions)))

	// If none, return
	if len(notifications) == 0 {
		return nil
	}

	// process notifications
	eng.processNotifications(notifications, subscriptions)
	return nil
}

// Start starts the timer for the webhook engine
// The process will run once every period to send pending notifications
// The period is defined in the Engine.PeriodInSeconds
func (eng *Engine) Start() error {

	// Set timer tick
	t := time.Tick(time.Duration(eng.PeriodInSeconds) * time.Second)

	// Run once prior to first wait period
	eng.run()
	// Run on each timer tick
	for range t {
		eng.run()
	}

	return nil
}
