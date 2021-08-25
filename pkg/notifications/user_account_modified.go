package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	"net/url"
	text "text/template"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/services"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/assets"
)

var (
	userAccountModifiedRawTextTemplate = string(assets.MustAsset("pkg/notifications/templates/user_account_modified_template.txt"))
	userAccountModifiedTextTemplate    = text.Must(text.New("text_template").Parse(userAccountModifiedRawTextTemplate))
	userAccountModifiedRawHTMLTemplate = string(assets.MustAsset("pkg/notifications/templates/user_account_modified_template.html"))
	userAccountModifiedHTMLTemplate    = html.Must(html.New("text_template").Parse(userAccountModifiedRawHTMLTemplate))
)

// UserAccountModified has notification content for alerting admins when a user account has been modified
type UserAccountModified struct {
	logger            Logger
	sysAdminEmail     string
	host              string
	action            string
	modifiedUserID    uuid.UUID
	responsibleUserID uuid.UUID
	modifiedAt        time.Time
	htmlTemplate      *html.Template
	textTemplate      *text.Template
}

// NewUserAccountCreated returns a new UserAccountModified notification for account creation
func NewUserAccountCreated(
	appCtx appcontext.AppContext,
	sysAdminEmail string,
	modifiedUserID uuid.UUID,
	modifiedAt time.Time,
) (*UserAccountModified, error) {
	return newUserAccountModified(appCtx, sysAdminEmail, "created", modifiedUserID, modifiedAt)
}

// NewUserAccountActivated returns a new UserAccountModified notification for account activation
func NewUserAccountActivated(
	appCtx appcontext.AppContext,
	sysAdminEmail string,
	modifiedUserID uuid.UUID,
	modifiedAt time.Time,
) (*UserAccountModified, error) {
	return newUserAccountModified(appCtx, sysAdminEmail, "activated", modifiedUserID, modifiedAt)
}

// NewUserAccountDeactivated returns a new UserAccountModified notification for account deactivation
func NewUserAccountDeactivated(
	appCtx appcontext.AppContext,
	sysAdminEmail string,
	modifiedUserID uuid.UUID,
	modifiedAt time.Time,
) (*UserAccountModified, error) {
	return newUserAccountModified(appCtx, sysAdminEmail, "deactivated", modifiedUserID, modifiedAt)
}

// NewUserAccountRemoved returns a new UserAccountModified notification for account removal
func NewUserAccountRemoved(
	appCtx appcontext.AppContext,
	sysAdminEmail string,
	modifiedUserID uuid.UUID,
	modifiedAt time.Time,
) (*UserAccountModified, error) {
	return newUserAccountModified(appCtx, sysAdminEmail, "removed", modifiedUserID, modifiedAt)
}

// newUserAccountModified returns a new UserAccountModified notification
func newUserAccountModified(
	appCtx appcontext.AppContext,
	sysAdminEmail string,
	action string,
	modifiedUserID uuid.UUID,
	modifiedAt time.Time,
) (*UserAccountModified, error) {
	logger := appCtx.Logger()
	session := auth.SessionFromContext(appCtx.RequestContext())
	if session == nil {
		return nil, services.NewContextError("Unable to find Session in Context")
	}
	responsibleUserID := session.UserID
	host := session.Hostname

	return &UserAccountModified{
		logger:            logger,
		sysAdminEmail:     sysAdminEmail,
		host:              host,
		action:            action,
		modifiedUserID:    modifiedUserID,
		responsibleUserID: responsibleUserID,
		modifiedAt:        modifiedAt,
		htmlTemplate:      userAccountModifiedHTMLTemplate,
		textTemplate:      userAccountModifiedTextTemplate,
	}, nil
}

// userAccountModifiedEmailData has content for email template
type userAccountModifiedEmailData struct {
	Action            string // created, activated, deactivated, or removed
	ActionSource      string // the host URL of the action, where it took place
	ModifiedUserID    string
	ResponsibleUserID string
	Timestamp         string
}

func (m UserAccountModified) emails() ([]emailContent, error) {
	var emails []emailContent

	actionSource := url.URL{
		Scheme: "https",
		Host:   m.host,
	}

	responsibleUserID := m.responsibleUserID
	if responsibleUserID == uuid.Nil {
		responsibleUserID = m.modifiedUserID
	}

	htmlBody, textBody, err := m.renderTemplates(userAccountModifiedEmailData{
		Action:            m.action,
		ActionSource:      actionSource.String(),
		ModifiedUserID:    m.modifiedUserID.String(),
		ResponsibleUserID: responsibleUserID.String(),
		Timestamp:         m.modifiedAt.String(),
	})

	if err != nil {
		m.logger.Error("error rendering template", zap.Error(err))
	}

	adminEmail := emailContent{
		subject:        "[MilMove] User Account Activity Alert",
		recipientEmail: m.sysAdminEmail,
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	m.logger.Info("generated user activity alert email to system admin",
		zap.String("responsibleUserID", responsibleUserID.String()),
		zap.String("modifiedUserID", m.modifiedUserID.String()),
	)

	return append(emails, adminEmail), nil
}

func (m UserAccountModified) renderTemplates(data userAccountModifiedEmailData) (string, string, error) {
	htmlBody, err := m.RenderHTML(data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering html template using %#v", data)
	}
	textBody, err := m.RenderText(data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering text template using %#v", data)
	}
	return htmlBody, textBody, nil
}

// RenderHTML renders the html for the email
func (m UserAccountModified) RenderHTML(data userAccountModifiedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m UserAccountModified) RenderText(data userAccountModifiedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
