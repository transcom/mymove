package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	shipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/shipment"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/event"
	"github.com/transcom/mymove/pkg/services/featureflag"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	ppmshipment "github.com/transcom/mymove/pkg/services/ppmshipment"
)

// ListMTOShipmentsHandler returns a list of MTO Shipments
type ListMTOShipmentsHandler struct {
	handlers.HandlerConfig
	services.MTOShipmentFetcher
	services.ShipmentSITStatus
}

// Handle listing mto shipments for the move task order
func (h ListMTOShipmentsHandler) Handle(params mtoshipmentops.ListMTOShipmentsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ListMTOShipmentsHandler error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewListMTOShipmentsNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return mtoshipmentops.NewListMTOShipmentsForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return mtoshipmentops.NewListMTOShipmentsInternalServerError(), err
				default:
					return mtoshipmentops.NewListMTOShipmentsInternalServerError(), err
				}
			}

			moveID := uuid.FromStringOrNil(params.MoveTaskOrderID.String())

			shipments, err := h.ListMTOShipments(appCtx, moveID)
			if err != nil {
				return handleError(err)
			}

			shipmentSITStatuses := h.CalculateShipmentsSITStatuses(appCtx, shipments)

			/** Feature Flag - Boat Shipment **/
			featureFlagName := "boat"
			isBoatFeatureOn := false
			flag, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
				isBoatFeatureOn = false
			} else {
				isBoatFeatureOn = flag.Match
			}

			// Remove Boat shipments if Boat FF is off
			if !isBoatFeatureOn {
				var filteredShipments models.MTOShipments
				if shipments != nil {
					filteredShipments = models.MTOShipments{}
				}
				for i, shipment := range shipments {
					if shipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway || shipment.ShipmentType == models.MTOShipmentTypeBoatTowAway {
						continue
					}

					filteredShipments = append(filteredShipments, shipments[i])
				}
				shipments = filteredShipments
			}
			/** End of Feature Flag **/

			/** Feature Flag - Mobile Home Shipment **/
			featureFlagNameMH := "mobile_home"
			isMobileHomeFeatureOn := false
			flagMH, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagNameMH, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureFlagNameMH), zap.Error(err))
				isMobileHomeFeatureOn = false
			} else {
				isMobileHomeFeatureOn = flagMH.Match
			}

			// Remove Mobile Home shipments if Mobile Home FF is off
			if !isMobileHomeFeatureOn {
				var filteredShipments models.MTOShipments
				if shipments != nil {
					filteredShipments = models.MTOShipments{}
				}
				for i, shipment := range shipments {
					if shipment.ShipmentType == models.MTOShipmentTypeMobileHome {
						continue
					}

					filteredShipments = append(filteredShipments, shipments[i])
				}
				shipments = filteredShipments
			}
			/** End of Feature Flag **/

			sitStatusPayload := payloads.SITStatuses(shipmentSITStatuses, h.FileStorer())
			payload := payloads.MTOShipments(h.FileStorer(), (*models.MTOShipments)(&shipments), sitStatusPayload)
			return mtoshipmentops.NewListMTOShipmentsOK().WithPayload(*payload), nil
		})
}

// GetMTOShipmentHandler is the handler to fetch a single MTO shipment by ID
type GetMTOShipmentHandler struct {
	handlers.HandlerConfig
	mtoShipmentFetcher services.MTOShipmentFetcher
}

