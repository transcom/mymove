package featureflag

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// Even though the FliptFetcher doesn't really need a database, the
// PopTestSuite has so many useful helper functions, include it
type FliptFetcherSuite struct {
	*testingsuite.PopTestSuite
}

func TestFliptFetcherSuite(t *testing.T) {
	hs := &FliptFetcherSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

type fliptRequestBody struct {
	RequestID      string            `json:"requestId"`
	EntityID       string            `json:"entityId"`
	RequestContext map[string]string `json:"requestContext"`
	Match          bool              `json:"match"`
	FlagKey        string            `json:"flagKey"`
	SegmentKey     string            `json:"segmentKey"`
	Timestamp      string            `json:"timestamp"`
}

// This test uses go-vcr to record actual API requests. If you are
// making a change that affects the API requests, you will need to
// re-record them.
//
//  1. Remove the recorded requests: `rm pkg/services/featureflag/testdata/flipt_*`
//  2. Start flipt: `make feature_flag_docker`
//  3. Run the tests to re-record the requests: `go test -count 1 ./pkg/services/featureflag/...`
//  4. Stop flipt
func (suite *FliptFetcherSuite) setupRecorder(path string) *recorder.Recorder {
	recorder, err := recorder.New(path)
	suite.NoError(err)
	suite.T().Cleanup(func() {
		suite.NoError(recorder.Stop())
	})

	recorder.SetReplayableInteractions(true)

	customMatcher := func(r *http.Request, expected cassette.Request) bool {
		suite.Logger().Info("Starting custom matcher")
		if !reflect.DeepEqual(r.Header, expected.Headers) {
			suite.Logger().Info("Header mismatch",
				zap.Any("expected", expected.Headers),
				zap.Any("actual", r.Header),
			)
			return false
		}
		if (r.Body == nil || r.Body == http.NoBody) && expected.Body == "" {
			return cassette.DefaultMatcher(r, expected)
		}

		reqBody, err := io.ReadAll(r.Body)
		suite.FatalNoError(err, "failed to read request body")
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		var actualBody fliptRequestBody
		err = json.Unmarshal(reqBody, &actualBody)
		suite.FatalNoError(err, "failed to decode actual request body")

		var expectedBody fliptRequestBody
		err = json.Unmarshal([]byte(expected.Body), &expectedBody)
		suite.FatalNoError(err, "failed to decode expected request body")

		if actualBody.EntityID != expectedBody.EntityID {
			suite.Logger().Info("EntityID mismatch")
			return false
		}
		if actualBody.Match != expectedBody.Match {
			suite.Logger().Info("Match mismatch")
			return false
		}
		if actualBody.FlagKey != expectedBody.FlagKey {
			suite.Logger().Info("FlagKey mismatch")
			return false
		}
		if actualBody.SegmentKey != expectedBody.SegmentKey {
			suite.Logger().Info("SegmentKey mismatch")
			return false
		}

		if !reflect.DeepEqual(actualBody.RequestContext, expectedBody.RequestContext) {
			suite.Logger().Info("RequestContext mismatch")
			return false
		}

		return cassette.DefaultMatcher(r, expected)
	}

	recorder.SetMatcher(customMatcher)
	return recorder
}

func (suite *FliptFetcherSuite) setupFliptFetcher(path string) *FliptFetcher {
	recorder := suite.setupRecorder(path)
	client := recorder.GetDefaultClient()
	ffConfig := cli.FeatureFlagConfig{
		URL:       "http://localhost:9080",
		Namespace: "development",
	}
	f, err := NewFliptFetcherWithClient(ffConfig, client)
	suite.NoError(err)
	return f
}

func (suite *FliptFetcherSuite) TestGetFlagForUserDisabledVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_user_disabled_variant")
	fakeSession := &auth.Session{
		UserID:          uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
		Email:           "foo@example.com",
		ApplicationName: auth.MilApp,
	}
	disabledVariantKey := "disabled_variant"
	flag, err := f.GetVariantFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(fakeSession),
		disabledVariantKey, map[string]string{})
	suite.NoError(err)
	suite.Equal(disabledVariantKey, flag.Key)
	suite.False(flag.Match)
	suite.Equal(fakeSession.UserID.String(), flag.Entity)
	suite.Equal("", flag.Variant)
	suite.Equal(f.config.Namespace, flag.Namespace)
}

func (suite *FliptFetcherSuite) TestGetFlagForUserBooleanFlag() {
	f := suite.setupFliptFetcher("testdata/flipt_user_boolean_flag")
	fakeSession := &auth.Session{
		UserID:          uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
		Email:           "foo@example.com",
		ApplicationName: auth.MilApp,
	}
	const booleanKey = "boolean_flag"
	flag, err := f.GetBooleanFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(fakeSession),
		booleanKey, map[string]string{})
	suite.NoError(err)
	suite.Equal(booleanKey, flag.Key)
	suite.True(flag.Match)
}

func (suite *FliptFetcherSuite) TestGetFlagForUserMultiVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_user_multi_variant")
	fakeSession := &auth.Session{
		UserID:          uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
		Email:           "foo@example.com",
		ApplicationName: auth.MilApp,
	}
	flag, err := f.GetVariantFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(fakeSession),
		"multi_variant", map[string]string{})
	suite.NoError(err)
	suite.True(flag.Match)
	suite.Equal("multi_variant", flag.Key)
	suite.Equal("one", flag.Variant)
	suite.Equal(fakeSession.UserID.String(), flag.Entity)
	suite.Equal(f.config.Namespace, flag.Namespace)
	suite.True(flag.IsVariant("one"))
	suite.False(flag.IsVariant("two"))
}

func (suite *FliptFetcherSuite) TestGetFlagSystemMultiVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_system_multi_variant")
	flag, err := f.GetVariantFlag(context.Background(),
		suite.Logger(),
		"system",
		"multi_variant", map[string]string{})
	suite.NoError(err)
	suite.True(flag.Match)
	suite.Equal("multi_variant", flag.Key)
	suite.Equal("two", flag.Variant)
	suite.Equal("system", flag.Entity)
	suite.Equal(f.config.Namespace, flag.Namespace)
	suite.True(flag.IsVariant("two"))
	suite.False(flag.IsVariant("one"))
}
