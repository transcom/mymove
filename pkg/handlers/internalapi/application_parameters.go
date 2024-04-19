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

// ApplicationParametersValidateHandler validates a code provided by the service member
type ApplicationParametersValidateHandler struct {
	handlers.HandlerConfig
}

// Handler receives a POST request containing a validation code
// if the code is present, it returns it back, if not, it returns an empty object
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
