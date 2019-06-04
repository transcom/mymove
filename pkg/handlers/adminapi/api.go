package adminapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/adminapi"
	adminops "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/user"
)

// NewAdminAPIHandler returns a handler for the admin API
func NewAdminAPIHandler(context handlers.HandlerContext) http.Handler {

	// Wire up the handlers to the publicAPIMux
	adminSpec, err := loads.Analyzed(adminapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	adminAPI := adminops.NewMymoveAPI(adminSpec)

	queryBuilder := query.NewQueryBuilder(context.DB())
	adminAPI.OfficeIndexOfficeUsersHandler = IndexOfficeUsersHandler{
		HandlerContext:        context,
		NewQueryFilter:        query.NewQueryFilter,
		OfficeUserListFetcher: user.NewOfficeUserListFetcher(queryBuilder),
	}

	return adminAPI.Serve(nil)
}
