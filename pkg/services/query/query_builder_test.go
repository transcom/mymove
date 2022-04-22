//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate test data for use in the unit test
//RA: Creation of test data generation for unit test consumption does not present any unexpected states and conditions
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package query

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type QueryBuilderSuite struct {
	testingsuite.PopTestSuite
}

func TestUserSuite(t *testing.T) {

	ts := &QueryBuilderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func defaultPagination() services.Pagination {
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	return pagination.NewPagination(&page, &perPage)
}

func defaultOrder() services.QueryOrder {
	return NewQueryOrder(nil, nil)
}

func defaultAssociations() services.QueryAssociations {
	return NewQueryAssociations([]services.QueryAssociation{})
}

func (suite *QueryBuilderSuite) TestFetchOne() {
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewQueryBuilder()
	var actualUser models.OfficeUser

	suite.T().Run("fetches one with filter", func(t *testing.T) {
		// create extra record to make sure we filter
		user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
		filters := []services.QueryFilter{
			NewQueryFilter("id", equals, user.ID.String()),
		}

		err := builder.FetchOne(suite.AppContextForTest(), &actualUser, filters)

		suite.NoError(err)
		suite.Equal(user.ID, actualUser.ID)

		// do the reverse to make sure we don't get the same record every time
		filters = []services.QueryFilter{
			NewQueryFilter("id", equals, user2.ID.String()),
		}

		err = builder.FetchOne(suite.AppContextForTest(), &actualUser, filters)

		suite.NoError(err)
		suite.Equal(user2.ID, actualUser.ID)
	})

	suite.T().Run("returns error on invalid column", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, user.ID.String()),
		}
		var actualUser models.OfficeUser

		err := builder.FetchOne(suite.AppContextForTest(), &actualUser, filters)

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Zero(actualUser)
	})

	suite.T().Run("returns error on invalid comparator", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", user.ID.String()),
		}
		var actualUser models.OfficeUser

		err := builder.FetchOne(suite.AppContextForTest(), &actualUser, filters)

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Zero(actualUser)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var actualUser models.OfficeUser

		err := builder.FetchOne(suite.AppContextForTest(), actualUser, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Zero(actualUser)
	})

	suite.T().Run("fails when not pointer to struct", func(t *testing.T) {
		var i int

		err := builder.FetchOne(suite.AppContextForTest(), &i, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
	})

}

func (suite *QueryBuilderSuite) TestFetchMany() {
	// this should be stubbed out with a model that is agnostic to our code
	// similar to how the pop repo tests might work
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewQueryBuilder()
	var actualUsers models.OfficeUsers

	suite.T().Run("fetches many with uuid filter", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("id", equals, user2.ID.String()),
		}

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user2.ID, actualUsers[0].ID)

		// do the reverse to make sure we don't get the same record every time
		filters = []services.QueryFilter{
			NewQueryFilter("id", equals, user.ID.String()),
		}
		var actualUsers models.OfficeUsers

		err = builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user.ID, actualUsers[0].ID)
	})

	suite.T().Run("fetches many with time filter", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("created_at", greaterThan, user.CreatedAt),
		}
		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user2.ID, actualUsers[0].ID)
	})

	suite.T().Run("fetches many with ilike filter", func(t *testing.T) {
		search := fmt.Sprintf("%%%s%%", "example.com")
		filters := []services.QueryFilter{
			NewQueryFilter("email", ilike, search),
		}
		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(actualUsers, 4)
	})

	suite.T().Run("fetches many with time sort desc", func(t *testing.T) {
		filters := []services.QueryFilter{}
		order, sort := "created_at", false
		ordering := NewQueryOrder(&order, &sort)

		testdatagen.MakeDefaultOfficeUser(suite.DB())
		testdatagen.MakeDefaultOfficeUser(suite.DB())

		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), ordering)

		suite.NoError(err)
		// check if we have at least two users
		suite.GreaterOrEqual(len(actualUsers), 2)
		suite.True(actualUsers[0].CreatedAt.After(actualUsers[1].CreatedAt), "First user created_at should be after second user created_at time")
	})

	suite.T().Run("fetches many with time sort asc", func(t *testing.T) {
		filters := []services.QueryFilter{}
		order, sort := "created_at", true
		ordering := NewQueryOrder(&order, &sort)

		testdatagen.MakeDefaultOfficeUser(suite.DB())
		testdatagen.MakeDefaultOfficeUser(suite.DB())

		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), ordering)

		suite.NoError(err)
		// check if we have at least two users
		suite.GreaterOrEqual(len(actualUsers), 2)
		suite.True(actualUsers[0].CreatedAt.Before(actualUsers[1].CreatedAt), "First user created_at should be before second user created_at time")
	})

	suite.T().Run("fails with invalid column", func(t *testing.T) {
		var actualUsers models.OfficeUsers
		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, user.ID.String()),
		}

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails with invalid comparator", func(t *testing.T) {
		var actualUsers models.OfficeUsers
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", user.ID.String()),
		}

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), actualUsers, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails when not pointer to slice", func(t *testing.T) {
		var actualUser models.OfficeUser

		err := builder.FetchMany(suite.AppContextForTest(), &actualUser, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Empty(actualUser)
	})

	suite.T().Run("fails when not pointer to slice of structs", func(t *testing.T) {
		var intSlice []int

		err := builder.FetchMany(suite.AppContextForTest(), &intSlice, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
	})
}

