package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	ffop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/feature_flags"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *HandlerSuite) TestBooleanFeatureFlagForUserHandler() {
	user := factory.BuildDefaultUser(suite.DB())

	req := httptest.NewRequest("POST", "/someurl", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := ffop.BooleanFeatureFlagForUserParams{
		HTTPRequest: req,
		Key:         "key",
		FlagContext: map[string]string{
			"thing": "one",
		},
	}

	handler := BooleanFeatureFlagsForUserHandler{suite.HandlerConfig()}

	response := handler.Handle(params)

	okResponse, ok := response.(*ffop.BooleanFeatureFlagForUserOK)
	suite.True(ok)
	suite.NoError(okResponse.Payload.Validate(strfmt.Default))
	expected := services.FeatureFlag{
		Entity:    "user@example.com",
		Key:       params.Key,
		Match:     true,
		Namespace: "test",
	}
	suite.Equal(expected.Entity, *okResponse.Payload.Entity)
	suite.Equal(expected.Key, *okResponse.Payload.Key)
	suite.Equal(expected.Match, *okResponse.Payload.Match)
	suite.Equal(expected.Namespace, *okResponse.Payload.Namespace)
}

func (suite *HandlerSuite) TestVariantFeatureFlagForUserHandler() {
	user := factory.BuildDefaultUser(suite.DB())

	req := httptest.NewRequest("POST", "/someurl", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := ffop.VariantFeatureFlagForUserParams{
		HTTPRequest: req,
		Key:         "key",
		FlagContext: map[string]string{
			"thing": "one",
		},
	}

	handler := VariantFeatureFlagsForUserHandler{suite.HandlerConfig()}

	response := handler.Handle(params)

	okResponse, ok := response.(*ffop.VariantFeatureFlagForUserOK)
	suite.True(ok)
	suite.NoError(okResponse.Payload.Validate(strfmt.Default))
	expected := services.FeatureFlag{
		Entity:    "user@example.com",
		Key:       params.Key,
		Match:     true,
		Variant:   "mockVariant",
		Namespace: "test",
	}
	suite.Equal(expected.Entity, *okResponse.Payload.Entity)
	suite.Equal(expected.Key, *okResponse.Payload.Key)
	suite.Equal(expected.Match, *okResponse.Payload.Match)
	suite.Equal(expected.Variant, *okResponse.Payload.Variant)
	suite.Equal(expected.Namespace, *okResponse.Payload.Namespace)
}
