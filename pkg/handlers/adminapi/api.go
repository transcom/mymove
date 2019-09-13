package adminapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/adminapi"
	adminops "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/handlers"
	electronicorder "github.com/transcom/mymove/pkg/services/electronic_order"
	"github.com/transcom/mymove/pkg/services/office"
	"github.com/transcom/mymove/pkg/services/pagination"
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
		context,
		user.NewOfficeUserListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OfficeGetOfficeUserHandler = GetOfficeUserHandler{
		context,
		user.NewOfficeUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.OfficeCreateOfficeUserHandler = CreateOfficeUserHandler{
		context,
		user.NewOfficeUserCreator(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.OfficePatchOfficeUserHandler = PatchOfficeUserHandler{
		context,
		user.NewOfficeUserUpdater(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.OfficeIndexOfficesHandler = IndexOfficesHandler{
		context,
		office.NewOfficeListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.ElectronicOrderIndexElectronicOrdersHandler = IndexElectronicOrdersHandler{
		context,
		electronicorder.NewElectronicOrderListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	return adminAPI.Serve(nil)
}
