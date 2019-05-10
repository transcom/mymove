package user

import (
	"github.com/transcom/mymove/pkg/db/query"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UserServiceSuite) TestFetchOfficeUserList() {
	testdatagen.MakeDefaultOfficeUser(suite.DB())
	email := "test@example.com"
	assertions := testdatagen.Assertions{OfficeUser: models.OfficeUser{Email: email}}
	user2 := testdatagen.MakeOfficeUser(suite.DB(), assertions)
	builder := query.NewPopQueryBuilder(suite.DB())
	fetcher := NewOfficeUserListFetcher(builder)
	filters := map[string]interface{}{
		"id": user2.ID,
		"email": email,
	}

	officeUsers, err := fetcher.FetchOfficeUserList(filters)

	suite.NoError(err)
	suite.Len(officeUsers, 1)
	suite.Equal(officeUsers[0].Email, email)
}
