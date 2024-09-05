package officeuser

import (
	"errors"
	"reflect"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
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

	setupTestData := func() (models.User, models.OfficeUser) {
		oID := uuid.Must(uuid.NewV4())
		existingUser := factory.BuildUser(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					OktaID:    oID.String(),
					OktaEmail: "spaceman+existing@leo.org",
					Active:    true,
				},
			}}, nil)
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					ID: uuid.Must(uuid.NewV4()),
				},
			}}, nil)
		officeUser := models.OfficeUser{
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  "spaceman@leo.org",
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
			EDIPI:                  models.StringPointer("1234567890"),
			OtherUniqueID:          models.StringPointer("1234567890"),
		}
		return existingUser, officeUser
	}

	// Happy path - creates a new User as well
	suite.Run("If the user is created successfully it should be returned", func() {
		_, userInfo := setupTestData()
		transportationOffice := userInfo.TransportationOffice

		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		queryBuilder := query.NewQueryBuilder()

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
		suite.Equal(userInfo.Email, officeUser.User.OktaEmail)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	// Reuses existing user if it's already been created for an admin or service member
	suite.Run("Finds existing user by email and associates with office user", func() {
		existingUser, userInfo := setupTestData()
		transportationOffice := userInfo.TransportationOffice
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		queryBuilder := query.NewQueryBuilder()

		existingUserInfo := models.OfficeUser{
			LastName:               "Spaceman",
			FirstName:              "Leo",
			Email:                  existingUser.OktaEmail,
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
				reflect.ValueOf(model).Elem().FieldByName("OktaID").Set(reflect.ValueOf(existingUser.OktaID))
				reflect.ValueOf(model).Elem().FieldByName("OktaEmail").Set(reflect.ValueOf(existingUserInfo.User.OktaEmail))
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

	suite.Run("Updates previously rejected office user instead of create", func() {
		rejectedStatus := models.OfficeUserStatusREJECTED
		requestedStatus := models.OfficeUserStatusREQUESTED

		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					ID: uuid.Must(uuid.NewV4()),
				},
			}}, nil)

		user := factory.BuildUser(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					ID:        uuid.Must(uuid.NewV4()),
					OktaEmail: "billy+existing@leo.org",
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					ID:        uuid.Must(uuid.NewV4()),
					FirstName: "Billy",
					LastName:  "Bob",
					Status:    &rejectedStatus,
					Email:     "billy+existing@leo.org",
				},
			},
			{
				Model:    user,
				LinkOnly: true,
			},
		}, []roles.RoleType{roles.RoleTypeTOO})

		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		queryBuilder := query.NewQueryBuilder()

		officeUserInfo := models.OfficeUser{
			LastName:               "Spaceman",
			FirstName:              "Billy",
			Email:                  officeUser.Email,
			TransportationOfficeID: transportationOffice.ID,
			Telephone:              "312-111-1111",
			TransportationOffice:   transportationOffice,
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			switch model.(type) {
			case *models.TransportationOffice:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(transportationOffice.ID))
			case *models.User:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(officeUser.User.ID))
				reflect.ValueOf(model).Elem().FieldByName("OktaID").Set(reflect.ValueOf(officeUser.User.OktaID))
				reflect.ValueOf(model).Elem().FieldByName("OktaEmail").Set(reflect.ValueOf(officeUser.User.OktaEmail))
			case *models.OfficeUser:
				reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(officeUser.ID))
				reflect.ValueOf(model).Elem().FieldByName("Email").Set(reflect.ValueOf(officeUser.Email))
				reflect.ValueOf(model).Elem().FieldByName("FirstName").Set(reflect.ValueOf(officeUser.FirstName))
				reflect.ValueOf(model).Elem().FieldByName("LastName").Set(reflect.ValueOf(officeUser.LastName))
				reflect.ValueOf(model).Elem().FieldByName("TransportationOfficeID").Set(reflect.ValueOf(officeUser.TransportationOfficeID))
				reflect.ValueOf(model).Elem().FieldByName("Telephone").Set(reflect.ValueOf(officeUser.Telephone))
				reflect.ValueOf(model).Elem().FieldByName("Status").Set(reflect.ValueOf(officeUser.Status))
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
		updatedOfficeUser, verrs, err := creator.CreateOfficeUser(appCtx, &officeUserInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(updatedOfficeUser)
		suite.Equal(updatedOfficeUser.ID, officeUser.ID)
		suite.Equal(updatedOfficeUser.Status, &requestedStatus)
		suite.Nil(updatedOfficeUser.RejectionReason)
	})

	// Bad transportation office ID
	suite.Run("If we are provided a transportation office that doesn't exist, the create should fail", func() {
		_, userInfo := setupTestData()
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})

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
		_, userInfo := setupTestData()
		transportationOffice := userInfo.TransportationOffice
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})

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
		_, userInfo := setupTestData()
		transportationOffice := userInfo.TransportationOffice
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})

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

	suite.Run("Test detailed uniqueness constraints being returned properly", func() {
		testCases := []struct {
			errorString              string
			shouldEdipiBeNil         bool
			shouldOtherUniqueIDBeNil bool
		}{
			{models.UniqueConstraintViolationOfficeUserEmailErrorString, false, false},
			{models.UniqueConstraintViolationOfficeUserEdipiErrorString, false, false},
			{models.UniqueConstraintViolationOfficeUserEdipiErrorString, true, false},
			{models.UniqueConstraintViolationOfficeUserOtherUniqueIDErrorString, false, false},
			{models.UniqueConstraintViolationOfficeUserOtherUniqueIDErrorString, false, true},
		}

		for _, tc := range testCases {
			_, userInfo := setupTestData()
			transportationOffice := userInfo.TransportationOffice
			if tc.shouldEdipiBeNil {
				userInfo.EDIPI = nil
			}
			if tc.shouldOtherUniqueIDBeNil {
				userInfo.OtherUniqueID = nil
			}
			appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})

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
					return nil, errors.New(tc.errorString)
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
			suite.EqualError(err, tc.errorString)

		}

	})

}
