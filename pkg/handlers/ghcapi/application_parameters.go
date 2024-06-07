package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/application_parameters"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForApplicationParametersModel(v models.ApplicationParameters) ghcmessages.ApplicationParameters {

	parameterValue := v.ParameterValue
	parameterName := v.ParameterName

	payload := ghcmessages.ApplicationParameters{
		ParameterValue: parameterValue,
		ParameterName:  parameterName,
	}
	return payload
}

// ApplicationParametersValidateHandler validates a value provided by the service member
type ApplicationParametersParamHandler struct {
	handlers.HandlerConfig
}

// Handler receives a GET request containing a parameter name
// if the name is present, it returns the value back, if not, it returns an empty object
func (h ApplicationParametersParamHandler) Handle(params application_parameters.GetParamParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// receive the value
			name := params.ParameterName

			// fetch the value, if not found it will be an empty string
			result, _ := models.FetchParameterValueByName(appCtx.DB(), name)

			parameterValuePayload := payloadForApplicationParametersModel(result)

			return application_parameters.NewGetParamOK().WithPayload(&parameterValuePayload), nil
		})
}
