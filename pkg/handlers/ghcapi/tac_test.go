package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	tacop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestTacValidation() {

	suite.Run("TAC validation", func() {
		user := factory.BuildOfficeUser(suite.DB(), nil, nil)
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

			// Validate incoming payload: no body to validate

			response := handler.Handle(params)

			suite.IsType(&tacop.TacValidationOK{}, response)
			okResponse := response.(*tacop.TacValidationOK)

			// Validate outgoing payload
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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&tacop.TacValidationUnauthorized{}, response)
		payload := response.(*tacop.TacValidationUnauthorized).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&tacop.TacValidationForbidden{}, response)
		payload := response.(*tacop.TacValidationForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}
