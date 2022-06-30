// This file is safe to edit. Once it exists it will not be overwritten

package internalapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/addresses"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/backup_contacts"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/calendar"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_locations"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/postal_codes"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/internal.yaml --api-package internaloperations --model-package internalmessages --server-package internalapi --principal interface{} --exclude-main

func configureFlags(api *internaloperations.MymoveAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *internaloperations.MymoveAPI) http.Handler {
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
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.BinProducer = runtime.ByteStreamProducer()
	api.JSONProducer = runtime.JSONProducer()

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// uploads.CreateUploadMaxParseMemory = 32 << 20
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// orders.UploadAmendedOrdersMaxParseMemory = 32 << 20

	if api.OfficeApproveMoveHandler == nil {
		api.OfficeApproveMoveHandler = office.ApproveMoveHandlerFunc(func(params office.ApproveMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation office.ApproveMove has not yet been implemented")
		})
	}
	if api.OfficeApprovePPMHandler == nil {
		api.OfficeApprovePPMHandler = office.ApprovePPMHandlerFunc(func(params office.ApprovePPMParams) middleware.Responder {
			return middleware.NotImplemented("operation office.ApprovePPM has not yet been implemented")
		})
	}
	if api.OfficeApproveReimbursementHandler == nil {
		api.OfficeApproveReimbursementHandler = office.ApproveReimbursementHandlerFunc(func(params office.ApproveReimbursementParams) middleware.Responder {
			return middleware.NotImplemented("operation office.ApproveReimbursement has not yet been implemented")
		})
	}
	if api.OfficeCancelMoveHandler == nil {
		api.OfficeCancelMoveHandler = office.CancelMoveHandlerFunc(func(params office.CancelMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation office.CancelMove has not yet been implemented")
		})
	}
	if api.DocumentsCreateDocumentHandler == nil {
		api.DocumentsCreateDocumentHandler = documents.CreateDocumentHandlerFunc(func(params documents.CreateDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation documents.CreateDocument has not yet been implemented")
		})
	}
	if api.MoveDocsCreateGenericMoveDocumentHandler == nil {
		api.MoveDocsCreateGenericMoveDocumentHandler = move_docs.CreateGenericMoveDocumentHandlerFunc(func(params move_docs.CreateGenericMoveDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation move_docs.CreateGenericMoveDocument has not yet been implemented")
		})
	}
	if api.MtoShipmentCreateMTOShipmentHandler == nil {
		api.MtoShipmentCreateMTOShipmentHandler = mto_shipment.CreateMTOShipmentHandlerFunc(func(params mto_shipment.CreateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.CreateMTOShipment has not yet been implemented")
		})
	}
	if api.MoveDocsCreateMovingExpenseDocumentHandler == nil {
		api.MoveDocsCreateMovingExpenseDocumentHandler = move_docs.CreateMovingExpenseDocumentHandlerFunc(func(params move_docs.CreateMovingExpenseDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation move_docs.CreateMovingExpenseDocument has not yet been implemented")
		})
	}
	if api.OrdersCreateOrdersHandler == nil {
		api.OrdersCreateOrdersHandler = orders.CreateOrdersHandlerFunc(func(params orders.CreateOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation orders.CreateOrders has not yet been implemented")
		})
	}
	if api.PpmCreatePPMAttachmentsHandler == nil {
		api.PpmCreatePPMAttachmentsHandler = ppm.CreatePPMAttachmentsHandlerFunc(func(params ppm.CreatePPMAttachmentsParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.CreatePPMAttachments has not yet been implemented")
		})
	}
	if api.PpmCreatePersonallyProcuredMoveHandler == nil {
		api.PpmCreatePersonallyProcuredMoveHandler = ppm.CreatePersonallyProcuredMoveHandlerFunc(func(params ppm.CreatePersonallyProcuredMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.CreatePersonallyProcuredMove has not yet been implemented")
		})
	}
	if api.ServiceMembersCreateServiceMemberHandler == nil {
		api.ServiceMembersCreateServiceMemberHandler = service_members.CreateServiceMemberHandlerFunc(func(params service_members.CreateServiceMemberParams) middleware.Responder {
			return middleware.NotImplemented("operation service_members.CreateServiceMember has not yet been implemented")
		})
	}
	if api.BackupContactsCreateServiceMemberBackupContactHandler == nil {
		api.BackupContactsCreateServiceMemberBackupContactHandler = backup_contacts.CreateServiceMemberBackupContactHandlerFunc(func(params backup_contacts.CreateServiceMemberBackupContactParams) middleware.Responder {
			return middleware.NotImplemented("operation backup_contacts.CreateServiceMemberBackupContact has not yet been implemented")
		})
	}
	if api.CertificationCreateSignedCertificationHandler == nil {
		api.CertificationCreateSignedCertificationHandler = certification.CreateSignedCertificationHandlerFunc(func(params certification.CreateSignedCertificationParams) middleware.Responder {
			return middleware.NotImplemented("operation certification.CreateSignedCertification has not yet been implemented")
		})
	}
	if api.UploadsCreateUploadHandler == nil {
		api.UploadsCreateUploadHandler = uploads.CreateUploadHandlerFunc(func(params uploads.CreateUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.CreateUpload has not yet been implemented")
		})
	}
	if api.MoveDocsCreateWeightTicketDocumentHandler == nil {
		api.MoveDocsCreateWeightTicketDocumentHandler = move_docs.CreateWeightTicketDocumentHandlerFunc(func(params move_docs.CreateWeightTicketDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation move_docs.CreateWeightTicketDocument has not yet been implemented")
		})
	}
	if api.MoveDocsDeleteMoveDocumentHandler == nil {
		api.MoveDocsDeleteMoveDocumentHandler = move_docs.DeleteMoveDocumentHandlerFunc(func(params move_docs.DeleteMoveDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation move_docs.DeleteMoveDocument has not yet been implemented")
		})
	}
	if api.MtoShipmentDeleteShipmentHandler == nil {
		api.MtoShipmentDeleteShipmentHandler = mto_shipment.DeleteShipmentHandlerFunc(func(params mto_shipment.DeleteShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.DeleteShipment has not yet been implemented")
		})
	}
	if api.UploadsDeleteUploadHandler == nil {
		api.UploadsDeleteUploadHandler = uploads.DeleteUploadHandlerFunc(func(params uploads.DeleteUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.DeleteUpload has not yet been implemented")
		})
	}
	if api.UploadsDeleteUploadsHandler == nil {
		api.UploadsDeleteUploadsHandler = uploads.DeleteUploadsHandlerFunc(func(params uploads.DeleteUploadsParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.DeleteUploads has not yet been implemented")
		})
	}
	if api.EntitlementsIndexEntitlementsHandler == nil {
		api.EntitlementsIndexEntitlementsHandler = entitlements.IndexEntitlementsHandlerFunc(func(params entitlements.IndexEntitlementsParams) middleware.Responder {
			return middleware.NotImplemented("operation entitlements.IndexEntitlements has not yet been implemented")
		})
	}
	if api.MoveDocsIndexMoveDocumentsHandler == nil {
		api.MoveDocsIndexMoveDocumentsHandler = move_docs.IndexMoveDocumentsHandlerFunc(func(params move_docs.IndexMoveDocumentsParams) middleware.Responder {
			return middleware.NotImplemented("operation move_docs.IndexMoveDocuments has not yet been implemented")
		})
	}
	if api.PpmIndexPersonallyProcuredMovesHandler == nil {
		api.PpmIndexPersonallyProcuredMovesHandler = ppm.IndexPersonallyProcuredMovesHandlerFunc(func(params ppm.IndexPersonallyProcuredMovesParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.IndexPersonallyProcuredMoves has not yet been implemented")
		})
	}
	if api.BackupContactsIndexServiceMemberBackupContactsHandler == nil {
		api.BackupContactsIndexServiceMemberBackupContactsHandler = backup_contacts.IndexServiceMemberBackupContactsHandlerFunc(func(params backup_contacts.IndexServiceMemberBackupContactsParams) middleware.Responder {
			return middleware.NotImplemented("operation backup_contacts.IndexServiceMemberBackupContacts has not yet been implemented")
		})
	}
	if api.CertificationIndexSignedCertificationHandler == nil {
		api.CertificationIndexSignedCertificationHandler = certification.IndexSignedCertificationHandlerFunc(func(params certification.IndexSignedCertificationParams) middleware.Responder {
			return middleware.NotImplemented("operation certification.IndexSignedCertification has not yet been implemented")
		})
	}
	if api.UsersIsLoggedInUserHandler == nil {
		api.UsersIsLoggedInUserHandler = users.IsLoggedInUserHandlerFunc(func(params users.IsLoggedInUserParams) middleware.Responder {
			return middleware.NotImplemented("operation users.IsLoggedInUser has not yet been implemented")
		})
	}
	if api.MtoShipmentListMTOShipmentsHandler == nil {
		api.MtoShipmentListMTOShipmentsHandler = mto_shipment.ListMTOShipmentsHandlerFunc(func(params mto_shipment.ListMTOShipmentsParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.ListMTOShipments has not yet been implemented")
		})
	}
	if api.MovesPatchMoveHandler == nil {
		api.MovesPatchMoveHandler = moves.PatchMoveHandlerFunc(func(params moves.PatchMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.PatchMove has not yet been implemented")
		})
	}
	if api.PpmPatchPersonallyProcuredMoveHandler == nil {
		api.PpmPatchPersonallyProcuredMoveHandler = ppm.PatchPersonallyProcuredMoveHandlerFunc(func(params ppm.PatchPersonallyProcuredMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.PatchPersonallyProcuredMove has not yet been implemented")
		})
	}
	if api.ServiceMembersPatchServiceMemberHandler == nil {
		api.ServiceMembersPatchServiceMemberHandler = service_members.PatchServiceMemberHandlerFunc(func(params service_members.PatchServiceMemberParams) middleware.Responder {
			return middleware.NotImplemented("operation service_members.PatchServiceMember has not yet been implemented")
		})
	}
	if api.PpmRequestPPMExpenseSummaryHandler == nil {
		api.PpmRequestPPMExpenseSummaryHandler = ppm.RequestPPMExpenseSummaryHandlerFunc(func(params ppm.RequestPPMExpenseSummaryParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.RequestPPMExpenseSummary has not yet been implemented")
		})
	}
	if api.PpmRequestPPMPaymentHandler == nil {
		api.PpmRequestPPMPaymentHandler = ppm.RequestPPMPaymentHandlerFunc(func(params ppm.RequestPPMPaymentParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.RequestPPMPayment has not yet been implemented")
		})
	}
	if api.DutyLocationsSearchDutyLocationsHandler == nil {
		api.DutyLocationsSearchDutyLocationsHandler = duty_locations.SearchDutyLocationsHandlerFunc(func(params duty_locations.SearchDutyLocationsParams) middleware.Responder {
			return middleware.NotImplemented("operation duty_locations.SearchDutyLocations has not yet been implemented")
		})
	}
	if api.AddressesShowAddressHandler == nil {
		api.AddressesShowAddressHandler = addresses.ShowAddressHandlerFunc(func(params addresses.ShowAddressParams) middleware.Responder {
			return middleware.NotImplemented("operation addresses.ShowAddress has not yet been implemented")
		})
	}
	if api.CalendarShowAvailableMoveDatesHandler == nil {
		api.CalendarShowAvailableMoveDatesHandler = calendar.ShowAvailableMoveDatesHandlerFunc(func(params calendar.ShowAvailableMoveDatesParams) middleware.Responder {
			return middleware.NotImplemented("operation calendar.ShowAvailableMoveDates has not yet been implemented")
		})
	}
	if api.DocumentsShowDocumentHandler == nil {
		api.DocumentsShowDocumentHandler = documents.ShowDocumentHandlerFunc(func(params documents.ShowDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation documents.ShowDocument has not yet been implemented")
		})
	}
	if api.TransportationOfficesShowDutyLocationTransportationOfficeHandler == nil {
		api.TransportationOfficesShowDutyLocationTransportationOfficeHandler = transportation_offices.ShowDutyLocationTransportationOfficeHandlerFunc(func(params transportation_offices.ShowDutyLocationTransportationOfficeParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_offices.ShowDutyLocationTransportationOffice has not yet been implemented")
		})
	}
	if api.UsersShowLoggedInUserHandler == nil {
		api.UsersShowLoggedInUserHandler = users.ShowLoggedInUserHandlerFunc(func(params users.ShowLoggedInUserParams) middleware.Responder {
			return middleware.NotImplemented("operation users.ShowLoggedInUser has not yet been implemented")
		})
	}
	if api.MovesShowMoveHandler == nil {
		api.MovesShowMoveHandler = moves.ShowMoveHandlerFunc(func(params moves.ShowMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.ShowMove has not yet been implemented")
		})
	}
	if api.MovesShowMoveDatesSummaryHandler == nil {
		api.MovesShowMoveDatesSummaryHandler = moves.ShowMoveDatesSummaryHandlerFunc(func(params moves.ShowMoveDatesSummaryParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.ShowMoveDatesSummary has not yet been implemented")
		})
	}
	if api.OfficeShowOfficeOrdersHandler == nil {
		api.OfficeShowOfficeOrdersHandler = office.ShowOfficeOrdersHandlerFunc(func(params office.ShowOfficeOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation office.ShowOfficeOrders has not yet been implemented")
		})
	}
	if api.OrdersShowOrdersHandler == nil {
		api.OrdersShowOrdersHandler = orders.ShowOrdersHandlerFunc(func(params orders.ShowOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation orders.ShowOrders has not yet been implemented")
		})
	}
	if api.PpmShowPPMEstimateHandler == nil {
		api.PpmShowPPMEstimateHandler = ppm.ShowPPMEstimateHandlerFunc(func(params ppm.ShowPPMEstimateParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.ShowPPMEstimate has not yet been implemented")
		})
	}
	if api.PpmShowPPMIncentiveHandler == nil {
		api.PpmShowPPMIncentiveHandler = ppm.ShowPPMIncentiveHandlerFunc(func(params ppm.ShowPPMIncentiveParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.ShowPPMIncentive has not yet been implemented")
		})
	}
	if api.PpmShowPPMSitEstimateHandler == nil {
		api.PpmShowPPMSitEstimateHandler = ppm.ShowPPMSitEstimateHandlerFunc(func(params ppm.ShowPPMSitEstimateParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.ShowPPMSitEstimate has not yet been implemented")
		})
	}
	if api.PpmShowPersonallyProcuredMoveHandler == nil {
		api.PpmShowPersonallyProcuredMoveHandler = ppm.ShowPersonallyProcuredMoveHandlerFunc(func(params ppm.ShowPersonallyProcuredMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.ShowPersonallyProcuredMove has not yet been implemented")
		})
	}
	if api.QueuesShowQueueHandler == nil {
		api.QueuesShowQueueHandler = queues.ShowQueueHandlerFunc(func(params queues.ShowQueueParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.ShowQueue has not yet been implemented")
		})
	}
	if api.ServiceMembersShowServiceMemberHandler == nil {
		api.ServiceMembersShowServiceMemberHandler = service_members.ShowServiceMemberHandlerFunc(func(params service_members.ShowServiceMemberParams) middleware.Responder {
			return middleware.NotImplemented("operation service_members.ShowServiceMember has not yet been implemented")
		})
	}
	if api.BackupContactsShowServiceMemberBackupContactHandler == nil {
		api.BackupContactsShowServiceMemberBackupContactHandler = backup_contacts.ShowServiceMemberBackupContactHandlerFunc(func(params backup_contacts.ShowServiceMemberBackupContactParams) middleware.Responder {
			return middleware.NotImplemented("operation backup_contacts.ShowServiceMemberBackupContact has not yet been implemented")
		})
	}
	if api.ServiceMembersShowServiceMemberOrdersHandler == nil {
		api.ServiceMembersShowServiceMemberOrdersHandler = service_members.ShowServiceMemberOrdersHandlerFunc(func(params service_members.ShowServiceMemberOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation service_members.ShowServiceMemberOrders has not yet been implemented")
		})
	}
	if api.MovesShowShipmentSummaryWorksheetHandler == nil {
		api.MovesShowShipmentSummaryWorksheetHandler = moves.ShowShipmentSummaryWorksheetHandlerFunc(func(params moves.ShowShipmentSummaryWorksheetParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.ShowShipmentSummaryWorksheet has not yet been implemented")
		})
	}
	if api.MovesSubmitAmendedOrdersHandler == nil {
		api.MovesSubmitAmendedOrdersHandler = moves.SubmitAmendedOrdersHandlerFunc(func(params moves.SubmitAmendedOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.SubmitAmendedOrders has not yet been implemented")
		})
	}
	if api.MovesSubmitMoveForApprovalHandler == nil {
		api.MovesSubmitMoveForApprovalHandler = moves.SubmitMoveForApprovalHandlerFunc(func(params moves.SubmitMoveForApprovalParams) middleware.Responder {
			return middleware.NotImplemented("operation moves.SubmitMoveForApproval has not yet been implemented")
		})
	}
	if api.PpmSubmitPersonallyProcuredMoveHandler == nil {
		api.PpmSubmitPersonallyProcuredMoveHandler = ppm.SubmitPersonallyProcuredMoveHandlerFunc(func(params ppm.SubmitPersonallyProcuredMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.SubmitPersonallyProcuredMove has not yet been implemented")
		})
	}
	if api.MtoShipmentUpdateMTOShipmentHandler == nil {
		api.MtoShipmentUpdateMTOShipmentHandler = mto_shipment.UpdateMTOShipmentHandlerFunc(func(params mto_shipment.UpdateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.UpdateMTOShipment has not yet been implemented")
		})
	}
	if api.MoveDocsUpdateMoveDocumentHandler == nil {
		api.MoveDocsUpdateMoveDocumentHandler = move_docs.UpdateMoveDocumentHandlerFunc(func(params move_docs.UpdateMoveDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation move_docs.UpdateMoveDocument has not yet been implemented")
		})
	}
	if api.OrdersUpdateOrdersHandler == nil {
		api.OrdersUpdateOrdersHandler = orders.UpdateOrdersHandlerFunc(func(params orders.UpdateOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation orders.UpdateOrders has not yet been implemented")
		})
	}
	if api.PpmUpdatePersonallyProcuredMoveHandler == nil {
		api.PpmUpdatePersonallyProcuredMoveHandler = ppm.UpdatePersonallyProcuredMoveHandlerFunc(func(params ppm.UpdatePersonallyProcuredMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.UpdatePersonallyProcuredMove has not yet been implemented")
		})
	}
	if api.BackupContactsUpdateServiceMemberBackupContactHandler == nil {
		api.BackupContactsUpdateServiceMemberBackupContactHandler = backup_contacts.UpdateServiceMemberBackupContactHandlerFunc(func(params backup_contacts.UpdateServiceMemberBackupContactParams) middleware.Responder {
			return middleware.NotImplemented("operation backup_contacts.UpdateServiceMemberBackupContact has not yet been implemented")
		})
	}
	if api.OrdersUploadAmendedOrdersHandler == nil {
		api.OrdersUploadAmendedOrdersHandler = orders.UploadAmendedOrdersHandlerFunc(func(params orders.UploadAmendedOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation orders.UploadAmendedOrders has not yet been implemented")
		})
	}
	if api.PostalCodesValidatePostalCodeWithRateDataHandler == nil {
		api.PostalCodesValidatePostalCodeWithRateDataHandler = postal_codes.ValidatePostalCodeWithRateDataHandlerFunc(func(params postal_codes.ValidatePostalCodeWithRateDataParams) middleware.Responder {
			return middleware.NotImplemented("operation postal_codes.ValidatePostalCodeWithRateData has not yet been implemented")
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
