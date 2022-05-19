package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	tacop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestTacValidation() {
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
		suite.T().Run("Successful TAC validation", func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tc.tacCode), nil)
			request = suite.AuthenticateOfficeRequest(request, user)
			params := tacop.TacValidationParams{
				HTTPRequest: request,
				Tac:         tc.tacCode,
			}
			handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
			handler := TacValidationHandler{handlerConfig}
			response := handler.Handle(params)

			suite.Assertions.IsType(&tacop.TacValidationOK{}, response)

			okResponse := response.(*tacop.TacValidationOK)
			suite.Equal(tc.isValid, *okResponse.Payload.IsValid,
				"Expected %v validation to return %v, got %v", tc.tacCode, tc.isValid, *okResponse.Payload.IsValid)
		})
	}

	suite.T().Run("Unknown user for TAC validation is unauthorized", func(t *testing.T) {
		tac := "4EVR"
		request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
		params := tacop.TacValidationParams{
			HTTPRequest: request,
			Tac:         tac,
		}
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		handler := TacValidationHandler{handlerConfig}
		response := handler.Handle(params)

		suite.Assertions.IsType(&tacop.TacValidationUnauthorized{}, response)
	})

	suite.T().Run("Unauthorized user for TAC validation is forbidden", func(t *testing.T) {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		unauthorizedUser := serviceMember.User
		tac := "4EVR"
		request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
		request = suite.AuthenticateUserRequest(request, unauthorizedUser)
		params := tacop.TacValidationParams{
			HTTPRequest: request,
			Tac:         tac,
		}
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		handler := TacValidationHandler{handlerConfig}
		response := handler.Handle(params)

		suite.Assertions.IsType(&tacop.TacValidationForbidden{}, response)
	})
}
