package models_test

import (
	"fmt"
	"strings"

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
	fmt.Println("Are we testing moves though?")

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

	theMove, err := GetMoveForUser(suite.db, user1.ID, move.ID)
	if err != nil {
		t.Error("Expected to get theMove back.", err)
	}
	if theMove.ID != move.ID {
		t.Error("Expected theMove to match move.")
	}

	_, err = GetMoveForUser(suite.db, user2.ID, move.ID)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "no rows in result set") {
			t.Error("Expected the error to end with 'no rows in result set'")
		}
	} else {
		t.Error("We should not have been able to retrieve this move")
	}

}
