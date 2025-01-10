package notifications

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type notificationReceiverSuite struct {
	*testingsuite.PopTestSuite
}

func TestNotificationReceiverSuite(t *testing.T) {

	hs := &notificationReceiverSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

// mock - Viper
type Viper struct {
	mock.Mock
}

func (_m *Viper) GetString(key string) string {
	switch key {
	case cli.ReceiverBackendFlag:
		return "sns&sqs"
	case cli.SNSRegionFlag:
		return "us-gov-west-1"
	case cli.SNSAccountId:
		return "12345"
	case cli.SNSTagsUpdatedTopicFlag:
		return "fake_sns_topic"
	}
	return ""
}
func (_m *Viper) SetEnvKeyReplacer(_ *strings.Replacer) {}

// mock - SNS
type MockSnsClient struct {
	mock.Mock
}

func (_m *MockSnsClient) Subscribe(ctx context.Context, params *sns.SubscribeInput, optFns ...func(*sns.Options)) (*sns.SubscribeOutput, error) {
	return &sns.SubscribeOutput{SubscriptionArn: aws.String("FakeSubscriptionArn")}, nil
}

func (_m *MockSnsClient) Unsubscribe(ctx context.Context, params *sns.UnsubscribeInput, optFns ...func(*sns.Options)) (*sns.UnsubscribeOutput, error) {
	return &sns.UnsubscribeOutput{}, nil
}

// mock - SQS
type MockSqsClient struct {
	mock.Mock
}

func (_m *MockSqsClient) CreateQueue(ctx context.Context, params *sqs.CreateQueueInput, optFns ...func(*sqs.Options)) (*sqs.CreateQueueOutput, error) {
	return &sqs.CreateQueueOutput{
		QueueUrl: aws.String("FakeQueueUrl"),
	}, nil
}
func (_m *MockSqsClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	messages := make([]types.Message, 0)
	messages = append(messages, types.Message{
		MessageId: aws.String("fakeMessageId"),
		Body:      aws.String(*params.QueueUrl + ":fakeMessageBody"),
	})
	return &sqs.ReceiveMessageOutput{
		Messages: messages,
	}, nil
}
func (_m *MockSqsClient) DeleteQueue(ctx context.Context, params *sqs.DeleteQueueInput, optFns ...func(*sqs.Options)) (*sqs.DeleteQueueOutput, error) {
	return &sqs.DeleteQueueOutput{}, nil
}

func (suite *notificationReceiverSuite) TestSuccessPath() {

	suite.Run("local backend - notification receiver stub", func() {
		v := viper.New()
		localReceiver, err := InitReceiver(v, suite.Logger())

		suite.NoError(err)
		suite.IsType(StubNotificationReceiver{}, localReceiver)

		defaultTopic, err := localReceiver.GetDefaultTopic()
		suite.Equal("stubDefaultTopic", defaultTopic)
		suite.NoError(err)

		queueParams := NotificationQueueParams{
			NamePrefix: "testPrefix",
		}
		createdQueueUrl, err := localReceiver.CreateQueueWithSubscription(suite.AppContextForTest(), queueParams)
		suite.NoError(err)
		suite.NotContains(createdQueueUrl, queueParams.NamePrefix)
		suite.Equal(createdQueueUrl, "stubQueueName")

		timerContext, cancelTimerContext := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelTimerContext()

		receivedMessages, err := localReceiver.ReceiveMessages(suite.AppContextForTest(), createdQueueUrl, timerContext)
		suite.NoError(err)
		suite.Len(receivedMessages, 1)
		suite.Equal(receivedMessages[0].MessageId, "stubMessageId")
		suite.Equal(*receivedMessages[0].Body, fmt.Sprintf("%s:stubMessageBody", createdQueueUrl))
	})

	suite.Run("aws backend - notification receiver init", func() {
		v := Viper{}

		receiver, _ := InitReceiver(&v, suite.Logger())
		suite.IsType(NotificationReceiverContext{}, receiver)
		defaultTopic, err := receiver.GetDefaultTopic()
		suite.Equal("fake_sns_topic", defaultTopic)
		suite.NoError(err)
	})

	suite.Run("aws backend - notification receiver with mock services", func() {
		v := Viper{}
		snsService := MockSnsClient{}
		sqsService := MockSqsClient{}

		receiver := NewNotificationReceiver(&v, &snsService, &sqsService, "", "")
		suite.IsType(NotificationReceiverContext{}, receiver)

		defaultTopic, err := receiver.GetDefaultTopic()
		suite.Equal("fake_sns_topic", defaultTopic)
		suite.NoError(err)

		queueParams := NotificationQueueParams{
			NamePrefix: "testPrefix",
		}
		createdQueueUrl, err := receiver.CreateQueueWithSubscription(suite.AppContextForTest(), queueParams)
		suite.NoError(err)
		suite.Equal("FakeQueueUrl", createdQueueUrl)

		timerContext, cancelTimerContext := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelTimerContext()

		receivedMessages, err := receiver.ReceiveMessages(suite.AppContextForTest(), createdQueueUrl, timerContext)
		suite.NoError(err)
		suite.Len(receivedMessages, 1)
		suite.Equal(receivedMessages[0].MessageId, "fakeMessageId")
		suite.Equal(*receivedMessages[0].Body, fmt.Sprintf("%s:fakeMessageBody", createdQueueUrl))

		err = receiver.CloseoutQueue(suite.AppContextForTest(), createdQueueUrl)
		suite.NoError(err)
	})
}
