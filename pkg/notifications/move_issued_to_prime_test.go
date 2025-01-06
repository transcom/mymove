package notifications

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *NotificationSuite) TestMoveIssuedToPrime() {
	move := factory.BuildMove(suite.DB(), nil, nil)
	notification := NewMoveIssuedToPrime(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Your personal property move has been ordered."

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := move.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
}

func (suite *NotificationSuite) TestMoveIssuedToPrimeHTMLTemplateRender() {
	suite.Run("Origin Duty location and Provides Government Counseling.", func() {
		approver := factory.BuildUser(nil, nil, nil)
		move := factory.BuildMove(suite.DB(), nil, nil)
		notification := NewMoveIssuedToPrime(move.ID)
		originDutyLocation := "origDutyLocation"

		s := moveIssuedToPrimeEmailData{
			MilitaryOneSourceLink:        OneSourceTransportationOfficeLink,
			OriginDutyLocation:           &originDutyLocation,
			DestinationLocation:          "destDutyLocation",
			Locator:                      "abc123",
			ProvidesGovernmentCounseling: true,
		}
		expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>

<p>
  This is a confirmation that a move task order has been placed for your move (Move Code abc123)
  from origDutyLocation to destDutyLocation.
</p>

<p>
  What this means to you:
</p>

<p>
  Your government-arranged shipment(s) will be managed by HomeSafe Alliance, the DoD contractor under the Global Household Goods Contract (GHC).
</p>

<h4>Next steps for your government-arranged shipment(s): </h4>

<p>
  HomeSafe will send you an e-mail invitation (check your spam or junk folder) to log in to their system, HomeSafe Connect.
</p>

<ul>
  <li>
    Log in to HomeSafe Connect as soon as possible to schedule your pre-move survey.
    You can request either a virtual, or in-person pre-move survey.
  </li>
</ul>

<p>HomeSafe Customer Care is Required to:</p>

<ul>
  <li>
    Reach out to you within one Government Business Day.
  </li>
  <li>
    Within 3-7 days of your receipt of this e-mail, contact you to provide a 7-day pickup date spread window.
    This spread window must contain your requested pickup date.
    (What this means: your requested pickup date may fall on the spread start date, the spread end date, or anywhere in between.)
  </li>
</ul>

<p>
  If you are requesting to move in 5 days or less, HomeSafe should assist you with scheduling within one day of your
  receipt of this email.
</p>

<p>Utilize your HomeSafe Customer Care Representative:</p>
<ul>
  <li>As your first contact if you have any questions during your move.</li>
  <li>To provide any updates on your shipment or status.</li>
</ul>

<p>
  If you are unsatisfied at any time, contact a government transportation office.
  You can see a listing of transportation offices on Military One Source here:
  <a href="` + OneSourceTransportationOfficeLink + `">` + OneSourceTransportationOfficeLink + `</a>.
</p>

<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

		htmlContent, err := notification.RenderHTML(suite.AppContextWithSessionForTest(&auth.Session{
			UserID:          approver.ID,
			ApplicationName: auth.OfficeApp,
		}), s)

		suite.NoError(err)
		suite.Equal(expectedHTMLContent, htmlContent)
	})

	suite.Run("No Origin Duty location and Provides Government Counseling.", func() {
		approver := factory.BuildUser(nil, nil, nil)
		move := factory.BuildMove(suite.DB(), nil, nil)
		notification := NewMoveIssuedToPrime(move.ID)

		s := moveIssuedToPrimeEmailData{
			MilitaryOneSourceLink:        OneSourceTransportationOfficeLink,
			DestinationLocation:          "destDutyLocation",
			Locator:                      "abc123",
			ProvidesGovernmentCounseling: true,
		}
		expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>

<p>
  This is a confirmation that a move task order has been placed for your move (Move Code abc123)
  to destDutyLocation.
</p>

<p>
  What this means to you:
</p>

<p>
  Your government-arranged shipment(s) will be managed by HomeSafe Alliance, the DoD contractor under the Global Household Goods Contract (GHC).
</p>

<h4>Next steps for your government-arranged shipment(s): </h4>

<p>
  HomeSafe will send you an e-mail invitation (check your spam or junk folder) to log in to their system, HomeSafe Connect.
</p>

<ul>
  <li>
    Log in to HomeSafe Connect as soon as possible to schedule your pre-move survey.
    You can request either a virtual, or in-person pre-move survey.
  </li>
</ul>

<p>HomeSafe Customer Care is Required to:</p>

<ul>
  <li>
    Reach out to you within one Government Business Day.
  </li>
  <li>
    Within 3-7 days of your receipt of this e-mail, contact you to provide a 7-day pickup date spread window.
    This spread window must contain your requested pickup date.
    (What this means: your requested pickup date may fall on the spread start date, the spread end date, or anywhere in between.)
  </li>
</ul>

<p>
  If you are requesting to move in 5 days or less, HomeSafe should assist you with scheduling within one day of your
  receipt of this email.
</p>

<p>Utilize your HomeSafe Customer Care Representative:</p>
<ul>
  <li>As your first contact if you have any questions during your move.</li>
  <li>To provide any updates on your shipment or status.</li>
</ul>

<p>
  If you are unsatisfied at any time, contact a government transportation office.
  You can see a listing of transportation offices on Military One Source here:
  <a href="` + OneSourceTransportationOfficeLink + `">` + OneSourceTransportationOfficeLink + `</a>.
</p>

<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

		htmlContent, err := notification.RenderHTML(suite.AppContextWithSessionForTest(&auth.Session{
			UserID:          approver.ID,
			ApplicationName: auth.OfficeApp,
		}), s)

		suite.NoError(err)
		suite.Equal(expectedHTMLContent, htmlContent)
	})

	suite.Run("Origin Duty location and Does Not Provide Government Counseling.", func() {
		approver := factory.BuildUser(nil, nil, nil)
		move := factory.BuildMove(suite.DB(), nil, nil)
		notification := NewMoveIssuedToPrime(move.ID)
		originDutyLocation := "origDutyLocation"

		s := moveIssuedToPrimeEmailData{
			MilitaryOneSourceLink:        OneSourceTransportationOfficeLink,
			OriginDutyLocation:           &originDutyLocation,
			DestinationLocation:          "destDutyLocation",
			Locator:                      "abc123",
			ProvidesGovernmentCounseling: false,
		}
		expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>

<p>
  This is a confirmation that a move task order has been placed for your move (Move Code abc123)
  from origDutyLocation to destDutyLocation.
</p>

<p>
  What this means to you:
</p>

<p>
  If you have requested a Personally Procured Move (PPM), <strong>DO NOT</strong> start your PPM until it has been approved by your counselor.
  You will receive an email when that is complete.
</p>

<p>
  Your government-arranged shipment(s) will be managed by HomeSafe Alliance, the DoD contractor under the Global Household Goods Contract (GHC).
</p>

<h4>Next steps for your government-arranged shipment(s): </h4>

<p>
  HomeSafe will send you an e-mail invitation (check your spam or junk folder) to log in to their system, HomeSafe Connect.
</p>

<ul>
  <li>
    Log in to HomeSafe Connect as soon as possible to complete counseling and schedule your pre-move survey.
    You can request either a virtual, or in-person pre-move survey.
  </li>
</ul>

<p>HomeSafe Customer Care is Required to:</p>

<ul>
  <li>
    Reach out to you within one Government Business Day.
  </li>
  <li>
    Within 3-7 days of your receipt of this e-mail, contact you to assist in completion of counseling
    and provide a 7-day pickup date spread window. This spread window must contain your requested pickup date.
    (What this means: your requested pickup date may fall on the spread start date, the spread end date, or anywhere in between.)
  </li>
</ul>

<p>
  If you are requesting to move in 5 days or less, HomeSafe should assist you with scheduling within one day of your
  receipt of this email.
</p>

<p>Utilize your HomeSafe Customer Care Representative:</p>
<ul>
  <li>As your first contact if you have any questions during your move.</li>
  <li>To provide any updates on your shipment or status.</li>
</ul>

<p>
  If you are unsatisfied at any time, contact a government transportation office.
  You can see a listing of transportation offices on Military One Source here:
  <a href="` + OneSourceTransportationOfficeLink + `">` + OneSourceTransportationOfficeLink + `</a>.
</p>

<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

		htmlContent, err := notification.RenderHTML(suite.AppContextWithSessionForTest(&auth.Session{
			UserID:          approver.ID,
			ApplicationName: auth.OfficeApp,
		}), s)

		suite.NoError(err)
		suite.Equal(expectedHTMLContent, htmlContent)
	})
}

func (suite *NotificationSuite) TestMoveIssuedToPrimeTextTemplateRender() {
	suite.Run("Origin Duty location and Provides Government Counseling.", func() {
		approver := factory.BuildUser(nil, nil, nil)
		move := factory.BuildMove(suite.DB(), nil, nil)
		notification := NewMoveIssuedToPrime(move.ID)
		originDutyLocation := "origDutyLocation"

		s := moveIssuedToPrimeEmailData{
			MilitaryOneSourceLink:        OneSourceTransportationOfficeLink,
			OriginDutyLocation:           &originDutyLocation,
			DestinationLocation:          "destDutyLocation",
			Locator:                      "abc123",
			ProvidesGovernmentCounseling: true,
		}
		expectedTextContent := `*** DO NOT REPLY directly to this email ***

This is a confirmation that a move task order has been placed for your move (Move Code abc123)
from origDutyLocation to destDutyLocation.

What this means to you:

Your government-arranged shipment(s) will be managed by HomeSafe Alliance,
the DoD contractor under the Global Household Goods Contract (GHC).

*** Next steps for your government-arranged shipment(s): ***

HomeSafe will send you an e-mail invitation (check your spam or junk folder) to log in to their system, HomeSafe Connect.

* Log in to HomeSafe Connect as soon as possible to schedule your pre-move survey. You can request either a virtual,
or in-person pre-move survey.

HomeSafe Customer Care is Required to:
* Reach out to you within one Government Business Day.
* Within 3-7 days of your receipt of this e-mail, contact you to provide a 7-day pickup date spread window.
This spread window must contain your requested pickup date.
(What this means: your requested pickup date may fall on the spread start date, the spread end date, or anywhere in between.)

If you are requesting to move in 5 days or less, HomeSafe should assist you with scheduling within one day of your receipt of this email.

Utilize your HomeSafe Customer Care Representative:
* As your first contact if you have any questions during your move.
* To provide any updates on your shipment or status.

If you are unsatisfied at any time, contact a government transportation office.
You can see a listing of transportation offices on Military One Source here:
` + OneSourceTransportationOfficeLink + `.

Thank you,

USTRANSCOM MilMove Team

The information contained in this email may contain Privacy Act information and is therefore protected
under the Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.
`

		textContent, err := notification.RenderText(suite.AppContextWithSessionForTest(&auth.Session{
			UserID:          approver.ID,
			ApplicationName: auth.OfficeApp,
		}), s)

		suite.NoError(err)
		suite.Equal(expectedTextContent, textContent)
	})

	suite.Run("No Origin Duty location and Provides Government Counseling.", func() {
		approver := factory.BuildUser(nil, nil, nil)
		move := factory.BuildMove(suite.DB(), nil, nil)
		notification := NewMoveIssuedToPrime(move.ID)

		s := moveIssuedToPrimeEmailData{
			MilitaryOneSourceLink:        OneSourceTransportationOfficeLink,
			DestinationLocation:          "destDutyLocation",
			Locator:                      "abc123",
			ProvidesGovernmentCounseling: true,
		}
		expectedTextContent := `*** DO NOT REPLY directly to this email ***

This is a confirmation that a move task order has been placed for your move (Move Code abc123)
to destDutyLocation.

What this means to you:

Your government-arranged shipment(s) will be managed by HomeSafe Alliance,
the DoD contractor under the Global Household Goods Contract (GHC).

*** Next steps for your government-arranged shipment(s): ***

HomeSafe will send you an e-mail invitation (check your spam or junk folder) to log in to their system, HomeSafe Connect.

* Log in to HomeSafe Connect as soon as possible to schedule your pre-move survey. You can request either a virtual,
or in-person pre-move survey.

HomeSafe Customer Care is Required to:
* Reach out to you within one Government Business Day.
* Within 3-7 days of your receipt of this e-mail, contact you to provide a 7-day pickup date spread window.
This spread window must contain your requested pickup date.
(What this means: your requested pickup date may fall on the spread start date, the spread end date, or anywhere in between.)

If you are requesting to move in 5 days or less, HomeSafe should assist you with scheduling within one day of your receipt of this email.

Utilize your HomeSafe Customer Care Representative:
* As your first contact if you have any questions during your move.
* To provide any updates on your shipment or status.

If you are unsatisfied at any time, contact a government transportation office.
You can see a listing of transportation offices on Military One Source here:
` + OneSourceTransportationOfficeLink + `.

Thank you,

USTRANSCOM MilMove Team

The information contained in this email may contain Privacy Act information and is therefore protected
under the Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.
`

		textContent, err := notification.RenderText(suite.AppContextWithSessionForTest(&auth.Session{
			UserID:          approver.ID,
			ApplicationName: auth.OfficeApp,
		}), s)

		suite.NoError(err)
		suite.Equal(expectedTextContent, textContent)
	})

	suite.Run("Origin Duty location and Does Not Provide Government Counseling.", func() {
		approver := factory.BuildUser(nil, nil, nil)
		move := factory.BuildMove(suite.DB(), nil, nil)
		notification := NewMoveIssuedToPrime(move.ID)
		originDutyLocation := "origDutyLocation"

		s := moveIssuedToPrimeEmailData{
			MilitaryOneSourceLink:        OneSourceTransportationOfficeLink,
			OriginDutyLocation:           &originDutyLocation,
			DestinationLocation:          "destDutyLocation",
			Locator:                      "abc123",
			ProvidesGovernmentCounseling: false,
		}
		expectedTextContent := `*** DO NOT REPLY directly to this email ***

This is a confirmation that a move task order has been placed for your move (Move Code abc123)
from origDutyLocation to destDutyLocation.

What this means to you:

If you have requested a Personally Procured Move (PPM), DO NOT start your PPM until it has been approved by your counselor.
You will receive an email when that is complete.

Your government-arranged shipment(s) will be managed by HomeSafe Alliance,
the DoD contractor under the Global Household Goods Contract (GHC).

*** Next steps for your government-arranged shipment(s): ***

HomeSafe will send you an e-mail invitation (check your spam or junk folder) to log in to their system, HomeSafe Connect.

* Log in to HomeSafe Connect as soon as possible to complete counseling and schedule your pre-move survey.
You can request either a virtual, or in-person pre-move survey.

HomeSafe Customer Care is Required to:
* Reach out to you within one Government Business Day.
* Within 3-7 days of your receipt of this e-mail, contact you to assist in completion of counseling
and provide a 7-day pickup date spread window. This spread window must contain your requested pickup date.
(What this means: your requested pickup date may fall on the spread start date, the spread end date, or anywhere in between.)

If you are requesting to move in 5 days or less, HomeSafe should assist you with scheduling within one day of your receipt of this email.

Utilize your HomeSafe Customer Care Representative:
* As your first contact if you have any questions during your move.
* To provide any updates on your shipment or status.

If you are unsatisfied at any time, contact a government transportation office.
You can see a listing of transportation offices on Military One Source here:
` + OneSourceTransportationOfficeLink + `.

Thank you,

USTRANSCOM MilMove Team

The information contained in this email may contain Privacy Act information and is therefore protected
under the Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.
`

		textContent, err := notification.RenderText(suite.AppContextWithSessionForTest(&auth.Session{
			UserID:          approver.ID,
			ApplicationName: auth.OfficeApp,
		}), s)

		suite.NoError(err)
		suite.Equal(expectedTextContent, textContent)
	})
}

func (suite *NotificationSuite) TestMoveIssuedToPrimeTOOApprovedMoveDetailsForSeparatee() {
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypeSEPARATION,
			},
		},
	}, nil)
	notification := NewMoveIssuedToPrime(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Your personal property move has been ordered."

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

func (suite *NotificationSuite) TestMoveIssuedToPrimeTOOApprovedMoveDetailsForRetiree() {
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypeRETIREMENT,
			},
		},
	}, nil)
	notification := NewMoveIssuedToPrime(move.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: move.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Your personal property move has been ordered."

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
