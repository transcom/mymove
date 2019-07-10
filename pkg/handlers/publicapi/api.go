package publicapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/rateengine"

	accesscodeservice "github.com/transcom/mymove/pkg/services/accesscode"
	paperworkservice "github.com/transcom/mymove/pkg/services/paperwork"
	postalcodeservice "github.com/transcom/mymove/pkg/services/postal_codes"
	shipmentservice "github.com/transcom/mymove/pkg/services/shipment"
	shipmentlineitemservice "github.com/transcom/mymove/pkg/services/shipment_line_item"
	sitservice "github.com/transcom/mymove/pkg/services/storage_in_transit"
)

// NewPublicAPIHandler returns a handler for the public API
func NewPublicAPIHandler(context handlers.HandlerContext, logger Logger) http.Handler {

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

	engine := rateengine.NewRateEngine(context.DB(), logger)
	publicAPI.ShipmentsDeliverShipmentHandler = DeliverShipmentHandler{
		context, shipmentservice.NewShipmentDeliverAndPricer(
			context.DB(),
			engine,
			context.Planner(),
		)}

	publicAPI.ShipmentsGetShipmentInvoicesHandler = GetShipmentInvoicesHandler{context}

	publicAPI.ShipmentsCompletePmSurveyHandler = CompletePmSurveyHandler{context}
	publicAPI.ShipmentsCreateGovBillOfLadingHandler = CreateGovBillOfLadingHandler{context, paperworkservice.NewFormCreator(context.FileStorer().TempFileSystem(), paperwork.NewFormFiller())}

	// Accessorials
	publicAPI.AccessorialsGetShipmentLineItemsHandler = GetShipmentLineItemsHandler{context, shipmentlineitemservice.NewShipmentLineItemFetcher(context.DB())}
	publicAPI.AccessorialsUpdateShipmentLineItemHandler = UpdateShipmentLineItemHandler{context}
	publicAPI.AccessorialsCreateShipmentLineItemHandler = CreateShipmentLineItemHandler{context}
	publicAPI.AccessorialsDeleteShipmentLineItemHandler = DeleteShipmentLineItemHandler{context}
	publicAPI.AccessorialsApproveShipmentLineItemHandler = ApproveShipmentLineItemHandler{context}
	publicAPI.AccessorialsRecalculateShipmentLineItemsHandler = RecalculateShipmentLineItemsHandler{context, shipmentlineitemservice.NewShipmentLineItemRecalculator(context.DB(), logger)}

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

	// Access Codes
	publicAPI.AccesscodeFetchAccessCodeHandler = FetchAccessCodeHandler{context, accesscodeservice.NewAccessCodeFetcher(context.DB())}
	publicAPI.AccesscodeValidateAccessCodeHandler = ValidateAccessCodeHandler{context, accesscodeservice.NewAccessCodeValidator(context.DB())}
	publicAPI.AccesscodeClaimAccessCodeHandler = ClaimAccessCodeHandler{context, accesscodeservice.NewAccessCodeClaimer(context.DB())}

	// Postal Codes
	publicAPI.PostalCodesValidatePostalCodeWithRateDataHandler = ValidatePostalCodeWithRateDataHandler{
		context,
		postalcodeservice.NewPostalCodeValidator(context.DB()),
	}

	return publicAPI.Serve(nil)
}
