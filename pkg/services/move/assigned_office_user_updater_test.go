package move

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *MoveServiceSuite) TestUpdateAssignedOfficeUser() {

	assignedOfficeUserUpdater := NewAssignedOfficeUserUpdater(NewMoveFetcher())
	scUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		ActiveRole:      scUser.User.Roles[0],
		OfficeUserID:    scUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	appCtx := suite.AppContextWithSessionForTest(&session)

	createdMove := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
	createdMove.SCCounselingAssignedID = &scUser.ID
	createdMove.SCCounselingAssignedUser = &scUser
	move, updateError := assignedOfficeUserUpdater.UpdateAssignedOfficeUser(appCtx, createdMove.ID, &scUser, models.QueueTypeCounseling)

	suite.NotNil(move)
	suite.FatalNoError(updateError)

	suite.Equal(createdMove.SCCounselingAssignedID, move.SCCounselingAssignedID)
	suite.Equal(createdMove.SCCounselingAssignedUser.ID, move.SCCounselingAssignedUser.ID)
	suite.Equal(createdMove.SCCounselingAssignedUser.FirstName, move.SCCounselingAssignedUser.FirstName)
	suite.Equal(createdMove.SCCounselingAssignedUser.LastName, move.SCCounselingAssignedUser.LastName)
}
