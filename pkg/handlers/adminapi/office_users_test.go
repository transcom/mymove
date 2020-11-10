package adminapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gobuffalo/validate/v3"

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
	officeUser := models.OfficeUser{
		ID:                     officeUserID,
		TransportationOfficeID: transportationOfficeID,
		UserID:                 nil,
		Active:                 true,
	}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	req := httptest.NewRequest("POST", "/office_users", nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := officeuserop.CreateOfficeUserParams{
		HTTPRequest: req,
		OfficeUser: &adminmessages.OfficeUserCreatePayload{
			FirstName:              officeUser.FirstName,
			LastName:               officeUser.LastName,
			Telephone:              officeUser.Telephone,
			TransportationOfficeID: strfmt.UUID(officeUser.TransportationOfficeID.String()),
		},
	}

	suite.T().Run("Successful create", func(t *testing.T) {
		officeUserCreator := &mocks.OfficeUserCreator{}

		officeUserCreator.On("CreateOfficeUser",
			&officeUser,
			mock.Anything).Return(&officeUser, nil, nil).Once()

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
			&officeUser,
			mock.Anything).Return(&officeUser, nil, nil).Once()

		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserCreator,
			newQueryFilter,
			usersroles.NewUsersRolesCreator(suite.DB()),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserCreated{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateOfficeUserHandler() {
	officeUserID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	officeUser := models.OfficeUser{ID: officeUserID, FirstName: "Leo", LastName: "Spaceman", Telephone: "206-555-0199"}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	endpoint := fmt.Sprintf("/office_users/%s", officeUserID)
	req := httptest.NewRequest("PUT", endpoint, nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := officeuserop.UpdateOfficeUserParams{
		HTTPRequest: req,
		OfficeUser: &adminmessages.OfficeUserUpdatePayload{
			FirstName:      &officeUser.FirstName,
			MiddleInitials: officeUser.MiddleInitials,
			LastName:       &officeUser.LastName,
			Telephone:      &officeUser.Telephone,
		},
	}
	suite.T().Run("If the user is updated successfully it should be returned", func(t *testing.T) {
		//TODO remove this test and replace w/ test for just update logic (i.e. don't include with handler)
		ou := testdatagen.MakeDefaultOfficeUser(suite.DB())
		officeUserUpdater := &mocks.OfficeUserUpdater{}
		updatedParams := officeuserop.UpdateOfficeUserParams{
			HTTPRequest: req,
			OfficeUser: &adminmessages.OfficeUserUpdatePayload{
				FirstName:      &ou.FirstName,
				MiddleInitials: ou.MiddleInitials,
				LastName:       &ou.LastName,
				Telephone:      &ou.Telephone,
			},
		}
		updatedParams.OfficeUserID = strfmt.UUID(ou.ID.String())

		id1, _ := uuid.NewV4()
		role1 := roles.Role{
			ID:       id1,
			RoleType: "role1",
		}
		id2, _ := uuid.NewV4()
		role2 := roles.Role{
			ID:       id2,
			RoleType: "role2",
		}
		rs := roles.Roles{role1, role2}
		err := suite.DB().Create(rs)
		suite.NoError(err)
		payloadRoleName1 := "name1"
		payloadRoleType1 := string(role1.RoleType)
		payloadRole1 := &adminmessages.OfficeUserRole{
			Name:     &payloadRoleName1,
			RoleType: &payloadRoleType1,
		}
		payloadRoleName2 := "name2"
		payloadRoleType2 := string(role2.RoleType)
		payloadRole2 := &adminmessages.OfficeUserRole{
			Name:     &payloadRoleName2,
			RoleType: &payloadRoleType2,
		}

		payloadRoles := []*adminmessages.OfficeUserRole{payloadRole1, payloadRole2}
		updatedParams.OfficeUser.Roles = payloadRoles

		officeUserUpdater.On("UpdateOfficeUser",
			mock.Anything,
			updatedParams.OfficeUser,
		).Return(&ou, nil, nil).Once()

		handler := UpdateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserUpdater,
			newQueryFilter,
			usersroles.NewUsersRolesCreator(suite.DB()),
		}

		response := handler.Handle(updatedParams)
		suite.IsType(&officeuserop.UpdateOfficeUserOK{}, response)

		ur := models.UsersRoles{}
		n, err := suite.DB().Count(&ur)
		suite.NoError(err)
		suite.Equal(2, n)

		user := models.User{}
		err = suite.DB().Eager("Roles").Find(&user, ou.UserID)
		suite.NoError(err)
		suite.Require().Len(user.Roles, 2)
		suite.Equal(user.Roles[0].ID, role1.ID)
		suite.Equal(user.Roles[1].ID, role2.ID)
	})

	suite.T().Run("Successful update", func(t *testing.T) {
		officeUserUpdater := &mocks.OfficeUserUpdater{}

		officeUserUpdater.On("UpdateOfficeUser",
			mock.Anything,
			params.OfficeUser,
		).Return(&officeUser, nil, nil).Once()

		handler := UpdateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserUpdater,
			newQueryFilter,
			usersroles.NewUsersRolesCreator(suite.DB()),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserOK{}, response)
	})

	suite.T().Run("Failed update", func(t *testing.T) {
		officeUserUpdater := &mocks.OfficeUserUpdater{}

		officeUserUpdater.On("UpdateOfficeUser",
			mock.Anything,
			params.OfficeUser,
		).Return(&officeUser, nil, nil).Once()

		handler := UpdateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserUpdater,
			newQueryFilter,
			usersroles.NewUsersRolesCreator(suite.DB()),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserOK{}, response)
	})

	officeUserUpdater := &mocks.OfficeUserUpdater{}
	err := validate.NewErrors()

	officeUserUpdater.On("UpdateOfficeUser",
		mock.Anything,
		params.OfficeUser,
	).Return(nil, err, nil).Once()

	handler := UpdateOfficeUserHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		officeUserUpdater,
		newQueryFilter,
		usersroles.NewUsersRolesCreator(suite.DB()),
	}

	handler.Handle(params)
	suite.Error(err, "Error saving user")

}
