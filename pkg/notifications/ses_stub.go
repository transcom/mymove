package notifications

import (
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

// StubSESClient mocks an SES client for local usage
type StubSESClient struct {
	sesiface.SESAPI
}

// SendRawEmail returns a dummy ID
func (m StubSESClient) SendRawEmail(input *ses.SendRawEmailInput) (*ses.SendRawEmailOutput, error) {
	notAnID := "twelve"
	output := ses.SendRawEmailOutput{
		MessageId: &notAnID,
	}

	return &output, nil
}
