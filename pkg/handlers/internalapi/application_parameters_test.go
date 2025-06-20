package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/application_parameters"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

func (suite *HandlerSuite) TestApplicationParametersValidateHandler() {
	suite.Run("can lookup application values if within the customer app", func() {
		user := factory.BuildDefaultUser(suite.DB())

		req := httptest.NewRequest("POST", "/application_parameters", nil)
		req = suite.AuthenticateUserRequest(req, user)

		validationCode := "validation_code"
		testCode := "Testcode123123"
		body := internalmessages.ApplicationParameters{
			ParameterName:  &validationCode,
			ParameterValue: &testCode,
		}

		params := application_parameters.ValidateParams{
			HTTPRequest: req,
			Body:        &body,
		}
		handler := ApplicationParametersValidateHandler{suite.NewHandlerConfig()}
		response := handler.Handle(params)

		suite.Assertions.IsType(&application_parameters.ValidateOK{}, response)
	})

	suite.Run("error for unauthenticated user outside the customer app", func() {
		req := httptest.NewRequest("POST", "/application_parameters", nil)
		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		validationCode := "validation_code"
		testCode := "Testcode123123"
		body := internalmessages.ApplicationParameters{
			ParameterName:  &validationCode,
			ParameterValue: &testCode,
		}

		params := application_parameters.ValidateParams{
			HTTPRequest: req.WithContext(ctx),
			Body:        &body,
		}

		handler := ApplicationParametersValidateHandler{suite.NewHandlerConfig()}
		response := handler.Handle(params)

		suite.Assertions.IsType(&application_parameters.ValidateUnauthorized{}, response)
	})
}
