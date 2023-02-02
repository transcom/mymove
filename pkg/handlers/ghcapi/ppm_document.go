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
				case apperror.ForbiddenError:
					return ppmdocumentops.NewGetPPMDocumentsForbidden().WithPayload(errPayload), nil
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
