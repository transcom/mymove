package ghcapi

import (
	"context"
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
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	uploadop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/uploads"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/upload"
	weightticketparser "github.com/transcom/mymove/pkg/services/weight_ticket_parser"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

const weightEstimatePages = 11

type CreateUploadHandler struct {
	handlers.HandlerConfig
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
				document, docErr := models.FetchDocument(appCtx.DB(), appCtx.Session(), documentID)
				if docErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), docErr), rollbackErr
				}
				docID = &document.ID
			}

			newUserUpload, url, verrs, createErr := uploader.CreateUserUploadForDocumentWrapper(
				appCtx,
				appCtx.Session().UserID,
				h.FileStorer(),
				file,
				file.Header.Filename,
				uploader.MaxCustomerUserUploadFileSizeLimit,
				uploader.AllowedTypesServiceMember,
				docID,
				models.UploadTypeOFFICE,
			)

			if verrs.HasAny() || createErr != nil {
				appCtx.Logger().Error("failed to create new user upload", zap.Error(createErr), zap.String("verrs", verrs.Error()))
				switch createErr.(type) {
				case uploader.ErrTooLarge:
					return uploadop.NewCreateUploadRequestEntityTooLarge(), rollbackErr
				case uploader.ErrFile:
					return uploadop.NewCreateUploadInternalServerError(), rollbackErr
				case uploader.ErrFailedToInitUploader:
					return uploadop.NewCreateUploadInternalServerError(), rollbackErr
				default:
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, createErr), rollbackErr
				}
			}

			uploadPayload := payloads.PayloadForUploadModel(h.FileStorer(), newUserUpload.Upload, url)
			return uploadop.NewCreateUploadCreated().WithPayload(uploadPayload), nil
		})
}

type UpdateUploadHandler struct {
	handlers.HandlerConfig
	services.UploadInformationFetcher
}

func (h UpdateUploadHandler) Handle(params uploadop.UpdateUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().IsOfficeApp() {
				forbiddenError := apperror.NewForbiddenError("User is not an Office User.")
				appCtx.Logger().Error(forbiddenError.Error())
				return uploadop.NewUpdateUploadForbidden(), forbiddenError
			}

			uploadID, _ := uuid.FromString(params.UploadID.String())
			updater := upload.NewUploadUpdater()
			newUpload, err := updater.UpdateUploadForRotation(appCtx, uploadID, params.Body.Rotation)
			if err != nil {
				return nil, apperror.NewBadDataError("unable to update upload")
			}

			url, err := h.FileStorer().PresignedURL(newUpload.StorageKey, newUpload.ContentType, newUpload.Filename)
			if err != nil {
				return nil, err
			}

			uploadPayload := payloads.Upload(h.FileStorer(), *newUpload, url)

			return uploadop.NewUpdateUploadCreated().WithPayload(uploadPayload), nil
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
			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().IsOfficeApp() {
				forbiddenError := apperror.NewForbiddenError("User is not an Office User.")
				appCtx.Logger().Error(forbiddenError.Error())
				return uploadop.NewDeleteUploadForbidden(), forbiddenError
			}

			uploadID, _ := uuid.FromString(params.UploadID.String())
			userUpload, err := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), uploadID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			userUploader, err := uploader.NewUserUploader(
				h.FileStorer(),
				uploader.MaxCustomerUserUploadFileSizeLimit,
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

// UploadStatusHandler returns status of an upload
type GetUploadStatusHandler struct {
	handlers.HandlerConfig
	services.UploadInformationFetcher
}

type CustomGetUploadStatusResponse struct {
	params     uploadop.GetUploadStatusParams
	storageKey string
	appCtx     appcontext.AppContext
	receiver   notifications.NotificationReceiver
	storer     storage.FileStorer
}

func (o *CustomGetUploadStatusResponse) writeEventStreamMessage(rw http.ResponseWriter, producer runtime.Producer, id int, event string, data string) {
	resProcess := []byte(fmt.Sprintf("id: %s\nevent: %s\ndata: %s\n\n", strconv.Itoa(id), event, data))
	if produceErr := producer.Produce(rw, resProcess); produceErr != nil {
		o.appCtx.Logger().Error(produceErr.Error())
	}
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}
}

