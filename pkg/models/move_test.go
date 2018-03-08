package models_test

import (
	"github.com/satori/go.uuid"

	. "github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func (suite *ModelSuite) TestGetOrCreateMove() {
	t := suite.T()

	// When: user ID is passed to create move func
	loginGovUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	user := User{
		LoginGovUUID:  loginGovUUID,
		LoginGovEmail: "sally@government.gov",
	}
	verrs, err := suite.db.ValidateAndCreate(&user)
	if verrs.HasAny() {
		t.Error("Error validating:", verrs)
	} else if err != nil {
		t.Error(err)
	}

	// And: user does not yet exist in the db
	newMove, err := GetOrCreateMove(suite.db, user.ID)
	if err != nil {
		t.Error("error querying or creating move.", err)
	}

	// Then: expect fields to be set on returned user
	if newMove.UserID != user.ID {
		t.Error("expected uuid to be set")
	}

	// When: The same UUID is passed in func
	sameUser, err := GetOrCreateMove(suite.db, user.ID)
	if err != nil {
		t.Error("error querying or creating move.")
	}

	// Then: expect the existing move to be returned
	if sameUser.UserID != newMove.UserID {
		t.Error("expected existing move to have been returned")
	}

	// And: no new move to have been created
	query := suite.db.Where("user_id = $1", user.ID)
	var moves []Move
	queryErr := query.All(&moves)
	if queryErr != nil {
		t.Error("DB Query Error", zap.Error(err))
	}
	if len(moves) > 1 {
		t.Error("1 move should have been returned")
	}
}
