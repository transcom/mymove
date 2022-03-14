package internalapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

//
// CREATE
//

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentCreator services.MTOShipmentCreator
	ppmShipmentCreator services.PPMShipmentCreator
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if appCtx.Session() == nil || (!appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil) {
				return mtoshipmentops.NewCreateMTOShipmentUnauthorized()
			}

			payload := params.Body
			if payload == nil {
				appCtx.Logger().Error("Invalid mto shipment: params Body is nil")
				return mtoshipmentops.NewCreateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The MTO Shipment request body cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest)))
			}
			mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
			var err error
			var ppmShipment *models.PPMShipment
			if payload.ShipmentType != nil && *payload.ShipmentType == internalmessages.MTOShipmentTypePPM {
				// Return a PPM Shipment with an MTO Shipment inside
				ppmShipment, err = h.ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, mtoShipment.PPMShipment)
			} else {
				// TODO: remove this status change once MB-3428 is implemented and can update to Submitted on second page
				mtoShipment.Status = models.MTOShipmentStatusSubmitted
				serviceItemsList := make(models.MTOServiceItems, 0)
				mtoShipment, err = h.mtoShipmentCreator.CreateMTOShipment(appCtx, mtoShipment, serviceItemsList)
			}

			if err != nil {
				appCtx.Logger().Error("internalapi.CreateMTOShipmentHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
				case apperror.InvalidInputError:
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors))
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("internalapi.CreateMTOServiceItemHandler error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
				default:
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
				}
			}

			if payload.ShipmentType != nil && *payload.ShipmentType == internalmessages.MTOShipmentTypePPM {
				// Return an mtoShipment that has a ppmShipment
				mtoShipment = &ppmShipment.Shipment
			}
			returnPayload := payloads.MTOShipment(mtoShipment)
			return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload)
		})
}

//
// UPDATE
//

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentUpdater services.MTOShipmentUpdater
	ppmShipmentUpdater services.PPMShipmentUpdater
}

// Handle updates the mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if appCtx.Session() == nil {
				return mtoshipmentops.NewUpdateMTOShipmentUnauthorized()
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				return mtoshipmentops.NewUpdateMTOShipmentForbidden()
			}

			payload := params.Body
			if payload == nil {
				appCtx.Logger().Error("Invalid mto shipment: params Body is nil")
				return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The MTO Shipment request body cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest)))
			}

			mtoShipment := payloads.MTOShipmentModelFromUpdate(payload)
			mtoShipment.ID = uuid.FromStringOrNil(params.MtoShipmentID.String())

			status := mtoShipment.Status
			if status != "" && status != models.MTOShipmentStatusDraft && status != models.MTOShipmentStatusSubmitted {
				appCtx.Logger().Error("Invalid mto shipment status: shipment in service member app can only have draft or submitted status")

				return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(
					payloads.ClientError(handlers.BadRequestErrMessage,
						"When present, the MTO Shipment status must be one of: DRAFT or SUBMITTED.",
						h.GetTraceIDFromRequest(params.HTTPRequest)))
			}

			var updatedMTOShipment *models.MTOShipment
			var updatedPPMShipment *models.PPMShipment
			var err error
			// We should move this logic out of the handler and into a composable service object
			err = appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
				updatedMTOShipment, err = h.mtoShipmentUpdater.UpdateMTOShipmentCustomer(txnAppCtx, mtoShipment, params.IfMatch)
				if err != nil {
					return err
				}

				updatedPPMShipment, err = h.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(txnAppCtx, mtoShipment.PPMShipment, mtoShipment.ID)
				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				appCtx.Logger().Error("internalapi.UpdateMTOShipmentHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors))
				case apperror.PreconditionFailedError:
					return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("internalapi.UpdateMTOServiceItemHandler error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
				default:
					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
				}
			}

			updatedMTOShipment.PPMShipment = updatedPPMShipment
			returnPayload := payloads.MTOShipment(updatedMTOShipment)

			return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(returnPayload)
		})
}

//
// GET ALL
//

// ListMTOShipmentsHandler returns a list of MTO Shipments
type ListMTOShipmentsHandler struct {
	handlers.HandlerContext
	services.ListFetcher
	services.Fetcher
}

// Handle listing mto shipments for the move task order
func (h ListMTOShipmentsHandler) Handle(params mtoshipmentops.ListMTOShipmentsParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if appCtx.Session() == nil || (!appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil) {
				return mtoshipmentops.NewListMTOShipmentsUnauthorized()
			}

			moveTaskOrderID, err := uuid.FromString(params.MoveTaskOrderID.String())
			// return any parsing error
			if err != nil {
				appCtx.Logger().Error("Invalid request: move task order ID not valid")
				return mtoshipmentops.NewListMTOShipmentsBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The MTO Shipment request body cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest)))
			}

			// check if move task order exists first
			queryFilters := []services.QueryFilter{
				query.NewQueryFilter("id", "=", moveTaskOrderID.String()),
				query.NewQueryFilter("show", "=", "TRUE"),
			}

			moveTaskOrder := &models.Move{}
			err = h.Fetcher.FetchRecord(appCtx, moveTaskOrder, queryFilters)
			if err != nil {
				appCtx.Logger().Error("Error fetching move task order: ", zap.Error(fmt.Errorf("Move Task Order ID: %s", moveTaskOrder.ID)), zap.Error(err))
				return mtoshipmentops.NewListMTOShipmentsNotFound()
			}

			queryFilters = []services.QueryFilter{
				query.NewQueryFilter("move_id", "=", moveTaskOrderID.String()),
			}

			// TODO: In some places, we used this unbound eager call accidentally and loaded all associations when the
			//   intention was to load no associations. In this instance, we get E2E failures if we change this to load
			//   no associations, so we'll keep it as is and can revisit later if we want to optimize further.  This is
			//   just loading shipments for a specific move (likely only 1 or 2 in most cases), so the impact of the
			//   additional loading shouldn't be too dramatic.
			queryAssociations := query.NewQueryAssociations([]services.QueryAssociation{})

			queryOrder := query.NewQueryOrder(swag.String("created_at"), swag.Bool(true))

			var shipments models.MTOShipments
			err = h.ListFetcher.FetchRecordList(appCtx, &shipments, queryFilters, queryAssociations, nil, queryOrder)
			// return any errors
			if err != nil {
				appCtx.Logger().Error("Error fetching mto shipments : ", zap.Error(err))

				return mtoshipmentops.NewListMTOShipmentsInternalServerError()
			}

			payload := payloads.MTOShipments(&shipments)
			return mtoshipmentops.NewListMTOShipmentsOK().WithPayload(*payload)
		})
}
