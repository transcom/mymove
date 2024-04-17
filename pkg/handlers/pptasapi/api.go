package pptasapi

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/transcom/mymove/pkg/gen/pptasapi"
	"github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations"
	pptasops "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

func NewPPTASApi(handlerConfig handlers.HandlerConfig) *pptasoperations.MymoveAPI {

	pptasSpec, err := loads.Analyzed(pptasapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	pptasApi := pptasops.NewMymoveAPI(pptasSpec)

	pptasApi.ServeError = handlers.ServeCustomError

	return pptasApi
}
