package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	gunsafeops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateGunSafeWeightTicketHandler
type CreateGunSafeWeightTicketHandler struct {
	handlers.HandlerConfig
	gunsafeCreator services.GunSafeWeightTicketCreator
}

// Handle creating a gunsafe weight ticket
func (h CreateGunSafeWeightTicketHandler) Handle(params gunsafeops.CreateGunSafeWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

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
				return gunsafeops.NewCreateGunSafeWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Feature flag for gun safe related endpoints has not been enabled.")
			}

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return gunsafeops.NewCreateGunSafeWeightTicketUnauthorized(), noSessionErr
			}
			if !appCtx.Session().IsOfficeApp() {
				return gunsafeops.NewCreateGunSafeWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				appCtx.Logger().Error("missing PPM Shipment ID", zap.Error(err))
				return gunsafeops.NewCreateGunSafeWeightTicketBadRequest(), nil
			}

			gunsafe, err := h.gunsafeCreator.CreateGunSafeWeightTicket(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.CreateGunSafeWeightTicketHandler", zap.Error(err))
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
			returnPayload := payloads.GunSafeWeightTicket(h.FileStorer(), gunsafe)

			if returnPayload == nil {
				appCtx.Logger().Error("Returned Payload is empty", zap.Error(err))
				return gunsafeops.NewCreateGunSafeWeightTicketInternalServerError(), nil
			}
			return gunsafeops.NewCreateGunSafeWeightTicketCreated().WithPayload(returnPayload), nil
		})
}

// UpdateGunSafeWeightTicketHandler
type UpdateGunSafeWeightTicketHandler struct {
	handlers.HandlerConfig
	gunsafeUpdater services.GunSafeWeightTicketUpdater
}

func (h UpdateGunSafeWeightTicketHandler) Handle(params gunsafeops.UpdateGunSafeWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.UpdateGunSafeWeightTicket
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

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
				return gunsafeops.NewUpdateGunSafeWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Feature flag for gun safe related endpoints has not been enabled.")
			}

			GunSafeWeightTicket := payloads.GunSafeWeightTicketModelFromUpdate(payload)

			if !appCtx.Session().IsOfficeApp() {
				return gunsafeops.NewUpdateGunSafeWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			GunSafeWeightTicket.ID = uuid.FromStringOrNil(params.GunSafeWeightTicketID.String())

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.UpdateWeightTicketHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return gunsafeops.NewUpdateGunSafeWeightTicketNotFound(), err
				case apperror.InvalidInputError:
					return gunsafeops.NewUpdateGunSafeWeightTicketUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return gunsafeops.NewUpdateGunSafeWeightTicketPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the error (usually a pq error) for better debugging
						appCtx.Logger().Error(
							"ghcapi.GetWeightTicketsHandler error",
							zap.Error(e.Unwrap()),
						)
					}
					return gunsafeops.NewUpdateGunSafeWeightTicketInternalServerError(), err
				default:
					return gunsafeops.NewUpdateGunSafeWeightTicketInternalServerError(), err
				}
			}

			updatedGunSafeWeightTicket, err := h.gunsafeUpdater.UpdateGunSafeWeightTicket(appCtx, *GunSafeWeightTicket, params.IfMatch)

			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.GunSafeWeightTicket(h.FileStorer(), updatedGunSafeWeightTicket)
			return gunsafeops.NewUpdateGunSafeWeightTicketOK().WithPayload(returnPayload), nil
		})
}

// DeleteGunSafeWeightTicketHandler
type DeleteGunSafeWeightTicketHandler struct {
	handlers.HandlerConfig
	gunsafeDeleter services.GunSafeWeightTicketDeleter
}

// Handle deletes a gun safe weight ticket
func (h DeleteGunSafeWeightTicketHandler) Handle(params gunsafeops.DeleteGunSafeWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

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
				return gunsafeops.NewDeleteGunSafeWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Feature flag for gun safe related endpoints has not been enabled.")
			}

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return gunsafeops.NewDeleteGunSafeWeightTicketUnauthorized(), noSessionErr
			}
			if !appCtx.Session().IsOfficeApp() {
				return gunsafeops.NewDeleteGunSafeWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			GunSafeWeightTicketID := uuid.FromStringOrNil(params.GunSafeWeightTicketID.String())
			err := h.gunsafeDeleter.DeleteGunSafeWeightTicket(appCtx, ppmID, GunSafeWeightTicketID)
			if err != nil {
				appCtx.Logger().Error("ghcapi.DeleteGunSafeWeightTicketHandler", zap.Error(err))

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
