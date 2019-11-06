package ghcapi

import (
	"log"
	"net/http"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/query"
	serviceitem "github.com/transcom/mymove/pkg/services/service_item"
)

// NewGhcAPIHandler returns a handler for the GHC API
func NewGhcAPIHandler(context handlers.HandlerContext) http.Handler {

	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)
	queryBuilder := query.NewQueryBuilder(context.DB())

	ghcAPI.EntitlementsGetEntitlementsHandler = GetEntitlementsHandler{context}
	ghcAPI.CustomerGetCustomerInfoHandler = GetCustomerInfoHandler{context}
	ghcAPI.MoveTaskOrderUpdateMoveTaskOrderActualWeightHandler = UpdateMoveTaskOrderActualWeightHandler{
		context,
		movetaskorder.NewMoveTaskOrderActualWeightUpdater(context.DB())}
	ghcAPI.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = UpdateMoveTaskOrderStatusHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderStatusUpdater(context.DB()),
	}
	ghcAPI.ServiceItemListServiceItemsHandler = ListServiceItemsHandler{
		context,
		serviceitem.NewServiceItemListFetcher(queryBuilder),
		query.NewQueryFilter,
	}
	ghcAPI.ServiceItemCreateServiceItemHandler = CreateServiceItemHandler{
		context,
		serviceitem.NewServiceItemCreator(queryBuilder),
		query.NewQueryFilter,
	}

	ghcAPI.CustomerGetAllCustomerMovesHandler = GetAllCustomerMovesHandler{context}
	return ghcAPI.Serve(nil)
}