// Handle handles the handling of fetching a single MTO shipment by ID.
func (h GetMTOShipmentHandler) Handle(params mtoshipmentops.GetShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetShipment error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewGetShipmentNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return mtoshipmentops.NewGetShipmentForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return mtoshipmentops.NewGetShipmentInternalServerError(), err
				default:
					return mtoshipmentops.NewGetShipmentInternalServerError(), err
				}
			}

			eagerAssociations := []string{"MoveTaskOrder",
				"PickupAddress",
				"DestinationAddress",
				"SecondaryPickupAddress",
				"SecondaryDeliveryAddress",
				"MTOServiceItems.CustomerContacts",
				"StorageFacility.Address",
				"PPMShipment",
				"BoatShipment",
				"MobileHome",
				"Distance"}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())

			mtoShipment, err := h.mtoShipmentFetcher.GetShipment(appCtx, shipmentID, eagerAssociations...)
			if err != nil {
				return handleError(err)
			}

			if mtoShipment.ShipmentType == models.MTOShipmentTypePPM {
				ppmEagerAssociations := []string{"PickupAddress",
					"DestinationAddress",
					"SecondaryPickupAddress",
					"SecondaryDestinationAddress",
				}

				ppmShipmentFetcher := ppmshipment.NewPPMShipmentFetcher()

				ppmShipment, err := ppmShipmentFetcher.GetPPMShipment(appCtx, mtoShipment.PPMShipment.ID, ppmEagerAssociations, nil)
				if err != nil {
					return handleError(err)
				}

				mtoShipment.PPMShipment.PickupAddress = ppmShipment.PickupAddress
				mtoShipment.PPMShipment.DestinationAddress = ppmShipment.DestinationAddress
				mtoShipment.PPMShipment.SecondaryPickupAddress = ppmShipment.SecondaryPickupAddress
				mtoShipment.PPMShipment.SecondaryDestinationAddress = ppmShipment.SecondaryDestinationAddress
			}

			var agents []models.MTOAgent
			err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Where("mto_shipment_id = ?", mtoShipment.ID).All(&agents)
			if err != nil {
				return handleError(err)
			}
			mtoShipment.MTOAgents = agents
			payload := payloads.MTOShipment(h.FileStorer(), mtoShipment, nil)
			return mtoshipmentops.NewGetShipmentOK().WithPayload(payload), nil
		})
}

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerConfig
	shipmentCreator services.ShipmentCreator
	shipmentStatus  services.ShipmentSITStatus
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body

			if payload == nil {
				invalidShipmentError := apperror.NewBadDataError("Invalid mto shipment: params Body is nil")
				appCtx.Logger().Error(invalidShipmentError.Error())
				return mtoshipmentops.NewCreateMTOShipmentBadRequest(), invalidShipmentError
			}

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.CreateMTOShipmentHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(&payload), err
				case apperror.EventError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(&payload), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"CreateMTOShipment",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payload), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.CreateMTOShipmentHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError(), err
				default:
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError(), err
				}
			}

			mtoShipment := payloads.MTOShipmentModelFromCreate(payload)

			if mtoShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom && mtoShipment.NTSRecordedWeight != nil {
				previouslyRecordedWeight := *mtoShipment.NTSRecordedWeight
				mtoShipment.PrimeEstimatedWeight = &previouslyRecordedWeight
			}

			var err error
			mtoShipment, err = h.shipmentCreator.CreateShipment(appCtx, mtoShipment)

			if err != nil {
				return handleError(err)
			}

			if mtoShipment == nil {
				shipmentNotCreatedError := apperror.NewInternalServerError("Unexpected nil shipment from CreateMTOShipment")
				appCtx.Logger().Error(shipmentNotCreatedError.Error())
				return mtoshipmentops.NewCreateMTOShipmentInternalServerError(), shipmentNotCreatedError
			}

			sitAllowance, err := h.shipmentStatus.CalculateShipmentSITAllowance(appCtx, *mtoShipment)
			if err != nil {
				return handleError(err)
			}

			mtoShipment.SITDaysAllowance = &sitAllowance

			returnPayload := payloads.MTOShipment(h.FileStorer(), mtoShipment, nil)
			return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

// UpdateShipmentHandler updates shipments
type UpdateShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentUpdater
	services.ShipmentSITStatus
}

