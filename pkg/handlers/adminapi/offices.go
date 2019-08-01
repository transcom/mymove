package adminapi

import (
	"github.com/go-openapi/runtime/middleware"

	officeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// IndexOfficesHandler returns a list of office users via GET /office_users
type IndexOfficesHandler struct {
	handlers.HandlerContext
	services.NewQueryFilter
}

// Handle retrieves a list of office users
func (h IndexOfficesHandler) Handle(params officeop.IndexOfficesParams) middleware.Responder {
	return officeop.NewIndexOfficesOK().WithPayload(adminmessages.TransportationOffices{})
}
