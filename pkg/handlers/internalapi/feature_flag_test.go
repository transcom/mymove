package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	ffop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/feature_flags"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *HandlerSuite) TestBooleanFeatureFlagUnauthenticatedHandler() {
	suite.Run("success for unauthenticated user in the customer app", func() {
		req := httptest.NewRequest("POST", "/open/feature-flags/boolean/test_ff", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		params := ffop.BooleanFeatureFlagUnauthenticatedParams{
			HTTPRequest: req.WithContext(ctx),
			Key:         "key",
			FlagContext: map[string]string{
				"thing": "one",
			},
		}

		handler := BooleanFeatureFlagsUnauthenticatedHandler{suite.NewHandlerConfig()}

		response := handler.Handle(params)

		okResponse, ok := response.(*ffop.BooleanFeatureFlagUnauthenticatedOK)
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
	})
	suite.Run("error for unauthenticated user outside the customer app", func() {
		req := httptest.NewRequest("POST", "/open/feature-flags/boolean/test_ff", nil)
		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		params := ffop.BooleanFeatureFlagUnauthenticatedParams{
			HTTPRequest: req.WithContext(ctx),
			Key:         "key",
			FlagContext: map[string]string{
				"thing": "one",
			},
		}

		handler := BooleanFeatureFlagsUnauthenticatedHandler{suite.NewHandlerConfig()}
		response := handler.Handle(params)
		res, ok := response.(*ffop.BooleanFeatureFlagUnauthenticatedUnauthorized)
		suite.True(ok)
		suite.IsType(&ffop.BooleanFeatureFlagUnauthenticatedUnauthorized{}, res)
	})
}

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

	handler := BooleanFeatureFlagsForUserHandler{suite.NewHandlerConfig()}

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

	handler := VariantFeatureFlagsForUserHandler{suite.NewHandlerConfig()}

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
