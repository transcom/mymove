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
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
)

// NotificationQueueParams stores the params for queue creation
type NotificationQueueParams struct {
	SubscriptionTopicName string
	NamePrefix            string
	FilterPolicy          string
}

// NotificationReceiver is an interface for receiving notifications
//
//go:generate mockery --name NotificationReceiver
type NotificationReceiver interface {
	CreateQueueWithSubscription(appCtx appcontext.AppContext, params NotificationQueueParams) (string, error)
	ReceiveMessages(appCtx appcontext.AppContext, queueUrl string) ([]types.Message, error)
	CloseoutQueue(appCtx appcontext.AppContext, queueUrl string) error
}

// NotificationReceiverConext provides context to a notification Receiver. Maps use queueUrl for key
type NotificationReceiverContext struct {
	snsService           *sns.Client
	sqsService           *sqs.Client
	awsRegion            string
	awsAccountId         string
	queueSubscriptionMap map[string]string
	receiverCancelMap    map[string]context.CancelFunc
}

// NewNotificationReceiver returns a new NotificationReceiverContext
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

// CreateQueueWithSubscription first creates a new queue, then subscribes an AWS topic to it
func (n NotificationReceiverContext) CreateQueueWithSubscription(appCtx appcontext.AppContext, params NotificationQueueParams) (string, error) {

	queueUUID := uuid.Must(uuid.NewV4())

	queueName := fmt.Sprintf("%s_%s", params.NamePrefix, queueUUID)
	queueArn := n.constructArn("sqs", queueName)
	topicArn := n.constructArn("sns", params.SubscriptionTopicName)

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

	subscribeInput := &sns.SubscribeInput{
		TopicArn: &topicArn,
		Protocol: aws.String("sqs"),
		Endpoint: &queueArn,
		Attributes: map[string]string{
			"FilterPolicy":      params.FilterPolicy,
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

// ReceiveMessages polls given queue continuously for messages for up to 20 seconds
func (n NotificationReceiverContext) ReceiveMessages(appCtx appcontext.AppContext, queueUrl string) ([]types.Message, error) {
	recCtx, cancelRecCtx := context.WithCancel(context.Background())
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

	if v.GetString(cli.ReceiverBackendFlag) == "sns&sqs" {
		// Setup notification receiver service with SNS & SQS backend dependencies
		awsSNSRegion := v.GetString(cli.AWSSNSRegionFlag)
		awsAccountId := v.GetString(cli.AWSSNSAccountId)

		logger.Info("Using aws sns&sqs receiver backend", zap.String("region", awsSNSRegion))

		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(awsSNSRegion),
		)
		if err != nil {
			logger.Fatal("error loading sns aws config", zap.Error(err))
		}

		snsService := sns.NewFromConfig(cfg)
		sqsService := sqs.NewFromConfig(cfg)

		return NewNotificationReceiver(snsService, sqsService, awsSNSRegion, awsAccountId), nil
	}

	return NewStubNotificationReceiver(), nil
}

func (n NotificationReceiverContext) constructArn(awsService string, endpointName string) string {
	return fmt.Sprintf("arn:aws-us-gov:%s:%s:%s:%s", awsService, n.awsRegion, n.awsAccountId, endpointName)
}