func (o *CustomGetUploadStatusResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// Check current tag before event-driven wait for anti-virus
	tags, err := o.storer.Tags(o.storageKey)
	var uploadStatus models.AVStatusType
	if err != nil {
		uploadStatus = models.AVStatusPROCESSING
	} else {
		uploadStatus = models.GetAVStatusFromTags(tags)
	}

	// Limitation: once the status code header has been written (first response), we are not able to update the status for subsequent responses.
	// Standard 200 OK used with common SSE paradigm
	rw.WriteHeader(http.StatusOK)
	if uploadStatus == models.AVStatusCLEAN || uploadStatus == models.AVStatusINFECTED || uploadStatus == models.ClamAVStatusCLEAN || uploadStatus == models.ClamAVStatusINFECTED {
		o.writeEventStreamMessage(rw, producer, 0, "message", string(uploadStatus))
		o.writeEventStreamMessage(rw, producer, 1, "close", "Connection closed")
		return // skip notification loop since object already tagged from anti-virus
	} else {
		o.writeEventStreamMessage(rw, producer, 0, "message", string(uploadStatus))
	}

	// Start waiting for tag updates
	topicName, err := o.receiver.GetDefaultTopic()
	if err != nil {
		o.appCtx.Logger().Error(err.Error())
	}

	filterPolicy := fmt.Sprintf(`{
		"detail": {
				"object": {
					"key": [
						{"suffix": "%s"}
					]
				}
			}
	}`, o.params.UploadID)

	notificationParams := notifications.NotificationQueueParams{
		SubscriptionTopicName: topicName,
		NamePrefix:            notifications.QueuePrefixObjectTagsAdded,
		FilterPolicy:          filterPolicy,
	}

	queueUrl, err := o.receiver.CreateQueueWithSubscription(o.appCtx, notificationParams)
	if err != nil {
		o.appCtx.Logger().Error(err.Error())
	}

	id_counter := 1

	// For loop over 120 seconds, cancel context when done and it breaks the loop
	totalReceiverContext, totalReceiverContextCancelFunc := context.WithTimeout(context.Background(), 120*time.Second)
	defer func() {
		id_counter++
		o.writeEventStreamMessage(rw, producer, id_counter, "close", "Connection closed")
		totalReceiverContextCancelFunc()
	}()

	// Cleanup if client closes connection
	go func() {
		<-o.params.HTTPRequest.Context().Done()
		totalReceiverContextCancelFunc()
	}()

	// Cleanup at end of work
	go func() {
		<-totalReceiverContext.Done()
		_ = o.receiver.CloseoutQueue(o.appCtx, queueUrl)
	}()

	for {
		o.appCtx.Logger().Info("Receiving Messages...")
		messages, errs := o.receiver.ReceiveMessages(o.appCtx, queueUrl, totalReceiverContext)

		if errors.Is(errs, context.Canceled) || errors.Is(errs, context.DeadlineExceeded) {
			return
		}
		if errs != nil {
			o.appCtx.Logger().Error(err.Error())
			return
		}

		if len(messages) != 0 {
			errTransaction := o.appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

				tags, err := o.storer.Tags(o.storageKey)

				if err != nil {
					uploadStatus = models.AVStatusPROCESSING
				} else {
					uploadStatus = models.GetAVStatusFromTags(tags)
				}

				o.writeEventStreamMessage(rw, producer, id_counter, "message", string(uploadStatus))

				if uploadStatus == models.AVStatusCLEAN || uploadStatus == models.AVStatusINFECTED || uploadStatus == models.ClamAVStatusCLEAN || uploadStatus == models.ClamAVStatusINFECTED {
					return errors.New("connection_closed")
				}

				return err
			})

			if errTransaction != nil && errTransaction.Error() == "connection_closed" {
				return
			}

			if errTransaction != nil {
				o.appCtx.Logger().Error(err.Error())
				return
			}
		}
		id_counter++

		select {
		case <-totalReceiverContext.Done():
			return
		default:
			time.Sleep(1 * time.Second) // Throttle as a precaution against hounding of the SDK
			continue
		}
	}
}

// Handle returns status of an upload
func (h GetUploadStatusHandler) Handle(params uploadop.GetUploadStatusParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetUploadStatusHandler error", zap.Error(err))
				switch errors.Cause(err) {
				case models.ErrFetchForbidden:
					return uploadop.NewGetUploadStatusForbidden(), err
				case models.ErrFetchNotFound:
					return uploadop.NewGetUploadStatusNotFound(), err
				default:
					return uploadop.NewGetUploadStatusInternalServerError(), err
				}
			}

			uploadId := params.UploadID.String()
			uploadUUID, err := uuid.FromString(uploadId)
			if err != nil {
				return handleError(err)
			}

			uploaded, err := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), uploadUUID)
			if err != nil {
				return handleError(err)
			}

			return &CustomGetUploadStatusResponse{
				params:     params,
				storageKey: uploaded.Upload.StorageKey,
				appCtx:     h.AppContextFromRequest(params.HTTPRequest),
				receiver:   h.NotificationReceiver(),
				storer:     h.FileStorer(),
			}, nil
		})
}

