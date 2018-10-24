package ordersapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ordersapi"
	ordersops "github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// Handler is a DI marker for the Orders http.Handler
type Handler http.Handler

// NewOrdersAPIHandler returns a handler for the Orders API
func NewOrdersAPIHandler(context handlers.HandlerContext) Handler {

	// Wire up the handlers to the ordersAPIMux
	ordersSpec, err := loads.Analyzed(ordersapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	ordersAPI := ordersops.NewMymoveAPI(ordersSpec)
	ordersAPI.GetOrdersHandler = GetOrdersHandler{context}
	ordersAPI.IndexOrdersHandler = IndexOrdersHandler{context}
	ordersAPI.PostRevisionHandler = PostRevisionHandler{context}
	ordersAPI.PostRevisionToOrdersHandler = PostRevisionToOrdersHandler{context}
	return ordersAPI.Serve(nil)
}
