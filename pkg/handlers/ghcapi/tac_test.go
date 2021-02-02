package ghcapi

import (
	"fmt"
	"net/http/httptest"

	tacop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestTacValidation() {
	user := testdatagen.MakeStubbedUser(suite.DB())
	tac := "RTUC"
	request := httptest.NewRequest("GET", fmt.Sprintf("/tac/valid?tac=%s", tac), nil)
	request = suite.AuthenticateUserRequest(request, user)
	params := tacop.TacValidationParams{
		HTTPRequest: request,
		Tac:         tac,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := TacValidationHandler{context}
	response := handler.Handle(params)
	suite.Assertions.IsType(&tacop.TacValidationOK{}, response)
}
