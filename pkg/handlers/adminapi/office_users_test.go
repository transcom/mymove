package adminapi

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/adminmessages"

	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexOfficeUsersHandler() {
	setupTestData := func() models.OfficeUsers {
		return models.OfficeUsers{
			testdatagen.MakeDefaultOfficeUser(suite.DB()),
			testdatagen.MakeDefaultOfficeUser(suite.DB()),
		}
	}

	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		officeUsers := setupTestData()
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficeUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(officeUsers[0].ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("fetch return an empty list", func() {
		setupTestData()
		// TEST:				IndexOfficeUserHandler, Fetcher
		// Set up:				Provide an invalid search that won't be found
		// Expected Outcome:	An empty list is returned and we get a 200 OK.
		fakeFilter := "{\"search\":\"something\"}"

		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &fakeFilter,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficeUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewQueryFilter: query.NewQueryFilter,
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)

		suite.Len(okResponse.Payload, 0)
	})
}

func (suite *HandlerSuite) TestGetOfficeUserHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		// Test:				GetOfficeUserHandler, Fetcher
		// Set up:				Provide a valid req with the office user ID to the endpoint
		// Expected Outcome:	The office user is returned and we get a 200 OK.
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			officeuser.NewOfficeUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetOfficeUserOK{}, response)
		okResponse := response.(*officeuserop.GetOfficeUserOK)
		suite.Equal(officeUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("500 error - Internal Server error. Unsuccessful fetch ", func() {
		// Test:				GetOfficeUserHandler, Fetcher
		// Set up:				Provide a valid req with the fake office user ID to the endpoint
		// Expected Outcome:	The office user is returned and we get a 404 NotFound.
		fakeID := "3b9c2975-4e54-40ea-a781-bab7d6e4a502"
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/office_users/%s", fakeID)),
			OfficeUserID: strfmt.UUID(fakeID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			officeuser.NewOfficeUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&handlers.ErrResponse{}, response)
		errResponse := response.(*handlers.ErrResponse)
		suite.Equal(http.StatusInternalServerError, errResponse.Code)
	})
}

