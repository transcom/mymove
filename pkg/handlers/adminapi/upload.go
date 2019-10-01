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
		// Question, it looks like a better flow exists from upload -> document -> moveDocument -> move, but it's
		// too recursive to use, so upload -> document -> SM -> orders[0] -> moves[0] will have to work at the moment.
		ID:          *handlers.FmtUUID(u.ID),
		MoveLocator: *swag.String(u.Document.ServiceMember.Orders[0].Moves[0].Locator),
		Upload: &adminmessages.Upload{
			ContentType: *swag.String(u.ContentType),
			CreatedAt:   *handlers.FmtDateTime(u.CreatedAt),
			Filename:    *swag.String(u.Filename),
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

	payload := payloadForUploadModel(uploads[0])

	return uploadop.NewGetUploadOK().WithPayload(payload)
}