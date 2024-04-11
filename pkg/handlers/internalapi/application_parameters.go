package internalapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/application_parameters"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForApplicationParametersModel(v models.ApplicationParameters) internalmessages.ValidationCode {
	payload := internalmessages.ValidationCode{
		ValidationCode: *handlers.FmtString(v.ValidationCode),
	}
	return payload
}

// GetOktaProfileHandler gets Okta Profile via GET /okta-profile
type ApplicationParametersValidateHandler struct {
	handlers.HandlerConfig
}

// Handle performs a POST request from Okta API, returns values in profile object from response
func (h ApplicationParametersValidateHandler) Handle(params application_parameters.ValidateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// receive the code
			code := params.Body.ValidationCode

			// fetch the code, if not found it will be an empty string
			result, _ := models.FetchValidationCode(appCtx.DB(), code)

			validationCodePayload := payloadForApplicationParametersModel(result)

			return application_parameters.NewValidateOK().WithPayload(&validationCodePayload), nil
		})
}
