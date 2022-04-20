// This file is safe to edit. Once it exists it will not be overwritten

package adminapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/admin_users"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/electronic_order"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/move"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/notification"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/organization"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/transportation_service_provider_performances"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/upload"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/webhook_subscriptions"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/admin.yaml --api-package adminoperations --model-package adminmessages --server-package adminapi --principal interface{} --exclude-main

func configureFlags(api *adminoperations.MymoveAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *adminoperations.MymoveAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.AdminUsersCreateAdminUserHandler == nil {
		api.AdminUsersCreateAdminUserHandler = admin_users.CreateAdminUserHandlerFunc(func(params admin_users.CreateAdminUserParams) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.CreateAdminUser has not yet been implemented")
		})
	}
	if api.OfficeUsersCreateOfficeUserHandler == nil {
		api.OfficeUsersCreateOfficeUserHandler = office_users.CreateOfficeUserHandlerFunc(func(params office_users.CreateOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.CreateOfficeUser has not yet been implemented")
		})
	}
	if api.WebhookSubscriptionsCreateWebhookSubscriptionHandler == nil {
		api.WebhookSubscriptionsCreateWebhookSubscriptionHandler = webhook_subscriptions.CreateWebhookSubscriptionHandlerFunc(func(params webhook_subscriptions.CreateWebhookSubscriptionParams) middleware.Responder {
			return middleware.NotImplemented("operation webhook_subscriptions.CreateWebhookSubscription has not yet been implemented")
		})
	}
	if api.AdminUsersGetAdminUserHandler == nil {
		api.AdminUsersGetAdminUserHandler = admin_users.GetAdminUserHandlerFunc(func(params admin_users.GetAdminUserParams) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.GetAdminUser has not yet been implemented")
		})
	}
	if api.ElectronicOrderGetElectronicOrdersTotalsHandler == nil {
		api.ElectronicOrderGetElectronicOrdersTotalsHandler = electronic_order.GetElectronicOrdersTotalsHandlerFunc(func(params electronic_order.GetElectronicOrdersTotalsParams) middleware.Responder {
			return middleware.NotImplemented("operation electronic_order.GetElectronicOrdersTotals has not yet been implemented")
		})
	}
	if api.MoveGetMoveHandler == nil {
		api.MoveGetMoveHandler = move.GetMoveHandlerFunc(func(params move.GetMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation move.GetMove has not yet been implemented")
		})
	}
	if api.OfficeUsersGetOfficeUserHandler == nil {
		api.OfficeUsersGetOfficeUserHandler = office_users.GetOfficeUserHandlerFunc(func(params office_users.GetOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.GetOfficeUser has not yet been implemented")
		})
	}
	if api.TransportationServiceProviderPerformancesGetTSPPHandler == nil {
		api.TransportationServiceProviderPerformancesGetTSPPHandler = transportation_service_provider_performances.GetTSPPHandlerFunc(func(params transportation_service_provider_performances.GetTSPPParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_service_provider_performances.GetTSPP has not yet been implemented")
		})
	}
	if api.UploadGetUploadHandler == nil {
		api.UploadGetUploadHandler = upload.GetUploadHandlerFunc(func(params upload.GetUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation upload.GetUpload has not yet been implemented")
		})
	}
	if api.UsersGetUserHandler == nil {
		api.UsersGetUserHandler = users.GetUserHandlerFunc(func(params users.GetUserParams) middleware.Responder {
			return middleware.NotImplemented("operation users.GetUser has not yet been implemented")
		})
	}
	if api.WebhookSubscriptionsGetWebhookSubscriptionHandler == nil {
		api.WebhookSubscriptionsGetWebhookSubscriptionHandler = webhook_subscriptions.GetWebhookSubscriptionHandlerFunc(func(params webhook_subscriptions.GetWebhookSubscriptionParams) middleware.Responder {
			return middleware.NotImplemented("operation webhook_subscriptions.GetWebhookSubscription has not yet been implemented")
		})
	}
	if api.AdminUsersIndexAdminUsersHandler == nil {
		api.AdminUsersIndexAdminUsersHandler = admin_users.IndexAdminUsersHandlerFunc(func(params admin_users.IndexAdminUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.IndexAdminUsers has not yet been implemented")
		})
	}
	if api.ElectronicOrderIndexElectronicOrdersHandler == nil {
		api.ElectronicOrderIndexElectronicOrdersHandler = electronic_order.IndexElectronicOrdersHandlerFunc(func(params electronic_order.IndexElectronicOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation electronic_order.IndexElectronicOrders has not yet been implemented")
		})
	}
	if api.MoveIndexMovesHandler == nil {
		api.MoveIndexMovesHandler = move.IndexMovesHandlerFunc(func(params move.IndexMovesParams) middleware.Responder {
			return middleware.NotImplemented("operation move.IndexMoves has not yet been implemented")
		})
	}
	if api.NotificationIndexNotificationsHandler == nil {
		api.NotificationIndexNotificationsHandler = notification.IndexNotificationsHandlerFunc(func(params notification.IndexNotificationsParams) middleware.Responder {
			return middleware.NotImplemented("operation notification.IndexNotifications has not yet been implemented")
		})
	}
	if api.OfficeUsersIndexOfficeUsersHandler == nil {
		api.OfficeUsersIndexOfficeUsersHandler = office_users.IndexOfficeUsersHandlerFunc(func(params office_users.IndexOfficeUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.IndexOfficeUsers has not yet been implemented")
		})
	}
	if api.OfficeIndexOfficesHandler == nil {
		api.OfficeIndexOfficesHandler = office.IndexOfficesHandlerFunc(func(params office.IndexOfficesParams) middleware.Responder {
			return middleware.NotImplemented("operation office.IndexOffices has not yet been implemented")
		})
	}
	if api.OrganizationIndexOrganizationsHandler == nil {
		api.OrganizationIndexOrganizationsHandler = organization.IndexOrganizationsHandlerFunc(func(params organization.IndexOrganizationsParams) middleware.Responder {
			return middleware.NotImplemented("operation organization.IndexOrganizations has not yet been implemented")
		})
	}
	if api.TransportationServiceProviderPerformancesIndexTSPPsHandler == nil {
		api.TransportationServiceProviderPerformancesIndexTSPPsHandler = transportation_service_provider_performances.IndexTSPPsHandlerFunc(func(params transportation_service_provider_performances.IndexTSPPsParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_service_provider_performances.IndexTSPPs has not yet been implemented")
		})
	}
	if api.UsersIndexUsersHandler == nil {
		api.UsersIndexUsersHandler = users.IndexUsersHandlerFunc(func(params users.IndexUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation users.IndexUsers has not yet been implemented")
		})
	}
	if api.WebhookSubscriptionsIndexWebhookSubscriptionsHandler == nil {
		api.WebhookSubscriptionsIndexWebhookSubscriptionsHandler = webhook_subscriptions.IndexWebhookSubscriptionsHandlerFunc(func(params webhook_subscriptions.IndexWebhookSubscriptionsParams) middleware.Responder {
			return middleware.NotImplemented("operation webhook_subscriptions.IndexWebhookSubscriptions has not yet been implemented")
		})
	}
	if api.AdminUsersUpdateAdminUserHandler == nil {
		api.AdminUsersUpdateAdminUserHandler = admin_users.UpdateAdminUserHandlerFunc(func(params admin_users.UpdateAdminUserParams) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.UpdateAdminUser has not yet been implemented")
		})
	}
	if api.MoveUpdateMoveHandler == nil {
		api.MoveUpdateMoveHandler = move.UpdateMoveHandlerFunc(func(params move.UpdateMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation move.UpdateMove has not yet been implemented")
		})
	}
	if api.OfficeUsersUpdateOfficeUserHandler == nil {
		api.OfficeUsersUpdateOfficeUserHandler = office_users.UpdateOfficeUserHandlerFunc(func(params office_users.UpdateOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.UpdateOfficeUser has not yet been implemented")
		})
	}
	if api.UsersUpdateUserHandler == nil {
		api.UsersUpdateUserHandler = users.UpdateUserHandlerFunc(func(params users.UpdateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation users.UpdateUser has not yet been implemented")
		})
	}
	if api.WebhookSubscriptionsUpdateWebhookSubscriptionHandler == nil {
		api.WebhookSubscriptionsUpdateWebhookSubscriptionHandler = webhook_subscriptions.UpdateWebhookSubscriptionHandlerFunc(func(params webhook_subscriptions.UpdateWebhookSubscriptionParams) middleware.Responder {
			return middleware.NotImplemented("operation webhook_subscriptions.UpdateWebhookSubscription has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
