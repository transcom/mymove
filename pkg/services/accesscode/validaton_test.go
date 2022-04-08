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
		MoveType: selectedMoveType,
	}

	suite.MustSave(&accessCode)
	err := checkAccessCode().Validate(suite.AppContextForTest(), &accessCode)
	suite.Empty(err.Error())
}

func (suite *AccessCodeServiceSuite) TestValidateAccessCode_InvalidAccessCode() {
	user := testdatagen.MakeDefaultServiceMember(suite.DB())
	selectedMoveType := models.SelectedMoveTypeHHG

	code := "CODE12"
	usedAccessCode := models.AccessCode{
		Code:            code,
		MoveType:        selectedMoveType,
		ServiceMemberID: &user.ID,
	}

	suite.MustSave(&usedAccessCode)
	err := checkAccessCode().Validate(suite.AppContextForTest(), &usedAccessCode)
	suite.NotEmpty(err.Error())
}
