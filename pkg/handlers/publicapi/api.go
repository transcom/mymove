package publicapi

import (
	"github.com/transcom/mymove/pkg/paperwork"
	paperworkservice "github.com/transcom/mymove/pkg/services/paperwork"
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/handlers"
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
	publicAPI.ShipmentsCreateGovBillOfLadingHandler = CreateGovBillOfLadingHandler{context, paperworkservice.NewCreateForm(context.FileStorer().TempFileSystem(), paperwork.NewFormFiller())}

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

	return publicAPI.Serve(nil)
}
