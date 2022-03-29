package internalapi

import (
	"regexp"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	postalcodesops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/postal_codes"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ValidatePostalCodeWithRateDataHandler has the service validator
type ValidatePostalCodeWithRateDataHandler struct {
	handlers.HandlerContext
	validatePostalCode services.PostalCodeValidator
}

// Handle should call the service validator and rescue expected errors and return false to valid
func (h ValidatePostalCodeWithRateDataHandler) Handle(params postalcodesops.ValidatePostalCodeWithRateDataParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			postalCode := params.PostalCode
			postalCodeType := params.PostalCodeType

			valid, err := h.validatePostalCode.ValidatePostalCode(appCtx,
				postalCode,
				services.PostalCodeType(postalCodeType),
			)
			latLongErrorRegex := regexp.MustCompile("Unsupported postal code lookup")

			if err != nil {
				switch {
				case latLongErrorRegex.MatchString(err.Error()):
					appCtx.Logger().Error("We don't have latlong for postal code", zap.Error(err))
				case err == models.ErrFetchNotFound && postalCodeType == "origin":
					appCtx.Logger().Error("We do not have rate area data for origin postal code", zap.Error(err))
				case err == models.ErrFetchNotFound && postalCodeType == "destination":
					appCtx.Logger().Error("We do not have region rate data for destination postal code", zap.Error(err))
				default:
					appCtx.Logger().Error("Validate postal code", zap.Error(err))
					return postalcodesops.NewValidatePostalCodeWithRateDataBadRequest(), err
				}
			}

			payload := internalmessages.RateEnginePostalCodePayload{
				Valid:          &valid,
				PostalCode:     &postalCode,
				PostalCodeType: &postalCodeType,
			}
			return postalcodesops.NewValidatePostalCodeWithRateDataOK().WithPayload(&payload), nil
		})
}
