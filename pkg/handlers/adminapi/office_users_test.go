package adminapi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/adminmessages"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

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

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		officeUserListFetcher := &mocks.ListFetcher{}
		officeUserListFetcher.On("FetchRecordList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		officeUserListFetcher.On("FetchRecordCount",
			mock.Anything,
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexOfficeUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: newQueryFilter,
			ListFetcher:    officeUserListFetcher,
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
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

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		officeUser := models.OfficeUser{ID: id}
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  req,
			OfficeUserID: strfmt.UUID(uuidString),
		}
		officeUserFetcher := &mocks.OfficeUserFetcher{}
		officeUserFetcher.On("FetchOfficeUser",
			mock.Anything,
		).Return(officeUser, nil).Once()
		handler := GetOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserFetcher,
			newQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetOfficeUserOK{}, response)
		okResponse := response.(*officeuserop.GetOfficeUserOK)
		suite.Equal(uuidString, okResponse.Payload.ID.String())
	})

	suite.T().Run("unsuccessful response when fetch fails", func(t *testing.T) {
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  req,
			OfficeUserID: strfmt.UUID(uuidString),
		}
		expectedError := models.ErrFetchNotFound
		officeUserFetcher := &mocks.OfficeUserFetcher{}
		officeUserFetcher.On("FetchOfficeUser",
			mock.Anything,
		).Return(models.OfficeUser{}, expectedError).Once()
		handler := GetOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserFetcher,
			newQueryFilter,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestCreateOfficeUserHandler() {
	transportationOfficeID, _ := uuid.NewV4()
	officeUserID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	userID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	officeUser := models.OfficeUser{
		ID:                     officeUserID,
		TransportationOfficeID: transportationOfficeID,
		UserID:                 &userID,
		Active:                 true,
	}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	req := httptest.NewRequest("POST", "/office_users", nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	tooRoleName := "Transportation Ordering Officer"
	tooRoleType := string(roles.RoleTypeTOO)

	tioRoleName := "Transportation Invoicing Officer"
	tioRoleType := string(roles.RoleTypeTIO)

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
			TransportationOfficeID: strfmt.UUID(officeUser.TransportationOfficeID.String()),
		},
	}

	suite.T().Run("Successful create", func(t *testing.T) {
		officeUserCreator := &mocks.OfficeUserCreator{}
		userRolesAssociator := &mocks.UserRoleAssociator{}

		officeUserCreator.On("CreateOfficeUser",
			mock.Anything,
			mock.Anything).Return(&officeUser, nil, nil).Once()

		userRolesAssociator.On("UpdateUserRoles",
			mock.Anything,
			mock.Anything).Return(nil, nil).Once()

		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserCreator,
			newQueryFilter,
			usersroles.NewUsersRolesCreator(suite.DB()),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserCreated{}, response)
	})

	suite.T().Run("Failed create", func(t *testing.T) {
		officeUserCreator := &mocks.OfficeUserCreator{}

		officeUserCreator.On("CreateOfficeUser",
			mock.Anything,
			mock.Anything).Return(&officeUser, nil, errors.New("Could not save user")).Once()

		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserCreator,
			newQueryFilter,
			usersroles.NewUsersRolesCreator(suite.DB()),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateOfficeUserHandler() {
	queryBuilder := query.NewQueryBuilder(suite.DB())
	updater := officeuser.NewOfficeUserUpdater(queryBuilder)
	handler := UpdateOfficeUserHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		updater,
		query.NewQueryFilter,
		usersroles.NewUsersRolesCreator(suite.DB()),
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
		transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())
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

	suite.T().Run("Update failes due to bad Transportation Office", func(t *testing.T) {
		officeUserUpdates := &adminmessages.OfficeUserUpdatePayload{
			TransportationOfficeID: strfmt.UUID("00000000-0000-0000-0000-000000000001"),
		}

		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  request,
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			OfficeUser:   officeUserUpdates,
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserInternalServerError{}, response)
	})
}