// Handle updates shipments
func (h UpdateShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			if payload == nil {
				appCtx.Logger().Error("Invalid mto shipment: params Body is nil")
				emptyBodyError := apperror.NewBadDataError("The MTO Shipment request body cannot be empty.")
				payload := payloadForValidationError(
					"Empty body error",
					emptyBodyError.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors(),
				)

				return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload), emptyBodyError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			oldShipment, err := mtoshipment.FindShipment(appCtx, shipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound(), err
				case apperror.EventError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(&payload), err
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				}
			}

			mtoShipment := payloads.MTOShipmentModelFromUpdate(payload)
			mtoShipment.ID = shipmentID
			isBoatShipment := mtoShipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway || mtoShipment.ShipmentType == models.MTOShipmentTypeBoatTowAway
			if !isBoatShipment {
				mtoShipment.ShipmentType = oldShipment.ShipmentType
			}

			//MTOShipmentModelFromUpdate defaults UsesExternalVendor to false if it's nil in the payload
			if payload.UsesExternalVendor == nil {
				mtoShipment.UsesExternalVendor = oldShipment.UsesExternalVendor
			}
			// booleans not passed will update to false
			mtoShipment.Diversion = oldShipment.Diversion

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound(), err
				case apperror.EventError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(&payload), err
				case apperror.ForbiddenError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return mtoshipmentops.NewUpdateMTOShipmentForbidden().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.UpdateShipmentHandler error", zap.Error(e.Unwrap()))
					}

					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				}
			}

			if mtoShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom && mtoShipment.NTSRecordedWeight != nil {
				previouslyRecordedWeight := *mtoShipment.NTSRecordedWeight
				mtoShipment.PrimeEstimatedWeight = &previouslyRecordedWeight
			}

			updatedMtoShipment, err := h.ShipmentUpdater.UpdateShipment(appCtx, mtoShipment, params.IfMatch, "ghc")
			if err != nil {
				return handleError(err)
			}

			_, err = event.TriggerEvent(event.Event{
				EndpointKey: event.GhcUpdateMTOShipmentEndpointKey,
				// Endpoint that is being handled
				EventKey:        event.MTOShipmentUpdateEventKey,    // Event that you want to trigger
				UpdatedObjectID: updatedMtoShipment.ID,              // ID of the updated logical object
				MtoID:           updatedMtoShipment.MoveTaskOrderID, // ID of the associated Move
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			// If the event trigger fails, just log the error.
			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMTOShipment could not generate the event")
			}

			shipmentSITStatus, _, err := h.CalculateShipmentSITStatus(appCtx, *updatedMtoShipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			returnPayload := payloads.MTOShipment(h.FileStorer(), updatedMtoShipment, sitStatusPayload)
			return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

// DeleteShipmentHandler soft deletes a shipment
type DeleteShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentDeleter
}

// Handle soft deletes a shipment
func (h DeleteShipmentHandler) Handle(params shipmentops.DeleteShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) && !appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				forbiddenError := apperror.NewForbiddenError("user is not authenticated with an office role")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewDeleteShipmentForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			moveID, err := h.DeleteShipment(appCtx, shipmentID)
			if err != nil {
				appCtx.Logger().Error("ghcapi.DeleteShipmentHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewDeleteShipmentNotFound(), err
				case apperror.ConflictError:
					return shipmentops.NewDeleteShipmentConflict(), err
				case apperror.ForbiddenError:
					return shipmentops.NewDeleteShipmentForbidden(), err
				case apperror.UnprocessableEntityError:
					return shipmentops.NewDeleteShipmentUnprocessableEntity(), err
				default:
					return shipmentops.NewDeleteShipmentInternalServerError(), err
				}
			}

			// Note that this is currently not sending any notifications because
			// the move isn't available to the Prime yet. See the objectEventHandler
			// function in pkg/services/event/notification.go.
			// We added this now because eventually, we will want to save events in
			// the DB for auditing purposes. When that happens, this code in the handler
			// will not change. However, we should make sure to add a test in
			// mto_shipment_test.go that verifies the audit got saved.
			h.triggerShipmentDeletionEvent(appCtx, shipmentID, moveID, params)

			return shipmentops.NewDeleteShipmentNoContent(), nil
		})
}

func (h DeleteShipmentHandler) triggerShipmentDeletionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.DeleteShipmentParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcDeleteShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentDeleteEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                   // ID of the updated logical object
		MtoID:           moveID,                       // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.DeleteShipmentHandler could not generate the event", zap.Error(err))
	}
}

// ApproveShipmentHandler approves a shipment
type ApproveShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentApprover
	services.ShipmentSITStatus
}

