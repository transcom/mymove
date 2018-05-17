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
	suite *NotificationSuite
}

func (m *mockSESClient) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	args := m.Called(input)

	m.suite.NotEmpty(input.Destination.ToAddresses)
	m.suite.NotEmpty(input.Message.Subject.Data)
	m.suite.NotEmpty(input.Message.Body.Html.Data)
	m.suite.NotEmpty(input.Message.Body.Text.Data)

	return args.Get(0).(*ses.SendEmailOutput), args.Error(1)
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

	approver, _ := testdatagen.MakeUser(suite.db)
	move, _ := testdatagen.MakeMove(suite.db)

	notification := MoveApproved{
		db:     suite.db,
		moveID: move.ID,
		reqApp: app.OfficeApp,
		user:   approver,
	}

	messageID := "a"
	mockSVC := mockSESClient{suite: suite}
	mockSVC.On("SendEmail", mock.Anything).Return(&ses.SendEmailOutput{MessageId: &messageID}, nil)

	err := SendNotification(notification, &mockSVC)
	if err != nil {
		t.Fatal(err)
	}

	mockSVC.AssertNumberOfCalls(t, "SendEmail", 1)
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
