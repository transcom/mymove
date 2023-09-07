package featureflag

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
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

func (suite *EnvFetcherSuite) TestGetBooleanFlagForUserEnvMissing() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	flag, err := f.GetBooleanFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{Email: "foo@example.com"}),
		"missing", map[string]string{})
	suite.NoError(err)
	suite.False(flag.Match)
}

func (suite *EnvFetcherSuite) TestGetVariantFlagForUserEnvMissing() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	flag, err := f.GetVariantFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{Email: "foo@example.com"}),
		"missing", map[string]string{})
	suite.NoError(err)
	suite.Equal("missing", flag.Key)
	suite.False(flag.Match)
}

func (suite *EnvFetcherSuite) TestGetFlagForUserEnvDisabled() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	suite.T().Setenv("FEATURE_FLAG_FOO", "0")
	flag, err := f.GetBooleanFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{Email: "foo@example.com"}),
		"foo", map[string]string{})
	suite.NoError(err)
	suite.Equal("foo", flag.Key)
	suite.False(flag.Match)
}

func (suite *EnvFetcherSuite) TestGetBooleanFlagForUserEnvEnabled() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	suite.T().Setenv("FEATURE_FLAG_FOO", envVariantEnabled)
	flag, err := f.GetBooleanFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{Email: "foo@example.com"}),
		"foo", map[string]string{})
	suite.NoError(err)
	suite.Equal("foo", flag.Key)
	suite.True(flag.Match)
}

func (suite *EnvFetcherSuite) TestGetBooleanFlagEnvEnabled() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	suite.T().Setenv("FEATURE_FLAG_FOO", envVariantEnabled)
	flag, err := f.GetBooleanFlag(context.Background(),
		suite.Logger(),
		"systemEntity",
		"foo", map[string]string{})
	suite.NoError(err)
	suite.Equal("foo", flag.Key)
	suite.True(flag.Match)
}

func (suite *EnvFetcherSuite) TestGetVariantFlagEnv() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	const myVariant = "myVariant"
	suite.T().Setenv("FEATURE_FLAG_FOO", myVariant)
	flag, err := f.GetVariantFlag(context.Background(),
		suite.Logger(),
		"systemEntity",
		"foo", map[string]string{})
	suite.NoError(err)
	suite.True(flag.Match)
	suite.True(flag.IsVariant(myVariant))
}

func (suite *EnvFetcherSuite) TestGetBooleanFlagForUserEnvEmailDisabled() {
	f, err := NewEnvFetcher(cli.FeatureFlagConfig{})
	suite.NoError(err)
	disabledEmail := "foo@example.com"
	suite.T().Setenv("FEATURE_FLAG_FOO", envVariantEnabled)
	suite.T().Setenv("FEATURE_FLAG_FOO_EMAIL", disabledEmail)
	suite.T().Setenv("FEATURE_FLAG_FOO_EMAIL_VALUE", "0")

	// any user another than the disabled one should be enabled
	flag, err := f.GetBooleanFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{
			UserID: uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
			Email:  "anyother@example.com",
		}),
		"foo", map[string]string{})
	suite.NoError(err)
	suite.Equal("foo", flag.Key)
	suite.True(flag.Match)

	// but the disabled email should not be enabled
	flag, err = f.GetBooleanFlagForUser(context.Background(),
		suite.AppContextWithSessionForTest(&auth.Session{
			UserID: uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
			Email:  disabledEmail,
		}),
		"foo", map[string]string{})
	suite.NoError(err)
	suite.Equal("foo", flag.Key)
	suite.False(flag.Match)
}
