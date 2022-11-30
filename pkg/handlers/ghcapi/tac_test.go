package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/go-openapi/strfmt"

	tacop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestTacValidation() {

	suite.Run("TAC validation", func() {
		user := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		transportationAccountingCode := testdatagen.MakeDefaultTransportationAccountingCode(suite.DB())
		tests := []struct {
			tacCode string
			isValid bool
		}{
			{tacCode: transportationAccountingCode.TAC, isValid: true},
			{tacCode: strings.ToLower(transportationAccountingCode.TAC), isValid: true}, // test case insensitivity
			{tacCode: "4EVR", isValid: false},
		}

		for _, tc := range tests {
			request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tc.tacCode), nil)
			request = suite.AuthenticateOfficeRequest(request, user)
			params := tacop.TacValidationParams{
				HTTPRequest: request,
				Tac:         tc.tacCode,
			}
			handlerConfig := suite.HandlerConfig()
			handler := TacValidationHandler{handlerConfig}
			response := handler.Handle(params)

			suite.Assertions.IsType(&tacop.TacValidationOK{}, response)
			okResponse := response.(*tacop.TacValidationOK)
			suite.NoError(okResponse.Payload.Validate(strfmt.Default))
			suite.Equal(tc.isValid, *okResponse.Payload.IsValid,
				"Expected %v validation to return %v, got %v", tc.tacCode, tc.isValid, *okResponse.Payload.IsValid)
		}

	})

	suite.Run("Unknown user for TAC validation is unauthorized", func() {
		tac := "4EVR"
		request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
		params := tacop.TacValidationParams{
			HTTPRequest: request,
			Tac:         tac,
		}
		handlerConfig := suite.HandlerConfig()
		handler := TacValidationHandler{handlerConfig}
		response := handler.Handle(params)

		suite.Assertions.IsType(&tacop.TacValidationUnauthorized{}, response)
		payload := response.(*tacop.TacValidationUnauthorized).Payload
		suite.Nil(payload) // No payload to validate
	})

	suite.Run("Unauthorized user for TAC validation is forbidden", func() {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		unauthorizedUser := serviceMember.User
		tac := "4EVR"
		request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
		request = suite.AuthenticateUserRequest(request, unauthorizedUser)
		params := tacop.TacValidationParams{
			HTTPRequest: request,
			Tac:         tac,
		}
		handlerConfig := suite.HandlerConfig()
		handler := TacValidationHandler{handlerConfig}
		response := handler.Handle(params)

		suite.Assertions.IsType(&tacop.TacValidationForbidden{}, response)
		payload := response.(*tacop.TacValidationForbidden).Payload
		suite.Nil(payload) // No payload to validate
	})
}
