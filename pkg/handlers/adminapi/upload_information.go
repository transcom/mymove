package adminapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	uploadop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/upload"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForUpload(u services.UploadInformation) *adminmessages.UploadInformation {
	return &adminmessages.UploadInformation{
		ID:          strfmt.UUID(u.UploadID.String()),
		MoveLocator: u.MoveLocator,
		Upload: &adminmessages.Upload{
			ContentType: u.ContentType,
			CreatedAt:   strfmt.DateTime(u.CreatedAt),
			Filename:    u.Filename,
			Size:        u.Bytes,
		},
		OfficeUserID:           handlers.FmtUUIDPtr(u.OfficeUserID),
		OfficeUserEmail:        u.OfficeUserEmail,
		OfficeUserFirstName:    u.OfficeUserFirstName,
		OfficeUserLastName:     u.OfficeUserLastName,
		OfficeUserPhone:        u.OfficeUserPhone,
		ServiceMemberID:        handlers.FmtUUIDPtr(u.ServiceMemberID),
		ServiceMemberEmail:     u.ServiceMemberEmail,
		ServiceMemberFirstName: u.ServiceMemberFirstName,
		ServiceMemberLastName:  u.ServiceMemberLastName,
		ServiceMemberPhone:     u.ServiceMemberPhone,
	}
}

// GetUploadHandler returns an upload via GET /uploads/{uploadID}
type GetUploadHandler struct {
	handlers.HandlerConfig
	services.UploadInformationFetcher
}

// Handle retrieves a specific upload
func (h GetUploadHandler) Handle(params uploadop.GetUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			uploadID := uuid.FromStringOrNil(params.UploadID.String())
			uploadInformation, err := h.FetchUploadInformation(appCtx, uploadID)
			if err != nil {
				switch err.(type) {
				case apperror.NotFoundError:
					appCtx.Logger().Error("adminapi.GetUploadHandler not found error:", zap.Error(err))
					return uploadop.NewGetUploadNotFound(), err
				default:
					appCtx.Logger().Error("adminapi.GetUploadHandler error:", zap.Error(err))
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
			}
			payload := payloadForUpload(uploadInformation)
			return uploadop.NewGetUploadOK().WithPayload(payload), nil
		})
}
