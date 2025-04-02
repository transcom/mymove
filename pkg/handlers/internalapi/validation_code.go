package internalapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	vcodeops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/validation_code"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// ValidationCodeValidationCodeHandler takes a customer provided
// validation code and returns if it is active or not
type ValidationCodeValidationCodeHandler struct {
	handlers.HandlerConfig
}

// Handler receives a POST request containing a parameter value
// if the value is present, it returns it back, if not, it returns an empty object
func (h ValidationCodeValidationCodeHandler) Handle(params vcodeops.ValidateCodeParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// we are only allowing this to be called from the customer app
			// since this is an open route outside of auth, we want to buckle down on validation here
			if !appCtx.Session().IsMilApp() {
				return vcodeops.NewValidateCodeUnauthorized(), apperror.NewSessionError("Request is not from the customer application")
			}

			// receive the value
			if params.Body.ValidationCode == nil || *params.Body.ValidationCode == "" {
				return vcodeops.NewValidateCodeBadRequest(), apperror.NewSessionError("Validation code must be provided")
			}

			// fetch the value, if not found it will be an empty string
			result, _ := models.FetchParameterValue(appCtx.DB(), "validation_code", *params.Body.ValidationCode)

			parameterValuePayload := payloadForApplicationParametersModel(result)

			return vcodeops.NewValidateCodeOK().WithPayload(&parameterValuePayload), nil
		})
}
