package featureflag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
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

func (suite *FliptFetcherSuite) TestGetFlagForUserDisabledVariant() {
	recorder, err := recorder.New("testdata/flipt_user_disabled_variant")
	suite.NoError(err)
	defer func() {
		suite.NoError(recorder.Stop())
	}()
	client := recorder.GetDefaultClient()
	ffConfig := cli.FeatureFlagConfig{
		URL:       "http://localhost:5050",
		Token:     "token",
		Namespace: "development",
	}
	f, err := NewFliptFetcherWithClient(ffConfig, client)
	suite.NoError(err)
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
	suite.Equal(ffConfig.Namespace, flag.Namespace)
}

func (suite *FliptFetcherSuite) TestIsEnabledForUserDisabledVariant() {
	recorder, err := recorder.New("testdata/flipt_user_disabled_variant")
	suite.NoError(err)
	defer func() {
		suite.NoError(recorder.Stop())
	}()
	client := recorder.GetDefaultClient()
	ffConfig := cli.FeatureFlagConfig{
		URL:       "http://localhost:5050",
		Token:     "token",
		Namespace: "development",
	}
	f, err := NewFliptFetcherWithClient(ffConfig, client)
	suite.NoError(err)
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
	recorder, err := recorder.New("testdata/flipt_user_boolean_variant")
	suite.NoError(err)
	defer func() {
		suite.NoError(recorder.Stop())
	}()
	client := recorder.GetDefaultClient()
	ffConfig := cli.FeatureFlagConfig{
		URL:       "http://localhost:5050",
		Token:     "token",
		Namespace: "development",
	}
	f, err := NewFliptFetcherWithClient(ffConfig, client)
	suite.NoError(err)
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
	suite.Equal(ffConfig.Namespace, flag.Namespace)
}

func (suite *FliptFetcherSuite) TestIsEnabledForUserBooleanVariant() {
	recorder, err := recorder.New("testdata/flipt_user_boolean_variant")
	suite.NoError(err)
	defer func() {
		suite.NoError(recorder.Stop())
	}()
	client := recorder.GetDefaultClient()
	ffConfig := cli.FeatureFlagConfig{
		URL:       "http://localhost:5050",
		Token:     "token",
		Namespace: "development",
	}
	f, err := NewFliptFetcherWithClient(ffConfig, client)
	suite.NoError(err)
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
	recorder, err := recorder.New("testdata/flipt_user_multi_variant")
	suite.NoError(err)
	defer func() {
		suite.NoError(recorder.Stop())
	}()
	client := recorder.GetDefaultClient()
	ffConfig := cli.FeatureFlagConfig{
		URL:       "http://localhost:5050",
		Token:     "token",
		Namespace: "development",
	}
	f, err := NewFliptFetcherWithClient(ffConfig, client)
	suite.NoError(err)
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
	suite.Equal(ffConfig.Namespace, flag.Namespace)
}

func (suite *FliptFetcherSuite) TestGetFlagSystemMultiVariant() {
	recorder, err := recorder.New("testdata/flipt_system_multi_variant")
	suite.NoError(err)
	defer func() {
		suite.NoError(recorder.Stop())
	}()
	client := recorder.GetDefaultClient()
	ffConfig := cli.FeatureFlagConfig{
		URL:       "http://localhost:5050",
		Token:     "token",
		Namespace: "development",
	}
	f, err := NewFliptFetcherWithClient(ffConfig, client)
	suite.NoError(err)
	flag, err := f.GetFlag(context.Background(),
		"system",
		"multi_variant", map[string]string{})
	suite.NoError(err)
	suite.True(flag.Enabled)
	suite.Equal("multi_variant", flag.Key)
	suite.Equal("two", flag.Value)
	suite.Equal("system", flag.Entity)
	suite.Equal(ffConfig.Namespace, flag.Namespace)
}
