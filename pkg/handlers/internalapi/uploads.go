package internalapi

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	weightticketparser "github.com/transcom/mymove/pkg/services/weight_ticket_parser"
	"github.com/transcom/mymove/pkg/uploader"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

const weightEstimatePages = 11

// CreateUploadHandler creates a new upload via POST /uploads?documentId={documentId}
type CreateUploadHandler struct {
	handlers.HandlerConfig
}

type CreatePPMUploadHandler struct {
	handlers.HandlerConfig
	services.WeightTicketGenerator
	services.WeightTicketComputer
	*uploader.UserUploader
}

// Handle creates a new UserUpload from a request payload
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			rollbackErr := fmt.Errorf("error creating upload")

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
				models.UploadTypeUSER,
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

// DeleteUploadHandler deletes an upload
type DeleteUploadHandler struct {
	handlers.HandlerConfig
	services.UploadInformationFetcher
}

// Handle deletes an upload
func (h DeleteUploadHandler) Handle(params uploadop.DeleteUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			uploadID, _ := uuid.FromString(params.UploadID.String())
			userUpload, err := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), uploadID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			var ppmShipmentStatus models.PPMShipmentStatus

			if params.PpmID != nil {
				ppmShipmentId, _ := uuid.FromString(params.PpmID.String())
				ppmShipment, err := models.FetchPPMShipmentByPPMShipmentID(appCtx.DB(), ppmShipmentId)
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
				ppmShipmentStatus = ppmShipment.Status
			}

			if params.OrderID != nil {
				orderID, _ := uuid.FromString(params.OrderID.String())
				move, e := models.FetchMoveByOrderID(appCtx.DB(), orderID)
				if e != nil {
					return handlers.ResponseForError(appCtx.Logger(), e), e
				}
				uploadInformation, e := h.FetchUploadInformationForDeletion(appCtx, uploadID, move.Locator)
				if e != nil {
					appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(e))
				}

				//If move status is not DRAFT, upload cannot be deleted
				if *uploadInformation.MoveStatus != models.MoveStatusDRAFT {
					return uploadop.NewDeleteUploadForbidden(), fmt.Errorf("deletion not permitted Move is not in 'DRAFT' status")
				}

				userUploader, e := uploaderpkg.NewUserUploader(
					h.FileStorer(),
					uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
				)
				if e != nil {
					appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(e))
				}
				if e = userUploader.DeleteUserUpload(appCtx, &userUpload); e != nil {
					return handlers.ResponseForError(appCtx.Logger(), e), e
				}

				return uploadop.NewDeleteUploadNoContent(), nil
			}

			if params.MoveID != nil {
				moveID, e := uuid.FromString(params.MoveID.String())
				if e != nil {
					appCtx.Logger().Error(fmt.Sprintf("UUID Parsing for %s", moveID.String()), zap.Error(err))
					return handlers.ResponseForError(appCtx.Logger(), e), e
				}

				userUploader, e := uploaderpkg.NewUserUploader(
					h.FileStorer(),
					uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
				)
				if e != nil {
					appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(e))
				}
				if e = userUploader.DeleteUserUpload(appCtx, &userUpload); e != nil {
					return handlers.ResponseForError(appCtx.Logger(), e), e
				}

				return uploadop.NewDeleteUploadNoContent(), nil
			}

			//Fetch upload information so we can retrieve the move status
			uploadInformation, err := h.FetchUploadInformation(appCtx, uploadID)
			if err != nil {
				appCtx.Logger().Error("error retrieving move associated with this upload", zap.Error(err))
			}

			//If move status is not DRAFT and customer is not uploading ppm docs, upload cannot be deleted
			if (*uploadInformation.MoveStatus != models.MoveStatusDRAFT) && (ppmShipmentStatus != models.PPMShipmentStatusWaitingOnCustomer) {
				return uploadop.NewDeleteUploadForbidden(), fmt.Errorf("deletion not permitted Move is not in 'DRAFT' status")
			}

			userUploader, err := uploaderpkg.NewUserUploader(
				h.FileStorer(),
				uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
			)
			if err != nil {
				appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			}
			if err = userUploader.DeleteUserUpload(appCtx, &userUpload); err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			return uploadop.NewDeleteUploadNoContent(), nil

		})
}

