package routing

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash"
	"net/http"
	"strconv"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
)

type contextKeyType string

const (
	gorillaCsrfName = "gorilla.csrf.Token"
	tokenLength     = 32
	maxAge          = 3600 * 12
	headerName      = "X-CSRF-Token"
)

var (
	safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}
)

type fakeCsrfState struct {
	clock   clock.Clock
	counter uint32
	hashKey []byte
	st      *cookieStore
}

var fakeState fakeCsrfState

type fakeCsrf struct {
	logger *zap.Logger
	orig   http.Handler
}

func GetFakeCSRFCookies() []http.Cookie {
	realToken := generateCounterBytes(tokenLength)
	maskedToken := maskToken(realToken)

	realCookie, _ := fakeState.st.GenerateCookie(realToken)

	// from auth/cookie.go
	maskedCookie := http.Cookie{
		Name:     auth.MaskedGorillaCSRFToken,
		Value:    maskedToken,
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}

	return []http.Cookie{*realCookie, maskedCookie}
}

// MaskedCSRFMiddleware handles setting the CSRF Token cookie
func NewFakeMaskedCSRFMiddleware(globalLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value(contextKeyType(gorillaCsrfName))
			maskedToken := ""
			if val != nil {
				if s, ok := val.(string); ok {
					maskedToken = s
				}
			}
			auth.WriteMaskedCSRFCookie(w, maskedToken, false)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}

func NewFakeCSRFMiddleware(clock clock.Clock, logger *zap.Logger) func(http.Handler) http.Handler {
	fakeHashKey := make([]byte, 32)
	binary.LittleEndian.PutUint32(fakeHashKey, 0xff00)
	sz := securecookie.JSONEncoder{}
	sc := securecookie.New(fakeHashKey, nil)
	sc.SetSerializer(sz)
	sc.MaxAge(maxAge)

	// initialize fake csrf state
	fakeState = fakeCsrfState{
		clock:   clock,
		counter: 0,
		hashKey: fakeHashKey,
		st: &cookieStore{
			name:     auth.GorillaCSRFToken,
			maxAge:   maxAge,
			secure:   false,
			httpOnly: true,
			sameSite: int(csrf.SameSiteLaxMode),
			path:     "/",
			domain:   "",
			sc:       sc,
			sz:       sz,
		},
	}

	return func(orig http.Handler) http.Handler {

		return &fakeCsrf{
			logger: logger,
			orig:   orig,
		}
	}
}

// Implements http.Handler for the csrf type.
func (cs *fakeCsrf) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	realToken, err := fakeState.st.Get(r)
	if err != nil || len(realToken) != tokenLength {
		realToken = generateCounterBytes(tokenLength)
		err = fakeState.st.Save(realToken, w)
		if err != nil {
			cs.logger.Error("error saving token", zap.Error(err))
			http.Error(w, "error saving token", http.StatusInternalServerError)
			return
		}
	}
	// Save the masked token to the request context
	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyType(gorillaCsrfName), maskToken(realToken))
	r = r.WithContext(ctx)

	if !contains(safeMethods, r.Method) {
		// Retrieve the combined token (pad + masked) token...
		maskedToken, err := cs.requestToken(r)
		if err != nil || maskedToken == nil {
			cs.logger.Error("error getting request token in fake csrf middleware", zap.Error(err))
			http.Error(w, "error getting request token", http.StatusInternalServerError)
			return
		}
		requestToken := unmaskToken(maskedToken)
		// Compare the request token against the real token
		if !compareTokens(requestToken, realToken) {
			cs.logger.Error("bad token", zap.Error(err))
			http.Error(w, "bad token", http.StatusInternalServerError)
			return
		}

	}
	// Set the Vary: Cookie header to protect clients from caching the response.
	w.Header().Add("Vary", "Cookie")
	cs.orig.ServeHTTP(w, r)
}

func generateCounterBytes(n int) []byte {
	b := make([]byte, tokenLength)
	binary.LittleEndian.PutUint32(b, fakeState.counter)
	fakeState.counter++
	return b
}

