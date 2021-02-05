package adminapi

import (
	"log"
	"net/http"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	usersroles "github.com/transcom/mymove/pkg/services/users_roles"

	"github.com/transcom/mymove/pkg/services/organization"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/adminapi"
	adminops "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/handlers"
	accesscodeservice "github.com/transcom/mymove/pkg/services/accesscode"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	electronicorder "github.com/transcom/mymove/pkg/services/electronic_order"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	move "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/office"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	tspop "github.com/transcom/mymove/pkg/services/tsp"
	"github.com/transcom/mymove/pkg/services/upload"
	user "github.com/transcom/mymove/pkg/services/user"
	webhooksubscription "github.com/transcom/mymove/pkg/services/webhook_subscription"
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
	adminAPI.ServeError = handlers.ServeCustomError

	adminAPI.OfficeUsersIndexOfficeUsersHandler = IndexOfficeUsersHandler{
		context,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OfficeUsersGetOfficeUserHandler = GetOfficeUserHandler{
		context,
		officeuser.NewOfficeUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	userRolesCreator := usersroles.NewUsersRolesCreator(context.DB())
	adminAPI.OfficeUsersCreateOfficeUserHandler = CreateOfficeUserHandler{
		context,
		officeuser.NewOfficeUserCreator(context.DB(), queryBuilder),
		query.NewQueryFilter,
		userRolesCreator,
	}

	adminAPI.OfficeUsersUpdateOfficeUserHandler = UpdateOfficeUserHandler{
		context,
		officeuser.NewOfficeUserUpdater(queryBuilder),
		query.NewQueryFilter,
		userRolesCreator,
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

	adminAPI.UsersUpdateUserHandler = UpdateUserHandler{
		context,
		user.NewUserSessionRevocation(queryBuilder),
		user.NewUserUpdater(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersGetAdminUserHandler = GetAdminUserHandler{
		context,
		adminuser.NewAdminUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersCreateAdminUserHandler = CreateAdminUserHandler{
		context,
		adminuser.NewAdminUserCreator(context.DB(), queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersUpdateAdminUserHandler = UpdateAdminUserHandler{
		context,
		adminuser.NewAdminUserUpdater(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.UsersGetUserHandler = GetUserHandler{
		context,
		user.NewUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.UsersIndexUsersHandler = IndexUsersHandler{
		context,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}
	adminAPI.UploadGetUploadHandler = GetUploadHandler{
		context,
		upload.NewUploadInformationFetcher(context.DB()),
	}

	adminAPI.NotificationIndexNotificationsHandler = IndexNotificationsHandler{
		context,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.MoveIndexMovesHandler = IndexMovesHandler{
		context,
		move.NewMoveListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.MoveUpdateMoveHandler = UpdateMoveHandler{
		context,
		movetaskorder.NewMoveTaskOrderUpdater(context.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)),
	}

	adminAPI.MoveGetMoveHandler = GetMoveHandler{
		context,
	}

	adminAPI.WebhookSubscriptionsIndexWebhookSubscriptionsHandler = IndexWebhookSubscriptionsHandler{
		context,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.WebhookSubscriptionsGetWebhookSubscriptionHandler = GetWebhookSubscriptionHandler{
		context,
		webhooksubscription.NewWebhookSubscriptionFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	return adminAPI.Serve(nil)
}
