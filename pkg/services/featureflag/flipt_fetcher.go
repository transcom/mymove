package featureflag

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"go.flipt.io/flipt/rpc/flipt/evaluation"
	sdk "go.flipt.io/flipt/sdk/go"
	sdkhttp "go.flipt.io/flipt/sdk/go/http"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/services"
)

const (
	applicationNameKey = "applicationName"
	isAdminUserKey     = "isAdminUser"
	isOfficeUserKey    = "isOfficeUser"
	isServiceMemberKey = "isServiceMember"
	emailKey           = "email"
)

type FliptFetcher struct {
	client sdk.SDK
	config cli.FeatureFlagConfig
}

func NewFliptFetcher(config cli.FeatureFlagConfig) (*FliptFetcher, error) {
	// For reasons I do not fully understand, trying to resolve the
	// service name in AWS ECS results in a panic(!!!) in the default
	// go resolver. Setting up this custom http client with the custom
	// resolver settings seems to work. ¯\_(ツ)_/¯
	//
	// ahobson - 2023-08-01
	client := &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			// Inspect the network connection type
			DialContext: (&net.Dialer{
				Resolver: &net.Resolver{
					PreferGo:     true,
					StrictErrors: false,
				},
			}).DialContext,
		},
	}

	return NewFliptFetcherWithClient(config, client)
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

// buildUserFlagContext creates the entityID and flag context from the
// user information in the appCtx
func buildUserFlagContext(appCtx appcontext.AppContext, flagContext map[string]string) (string, map[string]string, error) {
	if nil == appCtx.Session() {
		// if getting a flag for a user, a session must exist
		return "", flagContext, errors.New("Nil session when building user flag context")
	}

	entityID := appCtx.Session().UserID.String()

	// automatically set the context
	featureFlagContext := flagContext
	featureFlagContext[emailKey] = appCtx.Session().Email
	featureFlagContext[applicationNameKey] = string(appCtx.Session().ApplicationName)

	featureFlagContext[isAdminUserKey] = strconv.FormatBool(appCtx.Session().IsAdminUser())
	featureFlagContext[isOfficeUserKey] = strconv.FormatBool(appCtx.Session().IsOfficeUser())
	featureFlagContext[isServiceMemberKey] = strconv.FormatBool(appCtx.Session().IsServiceMember())

	// instead of sending roles, send permissions as that is more
	// granular and flexible
	permissions := appCtx.Session().Permissions
	for i := range permissions {
		featureFlagContext["permissions."+permissions[i]] = strconv.FormatBool(true)
	}

	return entityID, featureFlagContext, nil
}

func (ff *FliptFetcher) GetBooleanFlagForUser(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (services.FeatureFlag, error) {

	entityID, userFlagContext, err := buildUserFlagContext(appCtx, flagContext)
	if err != nil {
		return services.FeatureFlag{}, err
	}
	return ff.GetBooleanFlag(ctx, appCtx.Logger(), entityID, key, userFlagContext)
}

func (ff *FliptFetcher) GetBooleanFlag(ctx context.Context, logger *zap.Logger, entityID string, key string, flagContext map[string]string) (services.FeatureFlag, error) {

	// defaults in case the flag is not found
	featureFlag := services.FeatureFlag{
		Entity:    entityID,
		Key:       key,
		Match:     false,
		Namespace: ff.config.Namespace,
	}

	req := &evaluation.EvaluationRequest{
		RequestId:    uuid.Must(uuid.NewV4()).String(),
		NamespaceKey: ff.config.Namespace,
		FlagKey:      key,
		EntityId:     entityID,
		Context:      flagContext,
	}
	logger.Debug("flipt boolean evaluation request", zap.Any("req", req))
	result, err := ff.client.Evaluation().Boolean(ctx, req)

	if err != nil {
		logger.Warn("Flipt error", zap.Error(err))
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			// treat a missing feature flag as a disabled one
			logger.Warn("Feature flag variant not found",
				zap.String("key", key),
				zap.String("namespace", ff.config.Namespace))
			return featureFlag, nil
		}
		return featureFlag, err
	}

	featureFlag.Match = result.Enabled

	return featureFlag, nil
}

func (ff *FliptFetcher) GetVariantFlagForUser(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (services.FeatureFlag, error) {

	entityID, userFlagContext, err := buildUserFlagContext(appCtx, flagContext)
	if err != nil {
		return services.FeatureFlag{}, err
	}
	return ff.GetVariantFlag(ctx, appCtx.Logger(), entityID, key, userFlagContext)
}

func (ff *FliptFetcher) GetVariantFlag(ctx context.Context, logger *zap.Logger, entityID string, key string, flagContext map[string]string) (services.FeatureFlag, error) {

	// defaults in case the flag is not found
	featureFlag := services.FeatureFlag{
		Entity:    entityID,
		Key:       key,
		Match:     false,
		Variant:   "",
		Namespace: ff.config.Namespace,
	}
	req := &evaluation.EvaluationRequest{
		RequestId:    uuid.Must(uuid.NewV4()).String(),
		NamespaceKey: ff.config.Namespace,
		FlagKey:      key,
		EntityId:     entityID,
		Context:      flagContext,
	}
	logger.Debug("flipt variant evaluation request", zap.Any("req", req))
	result, err := ff.client.Evaluation().Variant(ctx, req)

	if err != nil {
		logger.Warn("Flipt error", zap.Error(err))
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			// treat a missing feature flag as a disabled one
			logger.Warn("Feature flag variant not found",
				zap.String("key", key),
				zap.String("namespace", ff.config.Namespace))
			return featureFlag, nil
		}
		return featureFlag, err
	}

	featureFlag.Match = result.Match
	featureFlag.Variant = result.VariantKey

	return featureFlag, nil
}

func (ff *FliptFetcher) GetConfig() cli.FeatureFlagConfig {
	return ff.config
}
