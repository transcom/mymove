package baselinetest

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/benbjohnson/clock"

	"github.com/transcom/mymove/pkg/auth"
)

func SetupFakeSessionManagers(clock clock.Clock, sessionStore scs.Store, useSecureCookie bool, idleTimeout time.Duration, lifetime time.Duration) auth.AppSessionManagers {
	var milSession, adminSession, officeSession *FakeSessionManager
	gob.Register(FakeSessionManager{})

	milSession = New(clock)
	milSession.ScsStore = sessionStore
	milSession.Cookie.Name = "mil_session_token"

	adminSession = New(clock)
	adminSession.ScsStore = sessionStore
	adminSession.Cookie.Name = "admin_session_token"

	officeSession = New(clock)
	officeSession.ScsStore = sessionStore
	officeSession.Cookie.Name = "office_session_token"

	// IdleTimeout controls the maximum length of time a session can be inactive
	// before it expires. The default is 15 minutes. To disable idle timeout in
	// a non-production environment, set SESSION_IDLE_TIMEOUT_IN_MINUTES to 0.
	milSession.IdleTimeout = idleTimeout
	adminSession.IdleTimeout = idleTimeout
	officeSession.IdleTimeout = idleTimeout

	// Lifetime controls the maximum length of time that a session is valid for
	// before it expires. The lifetime is an 'absolute expiry' which is set when
	// the session is first created or renewed (such as when a user signs in)
	// and does not change. The default value is 24 hours.
	milSession.Lifetime = lifetime
	adminSession.Lifetime = lifetime
	officeSession.Lifetime = lifetime

	milSession.Cookie.Path = "/"
	adminSession.Cookie.Path = "/"
	officeSession.Cookie.Path = "/"

	// A value of false means the session cookie will be deleted when the
	// browser is closed.
	milSession.Cookie.Persist = false
	adminSession.Cookie.Persist = false
	officeSession.Cookie.Persist = false

	if useSecureCookie {
		milSession.Cookie.Secure = true
		adminSession.Cookie.Secure = true
		officeSession.Cookie.Secure = true
	}

	return auth.NewAppSessionManagers(
		milSession,
		officeSession,
		adminSession,
	)
}

// this is mostly a copy from scs.SessionManager
type contextKey string

var (
	contextKeyID      uint64
	contextKeyIDMutex = &sync.Mutex{}
)

type Status int

const (
	Unmodified Status = iota

	Modified

	Destroyed
)

type FakeSessionManager struct {
	IdleTimeout time.Duration

	Lifetime time.Duration

	ScsStore scs.Store

	Cookie scs.SessionCookie

	Codec scs.Codec

	ErrorFunc func(http.ResponseWriter, *http.Request, error)

	contextKey contextKey

	clock clock.Clock

	counter uint32
}

func generateContextKey() contextKey {
	contextKeyIDMutex.Lock()
	defer contextKeyIDMutex.Unlock()
	atomic.AddUint64(&contextKeyID, 1)
	return contextKey(fmt.Sprintf("session.%d", contextKeyID))
}

func New(clock clock.Clock) *FakeSessionManager {
	s := &FakeSessionManager{
		IdleTimeout: 0,
		Lifetime:    24 * time.Hour,
		ScsStore:    memstore.New(),
		Codec:       scs.GobCodec{},
		ErrorFunc:   defaultErrorFunc,
		contextKey:  generateContextKey(),
		clock:       clock,
		Cookie: scs.SessionCookie{
			Name:     "session",
			Domain:   "",
			HttpOnly: true,
			Path:     "/",
			Persist:  true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		},
	}
	return s
}

func (s *FakeSessionManager) Store() scs.Store {
	return s.ScsStore
}

func (s *FakeSessionManager) Get(ctx context.Context, key string) interface{} {
	sd := s.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	return sd.values[key]
}

func (s *FakeSessionManager) Put(ctx context.Context, key string, val interface{}) {
	sd := s.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	sd.values[key] = val
	sd.status = Modified
	sd.mu.Unlock()
}

func (s *FakeSessionManager) Destroy(ctx context.Context) error {
	sd := s.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	err := s.doStoreDelete(ctx, sd.token)
	if err != nil {
		return err
	}

	sd.status = Destroyed

	// Reset everything else to defaults.
	sd.token = ""
	sd.deadline = s.clock.Now().Add(s.Lifetime).UTC()
	for key := range sd.values {
		delete(sd.values, key)
	}

	return nil
}

func (s *FakeSessionManager) RenewToken(ctx context.Context) error {
	sd := s.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	err := s.doStoreDelete(ctx, sd.token)
	if err != nil {
		return err
	}

	newToken, err := s.generateToken()
	if err != nil {
		return err
	}

	sd.token = newToken
	sd.deadline = s.clock.Now().Add(s.Lifetime).UTC()
	sd.status = Modified

	return nil
}

func (s *FakeSessionManager) Commit(ctx context.Context) (string, time.Time, error) {
	sd := s.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.token == "" {
		var err error
		if sd.token, err = s.generateToken(); err != nil {
			return "", time.Time{}, err
		}
	}

	b, err := s.Codec.Encode(sd.deadline, sd.values)
	if err != nil {
		return "", time.Time{}, err
	}

	expiry := sd.deadline
	if s.IdleTimeout > 0 {
		ie := s.clock.Now().Add(s.IdleTimeout).UTC()
		if ie.Before(expiry) {
			expiry = ie
		}
	}

	if err := s.doStoreCommit(ctx, sd.token, b, expiry); err != nil {
		return "", time.Time{}, err
	}

	return sd.token, expiry, nil
}

