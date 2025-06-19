package lockmove

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models/roles"
	movefetcher "github.com/transcom/mymove/pkg/services/move"
)

func (suite *MoveLockerServiceSuite) TestLockMove() {
	moveLocker := NewMoveLocker()

	suite.Run("successfully returns move with office user values and lockExpiresAt value", func() {
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

	suite.Run("locking a move doesn't change the moves updated_at value", func() {
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

		actualMove, err := moveLocker.LockMove(appCtx, &move, tooUser.ID)
		suite.FatalNoError(err)

		suite.Equal(move.ID, actualMove.ID)
		suite.Equal(move.LockedByOfficeUserID, &tooUser.ID)
		suite.Equal(move.LockedByOfficeUser.TransportationOffice.Name, tooUser.TransportationOffice.Name)
		suite.Equal(actualMove.UpdatedAt, move.UpdatedAt)
	})
}

func (suite *MoveLockerServiceSuite) TestLockMoves() {
	moveLocker := NewMoveLocker()

	suite.Run("successfully returns move with office user values and lockExpiresAt value", func() {
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

		move1 := factory.BuildMove(suite.DB(), nil, nil)
		move2 := factory.BuildMove(suite.DB(), nil, nil)
		move3 := factory.BuildMove(suite.DB(), nil, nil)
		move4 := factory.BuildMove(suite.DB(), nil, nil)

		idsToLock := make([]uuid.UUID, 3)
		idsToLock[0] = move1.ID
		idsToLock[1] = move2.ID
		idsToLock[2] = move3.ID

		err = moveLocker.LockMoves(appCtx, idsToLock, tooUser.ID)
		suite.FatalNoError(err)

		moveFetcher := movefetcher.NewMoveFetcher()
		ids := []ghcmessages.BulkAssignmentMoveData{
			ghcmessages.BulkAssignmentMoveData(move1.ID.String()),
			ghcmessages.BulkAssignmentMoveData(move2.ID.String()),
			ghcmessages.BulkAssignmentMoveData(move3.ID.String()),
			ghcmessages.BulkAssignmentMoveData(move4.ID.String()),
		}
		moves, err := moveFetcher.FetchMovesByIdArray(suite.AppContextForTest(), ids)

		suite.NoError(err)
		suite.Len(moves, 4)

		// saving time and rounding time values to nearest minute to avoid nanosecond differences when testing
		now := time.Now().UTC()
		expirationTime := now.Add(30 * time.Minute).Truncate(time.Minute)

		for _, move := range moves {
			if move.ID == move4.ID {
				suite.Nil(move.LockedByOfficeUserID)
				suite.Nil(move.LockExpiresAt)
			} else {
				suite.Equal(move.LockedByOfficeUserID, &tooUser.ID)
				suite.Equal(move.LockExpiresAt.Truncate(time.Minute).UTC(), expirationTime)
			}
		}
	})

	suite.Run("locking moves doesn't change the moves' updated_at value", func() {
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

		move1 := factory.BuildMove(suite.DB(), nil, nil)
		move2 := factory.BuildMove(suite.DB(), nil, nil)

		idsToLock := make([]uuid.UUID, 3)
		idsToLock[0] = move1.ID
		idsToLock[1] = move2.ID

		err = moveLocker.LockMoves(appCtx, idsToLock, tooUser.ID)
		suite.FatalNoError(err)

		moveFetcher := movefetcher.NewMoveFetcher()
		ids := []ghcmessages.BulkAssignmentMoveData{
			ghcmessages.BulkAssignmentMoveData(move1.ID.String()),
			ghcmessages.BulkAssignmentMoveData(move2.ID.String()),
		}
		moves, err := moveFetcher.FetchMovesByIdArray(suite.AppContextForTest(), ids)

		suite.NoError(err)
		suite.Len(moves, 2)

		for _, move := range moves {
			if move.ID == move1.ID {
				suite.Equal(move.UpdatedAt.UTC(), move1.UpdatedAt.UTC())
			} else {
				suite.Equal(move.UpdatedAt.UTC(), move2.UpdatedAt.UTC())
			}
		}
	})
}
