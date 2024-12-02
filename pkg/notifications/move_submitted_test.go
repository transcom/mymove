package notifications

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *NotificationSuite) TestMoveSubmitted() {
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveSubmitted(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
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

func (suite *NotificationSuite) TestMoveSubmittedoriginDSTransportInfoIsNil() {
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveSubmitted(move.ID)

	move.Orders.OriginDutyLocationID = nil

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Thank you for submitting your move details"

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
	suite.Contains(email.textBody, move.Orders.OriginDutyLocation.Name)
}

func (suite *NotificationSuite) TestMoveSubmittedDestinationIsFirstShipmentForSeparatee() {
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypeSEPARATION,
			},
		},
	}, nil)
	notification := NewMoveSubmitted(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Thank you for submitting your move details"

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

func (suite *NotificationSuite) TestMoveSubmittedDestinationIsFirstShipmentForRetiree() {
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypeRETIREMENT,
			},
		},
	}, nil)
	notification := NewMoveSubmitted(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Thank you for submitting your move details"

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

func (suite *NotificationSuite) TestMoveSubmittedDestinationIsDutyStationForPcsType() {
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
			},
		},
	}, nil)
	notification := NewMoveSubmitted(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Thank you for submitting your move details"

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
	suite.Contains(email.textBody, move.Orders.NewDutyLocation.Name)
	suite.NotContains(email.textBody, move.MTOShipments[0].DestinationAddress.StreetAddress1)
	suite.NotContains(email.textBody, *move.MTOShipments[0].DestinationAddress.StreetAddress2)
	suite.NotContains(email.textBody, *move.MTOShipments[0].DestinationAddress.StreetAddress3)
}

func SetupPpmMove(suite *NotificationSuite, ordersType internalmessages.OrdersType) models.Move {
	builtPpmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
			},
		},
	}, nil)
	builtMtoShipment := builtPpmShipment.Shipment
	move := builtMtoShipment.MoveTaskOrder
	move.MTOShipments = append(move.MTOShipments, builtMtoShipment)

	return move
}

func (suite *NotificationSuite) TestMoveSubmittedDestinationIsShipmentForPpmSeparatee() {
	move := SetupPpmMove(suite, internalmessages.OrdersTypeSEPARATION)
	notification := NewMoveSubmitted(move.ID)
	expectedSubject := "Thank you for submitting your move details"

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	mtoShipment := move.MTOShipments[0]
	ppmShipment := mtoShipment.PPMShipment
	destinationAddress := ppmShipment.DestinationAddress
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, expectedSubject)
	suite.Contains(email.htmlBody, "details for your move from\n  "+move.Orders.OriginDutyLocation.Name+" to "+destinationAddress.LineDisplayFormat()+".")
	suite.Contains(email.textBody, "details for your move from "+move.Orders.OriginDutyLocation.Name+" to "+destinationAddress.LineDisplayFormat()+".")
}

func (suite *NotificationSuite) TestMoveSubmittedDestinationIsShipmentForPpmRetiree() {
	move := SetupPpmMove(suite, internalmessages.OrdersTypeRETIREMENT)
	notification := NewMoveSubmitted(move.ID)
	expectedSubject := "Thank you for submitting your move details"

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	mtoShipment := move.MTOShipments[0]
	ppmShipment := mtoShipment.PPMShipment
	destinationAddress := ppmShipment.DestinationAddress
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, expectedSubject)
	suite.NotEmpty(email.htmlBody)
	suite.Contains(email.htmlBody, "details for your move from\n  "+move.Orders.OriginDutyLocation.Name+" to "+destinationAddress.LineDisplayFormat()+".")
	suite.Contains(email.textBody, "details for your move from "+move.Orders.OriginDutyLocation.Name+" to "+destinationAddress.LineDisplayFormat()+".")
}

func (suite *NotificationSuite) TestMoveSubmittedDestinationIsDutyStationForPpmPcsType() {
	move := SetupPpmMove(suite, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
	notification := NewMoveSubmitted(move.ID)
	expectedSubject := "Thank you for submitting your move details"

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	mtoShipment := move.MTOShipments[0]
	ppmShipment := mtoShipment.PPMShipment
	destinationAddress := ppmShipment.DestinationAddress
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, expectedSubject)
	suite.NotEmpty(email.htmlBody)
	suite.NotContains(email.htmlBody, destinationAddress.LineDisplayFormat())
	suite.NotContains(email.textBody, destinationAddress.LineDisplayFormat())
	suite.Contains(email.htmlBody, "details for your move from\n  "+move.Orders.OriginDutyLocation.Name+" to "+move.Orders.NewDutyLocation.Name+".")
	suite.Contains(email.textBody, "details for your move from "+move.Orders.OriginDutyLocation.Name+" to "+move.Orders.NewDutyLocation.Name+".")
}

