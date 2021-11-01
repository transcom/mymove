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

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/transcom/mymove/pkg/assets"
)

var (
	clientCertModifiedRawTextTemplate = string(assets.MustAsset("notifications/templates/client_cert_modified_template.txt"))
	clientCertModifiedTextTemplate    = text.Must(text.New("text_template").Parse(clientCertModifiedRawTextTemplate))
	clientCertModifiedRawHTMLTemplate = string(assets.MustAsset("notifications/templates/client_cert_modified_template.html"))
	clientCertModifiedHTMLTemplate    = html.Must(html.New("text_template").Parse(clientCertModifiedRawHTMLTemplate))
)

// ClientCertModified has notification content for alerting admins when a user account has been modified
type ClientCertModified struct {
	sysAdminEmail        string
	host                 string
	action               string
	modifiedClientCertID uuid.UUID
	responsibleUserID    uuid.UUID
	modifiedAt           time.Time
	htmlTemplate         *html.Template
	textTemplate         *text.Template
}

// NewClientCertCreated returns a new ClientCertModified notification for account creation
func NewClientCertCreated(
	sysAdminEmail string,
	modifiedClientCertID uuid.UUID,
	modifiedAt time.Time,
	responsibleUserID uuid.UUID,
	host string,
) (*ClientCertModified, error) {

	return newClientCertModified(sysAdminEmail, "created", modifiedClientCertID, modifiedAt, responsibleUserID, host)
}

// NewClientCertRemoved returns a new ClientCertModified notification for account removal
func NewClientCertRemoved(
	sysAdminEmail string,
	modifiedClientCertID uuid.UUID,
	modifiedAt time.Time,
	responsibleUserID uuid.UUID,
	host string,
) (*ClientCertModified, error) {
	return newClientCertModified(sysAdminEmail, "removed", modifiedClientCertID, modifiedAt, responsibleUserID, host)
}

// newClientCertModified returns a new ClientCertModified notification
func newClientCertModified(
	sysAdminEmail string,
	action string,
	modifiedClientCertID uuid.UUID,
	modifiedAt time.Time,
	responsibleUserID uuid.UUID,
	host string,
) (*ClientCertModified, error) {
	return &ClientCertModified{
		sysAdminEmail:        sysAdminEmail,
		action:               action,
		modifiedClientCertID: modifiedClientCertID,
		modifiedAt:           modifiedAt,
		responsibleUserID:    responsibleUserID,
		host:                 host,
		htmlTemplate:         clientCertModifiedHTMLTemplate,
		textTemplate:         clientCertModifiedTextTemplate,
	}, nil
}

// clientCertModifiedEmailData has content for email template
type clientCertModifiedEmailData struct {
	Action               string // created, activated, deactivated, or removed
	ActionSource         string // the host URL of the action, where it took place
	ModifiedClientCertID string
	ResponsibleUserID    string
	Timestamp            string
}

func (m ClientCertModified) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent

	actionSource := url.URL{
		Scheme: "https",
		Host:   m.host,
	}

	htmlBody, textBody, err := m.renderTemplates(appCtx, clientCertModifiedEmailData{
		Action:               m.action,
		ActionSource:         actionSource.String(),
		ModifiedClientCertID: m.modifiedClientCertID.String(),
		ResponsibleUserID:    m.responsibleUserID.String(),
		Timestamp:            m.modifiedAt.String(),
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	adminEmail := emailContent{
		subject:        "[MilMove] Client Cert Activity Alert",
		recipientEmail: m.sysAdminEmail,
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	appCtx.Logger().Info("generated client cert activity alert email to system admin",
		zap.String("responsibleUserID", m.responsibleUserID.String()),
		zap.String("modifiedClientCertID", m.modifiedClientCertID.String()),
	)

	return append(emails, adminEmail), nil
}

func (m ClientCertModified) renderTemplates(appCtx appcontext.AppContext, data clientCertModifiedEmailData) (string, string, error) {
	htmlBody, err := m.RenderHTML(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering html template using %#v", data)
	}
	textBody, err := m.RenderText(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering text template using %#v", data)
	}
	return htmlBody, textBody, nil
}

// RenderHTML renders the html for the email
func (m ClientCertModified) RenderHTML(appCtx appcontext.AppContext, data clientCertModifiedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m ClientCertModified) RenderText(appCtx appcontext.AppContext, data clientCertModifiedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
