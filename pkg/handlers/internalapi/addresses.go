package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	addressop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/addresses"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func addressModelFromPayload(rawAddress *internalmessages.Address) *models.Address {
	if rawAddress == nil {
		return nil
	}
	if rawAddress.County == nil {
		rawAddress.County = models.StringPointer("")
	}

	usPostRegionCitiesID := uuid.FromStringOrNil(rawAddress.UsPostRegionCitiesID.String())

	return &models.Address{
		StreetAddress1:     *rawAddress.StreetAddress1,
		StreetAddress2:     rawAddress.StreetAddress2,
		StreetAddress3:     rawAddress.StreetAddress3,
		City:               *rawAddress.City,
		State:              *rawAddress.State,
		PostalCode:         *rawAddress.PostalCode,
		County:             rawAddress.County,
		UsPostRegionCityId: &usPostRegionCitiesID,
	}
}

func updateAddressWithPayload(a *models.Address, payload *internalmessages.Address) {
	a.StreetAddress1 = *payload.StreetAddress1
	a.StreetAddress2 = payload.StreetAddress2
	a.StreetAddress3 = payload.StreetAddress3
	a.City = *payload.City
	a.State = *payload.State
	a.PostalCode = *payload.PostalCode
	usPostRegionCitiesID := uuid.FromStringOrNil(payload.UsPostRegionCitiesID.String())
	a.UsPostRegionCityId = &usPostRegionCitiesID
	if payload.County == nil {
		a.County = nil
	} else {
		a.County = payload.County
	}
}

// ShowAddressHandler returns an address
type ShowAddressHandler struct {
	handlers.HandlerConfig
}

// Handle returns a address given an addressId
func (h ShowAddressHandler) Handle(params addressop.ShowAddressParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			addressID, err := uuid.FromString(params.AddressID.String())

			if err != nil {
				appCtx.Logger().Error("Finding address", zap.Error(err))
			}
			address := models.FetchAddressByID(appCtx.DB(), &addressID)

			addressPayload := payloads.Address(address)
			return addressop.NewShowAddressOK().WithPayload(addressPayload), nil
		})
}

type GetLocationByZipCityStateHandler struct {
	handlers.HandlerConfig
	services.VLocation
}

func (h GetLocationByZipCityStateHandler) Handle(params addressop.GetLocationByZipCityStateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return addressop.NewGetLocationByZipCityStateForbidden(), noServiceMemberIDErr
			}

			locationList, err := h.GetLocationsByZipCityState(appCtx, params.Search)
			if err != nil {
				appCtx.Logger().Error("Error searching for Zip/City/State: ", zap.Error(err))
				return addressop.NewGetLocationByZipCityStateInternalServerError(), err
			}

			returnPayload := payloads.VLocations(*locationList)
			return addressop.NewGetLocationByZipCityStateOK().WithPayload(returnPayload), nil
		})
}
