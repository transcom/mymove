package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/application_parameters"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

func (suite *HandlerSuite) TestApplicationParametersValidateHandler() {
	user := factory.BuildDefaultUser(suite.DB())

	req := httptest.NewRequest("POST", "/validation_code", nil)
	req = suite.AuthenticateUserRequest(req, user)

	body := internalmessages.ValidationCode{
		ValidationCode: "TestCode123123",
	}

	params := application_parameters.ValidateParams{
		HTTPRequest: req,
		Body:        &body,
	}
	handler := ApplicationParametersValidateHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	suite.Assertions.IsType(&application_parameters.ValidateOK{}, response)
}
