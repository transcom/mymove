package query

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type QueryBuilderSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *QueryBuilderSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestUserSuite(t *testing.T) {

	hs := &QueryBuilderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
}

func (suite *QueryBuilderSuite) TestFetchOne() {
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewQueryBuilder(suite.DB())
	var actualUser models.OfficeUser

	suite.T().Run("fetches one with filter", func(t *testing.T) {
		// create extra record to make sure we filter
		user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
		filters := []services.QueryFilter{
			NewQueryFilter("id", equals, user.ID.String()),
		}

		err := builder.FetchOne(&actualUser, filters)

		suite.NoError(err)
		suite.Equal(user.ID, actualUser.ID)

		// do the reverse to make sure we don't get the same record every time
		filters = []services.QueryFilter{
			NewQueryFilter("id", equals, user2.ID.String()),
		}

		err = builder.FetchOne(&actualUser, filters)

		suite.NoError(err)
		suite.Equal(user2.ID, actualUser.ID)
	})

	suite.T().Run("returns error on invalid column", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, user.ID.String()),
		}
		var actualUser models.OfficeUser

		err := builder.FetchOne(&actualUser, filters)

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Zero(actualUser)
	})

	suite.T().Run("returns error on invalid comparator", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", user.ID.String()),
		}
		var actualUser models.OfficeUser

		err := builder.FetchOne(&actualUser, filters)

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Zero(actualUser)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var actualUser models.OfficeUser

		err := builder.FetchOne(actualUser, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Model should be pointer to struct", err.Error())
		suite.Zero(actualUser)
	})

	suite.T().Run("fails when not pointer to struct", func(t *testing.T) {
		var i int

		err := builder.FetchOne(&i, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Model should be pointer to struct", err.Error())
	})

}

func (suite *QueryBuilderSuite) TestFetchMany() {
	// this should be stubbed out with a model that is agnostic to our code
	// similar to how the pop repo tests might work
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewQueryBuilder(suite.DB())
	var actualUsers models.OfficeUsers

	suite.T().Run("fetches many with uuid filter", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("id", equals, user2.ID.String()),
		}

		err := builder.FetchMany(&actualUsers, filters)

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user2.ID, actualUsers[0].ID)

		// do the reverse to make sure we don't get the same record every time
		filters = []services.QueryFilter{
			NewQueryFilter("id", equals, user.ID.String()),
		}
		var actualUsers models.OfficeUsers

		err = builder.FetchMany(&actualUsers, filters)

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user.ID, actualUsers[0].ID)
	})

	suite.T().Run("fetches many with time filter", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("created_at", greaterThan, user.CreatedAt),
		}
		var actualUsers models.OfficeUsers

		pop.Debug = true
		err := builder.FetchMany(&actualUsers, filters)
		pop.Debug = false

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user2.ID, actualUsers[0].ID)
	})

	suite.T().Run("fails with invalid column", func(t *testing.T) {
		var actualUsers models.OfficeUsers
		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, user.ID.String()),
		}

		err := builder.FetchMany(&actualUsers, filters)

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails with invalid comparator", func(t *testing.T) {
		var actualUsers models.OfficeUsers
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", user.ID.String()),
		}

		err := builder.FetchMany(&actualUsers, filters)

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var actualUsers models.OfficeUsers

		err := builder.FetchMany(actualUsers, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails when not pointer to slice", func(t *testing.T) {
		var actualUser models.OfficeUser

		err := builder.FetchMany(&actualUser, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Empty(actualUser)
	})

	suite.T().Run("fails when not pointer to slice of structs", func(t *testing.T) {
		var intSlice []int

		err := builder.FetchMany(&intSlice, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
	})
}

