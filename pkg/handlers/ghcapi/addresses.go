package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	addressop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/addresses"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

type GetLocationByZipCityHandler struct {
	handlers.HandlerConfig
	services.VLocation
}

func (h GetLocationByZipCityHandler) Handle(params addressop.GetLocationByZipCityParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeApp() && appCtx.Session().OfficeUserID == uuid.Nil {
				noOfficeUserIDErr := apperror.NewSessionError("No office user ID")
				return addressop.NewGetLocationByZipCityForbidden(), noOfficeUserIDErr
			}

			locationList, err := h.GetLocationsByZipCity(appCtx, params.Search)
			if err != nil {
				appCtx.Logger().Error("Error searching for Zip/City/State: ", zap.Error(err))
				return addressop.NewGetLocationByZipCityInternalServerError(), err
			}

			returnPayload := payloads.VLocations(*locationList)
			return addressop.NewGetLocationByZipCityOK().WithPayload(returnPayload), nil
		})
}