// Handle approves a shipment
func (h ApproveShipmentHandler) Handle(params shipmentops.ApproveShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can approve shipments")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewApproveShipmentForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.ApproveShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewApproveShipmentNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Validation errors", "ApproveShipment", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return shipmentops.NewApproveShipmentUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewApproveShipmentPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ConflictError, mtoshipment.ConflictStatusError:
					return shipmentops.NewApproveShipmentConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewApproveShipmentInternalServerError(), err
				}
			}

			featureFlagValues, err := handlers.GetAllDomesticMHFlags(appCtx, h.HandlerConfig.FeatureFlagFetcher())
			if err != nil {
				return handleError(err)
			}
			shipment, err := h.ApproveShipment(appCtx, shipmentID, eTag, featureFlagValues)
			if err != nil {
				return handleError(err)
			}

			h.triggerShipmentApprovalEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			shipmentSITStatus, _, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			payload := payloads.MTOShipment(h.FileStorer(), shipment, sitStatusPayload)
			return shipmentops.NewApproveShipmentOK().WithPayload(payload), nil
		})
}

func (h ApproveShipmentHandler) triggerShipmentApprovalEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.ApproveShipmentParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcApproveShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentApproveEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                    // ID of the updated logical object
		MtoID:           moveID,                        // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.ApproveShipmentHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentDiversionHandler Requests a shipment diversion
type RequestShipmentDiversionHandler struct {
	handlers.HandlerConfig
	services.ShipmentDiversionRequester
	services.ShipmentSITStatus
}

// Handle Requests a shipment diversion
func (h RequestShipmentDiversionHandler) Handle(params shipmentops.RequestShipmentDiversionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can Request shipment diversions")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewRequestShipmentDiversionForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch
			diversionReason := params.Body.DiversionReason

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.RequestShipmentDiversionHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRequestShipmentDiversionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"RequestShipmentDiversion",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewRequestShipmentDiversionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewRequestShipmentDiversionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewRequestShipmentDiversionConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewRequestShipmentDiversionInternalServerError(), err
				}
			}

			shipment, err := h.RequestShipmentDiversion(appCtx, shipmentID, eTag, diversionReason)
			if err != nil {
				return handleError(err)
			}

			h.triggerRequestShipmentDiversionEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			shipmentSITStatus, _, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			payload := payloads.MTOShipment(h.FileStorer(), shipment, sitStatusPayload)
			return shipmentops.NewRequestShipmentDiversionOK().WithPayload(payload), nil
		})
}

func (h RequestShipmentDiversionHandler) triggerRequestShipmentDiversionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentDiversionParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentDiversionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestDiversionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                             // ID of the updated logical object
		MtoID:           moveID,                                 // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.RequestShipmentDiversionHandler could not generate the event", zap.Error(err))
	}
}

// ApproveShipmentDiversionHandler approves a shipment diversion
type ApproveShipmentDiversionHandler struct {
	handlers.HandlerConfig
	services.ShipmentDiversionApprover
	services.ShipmentSITStatus
}

// Handle approves a shipment diversion
func (h ApproveShipmentDiversionHandler) Handle(params shipmentops.ApproveShipmentDiversionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can approve shipment diversions")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewApproveShipmentDiversionForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.ApproveShipmentDiversionHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewApproveShipmentDiversionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"ApproveShipmentDiversion",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewApproveShipmentDiversionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewApproveShipmentDiversionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewApproveShipmentDiversionConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewApproveShipmentDiversionInternalServerError(), err
				}
			}

			shipment, err := h.ApproveShipmentDiversion(appCtx, shipmentID, eTag)
			if err != nil {
				return handleError(err)
			}

			h.triggerShipmentDiversionApprovalEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			shipmentSITStatus, _, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			payload := payloads.MTOShipment(h.FileStorer(), shipment, sitStatusPayload)
			return shipmentops.NewApproveShipmentDiversionOK().WithPayload(payload), nil
		})
}

func (h ApproveShipmentDiversionHandler) triggerShipmentDiversionApprovalEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.ApproveShipmentDiversionParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcApproveShipmentDiversionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentApproveDiversionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                             // ID of the updated logical object
		MtoID:           moveID,                                 // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.ApproveShipmentDiversionHandler could not generate the event", zap.Error(err))
	}
}

// RejectShipmentHandler rejects a shipment
type RejectShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentRejecter
}

