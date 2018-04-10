package handlers

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth/context"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestSubmitMoveHandlerAllValues() {
	t := suite.T()

	// Given: a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	// When: a new Move is posted
	newMovePayload := internalmessages.CreateMovePayload{}
	req := httptest.NewRequest("GET", "/moves", nil)

	params := moveop.CreateMoveParams{
		CreateMovePayload: &newMovePayload,
		HTTPRequest:       req,
	}

	// And: the context contains the auth values
	ctx := params.HTTPRequest.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	params.HTTPRequest = params.HTTPRequest.WithContext(ctx)

	handler := CreateMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*moveop.CreateMoveCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	// Then: we expect a move to have been created for the user
	query := suite.db.Where(fmt.Sprintf("user_id='%v'", user.ID))
	moves := []models.Move{}
	query.All(&moves)

	if len(moves) != 1 {
		t.Errorf("Expected to find 1 move but found %v", len(moves))
	}

}

func (suite *HandlerSuite) TestCreateMoveHandlerNoUserID() {
	t := suite.T()
	// Given: no authentication values in context
	// When: a new Move is posted
	var selectedType = internalmessages.SelectedMoveTypeHHG
	movePayload := internalmessages.CreateMovePayload{
		SelectedMoveType: &selectedType,
	}
	req := httptest.NewRequest("GET", "/moves", nil)
	params := moveop.CreateMoveParams{
		CreateMovePayload: &movePayload,
		HTTPRequest:       req,
	}

	handler := CreateMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*moveop.CreateMoveUnauthorized)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
	// Then: we expect no moves to have been created
	moves := []models.Move{}
	suite.db.All(&moves)

	if len(moves) > 0 {
		t.Errorf("Expected to find no moves but found %v", len(moves))
	}
}

func (suite *HandlerSuite) TestIndexMovesHandler() {
	t := suite.T()

	// Given: A move and a user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	var selectedType = internalmessages.SelectedMoveTypeHHG
	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move)

	req := httptest.NewRequest("GET", "/moves", nil)

	indexMovesParams := moveop.NewIndexMovesParams()

	// And: the context contains the auth values
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	indexMovesParams.HTTPRequest = req.WithContext(ctx)

	// And: All moves are queried
	indexHandler := IndexMovesHandler(NewHandlerContext(suite.db, suite.logger))
	indexResponse := indexHandler.Handle(indexMovesParams)

	// Then: Expect a 200 status code
	okResponse := indexResponse.(*moveop.IndexMovesOK)
	moves := okResponse.Payload

	// And: Returned query to include our added move
	moveExists := false
	for _, move := range moves {
		if move.UserID.String() == user.ID.String() {
			moveExists = true
			break
		}
	}

	if !moveExists {
		t.Errorf("Expected an move to have user ID '%v'. None do.", user.ID)
	}
}

func (suite *HandlerSuite) TestIndexMovesHandlerNoUser() {
	t := suite.T()

	// Given: A move with a user that isn't logged in
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	var selectedType = internalmessages.SelectedMoveTypeHHG
	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move)

	req := httptest.NewRequest("GET", "/moves", nil)
	indexMovesParams := moveop.NewIndexMovesParams()
	indexMovesParams.HTTPRequest = req

	// And: All moves are queried
	indexHandler := IndexMovesHandler(NewHandlerContext(suite.db, suite.logger))
	indexResponse := indexHandler.Handle(indexMovesParams)

	// Then: Expect a 401 unauthorized
	_, ok := indexResponse.(*moveop.IndexMovesUnauthorized)
	if !ok {
		t.Errorf("Expected to get an unauthorized response, but got something else.")
	}
}

func (suite *HandlerSuite) TestIndexMovesWrongUser() {
	t := suite.T()

	// Given: A move with a user and a separate logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	user2 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email2@example.com",
	}
	suite.mustSave(&user2)

	var selectedType = internalmessages.SelectedMoveTypeHHG
	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move)

	req := httptest.NewRequest("GET", "/moves", nil)
	indexMovesParams := moveop.NewIndexMovesParams()

	// And: the context contains the auth values for user 2
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user2.ID, "fake token")
	indexMovesParams.HTTPRequest = req.WithContext(ctx)

	// And: All moves are queried
	indexHandler := IndexMovesHandler(NewHandlerContext(suite.db, suite.logger))
	indexResponse := indexHandler.Handle(indexMovesParams)

	// Then: Expect a 200 status code
	okResponse := indexResponse.(*moveop.IndexMovesOK)
	moves := okResponse.Payload

	// And: No moves should be returned
	if len(moves) != 0 {
		t.Errorf("Expected no moves to be found, but found %v", len(moves))
	}
}

