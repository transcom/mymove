package ghcapi

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	uploadop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/uploads"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	weightticketparser "github.com/transcom/mymove/pkg/services/weight_ticket_parser"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

const weightEstimatePages = 11

type CreateUploadHandler struct {
	handlers.HandlerConfig
}

type CreatePPMUploadHandler struct {
	handlers.HandlerConfig
	services.WeightTicketGenerator
	services.WeightTicketComputer
	*uploaderpkg.UserUploader
}

func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			rollbackErr := apperror.NewBadDataError("error creating upload")

			file, ok := params.File.(*runtime.File)
			if !ok {
				appCtx.Logger().Error("This should always be a runtime.File, something has changed in go-swagger.")
				return uploadop.NewCreateUploadInternalServerError(), rollbackErr
			}

			appCtx.Logger().Info(
				"File uploader and size",
				zap.String("userID", appCtx.Session().UserID.String()),
				zap.String("serviceMemberID", appCtx.Session().ServiceMemberID.String()),
				zap.String("officeUserID", appCtx.Session().OfficeUserID.String()),
				zap.String("AdminUserID", appCtx.Session().AdminUserID.String()),
				zap.Int64("size", file.Header.Size),
			)

			var docID *uuid.UUID
			if params.DocumentID != nil {
				documentID, err := uuid.FromString(params.DocumentID.String())
				if err != nil {
					appCtx.Logger().Info("Badly formed UUID for document", zap.String("document_id", params.DocumentID.String()), zap.Error(err))
					return uploadop.NewCreateUploadBadRequest(), rollbackErr
				}

				// Fetch document to ensure user has access to it
				document, docErr := models.FetchDocument(appCtx.DB(), appCtx.Session(), documentID, true)
				if docErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), docErr), rollbackErr
				}
				docID = &document.ID
			}

			newUserUpload, url, verrs, createErr := uploaderpkg.CreateUserUploadForDocumentWrapper(
				appCtx,
				appCtx.Session().UserID,
				h.FileStorer(),
				file,
				file.Header.Filename,
				uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
				uploaderpkg.AllowedTypesServiceMember,
				docID,
			)

			if verrs.HasAny() || createErr != nil {
				appCtx.Logger().Error("failed to create new user upload", zap.Error(createErr), zap.String("verrs", verrs.Error()))
				switch createErr.(type) {
				case uploaderpkg.ErrTooLarge:
					return uploadop.NewCreateUploadRequestEntityTooLarge(), rollbackErr
				case uploaderpkg.ErrFile:
					return uploadop.NewCreateUploadInternalServerError(), rollbackErr
				case uploaderpkg.ErrFailedToInitUploader:
					return uploadop.NewCreateUploadInternalServerError(), rollbackErr
				default:
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, createErr), rollbackErr
				}
			}

			uploadPayload := payloads.PayloadForUploadModel(h.FileStorer(), newUserUpload.Upload, url)
			return uploadop.NewCreateUploadCreated().WithPayload(uploadPayload), nil
		})
}

