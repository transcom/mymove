package adminapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/adminmessages"

	"github.com/gofrs/uuid"

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
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID: id,
		},
	}
	// The commands MakeOfficeUser and MakeDefaultOfficeUser add a new Office User to the DB.
	// Don't use if writing a failing test for a User that should not be found.
	testdatagen.MakeOfficeUser(suite.DB(), assertions)
	testdatagen.MakeDefaultOfficeUser(suite.DB())

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/office_users", nil)
	req = suite.AuthenticateAdminRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: req,
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexOfficeUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("fetch return an empty list", func(t *testing.T) {
		// TEST:				IndexOfficeUserHandler, Fetcher
		// Set up:				Provide an invalid search that won't be found
		// Expected Outcome:	An empty list is returned and we get a 200 OK.
		fakeFilter := "{\"search\":\"something\"}"

		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: req,
			Filter:      &fakeFilter,
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexOfficeUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID: id,
		},
	}
	testdatagen.MakeOfficeUser(suite.DB(), assertions)

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/office_users/%s", id), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  req,
			OfficeUserID: strfmt.UUID(uuidString),
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := GetOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeuser.NewOfficeUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetOfficeUserOK{}, response)
		okResponse := response.(*officeuserop.GetOfficeUserOK)
		suite.Equal(uuidString, okResponse.Payload.ID.String())
	})

	suite.T().Run("successful response", func(t *testing.T) {
		// Test:				GetOfficeUserHandler, Fetcher
		// Set up:				Provide a valid req with the office user ID to the endpoint
		// Expected Outcome:	The office user is returned and we get a 200 OK.

		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  req,
			OfficeUserID: strfmt.UUID(uuidString),
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := GetOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeuser.NewOfficeUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetOfficeUserOK{}, response)
		okResponse := response.(*officeuserop.GetOfficeUserOK)
		suite.Equal(uuidString, okResponse.Payload.ID.String())
	})

	suite.T().Run("500 error - Internal Server error. Unsuccessful fetch ", func(t *testing.T) {
		// Test:				GetOfficeUserHandler, Fetcher
		// Set up:				Provide a valid req with the fake office user ID to the endpoint
		// Expected Outcome:	The office user is returned and we get a 404 NotFound.
		fakeID := "3b9c2975-4e54-40ea-a781-bab7d6e4a502"
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  req,
			OfficeUserID: strfmt.UUID(fakeID),
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := GetOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
	transportationOfficeID := testdatagen.MakeDefaultTransportationOffice(suite.DB()).ID
	officeUser := models.OfficeUser{
		TransportationOfficeID: transportationOfficeID,
		Active:                 true,
	}

	req := httptest.NewRequest("POST", "/office_users", nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	tooRoleName := "Transportation Ordering Officer"
	tooRoleType := string(roles.RoleTypeTOO)

	tioRoleName := "Transportation Invoicing Officer"
	tioRoleType := string(roles.RoleTypeTIO)

	suite.T().Run("200 - Successfully create Office User", func(t *testing.T) {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Create a new Office User, save new user to the DB
		// Expected Outcome:	The office user is created and we get a 200 OK.

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: req,
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
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeuser.NewOfficeUserCreator(suite.DB(), queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(suite.DB()),
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserCreated{}, response)
	})

	suite.T().Run("Failed create", func(t *testing.T) {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Add new Office User to the DB
		// Expected Outcome:	The office user is not created and we get a 500 internal server error.
		fakeTransportationOfficeID := "3b9c2975-4e54-40ea-a781-bab7d6e4a502"
		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: req,
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

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeuser.NewOfficeUserCreator(suite.DB(), queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(suite.DB()),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateOfficeUserHandler() {
	mockUpdater := mocks.OfficeUserUpdater{}
	handler := UpdateOfficeUserHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		&mockUpdater,
		query.NewQueryFilter,
		usersroles.NewUsersRolesCreator(suite.DB()), // a special can of worms, TODO mocked tests
	}

	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			TransportationOffice: models.TransportationOffice{
				Name: "Random Office",
			},
		},
	})
	requestUser := testdatagen.MakeStubbedUser(suite.DB())

	endpoint := fmt.Sprintf("/office_users/%s", officeUser.ID)
	request := suite.AuthenticateUserRequest(httptest.NewRequest("PUT", endpoint, nil), requestUser)

	suite.T().Run("Office user is successfully updated", func(t *testing.T) {
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
			HTTPRequest:  request,
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

		mockUpdater.On("UpdateOfficeUser", officeUser.ID, &expectedInput).Return(&expectedOfficeUser, nil, nil)

		response := handler.Handle(params)
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

	suite.T().Run("Update fails due to bad Transportation Office", func(t *testing.T) {
		officeUserUpdates := &adminmessages.OfficeUserUpdatePayload{
			TransportationOfficeID: strfmt.UUID("00000000-0000-0000-0000-000000000001"),
		}

		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  request,
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			OfficeUser:   officeUserUpdates,
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		expectedInput := *officeUserUpdates
		mockUpdater.On("UpdateOfficeUser", officeUser.ID, &expectedInput).Return(nil, nil, sql.ErrNoRows)

		response := handler.Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserInternalServerError{}, response)
	})
}
