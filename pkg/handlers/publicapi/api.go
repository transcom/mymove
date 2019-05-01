package publicapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/paperwork"
	paperworkservice "github.com/transcom/mymove/pkg/services/paperwork"
	sitservice "github.com/transcom/mymove/pkg/services/storage_in_transit"
)

// NewPublicAPIHandler returns a handler for the public API
func NewPublicAPIHandler(context handlers.HandlerContext) http.Handler {

	// Wire up the handlers to the publicAPIMux
	apiSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	publicAPI := publicops.NewMymoveAPI(apiSpec)

	// Blackouts

	// Documents
	publicAPI.MoveDocsCreateGenericMoveDocumentHandler = CreateGenericMoveDocumentHandler{context}
	publicAPI.MoveDocsIndexMoveDocumentsHandler = IndexMoveDocumentsHandler{context}
	publicAPI.MoveDocsUpdateMoveDocumentHandler = UpdateMoveDocumentHandler{context}
	publicAPI.UploadsCreateUploadHandler = CreateUploadHandler{context}
	publicAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{context}

	// Shipments
	publicAPI.ShipmentsIndexShipmentsHandler = IndexShipmentsHandler{context}
	publicAPI.ShipmentsGetShipmentHandler = GetShipmentHandler{context}
	publicAPI.ShipmentsPatchShipmentHandler = PatchShipmentHandler{context}
	publicAPI.ShipmentsAcceptShipmentHandler = AcceptShipmentHandler{context}
	publicAPI.ShipmentsTransportShipmentHandler = TransportShipmentHandler{context}
	publicAPI.ShipmentsDeliverShipmentHandler = DeliverShipmentHandler{context}
	publicAPI.ShipmentsGetShipmentInvoicesHandler = GetShipmentInvoicesHandler{context}

	publicAPI.ShipmentsCompletePmSurveyHandler = CompletePmSurveyHandler{context}
	publicAPI.ShipmentsCreateGovBillOfLadingHandler = CreateGovBillOfLadingHandler{context, paperworkservice.NewFormCreator(context.FileStorer().TempFileSystem(), paperwork.NewFormFiller())}

	// Accessorials
	publicAPI.AccessorialsGetShipmentLineItemsHandler = GetShipmentLineItemsHandler{context}
	publicAPI.AccessorialsUpdateShipmentLineItemHandler = UpdateShipmentLineItemHandler{context}
	publicAPI.AccessorialsCreateShipmentLineItemHandler = CreateShipmentLineItemHandler{context}
	publicAPI.AccessorialsDeleteShipmentLineItemHandler = DeleteShipmentLineItemHandler{context}
	publicAPI.AccessorialsApproveShipmentLineItemHandler = ApproveShipmentLineItemHandler{context}

	publicAPI.AccessorialsGetTariff400ngItemsHandler = GetTariff400ngItemsHandler{context}
	publicAPI.AccessorialsGetInvoiceHandler = GetInvoiceHandler{context}

	// Service Agents
	publicAPI.ServiceAgentsIndexServiceAgentsHandler = IndexServiceAgentsHandler{context}
	publicAPI.ServiceAgentsCreateServiceAgentHandler = CreateServiceAgentHandler{context}
	publicAPI.ServiceAgentsPatchServiceAgentHandler = PatchServiceAgentHandler{context}

	// TSPs
	publicAPI.TransportationServiceProviderGetTransportationServiceProviderHandler = GetTransportationServiceProviderHandler{context}
	publicAPI.TspsIndexTSPsHandler = TspsIndexTSPsHandler{context}
	publicAPI.TspsGetTspShipmentsHandler = TspsGetTspShipmentsHandler{context}

	// Storage In Transits
	publicAPI.StorageInTransitsCreateStorageInTransitHandler = CreateStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitCreator(context.DB()),
	}
	publicAPI.StorageInTransitsGetStorageInTransitHandler = GetStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitByIDFetcher(context.DB()),
	}
	publicAPI.StorageInTransitsIndexStorageInTransitsHandler = IndexStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitIndexer(context.DB()),
	}
	publicAPI.StorageInTransitsDeleteStorageInTransitHandler = DeleteStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitDeleter(context.DB()),
	}
	publicAPI.StorageInTransitsPatchStorageInTransitHandler = PatchStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitPatcher(context.DB()),
	}
	publicAPI.StorageInTransitsApproveStorageInTransitHandler = ApproveStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitApprover(context.DB()),
	}
	publicAPI.StorageInTransitsDenyStorageInTransitHandler = DenyStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitDenier(context.DB()),
	}
	publicAPI.StorageInTransitsInSitStorageInTransitHandler = InSitStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitInSITPlacer(context.DB()),
	}
	publicAPI.StorageInTransitsReleaseStorageInTransitHandler = ReleaseStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitInReleaser(context.DB()),
	}
	publicAPI.StorageInTransitsDeliverStorageInTransitHandler = DeliverStorageInTransitHandler{
		context,
		sitservice.NewStorageInTransitInDeliverer(context.DB()),
	}

	return publicAPI.Serve(nil)
}
