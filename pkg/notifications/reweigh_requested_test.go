package notifications

import (
	"strings"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *NotificationSuite) TestReweighRequestedOnSuccess() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	officeUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
		factory.GetTraitOfficeUserTOO,
	})

	notification := NewReweighRequested(move.ID, shipment)
	subject := "FYI: Your HHG should be reweighed before it is delivered"

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          officeUser.ID,
		ApplicationName: auth.OfficeApp,
	}))
	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
	suite.True(strings.Contains(email.textBody, "MilMove let your movers know that they need to reweigh your HHG shipment before they deliver it to your destination"))
}

func (suite *NotificationSuite) TestReweighRequestedHTMLTemplateRender() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	officeUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
		factory.GetTraitOfficeUserTOO,
	})
	notification := NewReweighRequested(move.ID, shipment)
	s := reweighRequestedEmailData{}
	expectedHTMLContent := `<p><strong>Essential information</strong></p>
<ul>
  <li>MilMove let your movers know that they need to reweigh your HHG shipment before they deliver it to your destination</li>
  <li>Your movers will tell you where and when they’ll perform the reweigh</li>
  <li>You have the right to witness the reweigh, but you do not have to</li>
</ul>

<p><strong>Why was a reweigh requested?</strong></p>
<p>When MilMove notices that you might move more weight than your weight allowance, the system will automatically request a reweigh.</p>
<p>This happens when a shipment’s weight brings you within 10% of your weight allowance.</p>
<p>It’s also possible for a transportation officer to manually request a reweigh for a shipment.</p>
<p>You can always request a reweigh yourself until the shipment has been unloaded.</p>
<p>The advantage for you: The official weight of your shipment will be the lower weight between the original weighing and the reweigh.</p>
<p>Remember, if you move more weight than your allowance, you have to pay for the extra. A reweigh makes sure you end up with the lowest weight total possible.</p>

<p><strong>What do I need to do?</strong></p>
<p>Your movers will let you know <strong>where</strong> and <strong>when</strong> they will reweigh your shipment.</p>
<p>Reweighs typically happen near your destination and close to your delivery date.</p>
<p>You’re entitled to witness the reweigh yourself, or send someone else to do it on your behalf. This is optional. You do not have to attend the reweigh.</p>

<p><strong>Who should I talk to about my reweigh?</strong></p>
<p>Your movers will be arranging the logistics of the reweigh. You should talk to them directly if you have questions or arrange to be present for the reweigh.</p>
<p>If your movers don’t give you a date and time for this reweigh soon, ask them for that info.</p>

<p><strong>Make sure your movers have reweighed your shipment before you accept delivery.</strong></p>
<p>The only reason they can reject a reweigh request is if they get the request after the shipment has been unloaded.</p>
<p>If you believe your shipment should be reweighed and has not been, do not accept delivery. Tell your movers that they need to reweigh the shipment before they unload it.</p>
`

	htmlContent, err := notification.RenderHTML(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          officeUser.ID,
		ApplicationName: auth.OfficeApp,
	}), s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestReweighRequestedTextTemplateRender() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	officeUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
		factory.GetTraitOfficeUserTOO,
	})
	notification := NewReweighRequested(move.ID, shipment)
	s := reweighRequestedEmailData{}
	expectedTextContent := `Essential information
* MilMove let your movers know that they need to reweigh your HHG shipment before they deliver it to your destination
* Your movers will tell you where and when they’ll perform the reweigh
* You have the right to witness the reweigh, but you do not have to

Why was a reweigh requested?
When MilMove notices that you might move more weight than your weight allowance, the system will automatically request a reweigh.

This happens when a shipment’s weight brings you within 10% of your weight allowance.

It’s also possible for a transportation officer to manually request a reweigh for a shipment.

You can always request a reweigh yourself until the shipment has been unloaded.

The advantage for you: The official weight of your shipment will be the lower weight between the original weighing and the reweigh.

Remember, if you move more weight than your allowance, you have to pay for the extra. A reweigh makes sure you end up with the lowest weight total possible.

What do I need to do?
Your movers will let you know where and when they will reweigh your shipment.

Reweighs typically happen near your destination and close to your delivery date.

You’re entitled to witness the reweigh yourself, or send someone else to do it on your behalf. This is optional. You do not have to attend the reweigh.

Who should I talk to about my reweigh?
Your movers will be arranging the logistics of the reweigh. You should talk to them directly if you have questions or arrange to be present for the reweigh.

If your movers don’t give you a date and time for this reweigh soon, ask them for that info.

Make sure your movers have reweighed your shipment before you accept delivery.

The only reason they can reject a reweigh request is if they get the request after the shipment has been unloaded.

If you believe your shipment should be reweighed and has not been, do not accept delivery. Tell your movers that they need to reweigh the shipment before they unload it.
`
	textContent, err := notification.RenderText(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          officeUser.ID,
		ApplicationName: auth.OfficeApp,
	}), s)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}
