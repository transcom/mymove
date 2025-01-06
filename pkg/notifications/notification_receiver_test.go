package notifications

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

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

func (suite *notificationReceiverSuite) TestNotificationReceiverLocalStub() {
	v := viper.New()
	localReceiver, err := InitReceiver(v, suite.Logger())

	suite.NoError(err)

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
}
