package models_test

import (
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
