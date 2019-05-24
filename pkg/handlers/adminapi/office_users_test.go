package adminapi

import (
	"net/http/httptest"

	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/query"
	user2 "github.com/transcom/mymove/pkg/services/user"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexOfficeUsersHandler() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	req := httptest.NewRequest("GET", "/office_users", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := officeuserop.IndexOfficeUsersParams{
		HTTPRequest: req,
	}

	queryBuilder := query.NewPopQueryBuilder(suite.DB())
	handler := IndexOfficeUsersHandler{
		HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		NewQueryFilter:        query.NewQueryFilter,
		OfficeUserListFetcher: user2.NewOfficeUserListFetcher(queryBuilder),
	}
	response := handler.Handle(params)

	suite.Assertions.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
}