func (suite *HandlerSuite) TestCreateOfficeUserHandler() {
	tooRoleName := "Transportation Ordering Officer"
	tooRoleType := string(roles.RoleTypeTOO)

	tioRoleName := "Transportation Invoicing Officer"
	tioRoleType := string(roles.RoleTypeTIO)

	suite.Run("200 - Successfully create Office User", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Create a new Office User, save new user to the DB
		// Expected Outcome:	The office user is created and we get a 200 OK.
		transportationOfficeID := testdatagen.MakeDefaultTransportationOffice(suite.DB()).ID

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreatePayload{
				FirstName: "Sam",
				LastName:  "Cook",
				Telephone: "555-555-5555",
				Email:     "fakeemail@gmail.com",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
					{
						Name:     &tioRoleName,
						RoleType: &tioRoleType,
					},
				},
				TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
			},
		}
		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserCreated{}, response)
	})

	suite.Run("Failed create", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Add new Office User to the DB
		// Expected Outcome:	The office user is not created and we get a 500 internal server error.
		fakeTransportationOfficeID := "3b9c2975-4e54-40ea-a781-bab7d6e4a502"
		officeUser := testdatagen.MakeStubbedOfficeUser(suite.DB())

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreatePayload{
				FirstName: officeUser.FirstName,
				LastName:  officeUser.LastName,
				Telephone: officeUser.Telephone,
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
					{
						Name:     &tioRoleName,
						RoleType: &tioRoleType,
					},
				},
				TransportationOfficeID: strfmt.UUID(fakeTransportationOfficeID),
			},
		}

		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateOfficeUserHandler() {
	setupHandler := func(updater services.OfficeUserUpdater, revoker services.UserSessionRevocation) UpdateOfficeUserHandler {
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.Logger())
		sessionManagers := setupSessionManagers()
		handlerContext.SetSessionManagers(sessionManagers)
		return UpdateOfficeUserHandler{
			handlerContext,
			updater,
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(), // a special can of worms, TODO mocked tests
			revoker,
		}
	}

	setupTestData := func() models.OfficeUser {
		return testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
			OfficeUser: models.OfficeUser{
				TransportationOffice: models.TransportationOffice{
					Name: "Random Office",
				},
			},
		})
	}

	suite.Run("Office user is successfully updated", func() {
		officeUser := setupTestData()
		transportationOffice := testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{Stub: true})
		firstName := "Riley"
		middleInitials := "RB"
		telephone := "865-555-5309"

		officeUserUpdates := &adminmessages.OfficeUserUpdatePayload{
			FirstName:              &firstName,
			MiddleInitials:         &middleInitials,
			Telephone:              &telephone,
			TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
		}

		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			OfficeUser:   officeUserUpdates,
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		// Mock DB update:
		expectedInput := *officeUserUpdates // make a copy so we can ensure our expected values don't change
		expectedOfficeUser := officeUser
		expectedOfficeUser.FirstName = *expectedInput.FirstName
		expectedOfficeUser.MiddleInitials = expectedInput.MiddleInitials
		expectedOfficeUser.Telephone = *expectedInput.Telephone
		expectedOfficeUser.TransportationOfficeID = transportationOffice.ID

		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.On("UpdateOfficeUser", mock.AnythingOfType("*appcontext.appContext"), officeUser.ID, &expectedInput).Return(&expectedOfficeUser, nil, nil)

		response := setupHandler(&mockUpdater, &mocks.UserSessionRevocation{}).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserOK{}, response)

		okResponse := response.(*officeuserop.UpdateOfficeUserOK)
		// Check updates:
		suite.Equal(firstName, *okResponse.Payload.FirstName)
		suite.Equal(middleInitials, *okResponse.Payload.MiddleInitials)
		suite.Equal(telephone, *okResponse.Payload.Telephone)
		suite.Equal(transportationOffice.ID.String(), okResponse.Payload.TransportationOfficeID.String())
		suite.Equal(officeUser.LastName, *okResponse.Payload.LastName) // should not have been updated
		suite.Equal(officeUser.Email, *okResponse.Payload.Email)       // should not have been updated
	})

	suite.Run("Update fails due to bad Transportation Office", func() {
		officeUser := setupTestData()
		officeUserUpdates := &adminmessages.OfficeUserUpdatePayload{
			TransportationOfficeID: strfmt.UUID("00000000-0000-0000-0000-000000000001"),
		}

		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			OfficeUser:   officeUserUpdates,
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		expectedInput := *officeUserUpdates
		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.On("UpdateOfficeUser", mock.AnythingOfType("*appcontext.appContext"), officeUser.ID, &expectedInput).Return(nil, nil, sql.ErrNoRows)

		response := setupHandler(&mockUpdater, &mocks.UserSessionRevocation{}).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserInternalServerError{}, response)
	})

	suite.Run("Office user session is revoked when roles are changed", func() {
		officeUser := setupTestData()
		// Setup payload to remove all roles for office user
		newRoles := []*adminmessages.OfficeUserRole{}
		officeUserUpdates := &adminmessages.OfficeUserUpdatePayload{
			Roles: newRoles,
		}
		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			OfficeUser:   officeUserUpdates,
		}

		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.
			On("UpdateOfficeUser", mock.AnythingOfType("*appcontext.appContext"), officeUser.ID, officeUserUpdates).
			Return(&officeUser, nil, nil)

		expectedSessionUpdate := &adminmessages.UserUpdatePayload{
			RevokeOfficeSession: swag.Bool(true),
		}
		mockRevoker := mocks.UserSessionRevocation{}
		mockRevoker.
			On("RevokeUserSession", mock.AnythingOfType("*appcontext.appContext"), *officeUser.UserID, expectedSessionUpdate, mock.Anything).
			Return(nil, nil, nil).
			Once()

		response := setupHandler(&mockUpdater, &mockRevoker).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserOK{}, response)
		mockRevoker.AssertNumberOfCalls(suite.T(), "RevokeUserSession", 1)
	})
}
