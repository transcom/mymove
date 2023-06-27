package featureflag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// Even though the EnvFetcher doesn't really need a database, the
// PopTestSuite has so many useful helper functions, include it
type EnvFetcherSuite struct {
	*testingsuite.PopTestSuite
}

func TestEnvFetcherSuite(t *testing.T) {
	hs := &EnvFetcherSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *EnvFetcherSuite) TestGetFlagForUserEnvMissing() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	flag, err := f.GetFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{Email: "foo@example.com"}),
		"missing", map[string]string{})
	suite.NoError(err)
	suite.False(flag.Enabled)
}

func (suite *EnvFetcherSuite) TestGetFlagForUserEnvDisabled() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	suite.T().Setenv("FEATURE_FLAG_FOO", "0")
	flag, err := f.GetFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{Email: "foo@example.com"}),
		"foo", map[string]string{})
	suite.NoError(err)
	suite.False(flag.Enabled)
}

func (suite *EnvFetcherSuite) TestGetFlagForUserEnvEnabled() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	suite.T().Setenv("FEATURE_FLAG_FOO", "true")
	flag, err := f.GetFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{Email: "foo@example.com"}),
		"foo", map[string]string{})
	suite.NoError(err)
	suite.True(flag.Enabled)
}

func (suite *EnvFetcherSuite) TestGetFlagEnvEnabled() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	suite.T().Setenv("FEATURE_FLAG_FOO", "true")
	flag, err := f.GetFlag(context.Background(),
		suite.Logger(),
		"systemEntity",
		"foo", map[string]string{})
	suite.NoError(err)
	suite.True(flag.Enabled)
}
