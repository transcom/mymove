package notifications

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
)

// Notification is an interface for creating emails
type NotificationQueueParams struct {
	// TODO: change to enum
	Action   string
	ObjectId string
}

// NotificationSender is an interface for sending notifications
//
//go:generate mockery --name NotificationSender
type NotificationReceiver interface {
	CreateQueueWithSubscription(appCtx appcontext.AppContext, topicArn string, params NotificationQueueParams) (string, error)
	ReceiveMessages(appCtx appcontext.AppContext, queueUrl string) ([]types.Message, error)
}

// NotificationSendingContext provides context to a notification sender
type NotificationReceiverContext struct {
	snsService   *sns.Client
	sqsService   *sqs.Client
	awsRegion    string
	awsAccountId string
}

// NewNotificationSender returns a new NotificationSendingContext
func NewNotificationReceiver(snsService *sns.Client, sqsService *sqs.Client, awsRegion string, awsAccountId string) NotificationReceiverContext {
	return NotificationReceiverContext{
		snsService:   snsService,
		sqsService:   sqsService,
		awsRegion:    awsRegion,
		awsAccountId: awsAccountId,
	}
}

func (n NotificationReceiverContext) CreateQueueWithSubscription(appCtx appcontext.AppContext, topicName string, params NotificationQueueParams) (string, error) {

	queueName := fmt.Sprintf("%s_%s", params.Action, params.ObjectId)
	queueArn := n.constructArn("sqs", queueName)
	topicArn := n.constructArn("sns", topicName)

	// Create queue

	accessPolicy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Sid": "AllowSNSPublish",
			"Effect": "Allow",
			"Principal": {
				"Service": "sns.amazonaws.com"
			},
			"Action": ["sqs:SendMessage"],
			"Resource": "%s",
			"Condition": {
				"ArnEquals": {
				"aws:SourceArn": "%s"
				}
      		}
		}]
	}`, queueArn, topicArn)

	input := &sqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: map[string]string{
			"MessageRetentionPeriod": "120",
			"Policy":                 accessPolicy,
		},
	}

	result, err := n.sqsService.CreateQueue(context.Background(), input)
	if err != nil {
		log.Fatalf("Failed to create SQS queue, %v", err)
	}

	// Create subscription

	filterPolicy := fmt.Sprintf(`{
		"detail": {
				"object": {
					"key": [
						{"suffix": "%s"}
					]
				}
			}
	}`, params.ObjectId)

	subscribeInput := &sns.SubscribeInput{
		TopicArn: &topicArn,
		Protocol: aws.String("sqs"),
		Endpoint: &queueArn,
		Attributes: map[string]string{
			"FilterPolicy":      filterPolicy,
			"FilterPolicyScope": "MessageBody",
		},
	}
	_, err = n.snsService.Subscribe(context.Background(), subscribeInput)
	if err != nil {
		log.Fatalf("Failed to create subscription, %v", err)
	}

	return *result.QueueUrl, err
}

// SendNotification sends a one or more notifications for all supported mediums
func (n NotificationReceiverContext) ReceiveMessages(appCtx appcontext.AppContext, queueUrl string) ([]types.Message, error) {
	result, err := n.sqsService.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		appCtx.Logger().Fatal("Couldn't get messages from queue. Here's why: %v\n", zap.Error(err))
	}
	return result.Messages, err
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

	// TODO: verify if we should change this param name to awsNotificationRegion
	awsSESRegion := v.GetString(cli.AWSSESRegionFlag)
	awsAccountId := v.GetString("aws-account-id")

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(awsSESRegion),
	)
	if err != nil {
		logger.Fatal("error loading ses aws config", zap.Error(err))
		return nil, err
	}

	snsService := sns.NewFromConfig(cfg)
	sqsService := sqs.NewFromConfig(cfg)

	return NewNotificationReceiver(snsService, sqsService, awsSESRegion, awsAccountId), nil
}

func (n NotificationReceiverContext) constructArn(awsService string, endpointName string) string {
	return fmt.Sprintf("arn:aws-us-gov:%s:%s:%s:%s", awsService, n.awsRegion, n.awsAccountId, endpointName)
}
