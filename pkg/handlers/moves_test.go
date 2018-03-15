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
