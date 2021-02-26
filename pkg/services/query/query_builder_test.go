package query

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gobuffalo/validate/v3"

	"github.com/gobuffalo/pop/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type QueryBuilderSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestUserSuite(t *testing.T) {

	ts := &QueryBuilderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
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

		err := builder.FetchMany(&actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user2.ID, actualUsers[0].ID)

		// do the reverse to make sure we don't get the same record every time
		filters = []services.QueryFilter{
			NewQueryFilter("id", equals, user.ID.String()),
		}
		var actualUsers models.OfficeUsers

		err = builder.FetchMany(&actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

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
		err := builder.FetchMany(&actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())
		pop.Debug = false

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

		pop.Debug = true
		err := builder.FetchMany(&actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())
		pop.Debug = false

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

		pop.Debug = true
		err := builder.FetchMany(&actualUsers, filters, defaultAssociations(), defaultPagination(), ordering)
		pop.Debug = false

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

		pop.Debug = true
		err := builder.FetchMany(&actualUsers, filters, defaultAssociations(), defaultPagination(), ordering)
		pop.Debug = false

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

		err := builder.FetchMany(&actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails with invalid comparator", func(t *testing.T) {
		var actualUsers models.OfficeUsers
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", user.ID.String()),
		}

		err := builder.FetchMany(&actualUsers, filters, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var actualUsers models.OfficeUsers

		err := builder.FetchMany(actualUsers, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails when not pointer to slice", func(t *testing.T) {
		var actualUser models.OfficeUser

		err := builder.FetchMany(&actualUser, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Empty(actualUser)
	})

	suite.T().Run("fails when not pointer to slice of structs", func(t *testing.T) {
		var intSlice []int

		err := builder.FetchMany(&intSlice, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrder())

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
	})
}

func (suite *QueryBuilderSuite) TestCount() {
	// this should be stubbed out with a model that is agnostic to our code
	// similar to how the pop repo tests might work
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewQueryBuilder(suite.DB())

	suite.T().Run("counts with uuid filter", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("id", equals, user2.ID.String()),
		}

		count, err := builder.Count(&models.OfficeUsers{}, filters)

		suite.NoError(err)
		suite.Equal(1, count)

		// do the reverse to make sure we don't get the same record every time
		filters = []services.QueryFilter{
			NewQueryFilter("id", equals, user.ID.String()),
		}

		count, err = builder.Count(&models.OfficeUsers{}, filters)

		suite.NoError(err)
		suite.Equal(1, count)
	})

	suite.T().Run("counts with time filter", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("created_at", greaterThan, user.CreatedAt),
		}

		pop.Debug = true
		count, err := builder.Count(&models.OfficeUsers{}, filters)
		pop.Debug = false
		suite.NoError(err)
		suite.Equal(1, count)
	})

	suite.T().Run("fails with invalid column", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("fake_column", equals, user.ID.String()),
		}

		count, err := builder.Count(&models.OfficeUsers{}, filters)

		suite.Error(err)
		suite.Equal("[fake_column =] is not valid input", err.Error())
		suite.Zero(count)
	})

	suite.T().Run("fails with invalid comparator", func(t *testing.T) {
		filters := []services.QueryFilter{
			NewQueryFilter("id", "*", user.ID.String()),
		}

		count, err := builder.Count(&models.OfficeUsers{}, filters)

		suite.Error(err)
		suite.Equal("[id *] is not valid input", err.Error())
		suite.Zero(count)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {

		count, err := builder.Count(models.OfficeUsers{}, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Zero(count)
	})

	suite.T().Run("fails when not pointer to slice", func(t *testing.T) {

		count, err := builder.Count(&models.OfficeUser{}, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Zero(count)
	})

	suite.T().Run("fails when not pointer to slice of structs", func(t *testing.T) {
		var intSlice []int

		count, err := builder.Count(&intSlice, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Zero(count)
	})
}

func (suite *QueryBuilderSuite) TestCreateOne() {
	builder := NewQueryBuilder(suite.DB())

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
		verrs, err := builder.CreateOne(&userInfo)
		suite.Nil(verrs)
		suite.Nil(err)
	})

	suite.T().Run("Rejects input that isn't a pointer to a struct", func(t *testing.T) {
		_, err := builder.CreateOne(userInfo)
		suite.Error(err, "Model should be a pointer to a struct")
	})

}

func (suite *QueryBuilderSuite) TestTransaction() {
	builder := NewQueryBuilder(suite.DB())

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
		txErr := builder.Transaction(func(tx *pop.Connection) error {
			txBuilder := NewQueryBuilder(tx)
			verrs, err = txBuilder.CreateOne(&userInfo)

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
		txErr := builder.Transaction(func(tx *pop.Connection) error {
			txBuilder := NewQueryBuilder(tx)
			verrs, err := txBuilder.CreateOne(&testUser)
			suite.Nil(verrs)
			suite.Nil(err)

			verrs, err = txBuilder.CreateOne(&models.ReService{})
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
	builder := NewQueryBuilder(suite.DB())

	transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())
	userInfo := models.OfficeUser{
		LastName:               "Spaceman",
		FirstName:              "Leo",
		Email:                  "spaceman@leo.org",
		TransportationOfficeID: transportationOffice.ID,
		Telephone:              "312-111-1111",
		TransportationOffice:   transportationOffice,
	}

	builder.CreateOne(&userInfo)

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

		verrs, err := builder.UpdateOne(&updatedOfficeUserInfo, nil)
		suite.Nil(verrs)
		suite.Nil(err)

		var filters []services.QueryFilter
		queryFilters := append(filters, NewQueryFilter("id", "=", updatedOfficeUserInfo.ID.String()))
		var record models.OfficeUser
		builder.FetchOne(&record, queryFilters)
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
		verrs, err := builder.UpdateOne(&updatedOfficeUserInfo, &eTag)
		suite.Nil(verrs)
		suite.Nil(err)

		var filters []services.QueryFilter
		queryFilters := append(filters, NewQueryFilter("id", "=", updatedOfficeUserInfo.ID.String()))
		var record models.OfficeUser
		builder.FetchOne(&record, queryFilters)
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
		_, err := builder.UpdateOne(&updatedOfficeUserInfo, &staleETag)
		suite.NotNil(err)
	})

	suite.T().Run("Rejects input that isn't a pointer to a struct", func(t *testing.T) {
		_, err := builder.UpdateOne(models.OfficeUser{}, nil)
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

	suite.T().Run("Successfully select some counts using an AND filter", func(t *testing.T) {
		counts, err := builder.FetchCategoricalCountsFromOneModel(electronicOrder, filters, &andFilters)
		suite.Nil(err)
		suite.Equal(0, counts[models.IssuerArmy])
		suite.Equal(0, counts[models.IssuerCoastGuard])
		suite.Equal(1, counts[models.IssuerMarineCorps])
		suite.Equal(0, counts[models.IssuerNavy])
		suite.Equal(0, counts[models.IssuerAirForce])
	})

	suite.T().Run("Unsuccessfully select some category counts", func(t *testing.T) {
		unsuccessfulFilter := []services.QueryFilter{NewQueryFilter("nonexisting-column", equals, "string")}

		_, err := builder.FetchCategoricalCountsFromOneModel(electronicOrder, unsuccessfulFilter, nil)
		suite.NotNil(err)

	})
}

func (suite *QueryBuilderSuite) TestQueryAssociations() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	selectedMoveType := models.SelectedMoveTypeHHG
	sm := models.ServiceMember{
		User:      user,
		UserID:    user.ID,
		FirstName: models.StringPointer("Travis"),
		LastName:  models.StringPointer("Wayfarer"),
	}
	suite.MustSave(&sm)
	// creates access code
	code := "TEST2"
	claimedTime := time.Now()
	invalidAccessCode := models.AccessCode{
		Code:            code,
		MoveType:        selectedMoveType,
		ServiceMemberID: &sm.ID,
		ClaimedAt:       &claimedTime,
	}
	suite.MustSave(&invalidAccessCode)
	code2 := "TEST10"
	accessCode2 := models.AccessCode{
		Code:     code2,
		MoveType: selectedMoveType,
	}
	suite.MustSave(&accessCode2)

	builder := NewQueryBuilder(suite.DB())

	suite.T().Run("fetches associated data", func(t *testing.T) {

		var accessCodes models.AccessCodes
		var filters []services.QueryFilter
		queryAssociations := []services.QueryAssociation{
			NewQueryAssociation("ServiceMember"),
		}
		associations := NewQueryAssociations(queryAssociations)

		err := builder.QueryForAssociations(&accessCodes, associations, filters, defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(accessCodes, 2)
		var names []string
		for _, v := range accessCodes {
			if v.ServiceMember.FirstName != nil {
				names = append(names, *v.ServiceMember.FirstName)
			}
		}
		suite.Contains(names, *sm.FirstName)
	})

	suite.T().Run("fetches associated data with filter", func(t *testing.T) {

		var filters []services.QueryFilter
		var accessCodes models.AccessCodes
		queryFilters := append(filters, NewQueryFilter("code", "=", code))
		queryAssociations := []services.QueryAssociation{
			NewQueryAssociation("ServiceMember"),
		}
		associations := NewQueryAssociations(queryAssociations)

		err := builder.QueryForAssociations(&accessCodes, associations, queryFilters, defaultPagination(), defaultOrder())

		suite.NoError(err)
		suite.Len(accessCodes, 1)
		suite.Equal(accessCodes[0].ServiceMember.FirstName, sm.FirstName)
	})

}