func (suite *QueryBuilderSuite) TestFetchManyAssociations() {
	// Create two default duty locations (with address and transportation office)
	testdatagen.MakeDefaultDutyLocation(suite.DB())
	testdatagen.MakeDefaultDutyLocation(suite.DB())
	builder := NewQueryBuilder()

	suite.T().Run("fetches many with default associations", func(t *testing.T) {
		var dutyLocations models.DutyLocations
		err := builder.FetchMany(suite.AppContextForTest(), &dutyLocations, nil, defaultAssociations(), nil, nil)
		suite.NoError(err)
		suite.Len(dutyLocations, 2)

		// Make sure every record has an address and transportation office loaded
		for _, dutyLocation := range dutyLocations {
			suite.NotEqual(uuid.Nil, dutyLocation.Address.ID)
			suite.NotEqual(uuid.Nil, dutyLocation.TransportationOffice.ID)
		}
	})

	suite.T().Run("fetches many with no associations", func(t *testing.T) {
		var dutyLocations models.DutyLocations
		err := builder.FetchMany(suite.AppContextForTest(), &dutyLocations, nil, nil, nil, nil)
		suite.NoError(err)
		suite.Len(dutyLocations, 2)

		// Make sure every record has no address or transportation office loaded
		for _, dutyLocation := range dutyLocations {
			suite.Equal(uuid.Nil, dutyLocation.Address.ID)
			suite.Equal(uuid.Nil, dutyLocation.TransportationOffice.ID)
		}
	})

	suite.T().Run("fetches many with one explicit non-preloaded association", func(t *testing.T) {
		var dutyLocations models.DutyLocations
		associations := NewQueryAssociations([]services.QueryAssociation{
			NewQueryAssociation("Address"),
		})

		err := builder.FetchMany(suite.AppContextForTest(), &dutyLocations, nil, associations, nil, nil)
		suite.NoError(err)
		suite.Len(dutyLocations, 2)

		// Make sure every record has an address loaded but not a transportation office
		for _, dutyLocation := range dutyLocations {
			suite.NotEqual(uuid.Nil, dutyLocation.Address.ID)
			suite.Equal(uuid.Nil, dutyLocation.TransportationOffice.ID)
		}
	})

	suite.T().Run("fetches many with one explicit preloaded two-level association", func(t *testing.T) {
		var dutyLocations models.DutyLocations
		associations := NewQueryAssociationsPreload([]services.QueryAssociation{
			NewQueryAssociation("TransportationOffice.Address"),
		})

		err := builder.FetchMany(suite.AppContextForTest(), &dutyLocations, nil, associations, nil, nil)
		suite.NoError(err)
		suite.Len(dutyLocations, 2)

		// Make sure every record does not have an address loaded but does have a transportation office and
		// its address loaded
		for _, dutyLocation := range dutyLocations {
			suite.Equal(uuid.Nil, dutyLocation.Address.ID)
			suite.NotEqual(uuid.Nil, dutyLocation.TransportationOffice.ID)
			suite.NotEqual(uuid.Nil, dutyLocation.TransportationOffice.Address.ID)
		}
	})
}

