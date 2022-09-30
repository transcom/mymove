package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gomodule/redigo/redis"
)

func (suite *authSuite) TestSetupSessionManagers() {
	idleTimeout := 15 * time.Minute
	lifetime := 24 * time.Hour
	useSecureCookie := true

	suite.Run("With redis enabled", func() {
		// on local dev machines, this shares the same redis server as
		// development. Should we spin up a separate server for test?
		// Use the same server but a different redis database?
		redisHost := os.Getenv("REDIS_HOST")
		redisPort, ok := os.LookupEnv("REDIS_PORT")
		if !ok {
			redisPort = "6379"
		}

		redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
		pool := &redis.Pool{
			Dial: func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
		}
		sessionManagers := SetupSessionManagers(
			pool, useSecureCookie, idleTimeout, lifetime,
		)
		ctx := context.Background()
		//		fakeSessionID := "fakeid"
		fakeSession := &Session{
			Hostname: "fake",
		}

		// make sure we can store and load from redis without error
		fakeToken := "fake_token"
		ctx, err := sessionManagers.Mil.Load(ctx, fakeToken)
		suite.NoError(err)
		sessionToken, _, err := sessionManagers.Mil.Commit(ctx)
		suite.NoError(err)
		sessionManagers.Mil.Put(ctx, "session", fakeSession)
		_, err = sessionManagers.Mil.Load(ctx, "session")
		suite.NoError(err)

		// make sure all stores are unique
		// before the change to use a separate redis prefix for each
		// session manager, this test would fail
		_, found, err := sessionManagers.Mil.Store().Find(sessionToken)
		suite.NoError(err)
		suite.True(found)

		_, found, err = sessionManagers.Admin.Store().Find(sessionToken)
		suite.NoError(err)
		suite.False(found)

		_, found, err = sessionManagers.Office.Store().Find(sessionToken)
		suite.NoError(err)
		suite.False(found)
	})

	suite.Run("With a supported scs.Store other than redisstore", func() {
		sessionManagers := SetupSessionManagers(
			nil, useSecureCookie, idleTimeout, lifetime,
		)
		_, ok := sessionManagers.Mil.Store().(*memstore.MemStore)
		suite.Require().True(ok)
		_, ok = sessionManagers.Admin.Store().(*memstore.MemStore)
		suite.Require().True(ok)
		_, ok = sessionManagers.Office.Store().(*memstore.MemStore)
		suite.Require().True(ok)
	})

	suite.Run("Session cookie names must be unique per app", func() {
		sessionManagers := SetupSessionManagers(
			nil, useSecureCookie, idleTimeout, lifetime,
		)
		milSession := sessionManagers.Mil.(ScsSessionManagerWrapper).ScsSessionManager
		adminSession := sessionManagers.Admin.(ScsSessionManagerWrapper).ScsSessionManager
		officeSession := sessionManagers.Office.(ScsSessionManagerWrapper).ScsSessionManager

		suite.Equal("mil_session_token", milSession.Cookie.Name)
		suite.Equal("admin_session_token", adminSession.Cookie.Name)
		suite.Equal("office_session_token", officeSession.Cookie.Name)
	})

	suite.Run("All session managers have the same secure cookie setting", func() {
		sessionManagers := SetupSessionManagers(
			nil, useSecureCookie, idleTimeout, lifetime,
		)
		milSession := sessionManagers.Mil.(ScsSessionManagerWrapper).ScsSessionManager
		adminSession := sessionManagers.Admin.(ScsSessionManagerWrapper).ScsSessionManager
		officeSession := sessionManagers.Office.(ScsSessionManagerWrapper).ScsSessionManager

		suite.Equal(useSecureCookie, milSession.Cookie.Secure)
		suite.Equal(useSecureCookie, adminSession.Cookie.Secure)
		suite.Equal(useSecureCookie, officeSession.Cookie.Secure)
	})

	suite.Run("All session managers have the same idleTimeout", func() {
		sessionManagers := SetupSessionManagers(
			nil, useSecureCookie, idleTimeout, lifetime,
		)
		milSession := sessionManagers.Mil.(ScsSessionManagerWrapper).ScsSessionManager
		adminSession := sessionManagers.Admin.(ScsSessionManagerWrapper).ScsSessionManager
		officeSession := sessionManagers.Office.(ScsSessionManagerWrapper).ScsSessionManager

		suite.Equal(idleTimeout, milSession.IdleTimeout)
		suite.Equal(idleTimeout, adminSession.IdleTimeout)
		suite.Equal(idleTimeout, officeSession.IdleTimeout)
	})

	suite.Run("All session managers have the same lifetime", func() {
		sessionManagers := SetupSessionManagers(
			nil, useSecureCookie, idleTimeout, lifetime,
		)
		milSession := sessionManagers.Mil.(ScsSessionManagerWrapper).ScsSessionManager
		adminSession := sessionManagers.Admin.(ScsSessionManagerWrapper).ScsSessionManager
		officeSession := sessionManagers.Office.(ScsSessionManagerWrapper).ScsSessionManager

		suite.Equal(lifetime, milSession.Lifetime)
		suite.Equal(lifetime, adminSession.Lifetime)
		suite.Equal(lifetime, officeSession.Lifetime)
	})

	suite.Run("All session managers have cookie path set to root", func() {
		sessionManagers := SetupSessionManagers(
			nil, useSecureCookie, idleTimeout, lifetime,
		)
		milSession := sessionManagers.Mil.(ScsSessionManagerWrapper).ScsSessionManager
		adminSession := sessionManagers.Admin.(ScsSessionManagerWrapper).ScsSessionManager
		officeSession := sessionManagers.Office.(ScsSessionManagerWrapper).ScsSessionManager

		suite.Equal("/", milSession.Cookie.Path)
		suite.Equal("/", adminSession.Cookie.Path)
		suite.Equal("/", officeSession.Cookie.Path)
	})

	suite.Run("All session managers do not persist cookie", func() {
		sessionManagers := SetupSessionManagers(
			nil, useSecureCookie, idleTimeout, lifetime,
		)
		milSession := sessionManagers.Mil.(ScsSessionManagerWrapper).ScsSessionManager
		adminSession := sessionManagers.Admin.(ScsSessionManagerWrapper).ScsSessionManager
		officeSession := sessionManagers.Office.(ScsSessionManagerWrapper).ScsSessionManager

		suite.Equal(false, milSession.Cookie.Persist)
		suite.Equal(false, adminSession.Cookie.Persist)
		suite.Equal(false, officeSession.Cookie.Persist)
	})
}
