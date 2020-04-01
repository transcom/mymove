package supportapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/supportapi"
	supportops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewSupportAPIHandler returns a handler for the Prime API
func NewSupportAPIHandler(context handlers.HandlerContext) http.Handler {
	// builder := query.NewQueryBuilder(context.DB())
	// fetcher := fetch.NewFetcher(builder)

	supportSpec, err := loads.Analyzed(supportapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	supportAPI := supportops.NewMymoveAPI(supportSpec)
	// queryBuilder := query.NewQueryBuilder(context.DB())

	// supportAPI.MoveTaskOrderFetchMTOUpdatesHandler = FetchMTOUpdatesHandler{
	// 	context,
	// }

	return supportAPI.Serve(nil)
}