func (suite *QueryBuilderSuite) TestCount() {
	// this should be stubbed out with a model that is agnostic to our code
	// similar to how the pop repo tests might work
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewQueryBuilder()

	suite.T().Run("counts with uuid filter", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("id", equals, user2.ID.String()),
		}

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUsers{}, filters)

		suite.NoError(err)
		suite.Equal(1, count)

		// do the reverse to make sure we don't get the same record every time
		filters = []services.QueryFilter{
			NewQueryFilter("id", equals, user.ID.String()),
		}

		count, err = builder.Count(suite.AppContextForTest(), &models.OfficeUsers{}, filters)

		suite.NoError(err)
		suite.Equal(1, count)
	})

	suite.T().Run("counts with time filter", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("created_at", greaterThan, user.CreatedAt),
		}

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUsers{}, filters)
		suite.NoError(err)
		suite.Equal(1, count)
	})

	suite.T().Run("fails with invalid column", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, user.ID.String()),
		}

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUsers{}, filters)

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Zero(count)
	})

	suite.T().Run("fails with invalid comparator", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", user.ID.String()),
		}

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUsers{}, filters)

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Zero(count)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {

		count, err := builder.Count(suite.AppContextForTest(), models.OfficeUsers{}, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Zero(count)
	})

	suite.T().Run("fails when not pointer to slice", func(t *testing.T) {

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUser{}, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Zero(count)
	})

	suite.T().Run("fails when not pointer to slice of structs", func(t *testing.T) {
		var intSlice []int

		count, err := builder.Count(suite.AppContextForTest(), &intSlice, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Zero(count)
	})
}

func (suite *QueryBuilderSuite) TestCreateOne() {
	builder := NewQueryBuilder()

	transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())
	userInfo := models.OfficeUser{
		LastName:               "Spaceman",
		FirstName:              "Leo",
		Email:                  "spaceman@leo.org",
		TransportationOfficeID: transportationOffice.ID,
		Telephone:              "312-111-1111",
		TransportationOffice:   transportationOffice,
	}

	suite.T().Run("Successfully creates a record", func(t *testing.T) {
		verrs, err := builder.CreateOne(suite.AppContextForTest(), &userInfo)
		suite.Nil(verrs)
		suite.Nil(err)
	})

	suite.T().Run("Rejects input that isn't a pointer to a struct", func(t *testing.T) {
		_, err := builder.CreateOne(suite.AppContextForTest(), userInfo)
		suite.Error(err, "Model should be a pointer to a struct")
	})

}

func (suite *QueryBuilderSuite) TestTransaction() {
	builder := NewQueryBuilder()

	transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())

	suite.T().Run("Successfully creates a record in a transaction", func(t *testing.T) {
		userInfo := models.OfficeUser{
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  "spacemannnnnnnnn@leo.org",
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}

		var verrs *validate.Errors
		var err error
		txErr := suite.AppContextForTest().NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			verrs, err = builder.CreateOne(txnAppCtx, &userInfo)

			return nil
		})
		suite.Nil(txErr)

		suite.Nil(verrs)
		suite.Nil(err)
		suite.NotZero(userInfo.ID)
	})

	suite.T().Run("Unsuccessfully creates a record in a transaction", func(t *testing.T) {
		testUser := models.OfficeUser{
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  "testman@test.com",
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}

		// rollback intentionally with a successful create and unsuccessful create
		txErr := suite.AppContextForTest().NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			verrs, err := builder.CreateOne(txnAppCtx, &testUser)
			suite.Nil(verrs)
			suite.Nil(err)

			verrs, err = builder.CreateOne(txnAppCtx, &models.ReService{})
			suite.NotNil(verrs)
			suite.Nil(err)

			return errors.New("testing")
		})
		suite.NotNil(txErr)

		err := suite.DB().Find(&models.OfficeUser{}, testUser.ID)
		suite.NotNil(err)
	})
}

