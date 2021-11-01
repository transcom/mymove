package notifications

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *NotificationSuite) TestClientCertModified() {
	responsibleUser := testdatagen.MakeDefaultUser(suite.DB())
	modifiedClientCert := testdatagen.MakeDevClientCert(suite.DB(),
		testdatagen.Assertions{
			ClientCert: models.ClientCert{
				UserID: responsibleUser.ID,
			},
		})
	session := &auth.Session{
		UserID:   responsibleUser.ID,
		Hostname: "adminlocal",
	}
	appCtx := suite.AppContextWithSessionForTest(session)
	subject := "[MilMove] Client Cert Activity Alert"
	sysAdminEmail := "admin@test.com"

	suite.Run("Create emails for all possible actions", func() {
		// Set up test cases for each action:
		testCases := map[string]struct {
			action     string
			newEmailer func(sysAdminEmail string,
				modifiedClientCertID uuid.UUID,
				modifiedAt time.Time,
				responsibleUserID uuid.UUID,
				host string,
			) (*ClientCertModified, error)
		}{
			"Success - Client cert created notification": {
				action:     "created",
				newEmailer: NewClientCertCreated,
			},
			"Success - Client cert removed notification": {
				action:     "removed",
				newEmailer: NewClientCertRemoved,
			},
		}

		// Loop through and run each test case:
		for name, tc := range testCases {
			suite.Run(name, func() {
				emailer, err := tc.newEmailer(sysAdminEmail, modifiedClientCert.ID, modifiedClientCert.UpdatedAt, appCtx.Session().UserID, appCtx.Session().Hostname)
				suite.Require().NoError(err)
				suite.Require().NotNil(emailer)

				emails, emailErr := emailer.emails(appCtx)
				suite.Require().NoError(emailErr)
				suite.Require().NotNil(emails)
				suite.Equal(len(emails), 1)

				email := emails[0]
				// Check expected values against received values:
				suite.Equal(sysAdminEmail, email.recipientEmail)
				suite.Equal(subject, email.subject)
				suite.NotEmpty(email.htmlBody)
				suite.NotEmpty(email.textBody)

				// Check email content:
				suite.Contains(email.textBody, session.Hostname)
				suite.Contains(email.textBody, fmt.Sprintf("Account %s", tc.action))
				suite.Contains(email.textBody, fmt.Sprintf("Modified client cert ID: %s", modifiedClientCert.ID))
				suite.Contains(email.textBody, fmt.Sprintf("Responsible user ID: %s", responsibleUser.ID))
			})
		}
	})

}

func (suite *NotificationSuite) TestClientCertModifiedHTMLTemplateRender() {
	responsibleUser := testdatagen.MakeDefaultUser(suite.DB())
	modifiedClientCert := testdatagen.MakeDevClientCert(suite.DB(),
		testdatagen.Assertions{
			ClientCert: models.ClientCert{
				UserID: responsibleUser.ID,
			},
		})
	session := &auth.Session{
		UserID:   responsibleUser.ID,
		Hostname: "adminlocal",
	}
	appCtx := suite.AppContextWithSessionForTest(session)

	emailer, err := NewClientCertCreated("", modifiedClientCert.ID, modifiedClientCert.UpdatedAt, appCtx.Session().UserID, appCtx.Session().Hostname)
	suite.Require().NoError(err)
	suite.Require().NotNil(emailer)

	emailData := clientCertModifiedEmailData{
		Action:               "created",
		ActionSource:         "https://adminlocal/",
		ModifiedClientCertID: "f7c602f9-1810-446d-9435-1f1e5cca89eb",
		ResponsibleUserID:    "47ba91f9-660f-4a64-bc41-36b2bf6added",
		Timestamp:            "2021-08-23 23:06:04.897745 +0000 UTC",
	}

	expectedHTMLContent := `<p>
  You are receiving this notification because you are listed as a System Administrator of the MilMove app.
  A MilMove client cert has been created.
</p>

<p><strong>Activity details:</strong></p>
<ul>
  <li>Modified client cert ID: f7c602f9-1810-446d-9435-1f1e5cca89eb</li>
  <li>Responsible user ID: 47ba91f9-660f-4a64-bc41-36b2bf6added</li>
  <li>Action: Account created</li>
  <li>Action source: https://adminlocal/</li>
  <li>Timestamp: 2021-08-23 23:06:04.897745 &#43;0000 UTC</li>
</ul>

<p>
  Please visit the AWS Console
  (<a href="https://dp3.atlassian.net/wiki/spaces/MT/pages/1250066433/0029+AWS+Organization+Authentication">instructions</a>)
  or the <a href="https://admin.move.mil">MilMove Admin Interface</a> to see more details about the above activity.
</p>
`

	htmlContent, err := emailer.RenderHTML(appCtx, emailData)
	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}

func (suite *NotificationSuite) TestClientCertModifiedTextTemplateRender() {
	responsibleUser := testdatagen.MakeDefaultUser(suite.DB())
	modifiedClientCert := testdatagen.MakeDevClientCert(suite.DB(),
		testdatagen.Assertions{
			ClientCert: models.ClientCert{
				UserID: responsibleUser.ID,
			},
		})
	session := &auth.Session{
		UserID:   responsibleUser.ID,
		Hostname: "adminlocal",
	}
	appCtx := suite.AppContextWithSessionForTest(session)

	emailer, err := NewClientCertCreated("", modifiedClientCert.ID, modifiedClientCert.UpdatedAt, appCtx.Session().UserID, appCtx.Session().Hostname)
	suite.Require().NoError(err)
	suite.Require().NotNil(emailer)

	emailData := clientCertModifiedEmailData{
		Action:               "created",
		ActionSource:         "https://adminlocal/",
		ModifiedClientCertID: "f7c602f9-1810-446d-9435-1f1e5cca89eb",
		ResponsibleUserID:    "47ba91f9-660f-4a64-bc41-36b2bf6added",
		Timestamp:            "2021-08-23 23:06:04.897745 +0000 UTC",
	}

	expectedTextContent := `You are receiving this notification because you are listed as a System Administrator of the MilMove app. A MilMove client cert has been created.

Activity details:
* Modified client cert ID: f7c602f9-1810-446d-9435-1f1e5cca89eb
* Responsible user ID: 47ba91f9-660f-4a64-bc41-36b2bf6added
* Action: Account created
* Action source: https://adminlocal/
* Timestamp: 2021-08-23 23:06:04.897745 +0000 UTC

Please visit the AWS Console ([instructions](https://dp3.atlassian.net/wiki/spaces/MT/pages/1250066433/0029+AWS+Organization+Authentication)) or the [MilMove Admin Interface](https://admin.move.mil) to see more details about the above activity.
`

	textContent, err := emailer.RenderText(appCtx, emailData)
	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}
