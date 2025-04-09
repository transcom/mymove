package notifications

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/cli"
	mocks "github.com/transcom/mymove/pkg/notifications/receiverMocks"
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

func (suite *notificationReceiverSuite) TestSuccessPath() {

	suite.Run("local backend - notification receiver stub", func() {
		// Setup mocks
		mockedViper := mocks.ViperType{}
		mockedViper.On("GetString", cli.ReceiverBackendFlag).Return("local")
		mockedViper.On("GetString", cli.SNSRegionFlag).Return("us-gov-west-1")
		mockedViper.On("GetString", cli.SNSAccountId).Return("12345")
		mockedViper.On("GetString", cli.SNSTagsUpdatedTopicFlag).Return("fake_sns_topic")
		localReceiver, err := InitReceiver(&mockedViper, suite.Logger(), true)

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

	suite.Run("aws backend - notification receiver InitReceiver", func() {
		// Setup mocks
		mockedViper := mocks.ViperType{}
		mockedViper.On("GetString", cli.ReceiverBackendFlag).Return("sns_sqs")
		mockedViper.On("GetString", cli.SNSRegionFlag).Return("us-gov-west-1")
		mockedViper.On("GetString", cli.SNSAccountId).Return("12345")
		mockedViper.On("GetString", cli.SNSTagsUpdatedTopicFlag).Return("fake_sns_topic")

		os.Unsetenv("AWS_PROFILE")

		receiver, err := InitReceiver(&mockedViper, suite.Logger(), false)

		suite.NoError(err)
		suite.IsType(NotificationReceiverContext{}, receiver)
		defaultTopic, err := receiver.GetDefaultTopic()
		suite.Equal("fake_sns_topic", defaultTopic)
		suite.NoError(err)
	})

	suite.Run("aws backend - notification receiver with mock services", func() {
		// Setup mocks
		mockedViper := mocks.ViperType{}
		mockedViper.On("GetString", cli.ReceiverBackendFlag).Return("sns_sqs")
		mockedViper.On("GetString", cli.SNSRegionFlag).Return("us-gov-west-1")
		mockedViper.On("GetString", cli.SNSAccountId).Return("12345")
		mockedViper.On("GetString", cli.SNSTagsUpdatedTopicFlag).Return("fake_sns_topic")

		mockedSns := mocks.SnsClient{}
		mockedSns.On("Subscribe", mock.Anything, mock.AnythingOfType("*sns.SubscribeInput")).Return(&sns.SubscribeOutput{
			SubscriptionArn: aws.String("FakeSubscriptionArn"),
		}, nil)
		mockedSns.On("Unsubscribe", mock.Anything, mock.AnythingOfType("*sns.UnsubscribeInput")).Return(&sns.UnsubscribeOutput{}, nil)
		mockedSns.On("ListSubscriptionsByTopic", mock.Anything, mock.AnythingOfType("*sns.ListSubscriptionsByTopicInput")).Return(&sns.ListSubscriptionsByTopicOutput{}, nil)

		mockedSqs := mocks.SqsClient{}
		mockedSqs.On("CreateQueue", mock.Anything, mock.AnythingOfType("*sqs.CreateQueueInput")).Return(&sqs.CreateQueueOutput{
			QueueUrl: aws.String("fakeQueueUrl"),
		}, nil)
		mockedSqs.On("ReceiveMessage", mock.Anything, mock.AnythingOfType("*sqs.ReceiveMessageInput")).Return(&sqs.ReceiveMessageOutput{
			Messages: []types.Message{
				{
					MessageId: aws.String("fakeMessageId"),
					Body:      aws.String("fakeQueueUrl:fakeMessageBody"),
				},
			},
		}, nil)
		mockedSqs.On("DeleteMessage", mock.Anything, mock.AnythingOfType("*sqs.DeleteMessageInput")).Return(&sqs.DeleteMessageOutput{}, nil)
		mockedSqs.On("DeleteQueue", mock.Anything, mock.AnythingOfType("*sqs.DeleteQueueInput")).Return(&sqs.DeleteQueueOutput{}, nil)
		mockedSqs.On("ListQueues", mock.Anything, mock.AnythingOfType("*sqs.ListQueuesInput")).Return(&sqs.ListQueuesOutput{}, nil)

		// Run test
		receiver := NewNotificationReceiver(&mockedViper, &mockedSns, &mockedSqs, "", "")
		suite.IsType(NotificationReceiverContext{}, receiver)

		defaultTopic, err := receiver.GetDefaultTopic()
		suite.Equal("fake_sns_topic", defaultTopic)
		suite.NoError(err)

		queueParams := NotificationQueueParams{
			NamePrefix: "testPrefix",
		}
		createdQueueUrl, err := receiver.CreateQueueWithSubscription(suite.AppContextForTest(), queueParams)
		suite.NoError(err)
		suite.Equal("fakeQueueUrl", createdQueueUrl)

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
