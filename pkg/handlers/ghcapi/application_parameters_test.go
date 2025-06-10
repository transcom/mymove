package ghcapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/application_parameters"
)

func (suite *HandlerSuite) TestApplicationParametersValidateHandler() {
	user := factory.BuildDefaultUser(suite.DB())

	req := httptest.NewRequest("GET", "/application_parameters", nil)
	req = suite.AuthenticateUserRequest(req, user)

	handlerConfig := suite.NewHandlerConfig()
	handler := ApplicationParametersParamHandler{
		HandlerConfig: handlerConfig,
	}

	params := application_parameters.GetParamParams{
		HTTPRequest:   req,
		ParameterName: "standaloneCrateCap",
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&application_parameters.GetParamOK{}, response)
}
