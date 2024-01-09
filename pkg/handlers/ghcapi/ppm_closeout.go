package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmcloseoutops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

type GetPPMCloseoutHandler struct {
	handlers.HandlerConfig
	ppmCloseoutFetcher services.PPMCloseout
}

// Handle handles the handling of fetching a single MTO shipment by ID.
func (h GetPPMCloseoutHandler) Handle(params ppmcloseoutops.FetchCloseoutCalculationsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetPPMCloseout error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return ppmcloseoutops.NewFetchCloseoutCalculationsNotFound().WithPayload(payload), err
				case apperror.QueryError:
					return ppmcloseoutops.NewFetchCloseoutCalculationsInternalServerError(), err
				default:
					return ppmcloseoutops.NewFetchCloseoutCalculationsInternalServerError(), err
				}
			}

			// eagerAssociations := []string{"MoveTaskOrder",
			// 	"PickupAddress",
			// 	"DestinationAddress",
			// 	"SecondaryPickupAddress",
			// 	"SecondaryDeliveryAddress",
			// 	"MTOServiceItems.CustomerContacts",
			// 	"StorageFacility.Address",
			// 	"PPMShipment"}

			// shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())

			// mtoShipment, err := h.mtoShipmentFetcher.GetShipment(appCtx, shipmentID, eagerAssociations...)
			// if err != nil {
			// 	return handleError(err)
			// }

			// var agents []models.MTOAgent
			// err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Where("mto_shipment_id = ?", mtoShipment.ID).All(&agents)
			// if err != nil {
			// 	return handleError(err)
			// }
			// mtoShipment.MTOAgents = agents
			// payload := payloads.MTOShipment(h.FileStorer(), mtoShipment, nil)
			return handleError(nil)
		})
}