package handlers

import (
	"fmt"
	"net/http/httptest"

	"github.com/satori/go.uuid"

	"github.com/transcom/mymove/pkg/auth/context"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestSubmitMoveHandlerAllValues() {
	t := suite.T()

	// Given: a logged in user
	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	// When: a new Move is posted
	var selectedType = internalmessages.SelectedMoveTypeHHG
	newMovePayload := internalmessages.CreateMovePayload{
		SelectedMoveType: selectedType,
	}
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
	query := suite.db.Where(fmt.Sprintf("user_id='%v'", user.ID)).Where(fmt.Sprintf("selected_move_type='%v'", selectedType))
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
		SelectedMoveType: selectedType,
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
	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: "HHG",
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
	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: "HHG",
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
	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	userUUID2, _ := uuid.FromString("3511d4d6-019d-4031-9c27-8a553e055543")
	user2 := models.User{
		LoginGovUUID:  userUUID2,
		LoginGovEmail: "email2@example.com",
	}
	suite.mustSave(&user2)

	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: "HHG",
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
