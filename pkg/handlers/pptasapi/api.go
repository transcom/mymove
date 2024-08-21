package pptasapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/pptasapi"
	pptasops "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

func NewPPTASAPI(handlerConfig handlers.HandlerConfig) *pptasops.MymoveAPI {
	pptasSpec, err := loads.Analyzed(pptasapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	pptasAPI := pptasops.NewMymoveAPI(pptasSpec)
	pptasAPI.ServeError = handlers.ServeCustomError

	pptasAPI.MovesListMovesHandler = ListMovesHandler{
		HandlerConfig:        handlerConfig,
		MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(),
	}

	return pptasAPI
}
