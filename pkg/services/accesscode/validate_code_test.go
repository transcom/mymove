package accesscode

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ValidateAccessCodeTestSuite struct {
	testingsuite.PopTestSuite
}

func (suite *ValidateAccessCodeTestSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestValidateAccessCodeTestSuite(t *testing.T) {
	ts := &ValidateAccessCodeTestSuite{
		testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, ts)
}

func (suite *ValidateAccessCodeTestSuite) TestValidateAccessCode_ValidAccessCode() {
	selectedMoveType := models.SelectedMoveTypePPM

	code := "CODE2"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: &selectedMoveType,
	}
	suite.MustSave(&accessCode)

	validateAccessCode := NewAccessCodeValidator(suite.DB())
	ac, valid, _ := validateAccessCode.ValidateAccessCode(code, selectedMoveType)
	suite.True(valid)
	suite.Equal(ac.Code, accessCode.Code, "expected CODE2")
}

func (suite *ValidateAccessCodeTestSuite) TestValidateAccessCode_InvalidAccessCode() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = models.OrderStatusSUBMITTED
	suite.MustSave(&orders)
	selectedMoveType := models.SelectedMoveTypePPM
	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	code := "CODE1"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: move.SelectedMoveType,
		MoveID:   move.ID,
	}
	suite.MustSave(&accessCode)
	validateAccessCode := NewAccessCodeValidator(suite.DB())
	_, valid, _ := validateAccessCode.ValidateAccessCode(code, selectedMoveType)
	suite.False(valid)
}
