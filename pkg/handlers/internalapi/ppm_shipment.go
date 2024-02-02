package internalapi

import (
	"fmt"
	"io"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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

// ResubmitPPMShipmentDocumentationHandler is the handler to resubmit PPM shipment documentation
type ResubmitPPMShipmentDocumentationHandler struct {
	handlers.HandlerConfig
	services.PPMShipmentUpdatedSubmitter
}

// Handle updates a customer's PPM shipment payment signature and re-routes the shipment to the service counselor.
func (h ResubmitPPMShipmentDocumentationHandler) Handle(params ppmops.ResubmitPPMShipmentDocumentationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsMilApp() {
				return ppmops.NewResubmitPPMShipmentDocumentationForbidden(), apperror.NewSessionError("Request is not from the customer app")
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

				return ppmops.NewResubmitPPMShipmentDocumentationBadRequest().WithPayload(errPayload), err
			}

			signedCertificationID, err := uuid.FromString(params.SignedCertificationID.String())
			if err != nil || signedCertificationID.IsNil() {
				appCtx.Logger().Error("error with signed certification ID", zap.Error(err))

				errDetail := "Invalid signed certification ID in URL"

				if err != nil {
					errDetail = errDetail + ": " + err.Error()
				}

				errPayload := payloads.ClientError(
					handlers.BadRequestErrMessage,
					errDetail,
					h.GetTraceIDFromRequest(params.HTTPRequest),
				)

				return ppmops.NewResubmitPPMShipmentDocumentationBadRequest().WithPayload(errPayload), err
			}

			signedCertification := payloads.ReSavePPMShipmentSignedCertification(ppmShipmentID, signedCertificationID, *params.SavePPMShipmentSignedCertificationPayload)

			ppmShipment, err := h.PPMShipmentUpdatedSubmitter.SubmitUpdatedCustomerCloseOut(appCtx, ppmShipmentID, signedCertification, params.IfMatch)

			if err != nil {
				appCtx.Logger().Error("internalapi.ResubmitPPMShipmentDocumentationHandler", zap.Error(err))

				switch e := err.(type) {
				case *apperror.BadDataError:
					errPayload := payloads.ClientError(
						handlers.BadRequestErrMessage,
						e.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
					)

					return ppmops.NewResubmitPPMShipmentDocumentationBadRequest().WithPayload(errPayload), err
				case apperror.NotFoundError:
					errPayload := payloads.ClientError(
						handlers.NotFoundMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
					)

					return ppmops.NewResubmitPPMShipmentDocumentationNotFound().WithPayload(errPayload), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error(
							"internalapi.ResubmitPPMShipmentDocumentationHandler error",
							zap.Error(e.Unwrap()),
						)
					}

					errPayload := payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))

					return ppmops.NewResubmitPPMShipmentDocumentationInternalServerError().WithPayload(errPayload), err
				case apperror.InvalidInputError:
					errPayload := payloads.ValidationError(
						handlers.ValidationErrMessage,
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors,
					)

					return ppmops.NewResubmitPPMShipmentDocumentationUnprocessableEntity().WithPayload(errPayload), err
				case apperror.ConflictError:
					errPayload := payloads.ClientError(
						handlers.ConflictErrMessage,
						e.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
					)

					return ppmops.NewResubmitPPMShipmentDocumentationConflict().WithPayload(errPayload), err
				case apperror.PreconditionFailedError:
					errPayload := payloads.ClientError(
						handlers.PreconditionErrMessage,
						e.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
					)

					return ppmops.NewResubmitPPMShipmentDocumentationPreconditionFailed().WithPayload(errPayload), err
				default:
					errPayload := payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))

					return ppmops.NewResubmitPPMShipmentDocumentationInternalServerError().WithPayload(errPayload), err
				}
			}

			returnPayload := payloads.PPMShipment(h.FileStorer(), ppmShipment)

			return ppmops.NewResubmitPPMShipmentDocumentationOK().WithPayload(returnPayload), nil
		})
}

// Handle returns a generated PDF
func (h ShowShipmentSummaryWorksheetHandler) Handle(params ppmops.ShowShipmentSummaryWorksheetParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			logger := appCtx.Logger()

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				logger.Error("Error fetching PPMShipment", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			ssfd, err := h.SSWPPMComputer.FetchDataShipmentSummaryWorksheetFormData(appCtx, appCtx.Session(), ppmShipmentID)
			if err != nil {
				logger.Error("Error fetching data for SSW", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			ssfd.Obligations, err = h.SSWPPMComputer.ComputeObligations(appCtx, *ssfd, h.DTODPlanner())
			if err != nil {
				logger.Error("Error calculating obligations ", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			page1Data, page2Data := h.SSWPPMComputer.FormatValuesShipmentSummaryWorksheet(*ssfd)
			if err != nil {
				logger.Error("Error formatting data for SSW", zap.Error(err))
				return handlers.ResponseForError(logger, err), err
			}

			SSWPPMWorksheet, SSWPDFInfo, err := h.SSWPPMGenerator.FillSSWPDFForm(page1Data, page2Data)
			if err != nil {
				return nil, err
			}
			if SSWPDFInfo.PageCount != 2 {
				return nil, errors.Wrap(err, "SSWGenerator output a corrupted or incorretly altered PDF")
			}
			payload := io.NopCloser(SSWPPMWorksheet)
			filename := fmt.Sprintf("inline; filename=\"%s-%s-ssw-%s.pdf\"", *ssfd.ServiceMember.FirstName, *ssfd.ServiceMember.LastName, time.Now().Format("01-02-2006"))

			return ppmops.NewShowShipmentSummaryWorksheetOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}

// ShowShipmentSummaryWorksheetHandler returns a Shipment Summary Worksheet PDF
type ShowShipmentSummaryWorksheetHandler struct {
	handlers.HandlerConfig
	services.SSWPPMComputer
	services.SSWPPMGenerator
}
