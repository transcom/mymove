package adminapi

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime/middleware"

	uploadop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/upload"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

//TODO:
// TODO update tests
// TODO extract fetchUploadInformation into existing service object maybe?
func payloadForUploadModel(u uploadInformation) *adminmessages.UploadInformation {
	return &adminmessages.UploadInformation{
		ID:          strfmt.UUID(u.UploadID.String()),
		MoveLocator: u.MoveLocator,
		Upload: &adminmessages.Upload{
			ContentType: u.ContentType,
			CreatedAt:   strfmt.DateTime(u.CreatedAt),
			Filename:    u.Filename,
			Size:        u.Bytes,
		},
		OfficeUserEmail: u.OfficeUserEmail,
		OfficeUserID:    handlers.FmtUUIDPtr(u.OfficeUserID),
		ServiceMemberID: handlers.FmtUUIDPtr(u.ServiceMemberID),
	}
}

// GetUploadHandler returns an upload via GET /uploads/{uploadID}
type GetUploadHandler struct {
	handlers.HandlerContext
	services.UploadFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a specific upload
func (h GetUploadHandler) Handle(params uploadop.GetUploadParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	uploadID := uuid.FromStringOrNil(params.UploadID.String())
	upload, err := h.fetchUploadInformation(uploadID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	payload := payloadForUploadModel(upload)

	return uploadop.NewGetUploadOK().WithPayload(payload)
}

type uploadInformation struct {
	UploadID        uuid.UUID `db:"upload_id"`
	Upload          models.Upload
	ContentType     string    `db:"content_type"`
	CreatedAt       time.Time `db:"created_at"`
	Filename        string
	Bytes           int64
	MoveLocator     string     `db:"locator"`
	ServiceMemberID *uuid.UUID `db:"service_member_id"`
	OfficeUserID    *uuid.UUID `db:"office_user_id"`
	OfficeUserEmail *string    `db:"office_user_id"`
}

// TODO extract to fetch upload
func (h GetUploadHandler) fetchUploadInformation(uploadID uuid.UUID) (uploadInformation, error) {
	q := `
SELECT uploads.id as upload_id,
       uploads.content_type,
       uploads.created_at,
       uploads.filename,
       uploads.bytes,
       moves.locator,
       sm.id AS service_member_id,
       ou.id AS office_user_id,
       ou.email AS office_user_id
FROM uploads
         JOIN users u ON uploads.uploader_id = u.id
         JOIN documents d ON uploads.document_id = d.id
         JOIN service_members documents_service_members ON d.service_member_id = documents_service_members.id
         JOIN orders ON documents_service_members.id = orders.service_member_id
         JOIN moves ON orders.id = moves.orders_id
         LEFT JOIN service_members sm ON u.id = sm.user_id
         LEFT JOIN office_users ou ON u.id = ou.user_id
WHERE uploads.id = ?`
	ui := uploadInformation{}
	err := h.DB().RawQuery(q, uploadID).First(&ui)
	if err != nil {
		return uploadInformation{}, err
	}
	return ui, nil
}
