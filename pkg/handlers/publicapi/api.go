package publicapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/handlers"
	accesscodeservice "github.com/transcom/mymove/pkg/services/accesscode"
)

// NewPublicAPIHandler returns a handler for the public API
func NewPublicAPIHandler(context handlers.HandlerContext) http.Handler {

	// Wire up the handlers to the publicAPIMux
	apiSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	publicAPI := publicops.NewMymoveAPI(apiSpec)

	// TSPs
	publicAPI.TspsIndexTSPsHandler = TspsIndexTSPsHandler{context}

	// Access Codes
	publicAPI.AccesscodeFetchAccessCodeHandler = FetchAccessCodeHandler{context, accesscodeservice.NewAccessCodeFetcher(context.DB())}
	publicAPI.AccesscodeValidateAccessCodeHandler = ValidateAccessCodeHandler{context, accesscodeservice.NewAccessCodeValidator(context.DB())}
	publicAPI.AccesscodeClaimAccessCodeHandler = ClaimAccessCodeHandler{context, accesscodeservice.NewAccessCodeClaimer(context.DB())}

	return publicAPI.Serve(nil)
}
