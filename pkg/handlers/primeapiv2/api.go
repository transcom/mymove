package primeapiv2

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/primev2api"
	"github.com/transcom/mymove/pkg/gen/primev2api/primev2operations"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorderfetcherv2 "github.com/transcom/mymove/pkg/services/move_task_order/move_task_order_fetcher/move_task_order_fetcher_v2"
)

// NewPrimeAPI returns the Prime API
func NewPrimeAPI(handlerConfig handlers.HandlerConfig) *primev2operations.MymoveAPI {
	// It is very important, if you copy the prime v1 api.go file, that you change every reference to primev2api
	// If you loads.Analyzed the wrong swagger file here, you will use the wrong definition to initialize the API.
	// Not that I did that and spent a few hours trying to figure out what was going on or anything... :face-palm:
	primeSpec, err := loads.Analyzed(primev2api.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPIV2 := primev2operations.NewMymoveAPI(primeSpec)

	primeAPIV2.ServeError = handlers.ServeCustomError

	primeAPIV2.MoveTaskOrderListMovesHandler = ListMovesHandler{
		handlerConfig,
		movetaskorderfetcherv2.NewMoveTaskOrderFetcher(),
	}

	return primeAPIV2
}
