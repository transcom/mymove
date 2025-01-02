package notifications

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gofrs/uuid"
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
	CloseoutQueue(appCtx appcontext.AppContext, queueUrl string) error
}

// NotificationSendingContext provides context to a notification sender. Maps use queueUrl
type NotificationReceiverContext struct {
	snsService           *sns.Client
	sqsService           *sqs.Client
	awsRegion            string
	awsAccountId         string
	queueSubscriptionMap map[string]string
	receiverCancelMap    map[string]context.CancelFunc
}

// NewNotificationSender returns a new NotificationSendingContext
func NewNotificationReceiver(snsService *sns.Client, sqsService *sqs.Client, awsRegion string, awsAccountId string) NotificationReceiverContext {
	return NotificationReceiverContext{
		snsService:           snsService,
		sqsService:           sqsService,
		awsRegion:            awsRegion,
		awsAccountId:         awsAccountId,
		queueSubscriptionMap: make(map[string]string),
		receiverCancelMap:    make(map[string]context.CancelFunc),
	}
}

func (n NotificationReceiverContext) CreateQueueWithSubscription(appCtx appcontext.AppContext, topicName string, params NotificationQueueParams) (string, error) {

	queueUUID := uuid.Must(uuid.NewV4())

	queueName := fmt.Sprintf("%s_%s", params.Action, queueUUID)
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
	subscribeOutput, err := n.snsService.Subscribe(context.Background(), subscribeInput)
	if err != nil {
		log.Fatalf("Failed to create subscription, %v", err)
	}

	n.queueSubscriptionMap[*result.QueueUrl] = *subscribeOutput.SubscriptionArn

	return *result.QueueUrl, err
}

// SendNotification sends a one or more notifications for all supported mediums
func (n NotificationReceiverContext) ReceiveMessages(appCtx appcontext.AppContext, queueUrl string) ([]types.Message, error) {
	recCtx, cancelRecCtx := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelRecCtx()
	n.receiverCancelMap[queueUrl] = cancelRecCtx

	result, err := n.sqsService.ReceiveMessage(recCtx, &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20,
	})
	if err != nil && recCtx.Err() != context.Canceled {
		appCtx.Logger().Info("Couldn't get messages from queue. Here's why: %v\n", zap.Error(err))
		return nil, err
	}

	if recCtx.Err() == context.Canceled {
		return nil, recCtx.Err()
	}

	return result.Messages, recCtx.Err()
}

// map of queueUrl to context

// CloseoutQueue stops receiving messages and cleans up the queue and its subscriptions
func (n NotificationReceiverContext) CloseoutQueue(appCtx appcontext.AppContext, queueUrl string) error {
	appCtx.Logger().Info("CLOSING OUT QUEUE CONTEXT")

	if cancelFunc, exists := n.receiverCancelMap[queueUrl]; exists {
		cancelFunc()
		delete(n.receiverCancelMap, queueUrl)
	}
	if subscriptionArn, exists := n.queueSubscriptionMap[queueUrl]; exists {
		_, err := n.snsService.Unsubscribe(context.Background(), &sns.UnsubscribeInput{
			SubscriptionArn: &subscriptionArn,
		})
		if err != nil {
			return err
		}
		delete(n.queueSubscriptionMap, queueUrl)
	}

	_, err := n.sqsService.DeleteQueue(context.Background(), &sqs.DeleteQueueInput{
		QueueUrl: &queueUrl,
	})
	if err != nil {
		return err
	}

	return nil
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
	}

	snsService := sns.NewFromConfig(cfg)
	sqsService := sqs.NewFromConfig(cfg)

	return NewNotificationReceiver(snsService, sqsService, awsSESRegion, awsAccountId), nil
}

func (n NotificationReceiverContext) constructArn(awsService string, endpointName string) string {
	return fmt.Sprintf("arn:aws-us-gov:%s:%s:%s:%s", awsService, n.awsRegion, n.awsAccountId, endpointName)
}