// Handle rejects a shipment
func (h RejectShipmentHandler) Handle(params shipmentops.RejectShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can reject shipments")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewRejectShipmentForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch
			rejectionReason := params.Body.RejectionReason
			shipment, err := h.RejectShipment(appCtx, shipmentID, eTag, rejectionReason)

			if err != nil {
				appCtx.Logger().Error("ghcapi.RejectShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRejectShipmentNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"RejectShipment",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewRejectShipmentUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewRejectShipmentPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewRejectShipmentConflict().
							WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}),
						err
				default:
					return shipmentops.NewRejectShipmentInternalServerError(), err
				}
			}

			h.triggerShipmentRejectionEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			payload := payloads.MTOShipment(h.FileStorer(), shipment, nil)
			return shipmentops.NewRejectShipmentOK().WithPayload(payload), nil
		})
}

func (h RejectShipmentHandler) triggerShipmentRejectionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RejectShipmentParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRejectShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRejectEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                   // ID of the updated logical object
		MtoID:           moveID,                       // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.RejectShipmentHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentCancellationHandler Requests a shipment cancellation
type RequestShipmentCancellationHandler struct {
	handlers.HandlerConfig
	services.ShipmentCancellationRequester
	services.ShipmentSITStatus
}

// Handle Requests a shipment cancellation
func (h RequestShipmentCancellationHandler) Handle(params shipmentops.RequestShipmentCancellationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can Request shipment diversions")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewRequestShipmentCancellationForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.RequestShipmentCancellationHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRequestShipmentCancellationNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"RequestShipmentCancellation",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewRequestShipmentCancellationUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewRequestShipmentCancellationPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.UpdateError:
					payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
					return shipmentops.NewRequestShipmentCancellationConflict().WithPayload(payload), err
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewRequestShipmentCancellationConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewRequestShipmentCancellationInternalServerError(), err
				}
			}

			shipment, err := h.RequestShipmentCancellation(appCtx, shipmentID, eTag)
			if err != nil {
				return handleError(err)
			}

			h.triggerRequestShipmentCancellationEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			shipmentSITStatus, _, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			payload := payloads.MTOShipment(h.FileStorer(), shipment, sitStatusPayload)
			return shipmentops.NewRequestShipmentCancellationOK().WithPayload(payload), nil
		})
}

func (h RequestShipmentCancellationHandler) triggerRequestShipmentCancellationEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentCancellationParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentCancellationEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestCancellationEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                                // ID of the updated logical object
		MtoID:           moveID,                                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.RequestShipmentCancellationHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentReweighHandler Requests a shipment reweigh
type RequestShipmentReweighHandler struct {
	handlers.HandlerConfig
	services.ShipmentReweighRequester
	services.ShipmentSITStatus
	services.MTOShipmentUpdater
}

// Handle Requests a shipment reweigh
func (h RequestShipmentReweighHandler) Handle(params shipmentops.RequestShipmentReweighParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can Request a shipment reweigh")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewRequestShipmentReweighForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			reweigh, err := h.RequestShipmentReweigh(appCtx, shipmentID, models.ReweighRequesterTOO)

			if err != nil {
				appCtx.Logger().Error("ghcapi.RequestShipmentReweighHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRequestShipmentReweighNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors", "RequestShipmentReweigh",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewRequestShipmentReweighUnprocessableEntity().WithPayload(payload), err
				case apperror.ConflictError:
					return shipmentops.NewRequestShipmentReweighConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewRequestShipmentReweighInternalServerError(), err
				}
			}

			handleError := func(err error) middleware.Responder {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound()
				case apperror.EventError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(&payload)
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					)
				}
			}

			shipment, err := mtoshipment.FindShipment(appCtx, shipmentID)
			if err != nil {
				return handleError(err), err
			}

			moveID := shipment.MoveTaskOrderID
			h.triggerRequestShipmentReweighEvent(appCtx, shipmentID, moveID, params)

			move, err := models.FetchMoveByMoveIDWithOrders(appCtx.DB(), shipment.MoveTaskOrderID)
			if err != nil {
				return nil, err
			}

			/* Don't send emails for BLUEBARK/SAFETY moves */
			if move.Orders.CanSendEmailWithOrdersType() {
				err = h.NotificationSender().SendNotification(appCtx,
					notifications.NewReweighRequested(moveID, *shipment),
				)
				if err != nil {
					appCtx.Logger().Error("problem sending email to user", zap.Error(err))
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
			}

			shipmentSITStatus, _, err := h.CalculateShipmentSITStatus(appCtx, reweigh.Shipment)
			if err != nil {
				return handleError(err), err
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			payload := payloads.Reweigh(reweigh, sitStatusPayload)
			return shipmentops.NewRequestShipmentReweighOK().WithPayload(payload), nil
		})
}

