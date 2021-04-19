package webhook

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/transcom/mymove/cmd/webhook-client/utils"
	"github.com/transcom/mymove/pkg/models"
)

// Engine encapsulates the services used by the webhook notification engine
type Engine struct {
	DB                  *pop.Connection
	Logger              utils.Logger
	Client              utils.WebhookClientPoster
	Cmd                 *cobra.Command
	PeriodInSeconds     int
	MaxImmediateRetries int
	SeverityThresholds  []int
	QuitChannel         chan os.Signal
	DoneChannel         chan bool
}

// processNotifications reads all the notifications and all the subscriptions and processes them one by one
func (eng *Engine) processNotifications(notifications []models.WebhookNotification, subscriptions []models.WebhookSubscription) {
	for _, notif := range notifications {
		notif := notif

		// search for subscription
		foundSub := false
		for _, sub := range subscriptions {
			sub := sub
			if sub.EventKey == notif.EventKey {
				foundSub = true
				stopLoop := false
				var sev int
				// If found, send  to subscription
				err := eng.sendOneNotification(&notif, &sub)
				// If notification send failed, we need to log the severity
				if err != nil {
					eng.Logger.Error("Webhook Notification send failed", zap.Error(err))
					if notif.FirstAttemptedAt == nil {
						eng.Logger.Error("FirstAttempted at time was not stored", zap.Error(err))
						// We should not ever get this error, so we trigger a sev1 failure immediately
						sev = 1
						eng.Logger.Error("Raising severity of failure",
							zap.String("subscriptionEvent", sub.EventKey),
							zap.Int("severityFrom", sub.Severity),
							zap.Int("severityTo", sev))
						notif.Status = models.WebhookNotificationFailed
						err = eng.updateNotification(&notif)
						if err != nil {
							eng.Logger.Error("Webhook Notification update failed", zap.Error(err))
						}
					} else {
						sev = eng.GetSeverity(time.Now(), *notif.FirstAttemptedAt)
						if sev != sub.Severity {
							eng.Logger.Error("Raising severity of failure",
								zap.String("subscriptionEvent", sub.EventKey),
								zap.Int("severityFrom", sub.Severity),
								zap.Int("severityTo", sev))
							if sev == 1 {
								notif.Status = models.WebhookNotificationFailed
								err = eng.updateNotification(&notif)
								if err != nil {
									eng.Logger.Error("Webhook Notification update failed", zap.Error(err))
								}
							}
						}
					}
					stopLoop = true
				}
				// Update subscription, needs to be done on success sometimes, hence it's out of the previous if
				errDB := eng.updateSubscriptionStatus(&notif, &sub, sev)
				if errDB != nil {
					eng.Logger.Error("Webhook Subscription update failed", zap.Error(err))
				}

				if stopLoop {
					return
				}

				// Return out of loop if quit signal recieved, otherwise, keep going
				select {
				case <-eng.QuitChannel:
					eng.Logger.Info("Interrupt signal recieved...")
					eng.DoneChannel <- true
					return
				default:
				}
			}
		}
		if foundSub == false {
			//If no subscription was found, update notification status to skipped.
			eng.Logger.Info("No subscription found for notification event, skipping.", zap.String("eventKey", notif.EventKey))
			notif.Status = models.WebhookNotificationSkipped
			err := eng.updateNotification(&notif)
			if err != nil {
				eng.Logger.Error("Notification update failed", zap.Error(err))
			}
		}
	}
}

