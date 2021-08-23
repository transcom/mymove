package accesscode

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AccessCodeServiceSuite) TestFetchAccessCode_FetchAccessCode() {
	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	code := "CODE12"
	serviceMemberID := &user.ID
	accessCode := models.AccessCode{
		Code:            code,
		MoveType:        models.SelectedMoveTypePPM,
		ServiceMemberID: serviceMemberID,
	}
	suite.MustSave(&accessCode)
	fetchAccessCode := NewAccessCodeFetcher()
	ac, _ := fetchAccessCode.FetchAccessCode(suite.TestAppContext(), *serviceMemberID)

	suite.Equal(ac.Code, accessCode.Code, "expected CODE12")
}

func (suite *AccessCodeServiceSuite) TestFetchAccessCode_FetchNotFound() {
	user := testdatagen.MakeDefaultServiceMember(suite.DB())
	serviceMemberID := &user.ID
	fetchAccessCode := NewAccessCodeFetcher()
	_, err := fetchAccessCode.FetchAccessCode(suite.TestAppContext(), *serviceMemberID)
	suite.Error(err)
	suite.Equal(sql.ErrNoRows, err)
}
