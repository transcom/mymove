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
	subject := "Thank you for submitting your move details"

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
		Link:                       "https://my.move.mil/",
		PpmLink:                    "https://office.move.mil/downloads/ppm_info_sheet.pdf",
		OriginDutyStation:          "origDutyStation",
		DestinationDutyStation:     "destDutyStation",
		OriginDutyStationPhoneLine: "555-555-5555",
		Locator:                    "abc123",
		// WeightAllowance:	entitlement total_weight_self
		// also fix references below
	}
	expectedHTMLContent := `<p>
		This is a confirmation that you’ve submitted the details for your move from origDutyStation to
		destDutyStation.
	</p>

  <p>
    <strong>We’ve assigned you a move code: abc123.</strong> You can use this code when talking to any
    representative about your move.
  </p>

  <p>
    To change any other information about your move, or to add or cancel shipments, you should tell your movers (if you’re
    using them) or your move counselor.
  </p>

  <p>
    <strong>Your weight allowance: {{.WeightAllowance}} pounds</strong>
    That’s how much the government will pay to ship for you on this move. You won’t owe anything if all your shipments
    combined weigh less than that.
  </p>

  <p>
    If you move more than {{.WeightAllowance}} pounds, you will owe the government the difference between that and the
    total amount you move.
  </p>

  <p>
    Your movers will estimate the total weight of your belongings, and you will be notified if it looks like you might
    exceed your weight allowance. But you’re ultimately responsible for the weight you move.
  </p>

  <p>
    <strong>For PPM (DITY, or do-it-yourself) shipments</strong>
    If you chose to do a full or partial PPM (DITY) move,
    <a href="https://office.move.mil/downloads/ppm_info_sheet.pdf"> review the Personally Procured Move (PPM) info sheet</a>
    for detailed instructions.
  </p>
  <ul>
    <li>Start your PPM shipment whenever you are ready</li>
    <li>You can wait until after you talk to your move counselor</li>
    <li>Getting your PPM shipment moved to your new home is entirely in your hands</li>
    <li>You can move everything yourself, hire help, or even hire your own movers</li>
    <li>You are responsible for any damage to your belongings</li>
    <li>
      <strong>Get certified weight tickets</strong> that show the empty and full weight for each vehicle used in each PPM
      shipment
    </li>
    <li>If you’re missing weight tickets, you may not get paid for your PPM</li>
    <li>
      <strong>Save receipts</strong> for PPM expenses to request reimbursement or to reduce taxes on your incentive
      payment
    </li>
  </ul>

  <p>
    If you have any questions about the PPM part of your move, call the origDutyStation PPPO at
    555-555-5555 and reference move code abc123.
  </p>

  <p>
    Once you’ve completed your PPM shipment, you can request payment by
    <a href="https://my.move.mil/">logging in to MilMove</a>.
  </p>

  <p>
    <strong>For HHG and other government-funded shipments</strong>
  </p>

  <p>Next steps:</p>
  <ul>
    <li>Talk to a move counselor</li>
    <li>Talk to your movers</li>
  </ul>

  <p>
    You can ask questions of your move counselor or your point of contact with the movers. They will both get in touch
    with you soon.
  </p>

  <p>Your counselor will:</p>
  <ul>
    <li>Verify the information you entered</li>
    <li>Give you moving-related advice</li>
    <li>Give you tips to avoid going over your weight allowance</li>
    <li>Identify things like pro-gear that won’t count against your weight allowance</li>
  </ul>

  <p>When the movers contact you, they’ll schedule a pre-move survey to estimate the total weight of your belongings.</p>

  <p>They’ll also finalize dates to pack and pick up your things, on or near the date you requested in MilMove.</p>

  <p>If any information about your move changes at any point during the move, let your movers know.</p>

  <p>Good luck on your move to destDutyStation!</p>
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
