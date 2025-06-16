package notifications

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *NotificationSuite) TestOfficeAccountPrivilegeRejectedSupervisor() {

	subject := "MilMove account supervisor privilege request was denied"

	modifiedUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

	notification := NewOfficeAccountPrivilegeRejectedSupervisor(modifiedUser.ID)
	suite.Require().NotNil(notification)

	emails, err := notification.emails(appCtx)
	suite.NoError(err)
	suite.Len(emails, 1)

	email := emails[0]
	suite.Equal(email.recipientEmail, modifiedUser.Email)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
}

func (suite *NotificationSuite) TestOfficeAccountPrivilegeRejectedSupervisorHTMLTemplateRender() {
	modifiedUser := factory.BuildOfficeUser(nil, nil, nil)
	appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

	emailer := NewOfficeAccountPrivilegeRejectedSupervisor(modifiedUser.ID)
	suite.Require().NotNil(emailer)

	emailData := OfficeAccountPrivilegeRejectedSupervisorEmailData{
		FirstName:     "Leo",
		LastName:      "Spaceman",
		HelpDeskEmail: "usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil",
	}

	expectedHTMLContent := `<img src="https://raw.githubusercontent.com/transcom/mymove-docs/main/static/img/logos/milmove-logo-black.png" alt="MilMove Logo" width="150" height="30">

<p>Hi Leo Spaceman,</p>

<p>&emsp;Your request for supervisor privilege for your MilMove Office user account has been denied.</p>

<p>&emsp;If you feel this denial was in error, please contact the Technical Help Desk at <a href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil">usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil</a> and provide adequate justification.</p>

<p>V/R,</p>
<p>MilMove Admin</p>
`
	htmlContent, err := emailer.RenderHTML(appCtx, emailData)
	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}

func (suite *NotificationSuite) TestOfficeAccountPrivilegeRejectedSupervisorTextTemplateRender() {
	modifiedUser := factory.BuildOfficeUser(nil, nil, nil)
	appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

	emailer := NewOfficeAccountPrivilegeRejectedSupervisor(modifiedUser.ID)
	suite.Require().NotNil(emailer)

	emailData := OfficeAccountPrivilegeRejectedSupervisorEmailData{
		FirstName:     "Leo",
		LastName:      "Spaceman",
		HelpDeskEmail: "usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil",
	}

	expectedTextContent := `Hi Leo Spaceman,

     Your request for supervisor privilege for your MilMove Office user account has been denied.

    If you feel this denial was in error, please contact the Technical Help Desk at usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil and provide adequate justification.

V/R,

MilMove Admin
`

	textContent, err := emailer.RenderText(appCtx, emailData)
	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}
