package featureflag

import (
	"context"
	"errors"
	"os"
	"regexp"
	"strings"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/services"
)

// EnvFetcher is a way to use environment variables as feature flags
// which is basically how we used to support feature flags. Also
// helpful for local testing
type EnvFetcher struct {
	config cli.FeatureFlagConfig
}

func NewEnvFetcher(config cli.FeatureFlagConfig) (*EnvFetcher, error) {
	return &EnvFetcher{config}, nil
}

func (ef *EnvFetcher) GetFlagForUser(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (services.FeatureFlag, error) {
	if nil == appCtx.Session() {
		featureFlag := services.FeatureFlag{}
		// if getting a flag for a user, a session must exist
		return featureFlag, errors.New("Nil session when calling GetFlagForUser")
	}
	entityID := appCtx.Session().UserID.String()
	flagContext[email] = appCtx.Session().Email
	return ef.GetFlag(ctx, appCtx.Logger(), entityID, key, flagContext)
}

func (ef *EnvFetcher) GetFlag(_ context.Context, _ *zap.Logger, entityID string, key string, flagContext map[string]string) (services.FeatureFlag, error) {
	// defaults in case the flag is not found
	featureFlag := services.FeatureFlag{
		Entity:    entityID,
		Key:       key,
		Match:     false,
		Namespace: ef.config.Namespace,
	}
	re, err := regexp.Compile("[^a-zA-Z0-9]")
	if err != nil {
		return featureFlag, err
	}
	envKey := "FEATURE_FLAG_" +
		strings.ToUpper(string(re.ReplaceAll([]byte(key), []byte("_"))))
	envVal := os.Getenv(envKey)

	// default to the value of the variable
	featureFlag.Match = envVal != ""
	featureFlag.Value = envVal

	// look for another environment variable with email addresses
	emailEnvKey := envKey + "_EMAIL"
	emailEnvVal := os.Getenv(emailEnvKey)
	if emailEnvVal != "" {
		// if the email environment variable is provided, now check to
		// see if we have a match
		flagEmails := strings.Split(emailEnvVal, ",")
		email, ok := flagContext[email]
		if ok {
			for i := range flagEmails {
				if email == flagEmails[i] {
					// if the email matches, set the value to the
					// custom one provided
					emailValueEnvKey := emailEnvKey + "_VALUE"
					featureFlag.Value = os.Getenv(emailValueEnvKey)
					return featureFlag, nil
				}
			}
		}
	}
	return featureFlag, nil
}