func (suite *NotificationSuite) TestMoveSubmittedHTMLTemplateRenderWithGovCounseling() {
	approver := factory.BuildUser(nil, nil, nil)
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveSubmitted(move.ID)

	originDutyLocation := "origDutyLocation"
	originDutyLocationPhoneLine := "555-555-5555"

	s := moveSubmittedEmailData{
		OriginDutyLocation:                &originDutyLocation,
		DestinationLocation:               "destDutyLocation",
		OriginDutyLocationPhoneLine:       &originDutyLocationPhoneLine,
		Locator:                           "abc123",
		WeightAllowance:                   "7,999",
		ProvidesGovernmentCounseling:      true,
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
	}
	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>

<p>
  This is a confirmation that you have submitted the details for your move from
  origDutyLocation to destDutyLocation.
</p>

<p>
  <strong>We have assigned you a move code: abc123.</strong> You can use this code when talking to any
  representative about your move.
</p>

<p> To change any information about your move, or to add or cancel shipments, you
  should contact 555-555-5555 or visit your
  <a href="` + OneSourceTransportationOfficeLink + `">local transportation office</a>.
</p>

<p>
  <strong>Your weight allowance: 7,999 pounds.</strong>
  That is how much combined weight the government will pay for all movements between authorized locations under your
  orders.
</p>

<p>
  If you move more than 7,999 pounds or ship to/from an other than authorized location, you may owe the
  government the difference in cost between what you are authorized and what you decide to move.
</p>

<p>
  If you are doing a Household Goods (HHG) shipment: The company responsible for managing your shipment, HomeSafe
  Alliance, will estimate the total weight of your personal property during a pre-move survey, and you will be notified
  if it looks like you might exceed your weight allowance. But you are responsible for excess costs associated with the
  weight moved, up to HomeSafe’s weight estimate plus 10%.
</p>

<h4>Next Steps for your Government-arranged Shipment(s):</h4>

<ul>
  <li>Your move request will be reviewed and a counselor will be assigned to brief you on your move entitlements.</li>
  <li>Your move counselor will get in touch with you soon.</li>
</ul>
<p>Your move counselor will, among other things:</p>
<ul>
  <li>Verify the information you entered and your entitlements</li>
  <li>Give you moving-related advice</li>
  <li>
    Give you tips to avoid excess costs (i.e., going over your weight allowance) or advise you if you are in an excess
    cost scenario
  </li>
</ul>
<p>
  Once your counseling is complete, your request will be reviewed by the responsible personal property shipping office,
  and a move task order will be placed with HomeSafe Alliance. Once this order is placed, you will receive an invitation
  to create an account in HomeSafe Connect. This is the system you will use to schedule your pre-move survey.
</p>

<p>
  HomeSafe is required to contact you within one Government Business Day. Once contact has been established, HomeSafe is
  your primary point of contact. If any information about your move changes at any point during the move, immediately
  notify your HomeSafe Customer Care Representative of the changes.
</p>

<p>
  If you have requested a PPM, <strong>DO NOT</strong> start your PPM until your counselor has approved it in MilMove.
  You will receive an email when that is complete.
</p>

<h4>IMPORTANT: Take the Customer Satisfaction Survey</h4>

<p>
  You will receive an invitation to take a quick customer satisfaction survey (CSS) at key stages of your move process.
  The first invitation will be sent shortly after counseling is complete.
</p>
<p>
  Taking the survey at each stage provides transparency and increases accountability of those assisting you with your
  relocation.
</p>

Thank you,<br />
USTRANSCOM MilMove Team
<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

	htmlContent, err := notification.RenderHTML(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}), s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}

