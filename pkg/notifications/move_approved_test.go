package notifications

import (
	"context"
	"fmt"
	"strings"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *NotificationSuite) TestMoveApproved() {
	ctx := context.Background()

	approver := testdatagen.MakeStubbedUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())

	notification := NewMoveApproved(suite.DB(), suite.logger, &auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}, "milmovelocal", move.ID)
	subject := fmt.Sprintf("[MilMove] Your Move is approved (move: %s)", move.Locator)

	emails, err := notification.emails(ctx)
	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
	suite.True(strings.Contains(email.textBody, notification.host))
}

func (suite *NotificationSuite) TestMoveApprovedHTMLTemplateRender() {
	approver := testdatagen.MakeStubbedUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())
	notification := NewMoveApproved(suite.DB(), suite.logger, &auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}, "milmovelocal", move.ID)

	s := moveApprovedEmailData{
		Link:                       "https://milmovelocal/downloads/ppm_info_sheet.pdf",
		OriginDutyStation:          "origDutyStation",
		DestinationDutyStation:     "destDutyStation",
		OriginDutyStationPhoneLine: "555-555-5555",
		Locator:                    "abc123",
	}
	expectedHTMLContent := `<p><strong>You're all set to move!</strong></p>

<p>
  The local transportation office <strong>approved your move</strong> from <strong>origDutyStation</strong> to
  <strong>destDutyStation</strong
  >.
</p>

<p>Please <a href="https://milmovelocal/downloads/ppm_info_sheet.pdf">review the Personally Procured Move (PPM) info sheet</a> for detailed instructions.</p>
<br />
<p>
  <strong>Next steps</strong> <br />Because you’ve chosen a do-it-yourself move, you can start whenever you are ready.
</p>

<p>
  Be sure to <strong>save your weight tickets and any receipts</strong> associated with your move. You’ll need them to
  request payment later in the process.
</p>

<p>If you have any questions, call the origDutyStation PPPO at 555-555-5555 and reference your move locator code: abc123</p>

<p>You can <a href="https://my.move.mil">check the status of your move</a> anytime at https://my.move.mil"</p>
`

	htmlContent, err := notification.RenderHTML(s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestMoveApprovedTextTemplateRender() {

	approver := testdatagen.MakeStubbedUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())
	notification := NewMoveApproved(suite.DB(), suite.logger, &auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}, "milmovelocal", move.ID)

	s := moveApprovedEmailData{
		Link:                       "https://milmovelocal/downloads/ppm_info_sheet.pdf",
		OriginDutyStation:          "origDutyStation",
		DestinationDutyStation:     "destDutyStation",
		OriginDutyStationPhoneLine: "555-555-5555",
		Locator:                    "abc123",
	}
	expectedTextContent := `You're all set to move!

The local transportation office approved your move from origDutyStation to destDutyStation.

Please review the Personally Procured Move (PPM) info sheet for detailed instructions at https://milmovelocal/downloads/ppm_info_sheet.pdf.


Next steps
Because you’ve chosen a do-it-yourself move, you can start whenever you are ready.

Be sure to save your weight tickets and any receipts associated with your move. You’ll need them to request payment later in the process.

If you have any questions, call the origDutyStation PPPO at 555-555-5555 and reference move locator code: abc123.

You can check the status of your move anytime at https://my.move.mil"
`

	textContent, err := notification.RenderText(s)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}
