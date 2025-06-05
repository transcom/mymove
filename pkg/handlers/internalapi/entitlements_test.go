package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/services/entitlements"
)

func (suite *HandlerSuite) TestIndexEntitlementsHandlerReturns200() {
	// Given: a set of orders, a move, user, servicemember and a PPM

	waf := entitlements.NewWeightAllotmentFetcher()
	ppm := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
	move := factory.BuildMove(suite.DB(), nil, nil)
	mtoShipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
	mtoShipment.PPMShipment = &ppm
	move.MTOShipments = append(move.MTOShipments, mtoShipment)
	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements", nil)
	request = suite.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.IndexEntitlementsParams{
		HTTPRequest: request,
	}

	// And: index entitlements endpoint is hit
	handler := IndexEntitlementsHandler{suite.NewHandlerConfig(), waf}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&entitlementop.IndexEntitlementsOK{}, response)
}