func (suite *HandlerSuite) TestPatchMoveHandler() {
	t := suite.T()

	// Given: a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	var origType = internalmessages.SelectedMoveTypeHHG
	var newType = internalmessages.SelectedMoveTypeCOMBO
	newMove := models.Move{
		UserID:           user.ID,
		SelectedMoveType: &origType,
	}
	suite.mustSave(&newMove)

	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	req = req.WithContext(ctx)

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(newMove.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	okResponse, ok := response.(*moveop.PatchMoveCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	patchPPMPayload := okResponse.Payload

	if *patchPPMPayload.SelectedMoveType != newType {
		t.Fatalf("SelectedMoveType should have been updated.")
	}
}

func (suite *HandlerSuite) TestPatchMoveHandlerWrongUser() {
	t := suite.T()

	// Given: a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	user2 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email2@example.com",
	}
	suite.mustSave(&user2)

	var origType = internalmessages.SelectedMoveTypeHHG
	var newType = internalmessages.SelectedMoveTypeCOMBO
	newMove := models.Move{
		UserID:           user.ID,
		SelectedMoveType: &origType,
	}
	suite.mustSave(&newMove)

	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user2.ID, "fake token")
	req = req.WithContext(ctx)

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(newMove.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*moveop.PatchMoveForbidden)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoMove() {
	t := suite.T()

	// Given: a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	moveUUID := uuid.Must(uuid.NewV4())

	var newType = internalmessages.SelectedMoveTypeCOMBO

	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	req = req.WithContext(ctx)

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(moveUUID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*moveop.PatchMoveNotFound)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoType() {
	t := suite.T()

	// Given: a logged in user with a move
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	var origType = internalmessages.SelectedMoveTypeHHG
	newMove := models.Move{
		UserID:           user.ID,
		SelectedMoveType: &origType,
	}
	suite.mustSave(&newMove)

	patchPayload := internalmessages.PatchMovePayload{}

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	req = req.WithContext(ctx)

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(newMove.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*moveop.PatchMoveCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
}

func (suite *HandlerSuite) TestShowMoveHandler() {
	t := suite.T()

	// Given: A move and a user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	newMove := models.Move{
		UserID: user.ID,
	}
	suite.mustSave(&newMove)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	req = req.WithContext(ctx)

	params := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(newMove.ID.String()),
	}
	// And: show Move is queried
	showHandler := ShowMoveHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*moveop.ShowMoveOK)
	move := okResponse.Payload

	// And: Returned query to include our added move
	if move.UserID.String() != user.ID.String() {
		t.Errorf("Expected an move to have user ID '%v'. None do.", user.ID)
	}

}

func (suite *HandlerSuite) TestShowMoveHandlerNoUser() {
	t := suite.T()

	// Given: A move with a user that isn't logged in
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	move := models.Move{
		UserID: user.ID,
	}
	suite.mustSave(&move)

	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	showMoveParams := moveop.NewShowMoveParams()
	showMoveParams.HTTPRequest = req

	// And: Show move is queried
	showHandler := ShowMoveHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(showMoveParams)

	// Then: Expect a 401 unauthorized
	_, ok := showResponse.(*moveop.ShowMoveUnauthorized)
	if !ok {
		t.Errorf("Expected to get an unauthorized response, but got something else.")
	}
}

func (suite *HandlerSuite) TestShowMoveWrongUser() {
	t := suite.T()

	// Given: A move with a not-logged-in user and a separate logged-in user
	notLoggedInUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&notLoggedInUser)

	loggedInUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email2@example.com",
	}
	suite.mustSave(&loggedInUser)

	// When: A move is created for not-logged-in-user
	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	newMove := models.Move{
		UserID:           notLoggedInUser.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&newMove)

	// And: the context contains the auth values for logged-in user
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, loggedInUser.ID, "fake token")
	req = req.WithContext(ctx)
	showMoveParams := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(newMove.ID.String()),
	}
	// And: Show move is queried
	showHandler := ShowMoveHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(showMoveParams)

	_, ok := showResponse.(*moveop.ShowMoveForbidden)
	if !ok {
		t.Fatalf("Request failed: %#v", showResponse)
	}
}
