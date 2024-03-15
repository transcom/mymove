package internalapi

import (
	"fmt"
	"io"
	"time"

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

// ShowAOAPacketHandler returns a Shipment Summary Worksheet PDF
type showAOAPacketHandler struct {
	handlers.HandlerConfig
	services.SSWPPMComputer
	services.SSWPPMGenerator
	services.AOAPacketCreator
}

// Handle returns a generated PDF
func (h showAOAPacketHandler) Handle(params ppmops.ShowAOAPacketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			logger := appCtx.Logger()

			// Ensures session
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return ppmops.NewShowAOAPacketForbidden(), noSessionErr
			}
			// Ensures service member ID is present
			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return ppmops.NewShowAOAPacketForbidden(), noServiceMemberIDErr
			}

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID)
			if err != nil {
				err := apperror.NewBadDataError("missing/empty required URI parameter: PPMShipmentID")
				appCtx.Logger().Error(err.Error())
				return ppmops.NewShowAOAPacketBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			// Ensures AOA is for the accessing member
			err = h.VerifyAOAPacketInternal(appCtx, ppmShipmentID)
			if err != nil {
				err := apperror.NewBadDataError("PPMShipment cannot be verified")
				appCtx.Logger().Error(err.Error())
				return ppmops.NewShowAOAPacketBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			AOAPacket, err := h.AOAPacketCreator.CreateAOAPacket(appCtx, ppmShipmentID)
			if err != nil {
				logger.Error("Error creating AOA", zap.Error(err))
				aoaError := err.Error()
				payload := payloads.InternalServerError(&aoaError, h.GetTraceIDFromRequest(params.HTTPRequest))
				return ppmops.NewShowAOAPacketInternalServerError().
					WithPayload(payload), err
			}

			payload := io.NopCloser(AOAPacket)
			filename := fmt.Sprintf("inline; filename=\"AOA-%s.pdf\"", time.Now().Format("01-02-2006_15-04-05"))

			return ppmops.NewShowAOAPacketOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}

// ShowPaymentPacketHandler returns a PPM Payment Packet PDF
type ShowPaymentPacketHandler struct {
	handlers.HandlerConfig
	services.PaymentPacketCreator
}

// Handle returns a generated PDF
func (h ShowPaymentPacketHandler) Handle(params ppmops.ShowPaymentPacketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			pdf, err := h.PaymentPacketCreator.GenerateDefault(appCtx, ppmShipmentID)
			if err != nil {
				switch err.(type) {
				case apperror.ForbiddenError:
					// this indicates user does not have access to PPM
					appCtx.Logger().Warn(fmt.Sprintf("internalapi.DownPaymentPacket ForbiddenError ppmShipmentID:%s", ppmShipmentID.String()), zap.Error(err))
					return ppmops.NewShowPaymentPacketForbidden(), err
				case apperror.NotFoundError:
					// this indicates ppm was not found
					appCtx.Logger().Warn(fmt.Sprintf("internalapi.DownPaymentPacket NotFoundError ppmShipmentID:%s", ppmShipmentID.String()), zap.Error(err))
					return ppmops.NewShowPaymentPacketNotFound(), err
				default:
					appCtx.Logger().Error(fmt.Sprintf("internalapi.DownPaymentPacket InternalServerError ppmShipmentID:%s", ppmShipmentID.String()), zap.Error(err))
					return ppmops.NewShowPaymentPacketInternalServerError(), err
				}
			}

			payload := io.NopCloser(pdf)
			filename := fmt.Sprintf("inline; filename=\"ppm_payment_packet-%s.pdf\"", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))

			return ppmops.NewShowPaymentPacketOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}
