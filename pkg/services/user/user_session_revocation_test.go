// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
// RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
// RA: in which this would be considered a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package user

import (
	"context"
	"reflect"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *UserServiceSuite) TestRevokeMilUserSession() {
	boolean := true
	payload := &adminmessages.UserUpdate{
		RevokeMilSession: &boolean,
	}
	newUUID, _ := uuid.NewV4()
	sessionManagers := auth.SetupSessionManagers(nil, false, time.Duration(180*time.Second), time.Duration(180*time.Second))

	ctx := context.Background()
	ctx, err := sessionManagers.Mil.Load(ctx, "fake_token")
	suite.NoError(err)
	sessionID, _, err := sessionManagers.Mil.Commit(ctx)
	suite.NoError(err)

	fakeUpdateOne := func(appcontext.AppContext, interface{}, *string) (*validate.Errors, error) {
		return nil, nil
	}
	fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
		reflect.ValueOf(model).Elem().FieldByName("CurrentMilSessionID").Set(reflect.ValueOf(sessionID))
		return nil
	}
	builder := &testUserQueryBuilder{
		fakeFetchOne:  fakeFetchOne,
		fakeUpdateOne: fakeUpdateOne,
	}
	updater := NewUserSessionRevocation(builder)

	suite.Run("Key is removed from Redis when boolean is true", func() {
		_, existsBefore, _ := sessionManagers.Mil.Store().Find(sessionID)

		suite.True(existsBefore)

		_, verrs, revokeErr := updater.RevokeUserSession(suite.AppContextForTest(), newUUID, payload, sessionManagers)
		_, existsAfter, _ := sessionManagers.Mil.Store().Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.False(existsAfter)
	})

	suite.Run("Key is not removed from Redis when boolean is false", func() {
		boolean = false
		payload = &adminmessages.UserUpdate{
			RevokeMilSession: &boolean,
		}
		sessionID, _, err := sessionManagers.Mil.Commit(ctx)
		suite.NoError(err)
		_, existsBefore, _ := sessionManagers.Mil.Store().Find(sessionID)

		suite.Equal(existsBefore, true)

		_, verrs, revokeErr := updater.RevokeUserSession(suite.AppContextForTest(), newUUID, payload, sessionManagers)
		_, exists, _ := sessionManagers.Mil.Store().Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.True(exists)
	})

	suite.Run("Returns an error if user is not found", func() {
		fakeUpdateOne := func(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return models.ErrFetchNotFound
		}

		builder := &testUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewUserSessionRevocation(builder)
		_, _, err := updater.RevokeUserSession(suite.AppContextForTest(), newUUID, payload, sessionManagers)

		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())
	})
}

func (suite *UserServiceSuite) TestRevokeAdminUserSession() {
	newUUID, _ := uuid.NewV4()
	sessionManagers := auth.SetupSessionManagers(nil, false, time.Duration(180*time.Second), time.Duration(180*time.Second))
	ctx := context.Background()
	ctx, err := sessionManagers.Admin.Load(ctx, "fake_token")
	suite.NoError(err)
	sessionID, _, err := sessionManagers.Admin.Commit(ctx)
	suite.NoError(err)

	fakeUpdateOne := func(appcontext.AppContext, interface{}, *string) (*validate.Errors, error) {
		return nil, nil
	}
	fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
		reflect.ValueOf(model).Elem().FieldByName("CurrentAdminSessionID").Set(reflect.ValueOf(sessionID))
		return nil
	}
	builder := &testUserQueryBuilder{
		fakeFetchOne:  fakeFetchOne,
		fakeUpdateOne: fakeUpdateOne,
	}
	updater := NewUserSessionRevocation(builder)

	boolean := true
	payload := &adminmessages.UserUpdate{
		RevokeAdminSession: &boolean,
	}

	suite.Run("Key is removed from Redis when boolean is true", func() {
		_, existsBefore, _ := sessionManagers.Admin.Store().Find(sessionID)

		suite.Equal(existsBefore, true)

		_, verrs, revokeErr := updater.RevokeUserSession(suite.AppContextForTest(), newUUID, payload, sessionManagers)
		_, existsAfter, _ := sessionManagers.Admin.Store().Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(existsAfter, false)
	})

	suite.Run("Key is not removed from Redis when boolean is false", func() {
		sessionID, _, err := sessionManagers.Admin.Commit(ctx)
		suite.NoError(err)
		_, existsBefore, _ := sessionManagers.Admin.Store().Find(sessionID)

		suite.Equal(existsBefore, true)

		boolean = false
		payload = &adminmessages.UserUpdate{
			RevokeAdminSession: &boolean,
		}

		_, verrs, revokeErr := updater.RevokeUserSession(suite.AppContextForTest(), newUUID, payload, sessionManagers)
		_, exists, _ := sessionManagers.Admin.Store().Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(exists, true)
	})
}

