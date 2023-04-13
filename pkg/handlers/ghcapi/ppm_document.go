package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmdocumentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetPPMDocumentsHandler is the handler that fetches all of the documents for a PPM shipment for the office api
type GetPPMDocumentsHandler struct {
	handlers.HandlerConfig
	services.PPMDocumentFetcher
}

// Handle retrieves all documents for a PPM shipment
func (h GetPPMDocumentsHandler) Handle(params ppmdocumentops.GetPPMDocumentsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))

			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppmdocumentops.NewGetPPMDocumentsForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())

			ppmDocuments, err := h.PPMDocumentFetcher.GetPPMDocuments(appCtx, shipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.GetPPMDocumentsHandler error", zap.Error(err))

				switch e := err.(type) {
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the error (usually a pq error) for better debugging
						appCtx.Logger().Error(
							"ghcapi.GetPPMDocumentsHandler error",
							zap.Error(e.Unwrap()),
						)
					}

					return ppmdocumentops.NewGetPPMDocumentsInternalServerError().WithPayload(errPayload), nil
				default:
					return ppmdocumentops.NewGetPPMDocumentsInternalServerError().WithPayload(errPayload), nil
				}
			}

			returnPayload := payloads.PPMDocuments(h.FileStorer(), ppmDocuments)

			return ppmdocumentops.NewGetPPMDocumentsOK().WithPayload(returnPayload), nil
		})
}

// FinishDocumentReviewHandler is the handler that updates a PPM shipment for the office api when documents have been reviewed

type FinishDocumentReviewHandler struct {
	handlers.HandlerConfig
	services.PPMShipmentReviewDocuments
}

// Handle finishes a review for a PPM shipment's documents
func (h FinishDocumentReviewHandler) Handle(params ppmdocumentops.FinishDocumentReviewParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			if appCtx.Session() == nil {
				return ppmdocumentops.NewFinishDocumentReviewUnauthorized(), apperror.NewSessionError("No user session")
			} else if !appCtx.Session().IsOfficeApp() {
				return ppmdocumentops.NewFinishDocumentReviewForbidden(), apperror.NewSessionError("Request is not from the customer app")
			} else if appCtx.Session().UserID.IsNil() {
				return ppmdocumentops.NewFinishDocumentReviewForbidden(), apperror.NewSessionError("No user ID in session")
			}

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil || ppmShipmentID.IsNil() {
				appCtx.Logger().Error("error with PPM Shipment ID", zap.Error(err))

				return ppmdocumentops.NewFinishDocumentReviewBadRequest().WithPayload(errPayload), err
			}
			ppmShipment, err := h.PPMShipmentReviewDocuments.SubmitReviewedDocuments(appCtx, ppmShipmentID)
			if err != nil {
				appCtx.Logger().Error("ghcapi.FinishDocumentReviewHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return ppmdocumentops.NewFinishDocumentReviewNotFound(), err
				case apperror.ConflictError:
					return ppmdocumentops.NewFinishDocumentReviewConflict(), err
				case apperror.InvalidInputError:
					return ppmdocumentops.NewFinishDocumentReviewUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return ppmdocumentops.NewFinishDocumentReviewPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the error (usually a pq error) for better debugging
						appCtx.Logger().Error(
							"ghcapi.FinishDocumentReviewHandler error",
							zap.Error(e.Unwrap()),
						)
					}
					return ppmdocumentops.NewFinishDocumentReviewInternalServerError(), err
				default:
					return ppmdocumentops.NewFinishDocumentReviewInternalServerError(), err
				}

			}

			returnPayload := payloads.PPMShipment(h.FileStorer(), ppmShipment)

			return ppmdocumentops.NewFinishDocumentReviewOK().WithPayload(returnPayload), nil
		})
}
