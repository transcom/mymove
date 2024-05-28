package adminapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/adminapi"
	adminops "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/handlers"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	"github.com/transcom/mymove/pkg/services/clientcert"
	electronicorder "github.com/transcom/mymove/pkg/services/electronic_order"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	move "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/office"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/organization"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	requestedofficeusers "github.com/transcom/mymove/pkg/services/requested_office_users"
	"github.com/transcom/mymove/pkg/services/roles"
	"github.com/transcom/mymove/pkg/services/upload"
	user "github.com/transcom/mymove/pkg/services/user"
	usersprivileges "github.com/transcom/mymove/pkg/services/users_privileges"
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

	adminAPI.RequestedOfficeUsersIndexRequestedOfficeUsersHandler = IndexRequestedOfficeUsersHandler{
		handlerConfig,
		requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	userRolesCreator := usersroles.NewUsersRolesCreator()
	newRolesFetcher := roles.NewRolesFetcher()

	adminAPI.RequestedOfficeUsersGetRequestedOfficeUserHandler = GetRequestedOfficeUserHandler{
		handlerConfig,
		requestedofficeusers.NewRequestedOfficeUserFetcher(queryBuilder),
		newRolesFetcher,
		query.NewQueryFilter,
	}

	adminAPI.RequestedOfficeUsersUpdateRequestedOfficeUserHandler = UpdateRequestedOfficeUserHandler{
		handlerConfig,
		requestedofficeusers.NewRequestedOfficeUserUpdater(queryBuilder),
		userRolesCreator,
		newRolesFetcher,
	}

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

	userPrivilegesCreator := usersprivileges.NewUsersPrivilegesCreator()
	adminAPI.OfficeUsersCreateOfficeUserHandler = CreateOfficeUserHandler{
		handlerConfig,
		officeuser.NewOfficeUserCreator(queryBuilder, handlerConfig.NotificationSender()),
		query.NewQueryFilter,
		userRolesCreator,
		newRolesFetcher,
		userPrivilegesCreator,
	}

	adminAPI.OfficeUsersUpdateOfficeUserHandler = UpdateOfficeUserHandler{
		handlerConfig,
		officeUpdater,
		query.NewQueryFilter,
		userRolesCreator,
		userPrivilegesCreator,
		user.NewUserSessionRevocation(queryBuilder),
	}

	adminAPI.TransportationOfficesIndexOfficesHandler = IndexOfficesHandler{
		handlerConfig,
		office.NewOfficeListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.OrganizationsIndexOrganizationsHandler = IndexOrganizationsHandler{
		handlerConfig,
		organization.NewOrganizationListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.ElectronicOrdersIndexElectronicOrdersHandler = IndexElectronicOrdersHandler{
		handlerConfig,
		electronicorder.NewElectronicOrderListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.ElectronicOrdersGetElectronicOrdersTotalsHandler = GetElectronicOrdersTotalsHandler{
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
	adminAPI.UploadsGetUploadHandler = GetUploadHandler{
		handlerConfig,
		upload.NewUploadInformationFetcher(),
	}

	adminAPI.NotificationsIndexNotificationsHandler = IndexNotificationsHandler{
		handlerConfig,
		fetch.NewListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.MovesIndexMovesHandler = IndexMovesHandler{
		handlerConfig,
		move.NewMoveListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	moveRouter := move.NewMoveRouter()
	adminAPI.MovesUpdateMoveHandler = UpdateMoveHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderUpdater(
			queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter,
		),
	}

	adminAPI.MovesGetMoveHandler = GetMoveHandler{
		handlerConfig,
	}

	adminAPI.ClientCertificatesIndexClientCertificatesHandler = IndexClientCertsHandler{
		handlerConfig,
		clientcert.NewClientCertListFetcher(queryBuilder),
		query.NewQueryFilter,
		pagination.NewPagination,
	}

	adminAPI.ClientCertificatesGetClientCertificateHandler = GetClientCertHandler{
		handlerConfig,
		clientcert.NewClientCertFetcher(queryBuilder),
		query.NewQueryFilter,
	}

	adminAPI.ClientCertificatesCreateClientCertificateHandler = CreateClientCertHandler{
		handlerConfig,
		clientcert.NewClientCertCreator(queryBuilder,
			userRolesCreator, handlerConfig.NotificationSender()),
	}

	adminAPI.ClientCertificatesUpdateClientCertificateHandler = UpdateClientCertHandler{
		handlerConfig,
		clientcert.NewClientCertUpdater(queryBuilder, userRolesCreator, handlerConfig.NotificationSender()),
		query.NewQueryFilter,
	}

	adminAPI.ClientCertificatesRemoveClientCertificateHandler = RemoveClientCertHandler{
		handlerConfig,
		clientcert.NewClientCertRemover(queryBuilder, userRolesCreator, handlerConfig.NotificationSender()),
		query.NewQueryFilter,
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
