package testhelpers

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/featureflag"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func MakeMobileHomeFFMap(initValue bool) map[string]bool {
	featureFlagValues := make(map[string]bool)
	featureFlagValues[featureflag.DomesticMobileHome] = true
	featureFlagValues[featureflag.DomesticMobileHomeDDPEnabled] = initValue
	featureFlagValues[featureflag.DomesticMobileHomeDOPEnabled] = initValue
	featureFlagValues[featureflag.DomesticMobileHomePackingEnabled] = initValue
	featureFlagValues[featureflag.DomesticMobileHomeUnpackingEnabled] = initValue

	return featureFlagValues
}

func MockGetFlagFunc(_ context.Context, _ *zap.Logger, entityID string, key string, _ map[string]string, mockVariant string, flagValue bool) (services.FeatureFlag, error) {
	return services.FeatureFlag{
		Entity:    entityID,
		Key:       key,
		Match:     flagValue,
		Variant:   mockVariant,
		Namespace: "test",
	}, nil
}

func SetupMockFeatureFlagFetcher(flagValue bool) *mocks.FeatureFlagFetcher {
	mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
	mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
		mock.Anything,
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(func(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (services.FeatureFlag, error) {
		return MockGetFlagFunc(ctx, appCtx.Logger(), "user@example.com", key, flagContext, "", flagValue)
	})

	mockFeatureFlagFetcher.On("GetBooleanFlag",
		mock.Anything,
		mock.AnythingOfType("*zap.Logger"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(func(ctx context.Context, logger *zap.Logger, entityID string, key string, flagContext map[string]string) (services.FeatureFlag, error) {
		return MockGetFlagFunc(ctx, nil, "user@example.com", key, flagContext, "", flagValue)
	})

	return mockFeatureFlagFetcher
}

func EditFeatureFlagMockValue(mockFeatureFlagFetcher *mocks.FeatureFlagFetcher, flagKey string, flagValue bool) {
	mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
		mock.MatchedBy(func(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) bool {
			return key == flagKey // Only mock this function call if FF key matches the one passed in
		}),
	).Return(func(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (services.FeatureFlag, error) {
		return MockGetFlagFunc(ctx, appCtx.Logger(), "user@example.com", key, flagContext, "", flagValue)
	})
}
