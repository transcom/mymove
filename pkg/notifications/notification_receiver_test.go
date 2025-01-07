package notifications

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

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

type Viper struct {
	mock.Mock
}

func (_m *Viper) GetString(key string) string {
	switch key {
	case cli.ReceiverBackendFlag:
		return "sns&sqs"
	case cli.AWSRegionFlag:
		return "us-gov-west-1"
	case cli.AWSSNSAccountId:
		return "12345"
	case cli.AWSSNSObjectTagsAddedTopicFlag:
		return "fake_sns_topic"
	}
	return ""
}

func (_m *Viper) SetEnvKeyReplacer(_ *strings.Replacer) {}

func (suite *notificationReceiverSuite) TestSuccessPath() {

	suite.Run("local backend - notification receiver stub", func() {
		v := viper.New()
		localReceiver, err := InitReceiver(v, suite.Logger())

		suite.NoError(err)
		suite.IsType(StubNotificationReceiver{}, localReceiver)

		queueParams := NotificationQueueParams{
			NamePrefix: "testPrefix",
		}
		createdQueueUrl, err := localReceiver.CreateQueueWithSubscription(suite.AppContextForTest(), queueParams)
		suite.NoError(err)
		suite.NotContains(createdQueueUrl, queueParams.NamePrefix)
		suite.Equal(createdQueueUrl, "stubQueueName")

		receivedMessages, err := localReceiver.ReceiveMessages(suite.AppContextForTest(), createdQueueUrl)
		suite.NoError(err)
		suite.Len(receivedMessages, 1)
		suite.Equal(receivedMessages[0].MessageId, "stubMessageId")
		suite.Equal(*receivedMessages[0].Body, fmt.Sprintf("%s:stubMessageBody", createdQueueUrl))
	})

	suite.Run("aws backend - notification receiver", func() {
		v := Viper{}

		rec, _ := InitReceiver(&v, suite.Logger())
		suite.IsType(NotificationReceiverContext{}, rec)
		defaultTopic, err := rec.GetDefaultTopic()
		suite.Logger().Error("%s", zap.String("default topic", defaultTopic))
		suite.Equal("fake_sns_topic", defaultTopic)
		suite.NoError(err)

		// queueParams := NotificationQueueParams{
		// 	NamePrefix: "testPrefix",
		// }
		// createdQueueUrl, err := localReceiver.CreateQueueWithSubscription(suite.AppContextForTest(), queueParams)
		// suite.NoError(err)
		// suite.NotContains(createdQueueUrl, queueParams.NamePrefix)
		// suite.Equal(createdQueueUrl, "stubQueueName")

		// receivedMessages, err := localReceiver.ReceiveMessages(suite.AppContextForTest(), createdQueueUrl)
		// suite.NoError(err)
		// suite.Len(receivedMessages, 1)
		// suite.Equal(receivedMessages[0].MessageId, "stubMessageId")
		// suite.Equal(*receivedMessages[0].Body, fmt.Sprintf("%s:stubMessageBody", createdQueueUrl))
	})

}

// func (suite *notificationReceiverSuite) TestNotificationReceiverInitReceiver

// func (suite *notificationReceiverSuite) TestNotificationReceiverAWS() {
// 	v := viper.New()
// 	v.Set(cli.ReceiverBackendFlag, "sns&sqs")
// 	v.Set(cli.AWSSNSRegionFlag, "us-gov-west-1")
// 	v.Set(cli.AWSSNSAccountId, "12345")

// 	awsReceiver, err := InitReceiver(v, suite.Logger())
// }
