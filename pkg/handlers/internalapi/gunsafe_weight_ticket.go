package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	gunsafeops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateGunSafeWeightTicketHandler
type CreateGunSafeWeightTicketHandler struct {
	handlers.HandlerConfig
	gunSafeCreator services.GunSafeWeightTicketCreator
}

// Handle creating a gunSafe weight ticket
func (h CreateGunSafeWeightTicketHandler) Handle(params gunsafeops.CreateGunSafeWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			/** Feature Flag - GUN_SAFE **/
			const featureFlagNameGunSafe = "gun_safe"
			isGunSafeFeatureOn := false
			flag, ffErr := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagNameGunSafe, map[string]string{})

			if ffErr != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagNameGunSafe), zap.Error(ffErr))
			} else {
				isGunSafeFeatureOn = flag.Match
			}

			if !isGunSafeFeatureOn {
				return gunsafeops.NewCreateGunSafeWeightTicketForbidden(), apperror.NewSessionError("Feature flag for gun safe related endpoints has not been enabled.")
			}

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				appCtx.Logger().Error("missing PPM Shipment ID", zap.Error(err))
				return gunsafeops.NewCreateGunSafeWeightTicketBadRequest(), nil
			}

			gunSafe, err := h.gunSafeCreator.CreateGunSafeWeightTicket(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("internalapi.CreateGunSafeWeightTicketHandler", zap.Error(err))
				switch err.(type) {
				case apperror.InvalidInputError:
					return gunsafeops.NewCreateGunSafeWeightTicketUnprocessableEntity(), err
				case apperror.ForbiddenError:
					return gunsafeops.NewCreateGunSafeWeightTicketForbidden(), err
				case apperror.NotFoundError:
					return gunsafeops.NewCreateGunSafeWeightTicketNotFound(), err
				default:
					return gunsafeops.NewCreateGunSafeWeightTicketInternalServerError(), err
				}
			}
			returnPayload := payloads.GunSafeWeightTicket(h.FileStorer(), gunSafe)

			if returnPayload == nil {
				return gunsafeops.NewCreateGunSafeWeightTicketInternalServerError(), err
			}
			return gunsafeops.NewCreateGunSafeWeightTicketCreated().WithPayload(returnPayload), nil
		})
}

// UpdateGunSafeWeightTicketHandler
type UpdateGunSafeWeightTicketHandler struct {
	handlers.HandlerConfig
	gunSafeUpdater services.GunSafeWeightTicketUpdater
}

func (h UpdateGunSafeWeightTicketHandler) Handle(params gunsafeops.UpdateGunSafeWeightTicketParams) middleware.Responder {
	// track every request with middleware:
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			/** Feature Flag - GUN_SAFE **/
			const featureFlagNameGunSafe = "gun_safe"
			isGunSafeFeatureOn := false
			flag, ffErr := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagNameGunSafe, map[string]string{})

			if ffErr != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagNameGunSafe), zap.Error(ffErr))
			} else {
				isGunSafeFeatureOn = flag.Match
			}

			if !isGunSafeFeatureOn {
				return gunsafeops.NewUpdateGunSafeWeightTicketForbidden(), apperror.NewSessionError("Feature flag for gun safe related endpoints has not been enabled.")
			}

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return gunsafeops.NewUpdateGunSafeWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return gunsafeops.NewUpdateGunSafeWeightTicketForbidden(), noServiceMemberIDErr
			}

			payload := params.UpdateGunSafeWeightTicket
			if payload == nil {
				noBodyErr := apperror.NewBadDataError("Invalid weight ticket: params UpdateGunSafePayload is nil")
				appCtx.Logger().Error(noBodyErr.Error())
				return gunsafeops.NewUpdateGunSafeWeightTicketBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The Weight Ticket request payload cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), noBodyErr
			}

			weightTicket := payloads.GunSafeWeightTicketModelFromUpdate(payload)
			weightTicket.ID = uuid.FromStringOrNil(params.GunSafeWeightTicketID.String())

			updateGunSafe, err := h.gunSafeUpdater.UpdateGunSafeWeightTicket(appCtx, *weightTicket, params.IfMatch)

			if err != nil {
				appCtx.Logger().Error("internalapi.UpdateGunSafeWeightTicketHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.InvalidInputError:
					return gunsafeops.NewUpdateGunSafeWeightTicketUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.PreconditionFailedError:
					return gunsafeops.NewUpdateGunSafeWeightTicketPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.NotFoundError:
					return gunsafeops.NewUpdateGunSafeWeightTicketNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.
							Logger().
							Error(
								"internalapi.UpdateGunSafeWeightTicketHandler error",
								zap.Error(e.Unwrap()),
							)
					}
					return gunsafeops.
						NewUpdateGunSafeWeightTicketInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				default:
					return gunsafeops.
						NewUpdateGunSafeWeightTicketInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				}

			}
			returnPayload := payloads.GunSafeWeightTicket(h.FileStorer(), updateGunSafe)
			return gunsafeops.NewUpdateGunSafeWeightTicketOK().WithPayload(returnPayload), nil
		})
}

// DeleteGunSafeWeightTicketHandler
type DeleteGunSafeWeightTicketHandler struct {
	handlers.HandlerConfig
	gunSafeDeleter services.GunSafeWeightTicketDeleter
}

// Handle deletes a gun safe weight ticket
func (h DeleteGunSafeWeightTicketHandler) Handle(params gunsafeops.DeleteGunSafeWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			/** Feature Flag - GUN_SAFE **/
			const featureFlagNameGunSafe = "gun_safe"
			isGunSafeFeatureOn := false
			flag, ffErr := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagNameGunSafe, map[string]string{})

			if ffErr != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagNameGunSafe), zap.Error(ffErr))
			} else {
				isGunSafeFeatureOn = flag.Match
			}

			if !isGunSafeFeatureOn {
				return gunsafeops.NewDeleteGunSafeWeightTicketForbidden(), apperror.NewSessionError("Feature flag for gun safe related endpoints has not been enabled.")
			}

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				appCtx.Logger().Error("internalapi.DeleteGunSafeWeightTicketHandler", zap.Error(noSessionErr))
				return gunsafeops.NewDeleteGunSafeWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() || appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				appCtx.Logger().Error("internalapi.DeleteGunSafeWeightTicketHandler", zap.Error(noServiceMemberIDErr))
				return gunsafeops.NewDeleteGunSafeWeightTicketForbidden(), noServiceMemberIDErr
			}

			// Make sure the service member is not modifying another service member's PPM
			ppmID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			gunSafeWeightTicketID := uuid.FromStringOrNil(params.GunSafeWeightTicketID.String())
			err := h.gunSafeDeleter.DeleteGunSafeWeightTicket(appCtx, ppmID, gunSafeWeightTicketID)
			if err != nil {
				appCtx.Logger().Error("internalapi.DeleteGunSafeWeightTicketHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return gunsafeops.NewDeleteGunSafeWeightTicketNotFound(), err
				case apperror.ConflictError:
					return gunsafeops.NewDeleteGunSafeWeightTicketConflict(), err
				case apperror.ForbiddenError:
					return gunsafeops.NewDeleteGunSafeWeightTicketForbidden(), err
				case apperror.UnprocessableEntityError:
					return gunsafeops.NewDeleteGunSafeWeightTicketUnprocessableEntity(), err
				default:
					return gunsafeops.NewDeleteGunSafeWeightTicketInternalServerError(), err
				}
			}

			return gunsafeops.NewDeleteGunSafeWeightTicketNoContent(), nil
		})
}
