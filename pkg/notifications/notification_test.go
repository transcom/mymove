package notifications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type NotificationSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *NotificationSuite) SetupTest() {
	suite.DB().TruncateAll()
}

type testNotification struct {
	email emailContent
}

func (n testNotification) emails() ([]emailContent, error) {
	return []emailContent{n.email}, nil
}

func (suite *NotificationSuite) TestMoveSubmitted() {
	ctx := context.Background()
	t := suite.T()

	move := testdatagen.MakeDefaultMove(suite.DB())
	notification := NewMoveSubmitted(suite.DB(), suite.logger, &auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}, move.ID)

	emails, err := notification.emails(ctx)
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

func (suite *NotificationSuite) getTestEmailContent() emailContent {
	return emailContent{
		recipientEmail: "lucky@winner.com",
		subject:        "This is a Test",
		htmlBody:       "Congrats!<br>You win!",
		textBody:       "Congrats! You win!",
	}
}

func TestNotificationSuite(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	ns := &NotificationSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, ns)
	ns.PopTestSuite.TearDown()
}
