package accesscode

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AccessCodeServiceSuite) TestClaimAccessCode_Success() {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	code := "CODE12"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: models.SelectedMoveTypePPM,
	}

	suite.MustSave(&accessCode)
	claimAccessCode := NewAccessCodeClaimer(suite.DB())
	ac, _, err := claimAccessCode.ClaimAccessCode(code, serviceMember.ID)

	suite.NoError(err)
	suite.Equal(ac.Code, accessCode.Code, "expected CODE2")
	suite.Equal(ac.ServiceMemberID, &serviceMember.ID)
	suite.NotNil(ac.ClaimedAt)
}

func (suite *AccessCodeServiceSuite) TestClaimAccessCode_Failed() {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	code := "CODE12"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: models.SelectedMoveTypePPM,
	}

	suite.MustSave(&accessCode)
	claimAccessCode := NewAccessCodeClaimer(suite.DB())
	ac1, _, err1 := claimAccessCode.ClaimAccessCode(code, serviceMember.ID)

	suite.Nil(err1)
	suite.Equal(ac1.Code, accessCode.Code, "expected CODE2")
	suite.Equal(ac1.ServiceMemberID, &serviceMember.ID)
	suite.NotNil(ac1.ClaimedAt)
	_, _, err2 := claimAccessCode.ClaimAccessCode(code, serviceMember.ID)

	suite.Equal(err2.Error(), "Access code already claimed")
}
func (suite *AccessCodeServiceSuite) TestClaimAccessCode_InvalidAccessCode() {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	code := "CODE12"

	claimAccessCode := NewAccessCodeClaimer(suite.DB())
	_, _, err := claimAccessCode.ClaimAccessCode(code, serviceMember.ID)

	suite.Equal(err.Error(), "Unable to find access code: "+models.RecordNotFoundErrorString)
}
