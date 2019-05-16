package publicapi

import (
	"regexp"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	postalcodesops "github.com/transcom/mymove/pkg/gen/restapi/apioperations/postal_codes"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ValidatePostalCodeHandler has the service validator
type ValidatePostalCodeHandler struct {
	handlers.HandlerContext
	validatePostalCode services.PostalCodeValidator
}

// Handle should call the service validator and rescue expected errors and return false to valid
func (h ValidatePostalCodeHandler) Handle(params postalcodesops.ValidatePostalCodeParams) middleware.Responder {
	postalCode := params.PostalCode
	postalCodeType := params.PostalCodeType

	valid, err := h.validatePostalCode.ValidatePostalCode(
		postalCode,
		services.PostalCodeType(postalCodeType),
	)
	latLongErrorRegex := regexp.MustCompile("Unsupported postal code lookup")

	if err != nil {
		if latLongErrorRegex.MatchString(err.Error()) {
			h.Logger().Error("We don't have latlong for postal code", zap.Error(err))
		} else if err == models.ErrFetchNotFound && postalCodeType == "origin" {
			h.Logger().Error("We do not have rate area data for origin postal code", zap.Error(err))
		} else if err == models.ErrFetchNotFound && postalCodeType == "destination" {
			h.Logger().Error("We do not have region rate data for destination postal code", zap.Error(err))
		} else {
			h.Logger().Error("Validate postal code", zap.Error(err))
			return postalcodesops.NewValidatePostalCodeBadRequest()
		}
	}

	payload := apimessages.ValidatePostalCodePayload{
		Valid:          &valid,
		PostalCode:     &postalCode,
		PostalCodeType: &postalCodeType,
	}
	return postalcodesops.NewValidatePostalCodeOK().WithPayload(&payload)
}
