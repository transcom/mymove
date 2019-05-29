package adminapi

import (
	"github.com/go-openapi/runtime/middleware"

	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/handlers"
)

// IndexOfficeUsersHandler returns a list of office users via GET /office_users
type IndexOfficeUsersHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a list of office users
func (h IndexOfficeUsersHandler) Handle(params officeuserop.IndexOfficeUsersParams) middleware.Responder {
	return officeuserop.NewIndexOfficeUsersOK()
}
