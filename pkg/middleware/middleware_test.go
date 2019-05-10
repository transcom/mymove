package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

const (
	testURL string = "http://mil.example.com/static/something"
)

var (
	errStatusCode    = "incorrect status code"
	errBody          = "incorrect response body"
	errMissingHeader = "missing header"
	errInvalidHeader = "invalid header"
)

type testSuite struct {
	suite.Suite
	logger  Logger
	ok      http.HandlerFunc
	reflect http.HandlerFunc
	panic   http.HandlerFunc
}

func TestSuite(t *testing.T) {

	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	ts := &testSuite{
		logger: logger,
		ok: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
		reflect: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(body))
		}),
		panic: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("foobar"))
		}),
	}

	suite.Run(t, ts)
}

// do makes the request given the middleware (mw), handler (h), response writer (w), and request (r).
func (s *testSuite) do(mw func(inner http.Handler) http.Handler, h http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	mw(h).ServeHTTP(w, r)
}
