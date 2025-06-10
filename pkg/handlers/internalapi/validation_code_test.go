package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	vcodeops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/validation_code"
)

func (suite *HandlerSuite) TestValidationValidateHander() {
	suite.Run("can lookup application values if within the customer app", func() {
		user := factory.BuildDefaultUser(suite.DB())

		req := httptest.NewRequest("POST", "/open/validation_code", nil)
		req = suite.AuthenticateUserRequest(req, user)

		testCode := "Testcode123123"
		body := vcodeops.ValidateCodeBody{
			ValidationCode: &testCode,
		}

		params := vcodeops.ValidateCodeParams{
			HTTPRequest: req,
			Body:        body,
		}
		handler := ValidationCodeValidationCodeHandler{suite.NewHandlerConfig()}
		response := handler.Handle(params)

		suite.Assertions.IsType(&vcodeops.ValidateCodeOK{}, response)
	})

	suite.Run("error for unauthenticated user outside the customer app", func() {
		req := httptest.NewRequest("POST", "/open/validation_code", nil)
		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		testCode := "Testcode123123"
		body := vcodeops.ValidateCodeBody{
			ValidationCode: &testCode,
		}

		params := vcodeops.ValidateCodeParams{
			HTTPRequest: req.WithContext(ctx),
			Body:        body,
		}

		handler := ValidationCodeValidationCodeHandler{suite.NewHandlerConfig()}
		response := handler.Handle(params)

		suite.Assertions.IsType(&vcodeops.ValidateCodeUnauthorized{}, response)
	})
}
