package lockmove

import (
	"time"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *MoveLockerServiceSuite) TestLockMove() {
	moveLocker := NewMoveLocker()

	suite.Run("successfully returns move with office user values and lockExpiresAt value", func() {
		tooUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           tooUser.User.Roles,
			OfficeUserID:    tooUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		})

		move := factory.BuildMove(suite.DB(), nil, nil)

		actualMove, err := moveLocker.LockMove(appCtx, &move, tooUser.ID)
		suite.FatalNoError(err)

		// saving time and rounding time values to nearest minute to avoid nanosecond differences when testing
		now := time.Now()
		expirationTime := now.Add(30 * time.Minute).Truncate(time.Minute)

		suite.Equal(move.ID, actualMove.ID)
		suite.Equal(move.LockedByOfficeUserID, &tooUser.ID)
		suite.Equal(move.LockedByOfficeUser.TransportationOffice.Name, tooUser.TransportationOffice.Name)
		suite.Equal(move.LockExpiresAt.Truncate(time.Minute), expirationTime)
	})
}
