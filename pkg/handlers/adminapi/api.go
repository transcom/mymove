package adminapi

import (
	"log"
	"net/http"

	"github.com/transcom/mymove/pkg/services/organization"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/adminapi"
	adminops "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/handlers"
	accesscodeservice "github.com/transcom/mymove/pkg/services/accesscode"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	electronicorder "github.com/transcom/mymove/pkg/services/electronic_order"
	"github.com/transcom/mymove/pkg/services/office"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	tspop "github.com/transcom/mymove/pkg/services/tsp"

	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/upload"
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

	adminAPI.OfficeUsersIndexOfficeUsersHandler = IndexOfficeUsersHandler{
		context,
		officeuser.NewOfficeUserListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OfficeUsersGetOfficeUserHandler = GetOfficeUserHandler{
		context,
		officeuser.NewOfficeUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.OfficeUsersCreateOfficeUserHandler = CreateOfficeUserHandler{
		context,
		officeuser.NewOfficeUserCreator(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.OfficeUsersUpdateOfficeUserHandler = UpdateOfficeUserHandler{
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

	adminAPI.OrganizationIndexOrganizationsHandler = IndexOrganizationsHandler{
		context,
		organization.NewOrganizationListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.TransportationServiceProviderPerformancesIndexTSPPsHandler = IndexTSPPsHandler{
		context,
		tspop.NewTransportationServiceProviderPerformanceListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.TransportationServiceProviderPerformancesGetTSPPHandler = GetTSPPHandler{
		context,
		tspop.NewTransportationServiceProviderPerformanceFetcher(queryBuilder),
		query.NewQueryFilter,
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

	adminAPI.AccessCodesIndexAccessCodesHandler = IndexAccessCodesHandler{
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

	adminAPI.AdminUsersGetAdminUserHandler = GetAdminUserHandler{
		context,
		adminuser.NewAdminUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersCreateAdminUserHandler = CreateAdminUserHandler{
		context,
		adminuser.NewAdminUserCreator(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersUpdateAdminUserHandler = UpdateAdminUserHandler{
		context,
		adminuser.NewAdminUserUpdater(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.UploadGetUploadHandler = GetUploadHandler{
		context,
		upload.NewUploadInformationFetcher(context.DB()),
	}

	return adminAPI.Serve(nil)
}
