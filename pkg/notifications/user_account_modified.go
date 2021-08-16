package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	"net/url"
	text "text/template"
	"time"

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
	host              string
	sysAdminEmail     string
	action            string
	modifiedUserID    uuid.UUID
	responsibleUserID uuid.UUID
	traceID           uuid.UUID
	modifiedAt        time.Time
	htmlTemplate      *html.Template
	textTemplate      *text.Template
}

// NewUserAccountCreated returns a new UserAccountModified notification for account creation
func NewUserAccountCreated(
	logger Logger,
	host string,
	sysAdminEmail string,
	modifiedUserID uuid.UUID,
	responsibleUserID uuid.UUID,
	traceID uuid.UUID,
	modifiedAt time.Time,
) *UserAccountModified {
	return newUserAccountModified(
		logger, host, sysAdminEmail, "created", modifiedUserID, responsibleUserID, traceID, modifiedAt)
}

// NewUserAccountActivated returns a new UserAccountModified notification for account activation
func NewUserAccountActivated(
	logger Logger,
	host string,
	sysAdminEmail string,
	modifiedUserID uuid.UUID,
	responsibleUserID uuid.UUID,
	traceID uuid.UUID,
	modifiedAt time.Time,
) *UserAccountModified {
	return newUserAccountModified(
		logger, host, sysAdminEmail, "activated", modifiedUserID, responsibleUserID, traceID, modifiedAt)
}

// NewUserAccountDeactivated returns a new UserAccountModified notification for account deactivation
func NewUserAccountDeactivated(
	logger Logger,
	host string,
	sysAdminEmail string,
	modifiedUserID uuid.UUID,
	responsibleUserID uuid.UUID,
	traceID uuid.UUID,
	modifiedAt time.Time,
) *UserAccountModified {
	return newUserAccountModified(
		logger, host, sysAdminEmail, "deactivated", modifiedUserID, responsibleUserID, traceID, modifiedAt)
}

// NewUserAccountRemoved returns a new UserAccountModified notification for account removal
func NewUserAccountRemoved(
	logger Logger,
	host string,
	sysAdminEmail string,
	modifiedUserID uuid.UUID,
	responsibleUserID uuid.UUID,
	traceID uuid.UUID,
	modifiedAt time.Time,
) *UserAccountModified {
	return newUserAccountModified(
		logger, host, sysAdminEmail, "removed", modifiedUserID, responsibleUserID, traceID, modifiedAt)
}

// newUserAccountModified returns a new UserAccountModified notification
func newUserAccountModified(
	logger Logger,
	host string,
	sysAdminEmail string,
	action string,
	modifiedUserID uuid.UUID,
	responsibleUserID uuid.UUID,
	traceID uuid.UUID,
	modifiedAt time.Time,
) *UserAccountModified {

	return &UserAccountModified{
		logger:            logger,
		host:              host,
		sysAdminEmail:     sysAdminEmail,
		action:            action,
		modifiedUserID:    modifiedUserID,
		responsibleUserID: responsibleUserID,
		traceID:           traceID,
		modifiedAt:        modifiedAt,
		htmlTemplate:      userAccountModifiedHTMLTemplate,
		textTemplate:      userAccountModifiedTextTemplate,
	}
}

// userAccountModifiedEmailData has content for email template
type userAccountModifiedEmailData struct {
	Action            string // created, activated, deactivated, or removed
	ModifiedUserID    string
	ResponsibleUserID string
	TraceID           string
	Timestamp         string
	AdminLink         string
}

func (m UserAccountModified) emails() ([]emailContent, error) {
	var emails []emailContent

	adminURL := url.URL{
		Scheme: "https",
		Host:   m.host,
	}

	htmlBody, textBody, err := m.renderTemplates(userAccountModifiedEmailData{
		Action:            m.action,
		ModifiedUserID:    m.modifiedUserID.String(),
		ResponsibleUserID: m.responsibleUserID.String(),
		TraceID:           m.traceID.String(),
		Timestamp:         m.modifiedAt.String(),
		AdminLink:         adminURL.String(),
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
		zap.String("responsibleUserID", m.responsibleUserID.String()),
		zap.String("modifiedUserID", m.modifiedUserID.String()),
		zap.String("traceID", m.traceID.String()),
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
