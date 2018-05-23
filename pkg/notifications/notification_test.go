package notifications

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type NotificationSuite struct {
	suite.Suite
	db *pop.Connection
}

type mockSESClient struct {
	sesiface.SESAPI
	mock.Mock
	Suite *NotificationSuite
}

func (m *mockSESClient) SendRawEmail(input *ses.SendRawEmailInput) (*ses.SendRawEmailOutput, error) {
	args := m.Called(input)

	testEmail := m.Suite.GetTestEmailContent()
	m.Suite.Equal(testEmail.recipientEmail, *input.Destinations[0])
	m.Suite.Equal(senderEmail(), *input.Source)

	message := string(input.RawMessage.Data)
	m.Suite.Contains(message, testEmail.subject)
	m.Suite.Contains(message, testEmail.htmlBody)
	m.Suite.Contains(message, testEmail.textBody)
	m.Suite.Contains(message, testEmail.recipientEmail)
	m.Suite.Contains(message, senderEmail())

	return args.Get(0).(*ses.SendRawEmailOutput), args.Error(1)
}

type testNotification struct {
	email emailContent
}

func (n testNotification) emails() ([]emailContent, error) {
	return []emailContent{n.email}, nil
}

func (suite *NotificationSuite) TestMoveApproved() {
	t := suite.T()

	approver, _ := testdatagen.MakeUser(suite.db)
	move, _ := testdatagen.MakeMove(suite.db)

	notification := MoveApproved{
		db:     suite.db,
		moveID: move.ID,
		reqApp: app.OfficeApp,
		user:   approver,
	}

	emails, err := notification.emails()
	if err != nil {
		t.Fatal(err)
	}

	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.NotEmpty(email.subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
}

func (suite *NotificationSuite) TestSendNotification() {
	t := suite.T()

	messageID := "a"
	mockSVC := mockSESClient{Suite: suite}
	mockSVC.On("SendRawEmail", mock.Anything).Return(&ses.SendRawEmailOutput{MessageId: &messageID}, nil)

	err := SendNotification(testNotification{email: suite.GetTestEmailContent()}, &mockSVC)
	if err != nil {
		t.Fatal(err)
	}

	mockSVC.AssertNumberOfCalls(t, "SendRawEmail", 1)
}

func (suite *NotificationSuite) GetTestEmailContent() emailContent {
	return emailContent{
		recipientEmail: "lucky@winner.com",
		subject:        "This is a Test",
		htmlBody:       "Congrats!<br>You win!",
		textBody:       "Congrats! You win!",
	}
}

func TestNotificationSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	s := &NotificationSuite{db: db}
	suite.Run(t, s)
}
