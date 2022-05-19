package adminapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/adminapi"
	adminops "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/handlers"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	electronicorder "github.com/transcom/mymove/pkg/services/electronic_order"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	move "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/office"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/organization"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	tspop "github.com/transcom/mymove/pkg/services/tsp"
	"github.com/transcom/mymove/pkg/services/upload"
	user "github.com/transcom/mymove/pkg/services/user"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
	webhooksubscription "github.com/transcom/mymove/pkg/services/webhook_subscription"
)

// NewAdminAPI returns the admin API
func NewAdminAPI(handlerConfig handlers.HandlerConfig) *adminops.MymoveAPI {

	// Wire up the handlers to the publicAPIMux
	adminSpec, err := loads.Analyzed(adminapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	adminAPI := adminops.NewMymoveAPI(adminSpec)
	queryBuilder := query.NewQueryBuilder()
	officeUpdater := officeuser.NewOfficeUserUpdater(queryBuilder)
	adminUpdater := adminuser.NewAdminUserUpdater(queryBuilder)

	adminAPI.ServeError = handlers.ServeCustomError

	adminAPI.OfficeUsersIndexOfficeUsersHandler = IndexOfficeUsersHandler{
		handlerConfig,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OfficeUsersGetOfficeUserHandler = GetOfficeUserHandler{
		handlerConfig,
		officeuser.NewOfficeUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	userRolesCreator := usersroles.NewUsersRolesCreator()
	adminAPI.OfficeUsersCreateOfficeUserHandler = CreateOfficeUserHandler{
		handlerConfig,
		officeuser.NewOfficeUserCreator(queryBuilder, handlerConfig.NotificationSender()),
		query.NewQueryFilter,
		userRolesCreator,
	}

	adminAPI.OfficeUsersUpdateOfficeUserHandler = UpdateOfficeUserHandler{
		handlerConfig,
		officeUpdater,
		query.NewQueryFilter,
		userRolesCreator,
		user.NewUserSessionRevocation(queryBuilder),
	}

	adminAPI.OfficeIndexOfficesHandler = IndexOfficesHandler{
		handlerConfig,
		office.NewOfficeListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OrganizationIndexOrganizationsHandler = IndexOrganizationsHandler{
		handlerConfig,
		organization.NewOrganizationListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.TransportationServiceProviderPerformancesIndexTSPPsHandler = IndexTSPPsHandler{
		handlerConfig,
		tspop.NewTransportationServiceProviderPerformanceListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.TransportationServiceProviderPerformancesGetTSPPHandler = GetTSPPHandler{
		handlerConfig,
		tspop.NewTransportationServiceProviderPerformanceFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.ElectronicOrderIndexElectronicOrdersHandler = IndexElectronicOrdersHandler{
		handlerConfig,
		electronicorder.NewElectronicOrderListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.ElectronicOrderGetElectronicOrdersTotalsHandler = GetElectronicOrdersTotalsHandler{
		handlerConfig,
		electronicorder.NewElectronicOrdersCategoricalCountsFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersIndexAdminUsersHandler = IndexAdminUsersHandler{
		handlerConfig,
		adminuser.NewAdminUserListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.UsersUpdateUserHandler = UpdateUserHandler{
		handlerConfig,
		user.NewUserSessionRevocation(queryBuilder),
		user.NewUserUpdater(queryBuilder, officeUpdater, adminUpdater, handlerConfig.NotificationSender()),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersGetAdminUserHandler = GetAdminUserHandler{
		handlerConfig,
		adminuser.NewAdminUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersCreateAdminUserHandler = CreateAdminUserHandler{
		handlerConfig,
		adminuser.NewAdminUserCreator(queryBuilder, handlerConfig.NotificationSender()),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersUpdateAdminUserHandler = UpdateAdminUserHandler{
		handlerConfig,
		adminUpdater,
		query.NewQueryFilter,
	}

	adminAPI.UsersGetUserHandler = GetUserHandler{
		handlerConfig,
		user.NewUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.UsersIndexUsersHandler = IndexUsersHandler{
		handlerConfig,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}
	adminAPI.UploadGetUploadHandler = GetUploadHandler{
		handlerConfig,
		upload.NewUploadInformationFetcher(),
	}

	adminAPI.NotificationIndexNotificationsHandler = IndexNotificationsHandler{
		handlerConfig,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.MoveIndexMovesHandler = IndexMovesHandler{
		handlerConfig,
		move.NewMoveListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	moveRouter := move.NewMoveRouter()
	adminAPI.MoveUpdateMoveHandler = UpdateMoveHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderUpdater(
			queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			moveRouter,
		),
	}

	adminAPI.MoveGetMoveHandler = GetMoveHandler{
		handlerConfig,
	}

	adminAPI.WebhookSubscriptionsIndexWebhookSubscriptionsHandler = IndexWebhookSubscriptionsHandler{
		handlerConfig,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.WebhookSubscriptionsGetWebhookSubscriptionHandler = GetWebhookSubscriptionHandler{
		handlerConfig,
		webhooksubscription.NewWebhookSubscriptionFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.WebhookSubscriptionsCreateWebhookSubscriptionHandler = CreateWebhookSubscriptionHandler{
		handlerConfig,
		webhooksubscription.NewWebhookSubscriptionCreator(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.WebhookSubscriptionsUpdateWebhookSubscriptionHandler = UpdateWebhookSubscriptionHandler{
		handlerConfig,
		webhooksubscription.NewWebhookSubscriptionUpdater(queryBuilder),
		query.NewQueryFilter,
	}

	return adminAPI
}
