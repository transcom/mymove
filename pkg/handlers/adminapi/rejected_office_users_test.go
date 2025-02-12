package adminapi

import (
	"github.com/transcom/mymove/pkg/factory"
	rejectedofficeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/rejected_office_users"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	rejectedofficeusers "github.com/transcom/mymove/pkg/services/rejected_office_users"
)

func (suite *HandlerSuite) TestIndexRejectedOfficeUsersHandler() {
	// test that everything is wired up
	suite.Run("rejected users result in ok response", func() {
		// building two office user with rejected status
		rejectedOfficeUsers := models.OfficeUsers{
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae}),
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})}

		params := rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRejectedOfficeUsersHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			NewQueryFilter:                query.NewQueryFilter,
			RejectedOfficeUserListFetcher: rejectedofficeusers.NewRejectedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                 pagination.NewPagination,
		}

		response := handler.Handle(params)

		// should get an ok response
		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse := response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(rejectedOfficeUsers[0].ID.String(), okResponse.Payload[0].ID.String())
	})
}
