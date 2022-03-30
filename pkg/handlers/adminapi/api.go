package adminapi

import (
	"log"

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

// NewAdminAPI returns the admin API
func NewAdminAPI(ctx handlers.HandlerContext) *adminops.MymoveAPI {

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
		ctx,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OfficeUsersGetOfficeUserHandler = GetOfficeUserHandler{
		ctx,
		officeuser.NewOfficeUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	userRolesCreator := usersroles.NewUsersRolesCreator()
	adminAPI.OfficeUsersCreateOfficeUserHandler = CreateOfficeUserHandler{
		ctx,
		officeuser.NewOfficeUserCreator(queryBuilder, ctx.NotificationSender()),
		query.NewQueryFilter,
		userRolesCreator,
	}

	adminAPI.OfficeUsersUpdateOfficeUserHandler = UpdateOfficeUserHandler{
		ctx,
		officeUpdater,
		query.NewQueryFilter,
		userRolesCreator,
		user.NewUserSessionRevocation(queryBuilder),
	}

	adminAPI.OfficeIndexOfficesHandler = IndexOfficesHandler{
		ctx,
		office.NewOfficeListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OrganizationIndexOrganizationsHandler = IndexOrganizationsHandler{
		ctx,
		organization.NewOrganizationListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.TransportationServiceProviderPerformancesIndexTSPPsHandler = IndexTSPPsHandler{
		ctx,
		tspop.NewTransportationServiceProviderPerformanceListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.TransportationServiceProviderPerformancesGetTSPPHandler = GetTSPPHandler{
		ctx,
		tspop.NewTransportationServiceProviderPerformanceFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.ElectronicOrderIndexElectronicOrdersHandler = IndexElectronicOrdersHandler{
		ctx,
		electronicorder.NewElectronicOrderListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.ElectronicOrderGetElectronicOrdersTotalsHandler = GetElectronicOrdersTotalsHandler{
		ctx,
		electronicorder.NewElectronicOrdersCategoricalCountsFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AccessCodesIndexAccessCodesHandler = IndexAccessCodesHandler{
		ctx,
		accesscodeservice.NewAccessCodeListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.AdminUsersIndexAdminUsersHandler = IndexAdminUsersHandler{
		ctx,
		adminuser.NewAdminUserListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.UsersUpdateUserHandler = UpdateUserHandler{
		ctx,
		user.NewUserSessionRevocation(queryBuilder),
		user.NewUserUpdater(queryBuilder, officeUpdater, adminUpdater, ctx.NotificationSender()),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersGetAdminUserHandler = GetAdminUserHandler{
		ctx,
		adminuser.NewAdminUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersCreateAdminUserHandler = CreateAdminUserHandler{
		ctx,
		adminuser.NewAdminUserCreator(queryBuilder, ctx.NotificationSender()),
		query.NewQueryFilter,
	}

	adminAPI.AdminUsersUpdateAdminUserHandler = UpdateAdminUserHandler{
		ctx,
		adminUpdater,
		query.NewQueryFilter,
	}

	adminAPI.UsersGetUserHandler = GetUserHandler{
		ctx,
		user.NewUserFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.UsersIndexUsersHandler = IndexUsersHandler{
		ctx,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}
	adminAPI.UploadGetUploadHandler = GetUploadHandler{
		ctx,
		upload.NewUploadInformationFetcher(),
	}

	adminAPI.NotificationIndexNotificationsHandler = IndexNotificationsHandler{
		ctx,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.MoveIndexMovesHandler = IndexMovesHandler{
		ctx,
		move.NewMoveListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	moveRouter := move.NewMoveRouter()
	adminAPI.MoveUpdateMoveHandler = UpdateMoveHandler{
		ctx,
		movetaskorder.NewMoveTaskOrderUpdater(
			queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			moveRouter,
		),
	}

	adminAPI.MoveGetMoveHandler = GetMoveHandler{
		ctx,
	}

	adminAPI.WebhookSubscriptionsIndexWebhookSubscriptionsHandler = IndexWebhookSubscriptionsHandler{
		ctx,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.WebhookSubscriptionsGetWebhookSubscriptionHandler = GetWebhookSubscriptionHandler{
		ctx,
		webhooksubscription.NewWebhookSubscriptionFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.WebhookSubscriptionsCreateWebhookSubscriptionHandler = CreateWebhookSubscriptionHandler{
		ctx,
		webhooksubscription.NewWebhookSubscriptionCreator(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.WebhookSubscriptionsUpdateWebhookSubscriptionHandler = UpdateWebhookSubscriptionHandler{
		ctx,
		webhooksubscription.NewWebhookSubscriptionUpdater(queryBuilder),
		query.NewQueryFilter,
	}

	return adminAPI
}
