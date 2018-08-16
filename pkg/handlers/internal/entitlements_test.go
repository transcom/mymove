package internal

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns200() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm := testdatagen.MakeDefaultPPM(suite.parent.Db)
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is under entitlement of 10500
	ppm.WeightEstimate = swag.Int64(10000)
	suite.parent.MustSave(&ppm)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.parent.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&entitlementop.ValidateEntitlementOK{}, response)

}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns409() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm := testdatagen.MakeDefaultPPM(suite.parent.Db)
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is over entitlement of 10500
	ppm.WeightEstimate = swag.Int64(14000)
	suite.parent.MustSave(&ppm)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.parent.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 409 status code
	suite.parent.Assertions.IsType(&utils.ErrResponse{}, response)
	errResponse := response.(*utils.ErrResponse)

	// Then: expect a 409 status code
	suite.parent.Assertions.Equal(http.StatusConflict, errResponse.Code)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoPpm() {
	// Given: a set of orders, a move, user, servicemember but NO ppm
	move := testdatagen.MakeDefaultMove(suite.parent.Db)

	// When: rank is E1, the orders have dependents and spouse gear
	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.parent.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.parent.Assertions.IsType(&entitlementop.ValidateEntitlementNotFound{}, response)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoMoveOrOrders() {
	// Given: a user, servicemember but NO Move
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.parent.Db)

	// When: rank is E1, the orders have dependents and spouse gear
	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.parent.AuthenticateRequest(request, serviceMember)

	badMoveID := uuid.Must(uuid.NewV4())

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(badMoveID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.parent.Assertions.IsType(&utils.ErrResponse{}, response)
	errResponse := response.(*utils.ErrResponse)

	// Then: expect a 404 status code
	suite.parent.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoRank() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm := testdatagen.MakeDefaultPPM(suite.parent.Db)
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is under entitlement of 10500
	ppm.WeightEstimate = swag.Int64(10000)
	suite.parent.MustSave(&ppm)

	move.Orders.ServiceMember.Rank = nil
	suite.parent.MustSave(&move.Orders.ServiceMember)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.parent.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.parent.Assertions.IsType(&entitlementop.ValidateEntitlementNotFound{}, response)
}
