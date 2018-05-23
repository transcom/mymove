package handlers

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestApproveMoveHandler() {
	// Given: a set of orders, a move, user and servicemember
	move, _ := testdatagen.MakeMove(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.authenticateOfficeRequest(req, move.Orders.ServiceMember.User)

	params := officeop.ApproveMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: a move is approved
	handler := ApproveMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApproveMoveOK{}, response)
	okResponse := response.(*officeop.ApproveMoveOK)

	// And: Returned query to have an approved status
	suite.Assertions.Equal(internalmessages.MoveStatusAPPROVED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApprovePPMHandler() {
	// Given: a set of orders, a move, user and servicemember
	ppm, _ := testdatagen.MakePPM(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/moves/some_id/approve", nil)
	req = suite.authenticateOfficeRequest(req, ppm.Move.Orders.ServiceMember.User)

	params := officeop.ApprovePPMParams{
		HTTPRequest:              req,
		PersonallyProcuredMoveID: strfmt.UUID(ppm.ID.String()),
	}

	// And: a ppm is approved
	handler := ApprovePPMHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&officeop.ApprovePPMOK{}, response)
	okResponse := response.(*officeop.ApprovePPMOK)

	// And: Returned query to have an approved status
	suite.Assertions.Equal(internalmessages.PPMStatusAPPROVED, okResponse.Payload.Status)
}
