package internalapi

import (
	"net/http/httptest"

	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexEntitlementsHandlerReturns200() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements", nil)
	request = suite.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.IndexEntitlementsParams{
		HTTPRequest: request,
	}

	// And: index entitlements endpoint is hit
	handler := IndexEntitlementsHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&entitlementop.IndexEntitlementsOK{}, response)
}
