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

	return publicAPI.Serve(nil)
}
