package publicapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/tsp"

	accesscodeservice "github.com/transcom/mymove/pkg/services/accesscode"
)

// NewPublicAPIHandler returns a handler for the public API
func NewPublicAPIHandler(context handlers.HandlerContext, logger Logger) http.Handler {

	// Wire up the handlers to the publicAPIMux
	apiSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	publicAPI := publicops.NewMymoveAPI(apiSpec)

	// Documents
	publicAPI.MoveDocsCreateGenericMoveDocumentHandler = CreateGenericMoveDocumentHandler{context}
	publicAPI.MoveDocsIndexMoveDocumentsHandler = IndexMoveDocumentsHandler{context}
	publicAPI.MoveDocsUpdateMoveDocumentHandler = UpdateMoveDocumentHandler{context}
	publicAPI.UploadsCreateUploadHandler = CreateUploadHandler{context}
	publicAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{context}

	publicAPI.AccessorialsGetTariff400ngItemsHandler = GetTariff400ngItemsHandler{context}

	// Service Agents
	publicAPI.ServiceAgentsIndexServiceAgentsHandler = IndexServiceAgentsHandler{context}
	publicAPI.ServiceAgentsCreateServiceAgentHandler = CreateServiceAgentHandler{context}
	publicAPI.ServiceAgentsPatchServiceAgentHandler = PatchServiceAgentHandler{context}

	// TSPs
	publicAPI.TransportationServiceProviderGetTransportationServiceProviderHandler = GetTransportationServiceProviderHandler{context}
	publicAPI.TspsIndexTSPsHandler = TspsIndexTSPsHandler{context}

	// Transportation Service Provider Performances
	queryBuilder := query.NewQueryBuilder(context.DB())
	publicAPI.TransportationServiceProviderPerformanceLogTransportationServiceProviderPerformanceHandler = LogTransportationServiceProviderPerformanceHandler{
		HandlerContext: context,
		NewQueryFilter: query.NewQueryFilter,
		TransportationServiceProviderPerformanceFetcher: tsp.NewTransportationServiceProviderPerformanceFetcher(queryBuilder),
	}

	// Access Codes
	publicAPI.AccesscodeFetchAccessCodeHandler = FetchAccessCodeHandler{context, accesscodeservice.NewAccessCodeFetcher(context.DB())}
	publicAPI.AccesscodeValidateAccessCodeHandler = ValidateAccessCodeHandler{context, accesscodeservice.NewAccessCodeValidator(context.DB())}
	publicAPI.AccesscodeClaimAccessCodeHandler = ClaimAccessCodeHandler{context, accesscodeservice.NewAccessCodeClaimer(context.DB())}

	return publicAPI.Serve(nil)
}
