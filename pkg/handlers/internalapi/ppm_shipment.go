package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// SubmitPPMShipmentDocumentationHandler is the handler to save a PPMShipment signature and route the PPM shipment to the office
type SubmitPPMShipmentDocumentationHandler struct {
	handlers.HandlerConfig
	services.PPMShipmentNewSubmitter
}

// Handle saves a new customer signature for PPMShipment documentation submission and routes PPM shipment to the
// service counselor.
func (h SubmitPPMShipmentDocumentationHandler) Handle(params ppmops.SubmitPPMShipmentDocumentationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				return ppmops.NewSubmitPPMShipmentDocumentationUnauthorized(), apperror.NewSessionError("No user session")
			} else if !appCtx.Session().IsMilApp() {
				return ppmops.NewSubmitPPMShipmentDocumentationForbidden(), apperror.NewSessionError("Request is not from the customer app")
			} else if appCtx.Session().UserID.IsNil() {
				return ppmops.NewSubmitPPMShipmentDocumentationForbidden(), apperror.NewSessionError("No user ID in session")
			}

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil || ppmShipmentID.IsNil() {
				appCtx.Logger().Error("error with PPM Shipment ID", zap.Error(err))

				errDetail := "Invalid PPM shipment ID in URL"

				if err != nil {
					errDetail = errDetail + ": " + err.Error()
				}

				errPayload := payloads.ClientError(
					handlers.BadRequestErrMessage,
					errDetail,
					h.GetTraceIDFromRequest(params.HTTPRequest),
				)

				return ppmops.NewSubmitPPMShipmentDocumentationBadRequest().WithPayload(errPayload), err
			}

			payload := params.SavePPMShipmentSignedCertificationPayload
			if payload == nil {
				noBodyErr := apperror.NewBadDataError("No body provided")

				appCtx.Logger().Error("No body provided", zap.Error(noBodyErr))

				errPayload := payloads.ClientError(
					handlers.BadRequestErrMessage,
					noBodyErr.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
				)

				return ppmops.NewSubmitPPMShipmentDocumentationBadRequest().WithPayload(errPayload), noBodyErr
			}

			signedCertification := payloads.SavePPMShipmentSignedCertification(ppmShipmentID, *payload)

			ppmShipment, err := h.PPMShipmentNewSubmitter.SubmitNewCustomerCloseOut(appCtx, ppmShipmentID, signedCertification)

			if err != nil {
				appCtx.Logger().Error("internalapi.SubmitPPMShipmentDocumentationHandler", zap.Error(err))

				switch e := err.(type) {
				case *apperror.BadDataError:
					errPayload := payloads.ClientError(
						handlers.BadRequestErrMessage,
						e.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
					)

					return ppmops.NewSubmitPPMShipmentDocumentationBadRequest().WithPayload(errPayload), err
				case apperror.NotFoundError:
					errPayload := payloads.ClientError(
						handlers.NotFoundMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
					)

					return ppmops.NewSubmitPPMShipmentDocumentationNotFound().WithPayload(errPayload), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error(
							"internalapi.SubmitPPMShipmentDocumentationHandler error",
							zap.Error(e.Unwrap()),
						)
					}

					errPayload := payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))

					return ppmops.NewSubmitPPMShipmentDocumentationInternalServerError().WithPayload(errPayload), err
				case apperror.InvalidInputError:
					errPayload := payloads.ValidationError(
						handlers.ValidationErrMessage,
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors,
					)

					return ppmops.NewSubmitPPMShipmentDocumentationUnprocessableEntity().WithPayload(errPayload), err
				case apperror.ConflictError:
					errPayload := payloads.ClientError(
						handlers.ConflictErrMessage,
						e.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
					)

					return ppmops.NewSubmitPPMShipmentDocumentationConflict().WithPayload(errPayload), err
				default:
					errPayload := payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))

					return ppmops.NewSubmitPPMShipmentDocumentationInternalServerError().WithPayload(errPayload), err
				}
			}

			returnPayload := payloads.PPMShipment(h.FileStorer(), ppmShipment)

			return ppmops.NewSubmitPPMShipmentDocumentationOK().WithPayload(returnPayload), nil
		})
}
