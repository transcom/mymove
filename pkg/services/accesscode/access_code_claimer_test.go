package accesscode

import (
	"github.com/transcom/mymove/pkg/appconfig"
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
	appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)

	claimAccessCode := NewAccessCodeClaimer()
	ac, _, err := claimAccessCode.ClaimAccessCode(appCfg, code, serviceMember.ID)

	suite.NoError(err)
	suite.Equal(ac.Code, accessCode.Code, "expected CODE2")
	suite.Equal(ac.ServiceMemberID, &serviceMember.ID)
	suite.NotNil(ac.ClaimedAt)
}

func (suite *AccessCodeServiceSuite) TestClaimAccessCode_Failed() {
	appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	code := "CODE12"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: models.SelectedMoveTypePPM,
	}

	suite.MustSave(&accessCode)
	claimAccessCode := NewAccessCodeClaimer()
	ac1, _, err1 := claimAccessCode.ClaimAccessCode(appCfg, code, serviceMember.ID)

	suite.Nil(err1)
	suite.Equal(ac1.Code, accessCode.Code, "expected CODE2")
	suite.Equal(ac1.ServiceMemberID, &serviceMember.ID)
	suite.NotNil(ac1.ClaimedAt)
	_, _, err2 := claimAccessCode.ClaimAccessCode(appCfg, code, serviceMember.ID)

	suite.Equal(err2.Error(), "Access code already claimed")
}
func (suite *AccessCodeServiceSuite) TestClaimAccessCode_InvalidAccessCode() {
	appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	code := "CODE12"

	claimAccessCode := NewAccessCodeClaimer()
	_, _, err := claimAccessCode.ClaimAccessCode(appCfg, code, serviceMember.ID)

	suite.Equal(err.Error(), "Unable to find access code: "+models.RecordNotFoundErrorString)
}
