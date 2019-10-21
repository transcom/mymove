package adminapi

import (
	"github.com/transcom/mymove/pkg/services/upload"

	"github.com/go-openapi/strfmt"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime/middleware"

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
	handlers.HandlerContext
	services.UploadInformationFetcher
}

// Handle retrieves a specific upload
func (h GetUploadHandler) Handle(params uploadop.GetUploadParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	uploadID := uuid.FromStringOrNil(params.UploadID.String())
	uploadInformation, err := h.FetchUploadInformation(uploadID)
	if err != nil {
		switch err.(type) {
		case upload.ErrNotFound:
			return uploadop.NewGetUploadNotFound()
		default:
			return handlers.ResponseForError(logger, err)
		}
	}
	payload := payloadForUpload(uploadInformation)
	return uploadop.NewGetUploadOK().WithPayload(payload)
}
