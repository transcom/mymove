package notifications

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *NotificationSuite) TestUserAccountModified() {
	var modifiedUser models.User
	var responsibleUser models.User

	suite.PreloadData(func() {
		modifiedUser = testdatagen.MakeStubbedUser(suite.DB())
		responsibleUser = testdatagen.MakeStubbedUser(suite.DB())
	})
	session := auth.Session{
		UserID:   responsibleUser.ID,
		Hostname: "adminlocal",
	}

	subject := "[MilMove] User Account Activity Alert"
	sysAdminEmail := "admin@test.com"

	suite.Run("Create emails for all possible actions", func() {
		// Set up test cases for each action:
		testCases := map[string]struct {
			action     string
			newEmailer func(appCtx appcontext.AppContext,
				sysAdminEmail string,
				modifiedUserID uuid.UUID,
				modifiedAt time.Time,
			) (*UserAccountModified, error)
		}{
			"Success - User account created notification": {
				action:     "created",
				newEmailer: NewUserAccountCreated,
			},
			"Success - User account activated notification": {
				action:     "activated",
				newEmailer: NewUserAccountActivated,
			},
			"Success - User account deactivated notification": {
				action:     "deactivated",
				newEmailer: NewUserAccountDeactivated,
			},
			"Success - User account removed notification": {
				action:     "removed",
				newEmailer: NewUserAccountRemoved,
			},
		}

		// Loop through and run each test case:
		for name, tc := range testCases {
			suite.Run(name, func() {
				emailer, err := tc.newEmailer(suite.AppContextWithSessionForTest(&session), sysAdminEmail, modifiedUser.ID, modifiedUser.UpdatedAt)
				suite.Require().NoError(err)
				suite.Require().NotNil(emailer)

				emails, emailErr := emailer.emails(suite.AppContextWithSessionForTest(&session))
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
				suite.Contains(email.textBody, fmt.Sprintf("Modified user ID: %s", modifiedUser.ID))
				suite.Contains(email.textBody, fmt.Sprintf("Responsible user ID: %s", responsibleUser.ID))
			})
		}
	})

	suite.Run("Success - User account creation with no user in session", func() {
		// Test case:   If a user just created their account, their userID information might not be in the session yet.
		// Expectation: The email should use the modified user ID as the responsible user ID as well.
		emptySessionCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		emailer, err := NewUserAccountCreated(emptySessionCtx, sysAdminEmail, modifiedUser.ID, modifiedUser.UpdatedAt)
		suite.Require().NoError(err)
		suite.Require().NotNil(emailer)

		emails, emailErr := emailer.emails(emptySessionCtx)
		suite.Require().NoError(emailErr)
		suite.Require().NotNil(emails)

		email := emails[0]
		suite.Require().NotEmpty(email.textBody)
		suite.Contains(email.textBody, fmt.Sprintf("Modified user ID: %s", modifiedUser.ID))
		suite.Contains(email.textBody, fmt.Sprintf("Responsible user ID: %s", modifiedUser.ID))
	})

	suite.Run("Fail - Session is nil", func() {
		// Test case:   The session wasn't set in the AppContext, for some reason. Possibly dev error.
		// Expectation: Initializing the UserAccountModified should return services.ContextError
		nilSessionCtx := suite.AppContextForTest()

		emailer, err := NewUserAccountCreated(nilSessionCtx, sysAdminEmail, modifiedUser.ID, modifiedUser.UpdatedAt)
		suite.Nil(emailer)
		suite.Error(err)
		suite.IsType(apperror.ContextError{}, err)
	})
}

func (suite *NotificationSuite) TestUserAccountModifiedHTMLTemplateRender() {
	modifiedUser := testdatagen.MakeStubbedUser(suite.DB())
	appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

	emailer, err := NewUserAccountCreated(appCtx, "", modifiedUser.ID, modifiedUser.UpdatedAt)
	suite.Require().NoError(err)
	suite.Require().NotNil(emailer)

	emailData := userAccountModifiedEmailData{
		Action:            "created",
		ActionSource:      "https://adminlocal/",
		ModifiedUserID:    "f7c602f9-1810-446d-9435-1f1e5cca89eb",
		ResponsibleUserID: "47ba91f9-660f-4a64-bc41-36b2bf6added",
		Timestamp:         "2021-08-23 23:06:04.897745 +0000 UTC",
	}

	expectedHTMLContent := `<p>
  You are receiving this notification because you are listed as a System Administrator of the MilMove app.
  A MilMove user account has been created.
</p>

<p><strong>Activity details:</strong></p>
<ul>
  <li>Modified user ID: f7c602f9-1810-446d-9435-1f1e5cca89eb</li>
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

func (suite *NotificationSuite) TestUserAccountModifiedTextTemplateRender() {
	modifiedUser := testdatagen.MakeStubbedUser(suite.DB())
	appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

	emailer, err := NewUserAccountCreated(appCtx, "", modifiedUser.ID, modifiedUser.UpdatedAt)
	suite.Require().NoError(err)
	suite.Require().NotNil(emailer)

	emailData := userAccountModifiedEmailData{
		Action:            "created",
		ActionSource:      "https://adminlocal/",
		ModifiedUserID:    "f7c602f9-1810-446d-9435-1f1e5cca89eb",
		ResponsibleUserID: "47ba91f9-660f-4a64-bc41-36b2bf6added",
		Timestamp:         "2021-08-23 23:06:04.897745 +0000 UTC",
	}

	expectedTextContent := `You are receiving this notification because you are listed as a System Administrator of the MilMove app. A MilMove user account has been created.

Activity details:
* Modified user ID: f7c602f9-1810-446d-9435-1f1e5cca89eb
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
