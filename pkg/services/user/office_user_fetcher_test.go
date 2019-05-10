package user

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/db/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UserServiceSuite) TestFetchOfficeUser() {
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := query.NewPopQueryBuilder(suite.DB())
	fetcher := NewOfficeUserFetcher(builder)
	pop.Debug = true
	officeUser, err := fetcher.FetchOfficeUser("id", user.ID.String())
	pop.Debug = false
	suite.NoError(err)
	suite.T().Log(officeUser)
}
