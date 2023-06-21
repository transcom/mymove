package featureflag

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

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
		URL:       "http://localhost:5050",
		Namespace: "development",
	}
	f, err := NewFliptFetcherWithClient(ffConfig, client)
	suite.NoError(err)
	return f
}

func (suite *FliptFetcherSuite) TestGetFlagForUserDisabledVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_user_disabled_variant")
	fakeSession := &auth.Session{
		Email:           "foo@example.com",
		ApplicationName: auth.MilApp,
	}
	flag, err := f.GetFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(fakeSession),
		"disabled_variant", map[string]string{})
	suite.NoError(err)
	suite.Equal("disabled_variant", flag.Key)
	suite.False(flag.Enabled)
	suite.Equal(fakeSession.Email, flag.Entity)
	suite.Equal("", flag.Value)
	suite.Equal(f.config.Namespace, flag.Namespace)
}

func (suite *FliptFetcherSuite) TestIsEnabledForUserDisabledVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_user_disabled_variant")
	fakeSession := &auth.Session{
		Email:           "foo@example.com",
		ApplicationName: auth.MilApp,
	}
	enabled, err := f.IsEnabledForUser(context.Background(),
		suite.AppContextWithSessionForTest(fakeSession),
		"disabled_variant")
	suite.NoError(err)
	suite.False(enabled)
}

func (suite *FliptFetcherSuite) TestGetFlagForUserBooleanVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_user_boolean_variant")
	fakeSession := &auth.Session{
		Email:           "foo@example.com",
		ApplicationName: auth.MilApp,
	}
	flag, err := f.GetFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(fakeSession),
		"boolean_variant", map[string]string{})
	suite.NoError(err)
	suite.Equal("boolean_variant", flag.Key)
	suite.True(flag.Enabled)
	suite.Equal(fakeSession.Email, flag.Entity)
	suite.Equal(enabledVariant, flag.Value)
	suite.Equal(f.config.Namespace, flag.Namespace)
}

func (suite *FliptFetcherSuite) TestIsEnabledForUserBooleanVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_user_boolean_variant")
	fakeSession := &auth.Session{
		Email:           "foo@example.com",
		ApplicationName: auth.MilApp,
	}
	enabled, err := f.IsEnabledForUser(context.Background(),
		suite.AppContextWithSessionForTest(fakeSession),
		"boolean_variant")
	suite.NoError(err)
	suite.True(enabled)
}

func (suite *FliptFetcherSuite) TestGetFlagForUserMultiVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_user_multi_variant")
	fakeSession := &auth.Session{
		Email:           "foo@example.com",
		ApplicationName: auth.MilApp,
	}
	flag, err := f.GetFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(fakeSession),
		"multi_variant", map[string]string{})
	suite.NoError(err)
	suite.True(flag.Enabled)
	suite.Equal("multi_variant", flag.Key)
	suite.Equal("one", flag.Value)
	suite.Equal(fakeSession.Email, flag.Entity)
	suite.Equal(f.config.Namespace, flag.Namespace)
}

func (suite *FliptFetcherSuite) TestGetFlagSystemMultiVariant() {
	f := suite.setupFliptFetcher("testdata/flipt_system_multi_variant")
	flag, err := f.GetFlag(context.Background(),
		"system",
		"multi_variant", map[string]string{})
	suite.NoError(err)
	suite.True(flag.Enabled)
	suite.Equal("multi_variant", flag.Key)
	suite.Equal("two", flag.Value)
	suite.Equal("system", flag.Entity)
	suite.Equal(f.config.Namespace, flag.Namespace)
}
