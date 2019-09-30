package adminapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	//"runtime"

	uploadop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/upload"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForUploadModel(u models.Upload) *adminmessages.UploadInformation {
	return &adminmessages.UploadInformation{
		// Question, if a service member can have multiple orders, can this break?
		MoveID: *handlers.FmtUUID(u.Document.ServiceMember.Orders[0].Moves[0].ID),
		Upload: &adminmessages.Upload{
			ContentType: *swag.String(u.ContentType),
			CreatedAt:   *handlers.FmtDateTime(u.CreatedAt),
			Filename:    *swag.String(u.Filename),
			ID:          *handlers.FmtUUID(u.ID),
			Size:        u.Bytes,
		},
	}
}

// GetUploadHandler returns an upload via GET /uploads/{uploadID}
type GetUploadHandler struct {
	handlers.HandlerContext
	services.UploadFetcher
	services.NewQueryFilter
}

// Handle retrieves a specific upload
func (h GetUploadHandler) Handle(params uploadop.GetUploadParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	uploadID := params.UploadID

	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", uploadID)}

	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("Document.ServiceMember.Orders.Moves"),
	}
	associations := query.NewQueryAssociations(queryAssociations)
	uploads, err := h.UploadFetcher.FetchUploads(queryFilters, associations)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	//fmt.Println("sm: ", uploads[0].Document.ServiceMember)
	//runtime.Breakpoint()
	payload := payloadForUploadModel(uploads[0])
	//payload := payloadForUploadModel(upload, upload.Document.ServiceMemberID)

	return uploadop.NewGetUploadOK().WithPayload(payload)
}