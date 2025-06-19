package lockmove

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *MoveLockerServiceSuite) TestMoveUnlocker() {
	moveLocker := NewMoveLocker()
	moveUnlocker := NewMoveUnlocker()

	suite.Run("successfully returns move with no values in locked_by or lock_expires_at column", func() {
		tooUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		defaultRole, err := tooUser.User.Roles.Default()
		suite.FatalNoError(err)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			ActiveRole:      *defaultRole,
			OfficeUserID:    tooUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		})

		// build the move so we can lock it and unlock it
		move := factory.BuildMove(suite.DB(), nil, nil)

		// lock the move first
		lockedMove, err := moveLocker.LockMove(appCtx, &move, tooUser.ID)
		suite.FatalNoError(err)
		suite.Equal(move.ID, lockedMove.ID)
		suite.Equal(lockedMove.LockedByOfficeUserID, &tooUser.ID)

		// now let's unlock it
		unlockedMove, err := moveUnlocker.UnlockMove(appCtx, lockedMove, tooUser.ID)
		suite.FatalNoError(err)

		// all values should now be nil
		suite.Equal(move.ID, unlockedMove.ID)
		suite.Nil(unlockedMove.LockedByOfficeUserID)
		suite.Nil(unlockedMove.LockedByOfficeUser)
		suite.Nil(unlockedMove.LockExpiresAt)
	})
}

func (suite *MoveLockerServiceSuite) TestCheckForLockedMovesAndUnlock() {
	moveLocker := NewMoveLocker()
	moveUnlocker := NewMoveUnlocker()

	suite.Run("successfully clears all moves that user has locked", func() {
		tooUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		defaultRole, err := tooUser.User.Roles.Default()
		suite.FatalNoError(err)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			ActiveRole:      *defaultRole,
			OfficeUserID:    tooUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		})

		// build some moves so we can lock them and unlock them all
		move := factory.BuildMove(suite.DB(), nil, nil)
		moveTwo := factory.BuildMove(suite.DB(), nil, nil)
		moveThree := factory.BuildMove(suite.DB(), nil, nil)

		// lock the moves
		lockedMove, err := moveLocker.LockMove(appCtx, &move, tooUser.ID)
		suite.FatalNoError(err)
		suite.Equal(move.ID, lockedMove.ID)
		suite.Equal(lockedMove.LockedByOfficeUserID, &tooUser.ID)

		lockedMoveTwo, err := moveLocker.LockMove(appCtx, &moveTwo, tooUser.ID)
		suite.FatalNoError(err)
		suite.Equal(moveTwo.ID, lockedMoveTwo.ID)
		suite.Equal(lockedMoveTwo.LockedByOfficeUserID, &tooUser.ID)

		lockedMoveThree, err := moveLocker.LockMove(appCtx, &moveThree, tooUser.ID)
		suite.FatalNoError(err)
		suite.Equal(moveThree.ID, lockedMoveThree.ID)
		suite.Equal(lockedMoveThree.LockedByOfficeUserID, &tooUser.ID)

		// now let's unlock them by calling CheckForUnlockedMoves
		err = moveUnlocker.CheckForLockedMovesAndUnlock(appCtx, tooUser.ID)
		suite.FatalNoError(err)

		// all values should now be nil in all the moves
		// find the moves in the database and verify
		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Nil(moveInDB.LockedByOfficeUserID)
		suite.Nil(moveInDB.LockedByOfficeUser)
		suite.Nil(moveInDB.LockExpiresAt)

		var moveTwoInDB models.Move
		err = suite.DB().Find(&moveTwoInDB, moveTwo.ID)
		suite.NoError(err)
		suite.Nil(moveTwoInDB.LockedByOfficeUserID)
		suite.Nil(moveTwoInDB.LockedByOfficeUser)
		suite.Nil(moveTwoInDB.LockExpiresAt)

		var moveThreeInDB models.Move
		err = suite.DB().Find(&moveThreeInDB, moveThree.ID)
		suite.NoError(err)
		suite.Nil(moveThreeInDB.LockedByOfficeUserID)
		suite.Nil(moveThreeInDB.LockedByOfficeUser)
		suite.Nil(moveThreeInDB.LockExpiresAt)
	})

	suite.Run("successfully unlock move without changing updated_at", func() {
		tooUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		defaultRole, err := tooUser.User.Roles.Default()
		suite.FatalNoError(err)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			ActiveRole:      *defaultRole,
			OfficeUserID:    tooUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		})

		move := factory.BuildMove(suite.DB(), nil, nil)

		lockedMove, err := moveLocker.LockMove(appCtx, &move, tooUser.ID)
		suite.FatalNoError(err)
		suite.Equal(move.ID, lockedMove.ID)
		suite.Equal(lockedMove.LockedByOfficeUserID, &tooUser.ID)
		suite.Equal(lockedMove.UpdatedAt, move.UpdatedAt)

		unlockedMove, err := moveUnlocker.UnlockMove(appCtx, lockedMove, tooUser.ID)
		suite.FatalNoError(err)
		suite.Equal(unlockedMove.UpdatedAt, move.UpdatedAt)
	})
}