func (h RequestShipmentReweighHandler) triggerRequestShipmentReweighEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentReweighParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentReweighEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestReweighEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                           // ID of the updated logical object
		MtoID:           moveID,                               // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.RequestShipmentReweighHandler could not generate the event", zap.Error(err))
	}
}

// ReviewShipmentAddressUpdateHandler Reviews a shipment address change
type ReviewShipmentAddressUpdateHandler struct {
	handlers.HandlerConfig
	services.ShipmentAddressUpdateRequester
}

// Handle ... reviews address update request
func (h ReviewShipmentAddressUpdateHandler) Handle(params shipmentops.ReviewShipmentAddressUpdateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			featureFlagValues := make(map[string]bool)

			isDomesticMobileHomeDOPOn, err := handlers.GetFeatureFlagValue(appCtx, h.HandlerConfig.FeatureFlagFetcher(), featureflag.DomesticMobileHomeDOPEnabled)
			if err != nil {
				return shipmentops.NewReviewShipmentAddressUpdateInternalServerError(), err
			}
			isDomesticMobileHomeDDPOn, err := handlers.GetFeatureFlagValue(appCtx, h.HandlerConfig.FeatureFlagFetcher(), featureflag.DomesticMobileHomeDDPEnabled)
			if err != nil {
				return shipmentops.NewReviewShipmentAddressUpdateInternalServerError(), err
			}
			isDomesticMobileHomeDPKOn, err := handlers.GetFeatureFlagValue(appCtx, h.HandlerConfig.FeatureFlagFetcher(), featureflag.DomesticMobileHomePackingEnabled)
			if err != nil {
				return shipmentops.NewReviewShipmentAddressUpdateInternalServerError(), err
			}
			isDomesticMobileHomeDUPKOn, err := handlers.GetFeatureFlagValue(appCtx, h.HandlerConfig.FeatureFlagFetcher(), featureflag.DomesticMobileHomeUnpackingEnabled)
			if err != nil {
				return shipmentops.NewReviewShipmentAddressUpdateInternalServerError(), err
			}
			featureFlagValues[featureflag.DomesticMobileHomeDDPEnabled] = isDomesticMobileHomeDDPOn
			featureFlagValues[featureflag.DomesticMobileHomeDOPEnabled] = isDomesticMobileHomeDOPOn
			featureFlagValues[featureflag.DomesticMobileHomePackingEnabled] = isDomesticMobileHomeDPKOn
			featureFlagValues[featureflag.DomesticMobileHomeUnpackingEnabled] = isDomesticMobileHomeDUPKOn

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			addressApprovalStatus := params.Body.Status
			remarks := params.Body.OfficeRemarks

			response, err := h.ShipmentAddressUpdateRequester.ReviewShipmentAddressChange(appCtx, shipmentID, models.ShipmentAddressUpdateStatus(*addressApprovalStatus), *remarks, featureFlagValues)
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.ReviewShipmentAddressUpdateHandler", zap.Error(err))
				payload := ghcmessages.Error{
					Message: handlers.FmtString(err.Error()),
				}

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewReviewShipmentAddressUpdateNotFound().WithPayload(&payload), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"ReviewShipmentAddressUpdate",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewReviewShipmentAddressUpdateUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewReviewShipmentAddressUpdatePreconditionFailed().
						WithPayload(&payload), err
				case apperror.ConflictError:
					return shipmentops.NewReviewShipmentAddressUpdateConflict().
						WithPayload(&payload), err
				default:
					return shipmentops.NewReviewShipmentAddressUpdateInternalServerError(), err
				}
			}
			if err != nil {
				return handleError(err)
			}
			payload := payloads.ShipmentAddressUpdate(response)
			return shipmentops.NewReviewShipmentAddressUpdateOK().WithPayload(payload), nil
		})
}

