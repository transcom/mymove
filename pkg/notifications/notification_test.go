package notifications

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
	userServices "github.com/transcom/mymove/pkg/services/user"
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

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

func serviceMemberService(db *pop.Connection, l *zap.Logger) services.FetchServiceMember {
	smDB := models.NewServiceMemberDB(db)
	return userServices.NewFetchServiceMemberService(smDB)
}

func (suite *NotificationSuite) TestMoveApproved() {
	t := suite.T()

	approver := testdatagen.MakeDefaultUser(suite.db)
	move := testdatagen.MakeDefaultMove(suite.db)
	session := &server.Session{
		UserID:          approver.ID,
		ApplicationName: server.OfficeApp,
	}
	notification := NewMoveApproved(suite.db,
		suite.logger,
		session,
		serviceMemberService(suite.db, suite.logger),
		move.ID)

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

func (suite *NotificationSuite) TestMoveSubmitted() {
	t := suite.T()

	move := testdatagen.MakeDefaultMove(suite.db)
	session := &server.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: server.MyApp,
	}
	notification := NewMoveSubmitted(suite.db,
		suite.logger,
		session,
		serviceMemberService(suite.db, suite.logger),
		move.ID)

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
