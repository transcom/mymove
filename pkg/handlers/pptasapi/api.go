package pptasapi

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/transcom/mymove/pkg/gen/pptasapi"
	pptasops "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/move"
)

func NewPPTASApiHandler(handlerConfig handlers.HandlerConfig) *pptasops.MymoveAPI {

	pptasSpec, err := loads.Analyzed(pptasapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	pptasApi := pptasops.NewMymoveAPI(pptasSpec)

	pptasApi.MovesMovesSinceHandler = IndexMovesHandler{
		HandlerConfig: handlerConfig,
		MoveSearcher:  move.NewMoveSearcher(),
	}

	pptasApi.ServeError = handlers.ServeCustomError

	return pptasApi
}
