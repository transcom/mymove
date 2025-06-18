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
	officeAccountPrivilegeRejectedSupervisorRawTextTemplate = string(assets.MustAsset("notifications/templates/office_account_privilege_rejected_supervisor_template.txt"))
	officeAccountPrivilegeRejectedSupervisorTextTemplate    = text.Must(text.New("text_template").Parse(officeAccountPrivilegeRejectedSupervisorRawTextTemplate))
	officeAccountPrivilegeRejectedSupervisorRawHTMLTemplate = string(assets.MustAsset("notifications/templates/office_account_privilege_rejected_supervisor_template.html"))
	officeAccountPrivilegeRejectedSupervisordHTMLTemplate   = html.Must(html.New("text_template").Parse(officeAccountPrivilegeRejectedSupervisorRawHTMLTemplate))
)

type OfficeAccountPrivilegeRejectedSupervisor struct {
	officeUserID uuid.UUID
	htmlTemplate *html.Template
	textTemplate *text.Template
}

func NewOfficeAccountPrivilegeRejectedSupervisor(officeUserID uuid.UUID) *OfficeAccountPrivilegeRejectedSupervisor {
	return &OfficeAccountPrivilegeRejectedSupervisor{
		officeUserID: officeUserID,
		htmlTemplate: officeAccountPrivilegeRejectedSupervisordHTMLTemplate,
		textTemplate: officeAccountPrivilegeRejectedSupervisorTextTemplate,
	}
}

func (o OfficeAccountPrivilegeRejectedSupervisor) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent

	officeUser, err := models.FetchOfficeUserByID(appCtx.DB(), o.officeUserID)
	if err != nil {
		return emails, err
	}

	htmlBody, textBody, err := o.renderTemplates(appCtx, OfficeAccountPrivilegeRejectedSupervisorEmailData{
		FirstName:     officeUser.FirstName,
		LastName:      officeUser.LastName,
		HelpDeskEmail: "usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil",
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	officeUserEmail := emailContent{
		recipientEmail: officeUser.Email,
		subject:        "MilMove account supervisor privilege request was denied",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	appCtx.Logger().Info("Generated office user supervisor privilege rejected email",
		zap.String("officeUserID", officeUser.ID.String()))

	return append(emails, officeUserEmail), nil
}

func (o OfficeAccountPrivilegeRejectedSupervisor) renderTemplates(appCtx appcontext.AppContext, data OfficeAccountPrivilegeRejectedSupervisorEmailData) (string, string, error) {
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

type OfficeAccountPrivilegeRejectedSupervisorEmailData struct {
	FirstName     string
	LastName      string
	HelpDeskEmail string
}

func (o OfficeAccountPrivilegeRejectedSupervisor) RenderHTML(appCtx appcontext.AppContext, data OfficeAccountPrivilegeRejectedSupervisorEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := o.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

func (o OfficeAccountPrivilegeRejectedSupervisor) RenderText(appCtx appcontext.AppContext, data OfficeAccountPrivilegeRejectedSupervisorEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := o.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