// updateSubscriptionStatus updates the subscription based on the status of the last notification.
// Returns nil if nothing to update or update succeeds, returns error if error found
func (eng *Engine) updateSubscriptionStatus(notif *models.WebhookNotification, sub *models.WebhookSubscription,
	newSeverity int) error {
	// Update subscription status if it has changed

	var doUpdate = false
	switch notif.Status {
	case models.WebhookNotificationFailed:
		// If the notification is set to failed, then we need to deactivate the subscription
		if sub.Status != models.WebhookSubscriptionStatusDisabled {
			sub.Status = models.WebhookSubscriptionStatusDisabled
			sub.Severity = newSeverity
			doUpdate = true
		}
	case models.WebhookNotificationFailing:
		// If the notification is failing, then we may need to update the status and/or severity
		if sub.Status != models.WebhookSubscriptionStatusFailing || sub.Severity != newSeverity {
			sub.Status = models.WebhookSubscriptionStatusFailing
			sub.Severity = newSeverity
			doUpdate = true
		}
	case models.WebhookNotificationSent:
		// If the notification sent, we may need to recover the status and severity back to a-ok
		if sub.Status != models.WebhookSubscriptionStatusActive || sub.Severity != 0 {
			sub.Status = models.WebhookSubscriptionStatusActive
			sub.Severity = 0
			doUpdate = true
		}
	}

	if doUpdate {
		verrs, err := eng.DB.ValidateAndUpdate(sub)
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
	return nil
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

// sendOneNotification sends the notification the max immediate retries. It updates the notification's status
// and stores it in the DB.
func (eng *Engine) sendOneNotification(notif *models.WebhookNotification, sub *models.WebhookSubscription) error {
	logger := eng.Logger

	// Construct notification to send
	message := GetWebhookNotificationPayload(notif)
	json, err := message.MarshalBinary()
	if err != nil {
		notif.Status = models.WebhookNotificationFailed
		updateNotificationErr := eng.updateNotification(notif)
		if updateNotificationErr != nil {
			eng.Logger.Error("Notification update failed", zap.Error(err))
		}
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
			now := time.Now()
			notif.FirstAttemptedAt = &now
			// Not writing to db, but should be written within
			// this function.
		}

		if err2 == nil && resp.StatusCode == 200 {
			// Update notification
			notif.Status = models.WebhookNotificationSent
			updateNotificationErr := eng.updateNotification(notif)
			if updateNotificationErr != nil {
				eng.Logger.Error("Notification update failed", zap.Error(err))
			}
			objectID := "<empty>"
			if message.ObjectID != nil {
				objectID = message.ObjectID.String()
			}
			mtoID := "<empty>"
			if message.MoveTaskOrderID != nil {
				mtoID = message.MoveTaskOrderID.String()
			}
			logger.Info("Notification successfully sent:",
				zap.String("id", message.ID.String()),
				zap.String("status", resp.Status),
				zap.String("eventKey", message.EventKey),
				zap.String("moveID", mtoID),
				zap.String("objectID", objectID),
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
			logger.Error("Received error on sending notification",
				zap.String("Response Status", resp.Status),
				zap.String("Response Body", string(body)),
				zap.String("notificationID", notif.ID.String()),
				zap.Int("Retry #", try))
		}
	}

	// Update Notification with failing if appropriate
	if try == eng.MaxImmediateRetries {
		notif.Status = models.WebhookNotificationFailing
		updateNotificationErr := eng.updateNotification(notif)
		if updateNotificationErr != nil {
			eng.Logger.Error("Notification update failed", zap.Error(err))
		}

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
	logger.Info("Notification Check:", zap.Int("Num notifications found", len(notifications)))

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
	logger.Info("Subscription Check!", zap.Int("Num subscriptions found", len(subscriptions)))

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

	logger := eng.Logger
	logger.Info("Starting engine", zap.Int("periodInSeconds", eng.PeriodInSeconds),
		zap.Int("maxImmediateRetries", eng.MaxImmediateRetries),
		zap.Any("SeverityThresholds", eng.SeverityThresholds))

	// Set timer tick
	t := time.Tick(time.Duration(eng.PeriodInSeconds) * time.Second)

	// Run once prior to first wait period
	err := eng.run()
	if err != nil {
		return err
	}

	// Run on each timer tick
	for range t {
		select {
		case <-eng.QuitChannel:
			eng.Logger.Info("Interrupt signal recieved...")
			eng.DoneChannel <- true
		default:
			err = eng.run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetSeverity is a function that returns the severity level of a single attempt given an array of severity thresholds
func (eng *Engine) GetSeverity(currentTime time.Time, firstAttempt time.Time) int {
	timeSinceFirstAttempt := int(currentTime.Sub(firstAttempt).Seconds())
	levels := len(eng.SeverityThresholds) + 1
	sev := 1 //if the loop condition is not met, then the severity is 1
	for index, threshold := range eng.SeverityThresholds {
		if timeSinceFirstAttempt < threshold {
			sev = levels - index
			break
		}
	}
	return sev
}
