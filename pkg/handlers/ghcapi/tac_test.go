package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	tacop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestTacValidation() {
	user := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})

	suite.T().Run("Successful Valid TAC validation", func(t *testing.T) {
		tac := "RTUC"
		request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
		request = suite.AuthenticateOfficeRequest(request, user)
		params := tacop.TacValidationParams{
			HTTPRequest: request,
			Tac:         tac,
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := TacValidationHandler{context}
		response := handler.Handle(params)

		suite.Assertions.IsType(&tacop.TacValidationOK{}, response)

		okResponse := response.(*tacop.TacValidationOK)
		suite.True(okResponse.Payload.IsValid)
	})

	suite.T().Run("Successful Invalid TAC validation", func(t *testing.T) {
		tac := "4EVR"
		request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
		request = suite.AuthenticateOfficeRequest(request, user)
		params := tacop.TacValidationParams{
			HTTPRequest: request,
			Tac:         tac,
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := TacValidationHandler{context}
		response := handler.Handle(params)

		suite.Assertions.IsType(&tacop.TacValidationOK{}, response)

		okResponse := response.(*tacop.TacValidationOK)
		suite.False(okResponse.Payload.IsValid)
	})

	suite.T().Run("Unknown user TAC validation", func(t *testing.T) {
		tac := "4EVR"
		request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
		params := tacop.TacValidationParams{
			HTTPRequest: request,
			Tac:         tac,
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := TacValidationHandler{context}
		response := handler.Handle(params)

		suite.Assertions.IsType(&tacop.TacValidationUnauthorized{}, response)
	})

	suite.T().Run("Unauthorized TAC validation", func(t *testing.T) {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		unauthorizedUser := serviceMember.User
		tac := "4EVR"
		request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
		request = suite.AuthenticateUserRequest(request, unauthorizedUser)
		params := tacop.TacValidationParams{
			HTTPRequest: request,
			Tac:         tac,
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := TacValidationHandler{context}
		response := handler.Handle(params)

		suite.Assertions.IsType(&tacop.TacValidationForbidden{}, response)
	})
}
