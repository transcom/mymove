package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
)

// ValidateEntitlementHandler validates a weight estimate based on entitlement
type ValidateEntitlementHandler HandlerContext

// Handle is the handler
func (h ValidateEntitlementHandler) Handle(params entitlementop.ValidateEntitlementParams) middleware.Responder {

	response := responseForConflictErrors(h.logger, fmt.Errorf())
	fmt.Println(response)
	return response
}
