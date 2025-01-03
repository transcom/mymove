package notifications

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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
	return "fakeQueueName", nil
}

func (n StubNotificationReceiver) ReceiveMessages(appCtx appcontext.AppContext, queueUrl string) ([]types.Message, error) {
	// TODO: sleep func here & swap types for return to avoid aws type
	messageId := "fakeMessageId"
	body := queueUrl + ":fakeMessageBody"
	mockMessages := make([]types.Message, 1)
	mockMessages = append(mockMessages, types.Message{
		MessageId: &messageId,
		Body:      &body,
	})
	appCtx.Logger().Debug("Receiving a fake message for queue: %v", zap.String("queueUrl", queueUrl))
	return mockMessages, nil
}

func (n StubNotificationReceiver) CloseoutQueue(appCtx appcontext.AppContext, queueUrl string) error {
	appCtx.Logger().Debug("Closing out the fake queue.")
	return nil
}
