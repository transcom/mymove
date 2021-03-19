package main

import (
	"errors"
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/cmd/webhook-client/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func initDbConnectionFlags(flag *pflag.FlagSet) {

	flag.SortFlags = false
}

func notificationCreate(db *pop.Connection, logger utils.Logger) (*models.WebhookNotification, error) {
	// Create a notification model
	notID := uuid.Must(uuid.NewV4())
	message := "{ \"message\" : \"A move task order was created.\" }"
	var notification = models.WebhookNotification{
		ID:       notID,
		EventKey: "MoveTaskOrder.Create",
		Payload:  message,
		Status:   models.WebhookNotificationPending,
	}

	// Save it to the db
	verrs, err := db.ValidateAndCreate(&notification)
	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return nil, err
	} else if verrs != nil && verrs.HasAny() {
		fmt.Println(verrs)
		return nil, errors.New(verrs.String())
	}

	return &notification, nil
}

func dbConnection(cmd *cobra.Command, args []string) error {
	v := viper.New()

	err := utils.ParseFlags(cmd, v, args)
	if err != nil {
		return err
	}

	db, logger, err := InitRootConfig(v)
	if err != nil {
		logger.Error("invalid configuration", zap.Error(err))
		return err
	}

	// Create notification
	notification, err := notificationCreate(db, logger)
	if err != nil {
		return err
	}
	logger.Info("Notification created", zap.String("ID", notification.ID.String()))

	// Update notification
	notification.Status = models.WebhookNotificationSent
	verrs, err := db.ValidateAndUpdate(notification)
	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	} else if verrs != nil && verrs.HasAny() {
		fmt.Println(verrs)
		return errors.New(verrs.String())
	}
	logger.Info("Notification updated", zap.String("ID", notification.ID.String()))

	err = db.Destroy(notification)
	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	}
	logger.Info("Notification deleted", zap.String("ID", notification.ID.String()))

	// Create a webhookSubscription
	subscription := testdatagen.MakeDefaultWebhookSubscription(db)
	verrs, err = db.ValidateAndSave(&subscription)
	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	} else if verrs != nil && verrs.HasAny() {
		fmt.Println(verrs)
		return errors.New(verrs.String())
	}
	logger.Info("Subscription created", zap.String("ID", subscription.ID.String()))

	subscription.Status = models.WebhookSubscriptionStatusDisabled
	verrs, err = db.ValidateAndUpdate(&subscription)
	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	} else if verrs != nil && verrs.HasAny() {
		fmt.Println(verrs)
		return errors.New(verrs.String())
	}
	logger.Info("Subscription updated", zap.String("ID", subscription.ID.String()))

	err = db.Destroy(&subscription)
	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	}
	logger.Info("Subscription deleted", zap.String("ID", subscription.ID.String()))

	return nil
}
