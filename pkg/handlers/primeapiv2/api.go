package primeapiv2

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/primev2api"
	"github.com/transcom/mymove/pkg/gen/primev2api/primev2operations"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

// NewPrimeAPI returns the Prime API
func NewPrimeAPI(handlerConfig handlers.HandlerConfig) *primev2operations.MymoveAPI {
	primeSpec, err := loads.Analyzed(primev2api.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPIV2 := primev2operations.NewMymoveAPI(primeSpec)

	primeAPIV2.ServeError = handlers.ServeCustomError

	primeAPIV2.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	return primeAPIV2
}
