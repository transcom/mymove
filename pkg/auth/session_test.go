package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gomodule/redigo/redis"
)

func (suite *authSuite) TestSetupSessionManagers() {
	idleTimeout := 15 * time.Minute
	lifetime := 24 * time.Hour
	useSecureCookie := true
	sessionStore := memstore.New()

	sessionManagers := SetupSessionManagers(
		sessionStore, useSecureCookie, idleTimeout, lifetime,
	)
	wrapper, ok := sessionManagers.mil.(ScsSessionManagerWrapper)
	if !ok {
		suite.Assert().FailNow("Cannot cast mil to ScsSessionManagerWrapper",
			sessionManagers.mil)
	}
	milSession := wrapper.ScsSessionManager
	wrapper, ok = sessionManagers.admin.(ScsSessionManagerWrapper)
	if !ok {
		suite.Assert().FailNow("Cannot cast admin to ScsSessionManagerWrapper",
			sessionManagers.admin)
	}
	adminSession := wrapper.ScsSessionManager
	wrapper, ok = sessionManagers.office.(ScsSessionManagerWrapper)
	if !ok {
		suite.Assert().FailNow("Cannot cast office to ScsSessionManagerWrapper",
			sessionManagers.office)
	}
	officeSession := wrapper.ScsSessionManager

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
		store := redisstore.New(pool)
		sessionManagers := SetupSessionManagers(
			store, useSecureCookie, idleTimeout, lifetime,
		)
		wrapper, ok := sessionManagers.mil.(ScsSessionManagerWrapper)
		if !ok {
			suite.Assert().FailNow("Cannot cast mil to ScsSessionManagerWrapper",
				sessionManagers.mil)
		}
		milSessionManager := wrapper.ScsSessionManager
		ctx := context.Background()
		fakeSessionID := "fakeid"
		fakeSession := &Session{
			Hostname: "fake",
		}
		expiry := time.Now().Add(30 * time.Minute).UTC()
		b, err := milSessionManager.Codec.Encode(expiry,
			map[string]interface{}{
				fakeSessionID: fakeSession,
			})
		suite.NoError(err)

		// make sure we can store and load from redis without error
		suite.NoError(milSessionManager.Store.Commit("token", b, expiry))
		_, err = milSessionManager.Load(ctx, "session")
		suite.NoError(err)
	})

	suite.Run("With a supported scs.Store other than redisstore", func() {
		sessionStore = memstore.New()
		sessionManagers := SetupSessionManagers(
			sessionStore, useSecureCookie, idleTimeout, lifetime,
		)
		suite.Equal(sessionStore, sessionManagers.mil.Store())
		suite.Equal(sessionStore, sessionManagers.admin.Store())
		suite.Equal(sessionStore, sessionManagers.office.Store())

	})

	suite.Run("Session cookie names must be unique per app", func() {
		suite.Equal("mil_session_token", milSession.Cookie.Name)
		suite.Equal("admin_session_token", adminSession.Cookie.Name)
		suite.Equal("office_session_token", officeSession.Cookie.Name)
	})

	suite.Run("All session managers have the same secure cookie setting", func() {
		suite.Equal(useSecureCookie, milSession.Cookie.Secure)
		suite.Equal(useSecureCookie, adminSession.Cookie.Secure)
		suite.Equal(useSecureCookie, officeSession.Cookie.Secure)
	})

	suite.Run("All session managers have the same idleTimeout", func() {
		suite.Equal(idleTimeout, milSession.IdleTimeout)
		suite.Equal(idleTimeout, adminSession.IdleTimeout)
		suite.Equal(idleTimeout, officeSession.IdleTimeout)
	})

	suite.Run("All session managers have the same lifetime", func() {
		suite.Equal(lifetime, milSession.Lifetime)
		suite.Equal(lifetime, adminSession.Lifetime)
		suite.Equal(lifetime, officeSession.Lifetime)
	})

	suite.Run("All session managers have cookie path set to root", func() {
		suite.Equal("/", milSession.Cookie.Path)
		suite.Equal("/", adminSession.Cookie.Path)
		suite.Equal("/", officeSession.Cookie.Path)
	})

	suite.Run("All session managers do not persist cookie", func() {
		suite.Equal(false, milSession.Cookie.Persist)
		suite.Equal(false, adminSession.Cookie.Persist)
		suite.Equal(false, officeSession.Cookie.Persist)
	})
}