// ApproveSITExtensionHandler approves a SIT extension
type ApproveSITExtensionHandler struct {
	handlers.HandlerConfig
	services.SITExtensionApprover
	services.ShipmentSITStatus
	services.ShipmentUpdater
}

// Handle ... approves the SIT extension
func (h ApproveSITExtensionHandler) Handle(params shipmentops.ApproveSITExtensionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error approving SIT extension", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewApproveSITExtensionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						handlers.ValidationErrMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewApproveSITExtensionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewApproveSITExtensionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return shipmentops.NewApproveSITExtensionForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewApproveSITExtensionInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("is not a TOO")
				return handleError(forbiddenError)
			}

			shipmentID := uuid.FromStringOrNil(string(params.ShipmentID))
			sitExtensionID := uuid.FromStringOrNil(string(params.SitExtensionID))
			approvedDays := int(*params.Body.ApprovedDays)
			requestReason := models.SITDurationUpdateRequestReason(params.Body.RequestReason)
			officeRemarks := params.Body.OfficeRemarks
			updatedShipment, err := h.SITExtensionApprover.ApproveSITExtension(appCtx, shipmentID, sitExtensionID, approvedDays, requestReason, officeRemarks, params.IfMatch)
			if err != nil {
				return handleError(err)
			}

			shipmentSITStatus, shipmentWithSITInfo, err := h.CalculateShipmentSITStatus(appCtx, *updatedShipment)
			if err != nil {
				return handleError(err)
			}

			existingETag := etag.GenerateEtag(updatedShipment.UpdatedAt)

			updatedShipment, err = h.UpdateShipment(appCtx, &shipmentWithSITInfo, existingETag, "ghc")
			if err != nil {
				return handleError(err)
			}

			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			shipmentPayload := payloads.MTOShipment(h.FileStorer(), updatedShipment, sitStatusPayload)

			h.triggerApproveSITExtensionEvent(appCtx, shipmentID, updatedShipment.MoveTaskOrderID, params)
			return shipmentops.NewApproveSITExtensionOK().WithPayload(shipmentPayload), nil
		})
}

func (h ApproveSITExtensionHandler) triggerApproveSITExtensionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.ApproveSITExtensionParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcApproveSITExtensionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ApproveSITExtensionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                        // ID of the updated logical object
		MtoID:           moveID,                            // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.ApproveSITExtensionHandler could not generate the event", zap.Error(err))
	}
}

// DenySITExtensionHandler denies a SIT extension
type DenySITExtensionHandler struct {
	handlers.HandlerConfig
	services.SITExtensionDenier
	services.ShipmentSITStatus
}

// Handle ... denies the SIT extension
func (h DenySITExtensionHandler) Handle(params shipmentops.DenySITExtensionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error denying SIT extension", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewDenySITExtensionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						handlers.ValidationErrMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewDenySITExtensionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewDenySITExtensionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return shipmentops.NewDenySITExtensionForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewDenySITExtensionInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("is not a TOO")
				return handleError(forbiddenError)
			}

			shipmentID := uuid.FromStringOrNil(string(params.ShipmentID))
			sitExtensionID := uuid.FromStringOrNil(string(params.SitExtensionID))
			officeRemarks := params.Body.OfficeRemarks
			convertToCustomerExpense := params.Body.ConvertToCustomerExpense
			updatedShipment, err := h.SITExtensionDenier.DenySITExtension(appCtx, shipmentID, sitExtensionID, officeRemarks, convertToCustomerExpense, params.IfMatch)
			if err != nil {
				return handleError(err)
			}

			shipmentSITStatus, _, err := h.CalculateShipmentSITStatus(appCtx, *updatedShipment)
			if err != nil {
				return handleError(err)
			}

			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())
			shipmentPayload := payloads.MTOShipment(h.FileStorer(), updatedShipment, sitStatusPayload)

			h.triggerDenySITExtensionEvent(appCtx, shipmentID, updatedShipment.MoveTaskOrderID, params)

			return shipmentops.NewDenySITExtensionOK().WithPayload(shipmentPayload), nil
		})
}

