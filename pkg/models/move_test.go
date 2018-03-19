package models_test

import (
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicMoveInstantiation() {
	move := &Move{}

	expErrors := map[string][]string{
		"selected_move_type": []string{"SelectedMoveType can not be blank."},
		"user_id":            []string{"UserID can not be blank."},
	}

	suite.verifyValidationErrors(move, expErrors)
}

func (suite *ModelSuite) TestGetMoveForUser() {
	t := suite.T()

	user1 := User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}

	user2 := User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "someoneelse@example.com",
	}

	verrs, err := suite.db.ValidateAndCreate(&user1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}
	verrs, err = suite.db.ValidateAndCreate(&user2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	move := Move{
		UserID:           user1.ID,
		SelectedMoveType: internalmessages.SelectedMoveTypeCOMBO,
	}
	verrs, err = suite.db.ValidateAndCreate(&move)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	fmt.Println(user1.ID, user2.ID, move.UserID)

	// All correct
	moveResult, err := GetMoveForUser(suite.db, user1.ID, move.ID)
	if err != nil {
		t.Error("Expected to get moveResult back.", err)
	}
	if !moveResult.IsValid() {
		t.Error("Expected the move to be valid")
	}
	if moveResult.Move().ID != move.ID {
		t.Error("Expected new move to match move.")
	}

	// Bad Move
	moveResult, err = GetMoveForUser(suite.db, user1.ID, uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Error("Expected to get a good moveResult back.", err)
	}
	if moveResult.IsValid() {
		t.Error("Expected the moveResult to be invalid")
	}
	if moveResult.ErrorCode() != FetchErrorNotFound {
		t.Error("Should have gotten a not found error")
	}

	// Bad User
	moveResult, err = GetMoveForUser(suite.db, user2.ID, move.ID)
	if err != nil {
		t.Error("Expected to get a good moveResult back.", err)
	}
	if moveResult.IsValid() {
		t.Error("Expected the moveResult to be invalid")
	}
	if moveResult.ErrorCode() != FetchErrorForbidden {
		t.Error("Should have gotten a forbidden error")
	}

}
