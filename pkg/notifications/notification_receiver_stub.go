package notifications

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

// StubNotificationReceiver mocks an SNS & SQS client for local usage
type StubNotificationReceiver NotificationReceiverContext

// NewStubNotificationReceiver returns a new StubNotificationReceiver
func NewStubNotificationReceiver() StubNotificationReceiver {
	return StubNotificationReceiver{
		snsService:           nil,
		sqsService:           nil,
		awsRegion:            "",
		awsAccountId:         "",
		queueSubscriptionMap: make(map[string]string),
		receiverCancelMap:    make(map[string]context.CancelFunc),
	}
}

func (n StubNotificationReceiver) CreateQueueWithSubscription(appCtx appcontext.AppContext, params NotificationQueueParams) (string, error) {
	return "stubQueueName", nil
}

func (n StubNotificationReceiver) ReceiveMessages(appCtx appcontext.AppContext, queueUrl string, timerContext context.Context) ([]ReceivedMessage, error) {
	time.Sleep(2 * time.Second)
	messageId := "stubMessageId"
	body := queueUrl + ":stubMessageBody"
	mockMessages := make([]ReceivedMessage, 1)
	mockMessages[0] = ReceivedMessage{
		MessageId: messageId,
		Body:      &body,
	}
	appCtx.Logger().Debug("Receiving a stubbed message for queue: %v", zap.String("queueUrl", queueUrl))
	return mockMessages, nil
}

func (n StubNotificationReceiver) CloseoutQueue(appCtx appcontext.AppContext, queueUrl string) error {
	appCtx.Logger().Debug("Closing out the stubbed queue.")
	return nil
}

func (n StubNotificationReceiver) GetDefaultTopic() (string, error) {
	return "stubDefaultTopic", nil
}
