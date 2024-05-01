package pptasapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/transcom/mymove/pkg/gen/pptasapi"
	pptasops "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

func NewPPTASApiHandler(handlerConfig handlers.HandlerConfig) http.Handler {

	pptasSpec, err := loads.Analyzed(pptasapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	pptasApi := pptasops.NewMymoveAPI(pptasSpec)
	pptasApi.ServeError = handlers.ServeCustomError

	pptasApi.MovesListMovesHandler = ListMovesHandler{
		HandlerConfig:        handlerConfig,
		MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(),
	}

	return pptasApi.Serve(nil)
}