func (s *FakeSessionManager) LoadAndSave(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		cookie, err := r.Cookie(s.Cookie.Name)
		if err == nil {
			token = cookie.Value
		}

		ctx, err := s.Load(r.Context(), token)
		if err != nil {
			s.ErrorFunc(w, r, err)
			return
		}

		sr := r.WithContext(ctx)
		bw := &bufferedResponseWriter{ResponseWriter: w}
		next.ServeHTTP(bw, sr)

		if sr.MultipartForm != nil {
			err = sr.MultipartForm.RemoveAll()
			if err != nil {
				s.ErrorFunc(w, r, err)
				return
			}
		}

		switch s.Status(ctx) {
		case Modified:
			token, expiry, cerr := s.Commit(ctx)
			if err != nil {
				s.ErrorFunc(w, r, cerr)
				return
			}

			s.WriteSessionCookie(ctx, w, token, expiry)
		case Destroyed:
			s.WriteSessionCookie(ctx, w, "", time.Time{})
		}

		w.Header().Add("Vary", "Cookie")

		if bw.code != 0 {
			w.WriteHeader(bw.code)
		}
		_, err = w.Write(bw.buf.Bytes())
		if err != nil {
			// too late to respond with an error
			_ = log.Output(2, err.Error())
		}
	})
}

func (s *FakeSessionManager) Load(ctx context.Context, token string) (context.Context, error) {
	if _, ok := ctx.Value(s.contextKey).(*sessionData); ok {
		return ctx, nil
	}

	if token == "" {
		return s.addSessionDataToContext(ctx, newSessionData(s.clock, s.Lifetime)), nil
	}

	b, found, err := s.doStoreFind(ctx, token)
	if err != nil {
		return nil, err
	} else if !found {
		return s.addSessionDataToContext(ctx, newSessionData(s.clock, s.Lifetime)), nil
	}

	sd := &sessionData{
		status: Unmodified,
		token:  token,
	}
	if sd.deadline, sd.values, err = s.Codec.Decode(b); err != nil {
		return nil, err
	}

	// Mark the session data as modified if an idle timeout is being used. This
	// will force the session data to be re-committed to the session store with
	// a new expiry time.
	if s.IdleTimeout > 0 {
		sd.status = Modified
	}

	return s.addSessionDataToContext(ctx, sd), nil
}

func (s *FakeSessionManager) Status(ctx context.Context) Status {
	sd := s.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	return sd.status
}

func (s *FakeSessionManager) WriteSessionCookie(ctx context.Context, w http.ResponseWriter, token string, expiry time.Time) {
	cookie := &http.Cookie{
		Name:     s.Cookie.Name,
		Value:    token,
		Path:     s.Cookie.Path,
		Domain:   s.Cookie.Domain,
		Secure:   s.Cookie.Secure,
		HttpOnly: s.Cookie.HttpOnly,
		SameSite: s.Cookie.SameSite,
	}

	if expiry.IsZero() {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	} else if s.Cookie.Persist || s.GetBool(ctx, "__rememberMe") {
		cookie.Expires = time.Unix(expiry.Unix()+1, 0)        // Round up to the nearest second.
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1) // Round up to the nearest second.
	}

	w.Header().Add("Set-Cookie", cookie.String())
	w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
}

func (s *FakeSessionManager) GetBool(ctx context.Context, key string) bool {
	val := s.Get(ctx, key)
	b, ok := val.(bool)
	if !ok {
		return false
	}
	return b
}

type sessionData struct {
	deadline time.Time
	status   Status
	token    string
	values   map[string]interface{}
	mu       sync.Mutex
}

func newSessionData(clock clock.Clock, lifetime time.Duration) *sessionData {
	return &sessionData{
		deadline: clock.Now().Add(lifetime).UTC(),
		status:   Unmodified,
		values:   make(map[string]interface{}),
	}
}

// DETERMINISTIC generateToken
func (s *FakeSessionManager) generateToken() (string, error) {
	b := make([]byte, 32)
	binary.LittleEndian.PutUint32(b, s.counter)
	s.counter++
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *FakeSessionManager) getSessionDataFromContext(ctx context.Context) *sessionData {
	c, ok := ctx.Value(s.contextKey).(*sessionData)
	if !ok {
		panic("scs: no session data in context")
	}
	return c
}

func (s *FakeSessionManager) doStoreDelete(ctx context.Context, token string) (err error) {
	c, ok := s.ScsStore.(interface {
		DeleteCtx(context.Context, string) error
	})
	if ok {
		return c.DeleteCtx(ctx, token)
	}
	return s.ScsStore.Delete(token)
}

func (s *FakeSessionManager) addSessionDataToContext(ctx context.Context, sd *sessionData) context.Context {
	return context.WithValue(ctx, s.contextKey, sd)
}

func (s *FakeSessionManager) doStoreFind(ctx context.Context, token string) (b []byte, found bool, err error) {
	c, ok := s.ScsStore.(interface {
		FindCtx(context.Context, string) ([]byte, bool, error)
	})
	if ok {
		return c.FindCtx(ctx, token)
	}
	return s.ScsStore.Find(token)
}

func (s *FakeSessionManager) doStoreCommit(ctx context.Context, token string, b []byte, expiry time.Time) (err error) {
	c, ok := s.ScsStore.(interface {
		CommitCtx(context.Context, string, []byte, time.Time) error
	})
	if ok {
		return c.CommitCtx(ctx, token, b, expiry)
	}
	return s.ScsStore.Commit(token, b, expiry)
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf         bytes.Buffer
	code        int
	wroteHeader bool
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	if !bw.wroteHeader {
		bw.code = code
		bw.wroteHeader = true
	}
}

func (bw *bufferedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := bw.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (bw *bufferedResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := bw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

func defaultErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	_ = log.Output(2, err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
