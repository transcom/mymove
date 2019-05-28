package adminapi

import (
	"net/http/httptest"

	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexOfficeUsersHandler() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	req := httptest.NewRequest("GET", "/office_users", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := officeuserop.IndexOfficeUsersParams{
		HTTPRequest: req,
	}

	handler := IndexOfficeUsersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
}
