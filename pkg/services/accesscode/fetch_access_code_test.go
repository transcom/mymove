package accesscode

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type FetchAccessCodeTestSuite struct {
	testingsuite.PopTestSuite
}

func (suite *FetchAccessCodeTestSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestFetchAccessCodeTestSuite(t *testing.T) {
	ts := &FetchAccessCodeTestSuite{
		testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, ts)
}

func (suite *FetchAccessCodeTestSuite) TestFetchAccessCode_FetchAccessCode() {
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

func (suite *FetchAccessCodeTestSuite) TestFetchAccessCode_FetchEmptyAccessCode() {
	user := testdatagen.MakeDefaultServiceMember(suite.DB())
	serviceMemberID := &user.ID
	fetchAccessCode := NewAccessCodeFetcher(suite.DB())
	ac, _ := fetchAccessCode.FetchAccessCode(*serviceMemberID)
	suite.Equal(ac.Code, "")
}
