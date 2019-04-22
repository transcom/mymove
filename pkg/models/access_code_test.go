package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestValidateInvalidAccessCode() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED
	suite.MustSave(&orders)
	selectedMoveType := SelectedMoveTypePPM
	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	code := "CODE1"
	accessCode := AccessCode{
		Code:     code,
		MoveType: *move.SelectedMoveType,
		MoveID:   move.ID,
	}
	suite.MustSave(&accessCode)
	_, valid := ValidateAccessCode(suite.DB(), code, *move.SelectedMoveType)
	suite.False(valid)
}

func (suite *ModelSuite) TestValidateValidAccessCode() {
	selectedMoveType := SelectedMoveTypePPM

	code := "CODE2"
	accessCode := AccessCode{
		Code:     code,
		MoveType: selectedMoveType,
	}
	suite.MustSave(&accessCode)
	ac, valid := ValidateAccessCode(suite.DB(), code, selectedMoveType)
	suite.True(valid)
	suite.Equal(ac.Code, accessCode.Code, "expected CODE2")
}
