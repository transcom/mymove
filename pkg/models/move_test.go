package models_test

import (
	"fmt"

	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicMoveInstantiation() {
	move := &Move{}

	expErrors := map[string][]string{
		"user_id": {"UserID can not be blank."},
	}

	suite.verifyValidationErrors(move, expErrors)
}

func (suite *ModelSuite) TestFetchMove() {
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
	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move := Move{
		UserID:           user1.ID,
		SelectedMoveType: &selectedType,
	}
	verrs, err = suite.db.ValidateAndCreate(&move)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	fmt.Println(user1.ID, user2.ID, move.UserID)

	// All correct
	fetchedMove, err := FetchMove(suite.db, user1, move.ID)
	if err != nil {
		t.Error("Expected to get moveResult back.", err)
	}
	if fetchedMove.ID != move.ID {
		t.Error("Expected new move to match move.")
	}

	// Bad Move
	fetchedMove, err = FetchMove(suite.db, user1, uuid.Must(uuid.NewV4()))
	if err != ErrFetchNotFound {
		t.Error("Expected to get fetchnotfound.", err)
	}

	// Bad User
	fetchedMove, err = FetchMove(suite.db, user2, move.ID)
	if err != ErrFetchForbidden {
		t.Error("Expected to get a Forbidden back.", err)
	}

}
