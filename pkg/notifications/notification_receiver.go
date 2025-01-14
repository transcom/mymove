package notifications

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
)

// NotificationQueueParams stores the params for queue creation
type NotificationQueueParams struct {
	SubscriptionTopicName string
	NamePrefix            QueuePrefixType
	FilterPolicy          string
}

// NotificationReceiver is an interface for receiving notifications
type NotificationReceiver interface {
	CreateQueueWithSubscription(appCtx appcontext.AppContext, params NotificationQueueParams) (string, error)
	ReceiveMessages(appCtx appcontext.AppContext, queueUrl string, timerContext context.Context) ([]ReceivedMessage, error)
	CloseoutQueue(appCtx appcontext.AppContext, queueUrl string) error
	GetDefaultTopic() (string, error)
}

// NotificationReceiverConext provides context to a notification Receiver. Maps use queueUrl for key
type NotificationReceiverContext struct {
	viper                ViperType
	snsService           SnsClient
	sqsService           SqsClient
	awsRegion            string
	awsAccountId         string
	queueSubscriptionMap map[string]string
	receiverCancelMap    map[string]context.CancelFunc
}

// QueuePrefixType represents a prefix identifier given to a name of dynamic notification queues
type QueuePrefixType string

const (
	QueuePrefixObjectTagsAdded QueuePrefixType = "ObjectTagsAdded"
)

//go:generate mockery --name SnsClient --output ./receiverMocks
type SnsClient interface {
	Subscribe(ctx context.Context, params *sns.SubscribeInput, optFns ...func(*sns.Options)) (*sns.SubscribeOutput, error)
	Unsubscribe(ctx context.Context, params *sns.UnsubscribeInput, optFns ...func(*sns.Options)) (*sns.UnsubscribeOutput, error)
	ListSubscriptionsByTopic(context.Context, *sns.ListSubscriptionsByTopicInput, ...func(*sns.Options)) (*sns.ListSubscriptionsByTopicOutput, error)
}

//go:generate mockery --name SqsClient --output ./receiverMocks
type SqsClient interface {
	CreateQueue(ctx context.Context, params *sqs.CreateQueueInput, optFns ...func(*sqs.Options)) (*sqs.CreateQueueOutput, error)
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteQueue(ctx context.Context, params *sqs.DeleteQueueInput, optFns ...func(*sqs.Options)) (*sqs.DeleteQueueOutput, error)
	ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error)
}

//go:generate mockery --name ViperType --output ./receiverMocks
type ViperType interface {
	GetString(string) string
	SetEnvKeyReplacer(*strings.Replacer)
}

// ReceivedMessage standardizes the format of the received message
type ReceivedMessage struct {
	MessageId string
	Body      *string
}

// NewNotificationReceiver returns a new NotificationReceiverContext
func NewNotificationReceiver(v ViperType, snsService SnsClient, sqsService SqsClient, awsRegion string, awsAccountId string) NotificationReceiverContext {
	return NotificationReceiverContext{
		viper:                v,
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
		appCtx.Logger().Error("Failed to create SQS queue, %v", zap.Error(err))
		return "", err
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
		appCtx.Logger().Error("Failed to create subscription, %v", zap.Error(err))
		return "", err
	}

	n.queueSubscriptionMap[*result.QueueUrl] = *subscribeOutput.SubscriptionArn

	return *result.QueueUrl, nil
}

// ReceiveMessages polls given queue continuously for messages for up to 20 seconds
func (n NotificationReceiverContext) ReceiveMessages(appCtx appcontext.AppContext, queueUrl string, timerContext context.Context) ([]ReceivedMessage, error) {
	recCtx, cancelRecCtx := context.WithCancel(timerContext)
	defer cancelRecCtx()
	n.receiverCancelMap[queueUrl] = cancelRecCtx

	result, err := n.sqsService.ReceiveMessage(recCtx, &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20,
	})
	if errors.Is(recCtx.Err(), context.Canceled) || errors.Is(recCtx.Err(), context.DeadlineExceeded) {
		return nil, recCtx.Err()
	}

	if err != nil {
		appCtx.Logger().Info("Couldn't get messages from queue. Error: %v\n", zap.Error(err))
		return nil, err
	}

	receivedMessages := make([]ReceivedMessage, len(result.Messages))
	for index, value := range result.Messages {
		receivedMessages[index] = ReceivedMessage{
			MessageId: *value.MessageId,
			Body:      value.Body,
		}
	}

	return receivedMessages, recCtx.Err()
}

