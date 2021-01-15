//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
//RA Validator: jneuner@mitre.org
//RA Modified Severity:
// nolint:errcheck
package user

import (
	"reflect"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *UserServiceSuite) TestRevokeMilUserSession() {
	newUUID, _ := uuid.NewV4()
	sessionStore := memstore.New()
	sessionID := "mil_session_token"
	sessionStore.Commit(sessionID, []byte("encoded_data"), time.Now().Add(time.Minute))

	fakeUpdateOne := func(interface{}, *string) (*validate.Errors, error) {
		return nil, nil
	}
	fakeFetchOne := func(model interface{}) error {
		reflect.ValueOf(model).Elem().FieldByName("CurrentMilSessionID").Set(reflect.ValueOf(sessionID))
		return nil
	}
	builder := &testUserQueryBuilder{
		fakeFetchOne:  fakeFetchOne,
		fakeUpdateOne: fakeUpdateOne,
	}
	updater := NewUserSessionRevocation(builder)

	boolean := true
	payload := &adminmessages.UserRevokeSessionPayload{
		RevokeMilSession: &boolean,
	}

	suite.T().Run("Key is removed from Redis when boolean is true", func(t *testing.T) {
		_, existsBefore, _ := sessionStore.Find(sessionID)

		suite.Equal(existsBefore, true)

		_, verrs, revokeErr := updater.RevokeUserSession(newUUID, payload, sessionStore)
		_, existsAfter, _ := sessionStore.Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(existsAfter, false)
	})

	suite.T().Run("Key is not removed from Redis when boolean is false", func(t *testing.T) {
		sessionStore.Commit(sessionID, []byte("encoded_data"), time.Now().Add(time.Minute))
		boolean = false
		payload = &adminmessages.UserRevokeSessionPayload{
			RevokeMilSession: &boolean,
		}

		_, verrs, revokeErr := updater.RevokeUserSession(newUUID, payload, sessionStore)
		_, exists, _ := sessionStore.Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(exists, true)
	})

	suite.T().Run("Returns an error if user is not found", func(t *testing.T) {
		fakeUpdateOne := func(model interface{}, eTag *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(model interface{}) error {
			return models.ErrFetchNotFound
		}

		builder := &testUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewUserSessionRevocation(builder)
		_, _, err := updater.RevokeUserSession(newUUID, payload, sessionStore)

		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())
	})
}

func (suite *UserServiceSuite) TestRevokeAdminUserSession() {
	newUUID, _ := uuid.NewV4()
	sessionStore := memstore.New()
	sessionID := "admin_session_token"
	sessionStore.Commit(sessionID, []byte("encoded_data"), time.Now().Add(time.Minute))

	fakeUpdateOne := func(interface{}, *string) (*validate.Errors, error) {
		return nil, nil
	}
	fakeFetchOne := func(model interface{}) error {
		reflect.ValueOf(model).Elem().FieldByName("CurrentAdminSessionID").Set(reflect.ValueOf(sessionID))
		return nil
	}
	builder := &testUserQueryBuilder{
		fakeFetchOne:  fakeFetchOne,
		fakeUpdateOne: fakeUpdateOne,
	}
	updater := NewUserSessionRevocation(builder)

	boolean := true
	payload := &adminmessages.UserRevokeSessionPayload{
		RevokeAdminSession: &boolean,
	}

	suite.T().Run("Key is removed from Redis when boolean is true", func(t *testing.T) {
		_, existsBefore, _ := sessionStore.Find(sessionID)

		suite.Equal(existsBefore, true)

		_, verrs, revokeErr := updater.RevokeUserSession(newUUID, payload, sessionStore)
		_, existsAfter, _ := sessionStore.Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(existsAfter, false)
	})

	suite.T().Run("Key is not removed from Redis when boolean is false", func(t *testing.T) {
		sessionStore.Commit(sessionID, []byte("encoded_data"), time.Now().Add(time.Minute))
		boolean = false
		payload = &adminmessages.UserRevokeSessionPayload{
			RevokeAdminSession: &boolean,
		}

		_, verrs, revokeErr := updater.RevokeUserSession(newUUID, payload, sessionStore)
		_, exists, _ := sessionStore.Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(exists, true)
	})
}