func (h CreatePPMUploadHandler) Handle(params ppmop.CreatePPMUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			rollbackErr := fmt.Errorf("error creating upload")

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().IsOfficeApp() {
				return ppmop.NewCreatePPMUploadForbidden(), apperror.NewForbiddenError("is not an Office User.")
			}

			file, ok := params.File.(*runtime.File)
			if !ok {
				appCtx.Logger().Error("This should always be a runtime.File, something has changed in go-swagger.")
				return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
			}

			appCtx.Logger().Info(
				"File uploader and size",
				zap.String("userID", appCtx.Session().UserID.String()),
				zap.String("serviceMemberID", appCtx.Session().ServiceMemberID.String()),
				zap.Int64("size", file.Header.Size),
			)

			documentID := uuid.FromStringOrNil(params.DocumentID.String())

			// Fetch document to ensure user has access to it
			document, docErr := models.FetchDocument(appCtx.DB(), appCtx.Session(), documentID, true)
			if docErr != nil {
				docNotFoundErr := fmt.Errorf("documentId %q was not found for this user", documentID)
				// errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
				errString := docNotFoundErr.Error()
				errPayload := &ghcmessages.Error{Message: &errString}
				return ppmop.NewCreatePPMUploadNotFound().WithPayload(errPayload), docNotFoundErr
			}

			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			// Ensure the document belongs to an association of the PPM shipment
			shipErr := ppmshipment.FindPPMShipmentWithDocument(appCtx, ppmShipmentID, documentID)
			if shipErr != nil {
				docNotFoundErr := fmt.Errorf("documentId %q was not found for this shipment", documentID)
				errString := docNotFoundErr.Error()
				errPayload := &ghcmessages.Error{Message: &errString}
				return ppmop.NewCreatePPMUploadNotFound().WithPayload(errPayload), docNotFoundErr
			}

			var newUserUpload *models.UserUpload
			var verrs *validate.Errors
			var url string
			var createErr error
			isWeightEstimatorFile := false

			uploadedFile := file

			// check if this is an excel file and parse if it is
			extension := filepath.Ext(file.Header.Filename)

			if extension == ".xlsx" {
				var err error

				isWeightEstimatorFile, err = weightticketparser.IsWeightEstimatorFile(appCtx, file)

				if err != nil {
					return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
				}

				_, err = file.Data.Seek(0, io.SeekStart)

				if err != nil {
					return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
				}
			}

			if params.WeightReceipt && isWeightEstimatorFile {
				pageValues, err := h.WeightTicketComputer.ParseWeightEstimatorExcelFile(appCtx, file)

				if err != nil {
					return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
				}

				pdfFileName := strings.TrimSuffix(file.Header.Filename, filepath.Ext(file.Header.Filename)) + ".pdf"
				aFile, pdfInfo, err := h.WeightTicketGenerator.FillWeightEstimatorPDFForm(*pageValues, pdfFileName)

				// Ensure weight receipt PDF is not corrupted
				if err != nil || pdfInfo.PageCount != weightEstimatePages {
					return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
				}

				// we already generated an afero file so we can skip that process the wrapper method does
				newUserUpload, verrs, createErr = h.UserUploader.CreateUserUploadForDocument(appCtx, &document.ID, appCtx.Session().UserID, uploaderpkg.File{File: aFile}, uploaderpkg.AllowedTypesPPMDocuments)
				if verrs.HasAny() || createErr != nil {
					appCtx.Logger().Error("failed to create new user upload", zap.Error(createErr), zap.String("verrs", verrs.Error()))
					switch createErr.(type) {
					case uploaderpkg.ErrUnsupportedContentType:
						return ppmop.NewCreatePPMUploadUnprocessableEntity().WithPayload(payloads.ValidationError(createErr.Error(), uuid.Nil, verrs)), createErr
					case uploaderpkg.ErrTooLarge:
						return ppmop.NewCreatePPMUploadRequestEntityTooLarge(), createErr
					case uploaderpkg.ErrFile:
						return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
					case uploaderpkg.ErrFailedToInitUploader:
						return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
					default:
						return handlers.ResponseForVErrors(appCtx.Logger(), verrs, createErr), createErr
					}
				}

				url, err = h.UserUploader.PresignedURL(appCtx, newUserUpload)

				if err != nil {
					return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
				}
			} else {
				newUserUpload, url, verrs, createErr = uploaderpkg.CreateUserUploadForDocumentWrapper(
					appCtx,
					appCtx.Session().UserID,
					h.FileStorer(),
					uploadedFile,
					uploadedFile.Header.Filename,
					uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
					uploaderpkg.AllowedTypesPPMDocuments,
					&document.ID,
				)

				if verrs.HasAny() || createErr != nil {
					appCtx.Logger().Error("failed to create new user upload", zap.Error(createErr), zap.String("verrs", verrs.Error()))
					switch createErr.(type) {
					case uploaderpkg.ErrUnsupportedContentType:
						return ppmop.NewCreatePPMUploadUnprocessableEntity().WithPayload(payloads.ValidationError(createErr.Error(), uuid.Nil, verrs)), createErr
					case uploaderpkg.ErrTooLarge:
						return ppmop.NewCreatePPMUploadRequestEntityTooLarge(), createErr
					case uploaderpkg.ErrFile:
						return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
					case uploaderpkg.ErrFailedToInitUploader:
						return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
					default:
						return handlers.ResponseForVErrors(appCtx.Logger(), verrs, createErr), createErr
					}
				}
			}

			uploadPayload := payloads.PayloadForUploadModel(h.FileStorer(), newUserUpload.Upload, url)
			return ppmop.NewCreatePPMUploadCreated().WithPayload(uploadPayload), nil
		})
}
