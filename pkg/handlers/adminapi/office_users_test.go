package adminapi

import (
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"

	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/user"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexOfficeUsersHandler() {
	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		// replace this with generated UUID when filter param is built out
		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		id, _ := uuid.FromString(uuidString)
		assertions := testdatagen.Assertions{
			OfficeUser: models.OfficeUser{
				ID: id,
			},
		}
		testdatagen.MakeOfficeUser(suite.DB(), assertions)
		testdatagen.MakeDefaultOfficeUser(suite.DB())

		requestUser := testdatagen.MakeDefaultUser(suite.DB())
		req := httptest.NewRequest("GET", "/office_users", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: req,
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexOfficeUsersHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:        query.NewQueryFilter,
			OfficeUserListFetcher: user.NewOfficeUserListFetcher(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})
}
