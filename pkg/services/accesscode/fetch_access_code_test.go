package accesscode

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AccessCodeServiceSuite) TestFetchAccessCode_FetchAccessCode() {
	user := testdatagen.MakeDefaultServiceMember(suite.DB())
	selectedMoveType := models.SelectedMoveTypePPM

	code := "CODE12"
	serviceMemberID := &user.ID
	accessCode := models.AccessCode{
		Code:            code,
		MoveType:        &selectedMoveType,
		ServiceMemberID: serviceMemberID,
	}
	suite.MustSave(&accessCode)
	fetchAccessCode := NewAccessCodeFetcher(suite.DB())
	ac, _ := fetchAccessCode.FetchAccessCode(*serviceMemberID)

	suite.Equal(ac.Code, accessCode.Code, "expected CODE12")
}

func (suite *AccessCodeServiceSuite) TestFetchAccessCode_FetchEmptyAccessCode() {
	user := testdatagen.MakeDefaultServiceMember(suite.DB())
	serviceMemberID := &user.ID
	fetchAccessCode := NewAccessCodeFetcher(suite.DB())
	ac, _ := fetchAccessCode.FetchAccessCode(*serviceMemberID)
	suite.Equal(ac.Code, "")
}