func (suite *UserServiceSuite) TestRevokeOfficeUserSession() {
	newUUID, _ := uuid.NewV4()
	sessionStore := memstore.New()
	sessionID := "office_session_token"
	sessionStore.Commit(sessionID, []byte("encoded_data"), time.Now().Add(time.Minute))

	fakeUpdateOne := func(interface{}, *string) (*validate.Errors, error) {
		return nil, nil
	}
	fakeFetchOne := func(model interface{}) error {
		reflect.ValueOf(model).Elem().FieldByName("CurrentOfficeSessionID").Set(reflect.ValueOf(sessionID))
		return nil
	}
	builder := &testUserQueryBuilder{
		fakeFetchOne:  fakeFetchOne,
		fakeUpdateOne: fakeUpdateOne,
	}
	updater := NewUserSessionRevocation(builder)

	boolean := true
	payload := &adminmessages.UserRevokeSessionPayload{
		RevokeOfficeSession: &boolean,
	}

	suite.T().Run("Key is removed from Redis when boolean is true", func(t *testing.T) {
		_, existsBefore, _ := sessionStore.Find(sessionID)

		suite.Equal(existsBefore, true)

		_, verrs, revokeErr := updater.RevokeUserSession(newUUID, payload, sessionStore)
		_, existsAfter, _ := sessionStore.Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(existsAfter, false)
	})

	suite.T().Run("Key is not removed from Redis when boolean is false", func(t *testing.T) {
		sessionStore.Commit(sessionID, []byte("encoded_data"), time.Now().Add(time.Minute))
		boolean = false
		payload = &adminmessages.UserRevokeSessionPayload{
			RevokeOfficeSession: &boolean,
		}

		_, verrs, revokeErr := updater.RevokeUserSession(newUUID, payload, sessionStore)
		_, exists, _ := sessionStore.Find(sessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(exists, true)
	})
}

func (suite *UserServiceSuite) TestRevokeMultipleSessions() {
	newUUID, _ := uuid.NewV4()
	sessionStore := memstore.New()
	officeSessionID := "office_session_token"
	milSessionID := "mil_session_token"
	adminSessionID := "admin_session_token"
	sessionStore.Commit(officeSessionID, []byte("encoded_data"), time.Now().Add(time.Minute))
	sessionStore.Commit(milSessionID, []byte("encoded_data"), time.Now().Add(time.Minute))
	sessionStore.Commit(adminSessionID, []byte("encoded_data"), time.Now().Add(time.Minute))

	fakeUpdateOne := func(interface{}, *string) (*validate.Errors, error) {
		return nil, nil
	}
	fakeFetchOne := func(model interface{}) error {
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
	payload := &adminmessages.UserRevokeSessionPayload{
		RevokeOfficeSession: &boolean,
		RevokeAdminSession:  &boolean,
		RevokeMilSession:    &boolean,
	}

	suite.T().Run("All keys are removed from Redis when boolean is true", func(t *testing.T) {
		_, adminExistsBefore, _ := sessionStore.Find(adminSessionID)
		_, officeExistsBefore, _ := sessionStore.Find(officeSessionID)
		_, milExistsBefore, _ := sessionStore.Find(milSessionID)

		suite.Equal(adminExistsBefore, true)
		suite.Equal(officeExistsBefore, true)
		suite.Equal(milExistsBefore, true)

		_, verrs, revokeErr := updater.RevokeUserSession(newUUID, payload, sessionStore)
		_, adminExistsAfter, _ := sessionStore.Find(adminSessionID)
		_, officeExistsAfter, _ := sessionStore.Find(officeSessionID)
		_, milExistsAfter, _ := sessionStore.Find(milSessionID)

		suite.NoError(revokeErr)
		suite.Nil(verrs)
		suite.Equal(adminExistsAfter, false)
		suite.Equal(officeExistsAfter, false)
		suite.Equal(milExistsAfter, false)
	})
}
