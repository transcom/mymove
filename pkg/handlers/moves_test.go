package handlers

import (
	// "fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	// "github.com/transcom/mymove/pkg/auth"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	// "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestSubmitMoveHandlerAllValues() {
	// Given: a set of orders, user and servicemember
	orders, _ := testdatagen.MakeOrder(suite.db)

	req := httptest.NewRequest("POST", "/orders/orderid/moves", nil)
	req = suite.authenticateRequest(req, orders.ServiceMember.User)

	// When: a new Move is posted
	var selectedType = internalmessages.SelectedMoveTypePPM
	newMovePayload := &internalmessages.CreateMovePayload{
		SelectedMoveType: &selectedType,
	}
	params := moveop.CreateMoveParams{
		OrdersID:          strfmt.UUID(orders.ID.String()),
		CreateMovePayload: newMovePayload,
		HTTPRequest:       req,
	}
	// Then: we expect a move to have been created based on orders
	handler := CreateMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.CreateMoveCreated{}, response)
	okResponse := response.(*moveop.CreateMoveCreated)

	suite.Assertions.Equal(orders.ID.String(), okResponse.Payload.OrdersID.String())
}

// func (suite *HandlerSuite) TestIndexMovesHandler() {
// 	t := suite.T()

// 	// Given: A move and a user
// 	user := models.User{
// 		LoginGovUUID:  uuid.Must(uuid.NewV4()),
// 		LoginGovEmail: "email@example.com",
// 	}
// 	suite.mustSave(&user)

// 	var selectedType = internalmessages.SelectedMoveTypeHHG
// 	move := models.Move{
// 		UserID:           user.ID,
// 		SelectedMoveType: &selectedType,
// 	}
// 	suite.mustSave(&move)

// 	req := httptest.NewRequest("GET", "/moves", nil)

// 	indexMovesParams := moveop.NewIndexMovesParams()

// 	// And: the context contains the auth values
// 	ctx := req.Context()
// 	ctx = auth.PopulateAuthContext(ctx, user.ID, "fake token")
// 	ctx = auth.PopulateUserModel(ctx, user)
// 	indexMovesParams.HTTPRequest = req.WithContext(ctx)

// 	// And: All moves are queried
// 	indexHandler := IndexMovesHandler(NewHandlerContext(suite.db, suite.logger))
// 	indexResponse := indexHandler.Handle(indexMovesParams)

// 	// Then: Expect a 200 status code
// 	okResponse := indexResponse.(*moveop.IndexMovesOK)
// 	moves := okResponse.Payload

// 	// And: Returned query to include our added move
// 	moveExists := false
// 	for _, move := range moves {
// 		if move.UserID.String() == user.ID.String() {
// 			moveExists = true
// 			break
// 		}
// 	}

// 	if !moveExists {
// 		t.Errorf("Expected a move to have user ID '%v'. None do.", user.ID)
// 	}
// }

// func (suite *HandlerSuite) TestIndexMovesWrongUser() {
// 	t := suite.T()

// 	// Given: A move with a user and a separate logged in user
// 	user := models.User{
// 		LoginGovUUID:  uuid.Must(uuid.NewV4()),
// 		LoginGovEmail: "email@example.com",
// 	}
// 	suite.mustSave(&user)

// 	user2 := models.User{
// 		LoginGovUUID:  uuid.Must(uuid.NewV4()),
// 		LoginGovEmail: "email2@example.com",
// 	}
// 	suite.mustSave(&user2)

// 	var selectedType = internalmessages.SelectedMoveTypeHHG
// 	move := models.Move{
// 		UserID:           user.ID,
// 		SelectedMoveType: &selectedType,
// 	}
// 	suite.mustSave(&move)

// 	req := httptest.NewRequest("GET", "/moves", nil)
// 	indexMovesParams := moveop.NewIndexMovesParams()

// 	// And: the context contains the auth values for user 2
// 	ctx := req.Context()
// 	ctx = auth.PopulateAuthContext(ctx, user2.ID, "fake token")
// 	ctx = auth.PopulateUserModel(ctx, user2)
// 	indexMovesParams.HTTPRequest = req.WithContext(ctx)

// 	// And: All moves are queried
// 	indexHandler := IndexMovesHandler(NewHandlerContext(suite.db, suite.logger))
// 	indexResponse := indexHandler.Handle(indexMovesParams)

// 	// Then: Expect a 200 status code
// 	okResponse := indexResponse.(*moveop.IndexMovesOK)
// 	moves := okResponse.Payload

// 	// And: No moves should be returned
// 	if len(moves) != 0 {
// 		t.Errorf("Expected no moves to be found, but found %v", len(moves))
// 	}
// }

func (suite *HandlerSuite) TestPatchMoveHandler() {
	// Given: a set of orders, a move, user and servicemember
	move, _ := testdatagen.MakeMove(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember.User)

	var newType = internalmessages.SelectedMoveTypeCOMBO
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}
	// And: a move is patched
	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&moveop.PatchMoveCreated{}, response)
	okResponse := response.(*moveop.PatchMoveCreated)

	// And: Returned query to include our added move
	suite.Assertions.Equal(&newType, okResponse.Payload.SelectedMoveType)
}

func (suite *HandlerSuite) TestPatchMoveHandlerWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move, _ := testdatagen.MakeMove(suite.db)
	// And: a not logged in user
	unAuthedUser, _ := testdatagen.MakeUser(suite.db)

	// And: the context contains a different user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, unAuthedUser)

	var newType = internalmessages.SelectedMoveTypeCOMBO
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.checkResponseForbidden(response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoMove() {
	// Given: a logged in user and no Move
	user, _ := testdatagen.MakeUser(suite.db)

	moveUUID := uuid.Must(uuid.NewV4())

	// And: the context contains a logged in user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, user)

	var newType = internalmessages.SelectedMoveTypeCOMBO
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(moveUUID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.checkResponseNotFound(response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoType() {
	// Given: a set of orders, a move, user and servicemember
	move, _ := testdatagen.MakeMove(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember.User)

	patchPayload := internalmessages.PatchMovePayload{}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.PatchMoveCreated{}, response)
	okResponse := response.(*moveop.PatchMoveCreated)

	suite.Assertions.Equal(move.ID.String(), okResponse.Payload.ID.String())
}

func (suite *HandlerSuite) TestShowMoveHandler() {

	// Given: a set of orders, a move, user and servicemember
	move, _ := testdatagen.MakeMove(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember.User)

	params := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: show Move is queried
	showHandler := ShowMoveHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.Assertions.IsType(&moveop.ShowMoveOK{}, showResponse)
	okResponse := showResponse.(*moveop.ShowMoveOK)

	// And: Returned query to include our added move
	suite.Assertions.Equal(move.OrdersID.String(), okResponse.Payload.OrdersID.String())

}

func (suite *HandlerSuite) TestShowMoveWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move, _ := testdatagen.MakeMove(suite.db)
	// And: a not logged in user
	unAuthedUser, _ := testdatagen.MakeUser(suite.db)

	// And: the context contains the auth values for not logged-in user
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.authenticateRequest(req, unAuthedUser)

	showMoveParams := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: Show move is queried
	showHandler := ShowMoveHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(showMoveParams)
	// Then: expect a forbidden response
	suite.checkResponseForbidden(showResponse)

}
