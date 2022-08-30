package officeuser

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
	// The OfficeUserCreator needs a NotificationSender for sending user activity emails to admins.
	// This function allows us to set up a fresh mock for each test so we can check the number of calls it has.
	mockSender := mocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*notifications.UserAccountModified"),
	).Return(nil)

	return &mockSender
}

func (suite *OfficeUserServiceSuite) TestCreateOfficeUser() {
	appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})
	queryBuilder := query.NewQueryBuilder()

	loginGovUUID := uuid.Must(uuid.NewV4())
	existingUser := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
		User: models.User{
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: "spaceman+existing@leo.org",
			Active:        true,
		},
	})

	transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())
	userInfo := models.OfficeUser{
		LastName:               "Spaceman",
		FirstName:              "Leo",
		Email:                  "spaceman@leo.org",
		TransportationOfficeID: transportationOffice.ID,
		Telephone:              "312-111-1111",
		TransportationOffice:   transportationOffice,
	}

	// Happy path - creates a new User as well
	suite.Run("If the user is created successfully it should be returned", func() {
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.TransportationOffice:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(transportationOffice.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}
		fakeQueryAssociations := func(appCtx appcontext.AppContext, model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error {
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", transportationOffice.ID)}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:             fakeFetchOne,
			fakeCreateOne:            queryBuilder.CreateOne,
			fakeQueryForAssociations: fakeQueryAssociations,
		}
		mockSender := setUpMockNotificationSender()

		creator := NewOfficeUserCreator(builder, mockSender)
		officeUser, verrs, err := creator.CreateOfficeUser(appCtx, &userInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(officeUser.User)
		suite.Equal(officeUser.User.ID, *officeUser.UserID)
		suite.Equal(userInfo.Email, officeUser.User.LoginGovEmail)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	// Reuses existing user if it's already been created for an admin or service member
	suite.Run("Finds existing user by email and associates with office user", func() {
		existingUserInfo := models.OfficeUser{
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  existingUser.LoginGovEmail,
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.TransportationOffice:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(transportationOffice.ID))
			case *models.User:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(existingUser.ID))
				reflect.ValueOf(model).Elem().FieldByName("LoginGovUUID").Set(reflect.ValueOf(existingUser.LoginGovUUID))
				reflect.ValueOf(model).Elem().FieldByName("LoginGovEmail").Set(reflect.ValueOf(existingUserInfo.User.LoginGovEmail))
			}
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", transportationOffice.ID)}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: queryBuilder.CreateOne,
		}
		mockSender := setUpMockNotificationSender()

		creator := NewOfficeUserCreator(builder, mockSender)
		officeUser, verrs, err := creator.CreateOfficeUser(appCtx, &existingUserInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(officeUser.User)
		suite.Equal(officeUser.User.ID, *officeUser.UserID)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 0)
	})

	// Bad transportation office ID
	suite.Run("If we are provided a transportation office that doesn't exist, the create should fail", func() {
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return models.ErrFetchNotFound
		}
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", "b9c41d03-c730-4580-bd37-9ccf4845af6c")}
		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		creator := NewOfficeUserCreator(builder, setUpMockNotificationSender())
		_, _, err := creator.CreateOfficeUser(appCtx, &userInfo, filter)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

	// Transaction rollback on createOne validation failure
	suite.Run("CreateOne validation error should rollback transaction", func() {
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.TransportationOffice:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(transportationOffice.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}
		fakeCreateOne := func(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
			// Fail on the OfficeUser call to CreateOne but let User succeed
			switch model.(type) {
			case *models.OfficeUser:
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
		fakeQueryAssociations := func(appCtx appcontext.AppContext, model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error {
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", transportationOffice.ID)}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:             fakeFetchOne,
			fakeCreateOne:            fakeCreateOne,
			fakeQueryForAssociations: fakeQueryAssociations,
		}

		creator := NewOfficeUserCreator(builder, setUpMockNotificationSender())
		_, verrs, _ := creator.CreateOfficeUser(appCtx, &userInfo, filter)
		suite.NotNil(verrs)
		suite.Equal("violation message", verrs.Errors["errorKey"][0])
	})

	// Transaction rollback on createOne error failure
	suite.Run("CreateOne error should rollback transaction", func() {
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.TransportationOffice:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(transportationOffice.ID))
			case *models.User:
				return errors.New("User Not Found")
			}
			return nil
		}
		fakeCreateOne := func(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
			// Fail on the second createOne call with OfficeUser
			switch model.(type) {
			case *models.OfficeUser:
				return nil, errors.New("uniqueness constraint conflict")
			default:
				return nil, nil
			}
		}
		fakeQueryAssociations := func(appCtx appcontext.AppContext, model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error {
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", transportationOffice.ID)}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:             fakeFetchOne,
			fakeCreateOne:            fakeCreateOne,
			fakeQueryForAssociations: fakeQueryAssociations,
		}

		creator := NewOfficeUserCreator(builder, setUpMockNotificationSender())
		_, _, err := creator.CreateOfficeUser(appCtx, &userInfo, filter)
		suite.EqualError(err, "uniqueness constraint conflict")
	})
}