// DeleteUploadsHandler deletes a collection of uploads
type DeleteUploadsHandler struct {
	handlers.HandlerConfig
}

// Handle deletes uploads
func (h DeleteUploadsHandler) Handle(params uploadop.DeleteUploadsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			userUploader, err := uploaderpkg.NewUserUploader(
				h.FileStorer(),
				uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
			)
			if err != nil {
				appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			}

			for _, uploadID := range params.UploadIds {
				uploadUUID, _ := uuid.FromString(uploadID.String())
				userUpload, err := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), uploadUUID)
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}

				if err = userUploader.DeleteUserUpload(appCtx, &userUpload); err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
			}

			return uploadop.NewDeleteUploadsNoContent(), nil
		})
}

// UploadStatusHandler returns status of an upload
type GetUploadStatusHandler struct {
	handlers.HandlerConfig
	services.UploadInformationFetcher
}

type CustomNewUploadStatusOK struct {
	params uploadop.GetUploadStatusParams
	appCtx appcontext.AppContext
}

func (o *CustomNewUploadStatusOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	id_counter := 0

	// TODO: add check for permissions to view upload

	err := o.appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		uploadId, err := uuid.FromString(o.params.UploadID.String())
		if err != nil {
			panic(err)
		}
		uploaded, err := models.FetchUserUploadFromUploadID(txnAppCtx.DB(), txnAppCtx.Session(), uploadId)
		if err != nil {
			txnAppCtx.Logger().Error(err.Error())
		}

		txnAppCtx.Logger().Info("HELLOW: " + uploaded.UploadID.String())

		return err
	})

	if err != nil {
		o.appCtx.Logger().Error(err.Error())
	}

	for range 2 {
		resProcess := []byte("id: " + strconv.Itoa(id_counter) + "\nevent: message\ndata: PROCESSING\n\n")
		if err := producer.Produce(rw, resProcess); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
		if f, ok := rw.(http.Flusher); ok {
			f.Flush()
		}

		time.Sleep(4 * time.Second)
		id_counter++
	}

	resClean := []byte("id: " + strconv.Itoa(id_counter) + "\nevent: message\ndata: CLEAN\n\n")
	if err := producer.Produce(rw, resClean); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// Handle returns status of an upload
func (h GetUploadStatusHandler) Handle(params uploadop.GetUploadStatusParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			return &CustomNewUploadStatusOK{
				params: params,
				appCtx: h.AppContextFromRequest(params.HTTPRequest),
			}, nil
		})
}

func (h CreatePPMUploadHandler) Handle(params ppmop.CreatePPMUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			rollbackErr := fmt.Errorf("error creating upload")

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
				return ppmop.NewCreatePPMUploadNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, docNotFoundErr.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), docNotFoundErr
			}

			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			// Ensure the document belongs to an association of the PPM shipment
			shipErr := ppmshipment.FindPPMShipmentWithDocument(appCtx, ppmShipmentID, documentID)
			if shipErr != nil {
				docNotFoundErr := fmt.Errorf("documentId %q was not found for this shipment", documentID)
				return ppmop.NewCreatePPMUploadNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, docNotFoundErr.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), docNotFoundErr
			}

			var newUserUpload *models.UserUpload
			var verrs *validate.Errors
			var url string
			var createErr error
			isWeightEstimatorFile := false

			uploadedFile := file

			// extract extension from filename
			filename := file.Header.Filename
			timestampPattern := regexp.MustCompile(`-(\d{14})$`)

			timestamp := ""
			filenameWithoutTimestamp := ""
			if matches := timestampPattern.FindStringSubmatch(filename); len(matches) > 1 {
				timestamp = matches[1]
				filenameWithoutTimestamp = strings.TrimSuffix(filename, "-"+timestamp)
			} else {
				filenameWithoutTimestamp = filename
			}

			extension := filepath.Ext(filenameWithoutTimestamp)
			extensionLower := strings.ToLower(extension)

			// check if file is an excel file
			if extensionLower == ".xlsx" {
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

				pdfFileName := strings.TrimSuffix(filenameWithoutTimestamp, filepath.Ext(filenameWithoutTimestamp)) + ".pdf" + "-" + timestamp
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
					models.UploadTypeUSER,
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
