package notifications

import (
	"context"
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type NotificationSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

type testNotification struct {
	email emailContent
}

func (n testNotification) emails() ([]emailContent, error) {
	return []emailContent{n.email}, nil
}

func (suite *NotificationSuite) TestMoveApproved() {
	ctx := context.Background()
	t := suite.T()

	approver := testdatagen.MakeDefaultUser(suite.db)
	move := testdatagen.MakeDefaultMove(suite.db)
	notification := MoveApproved{
		db:     suite.db,
		logger: suite.logger,
		moveID: move.ID,
		session: &auth.Session{
			UserID:          approver.ID,
			ApplicationName: auth.OfficeApp,
		},
	}

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

func (suite *NotificationSuite) TestMoveSubmitted() {
	ctx := context.Background()
	t := suite.T()

	move := testdatagen.MakeDefaultMove(suite.db)
	notification := MoveSubmitted{
		db:     suite.db,
		logger: suite.logger,
		moveID: move.ID,
		session: &auth.Session{
			ServiceMemberID: move.Orders.ServiceMember.ID,
			ApplicationName: auth.MyApp,
		},
	}

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

	logger, _ := zap.NewDevelopment()

	s := &NotificationSuite{
		db:     db,
		logger: logger,
	}
	suite.Run(t, s)
}
