package notifications

import (
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

var member = primemessages.Customer{Email: "test@example.com"}
var primeOrder = primemessages.Order{
	OriginDutyLocation:      &primemessages.DutyLocation{Name: "Fort Origin"},
	DestinationDutyLocation: &primemessages.DutyLocation{Name: "Fort Destination"},
	Customer:                &member,
}
var payload = primemessages.MoveTaskOrder{
	MoveCode: "TEST00",
	Order:    &primeOrder,
}
var correctPrimeCounselingData = PrimeCounselingCompleteData{
	CustomerEmail:                     member.Email,
	Locator:                           payload.MoveCode,
	OriginDutyLocation:                primeOrder.OriginDutyLocation.Name,
	DestinationDutyLocation:           primeOrder.DestinationDutyLocation.Name,
	OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
	MyMoveLink:                        MyMoveLink,
}

func (suite *NotificationSuite) TestPrimeCounselingComplete() {
	notification := NewPrimeCounselingComplete(payload)

	primeCounselingEmailData, err := notification.GetEmailData(notification.moveTaskOrder, suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(primeCounselingEmailData)
	suite.Equal(primeCounselingEmailData, correctPrimeCounselingData)

	suite.EqualExportedValues(primeCounselingEmailData, PrimeCounselingCompleteData{
		CustomerEmail:                     member.Email,
		OriginDutyLocation:                primeOrder.OriginDutyLocation.Name,
		DestinationDutyLocation:           primeOrder.DestinationDutyLocation.Name,
		Locator:                           payload.MoveCode,
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
		MyMoveLink:                        MyMoveLink,
	})

	expectedHTMLContent := getCorrectEmailTemplate(primeCounselingEmailData)

	htmlContent, err := notification.RenderHTML(suite.AppContextForTest(), primeCounselingEmailData)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}

func (suite *NotificationSuite) TestPrimeCounselingCompleteTextTemplateRender() {
	notification := NewPrimeCounselingComplete(payload)

	primeCounselingEmailData, err := notification.GetEmailData(notification.moveTaskOrder, suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(primeCounselingEmailData)
	suite.Equal(primeCounselingEmailData, correctPrimeCounselingData)

	suite.EqualExportedValues(primeCounselingEmailData, PrimeCounselingCompleteData{
		CustomerEmail:                     member.Email,
		OriginDutyLocation:                primeOrder.OriginDutyLocation.Name,
		DestinationDutyLocation:           primeOrder.DestinationDutyLocation.Name,
		Locator:                           payload.MoveCode,
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
		MyMoveLink:                        MyMoveLink,
	})

	expectedTextContent := getCorrectTextTemplate(primeCounselingEmailData)

	textContent, err := notification.RenderText(suite.AppContextForTest(), primeCounselingEmailData)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}

func getCorrectEmailTemplate(emailData PrimeCounselingCompleteData) string {
	return `<p>*** DO NOT REPLY directly to this email ***</p>
<p>This is a confirmation that your counselor has approved move details for the <strong>assigned move code ` + emailData.Locator + `</strong> from <strong>` + emailData.OriginDutyLocation + `</strong> to <strong>` + emailData.DestinationDutyLocation + `</strong> in the MilMove system.</p>
<p>What this means to you:</p>
<p>If you are doing a Personally Procured Move (PPM), you can start moving your personal property.</p>
<h4>Next steps for a PPM:</h4>
<ul>
<li>Remember to get legible certified weight tickets for both the empty and full weights for every trip you perform.  If you do not upload legible certified weight tickets, your PPM incentive could be affected.</li>

<li>If you are requesting an Advance Operating Allowance (AOA, or cash advance) for a PPM, log into <a href="` + MyMoveLink + `">MilMove</a> to download your AOA packet. You must obtain signature approval on the AOA packet from a government transportation office before submitting it to finance. If you have been directed to use your government travel charge card (GTCC) for expenses no further action is required.</li>

<li>If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href="` + OneSourceTransportationOfficeLink + `">` + OneSourceTransportationOfficeLink + `</a></li>

<li>Once you complete your PPM, log into <a href="` + MyMoveLink + `">MilMove</a>, upload your receipts and weight tickets, and submit your PPM for review.</li>
</ul>
<h4>Next steps for government arranged shipments:</h4>
<ul>
<li>If additional services were identified during counseling, HomeSafe will send the request to the responsible government transportation office for review. Your HomeSafe Customer Care Representative should keep you informed on the status of the request.</li>

<li>If you have not already done so, please schedule a pre-move survey using HomeSafe Connect or by contacting a HomeSafe Customer Care Representative.</li>

<li>HomeSafe is your primary point of contact. If any information changes during the move, immediately notify your HomeSafe Customer Care Representative of the changes. Remember to keep your contact information updated in MilMove.</li>
</ul>
<p>If you are unsatisfied at any time, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href="` + OneSourceTransportationOfficeLink + `">` + OneSourceTransportationOfficeLink + `</a></p>
<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>The information contained in this email may contain Privacy Act information and is therefore protected under the Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.</p>`
}

func getCorrectTextTemplate(emailData PrimeCounselingCompleteData) string {
	return `*** DO NOT REPLY directly to this email ***
This is a confirmation that your counselor has approved move details for the assigned move code ` + emailData.Locator + ` from ` + emailData.OriginDutyLocation + ` to ` + emailData.DestinationDutyLocation + ` in the MilMove system.

What this means to you:
If you are doing a Personally Procured Move (PPM), you can start moving your personal property.

Next steps for a PPM:
• Remember to get legible certified weight tickets for both the empty and full weights for every trip you perform.  If you do not upload legible certified weight tickets, your PPM incentive could be affected.

• If you are requesting an Advance Operating Allowance (AOA, or cash advance) for a PPM, log into MilMove <` + MyMoveLink + `/> to download your AOA packet. You must obtain signature approval on the AOA packet from a government transportation office before submitting it to finance. If you have been directed to use your government travel charge card (GTCC) for expenses no further action is required.

• If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <` + OneSourceTransportationOfficeLink + `>

• Once you complete your PPM, log into MilMove <` + MyMoveLink + `/>, upload your receipts and weight tickets, and submit your PPM for review.

Next steps for government arranged shipments:
• If additional services were identified during counseling, HomeSafe will send the request to the responsible government transportation office for review. Your HomeSafe Customer Care Representative should keep you informed on the status of the request.

• If you have not already done so, please schedule a pre-move survey using HomeSafe Connect or by contacting a HomeSafe Customer Care Representative.

• HomeSafe is your primary point of contact. If any information changes during the move, immediately notify your HomeSafe Customer Care Representative of the changes. Remember to keep your contact information updated in MilMove.

If you are unsatisfied at any time, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: ` + OneSourceTransportationOfficeLink + `.

Thank you,

USTRANSCOM MilMove Team

The information contained in this email may contain Privacy Act information and is therefore protected under the Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.`
}
