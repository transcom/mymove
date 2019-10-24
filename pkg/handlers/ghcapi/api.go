package ghcapi

import (
	"log"
	"net/http"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewGhcAPIHandler returns a handler for the GHC API
func NewGhcAPIHandler(context handlers.HandlerContext) http.Handler {

	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)

	ghcAPI.EntitlementsGetEntitlementsHandler = GetEntitlementsHandler{context}
	ghcAPI.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = UpdateMoveTaskOrderStatusHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB()),
	}

	return ghcAPI.Serve(nil)
}
