package accesscode

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ClaimAccessCodeTestSuite struct {
	testingsuite.PopTestSuite
}

func (suite *ClaimAccessCodeTestSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestClaimAccessCodeTestSuite(t *testing.T) {
	ts := &ClaimAccessCodeTestSuite{
		testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, ts)
}

func (suite *ClaimAccessCodeTestSuite) TestClaimAccessCode_Success() {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	selectedMoveType := models.SelectedMoveTypePPM

	code := "CODE12"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: &selectedMoveType,
	}

	suite.MustSave(&accessCode)
	claimAccessCode := NewAccessCodeClaimer(suite.DB())
	ac, err := claimAccessCode.ClaimAccessCode(code, serviceMember.ID)

	suite.Nil(err)
	suite.Equal(ac.Code, accessCode.Code, "expected CODE2")
	suite.Equal(ac.ServiceMemberID, &serviceMember.ID)
	suite.NotNil(ac.ClaimedAt)
}

func (suite *ClaimAccessCodeTestSuite) TestClaimAccessCode_Failed() {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	selectedMoveType := models.SelectedMoveTypePPM

	code := "CODE12"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: &selectedMoveType,
	}

	suite.MustSave(&accessCode)
	claimAccessCode := NewAccessCodeClaimer(suite.DB())
	ac1, err1 := claimAccessCode.ClaimAccessCode(code, serviceMember.ID)

	suite.Nil(err1)
	suite.Equal(ac1.Code, accessCode.Code, "expected CODE2")
	suite.Equal(ac1.ServiceMemberID, &serviceMember.ID)
	suite.NotNil(ac1.ClaimedAt)

	_, err2 := claimAccessCode.ClaimAccessCode(code, serviceMember.ID)

	suite.Equal(err2.Error(), "Unable to claim access code: Access code already claimed")
}

func (suite *ClaimAccessCodeTestSuite) TestClaimAccessCode_InvalidAccessCode() {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	code := "CODE12"

	claimAccessCode := NewAccessCodeClaimer(suite.DB())
	_, err := claimAccessCode.ClaimAccessCode(code, serviceMember.ID)

	suite.Equal(err.Error(), "Unable to claim access code: Unable to find access code: sql: no rows in result set")
}