func (suite *UserServiceSuite) TestRevokeOfficeUserSession() {
	newUUID, _ := uuid.NewV4()
	sessionManagers := auth.SetupSessionManagers(nil, false, time.Duration(180*time.Second), time.Duration(180*time.Second))

	ctx := context.Background()
	ctx, err := sessionManagers.Office.Load(ctx, "fake_token")
	suite.NoError(err)
	sessionID, _, err := sessionManagers.Office.Commit(ctx)
	suite.NoError(err)

	fakeUpdateOne := func(appcontext.AppContext, interface{}, *string) (*validate.Errors, error) {
		return nil, nil
	}
	fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
		reflect.ValueOf(model).Elem().FieldByName("CurrentOfficeSessionID").Set(reflect.ValueOf(sessionID))
		return nil
	}
	builder := &testUserQueryBuilder{
		fakeFetchOne:  fakeFetchOne,
		fakeUpdateOne: fakeUpdateOne,
	}
	updater := NewUserSessionRevocation(builder)

	boolean := true
	payload := &adminmessages.UserUpdate{
		RevokeOfficeSession: &boolean,
	}

	suite.Run("Key is removed from Redis when boolean is true", func() {
		_, existsBefore, _ := sessionManagers.Office.Store().Find(sessionID)

		suite.Equal(existsBefore, true)

		_, verrs, revokeErr := updater.RevokeUserSession(suite.AppContextForTest(), newUUID, payload, sessionManagers)
		_, existsAfter, _ := sessionManagers.Office.Store().Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(existsAfter, false)
	})

	suite.Run("Key is not removed from Redis when boolean is false", func() {
		boolean = false
		payload = &adminmessages.UserUpdate{
			RevokeOfficeSession: &boolean,
		}
		sessionID, _, err := sessionManagers.Office.Commit(ctx)
		suite.NoError(err)

		_, existsBefore, _ := sessionManagers.Office.Store().Find(sessionID)

		suite.Equal(existsBefore, true)

		_, verrs, revokeErr := updater.RevokeUserSession(suite.AppContextForTest(), newUUID, payload, sessionManagers)
		_, exists, _ := sessionManagers.Office.Store().Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(exists, true)
	})
}

func (suite *UserServiceSuite) TestRevokeMultipleSessions() {
	newUUID, _ := uuid.NewV4()
	sessionManagers := auth.SetupSessionManagers(nil, false, time.Duration(180*time.Second), time.Duration(180*time.Second))

	officeCtx := context.Background()
	officeCtx, err := sessionManagers.Office.Load(officeCtx, "fake_token")
	suite.NoError(err)
	officeSessionID, _, err := sessionManagers.Office.Commit(officeCtx)
	suite.NoError(err)

	milCtx := context.Background()
	milCtx, err = sessionManagers.Mil.Load(milCtx, "fake_token")
	suite.NoError(err)
	milSessionID, _, err := sessionManagers.Mil.Commit(milCtx)
	suite.NoError(err)

	adminCtx := context.Background()
	adminCtx, err = sessionManagers.Admin.Load(adminCtx, "fake_token")
	suite.NoError(err)
	adminSessionID, _, err := sessionManagers.Admin.Commit(adminCtx)
	suite.NoError(err)

	fakeUpdateOne := func(appcontext.AppContext, interface{}, *string) (*validate.Errors, error) {
		return nil, nil
	}
	fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
		reflect.ValueOf(model).Elem().FieldByName("CurrentOfficeSessionID").Set(reflect.ValueOf(officeSessionID))
		reflect.ValueOf(model).Elem().FieldByName("CurrentMilSessionID").Set(reflect.ValueOf(milSessionID))
		reflect.ValueOf(model).Elem().FieldByName("CurrentAdminSessionID").Set(reflect.ValueOf(adminSessionID))
		return nil
	}
	builder := &testUserQueryBuilder{
		fakeFetchOne:  fakeFetchOne,
		fakeUpdateOne: fakeUpdateOne,
	}
	updater := NewUserSessionRevocation(builder)

	boolean := true
	payload := &adminmessages.UserUpdate{
		RevokeOfficeSession: &boolean,
		RevokeAdminSession:  &boolean,
		RevokeMilSession:    &boolean,
	}

	suite.Run("All keys are removed from Redis when boolean is true", func() {
		_, adminExistsBefore, _ := sessionManagers.Admin.Store().Find(adminSessionID)
		_, officeExistsBefore, _ := sessionManagers.Office.Store().Find(officeSessionID)
		_, milExistsBefore, _ := sessionManagers.Mil.Store().Find(milSessionID)

		suite.True(adminExistsBefore)
		suite.True(officeExistsBefore)
		suite.True(milExistsBefore)

		_, verrs, revokeErr := updater.RevokeUserSession(suite.AppContextForTest(), newUUID, payload, sessionManagers)
		_, adminExistsAfter, _ := sessionManagers.Admin.Store().Find(adminSessionID)
		_, officeExistsAfter, _ := sessionManagers.Office.Store().Find(officeSessionID)
		_, milExistsAfter, _ := sessionManagers.Mil.Store().Find(milSessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.False(adminExistsAfter)
		suite.False(officeExistsAfter)
		suite.False(milExistsAfter)
	})
}
