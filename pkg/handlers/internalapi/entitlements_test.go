package internalapi

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns200() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is under entitlement of 10500
	ppm.WeightEstimate = swag.Int64(10000)
	suite.MustSave(&ppm)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&entitlementop.ValidateEntitlementOK{}, response)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns409IfPPM() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is over entitlement of 10500
	ppm.WeightEstimate = swag.Int64(14000)
	suite.MustSave(&ppm)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 409 status code
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	// Then: expect a 409 status code
	suite.Assertions.Equal(http.StatusConflict, errResponse.Code)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns409IfHHG() {
	// Given: a set of orders, a move, user, servicemember and an HHG
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	move := shipment.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is over entitlement of 10500
	weight := unit.Pound(14000)
	shipment.WeightEstimate = &weight
	suite.MustSave(&shipment)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 409 status code
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	// Then: expect a 409 status code
	suite.Assertions.Equal(http.StatusConflict, errResponse.Code)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns409IfComboMove() {
	// Given: a set of orders, a move, user, servicemember and Combo HHG and PPM
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	move := shipment.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is over entitlement of 10500
	weight := unit.Pound(14000)
	shipment.WeightEstimate = &weight
	suite.MustSave(&shipment)

	ppmWeight := int64(weight)
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Move:           move,
			WeightEstimate: &ppmWeight,
		},
	})

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 409 status code
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	// Then: expect a 409 status code
	suite.Assertions.Equal(http.StatusConflict, errResponse.Code)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoPpmOrHhg() {
	// Given: a set of orders, a move, user, servicemember but NO ppm and NO hhg
	move := testdatagen.MakeDefaultMove(suite.DB())

	// When: rank is E1, the orders have dependents and spouse gear
	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.Assertions.IsType(&entitlementop.ValidateEntitlementNotFound{}, response)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoMoveOrOrders() {
	// Given: a user, servicemember but NO Move
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	// When: rank is E1, the orders have dependents and spouse gear
	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.AuthenticateRequest(request, serviceMember)

	badMoveID := uuid.Must(uuid.NewV4())

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(badMoveID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	// Then: expect a 404 status code
	suite.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoRank() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is under entitlement of 10500
	ppm.WeightEstimate = swag.Int64(10000)
	suite.MustSave(&ppm)

	move.Orders.ServiceMember.Rank = nil
	suite.MustSave(&move.Orders.ServiceMember)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.AuthenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.Assertions.IsType(&entitlementop.ValidateEntitlementNotFound{}, response)
}