func (cs *fakeCsrf) requestToken(r *http.Request) ([]byte, error) {
	// 1. Check the HTTP header first.
	issued := r.Header.Get(headerName)

	// 2. Fall back to the POST (form) value.
	if issued == "" {
		issued = r.PostFormValue(gorillaCsrfName)
	}

	// 3. Finally, fall back to the multipart form (if set).
	if issued == "" && r.MultipartForm != nil {
		vals := r.MultipartForm.Value[gorillaCsrfName]

		if len(vals) > 0 {
			issued = vals[0]
		}
	}

	// Return nil (equivalent to empty byte slice) if no token was found
	if issued == "" {
		return nil, nil
	}

	// Decode the "issued" (pad + masked) token sent in the request. Return a
	// nil byte slice on a decoding error (this will fail upstream).
	decoded, err := base64.StdEncoding.DecodeString(issued)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

// WARNING: This is a FAKE VERSION modified from gorilla/csrf
//
// mask should returns a unique-per-request token to mitigate the BREACH attack
// as per http://breachattack.com/#mitigations
//
// WARNING: This is a FAKE VERSION modified from gorilla/csrf
//
// The token is generated by XOR'ing a one-time-pad and the base (session) CSRF
// token and returning them together as a 64-byte slice. This effectively
// randomises the token on a per-request basis without breaking multiple browser
// tabs/windows.
func maskToken(realToken []byte) string {
	otp := generateCounterBytes(tokenLength)

	// XOR the OTP with the real token to generate a masked token. Append the
	// OTP to the front of the masked token to allow unmasking in the subsequent
	// request.
	return base64.StdEncoding.EncodeToString(append(otp, xorToken(otp, realToken)...))
}

// unmask splits the issued token (one-time-pad + masked token) and returns the
// unmasked request token for comparison.
func unmaskToken(issued []byte) []byte {
	// Issued tokens are always masked and combined with the pad.
	if len(issued) != tokenLength*2 {
		return nil
	}

	// We now know the length of the byte slice.
	otp := issued[tokenLength:]
	masked := issued[:tokenLength]

	// Unmask the token by XOR'ing it against the OTP used to mask it.
	return xorToken(otp, masked)
}

// xorToken XORs tokens ([]byte) to provide unique-per-request CSRF tokens. It
// will return a masked token if the base token is XOR'ed with a one-time-pad.
// An unmasked token will be returned if a masked token is XOR'ed with the
// one-time-pad used to mask it.
func xorToken(a, b []byte) []byte {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}

	res := make([]byte, n)

	for i := 0; i < n; i++ {
		res[i] = a[i] ^ b[i]
	}

	return res
}

// contains is a helper function to check if a string exists in a slice - e.g.
// whether a HTTP method exists in a list of safe methods.
func contains(vals []string, s string) bool {
	for _, v := range vals {
		if v == s {
			return true
		}
	}

	return false
}

// compare securely (constant-time) compares the unmasked token from the request
// against the real token from the session.
func compareTokens(a, b []byte) bool {
	// This is required as subtle.ConstantTimeCompare does not check for equal
	// lengths in Go versions prior to 1.3.
	if len(a) != len(b) {
		return false
	}

	return subtle.ConstantTimeCompare(a, b) == 1
}

// cookieStore is a signed cookie session store for CSRF tokens.
type cookieStore struct {
	name     string
	maxAge   int
	secure   bool
	httpOnly bool
	path     string
	domain   string
	sc       *securecookie.SecureCookie
	sameSite int
	sz       securecookie.Serializer
}