func (suite *QueryBuilderSuite) TestUpdateOne() {
	builder := NewQueryBuilder()

	transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())
	userInfo := models.OfficeUser{
		LastName:               "Spaceman",
		FirstName:              "Leo",
		Email:                  "spaceman@leo.org",
		TransportationOfficeID: transportationOffice.ID,
		Telephone:              "312-111-1111",
		TransportationOffice:   transportationOffice,
	}

	builder.CreateOne(suite.AppContextForTest(), &userInfo)

	suite.T().Run("Successfully updates a record", func(t *testing.T) {
		officeUser := models.OfficeUser{}
		suite.DB().Last(&officeUser)

		updatedOfficeUserInfo := models.OfficeUser{
			ID:                     officeUser.ID,
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  "leo@spaceman.org", // updated the email
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}

		verrs, err := builder.UpdateOne(suite.AppContextForTest(), &updatedOfficeUserInfo, nil)
		suite.Nil(verrs)
		suite.Nil(err)

		var filters []services.QueryFilter
		queryFilters := append(filters, NewQueryFilter("id", "=", updatedOfficeUserInfo.ID.String()))
		var record models.OfficeUser
		builder.FetchOne(suite.AppContextForTest(), &record, queryFilters)
		suite.Equal("leo@spaceman.org", record.Email)
	})

	suite.T().Run("Successfully updates a record with an eTag for optimistic locking", func(t *testing.T) {
		officeUser := models.OfficeUser{}
		suite.DB().Last(&officeUser)

		updatedOfficeUserInfo := models.OfficeUser{
			ID:                     officeUser.ID,
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  "leo@spaceman.org", // updated the email
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}

		eTag := etag.GenerateEtag(officeUser.UpdatedAt)
		verrs, err := builder.UpdateOne(suite.AppContextForTest(), &updatedOfficeUserInfo, &eTag)
		suite.Nil(verrs)
		suite.Nil(err)

		var filters []services.QueryFilter
		queryFilters := append(filters, NewQueryFilter("id", "=", updatedOfficeUserInfo.ID.String()))
		var record models.OfficeUser
		builder.FetchOne(suite.AppContextForTest(), &record, queryFilters)
		suite.Equal("leo@spaceman.org", record.Email)
	})

	suite.T().Run("Reject the update when a stale eTag is used", func(t *testing.T) {
		officeUser := models.OfficeUser{}
		suite.DB().Last(&officeUser)

		updatedOfficeUserInfo := models.OfficeUser{
			ID:                     officeUser.ID,
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  "leo@spaceman.org", // updated the email
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}

		staleETag := etag.GenerateEtag(time.Now())
		_, err := builder.UpdateOne(suite.AppContextForTest(), &updatedOfficeUserInfo, &staleETag)
		suite.NotNil(err)
	})

	suite.T().Run("Rejects input that isn't a pointer to a struct", func(t *testing.T) {
		_, err := builder.UpdateOne(suite.AppContextForTest(), models.OfficeUser{}, nil)
		suite.Error(err, "Model should be a pointer to a struct")
	})
}

func (suite *QueryBuilderSuite) TestFetchCategoricalCountsFromOneModel() {
	builder := NewQueryBuilder()
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
		counts, err := builder.FetchCategoricalCountsFromOneModel(suite.AppContextForTest(), electronicOrder, filters, nil)
		suite.Nil(err)
		suite.Equal(1, counts[models.IssuerArmy])
		suite.Equal(1, counts[models.IssuerCoastGuard])
		suite.Equal(1, counts[models.IssuerMarineCorps])
		suite.Equal(1, counts[models.IssuerNavy])
		suite.Equal(1, counts[models.IssuerAirForce])

		counts, err = builder.FetchCategoricalCountsFromOneModel(suite.AppContextForTest(), electronicOrder, andFilters, nil)
		suite.Nil(err)
		suite.Equal(1, counts[marineCorpsOrders.UpdatedAt])

	})

	suite.T().Run("Successfully select some counts using an AND filter", func(t *testing.T) {
		counts, err := builder.FetchCategoricalCountsFromOneModel(suite.AppContextForTest(), electronicOrder, filters, &andFilters)
		suite.Nil(err)
		suite.Equal(0, counts[models.IssuerArmy])
		suite.Equal(0, counts[models.IssuerCoastGuard])
		suite.Equal(1, counts[models.IssuerMarineCorps])
		suite.Equal(0, counts[models.IssuerNavy])
		suite.Equal(0, counts[models.IssuerAirForce])
	})

	suite.T().Run("Unsuccessfully select some category counts", func(t *testing.T) {
		unsuccessfulFilter := []services.QueryFilter{NewQueryFilter("nonexisting-column", equals, "string")}

		_, err := builder.FetchCategoricalCountsFromOneModel(suite.AppContextForTest(), electronicOrder, unsuccessfulFilter, nil)
		suite.NotNil(err)

	})
}
