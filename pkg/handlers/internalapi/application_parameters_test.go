package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/application_parameters"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

func (suite *HandlerSuite) TestApplicationParametersValidateHandler() {
	user := factory.BuildDefaultUser(suite.DB())

	req := httptest.NewRequest("POST", "/application_parameters", nil)
	req = suite.AuthenticateUserRequest(req, user)

	body := internalmessages.ApplicationParameters{
		ParameterValue: "TestCode123123",
		ParameterName:  "validation_code",
	}

	params := application_parameters.ValidateParams{
		HTTPRequest: req,
		Body:        &body,
	}
	handler := ApplicationParametersValidateHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	suite.Assertions.IsType(&application_parameters.ValidateOK{}, response)
}
