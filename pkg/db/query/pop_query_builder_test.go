package query

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

func (suite *PopQueryBuilderSuite) TestFetchMany() {
	// this should be stubbed out with a model that is agnostic to our code
	// similar to how the pop repo tests might work
	testdatagen.MakeDefaultOfficeUser(suite.DB())
	email := "test@example.com"
	assertions := testdatagen.Assertions{OfficeUser: models.OfficeUser{Email: email}}
	user2 := testdatagen.MakeOfficeUser(suite.DB(), assertions)
	builder := NewPopQueryBuilder(suite.DB())
	var officeUsers models.OfficeUsers
	suite.T().Run("fetches many with filter", func(t *testing.T) {
		filters := map[string]interface{}{
			"id":    user2.ID,
			"email": email,
		}
		err := builder.FetchMany(&officeUsers, filters)

		suite.NoError(err)
		suite.Len(officeUsers, 1)
		suite.Equal(officeUsers[0].Email, email)
	})

	suite.T().Run("fails with invalid column", func(t *testing.T) {
		filters := map[string]interface{}{
			"id":          user2.ID,
			"fake_column": email,
		}
		err := builder.FetchMany(&officeUsers, filters)

		suite.Error(err)
		suite.Equal(err.Error(), "[fake_column] is not valid input")
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var officeUsers models.OfficeUsers
		filters := map[string]interface{}{}

		err := builder.FetchMany(officeUsers, filters)

		suite.Error(err)
		suite.Equal(err.Error(), "Model should be pointer to slice of structs")
	})

	suite.T().Run("fails when not pointer to slice", func(t *testing.T) {
		var officeUser models.OfficeUser
		filters := map[string]interface{}{}

		err := builder.FetchMany(&officeUser, filters)

		suite.Error(err)
		suite.Equal(err.Error(), "Model should be pointer to slice of structs")
	})

	suite.T().Run("fails when not pointer to slice of structs", func(t *testing.T) {
		var intSlice []int
		filters := map[string]interface{}{}

		err := builder.FetchMany(&intSlice, filters)

		suite.Error(err)
		suite.Equal(err.Error(), "Model should be pointer to slice of structs")
	})
}

func (suite *PopQueryBuilderSuite) TestFetchOne() {
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewPopQueryBuilder(suite.DB())
	var officeUser models.OfficeUser

	suite.T().Run("fetches one", func(t *testing.T) {
		err := builder.FetchOne(&officeUser, "id", user.ID.String())

		suite.NoError(err)
		suite.Equal(officeUser.ID, user.ID)
	})

	suite.T().Run("returns error on invalid column", func(t *testing.T) {
		err := builder.FetchOne(&officeUser, "fake_column", user.ID.String())

		suite.Error(err)
		suite.Equal(err.Error(), "[fake_column] is not valid input")
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var officeUser models.OfficeUser

		err := builder.FetchOne(officeUser, "id", 1)

		suite.Error(err)
		suite.Equal(err.Error(), "Model should be pointer to slice of structs")
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var i int

		err := builder.FetchOne(&i, "id", 1)

		suite.Error(err)
		suite.Equal(err.Error(), "Model should be pointer to slice of structs")
	})

}

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