func (suite *QueryBuilderSuite) TestCreateOne() {
	builder := NewQueryBuilder(suite.DB())

	transportationOffice := testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{})
	userInfo := models.OfficeUser{
		LastName:               "Spaceman",
		FirstName:              "Leo",
		Email:                  "spaceman@leo.org",
		TransportationOfficeID: transportationOffice.ID,
		Telephone:              "312-111-1111",
		TransportationOffice:   transportationOffice,
	}

	suite.T().Run("Successfully creates a record", func(t *testing.T) {
		verrs, err := builder.CreateOne(&userInfo)
		suite.Nil(verrs)
		suite.Nil(err)
	})

	suite.T().Run("Rejects input that isn't a pointer to a struct", func(t *testing.T) {
		_, err := builder.CreateOne(userInfo)
		suite.Error(err, "Model should be a pointer to a struct")
	})

}

func (suite *QueryBuilderSuite) TestFetchCategoricalCountsFromOneModel() {
	builder := NewQueryBuilder(suite.DB())
	var electronicOrder models.ElectronicOrder
	ordersAssertion := testdatagen.Assertions{
		ElectronicOrder: models.ElectronicOrder{},
	}
	// Let's make a some electronic orders to test this with
	ordersAssertion.ElectronicOrder.Issuer = models.IssuerNavy
	ordersAssertion.ElectronicOrder.OrdersNumber = "8675308"
	testdatagen.MakeElectronicOrder(suite.DB(), ordersAssertion)

	ordersAssertion.ElectronicOrder.Issuer = models.IssuerArmy
	ordersAssertion.ElectronicOrder.OrdersNumber = "8675310"
	testdatagen.MakeElectronicOrder(suite.DB(), ordersAssertion)

	ordersAssertion.ElectronicOrder.Issuer = models.IssuerMarineCorps
	ordersAssertion.ElectronicOrder.OrdersNumber = "8675311"
	marineCorpsOrders := testdatagen.MakeElectronicOrder(suite.DB(), ordersAssertion)

	ordersAssertion.ElectronicOrder.Issuer = models.IssuerAirForce
	ordersAssertion.ElectronicOrder.OrdersNumber = "8675312"
	testdatagen.MakeElectronicOrder(suite.DB(), ordersAssertion)

	ordersAssertion.ElectronicOrder.Issuer = models.IssuerCoastGuard
	ordersAssertion.ElectronicOrder.OrdersNumber = "8675313"
	testdatagen.MakeElectronicOrder(suite.DB(), ordersAssertion)

	filters := []services.QueryFilter{
		NewQueryFilter("issuer", equals, models.IssuerArmy),
		NewQueryFilter("issuer", equals, models.IssuerCoastGuard),
		NewQueryFilter("issuer", equals, models.IssuerMarineCorps),
		NewQueryFilter("issuer", equals, models.IssuerNavy),
		NewQueryFilter("issuer", equals, models.IssuerAirForce),
	}

	andFilters := []services.QueryFilter{
		NewQueryFilter("updated_at", equals, marineCorpsOrders.UpdatedAt),
	}

	suite.T().Run("Successfully select some category counts", func(t *testing.T) {
		counts, err := builder.FetchCategoricalCountsFromOneModel(electronicOrder, filters, nil)
		suite.Nil(err)
		suite.Equal(1, counts[models.IssuerArmy])
		suite.Equal(1, counts[models.IssuerCoastGuard])
		suite.Equal(1, counts[models.IssuerMarineCorps])
		suite.Equal(1, counts[models.IssuerNavy])
		suite.Equal(1, counts[models.IssuerAirForce])

		counts, err = builder.FetchCategoricalCountsFromOneModel(electronicOrder, andFilters, nil)
		suite.Nil(err)
		suite.Equal(1, counts[marineCorpsOrders.UpdatedAt])

	})

	suite.T().Run("Successfully select some counts using and AN filter", func(t *testing.T) {
		counts, err := builder.FetchCategoricalCountsFromOneModel(electronicOrder, filters, &andFilters)
		suite.Nil(err)
		suite.Equal(0, counts[models.IssuerArmy])
		suite.Equal(0, counts[models.IssuerCoastGuard])
		suite.Equal(1, counts[models.IssuerMarineCorps])
		suite.Equal(0, counts[models.IssuerNavy])
		suite.Equal(0, counts[models.IssuerAirForce])
	})

}