func (h DenySITExtensionHandler) triggerDenySITExtensionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.DenySITExtensionParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcDenySITExtensionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.DenySITExtensionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                     // ID of the updated logical object
		MtoID:           moveID,                         // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.DenySITExtensionHandler could not generate the event", zap.Error(err))
	}
}

// UpdateSITServiceItemCustomerExpenseHandler converts a SIT to customer expense
type UpdateSITServiceItemCustomerExpenseHandler struct {
	handlers.HandlerConfig
	services.MTOServiceItemUpdater
	services.MTOShipmentFetcher
	services.ShipmentSITStatus
}

// Handle ... converts the SIT to customer expense
func (h UpdateSITServiceItemCustomerExpenseHandler) Handle(params shipmentops.UpdateSITServiceItemCustomerExpenseParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error converting SIT to customer expense", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewUpdateSITServiceItemCustomerExpenseNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						handlers.ValidationErrMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewUpdateSITServiceItemCustomerExpenseUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewUpdateSITServiceItemCustomerExpensePreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return shipmentops.NewUpdateSITServiceItemCustomerExpenseForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewUpdateSITServiceItemCustomerExpenseInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("is not a TOO")
				return handleError(forbiddenError)
			}

			shipmentID := uuid.FromStringOrNil(string(params.ShipmentID))
			convertToCustomerExpense := params.Body.ConvertToCustomerExpense
			customerExpenseReason := params.Body.CustomerExpenseReason
			eagerAssociations := []string{"SITDurationUpdates",
				"MTOServiceItems",
				"MTOServiceItems.ReService.Code"}

			shipment, err := h.MTOShipmentFetcher.GetShipment(appCtx, shipmentID, eagerAssociations...)
			if err != nil {
				return handleError(err)
			}
			if *convertToCustomerExpense {
				_, err = h.MTOServiceItemUpdater.ConvertItemToCustomerExpense(appCtx, shipment, customerExpenseReason, true)
				if err != nil {
					return handleError(err)
				}
			}
			shipmentSITStatus, _, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			payload := payloads.MTOShipment(h.FileStorer(), shipment, sitStatusPayload)
			return shipmentops.NewUpdateSITServiceItemCustomerExpenseOK().WithPayload(payload), nil
		})
}

// CreateApprovedSITDurationUpdateHandler creates a SIT Duration Update in the approved state
type CreateApprovedSITDurationUpdateHandler struct {
	handlers.HandlerConfig
	services.ApprovedSITDurationUpdateCreator
	services.ShipmentSITStatus
	services.ShipmentUpdater
}

// Handle creates the approved SIT extension
func (h CreateApprovedSITDurationUpdateHandler) Handle(params shipmentops.CreateApprovedSITDurationUpdateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			shipmentID := params.ShipmentID

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.CreateApprovedSITExtension error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return shipmentops.NewCreateApprovedSITDurationUpdateNotFound().WithPayload(&payload), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"CreateApprovedSITExtension",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewCreateApprovedSITDurationUpdateUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewDenySITExtensionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.CreateApprovedSITExtension query error", zap.Error(e.Unwrap()))
					}
					return shipmentops.NewCreateApprovedSITDurationUpdateInternalServerError(), err
				case apperror.ForbiddenError:
					return shipmentops.NewCreateApprovedSITDurationUpdateForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewCreateApprovedSITDurationUpdateInternalServerError(), err
				}
			}

			sitExtension := payloads.ApprovedSITExtensionFromCreate(payload, shipmentID)
			shipment, err := h.ApprovedSITDurationUpdateCreator.CreateApprovedSITDurationUpdate(appCtx, sitExtension, sitExtension.MTOShipmentID, params.IfMatch)
			if err != nil {
				return handleError(err)
			}

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				return handleError(apperror.NewForbiddenError("is not a TOO"))
			}

			shipmentSITStatus, shipmentWithSITInfo, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}

			existingETag := etag.GenerateEtag(shipment.UpdatedAt)

			shipment, err = h.UpdateShipment(appCtx, &shipmentWithSITInfo, existingETag, "ghc")
			if err != nil {
				return handleError(err)
			}

			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())
			returnPayload := payloads.MTOShipment(h.FileStorer(), shipment, sitStatusPayload)
			return shipmentops.NewCreateApprovedSITDurationUpdateOK().WithPayload(returnPayload), nil
		})
}
