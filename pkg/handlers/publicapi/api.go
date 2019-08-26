package publicapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/handlers"
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

	// TSPs
	publicAPI.TspsIndexTSPsHandler = TspsIndexTSPsHandler{context}

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

	return publicAPI.Serve(nil)
}