func (suite *NotificationSuite) TestMoveSubmittedHTMLTemplateRenderWithoutGovCounseling() {
	approver := factory.BuildUser(nil, nil, nil)
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveSubmitted(move.ID)

	originDutyLocation := "origDutyLocation"
	originDutyLocationPhoneLine := "555-555-5555"

	s := moveSubmittedEmailData{
		OriginDutyLocation:                &originDutyLocation,
		DestinationLocation:               "destDutyLocation",
		OriginDutyLocationPhoneLine:       &originDutyLocationPhoneLine,
		Locator:                           "abc123",
		WeightAllowance:                   "7,999",
		ProvidesGovernmentCounseling:      false,
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
	}
	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>

<p>
  This is a confirmation that you have submitted the details for your move from
  origDutyLocation to destDutyLocation.
</p>

<p>
  <strong>We have assigned you a move code: abc123.</strong> You can use this code when talking to any
  representative about your move.
</p>

<p> To change any information about your move, or to add or cancel shipments, you
  should contact 555-555-5555 or visit your
  <a href="` + OneSourceTransportationOfficeLink + `">local transportation office</a>.
</p>

<p>
  <strong>Your weight allowance: 7,999 pounds.</strong>
  That is how much combined weight the government will pay for all movements between authorized locations under your
  orders.
</p>

<p>
  If you move more than 7,999 pounds or ship to/from an other than authorized location, you may owe the
  government the difference in cost between what you are authorized and what you decide to move.
</p>

<p>
  If you are doing a Household Goods (HHG) shipment: The company responsible for managing your shipment, HomeSafe
  Alliance, will estimate the total weight of your personal property during a pre-move survey, and you will be notified
  if it looks like you might exceed your weight allowance. But you are responsible for excess costs associated with the
  weight moved, up to HomeSafe’s weight estimate plus 10%.
</p>

<h4>Next Steps for your Government-arranged Shipment(s):</h4>

<p>
  Your move request will be reviewed by the responsible personal property shipping office and a move task order for
  services will be placed with HomeSafe Alliance.
</p>
<p>
  Once this order is placed, you will receive an invitation to create an account in HomeSafe Connect. This is the system
  you will use for your counseling session. You will also schedule your pre-move survey during this session.
</p>

<p>
  HomeSafe is required to contact you within one Government Business Day. Once contact has been established, HomeSafe is
  your primary point of contact. If any information about your move changes at any point during the move, immediately
  notify your HomeSafe Customer Care Representative of the changes.
</p>

<p>
  If you have requested a PPM, <strong>DO NOT</strong> start your PPM until your counselor has approved it in MilMove.
  You will receive an email when that is complete.
</p>

<h4>IMPORTANT: Take the Customer Satisfaction Survey</h4>

<p>
  You will receive an invitation to take a quick customer satisfaction survey (CSS) at key stages of your move process.
  The first invitation will be sent shortly after counseling is complete.
</p>
<p>
  Taking the survey at each stage provides transparency and increases accountability of those assisting you with your
  relocation.
</p>

Thank you,<br />
USTRANSCOM MilMove Team
<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

	htmlContent, err := notification.RenderHTML(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}), s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}

