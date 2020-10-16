package notifications

import (
	"context"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *NotificationSuite) TestMoveSubmitted() {
	ctx := context.Background()

	move := testdatagen.MakeDefaultMove(suite.DB())
	notification := NewMoveSubmitted(suite.DB(), suite.logger, &auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}, move.ID)

	emails, err := notification.emails(ctx)
	subject := "[MilMove] You’ve submitted your move details"

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
}

func (suite *NotificationSuite) TestMoveSubmittedHTMLTemplateRender() {
	approver := testdatagen.MakeStubbedUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())
	notification := NewMoveSubmitted(suite.DB(), suite.logger, &auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}, move.ID)

	s := moveSubmittedEmailData{
		Link:                       "https://milmovelocal/downloads/ppm_info_sheet.pdf",
		OriginDutyStation:          "origDutyStation",
		DestinationDutyStation:     "destDutyStation",
		OriginDutyStationPhoneLine: "555-555-5555",
		Locator:                    "abc123",
	}
	expectedHTMLContent := `<p>
  Your move from origDutyStation to destDutyStation has been submitted to your local transportation
  office for review.
</p>

<p>This can take up to 3 business days. The office will email you once your move has been approved.</p>

<p>Your move locator code is abc123. Use this code when communicating with the office about your move.</p>

<p>
  In the meantime, if you have questions or need expedited processing, call the origDutyStation PPPO at
  555-555-5555.
</p>

<p>You can check the status of your move at any time at https://my.move.mil/</p>

<p>
  <strong>Let us know how we’re doing.</strong> <a href="https://milmovelocal/downloads/ppm_info_sheet.pdf">Please take a brief survey</a> and share how well
  we’re handling your move so far at https://milmovelocal/downloads/ppm_info_sheet.pdf.
</p>
`

	htmlContent, err := notification.RenderHTML(s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestMoveSubmittedTextTemplateRender() {

	approver := testdatagen.MakeStubbedUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())
	notification := NewMoveSubmitted(suite.DB(), suite.logger, &auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}, move.ID)

	s := moveSubmittedEmailData{
		Link:                       "https://milmovelocal/downloads/ppm_info_sheet.pdf",
		OriginDutyStation:          "origDutyStation",
		DestinationDutyStation:     "destDutyStation",
		OriginDutyStationPhoneLine: "555-555-5555",
		Locator:                    "abc123",
	}
	expectedTextContent := `Your move from origDutyStation to destDutyStation has been submitted to your local transportation office for review.

This can take up to 3 business days. The office will email you once your move has been approved.

Your move locator code is abc123. Use this code when communicating with the office about your move.

In the meantime, if you have questions or need expedited processing, call the origDutyStation PPPO at 555-555-5555.

You can check the status of your move at any time at https://my.move.mil/

Let us know how we’re doing. Please take a brief survey and share how well we’re handling your move so far at https://milmovelocal/downloads/ppm_info_sheet.pdf.
`

	textContent, err := notification.RenderText(s)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}
