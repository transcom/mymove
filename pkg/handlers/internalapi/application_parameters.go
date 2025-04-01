package internalapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/application_parameters"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForApplicationParametersModel(v models.ApplicationParameters) internalmessages.ApplicationParameters {

	parameterValue := v.ParameterValue
	parameterName := v.ParameterName

	payload := internalmessages.ApplicationParameters{
		ParameterValue: parameterValue,
		ParameterName:  parameterName,
	}
	return payload
}

// ApplicationParametersValidateHandler validates a value provided by the service member
type ApplicationParametersValidateHandler struct {
	handlers.HandlerConfig
}

// Handler receives a POST request containing a parameter value
// if the value is present, it returns it back, if not, it returns an empty object
func (h ApplicationParametersValidateHandler) Handle(params application_parameters.ValidateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// we are only allowing this to be called from the customer app
			// since this is an open route outside of auth, we want to buckle down on validation here
			if !appCtx.Session().IsMilApp() {
				return application_parameters.NewValidateUnauthorized(), apperror.NewSessionError("Request is not from the customer application")
			}

			// receive the value
			value := params.Body.ParameterValue
			name := params.Body.ParameterName

			// fetch the value, if not found it will be an empty string
			result, _ := models.FetchParameterValue(appCtx.DB(), *name, *value)

			parameterValuePayload := payloadForApplicationParametersModel(result)

			return application_parameters.NewValidateOK().WithPayload(&parameterValuePayload), nil
		})
}
