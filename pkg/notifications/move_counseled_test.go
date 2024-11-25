package notifications

import (
	"regexp"
	"strings"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *NotificationSuite) TestMoveCounseled() {
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveCounseled(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Your counselor has approved your move details"

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
}

func (suite *NotificationSuite) TestMoveCounseledHTMLTemplateRender() {
	approver := factory.BuildUser(nil, nil, nil)
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveCounseled(move.ID)

	originDutyLocation := "origDutyLocation"

	s := MoveCounseledEmailData{
		OriginDutyLocation:         &originDutyLocation,
		DestinationLocation:        "destDutyLocation",
		Locator:                    "abc123",
		MyMoveLink:                 MyMoveLink,
		ActualExpenseReimbursement: true,
	}

	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>
<p>This is a confirmation that your counselor has approved move details for the assigned move code abc123 from origDutyLocation to destDutyLocation in the MilMove system.</p>
<p>What this means to you:</br>
If you are doing a Personally Procured Move (PPM), you can start moving your personal property.</p>
<p><strong>Next steps for a PPM:</strong>
<ul>
  <li>Please Note: Your PPM has been designated as Actual Expense Reimbursement. This is the standard entitlement for Civilian employees. For uniformed Service Members, your PPM may have been designated as Actual Expense Reimbursement due to failure to receive authorization prior to movement or failure to obtain certified weight tickets. Actual Expense Reimbursement means reimbursement for expenses not to exceed the Government Constructed Cost (GCC).</li>
  <li>Remember to get legible certified weight tickets for both the empty and full weights for every trip you perform. If you do not upload legible certified weight tickets, your PPM incentive (or Actual Expense Reimbursement for Civilians) could be affected. Failure to obtain weight tickets will result in losing eligibility to receive your incentive.</li>
<p>Note: To receive allowance for Pro-Gear, you must identify allowable items and provide weight tickets separately for Pro-Gear.</p>
  <li>For authorized storage:</li>
    <ul>
      <li>You will need to get weight ticket(s) for the items you store.</li>
      <li>Storage costs cannot be paid in advance.</li>
    </ul>
  <li>If your counselor approved an Advance Operating Allowance (AOA, or cash advance) for a PPM, log into <a href="https://my.move.mil/">MilMove</a> to download your AOA Packet, and submit it to finance according to the instructions provided by your counselor. If you have been directed to use your government travel charge card (GTCC) for expenses no further action is required.</li>
  <li>Once you complete your PPM, log into <a href="https://my.move.mil/">MilMove</a>, upload your receipts and weight tickets, and submit your PPM for review.</li>
</ul>
<p><strong>Next steps for government arranged shipments:</strong></br>
<ul>
  <li>Your move request will be reviewed by the responsible personal property shipping office and a move task order for services will be placed with HomeSafe Alliance.</li>
  <li>Once this order is placed, you will receive an e-mail invitation to create an account in HomeSafe Connect (check your spam or junk folder). This is the system you will use to schedule your pre-move survey.</li>
  <li>HomeSafe is required to contact you within one Government Business Day. Once contact has been established, HomeSafe is your primary point of contact. If any information about your move changes at any point during the move, immediately notify your HomeSafe Customer Care Representative of the changes. Remember to keep your contact information updated in MilMove.</li>
</ul>
<p>Thank you,<br>
USTRANSCOM MilMove Team</p>
<p>The information contained in this email may contain Privacy Act information and is therefore protected under the Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.</p>`

	htmlContent, err := notification.RenderHTML(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}), s)

	suite.NoError(err)
	suite.Equal(trimExtraSpaces(expectedHTMLContent), trimExtraSpaces(htmlContent))
}

func (suite *NotificationSuite) TestMoveCounseledTextTemplateRender() {

	approver := factory.BuildUser(nil, nil, nil)
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveCounseled(move.ID)

	originDutyLocation := "origDutyLocation"

	s := MoveCounseledEmailData{
		OriginDutyLocation:         &originDutyLocation,
		DestinationLocation:        "destDutyLocation",
		Locator:                    "abc123",
		MyMoveLink:                 MyMoveLink,
		ActualExpenseReimbursement: false,
	}

	expectedTextContent := `*** DO NOT REPLY directly to this email ***
This is a confirmation that your counselor has approved move details for the assigned move code abc123 from origDutyLocation to destDutyLocation in the MilMove system.
What this means to you:
If you are doing a Personally Procured Move (PPM), you can start moving your personal property.
Next steps for a PPM:
    * Remember to get legible certified weight tickets for both the empty and full weights for every trip you perform. If you do not upload legible certified weight tickets, your PPM incentive (or Actual Expense Reimbursement for Civilians) could be affected. Failure to obtain weight tickets will result in losing eligibility to receive your incentive.
Note: To receive allowance for Pro-Gear, you must identify allowable items and provide weight tickets separately for Pro-Gear.
    * For authorized storage:
        * You will need to get weight ticket(s) for the items you store.
        * Storage costs cannot be paid in advance.
    * If your counselor approved an Advance Operating Allowance (AOA, or cash advance) for a PPM, log into MilMove <https://my.move.mil/> to download your AOA Packet, and submit it to finance according to the instructions provided by your counselor. If you have been directed to use your government travel charge card (GTCC) for expenses no further action is required.
    * Once you complete your PPM, log into MilMove <https://my.move.mil/>, upload your receipts and weight tickets, and submit your PPM for review.
Next steps for government arranged shipments:
    * Your move request will be reviewed by the responsible personal property shipping office and a move task order for services will be placed with HomeSafe Alliance.
    * Once this order is placed, you will receive an e-mail invitation to create an account in HomeSafe Connect (check your spam or junk folder). This is the system you will use to schedule your pre-move survey.
    * HomeSafe is required to contact you within one Government Business Day. Once contact has been established, HomeSafe is your primary point of contact. If any information about your move changes at any point during the move, immediately notify your HomeSafe Customer Care Representative of the changes. Remember to keep your contact information updated in MilMove.
Thank you,
USTRANSCOM MilMove Team
The information contained in this email may contain Privacy Act information and is therefore protected under the Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.`

	textContent, err := notification.RenderText(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}), s)

	suite.NoError(err)
	suite.Equal(trimExtraSpaces(expectedTextContent), trimExtraSpaces(textContent))
}

func trimExtraSpaces(input string) string {
	// Replace consecutive white spaces with a single space
	re := regexp.MustCompile(`\s+`)
	// return the result without leading or trailing spaces
	return strings.TrimSpace(re.ReplaceAllString(input, " "))
}

func (suite *NotificationSuite) TestCounselorApprovedMoveForSeparatee() {
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypeSEPARATION,
			},
		},
	}, nil)
	notification := NewMoveCounseled(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Your counselor has approved your move details"

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
	suite.Contains(email.textBody, move.MTOShipments[0].DestinationAddress.StreetAddress1)
	suite.Contains(email.textBody, *move.MTOShipments[0].DestinationAddress.StreetAddress2)
	suite.Contains(email.textBody, *move.MTOShipments[0].DestinationAddress.StreetAddress3)
	suite.Contains(email.textBody, move.MTOShipments[0].DestinationAddress.City)
	suite.Contains(email.textBody, move.MTOShipments[0].DestinationAddress.State)
	suite.Contains(email.textBody, move.MTOShipments[0].DestinationAddress.PostalCode)
}

func (suite *NotificationSuite) TestCounselorApprovedMoveForRetiree() {
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypeRETIREMENT,
			},
		},
	}, nil)
	notification := NewMoveCounseled(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Your counselor has approved your move details"

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
	suite.Contains(email.textBody, move.MTOShipments[0].DestinationAddress.StreetAddress1)
	suite.Contains(email.textBody, *move.MTOShipments[0].DestinationAddress.StreetAddress2)
	suite.Contains(email.textBody, *move.MTOShipments[0].DestinationAddress.StreetAddress3)
	suite.Contains(email.textBody, move.MTOShipments[0].DestinationAddress.City)
	suite.Contains(email.textBody, move.MTOShipments[0].DestinationAddress.State)
	suite.Contains(email.textBody, move.MTOShipments[0].DestinationAddress.PostalCode)
}