type CreatePPMUploadHandler struct {
	handlers.HandlerConfig
	services.WeightTicketGenerator
	services.WeightTicketComputer
	*uploader.UserUploader
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
				zap.String("OfficeUserID", appCtx.Session().OfficeUserID.String()),
				zap.Int64("size", file.Header.Size),
			)

			documentID := uuid.FromStringOrNil(params.DocumentID.String())

			document, docErr := models.FetchDocumentWithNoRestrictions(appCtx.DB(), appCtx.Session(), documentID)
			if docErr != nil {
				docNotFoundErr := fmt.Errorf("documentId %q was not found", documentID)
				return ppmop.NewCreatePPMUploadNotFound().WithPayload(&ghcmessages.Error{
					Message: handlers.FmtString(docNotFoundErr.Error()),
				}), docNotFoundErr
			}

			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			// Ensure the document belongs to an association of the PPM shipment
			shipErr := ppmshipment.FindPPMShipmentWithDocument(appCtx, ppmShipmentID, documentID)
			if shipErr != nil {
				docNotFoundErr := fmt.Errorf("documentId %q was not found for this shipment", documentID)
				return ppmop.NewCreatePPMUploadNotFound().WithPayload(&ghcmessages.Error{
					Message: handlers.FmtString(docNotFoundErr.Error()),
				}), docNotFoundErr
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
				// if isWeightEstimatorFile is false, throw an error, and send message to front end to let user know.
				if !isWeightEstimatorFile {
					message := "The uploaded .xlsx file does not match the expected weight estimator file format."
					return ppmop.NewCreatePPMUploadForbidden().WithPayload(&ghcmessages.Error{
						Message: &message,
					}), nil
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
					cleanupErr := h.WeightTicketGenerator.CleanupFile(aFile)

					if cleanupErr != nil {
						appCtx.Logger().Warn("failed to cleanup weight ticket file", zap.Error(cleanupErr), zap.String("verrs", verrs.Error()))
					}

					return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
				}

				// we already generated an afero file so we can skip that process the wrapper method does
				newUserUpload, verrs, createErr = h.UserUploader.CreateUserUploadForDocument(appCtx, &document.ID, appCtx.Session().UserID, uploaderpkg.File{File: aFile}, uploaderpkg.AllowedTypesPPMDocuments)

				if verrs.HasAny() || createErr != nil {
					appCtx.Logger().Error("failed to create new user upload", zap.Error(createErr), zap.String("verrs", verrs.Error()))
					cleanupErr := h.WeightTicketGenerator.CleanupFile(aFile)

					if cleanupErr != nil {
						appCtx.Logger().Warn("failed to cleanup weight ticket file", zap.Error(cleanupErr), zap.String("verrs", verrs.Error()))
						return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
					}
					switch createErr.(type) {
					case uploaderpkg.ErrUnsupportedContentType:
						return ppmop.NewCreatePPMUploadUnprocessableEntity().WithPayload(payloadForValidationError(
							createErr.Error(),
							"createPPMUpload",
							uuid.Nil,
							verrs)), createErr
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

				newUserUpload, verrs1, err := h.UserUploader.UpdateUserXlsxUploadFilename(appCtx, newUserUpload, filename)

				if verrs1.HasAny() || err != nil {
					appCtx.Logger().Error("failed to rename uploaded filename", zap.Error(createErr), zap.String("verrs", verrs.Error()))
				}

				err = h.WeightTicketGenerator.CleanupFile(aFile)

				if err != nil {
					appCtx.Logger().Warn("failed to cleanup weight ticket file", zap.Error(err), zap.String("verrs", verrs.Error()))
					return ppmop.NewCreatePPMUploadInternalServerError(), rollbackErr
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
						return ppmop.NewCreatePPMUploadUnprocessableEntity().WithPayload(payloadForValidationError(
							"createPPMUpload",
							createErr.Error(),
							uuid.Nil,
							verrs)), createErr
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
