package accesscode

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AccessCodeServiceSuite) TestValidateAccessCode_ValidAccessCode() {
	selectedMoveType := models.SelectedMoveTypePPM

	code := "CODE12"
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

func (suite *AccessCodeServiceSuite) TestValidateAccessCode_InvalidAccessCode() {
	user := testdatagen.MakeDefaultServiceMember(suite.DB())
	selectedMoveType := models.SelectedMoveTypeHHG

	code := "CODE12"
	usedAccessCode := models.AccessCode{
		Code:            code,
		MoveType:        &selectedMoveType,
		ServiceMemberID: &user.ID,
	}
	suite.MustSave(&usedAccessCode)
	validateAccessCode := NewAccessCodeValidator(suite.DB())
	_, valid, _ := validateAccessCode.ValidateAccessCode(code, selectedMoveType)
	suite.False(valid)
}
