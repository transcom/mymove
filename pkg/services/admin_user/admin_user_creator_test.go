package adminuser

import (
	"errors"
	"reflect"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func setUpMockNotificationSender() notifications.NotificationSender {
	// The AdminUserCreator needs a NotificationSender for sending user activity emails to system admins.
	// This function allows us to set up a fresh mock for each test so we can check the number of calls it has.
	mockSender := mocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*notifications.UserAccountModified"),
	).Return(nil)

	return &mockSender
}

func (suite *AdminUserServiceSuite) TestCreateAdminUser() {
	queryBuilder := query.NewQueryBuilder()

	setupTestData := func() (models.User, models.AdminUser) {

		loginGovUUID := uuid.Must(uuid.NewV4())
		existingUser := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				LoginGovUUID:  &loginGovUUID,
				LoginGovEmail: "spaceman+existing@leo.org",
				Active:        true,
			},
		})

		organization := testdatagen.MakeDefaultOrganization(suite.DB())
		userInfo := models.AdminUser{
			LastName:       "Spaceman",
			FirstName:      "Leo",
			Email:          "spaceman@leo.org",
			OrganizationID: &organization.ID,
			Organization:   organization,
			Role:           models.SystemAdminRole,
		}
		return existingUser, userInfo
	}

	// Happy path - creates a new User as well
	suite.Run("If the user is created successfully it should be returned", func() {
		_, userInfo := setupTestData()
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		fakeFetchOne := func(appConfig appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.Organization:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(userInfo.Organization.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", userInfo.Organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: queryBuilder.CreateOne,
		}
		mockSender := setUpMockNotificationSender()

		creator := NewAdminUserCreator(builder, mockSender)
		adminUser, verrs, err := creator.CreateAdminUser(appCtx, &userInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(adminUser.User)
		suite.Equal(adminUser.User.ID, *adminUser.UserID)
		suite.Equal(userInfo.Email, adminUser.User.LoginGovEmail)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	// Reuses existing user if it's already been created for an office or service member
	suite.Run("Finds existing user by email and associates with admin user", func() {
		existingUser, userInfo := setupTestData()
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		existingUserInfo := models.AdminUser{
			LastName:       "Spaceman",
			FirstName:      "Leo",
			Email:          existingUser.LoginGovEmail,
			OrganizationID: &userInfo.Organization.ID,
			Organization:   userInfo.Organization,
			Role:           models.SystemAdminRole,
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.Organization:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(userInfo.Organization.ID))
			case *models.User:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(existingUser.ID))
				reflect.ValueOf(model).Elem().FieldByName("LoginGovUUID").Set(reflect.ValueOf(existingUser.LoginGovUUID))
				reflect.ValueOf(model).Elem().FieldByName("LoginGovEmail").Set(reflect.ValueOf(existingUserInfo.User.LoginGovEmail))
			}
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", userInfo.Organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: queryBuilder.CreateOne,
		}
		mockSender := setUpMockNotificationSender()

		creator := NewAdminUserCreator(builder, mockSender)
		adminUser, verrs, err := creator.CreateAdminUser(appCtx, &existingUserInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(adminUser.User)
		suite.Equal(adminUser.User.ID, *adminUser.UserID)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 0)
	})

	// Bad organization ID
	suite.Run("If we are provided a organization that doesn't exist, the create should fail", func() {
		_, userInfo := setupTestData()
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return models.ErrFetchNotFound
		}
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", "b9c41d03-c730-4580-bd37-9ccf4845af6c")}
		builder := &testAdminUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		creator := NewAdminUserCreator(builder, setUpMockNotificationSender())
		_, _, err := creator.CreateAdminUser(appCtx, &userInfo, filter)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())
	})

	// Transaction rollback on createOne validation failure
	suite.Run("CreateOne validation error should rollback transaction", func() {
		_, userInfo := setupTestData()
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.Organization:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(userInfo.Organization.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}
		fakeCreateOne := func(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
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
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", userInfo.Organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		creator := NewAdminUserCreator(builder, setUpMockNotificationSender())
		_, verrs, _ := creator.CreateAdminUser(appCtx, &userInfo, filter)
		suite.NotNil(verrs)
		suite.Equal("violation message", verrs.Errors["errorKey"][0])
	})

	// Transaction rollback on createOne error failure
	suite.Run("CreateOne error should rollback transaction", func() {
		_, userInfo := setupTestData()
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.Organization:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(userInfo.Organization.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}
		fakeCreateOne := func(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
			// Fail on the second createOne call with OfficeUser
			switch model.(type) {
			case *models.AdminUser:
				return nil, errors.New("uniqueness constraint conflict")
			default:
				return nil, nil
			}
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", userInfo.Organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		creator := NewAdminUserCreator(builder, setUpMockNotificationSender())
		_, _, err := creator.CreateAdminUser(appCtx, &userInfo, filter)
		suite.EqualError(err, "uniqueness constraint conflict")
	})
}
