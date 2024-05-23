package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	text "text/template"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/models"
)

var (
	officeAccountRejectedRawTextTemplate = string(assets.MustAsset("notifications/templates/office_account_rejected_template.txt"))
	officeAccountRejectedTextTemplate    = text.Must(text.New("text_template").Parse(officeAccountRejectedRawTextTemplate))
	officeAccountRejectedRawHTMLTemplate = string(assets.MustAsset("notifications/templates/office_account_rejected_template.html"))
	officeAccountRejecteddHTMLTemplate   = html.Must(html.New("text_template").Parse(officeAccountRejectedRawHTMLTemplate))
)

// OfficeAccountRejected has notification content for rejected office users
type OfficeAccountRejected struct {
	officeUserID uuid.UUID
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewOfficeAccountRejected returns a new office user rejected notification
func NewOfficeAccountRejected(officeUserID uuid.UUID) *OfficeAccountRejected {

	return &OfficeAccountRejected{
		officeUserID: officeUserID,
		htmlTemplate: officeAccountRejecteddHTMLTemplate,
		textTemplate: officeAccountRejectedTextTemplate,
	}
}

func (o OfficeAccountRejected) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent

	officeUser, err := models.FetchOfficeUserByID(appCtx.DB(), o.officeUserID)
	if err != nil {
		return emails, err
	}

	htmlBody, textBody, err := o.renderTemplates(appCtx, officeAccountRejectedEmailData{
		FirstName:       officeUser.FirstName,
		LastName:        officeUser.LastName,
		RejectionReason: *officeUser.RejectionReason,
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	officeUserEmail := emailContent{
		recipientEmail: officeUser.Email,
		subject:        "MilMove account request was denied",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	appCtx.Logger().Info("Generated office user rejected email",
		zap.String("officeUserID", officeUser.ID.String()))

	return append(emails, officeUserEmail), nil
}

func (o OfficeAccountRejected) renderTemplates(appCtx appcontext.AppContext, data officeAccountRejectedEmailData) (string, string, error) {
	htmlBody, err := o.RenderHTML(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering html template using %#v", data)
	}
	textBody, err := o.RenderText(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering text template using %#v", data)
	}
	return htmlBody, textBody, nil
}

type officeAccountRejectedEmailData struct {
	FirstName       string
	LastName        string
	RejectionReason string
}

// RenderHTML renders the html for the email
func (o OfficeAccountRejected) RenderHTML(appCtx appcontext.AppContext, data officeAccountRejectedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := o.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (o OfficeAccountRejected) RenderText(appCtx appcontext.AppContext, data officeAccountRejectedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := o.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
