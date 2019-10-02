package adminapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/adminapi"
	adminops "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/handlers"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	electronicorder "github.com/transcom/mymove/pkg/services/electronic_order"
	"github.com/transcom/mymove/pkg/services/office"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"

	accesscodeservice "github.com/transcom/mymove/pkg/services/accesscode"
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
		officeuser.NewOfficeUserListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OfficeGetOfficeUserHandler = GetOfficeUserHandler{
		context,
		officeuser.NewOfficeUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.OfficeCreateOfficeUserHandler = CreateOfficeUserHandler{
		context,
		officeuser.NewOfficeUserCreator(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.OfficeUpdateOfficeUserHandler = UpdateOfficeUserHandler{
		context,
		officeuser.NewOfficeUserUpdater(queryBuilder),
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

	adminAPI.ElectronicOrderGetElectronicOrdersTotalsHandler = GetElectronicOrdersTotalsHandler{
		context,
		electronicorder.NewElectronicOrdersCategoricalCountsFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.OfficeIndexAccessCodesHandler = IndexAccessCodesHandler{
		context,
		accesscodeservice.NewAccessCodeListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.AdminUsersIndexAdminUsersHandler = IndexAdminUsersHandler{
		context,
		adminuser.NewAdminUserListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	return adminAPI.Serve(nil)
}
