package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	text "text/template"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/models"
)

var (
	paymentRequestFailedRawTextTemplate = string(assets.MustAsset("notifications/templates/failed_payment_request_template.txt"))
	paymentRequestFailedTextTemplate    = text.Must(text.New("text_template").Parse(paymentRequestFailedRawTextTemplate))
	paymentRequestFailedRawHTMLTemplate = string(assets.MustAsset("notifications/templates/failed_payment_request_template.html"))
	paymentRequestFailedHTMLTemplate    = html.Must(html.New("text_template").Parse(paymentRequestFailedRawHTMLTemplate))
)

type PaymentRequestFailed struct {
	paymentRequest models.PaymentRequest
	htmlTemplate   *html.Template
	textTemplate   *text.Template
}

func NewPaymentRequestFailed(paymentRequest models.PaymentRequest) *PaymentRequestFailed {
	return &PaymentRequestFailed{
		paymentRequest: paymentRequest,
		htmlTemplate:   paymentRequestFailedHTMLTemplate,
		textTemplate:   paymentRequestFailedTextTemplate,
	}
}

func (p PaymentRequestFailed) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent
	srcEmail, err := models.FetchParameterValueByName(appCtx.DB(), "src_email")
	if err != nil {
		return emails, nil
	}
	dstEmail, err := models.FetchParameterValueByName(appCtx.DB(), "transcom_distro_email")
	if err != nil {
		return emails, nil
	}
	opsEmail, err := models.FetchParameterValueByName(appCtx.DB(), "milmove_ops_email")
	if err != nil {
		return emails, nil
	}

	recipients := []string{}

	if srcEmail.ParameterValue != nil {
		recipients = append(recipients, *srcEmail.ParameterValue)
	}

	if dstEmail.ParameterValue != nil {
		recipients = append(recipients, *dstEmail.ParameterValue)
	}

	if opsEmail.ParameterValue != nil {
		recipients = append(recipients, *opsEmail.ParameterValue)
	}

	ediError, err := models.FetchEdiErrorByPaymentRequestID(appCtx.DB(), p.paymentRequest.ID)
	if err != nil {
		return emails, nil
	}

	if ediError.Code == nil {
		defaultCode := "000"
		ediError.Code = &defaultCode
	}

	htmlBody, textBody, err := p.renderTemplates(appCtx, emailData{
		PaymentRequestNumber: p.paymentRequest.PaymentRequestNumber,
		ErrorCode:            *ediError.Code,
		ErrorDescription:     *ediError.Description,
	})
	if err != nil {
		return emails, err
	}

	email := emailContent{
		recipientEmails: recipients,
		subject:         "Payment Request Failed",
		htmlBody:        htmlBody,
		textBody:        textBody,
	}
	if len(recipients) == 0 {
		return nil, fmt.Errorf("no email found for payment request")
	}
	if len(recipients) == 0 {
		return nil, fmt.Errorf("no email found for payment request")
	}
	emails = append(emails, email)
	return emails, nil
}

func (p PaymentRequestFailed) renderTemplates(appCtx appcontext.AppContext, data emailData) (string, string, error) {
	htmlBody, err := p.RenderHTML(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering html template using %#v", data)
	}
	textBody, err := p.RenderText(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering text template using %#v", data)
	}
	return htmlBody, textBody, nil
}

type emailData struct {
	PaymentRequestNumber string
	ErrorCode            string
	ErrorDescription     string
}

// RenderHTML renders the html for the email
func (p PaymentRequestFailed) RenderHTML(appCtx appcontext.AppContext, data emailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := p.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (p PaymentRequestFailed) RenderText(appCtx appcontext.AppContext, data emailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := p.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
