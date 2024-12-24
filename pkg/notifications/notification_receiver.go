package notifications

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
)

// Notification is an interface for creating emails
type NotificationFilter struct {
}

// NotificationSender is an interface for sending notifications
//
//go:generate mockery --name NotificationSender
type NotificationReceiver interface {
	SubscribeToTopic(appCtx appcontext.AppContext, filter NotificationFilter) error
}

// NotificationSendingContext provides context to a notification sender
type NotificationReceiverContext struct {
	svc *sqs.Client
}

// NewNotificationSender returns a new NotificationSendingContext
func NewNotificationReceiver(svc *sqs.Client) NotificationReceiverContext {
	return NotificationReceiverContext{
		svc: svc,
	}
}

// SendNotification sends a one or more notifications for all supported mediums
func (n NotificationReceiverContext) SubscribeToTopic(appCtx appcontext.AppContext, filter NotificationFilter) error {
	queueRaw := "testQueue"
	queue := &queueRaw

	urlResult, _ := n.svc.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: queue,
	})

	result, err := n.svc.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl:            urlResult.QueueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     5,
	})
	if err != nil {
		appCtx.Logger().Fatal("Couldn't get messages from queue. Here's why: %v\n", zap.Error(err))
	} else {
		for _, val := range result.Messages {
			appCtx.Logger().Info(*val.MessageId)
		}
	}
	return err
}

// InitEmail initializes the email backend
func InitReceiver(v *viper.Viper, logger *zap.Logger) (NotificationReceiver, error) {
	// if v.GetString(cli.EmailBackendFlag) == "ses" {
	// 	// Setup Amazon SES (email) service TODO: This might be able
	// 	// to be combined with the AWS Session that we're using for S3
	// 	// down below.

	// 	awsSESRegion := v.GetString(cli.AWSSESRegionFlag)
	// 	awsSESDomain := v.GetString(cli.AWSSESDomainFlag)
	// 	sysAdminEmail := v.GetString(cli.SysAdminEmail)
	// 	logger.Info("Using ses email backend",
	// 		zap.String("region", awsSESRegion),
	// 		zap.String("domain", awsSESDomain))
	// 	cfg, err := config.LoadDefaultConfig(context.Background(),
	// 		config.WithRegion(awsSESRegion),
	// 	)
	// 	if err != nil {
	// 		logger.Fatal("error loading ses aws config", zap.Error(err))
	// 	}

	// 	sesService := ses.NewFromConfig(cfg)
	// 	input := &ses.GetAccountSendingEnabledInput{}
	// 	result, err := sesService.GetAccountSendingEnabled(context.Background(), input)
	// 	if err != nil || result == nil || !result.Enabled {
	// 		logger.Error("email sending not enabled", zap.Error(err))
	// 		return NewNotificationSender(nil, awsSESDomain, sysAdminEmail), err
	// 	}
	// 	return NewNotificationSender(sesService, awsSESDomain, sysAdminEmail), nil
	// }

	// domain := "milmovelocal"
	// logger.Info("Using local email backend", zap.String("domain", domain))
	// return NewStubNotificationSender(domain), nil

	// Setup Amazon SES (email) service TODO: This might be able
	// to be combined with the AWS Session that we're using for S3
	// down below.

	awsSESRegion := v.GetString(cli.AWSSESRegionFlag)

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(awsSESRegion),
	)
	if err != nil {
		logger.Fatal("error loading ses aws config", zap.Error(err))
	}

	sqsService := sqs.NewFromConfig(cfg)

	return NewNotificationReceiver(sqsService), nil
}
