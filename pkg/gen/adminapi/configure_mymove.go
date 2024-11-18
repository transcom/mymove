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
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/client_certificates"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/electronic_orders"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/moves"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/notifications"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/organizations"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/payment_request_syncada_file"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/payment_request_syncada_files"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/requested_office_users"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/transportation_offices"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/uploads"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/user"
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
	if api.ClientCertificatesCreateClientCertificateHandler == nil {
		api.ClientCertificatesCreateClientCertificateHandler = client_certificates.CreateClientCertificateHandlerFunc(func(params client_certificates.CreateClientCertificateParams) middleware.Responder {
			return middleware.NotImplemented("operation client_certificates.CreateClientCertificate has not yet been implemented")
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
	if api.ClientCertificatesGetClientCertificateHandler == nil {
		api.ClientCertificatesGetClientCertificateHandler = client_certificates.GetClientCertificateHandlerFunc(func(params client_certificates.GetClientCertificateParams) middleware.Responder {
			return middleware.NotImplemented("operation client_certificates.GetClientCertificate has not yet been implemented")
		})
	}
	if api.ElectronicOrdersGetElectronicOrdersTotalsHandler == nil {
		api.ElectronicOrdersGetElectronicOrdersTotalsHandler = electronic_orders.GetElectronicOrdersTotalsHandlerFunc(func(params electronic_orders.GetElectronicOrdersTotalsParams) middleware.Responder {
			return middleware.NotImplemented("operation electronic_orders.GetElectronicOrdersTotals has not yet been implemented")
		})
	}
	if api.UserGetLoggedInAdminUserHandler == nil {
		api.UserGetLoggedInAdminUserHandler = user.GetLoggedInAdminUserHandlerFunc(func(params user.GetLoggedInAdminUserParams) middleware.Responder {
			return middleware.NotImplemented("operation user.GetLoggedInAdminUser has not yet been implemented")
		})
	}
	if api.MovesGetMoveHandler == nil {
		api.MovesGetMoveHandler = moves.GetMoveHandlerFunc(func(params moves.GetMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.GetMove has not yet been implemented")
		})
	}
	if api.TransportationOfficesGetOfficeByIDHandler == nil {
		api.TransportationOfficesGetOfficeByIDHandler = transportation_offices.GetOfficeByIDHandlerFunc(func(params transportation_offices.GetOfficeByIDParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_offices.GetOfficeByID has not yet been implemented")
		})
	}
	if api.OfficeUsersGetOfficeUserHandler == nil {
		api.OfficeUsersGetOfficeUserHandler = office_users.GetOfficeUserHandlerFunc(func(params office_users.GetOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.GetOfficeUser has not yet been implemented")
		})
	}
	if api.RequestedOfficeUsersGetRequestedOfficeUserHandler == nil {
		api.RequestedOfficeUsersGetRequestedOfficeUserHandler = requested_office_users.GetRequestedOfficeUserHandlerFunc(func(params requested_office_users.GetRequestedOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation requested_office_users.GetRequestedOfficeUser has not yet been implemented")
		})
	}
	if api.UploadsGetUploadHandler == nil {
		api.UploadsGetUploadHandler = uploads.GetUploadHandlerFunc(func(params uploads.GetUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.GetUpload has not yet been implemented")
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
	if api.ClientCertificatesIndexClientCertificatesHandler == nil {
		api.ClientCertificatesIndexClientCertificatesHandler = client_certificates.IndexClientCertificatesHandlerFunc(func(params client_certificates.IndexClientCertificatesParams) middleware.Responder {
			return middleware.NotImplemented("operation client_certificates.IndexClientCertificates has not yet been implemented")
		})
	}
	if api.ElectronicOrdersIndexElectronicOrdersHandler == nil {
		api.ElectronicOrdersIndexElectronicOrdersHandler = electronic_orders.IndexElectronicOrdersHandlerFunc(func(params electronic_orders.IndexElectronicOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation electronic_orders.IndexElectronicOrders has not yet been implemented")
		})
	}
	if api.MovesIndexMovesHandler == nil {
		api.MovesIndexMovesHandler = moves.IndexMovesHandlerFunc(func(params moves.IndexMovesParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.IndexMoves has not yet been implemented")
		})
	}
	if api.NotificationsIndexNotificationsHandler == nil {
		api.NotificationsIndexNotificationsHandler = notifications.IndexNotificationsHandlerFunc(func(params notifications.IndexNotificationsParams) middleware.Responder {
			return middleware.NotImplemented("operation notifications.IndexNotifications has not yet been implemented")
		})
	}
	if api.OfficeUsersIndexOfficeUsersHandler == nil {
		api.OfficeUsersIndexOfficeUsersHandler = office_users.IndexOfficeUsersHandlerFunc(func(params office_users.IndexOfficeUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.IndexOfficeUsers has not yet been implemented")
		})
	}
	if api.TransportationOfficesIndexOfficesHandler == nil {
		api.TransportationOfficesIndexOfficesHandler = transportation_offices.IndexOfficesHandlerFunc(func(params transportation_offices.IndexOfficesParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_offices.IndexOffices has not yet been implemented")
		})
	}
	if api.OrganizationsIndexOrganizationsHandler == nil {
		api.OrganizationsIndexOrganizationsHandler = organizations.IndexOrganizationsHandlerFunc(func(params organizations.IndexOrganizationsParams) middleware.Responder {
			return middleware.NotImplemented("operation organizations.IndexOrganizations has not yet been implemented")
		})
	}
	if api.PaymentRequestSyncadaFilesIndexPaymentRequestSyncadaFilesHandler == nil {
		api.PaymentRequestSyncadaFilesIndexPaymentRequestSyncadaFilesHandler = payment_request_syncada_files.IndexPaymentRequestSyncadaFilesHandlerFunc(func(params payment_request_syncada_files.IndexPaymentRequestSyncadaFilesParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_request_syncada_files.IndexPaymentRequestSyncadaFiles has not yet been implemented")
		})
	}
	if api.RequestedOfficeUsersIndexRequestedOfficeUsersHandler == nil {
		api.RequestedOfficeUsersIndexRequestedOfficeUsersHandler = requested_office_users.IndexRequestedOfficeUsersHandlerFunc(func(params requested_office_users.IndexRequestedOfficeUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation requested_office_users.IndexRequestedOfficeUsers has not yet been implemented")
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
	if api.PaymentRequestSyncadaFilePaymentRequestSyncadaFileHandler == nil {
		api.PaymentRequestSyncadaFilePaymentRequestSyncadaFileHandler = payment_request_syncada_file.PaymentRequestSyncadaFileHandlerFunc(func(params payment_request_syncada_file.PaymentRequestSyncadaFileParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_request_syncada_file.PaymentRequestSyncadaFile has not yet been implemented")
		})
	}
	if api.ClientCertificatesRemoveClientCertificateHandler == nil {
		api.ClientCertificatesRemoveClientCertificateHandler = client_certificates.RemoveClientCertificateHandlerFunc(func(params client_certificates.RemoveClientCertificateParams) middleware.Responder {
			return middleware.NotImplemented("operation client_certificates.RemoveClientCertificate has not yet been implemented")
		})
	}
	if api.AdminUsersUpdateAdminUserHandler == nil {
		api.AdminUsersUpdateAdminUserHandler = admin_users.UpdateAdminUserHandlerFunc(func(params admin_users.UpdateAdminUserParams) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.UpdateAdminUser has not yet been implemented")
		})
	}
	if api.ClientCertificatesUpdateClientCertificateHandler == nil {
		api.ClientCertificatesUpdateClientCertificateHandler = client_certificates.UpdateClientCertificateHandlerFunc(func(params client_certificates.UpdateClientCertificateParams) middleware.Responder {
			return middleware.NotImplemented("operation client_certificates.UpdateClientCertificate has not yet been implemented")
		})
	}
	if api.MovesUpdateMoveHandler == nil {
		api.MovesUpdateMoveHandler = moves.UpdateMoveHandlerFunc(func(params moves.UpdateMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.UpdateMove has not yet been implemented")
		})
	}
	if api.OfficeUsersUpdateOfficeUserHandler == nil {
		api.OfficeUsersUpdateOfficeUserHandler = office_users.UpdateOfficeUserHandlerFunc(func(params office_users.UpdateOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.UpdateOfficeUser has not yet been implemented")
		})
	}
	if api.RequestedOfficeUsersUpdateRequestedOfficeUserHandler == nil {
		api.RequestedOfficeUsersUpdateRequestedOfficeUserHandler = requested_office_users.UpdateRequestedOfficeUserHandlerFunc(func(params requested_office_users.UpdateRequestedOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation requested_office_users.UpdateRequestedOfficeUser has not yet been implemented")
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
