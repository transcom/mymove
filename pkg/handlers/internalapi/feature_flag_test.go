package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	ffop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/feature_flags"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *HandlerSuite) TestFeatureFlagForUserHandler() {
	user := factory.BuildDefaultUser(suite.DB())

	req := httptest.NewRequest("GET", "/someurl", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := ffop.FeatureFlagForUserParams{
		HTTPRequest: req,
		Key:         "key",
		FlagContext: map[string]string{
			"thing": "one",
		},
	}

	handler := FeatureFlagsForUserHandler{suite.HandlerConfig()}

	response := handler.Handle(params)

	okResponse, ok := response.(*ffop.FeatureFlagForUserOK)
	suite.True(ok)
	suite.NoError(okResponse.Payload.Validate(strfmt.Default))
	expected := services.FeatureFlag{
		Entity:    "user@example.com",
		Key:       params.Key,
		Enabled:   true,
		Value:     "mock",
		Namespace: "test",
	}
	suite.Equal(expected.Entity, *okResponse.Payload.Entity)
	suite.Equal(expected.Key, *okResponse.Payload.Key)
	suite.Equal(expected.Enabled, *okResponse.Payload.Enabled)
	suite.Equal(expected.Value, *okResponse.Payload.Value)
	suite.Equal(expected.Namespace, *okResponse.Payload.Namespace)
}