// CloseoutQueue stops receiving messages and cleans up the queue and its subscriptions
func (n NotificationReceiverContext) CloseoutQueue(appCtx appcontext.AppContext, queueUrl string) error {
	appCtx.Logger().Info("Closing out queue: %v", zap.String("queueUrl", queueUrl))

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

	return err
}

// GetDefaultTopic returns the topic value set within the environment
func (n NotificationReceiverContext) GetDefaultTopic() (string, error) {
	topicName := n.viper.GetString(cli.SNSTagsUpdatedTopicFlag)
	receiverBackend := n.viper.GetString(cli.ReceiverBackendFlag)
	if topicName == "" && receiverBackend == "sns&sqs" {
		return "", errors.New("sns_tags_updated_topic key not available")
	}
	return topicName, nil
}

// InitReceiver initializes the receiver backend, only call this once
func InitReceiver(v ViperType, logger *zap.Logger, wipeAllNotificationQueues bool) (NotificationReceiver, error) {

	if v.GetString(cli.ReceiverBackendFlag) == "sns&sqs" {
		// Setup notification receiver service with SNS & SQS backend dependencies
		awsSNSRegion := v.GetString(cli.SNSRegionFlag)
		awsAccountId := v.GetString(cli.SNSAccountId)

		logger.Info("Using aws sns&sqs receiver backend", zap.String("region", awsSNSRegion))

		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(awsSNSRegion),
		)
		if err != nil {
			logger.Fatal("error loading sns aws config", zap.Error(err))
			return nil, err
		}

		snsService := sns.NewFromConfig(cfg)
		sqsService := sqs.NewFromConfig(cfg)

		notificationReceiver := NewNotificationReceiver(v, snsService, sqsService, awsSNSRegion, awsAccountId)

		// Remove any remaining previous notification queues on server start
		if wipeAllNotificationQueues {
			err = notificationReceiver.wipeAllNotificationQueues(logger)
			if err != nil {
				return nil, err
			}
		}

		return notificationReceiver, nil
	}

	return NewStubNotificationReceiver(), nil
}

func (n NotificationReceiverContext) constructArn(awsService string, endpointName string) string {
	return fmt.Sprintf("arn:aws-us-gov:%s:%s:%s:%s", awsService, n.awsRegion, n.awsAccountId, endpointName)
}

// Removes ALL previously created notification queues
func (n *NotificationReceiverContext) wipeAllNotificationQueues(logger *zap.Logger) error {
	defaultTopic, err := n.GetDefaultTopic()
	if err != nil {
		return err
	}

	logger.Info("Removing previous subscriptions...")
	paginator := sns.NewListSubscriptionsByTopicPaginator(n.snsService, &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(n.constructArn("sns", defaultTopic)),
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return err
		}
		for _, subscription := range output.Subscriptions {
			if strings.Contains(*subscription.Endpoint, string(QueuePrefixObjectTagsAdded)) {
				logger.Info("Subscription ARN: ", zap.String("subscription arn", *subscription.SubscriptionArn))
				logger.Info("Endpoint ARN: ", zap.String("endpoint arn", *subscription.Endpoint))
				_, err = n.snsService.Unsubscribe(context.Background(), &sns.UnsubscribeInput{
					SubscriptionArn: subscription.SubscriptionArn,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	logger.Info("Removing previous queues...")
	result, err := n.sqsService.ListQueues(context.Background(), &sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(string(QueuePrefixObjectTagsAdded)),
	})
	if err != nil {
		return err
	}

	for _, url := range result.QueueUrls {
		_, err = n.sqsService.DeleteQueue(context.Background(), &sqs.DeleteQueueInput{
			QueueUrl: &url,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
