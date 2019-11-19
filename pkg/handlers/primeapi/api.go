package primeapi

import (
	"log"
	"net/http"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/primeapi"
	primeops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewPrimeAPIHandler returns a handler for the Prime API
func NewPrimeAPIHandler(context handlers.HandlerContext) http.Handler {

	primeSpec, err := loads.Analyzed(primeapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPI := primeops.NewMymoveAPI(primeSpec)

	primeAPI.MoveTaskOrderListMoveTaskOrdersHandler = ListMoveTaskOrdersHandler{
		context,
	}
	primeAPI.MoveTaskOrderUpdateMoveTaskOrderEstimatedWeightHandler = UpdateMoveTaskOrderEstimatedWeightHandler{
		context,
		movetaskorder.NewMoveTaskOrderEstimatedWeightUpdater(context.DB()),
	}
	primeAPI.MoveTaskOrderUpdateMoveTaskOrderPostCounselingInformationHandler = UpdateMoveTaskOrderPostCounselingInformationHandler{
		context,
		movetaskorder.NewMoveTaskOrderPostCounselingInformationUpdater(context.DB()),
	}
	primeAPI.MoveTaskOrderUpdateMoveTaskOrderDestinationAddressHandler = UpdateMoveTaskOrderDestinationAddressHandler{
		context,
		movetaskorder.NewMoveTaskOrderDestinationAddressUpdater(context.DB()),
	}
	primeAPI.MoveTaskOrderGetPrimeEntitlementsHandler = GetPrimeEntitlementsHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB()),
	}
	primeAPI.MoveTaskOrderUpdateMoveTaskOrderActualWeightHandler = UpdateMoveTaskOrderActualWeightHandler{
		context,
		movetaskorder.NewMoveTaskOrderActualWeightUpdater(context.DB()),
	}
	primeAPI.MoveTaskOrderGetMoveTaskOrderCustomerHandler = GetMoveTaskOrderCustomerHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB()),
	}
	primeAPI.MoveTaskOrderUpdateMoveTaskOrderPostCounselingInformationHandler = UpdateMoveTaskOrderPostCounselingInformationHandler{
		context,
		movetaskorder.NewMoveTaskOrderPostCounselingInformationUpdater(context.DB()),
	}
	primeAPI.MoveTaskOrderUpdateMoveTaskOrderDestinationAddressHandler = UpdateMoveTaskOrderDestinationAddressHandler{
		context,
		movetaskorder.NewMoveTaskOrderDestinationAddressUpdater(context.DB()),
	}

	return primeAPI.Serve(nil)
}
