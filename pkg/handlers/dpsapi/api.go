package dpsapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/dpsapi"
	dpsops "github.com/transcom/mymove/pkg/gen/dpsapi/dpsoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewDPSAPI returns the DPS API
func NewDPSAPI(context handlers.HandlerConfig) *dpsops.MymoveAPI {
	dpsSpec, err := loads.Analyzed(dpsapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	dpsAPI := dpsops.NewMymoveAPI(dpsSpec)
	dpsAPI.DpsGetUserHandler = GetUserHandler{context}
	return dpsAPI
}
