package featureflag

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"go.flipt.io/flipt/rpc/flipt"
	sdk "go.flipt.io/flipt/sdk/go"
	sdkhttp "go.flipt.io/flipt/sdk/go/http"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/services"
)

const (
	applicationName = "applicationName"
	isAdminUser     = "isAdminUser"
	isOfficeUser    = "isOfficeUser"
	isServiceMember = "isServiceMember"

	// use a convention in flipt where the name of the variant is
	// enabled or disabled for booleans
	enabledVariant  = "enabled"
	disabledVariant = "disabled"
)

type FliptFetcher struct {
	client sdk.SDK
	config cli.FeatureFlagConfig
}

func NewFliptFetcher(config cli.FeatureFlagConfig) (*FliptFetcher, error) {
	return NewFliptFetcherWithClient(config, nil)
}

func NewFliptFetcherWithClient(config cli.FeatureFlagConfig, httpClient *http.Client) (*FliptFetcher, error) {
	if config.URL == "" {
		return nil, errors.New("FliptFetcher needs a non-empty Endpoint")
	}
	sdkOptions := []sdk.Option{}
	if config.Token != "" {
		// if flipt is not exposed to the internet, we can run it
		// without authentication
		provider := sdk.StaticClientTokenProvider(config.Token)
		sdkOptions = append(sdkOptions, sdk.WithClientTokenProvider(provider))
	}
	transportOptions := []sdkhttp.Option{}
	if httpClient != nil {
		transportOptions = append(transportOptions, sdkhttp.WithHTTPClient(httpClient))
	}
	transport := sdkhttp.NewTransport(config.URL, transportOptions...)
	client := sdk.New(transport, sdkOptions...)
	return &FliptFetcher{
		client: client,
		config: config,
	}, nil
}

func (ff *FliptFetcher) GetFlagForUser(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (services.FeatureFlag, error) {
	if nil == appCtx.Session() {
		featureFlag := services.FeatureFlag{}
		// if getting a flag for a user, a session must exist
		return featureFlag, errors.New("Nil session when calling GetFlagForUser")
	}

	// use email for entityID as that makes the feature flags easier
	// to reason about
	entityID := appCtx.Session().Email

	// automatically set the context
	featureFlagContext := flagContext
	featureFlagContext[applicationName] = string(appCtx.Session().ApplicationName)

	featureFlagContext[isAdminUser] = strconv.FormatBool(appCtx.Session().IsAdminUser())
	featureFlagContext[isOfficeUser] = strconv.FormatBool(appCtx.Session().IsOfficeUser())
	featureFlagContext[isServiceMember] = strconv.FormatBool(appCtx.Session().IsServiceMember())

	// instead of sending roles, send permissions as that is more
	// granular and flexible
	permissions := appCtx.Session().Permissions
	for i := range permissions {
		featureFlagContext["permissions."+permissions[i]] = strconv.FormatBool(true)
	}

	return ff.GetFlag(ctx, entityID, key, flagContext)
}

// IsEnabledForUser is a wrapper around GetFlag for boolean flags
func (ff *FliptFetcher) IsEnabledForUser(ctx context.Context, appCtx appcontext.AppContext, key string) (bool, error) {
	flag, err := ff.GetFlagForUser(ctx, appCtx, key, map[string]string{})
	if err != nil {
		return false, err
	}
	// if the flag is not enabled at all, nothing more to do
	if !flag.Enabled {
		return false, nil
	}

	// Check for a variant specifically called 'enabled'
	return flag.Value == enabledVariant, nil
}

func (ff *FliptFetcher) GetFlag(ctx context.Context, entityID string, key string, flagContext map[string]string) (services.FeatureFlag, error) {

	featureFlag := services.FeatureFlag{}
	result, err := ff.client.Flipt().Evaluate(ctx, &flipt.EvaluationRequest{
		RequestId:    uuid.Must(uuid.NewV4()).String(),
		NamespaceKey: ff.config.Namespace,
		FlagKey:      key,
		EntityId:     entityID,
		Context:      flagContext,
	})
	if err != nil {
		return featureFlag, err
	}

	featureFlag.Entity = entityID
	featureFlag.Key = key
	featureFlag.Enabled = result.Match
	featureFlag.Value = result.Value
	featureFlag.Namespace = result.NamespaceKey

	return featureFlag, nil
}
