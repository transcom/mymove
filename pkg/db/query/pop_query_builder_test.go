package query

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PopQueryBuilderSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *PopQueryBuilderSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestUserSuite(t *testing.T) {

	hs := &PopQueryBuilderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
}

func (suite *PopQueryBuilderSuite) TestFetchOne() {
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewPopQueryBuilder(suite.DB())
	var actualUser models.OfficeUser

	suite.T().Run("fetches one with filter", func(t *testing.T) {
		// create extra record to make sure we filter
		user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
		filters := map[string]interface{}{
			"id": user.ID,
		}

		err := builder.FetchOne(&actualUser, filters)

		suite.NoError(err)
		suite.Equal(user.ID, actualUser.ID)

		// do the reverse to make sure we don't get the same record every time
		filters = map[string]interface{}{
			"id": user2.ID,
		}

		err = builder.FetchOne(&actualUser, filters)

		suite.NoError(err)
		suite.Equal(user2.ID, actualUser.ID)
	})

	suite.T().Run("returns error on invalid column", func(t *testing.T) {
		filters := map[string]interface{}{
			"id":          user.ID,
			"fake_column": "test@example.com",
		}

		err := builder.FetchOne(&actualUser, filters)

		suite.Error(err)
		suite.Equal("[fake_column] is not valid input", err.Error())
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		filters := map[string]interface{}{}
		var officeUser models.OfficeUser

		err := builder.FetchOne(officeUser, filters)

		suite.Error(err)
		suite.Equal("Model should be pointer to struct", err.Error())
	})

	suite.T().Run("fails when not pointer to struct", func(t *testing.T) {
		var i int
		filters := map[string]interface{}{}

		err := builder.FetchOne(&i, filters)

		suite.Error(err)
		suite.Equal("Model should be pointer to struct", err.Error())
	})

}

func (suite *PopQueryBuilderSuite) TestFetchMany() {
	// this should be stubbed out with a model that is agnostic to our code
	// similar to how the pop repo tests might work
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	email := "test@example.com"
	assertions := testdatagen.Assertions{OfficeUser: models.OfficeUser{Email: email}}
	user2 := testdatagen.MakeOfficeUser(suite.DB(), assertions)
	builder := NewPopQueryBuilder(suite.DB())
	var actualUsers models.OfficeUsers

	suite.T().Run("fetches many with filter", func(t *testing.T) {
		filters := map[string]interface{}{
			"id":    user2.ID,
			"email": email,
		}
		err := builder.FetchMany(&actualUsers, filters)

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(email, actualUsers[0].Email)
	})

	suite.T().Run("fails with invalid column", func(t *testing.T) {
		filters := map[string]interface{}{
			"id":          user2.ID,
			"fake_column": email,
		}
		err := builder.FetchMany(&actualUsers, filters)

		suite.Error(err)
		suite.Equal("[fake_column] is not valid input", err.Error())
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var officeUsers models.OfficeUsers
		filters := map[string]interface{}{}

		err := builder.FetchMany(officeUsers, filters)

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
	})

	suite.T().Run("fails when not pointer to slice", func(t *testing.T) {
		var officeUser models.OfficeUser
		filters := map[string]interface{}{}

		err := builder.FetchMany(&officeUser, filters)

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
	})

	suite.T().Run("fails when not pointer to slice of structs", func(t *testing.T) {
		var intSlice []int
		filters := map[string]interface{}{
			"id": user.ID,
		}

		err := builder.FetchMany(&intSlice, filters)

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
	})
}