func (suite *NotificationSuite) TestMoveSubmittedHTMLTemplateRenderNoDutyLocation() {
	approver := factory.BuildUser(nil, nil, nil)
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveSubmitted(move.ID)

	s := moveSubmittedEmailData{
		OriginDutyLocation:                nil,
		DestinationLocation:               "destDutyLocation",
		OriginDutyLocationPhoneLine:       nil,
		Locator:                           "abc123",
		WeightAllowance:                   "7,999",
		ProvidesGovernmentCounseling:      false,
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
	}
	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>

<p>
  This is a confirmation that you have submitted the details for your move to destDutyLocation.
</p>

<p>
  <strong>We have assigned you a move code: abc123.</strong> You can use this code when talking to any
  representative about your move.
</p>

<p> To change any information about your move, or to add or cancel shipments, you should
  contact your nearest transportation office. You can find the contact information using the
  <a href="` + OneSourceTransportationOfficeLink + `">directory of PCS-related contacts</a>.
</p>

<p>
  <strong>Your weight allowance: 7,999 pounds.</strong>
  That is how much combined weight the government will pay for all movements between authorized locations under your
  orders.
</p>

<p>
  If you move more than 7,999 pounds or ship to/from an other than authorized location, you may owe the
  government the difference in cost between what you are authorized and what you decide to move.
</p>

<p>
  If you are doing a Household Goods (HHG) shipment: The company responsible for managing your shipment, HomeSafe
  Alliance, will estimate the total weight of your personal property during a pre-move survey, and you will be notified
  if it looks like you might exceed your weight allowance. But you are responsible for excess costs associated with the
  weight moved, up to HomeSafe’s weight estimate plus 10%.
</p>

<h4>Next Steps for your Government-arranged Shipment(s):</h4>

<p>
  Your move request will be reviewed by the responsible personal property shipping office and a move task order for
  services will be placed with HomeSafe Alliance.
</p>
<p>
  Once this order is placed, you will receive an invitation to create an account in HomeSafe Connect. This is the system
  you will use for your counseling session. You will also schedule your pre-move survey during this session.
</p>

<p>
  HomeSafe is required to contact you within one Government Business Day. Once contact has been established, HomeSafe is
  your primary point of contact. If any information about your move changes at any point during the move, immediately
  notify your HomeSafe Customer Care Representative of the changes.
</p>

<p>
  If you have requested a PPM, <strong>DO NOT</strong> start your PPM until your counselor has approved it in MilMove.
  You will receive an email when that is complete.
</p>

<h4>IMPORTANT: Take the Customer Satisfaction Survey</h4>

<p>
  You will receive an invitation to take a quick customer satisfaction survey (CSS) at key stages of your move process.
  The first invitation will be sent shortly after counseling is complete.
</p>
<p>
  Taking the survey at each stage provides transparency and increases accountability of those assisting you with your
  relocation.
</p>

Thank you,<br />
USTRANSCOM MilMove Team
<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

	htmlContent, err := notification.RenderHTML(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}), s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestMoveSubmittedTextTemplateRender() {

	approver := factory.BuildUser(nil, nil, nil)
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveSubmitted(move.ID)

	originDutyLocation := "origDutyLocation"
	originDutyLocationPhoneLine := "555-555-5555"

	s := moveSubmittedEmailData{
		OriginDutyLocation:                &originDutyLocation,
		DestinationLocation:               "destDutyLocation",
		OriginDutyLocationPhoneLine:       &originDutyLocationPhoneLine,
		Locator:                           "abc123",
		WeightAllowance:                   "7,999",
		ProvidesGovernmentCounseling:      true,
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
	}

	expectedTextContent := `*** DO NOT REPLY directly to this email ***

This is a confirmation that you have submitted the details for your move from origDutyLocation to destDutyLocation.

We have assigned you a move code: abc123. You can use this code when talking to any representative about your move.

To change any information about your move, or to add or cancel shipments, you should contact 555-555-5555 or visit your local transportation office (` + OneSourceTransportationOfficeLink + `) .

Your weight allowance: 7,999 pounds. That is how much combined weight the government will pay for all movements between authorized locations under your orders.

If you move more than 7,999 pounds or ship to/from an other than authorized location, you may owe the government the difference in cost between what you are authorized and what you decide to move.

If you are doing a Household Goods (HHG) shipment: The company responsible for managing your shipment, HomeSafe Alliance, will estimate the total weight of your personal property during a pre-move survey, and you will be notified if it looks like you might exceed your weight allowance. But you are responsible for excess costs associated with the weight moved, up to HomeSafe’s weight estimate plus 10%.


** Next Steps for your Government-arranged Shipment(s):
------------------------------------------------------------

* Your move request will be reviewed and a counselor will be assigned to brief you on your move entitlements.
* Your move counselor will get in touch with you soon.

Your move counselor will, among other things:
* Verify the information you entered and your entitlements
* Give you moving-related advice
* Give you tips to avoid excess costs (i.e., going over your weight allowance) or advise you if you are in an excess cost scenario

Once your counseling is complete, your request will be reviewed by the responsible personal property shipping office, and a move task order will be placed with HomeSafe Alliance. Once this order is placed, you will receive an invitation to create an account in HomeSafe Connect. This is the system you will use to schedule your pre-move survey.

HomeSafe is required to contact you within one Government Business Day. Once contact has been established, HomeSafe is your primary point of contact. If any information about your move changes at any point during the move, immediately notify your HomeSafe Customer Care Representative of the changes.

If you have requested a PPM, DO NOT start your PPM until your counselor has approved it in MilMove. You will receive an email when that is complete.


** IMPORTANT: Take the Customer Satisfaction Survey
------------------------------------------------------------

You will receive an invitation to take a quick customer satisfaction survey (CSS) at key stages of your move process. The first invitation will be sent shortly after counseling is complete.

Taking the survey at each stage provides transparency and increases accountability of those assisting you with your relocation.

Thank you,
USTRANSCOM MilMove Team

The information contained in this email may contain Privacy Act information and is therefore protected under the Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.
`

	textContent, err := notification.RenderText(suite.AppContextWithSessionForTest(&auth.Session{
		UserID:          approver.ID,
		ApplicationName: auth.OfficeApp,
	}), s)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}
