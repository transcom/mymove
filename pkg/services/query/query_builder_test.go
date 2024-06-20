package query

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type QueryBuilderSuite struct {
	*testingsuite.PopTestSuite
}

func TestUserSuite(t *testing.T) {

	ts := &QueryBuilderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
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
	builder := NewQueryBuilder()
	validUUID := uuid.Must(uuid.NewV4())
	var actualUser models.OfficeUser

	suite.Run("fetches one with filter", func() {
		// create extra record to make sure we filter
		user := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		user2 := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
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

	suite.Run("returns error on invalid column", func() {

		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, validUUID.String()),
		}
		var actualUser models.OfficeUser

		err := builder.FetchOne(suite.AppContextForTest(), &actualUser, filters)

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Zero(actualUser)
	})

	suite.Run("returns error on invalid comparator", func() {

		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", validUUID.String()),
		}
		var actualUser models.OfficeUser

		err := builder.FetchOne(suite.AppContextForTest(), &actualUser, filters)

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Zero(actualUser)
	})

	suite.Run("fails when not pointer", func() {
		var actualUser models.OfficeUser

		err := builder.FetchOne(suite.AppContextForTest(), actualUser, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Zero(actualUser)
	})

	suite.Run("fails when not pointer to struct", func() {
		var i int

		err := builder.FetchOne(suite.AppContextForTest(), &i, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
	})

}

func (suite *QueryBuilderSuite) TestFetchMany() {
	builder := NewQueryBuilder()
	var actualUsers models.OfficeUsers

	suite.Run("fetches many with uuid filter", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Create 2 users, fetch based on first ID, fetch based on second ID
		// Expected outcome: Each search returns the single matching record
		user := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		user2 := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

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

	suite.Run("fetches many with time filter", func() {
		// Under test: FetchMany function
		// Mocked: 	None
		// Set up: 	Create 2 users, fetch based createdAt timestamp being greater
		//			than that recorded for first user
		// Expected outcome: Fetch returns the single matching record

		user := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		user2 := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		filters := []services.QueryFilter{
			NewQueryFilter("created_at", greaterThan, user.CreatedAt),
		}
		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user2.ID, actualUsers[0].ID)
	})

	suite.Run("fetches many with ilike filter", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Create 2 users, search for all users with email including something.com
		// Expected outcome: Expect to find 1 of the 2 users
		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{Model: models.OfficeUser{
				Email: "email@something.com",
			}},
		}, nil)

		search := fmt.Sprintf("%%%s%%", "something.com")
		filters := []services.QueryFilter{
			NewQueryFilter("email", ilike, search),
		}
		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(actualUsers, 1)
	})

	suite.Run("fetches many with time sort desc", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Create 2 users, search with sort order by created_at, descending
		// Expected outcome: Expect that they will be returned sorted by created_at
		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		filters := []services.QueryFilter{}
		order, sort := "created_at", false
		ordering := NewQueryOrder(&order, &sort)

		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), ordering)

		suite.NoError(err)
		// check if we have at least two users
		suite.GreaterOrEqual(len(actualUsers), 2)
		suite.True(actualUsers[0].CreatedAt.After(actualUsers[1].CreatedAt), "First user created_at should be after second user created_at time")
	})

	suite.Run("fetches many with time sort asc", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Create 2 users, search with sort order by created_at, ascending
		// Expected outcome: Expect that they will be returned sorted by created_at

		filters := []services.QueryFilter{}
		order, sort := "created_at", true
		ordering := NewQueryOrder(&order, &sort)

		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), ordering)

		suite.NoError(err)
		// check if we have at least two users
		suite.GreaterOrEqual(len(actualUsers), 2)
		suite.True(actualUsers[0].CreatedAt.Before(actualUsers[1].CreatedAt), "First user created_at should be before second user created_at time")
	})

	suite.Run("fails with invalid column", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Create 2 users, search for a fake column
		// Expected outcome: Expect an error related to the fake column

		user := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		var actualUsers models.OfficeUsers
		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, user.ID.String()),
		}

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.Run("fails with invalid comparator", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Create 2 users, search for users using invalid id *
		// Expected outcome: Expect error about the invalid id

		user := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		var actualUsers models.OfficeUsers
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", user.ID.String()),
		}

		err := builder.FetchMany(suite.AppContextForTest(), &actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.Run("fails when not pointer", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Call FetchMany without providing a pointer for the return value
		// Expected outcome: Expect Data error (could be better phrased)

		var actualUsers models.OfficeUsers

		err := builder.FetchMany(suite.AppContextForTest(), actualUsers, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Empty(actualUsers)
	})

	suite.Run("fails when not pointer to slice", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Call FetchMany without providing a slice pointer for the return value
		// Expected outcome: Expect Data error
		var actualUser models.OfficeUser

		err := builder.FetchMany(suite.AppContextForTest(), &actualUser, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Empty(actualUser)
	})

	suite.Run("fails when not pointer to slice of structs", func() {
		// Under test: FetchMany function
		// Mocked: None
		// Set up: Call FetchMany without providing a slice pointer for the return value
		// Expected outcome: Expect Data error
		var intSlice []int

		err := builder.FetchMany(suite.AppContextForTest(), &intSlice, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
	})
}

func (suite *QueryBuilderSuite) TestFetchManyAssociations() {
	setupTestData := func() {
		// Create two default duty locations (with address and transportation office)
		factory.BuildDutyLocation(suite.DB(), nil, nil)
		factory.BuildDutyLocation(suite.DB(), nil, nil)
	}
	builder := NewQueryBuilder()

	suite.Run("fetches many with default associations", func() {
		setupTestData()
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

	suite.Run("fetches many with no associations", func() {
		setupTestData()
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

	suite.Run("fetches many with one explicit non-preloaded association", func() {
		setupTestData()
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

	suite.Run("fetches many with one explicit preloaded two-level association", func() {
		setupTestData()
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
	builder := NewQueryBuilder()

	suite.Run("counts with uuid filter", func() {
		user := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		user2 := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
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

	suite.Run("counts with time filter", func() {
		user := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})
		filters := []services.QueryFilter{
			NewQueryFilter("created_at", greaterThan, user.CreatedAt),
		}

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUsers{}, filters)
		suite.NoError(err)
		suite.Equal(1, count)
	})

	suite.Run("fails with invalid column", func() {
		user := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, user.ID.String()),
		}

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUsers{}, filters)

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Zero(count)
	})

	suite.Run("fails with invalid comparator", func() {
		validUUID := uuid.Must(uuid.NewV4())
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", validUUID.String()),
		}

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUsers{}, filters)

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Zero(count)
	})

	suite.Run("fails when not pointer", func() {

		count, err := builder.Count(suite.AppContextForTest(), models.OfficeUsers{}, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Zero(count)
	})

	suite.Run("fails when not pointer to slice", func() {

		count, err := builder.Count(suite.AppContextForTest(), &models.OfficeUser{}, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Zero(count)
	})

	suite.Run("fails when not pointer to slice of structs", func() {
		var intSlice []int

		count, err := builder.Count(suite.AppContextForTest(), &intSlice, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Data error encountered", err.Error())
		suite.Zero(count)
	})
}

func (suite *QueryBuilderSuite) TestCreateOne() {
	builder := NewQueryBuilder()

	suite.Run("Successfully creates a record", func() {
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())
		userInfo := models.OfficeUser{
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  "spaceman@leo.org",
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}
		verrs, err := builder.CreateOne(suite.AppContextForTest(), &userInfo)
		suite.Nil(verrs)
		suite.Nil(err)
	})

	suite.Run("Rejects input that isn't a pointer to a struct", func() {
		userInfo := models.OfficeUser{}
		_, err := builder.CreateOne(suite.AppContextForTest(), userInfo)
		suite.Error(err, "Model should be a pointer to a struct")
	})

}

func (suite *QueryBuilderSuite) TestTransaction() {
	builder := NewQueryBuilder()

	suite.Run("Successfully creates a record in a transaction", func() {
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())
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

	suite.Run("Unsuccessfully creates a record in a transaction", func() {
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())
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

	setupTestData := func() models.TransportationOffice {
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())
		userInfo := models.OfficeUser{
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  "spaceman@leo.org",
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}

		_, err := builder.CreateOne(suite.AppContextForTest(), &userInfo)
		suite.NoError(err)
		return transportationOffice
	}

	suite.Run("Successfully updates a record", func() {
		transportationOffice := setupTestData()
		officeUser := models.OfficeUser{}
		suite.NoError(suite.DB().Last(&officeUser))

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
		suite.NoError(builder.FetchOne(suite.AppContextForTest(), &record, queryFilters))
		suite.Equal("leo@spaceman.org", record.Email)
	})

	suite.Run("Successfully updates a record with an eTag for optimistic locking", func() {
		transportationOffice := setupTestData()
		officeUser := models.OfficeUser{}
		suite.NoError(suite.DB().Last(&officeUser))

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
		suite.NoError(builder.FetchOne(suite.AppContextForTest(), &record, queryFilters))
		suite.Equal("leo@spaceman.org", record.Email)
	})

	suite.Run("Reject the update when a stale eTag is used", func() {
		transportationOffice := setupTestData()
		officeUser := models.OfficeUser{}
		suite.NoError(suite.DB().Last(&officeUser))

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

	suite.Run("Rejects input that isn't a pointer to a struct", func() {
		setupTestData()
		_, err := builder.UpdateOne(suite.AppContextForTest(), models.OfficeUser{}, nil)
		suite.Error(err, "Model should be a pointer to a struct")
	})
}

func (suite *QueryBuilderSuite) TestFetchCategoricalCountsFromOneModel() {
	builder := NewQueryBuilder()
	var electronicOrder models.ElectronicOrder
	setupTestData := func() models.ElectronicOrder {

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
		return marineCorpsOrders
	}

	filters := []services.QueryFilter{
		NewQueryFilter("issuer", equals, models.IssuerArmy),
		NewQueryFilter("issuer", equals, models.IssuerCoastGuard),
		NewQueryFilter("issuer", equals, models.IssuerMarineCorps),
		NewQueryFilter("issuer", equals, models.IssuerNavy),
		NewQueryFilter("issuer", equals, models.IssuerAirForce),
	}

	suite.Run("Successfully select some category counts", func() {

		marineCorpsOrders := setupTestData()
		counts, err := builder.FetchCategoricalCountsFromOneModel(suite.AppContextForTest(), electronicOrder, filters, nil)
		suite.Nil(err)
		suite.Equal(1, counts[models.IssuerArmy])
		suite.Equal(1, counts[models.IssuerCoastGuard])
		suite.Equal(1, counts[models.IssuerMarineCorps])
		suite.Equal(1, counts[models.IssuerNavy])
		suite.Equal(1, counts[models.IssuerAirForce])

		andFilters := []services.QueryFilter{
			NewQueryFilter("updated_at", equals, marineCorpsOrders.UpdatedAt),
		}

		counts, err = builder.FetchCategoricalCountsFromOneModel(suite.AppContextForTest(), electronicOrder, andFilters, nil)
		suite.Nil(err)
		suite.Equal(1, counts[marineCorpsOrders.UpdatedAt])

	})

	suite.Run("Successfully select some counts using an AND filter", func() {
		marineCorpsOrders := setupTestData()
		andFilters := []services.QueryFilter{
			NewQueryFilter("updated_at", equals, marineCorpsOrders.UpdatedAt),
		}

		counts, err := builder.FetchCategoricalCountsFromOneModel(suite.AppContextForTest(), electronicOrder, filters, &andFilters)
		suite.Nil(err)
		suite.Equal(0, counts[models.IssuerArmy])
		suite.Equal(0, counts[models.IssuerCoastGuard])
		suite.Equal(1, counts[models.IssuerMarineCorps])
		suite.Equal(0, counts[models.IssuerNavy])
		suite.Equal(0, counts[models.IssuerAirForce])
	})

	suite.Run("Unsuccessfully select some category counts", func() {

		unsuccessfulFilter := []services.QueryFilter{NewQueryFilter("nonexisting-column", equals, "string")}

		_, err := builder.FetchCategoricalCountsFromOneModel(suite.AppContextForTest(), electronicOrder, unsuccessfulFilter, nil)
		suite.NotNil(err)

	})
}

func (suite *QueryBuilderSuite) TestDeleteOne() {
	builder := NewQueryBuilder()

	suite.Run("Successfully deletes a record", func() {
		clientCert := factory.BuildClientCert(suite.DB(), nil, nil)

		suite.NoError(builder.DeleteOne(suite.AppContextForTest(), &clientCert))

		var filters []services.QueryFilter
		queryFilters := append(filters, NewQueryFilter("id", "=", clientCert.ID.String()))
		var record models.ClientCert
		suite.Equal(sql.ErrNoRows,
			builder.FetchOne(suite.AppContextForTest(), &record, queryFilters))
	})

	suite.Run("Rejects input that isn't a pointer to a struct", func() {
		err := builder.DeleteOne(suite.AppContextForTest(), models.OfficeUser{})
		suite.Error(err, "Model should be a pointer to a struct")
	})
}
