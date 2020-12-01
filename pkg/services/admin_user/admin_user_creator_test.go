package adminuser

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AdminUserServiceSuite) TestCreateAdminUser() {
	queryBuilder := query.NewQueryBuilder(suite.DB())
	organization := testdatagen.MakeDefaultOrganization(suite.DB())
	loginGovUUID := uuid.Must(uuid.NewV4())
	existingUser := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
		User: models.User{
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: "spaceman+existing@leo.org",
			Active:        true,
		},
	})

	userInfo := models.AdminUser{
		LastName:       "Spaceman",
		FirstName:      "Leo",
		Email:          "spaceman@leo.org",
		OrganizationID: &organization.ID,
		Organization:   organization,
		Role:           models.SystemAdminRole,
	}

	// Happy path
	suite.T().Run("If the user is created successfully it should be returned", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			switch model.(type) {
			case *models.Organization:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(organization.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: queryBuilder.CreateOne,
		}

		creator := NewAdminUserCreator(suite.DB(), builder)
		adminUser, verrs, err := creator.CreateAdminUser(&userInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(adminUser.User)
		suite.Equal(adminUser.User.ID, *adminUser.UserID)
		suite.Equal(userInfo.Email, adminUser.User.LoginGovEmail)
	})

	// Reuses existing user if it's already been created for an office or service member
	suite.T().Run("Finds existing user by email and associates with admin user", func(t *testing.T) {
		existingUserInfo := models.AdminUser{
			LastName:       "Spaceman",
			FirstName:      "Leo",
			Email:          existingUser.LoginGovEmail,
			OrganizationID: &organization.ID,
			Organization:   organization,
			Role:           models.SystemAdminRole,
		}

		fakeFetchOne := func(model interface{}) error {
			switch model.(type) {
			case *models.Organization:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(organization.ID))
			case *models.User:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(existingUser.ID))
				reflect.ValueOf(model).Elem().FieldByName("LoginGovUUID").Set(reflect.ValueOf(existingUser.LoginGovUUID))
				reflect.ValueOf(model).Elem().FieldByName("LoginGovEmail").Set(reflect.ValueOf(existingUserInfo.User.LoginGovEmail))
			}
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: queryBuilder.CreateOne,
		}

		creator := NewAdminUserCreator(suite.DB(), builder)
		adminUser, verrs, err := creator.CreateAdminUser(&existingUserInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(adminUser.User)
		suite.Equal(adminUser.User.ID, *adminUser.UserID)
	})

	// Bad organization ID
	suite.T().Run("If we are provided a organization that doesn't exist, the create should fail", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return models.ErrFetchNotFound
		}
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", "b9c41d03-c730-4580-bd37-9ccf4845af6c")}
		builder := &testAdminUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		creator := NewAdminUserCreator(suite.DB(), builder)
		_, _, err := creator.CreateAdminUser(&userInfo, filter)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())
	})

	// Transaction rollback on createOne validation failure
	suite.T().Run("CreateOne validation error should rollback transaction", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			switch model.(type) {
			case *models.Organization:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(organization.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}
		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			// Fail on the OfficeUser call to CreateOne but let User succeed
			switch model.(type) {
			case *models.AdminUser:
				return &validate.Errors{
					Errors: map[string][]string{
						"errorKey": {"violation message"},
					},
				}, nil
			default:
				{
					return nil, nil
				}
			}
		}
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		creator := NewAdminUserCreator(suite.DB(), builder)
		_, verrs, _ := creator.CreateAdminUser(&userInfo, filter)
		suite.NotNil(verrs)
		suite.Equal("violation message", verrs.Errors["errorKey"][0])
	})

	// Transaction rollback on createOne error failure
	suite.T().Run("CreateOne error should rollback transaction", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			switch model.(type) {
			case *models.Organization:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(organization.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}
		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			// Fail on the second createOne call with OfficeUser
			switch model.(type) {
			case *models.AdminUser:
				return nil, errors.New("uniqueness constraint conflict")
			default:
				return nil, nil
			}
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		creator := NewAdminUserCreator(suite.DB(), builder)
		_, _, err := creator.CreateAdminUser(&userInfo, filter)
		suite.EqualError(err, "uniqueness constraint conflict")
	})
}