// Get retrieves a CSRF token from the session cookie. It returns an empty token
// if decoding fails (e.g. HMAC validation fails or the named cookie doesn't exist).
func (cs *cookieStore) Get(r *http.Request) ([]byte, error) {
	// Retrieve the cookie from the request
	cookie, err := r.Cookie(cs.name)
	if err != nil {
		return nil, err
	}

	token := make([]byte, tokenLength)
	// Decode the HMAC authenticated cookie.
	err = cs.FakeDecode(cs.name, cookie.Value, &token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (cs *cookieStore) GenerateCookie(token []byte) (*http.Cookie, error) {
	// Generate an encoded cookie value with the CSRF token.
	encoded, err := cs.FakeEncode(cs.name, token)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     cs.name,
		Value:    encoded,
		MaxAge:   cs.maxAge,
		HttpOnly: cs.httpOnly,
		Secure:   cs.secure,
		SameSite: http.SameSite(cs.sameSite),
		Path:     cs.path,
		Domain:   cs.domain,
	}

	// Set the Expires field on the cookie based on the MaxAge
	// If MaxAge <= 0, we don't set the Expires attribute, making the cookie
	// session-only.
	if cs.maxAge > 0 {
		cookie.Expires = fakeState.clock.Now().Add(
			time.Duration(cs.maxAge) * time.Second)
	}
	return cookie, nil
}

// Save stores the CSRF token in the session cookie.
func (cs *cookieStore) Save(token []byte, w http.ResponseWriter) error {
	cookie, err := cs.GenerateCookie(token)
	if err != nil {
		return err
	}
	// Write the authenticated cookie to the response.
	http.SetCookie(w, cookie)

	return nil
}

func (cs *cookieStore) FakeEncode(name string, value interface{}) (string, error) {
	var err error
	var b []byte
	// 1. Serialize.
	if b, err = cs.sz.Serialize(value); err != nil {
		return "", err
	}
	b = encode(b)
	// 3. Create MAC for "name|date|value". Extra pipe to be used later.
	b = []byte(fmt.Sprintf("%s|%d|%s|", name, fakeState.clock.Now().UTC().Unix(), b))
	hashFunc := sha256.New
	mac := createMac(hmac.New(hashFunc, fakeState.hashKey), b[:len(b)-1])
	// Append mac, remove name.
	b = append(b, mac...)[len(name)+1:]
	// 4. Encode to base64.
	b = encode(b)
	// 5. skip Check length.
	// Done.
	return string(b), nil
}

// encode encodes a value using base64.
func encode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

// decode decodes a cookie using base64.
func decode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed")
	}
	return decoded[:b], nil
}

// createMac creates a message authentication code (MAC).
func createMac(h hash.Hash, value []byte) []byte {
	h.Write(value)
	return h.Sum(nil)
}

// verifyMac verifies that a message authentication code (MAC) is valid.
func verifyMac(h hash.Hash, value []byte, mac []byte) error {
	mac2 := createMac(h, value)
	// Check that both MACs are of equal length, as subtle.ConstantTimeCompare
	// does not do this prior to Go 1.4.
	if len(mac) == len(mac2) && subtle.ConstantTimeCompare(mac, mac2) == 1 {
		return nil
	}
	return fmt.Errorf("Cookie MAC invalid")
}

func (cs *cookieStore) FakeDecode(name string, value string, dst interface{}) error {

	maxLength := 4096
	// 1. Check length.
	if len(value) > maxLength {
		return fmt.Errorf("The cooke value is too long")
	}
	// 2. Decode from base64.
	b, err := decode([]byte(value))
	if err != nil {
		return err
	}
	// 3. Verify MAC. Value is "date|value|mac".
	parts := bytes.SplitN(b, []byte("|"), 3)
	if len(parts) != 3 {
		return fmt.Errorf("Cookie MAC invalid")
	}
	hashFunc := sha256.New
	h := hmac.New(hashFunc, fakeState.hashKey)
	b = append([]byte(name+"|"), b[:len(b)-len(parts[2])-1]...)
	if err = verifyMac(h, b, parts[2]); err != nil {
		return err
	}
	// 4. Verify date ranges.
	var t1 int64
	if t1, err = strconv.ParseInt(string(parts[0]), 10, 64); err != nil {
		return fmt.Errorf("cookie timestamp invalid")
	}
	t2 := fakeState.clock.Now().UTC().Unix()
	if t1 < t2-int64(cs.maxAge) {
		return fmt.Errorf("cookie expired")
	}
	// 5. Decrypt (optional).
	b, err = decode(parts[1])
	if err != nil {
		return err
	}
	// 6. Deserialize.
	if err = cs.sz.Deserialize(b, dst); err != nil {
		return err
	}
	// Done.
	return nil
}
